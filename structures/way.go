package structures

import (
	"errors"

	"github.com/dernise/venise/geo"
)

type Way struct {
	ID      int64 `bson:"_id"`
	Tags    map[string]string
	NodeIDs []int64
	Info    Info

	Center *Point      `bson:"center,omitempty"`
	Loc    interface{} `bson:"loc,omitempty"`
}

func (w *Way) FeedCoordinates(nodes []Node) error {
	if len(w.NodeIDs) == 0 {
		return errors.New("one of the ways has no node")
	}

	closed := w.IsClosed()

	if closed {
		coordinates := make([][]geo.Point, 1)
		polygonCoordinates := make([][]geo.Point, 1)

		// Get all the coordinates
		for _, node := range nodes {
			coordinates[0] = append(coordinates[0], geo.Point{node.Loc.Coordinates[0], node.Loc.Coordinates[1]})
			polygonCoordinates[0] = append(polygonCoordinates[0], geo.Point{node.Loc.Coordinates[0], node.Loc.Coordinates[1]})
		}

		// Create the polygon
		p := &geo.Polygon{
			Rings: polygonCoordinates,
		}

		// Transform the polygon
		p.Transform(geo.Mercator.Project)

		// Calculate the centre
		center := geo.Polylabel(*p, 1, true)

		// Transform the center back
		center.Transform(geo.Mercator.Inverse)

		w.Center = &Point{
			"Point",
			center,
		}

		w.Loc = Polygon{"Polygon", coordinates}
	} else {
		coordinates := make([]geo.Point, len(nodes))
		for _, node := range nodes {
			coordinates = append(coordinates, geo.Point{node.Loc.Coordinates[0], node.Loc.Coordinates[1]})
		}

		w.Loc = LineString{"LineString", coordinates}
	}

	return nil
}

func (w *Way) IsClosed() bool {
	return len(w.NodeIDs) >= 4 && w.NodeIDs[0] == w.NodeIDs[len(w.NodeIDs)-1]
}
