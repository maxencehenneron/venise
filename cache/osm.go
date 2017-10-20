package cache

import (
	"encoding/binary"
	"path/filepath"
)

type OSM struct {
	dir    string
	Nodes  *Nodes
	Coords *Coords
	opened bool
}

func NewOSMCache(dir string) *OSM {
	cache := &OSM{dir: dir}
	return cache
}

func (c *OSM) Open() error {
	var err error
	c.Nodes, err = NewNodesCache(filepath.Join(c.dir, "nodes"))
	c.Coords, err = NewCoordsCache(filepath.Join(c.dir, "coords"))
	if err != nil {
		c.Close()
		return err
	}

	c.opened = true
	return nil
}

func (c *OSM) Close() {
	if c.Nodes != nil {
		c.Close()
		c.Nodes = nil
	}
}

func idToKeyBuf(id int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(id))
	return b[:8]
}
