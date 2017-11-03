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

func NewRelationsWriter(cache *cache.OSM, database *database.MongoDatabase, tags map[string][]string) *RelationsWriter {
	return &RelationsWriter{
		OSM{cache, database, tags},
	}
}

func (rw *RelationsWriter) WriteRelations() {
	relations := rw.cache.Relations.Iterate()

	for relation := range relations {
		if rw.HasInterrestingTag(relation.Tags) {
			rw.transformToGeoJson(relation)
		}
	}
}

func (rw *RelationsWriter) transformToGeoJson(relation *structures.Relation) interface{} {
	// Relation is a route
	if len(relation.Members) == 0 {
		fmt.Printf("Weird relation: %v and no members\n", relation.ID)
		return nil
	}

	if relation.Tags["type"] == "route" || relation.Tags["type"] == "waterway" {
		return nil // Not needed at this time
	}

	if relation.Tags["type"] == "multipolygon" || relation.Tags["type"] == "boundary" {
		outerCount := 0
		for _, member := range relation.Members {
			if member.Role == "outer" {
				outerCount++
			} else if member.Role != "inner" {
				fmt.Printf("Ignored member because it has an invalid role : %v\n", member.Role)
			}
		}

		if outerCount > 1 {
			fmt.Println(relation.ID)
		}

		// Checks if the polygon is simple
		simpleMp := false
		_, hasType := relation.Tags["type"]
		if outerCount == 1 && !hasType {
			simpleMp = true
		}

		if simpleMp {
			outerWays, err := getOuterWays(relation.Members, rw.cache)
			if err != nil {
				panic(err)
			}

			way := outerWays[0] // Has only one way
			fmt.Printf("%v\n", way.ID)
		}

	}

	return nil
}

func getOuterWays(members []structures.Member, cache *cache.OSM) ([]*structures.Way, error) {
	var outerWays []*structures.Way
	for _, member := range members {
		if member.Role == "outer" {
			way, err := cache.Ways.GetWay(member.ID)
			if err != nil {
				return nil, err
			}
			outerWays = append(outerWays, way)
		}
	}
	return outerWays, nil
}
