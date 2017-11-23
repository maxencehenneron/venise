package cache

import (
	"github.com/dernise/venise/cache/binary"
	"github.com/dernise/venise/structures"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type Nodes struct {
	*leveldb.DB
}

func NewNodesCache(path string, options *opt.Options) (*Nodes, error) {
	db, err := leveldb.OpenFile(path, options)
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

func (n *Nodes) Iterate() chan *structures.Node {
	nodes := make(chan *structures.Node)
	go func() {
		defer close(nodes)
		iterator := n.NewIterator(nil, nil)
		iterator.First()
		for ; iterator.Valid(); iterator.Next() {
			node, err := binary.UnmarshalNode(iterator.Value())
			if err != nil {
				panic(err)
			}
			node.ID = idFromKeyBuf(iterator.Key())

			nodes <- node
		}
	}()
	return nodes
}
