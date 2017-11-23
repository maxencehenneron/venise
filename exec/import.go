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

	// Setup cache
	cache := cache.NewOSMCache("bin")
	cache.Open()

	pbfReader := reader.NewPbfReader("france-latest.osm.pbf", cache)

	pbfReader.Read()
}
