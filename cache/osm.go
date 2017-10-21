package cache

import (
	"encoding/binary"
	"path/filepath"
)

type OSM struct {
	dir       string
	Nodes     *Nodes
	Coords    *Coords
	Ways      *Ways
	Relations *Relations
	opened    bool
}

func NewOSMCache(dir string) *OSM {
	cache := &OSM{dir: dir}
	return cache
}

func (c *OSM) Open() error {
	var err error
	c.Nodes, err = NewNodesCache(filepath.Join(c.dir, "nodes"))
	if err != nil {
		c.Close()
		return err
	}

	c.Coords, err = NewCoordsCache(filepath.Join(c.dir, "coords"))
	if err != nil {
		c.Close()
		return err
	}

	c.Ways, err = NewWaysCache(filepath.Join(c.dir, "ways"))
	if err != nil {
		c.Close()
		return err
	}

	c.Relations, err = NewRelationsCache(filepath.Join(c.dir, "relations"))
	if err != nil {
		c.Close()
		return err
	}

	c.opened = true
	return nil
}

func (c *OSM) Close() {
	if c.Nodes != nil {
		c.Nodes.Close()
		c.Nodes = nil
	}

	if c.Coords != nil {
		c.Coords.Close()
		c.Coords = nil
	}

	if c.Ways != nil {
		c.Ways.Close()
		c.Ways = nil
	}

	if c.Relations != nil {
		c.Relations.Close()
		c.Relations = nil
	}
}

func idToKeyBuf(id int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(id))
	return b[:8]
}
