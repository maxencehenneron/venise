package database

import (
	"github.com/dernise/venise/structures"
	"gopkg.in/mgo.v2"
)

type MongoDatabase struct {
	*mgo.Database
	collections []string
}

func New(mongoDatabase *mgo.Database, collections []string) (*MongoDatabase, error) {
	database := &MongoDatabase{
		mongoDatabase,
		collections,
	}

	if err := database.EnsureIndexes(); err != nil {
		return nil, err
	}

	return database, nil
}

func (mg *MongoDatabase) EnsureIndexes() error {
	session := mg.Session.Copy()
	defer session.Close()

	for _, c := range mg.collections {
		collection := mg.C(c)

		indexes := []mgo.Index{
			{
				Key:  []string{"$2dsphere:loc"},
				Bits: 26,
			},
			{
				Key:    []string{"osm_id", "loc.type"},
				Unique: true,
			},
		}

		for _, index := range indexes {
			err := collection.With(session).EnsureIndex(index)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (mg *MongoDatabase) Insert(collectionName string, object structures.Object) error {
	session := mg.Session.Copy()
	defer session.Close()

	collection := mg.C(collectionName)

	err := collection.With(session).Insert(object)
	if err != nil {
		return err
	}

	return nil
}
