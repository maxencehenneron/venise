package structures

type Member struct {
	ID   int64 `bson:"_id"`
	Type MemberType
	Role string
	Way  *Way  `json:"-"`
	Node *Node `json:"-"`
}

type MemberType int

const (
	NodeType MemberType = iota
	WayType
	RelationType
)
