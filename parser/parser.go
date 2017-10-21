package parser

import (
	"io"
	"log"

	"fmt"
	"time"

	"github.com/dernise/venise/decoder"
	"github.com/dernise/venise/structures"
)

var n, w, r = 0, 0, 0
var header *decoder.Header

type Pbf struct {
	numProcs  int
	decoder   *decoder.Pbf
	coords    chan structures.Node
	nodes     chan structures.Node
	ways      chan structures.Way
	relations chan structures.Relation
}

func NewPbfParser(numProcs int, decoder *decoder.Pbf, coords chan structures.Node, nodes chan structures.Node, ways chan structures.Way, relations chan structures.Relation) *Pbf {
	return &Pbf{
		numProcs,
		decoder,
		coords,
		nodes,
		ways,
		relations,
	}
}

func (p *Pbf) DecodeHeader() (*decoder.Header, error) {
	header, err := p.decoder.Header()
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (p *Pbf) DecodePbfData() error {
	// use more memory from the start, it is faster
	p.decoder.SetBufferSize(decoder.MaxBlobSize)

	// start decoding with several goroutines, it is faster
	err := p.decoder.Start(p.numProcs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		if v, err := p.decoder.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *structures.Node:
				// Process Node v.
				p.coords <- *v
				p.nodes <- *v
				n++
			case *structures.Way:
				// Process Way v.
				if w == 0 { // First way is being processed, close the other chans
					close(p.nodes)
					close(p.coords)
				}
				p.ways <- *v
				w++
			case *structures.Relation:
				// Process Relation v.
				p.relations <- *v
				r++
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}

	close(p.ways)
	close(p.relations)

	return nil
}

func StartDetailsRoutine() {
	for {
		fmt.Printf("Nodes : %d, Ways: %d, Relations: %d\n", n, w, r)
		time.Sleep(10 * time.Millisecond)
	}
}
