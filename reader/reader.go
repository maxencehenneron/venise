package reader

import (
	"log"
	"os"
	"runtime"

	"github.com/dernise/venise/cache"
	"github.com/dernise/venise/decoder"
	"github.com/dernise/venise/parser"
	"github.com/dernise/venise/structures"
)

type Pbf struct {
	fileName string
	cache    *cache.OSM
}

func NewPbfReader(fileName string, cache *cache.OSM) *Pbf {
	return &Pbf{fileName, cache}
}

func (pbf *Pbf) Read() error {
	nodes := make(chan structures.Node, 100)
	coords := make(chan structures.Node, 100)
	ways := make(chan structures.Way, 100)
	relations := make(chan structures.Relation, 100)

	f, err := os.Open(pbf.fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := decoder.NewDecoder(f)

	pbfParser := parser.NewPbfParser(runtime.NumCPU(), decoder, coords, nodes, ways, relations)

	cache := cache.NewOSMCache("bin")
	cache.Open()

	go func() {
		for range coords {

		}
	}()

	go func() {
		for node := range nodes {
			cache.Nodes.PutNode(node)
		}
	}()

	go func() {
		for range relations {

		}
	}()

	go func() {
		for range ways {

		}
	}()

	pbfParser.DecodePbfData()

	return nil
}
