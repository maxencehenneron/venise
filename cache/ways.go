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

func (w *Ways) GetWay(wayId int64) (*structures.Way, error) {
	bytes, err := w.Get(idToKeyBuf(wayId), nil)
	if err != nil {
		return nil, err
	}

	way, err := binary.UnmarshalWay(bytes)
	if err != nil {
		return nil, err
	}

	way.ID = wayId

	return way, err
}

func (w *Ways) Iterate() chan *structures.Way {
	ways := make(chan *structures.Way)
	go func() {
		defer close(ways)
		iterator := w.NewIterator(nil, nil)
		iterator.First()
		for ; iterator.Valid(); iterator.Next() {
			way, err := binary.UnmarshalWay(iterator.Value())
			if err != nil {
				panic(err)
			}
			way.ID = idFromKeyBuf(iterator.Key())

			ways <- way
		}

	}()
	return ways
}
