package main

import (
	"fmt"
	"runtime"

	"github.com/dernise/venise/parser"
	"github.com/dernise/venise/reader"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(" ~===== Venise parser version 0.1 Alpha =====~ ")

	//mongoDatabase, err := mgo.Dial("localhost:27017")
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	//db, err := database.New(mongoDatabase.DB("venise"))
	//if err != nil {
	//	log.Fatal(err.Error())
	//}

	go parser.StartDetailsRoutine()

	pbfReader := reader.NewPbfReader("nord-pas-de-calais-latest.osm.pbf", nil)
	pbfReader.Read()
}
