package database

import (
	"fmt"

	"reflect"

	"log"

	"github.com/dernise/venise/geo"
	"github.com/dernise/venise/structures"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	BULK_SIZE = 5000
)

var collectionTypeMap = map[string]string{
	"structures.Node":     "nodes",
	"structures.Way":      "ways",
	"structures.Relation": "relations",
}

type MongoDatabase struct {
	*mgo.Database
	bulk            *mgo.Bulk
	currentSession  *mgo.Session
	currentBulkSize int
	bulkCollection  string
}

func New(mongoDatabase *mgo.Database) (*MongoDatabase, error) {
	database := &MongoDatabase{
		mongoDatabase,
		nil,
		nil,
		0,
		"",
	}

	if err := database.EnsureIndexes(); err != nil {
		return nil, err
	}

	return database, nil
}

func (mg *MongoDatabase) EnsureIndexes() error {
	session := mg.Session.Copy()
	defer session.Close()

	collectionIndexes := make(map[*mgo.Collection][]mgo.Index)

	// Node indexes
	nodes := mg.C(NodeCollection)
	collectionIndexes[nodes] = []mgo.Index{
		{
			Key:  []string{"$2dsphere:loc"},
			Bits: 26,
		},
	}

	ways := mg.C(WayCollection)
	collectionIndexes[ways] = []mgo.Index{
		{
			Key:  []string{"$2dsphere:loc"},
			Bits: 26,
		},
		{
			Key:  []string{"$2dsphere:center"},
			Bits: 26,
		},
	}

	relations := mg.C(RelationCollection)
	collectionIndexes[relations] = []mgo.Index{
		{
			Key:  []string{"$2dsphere:loc"},
			Bits: 26,
		},
	}

	for collection, indexes := range collectionIndexes {
		for _, index := range indexes {
			err := collection.EnsureIndex(index)

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (mg *MongoDatabase) checkBulk(object interface{}) {
	objectType := reflect.TypeOf(object).String()
	objectCollection := collectionTypeMap[objectType]

	// First check : bulk isn't nil
	if mg.bulk == nil {
		mg.createNewBulk(objectType, objectCollection)
		return
	}

	// Second check : check if object didn't change
	if mg.bulkCollection != objectCollection {
		mg.InsertCurrentBulk()
		mg.createNewBulk(objectType, objectCollection)
		return
	}

	// Third check : Check if the limit isn't reached
	if mg.currentBulkSize >= BULK_SIZE {
		mg.InsertCurrentBulk()
		mg.createNewBulk(objectType, objectCollection)
		return
	}
}

func (mg *MongoDatabase) InsertCurrentBulk() {
	defer mg.currentSession.Close()

	_, bulkErr := mg.bulk.Run()
	if bulkErr != nil {
		fmt.Println(bulkErr)
	}

	mg.bulk = nil
}

func (mg *MongoDatabase) createNewBulk(objectType string, objectCollection string) {
	session := mg.Session.Clone()

	mg.bulkCollection = objectCollection
	mg.bulk = mg.C(objectCollection).With(session).Bulk()
	mg.currentBulkSize = 0
	mg.currentSession = session
	mg.bulk.Unordered()
}

func (mg *MongoDatabase) insertObjectInBulk(object interface{}) {
	mg.currentBulkSize++
	mg.bulk.Insert(object)
}

func (mg *MongoDatabase) Insert(object interface{}) {
	mg.checkBulk(object)

	switch o := object.(type) {
	case *structures.Node:
		o.Loc = &structures.Point{
			Type:        "Point",
			Coordinates: geo.Point{o.Lon, o.Lat},
		}

		mg.insertObjectInBulk(o)
		break
	case *structures.Way:
		nodes := make([]structures.Node, len(o.NodeIDs))

		// Gets all the nodes associated to the way
		searchNodes, err := mg.FindNodes(o.NodeIDs)
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		//	Creates an array of node from their ID and the search result
		for index, nodeId := range o.NodeIDs {
			node := structures.GetNodeById(nodeId, searchNodes)
			if node == nil {
				log.Fatal("node wasn't found in the search result")
			}

			nodes[index] = *node
		}

		// Adds the center and locations
		err = o.FeedCoordinates(nodes)
		if err != nil {
			log.Fatal(err.Error())
		}

		mg.insertObjectInBulk(o)
		break
	default:

	}
}

func (mg *MongoDatabase) FindNodes(ids []int64) ([]structures.Node, error) {
	session := mg.Session.Clone()
	defer session.Close()

	var nodes []structures.Node
	err := mg.C(NodeCollection).With(session).Find(bson.M{"_id": bson.M{"$in": ids}}).All(&nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}
