package main

import (
	"fmt"
	"runtime"

	"github.com/dernise/venise/cache"
	"github.com/dernise/venise/parser"
	"github.com/dernise/venise/reader"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(" ~===== Venise parser version 0.1 Alpha =====~ ")

	go parser.StartDetailsRoutine()

	tags := make(map[string][]string)
	tags["amenity"] = []string{
		"bicycle_rental",
	}

	// Setup cache
	cache := cache.NewOSMCache("bin")
	cache.Open()

	pbfReader := reader.NewPbfReader("nord-pas-de-calais-latest.osm.pbf", cache, tags)

	pbfReader.Read()
}
