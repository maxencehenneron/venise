package structures

type Object struct {
	OsmId int64             `bson:"osm_id"`
	Tags  map[string]string `bson:"tags"`

	//Location
	Loc interface{} `bson:"loc"`
}
