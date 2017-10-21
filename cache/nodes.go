package cache

import (
	"github.com/dernise/venise/cache/binary"
	"github.com/dernise/venise/structures"
	"github.com/syndtr/goleveldb/leveldb"
)

type Nodes struct {
	*leveldb.DB
}

func NewNodesCache(path string) (*Nodes, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	cache := Nodes{db}

	return &cache, err
}

func (n *Nodes) PutNode(node structures.Node) error {
	if len(node.Tags) == 0 {
		return nil
	}

	bytes, err := binary.MarshalNode(&node)
	if err != nil {
		return err
	}

	n.Put(idToKeyBuf(node.ID), bytes, nil)
	return nil
}
