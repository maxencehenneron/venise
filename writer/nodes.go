package writer

import (
	"github.com/dernise/venise/cache"
	"github.com/dernise/venise/database"
	"github.com/dernise/venise/geo"
	"github.com/dernise/venise/structures"
)

type NodesWriter struct {
	OSM
}

func NewNodesWriter(cache *cache.OSM, database *database.MongoDatabase, tags map[string][]string) *NodesWriter {
	return &NodesWriter{
		OSM{cache, database, tags},
	}
}

func (nw *NodesWriter) WriteNodes() {
	nodes := nw.cache.Nodes.Iterate()

	for node := range nodes {
		if nw.HasInterrestingTag(node.Tags) {
			object := structures.Object{
				OsmId: node.ID,
				Tags:  node.Tags,
				Loc: structures.Point{
					Type:        "Point",
					Coordinates: *geo.NewPointFromLatLng(node.Lat, node.Lon),
				},
			}
			nw.database.Insert(node.Tags["amenity"], object)
		}
	}
}
