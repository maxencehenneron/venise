package structures

type Node struct {
	ID   int64 `bson:"_id"`
	Lat  float64
	Lon  float64
	Tags map[string]string
	Info Info

	//BSON Data (Needs to be set before inserting the object)
	Loc *Point `bson:"loc"`
}

func GetNodeById(id int64, nodes []Node) *Node {
	for _, node := range nodes {
		if node.ID == id {
			return &node
		}
	}
	return nil
}
