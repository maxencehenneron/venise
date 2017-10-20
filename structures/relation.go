package structures

type Relation struct {
	ID      int64 `bson:"_id"`
	Tags    map[string]string
	Members []Member
	Info    Info
	Loc     *struct {
		Type        string      `bson:"type"`
		Coordinates [][]float64 `bson:"coordinates"`
	} `bson:"loc,omitempty"`
}

func (r *Relation) FeedCoordinates(members []interface{}) {

}
