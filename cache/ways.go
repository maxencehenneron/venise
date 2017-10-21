package cache

import (
	"github.com/dernise/venise/cache/binary"
	"github.com/dernise/venise/structures"
	"github.com/syndtr/goleveldb/leveldb"
)

type Ways struct {
	*leveldb.DB
}

func NewWaysCache(path string) (*Ways, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	cache := Ways{db}

	return &cache, err
}

func (w *Ways) PutWay(way structures.Way) error {
	bytes, err := binary.MarshalWay(&way)
	if err != nil {
		return err
	}

	w.Put(idToKeyBuf(way.ID), bytes, nil)
	return nil
}
