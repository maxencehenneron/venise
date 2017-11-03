package writer

import (
	"github.com/dernise/venise/cache"
	"github.com/dernise/venise/database"
)

type WaysWriter struct {
	OSM
}

func NewWaysWriter(cache *cache.OSM, database *database.MongoDatabase, tags map[string][]string) *WaysWriter {
	return &WaysWriter{
		OSM{cache, database, tags},
	}
}

func (nw *WaysWriter) WriteWays() {
	ways := nw.cache.Ways.Iterate()

	for way := range ways {
		if nw.HasInterrestingTag(way.Tags) {
			nw.cache.Coords.FillWay(way)
		}
	}
}
