package structures

import "github.com/dernise/venise/geo"

type Point struct {
	Type        string    `bson:"type"`
	Coordinates geo.Point `bson:"coordinates"`
}

type LineString struct {
	Type        string      `bson:"type"`
	Coordinates []geo.Point `bson:"coordinates"`
}

type Polygon struct {
	Type        string        `bson:"type"`
	Coordinates [][]geo.Point `bson:"coordinates"`
}
