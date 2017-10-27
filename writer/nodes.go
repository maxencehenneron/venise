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

func NewNodesWriter(cache *cache.OSM, database *database.MongoDatabase) *NodesWriter {
	return &NodesWriter{
		OSM{cache, database},
	}
}

func (nw *NodesWriter) WriteNodes(tags map[string][]string) {
	nodes := nw.cache.Nodes.Iterate()

	for node := range nodes {
		if ShouldInsertNode(tags, *node) {
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

// Verifies that the way is in the list of wanted nodes
func ShouldInsertNode(tags map[string][]string, way structures.Node) bool {
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
