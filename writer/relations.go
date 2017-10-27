package writer

import (
	"fmt"

	"github.com/dernise/venise/cache"
	"github.com/dernise/venise/database"
	"github.com/dernise/venise/structures"
)

type RelationsWriter struct {
	OSM
}

func NewRelationsWriter(cache *cache.OSM, database *database.MongoDatabase) *RelationsWriter {
	return &RelationsWriter{
		OSM{cache, database},
	}
}

func (nw *RelationsWriter) WriteRelations(tags map[string][]string) {
	relations := nw.cache.Relations.Iterate()

	for relation := range relations {
		if ShouldInsertRelation(tags, *relation) {
			fmt.Printf("relation : %+v\n", relation.Tags["amenity"])
		}
	}
}

// Verifies that the way is in the list of wanted relations
func ShouldInsertRelation(tags map[string][]string, way structures.Relation) bool {
	shouldInsert := false
	for key, values := range tags {
		if val, ok := way.Tags[key]; ok {
			fmt.Printf("%v LOOOL\n", val)
			for _, value := range values {
				if val == value {
					shouldInsert = true
				}
			}
		}
	}
	return shouldInsert
}
