package writer

import (
	"fmt"

	"github.com/dernise/venise/cache"
	"github.com/dernise/venise/structures"
	"github.com/dernise/venise/database"
)

type WaysWriter struct {
	OSM
}

func NewWaysWriter(cache *cache.OSM, database *database.MongoDatabase) *WaysWriter {
	return &WaysWriter{
		OSM{cache, database},
	}
}

func (nw *WaysWriter) WriteWays(tags map[string][]string) {
	ways := nw.cache.Ways.Iterate()

	for way := range ways {
		if ShouldInsertWay(tags, *way) {
			fmt.Printf("way : %+v - %+v\n", way.Tags["amenity"], way.Tags["name"])
			nw.cache.Coords.FillWay(way)
		}
	}
}

// Verifies that the way is in the list of wanted nodes
func ShouldInsertWay(tags map[string][]string, way structures.Way) bool {
	shouldInsert := false
	for key, values := range tags {
		if val, ok := way.Tags[key]; ok {
			for _, value := range values {
				if val == value {
					shouldInsert = true
				}
			}
		}
	}
	return shouldInsert
}
