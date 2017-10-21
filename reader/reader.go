package reader

import (
	"log"
	"os"
	"runtime"

	"sync"

	"fmt"

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

	coordsSync := sync.WaitGroup{}
	waysSync := sync.WaitGroup{}

	coordsSync.Add(1)
	go func() {
		for node := range coords {
			pbf.cache.Coords.PutCoord(node) // Every node should be put in cache.
		}
		fmt.Println("Finished parsing coords")
		coordsSync.Done()
	}()

	coordsSync.Add(1)
	go func() {
		for node := range nodes {
			if len(node.Tags) > 0 {
				pbf.cache.Nodes.PutNode(node)
			}
		}
		fmt.Println("Finished parsing nodes")
		coordsSync.Done()
	}()

	waysSync.Add(1)
	go func() {
		coordsSync.Wait()
		fmt.Println("Started parsing ways")
		for way := range ways {
			pbf.cache.Ways.PutWay(way)
		}
		waysSync.Done()
	}()

	go func() {
		waysSync.Wait()
		fmt.Println("Started parsing relations")
		for relation := range relations {
			pbf.cache.Relations.PutRelation(relation)
		}
	}()

	pbfParser.DecodePbfData()

	return nil
}

// Verifies that the node is in the list of wanted nodes
//func (pbf *Pbf) ShouldInsert(node structures.Node) bool {
//	shouldInsert := false
//	for key, values := range pbf.tags {
//		if val, ok := node.Tags[key]; ok {
//			for _, value := range values {
//				if val == value {
//					shouldInsert = true
//				}
//			}
//		}
//	}
//	return shouldInsert
//}
