package cache

import (
	"github.com/dernise/venise/cache/binary"
	"github.com/dernise/venise/structures"
	"github.com/syndtr/goleveldb/leveldb"
)

type Relations struct {
	*leveldb.DB
}

func NewRelationsCache(path string) (*Relations, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	cache := Relations{db}

	return &cache, err
}

func (r *Relations) PutRelation(relation structures.Relation) error {
	if len(relation.Tags) == 0 {
		return nil
	}

	bytes, err := binary.MarshalRelation(&relation)
	if err != nil {
		return err
	}

	r.Put(idToKeyBuf(relation.ID), bytes, nil)
	return nil
}

func (r *Relations) Iterate() chan *structures.Relation {
	relations := make(chan *structures.Relation)
	go func() {
		defer close(relations)
		iterator := r.NewIterator(nil, nil)
		iterator.First()
		for ; iterator.Valid(); iterator.Next() {
			relation, err := binary.UnmarshalRelation(iterator.Value())
			if err != nil {
				panic(err)
			}
			relation.ID = idFromKeyBuf(iterator.Key())

			relations <- relation
		}

	}()
	return relations
}
