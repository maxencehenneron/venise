package binary

import (
	"github.com/dernise/imposm3/element"
	"github.com/dernise/venise/structures"
	"github.com/golang/protobuf/proto"
)

const COORD_FACTOR float64 = 11930464.7083 // ((2<<31)-1)/360.0

func CoordToInt(coord float64) uint32 {
	return uint32((coord + 180.0) * COORD_FACTOR)
}

func IntToCoord(coord uint32) float64 {
	return float64((float64(coord) / COORD_FACTOR) - 180.0)
}

func MarshalNode(node *structures.Node) ([]byte, error) {
	pbfNode := &Node{}
	pbfNode.fromWgsCoord(node.Lon, node.Lat)
	pbfNode.Tags = tagsAsArray(node.Tags)
	return proto.Marshal(pbfNode)
}

func UnmarshalNode(data []byte) (node *structures.Node, err error) {
	pbfNode := &Node{}
	err = proto.Unmarshal(data, pbfNode)
	if err != nil {
		return nil, err
	}

	node = &structures.Node{}
	node.Lon, node.Lat = pbfNode.wgsCoord()
	node.Tags = tagsFromArray(pbfNode.Tags)
	return node, nil
}

func deltaPack(data []int64) {
	if len(data) < 2 {
		return
	}
	lastVal := data[0]
	for i := 1; i < len(data); i++ {
		data[i], lastVal = data[i]-lastVal, data[i]
	}
}

func deltaUnpack(data []int64) {
	if len(data) < 2 {
		return
	}
	for i := 1; i < len(data); i++ {
		data[i] = data[i] + data[i-1]
	}
}

func MarshalWay(way *structures.Way) ([]byte, error) {
	pbfWay := &Way{}
	deltaPack(way.NodeIDs)
	pbfWay.Refs = way.NodeIDs
	pbfWay.Tags = tagsAsArray(way.Tags)
	return proto.Marshal(pbfWay)
}

func UnmarshalWay(data []byte) (way *element.Way, err error) {
	pbfWay := &Way{}
	err = proto.Unmarshal(data, pbfWay)
	if err != nil {
		return nil, err
	}

	way = &element.Way{}
	deltaUnpack(pbfWay.Refs)
	way.Refs = pbfWay.Refs
	way.Tags = tagsFromArray(pbfWay.Tags)
	return way, nil
}

func MarshalRelation(relation *structures.Relation) ([]byte, error) {
	pbfRelation := &Relation{}
	pbfRelation.MemberIds = make([]int64, len(relation.Members))
	pbfRelation.MemberTypes = make([]Relation_MemberType, len(relation.Members))
	pbfRelation.MemberRoles = make([]string, len(relation.Members))
	for i, m := range relation.Members {
		pbfRelation.MemberIds[i] = m.ID
		pbfRelation.MemberTypes[i] = Relation_MemberType(m.Type)
		pbfRelation.MemberRoles[i] = m.Role
	}
	pbfRelation.Tags = tagsAsArray(relation.Tags)
	return proto.Marshal(pbfRelation)
}

func UnmarshalRelation(data []byte) (relation *structures.Relation, err error) {
	pbfRelation := &Relation{}
	err = proto.Unmarshal(data, pbfRelation)
	if err != nil {
		return nil, err
	}

	relation = &structures.Relation{}
	relation.Members = make([]structures.Member, len(pbfRelation.MemberIds))
	for i, _ := range pbfRelation.MemberIds {
		relation.Members[i].ID = pbfRelation.MemberIds[i]
		relation.Members[i].Type = structures.MemberType(pbfRelation.MemberTypes[i])
		relation.Members[i].Role = pbfRelation.MemberRoles[i]
	}
	//relation.Nodes = pbfRelation.Node
	relation.Tags = tagsFromArray(pbfRelation.Tags)
	return relation, nil
}
