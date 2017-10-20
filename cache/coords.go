package cache

import (
	"github.com/dernise/venise/cache/binary"
	"github.com/dernise/venise/structures"
	"github.com/syndtr/goleveldb/leveldb"
)

// The coord cache is a cache of every nodes' coordinates.
type Coords struct {
	*leveldb.DB
}

func NewCoordsCache(path string) (*Coords, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	cache := Coords{db}

	return &cache, err
}

func (c *Coords) PutCoord(node structures.Node) error {
	bytes, err := binary.MarshalCoords(node)
	if err != nil {
		return err
	}

	c.Put(idToKeyBuf(node.ID), bytes, nil)
	return nil
}
