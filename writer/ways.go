package writer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dernise/venise/cache"
	"github.com/dernise/venise/database"
	"github.com/dernise/venise/geo"
	"github.com/dernise/venise/structures"
)

type WaysWriter struct {
	OSM
}

func NewWaysWriter(cache *cache.OSM, database *database.MongoDatabase, tags map[string][]string) *WaysWriter {
	return &WaysWriter{
		OSM{cache, database, tags},
	}
}

func (ww *WaysWriter) WriteWays() {
	ways := ww.cache.Ways.Iterate()
	polygonDetecter := getPolygonDetection()

	for way := range ways {
		if ww.HasInterrestingTag(way.Tags) {
			err := ww.cache.Coords.FillWay(way)
			if err != nil {
				panic(err)
			}

			if len(way.Nodes) != len(way.NodeIDs) {
				continue
			}

			geoJson := ww.transformToGeoJson(way, polygonDetecter)
			err = ww.database.Insert(way.Tags["amenity"], geoJson)
			if err != nil {
				fmt.Printf("Error when adding the polygon %v in database\n", way.ID)
			}
		}
	}
}

func (ww *WaysWriter) transformToGeoJson(way *structures.Way, detection Detection) structures.Object {
	wayType := "LineString"
	if detection.IsPolygon(way) && way.IsClosed() {
		wayType = "Polygon"
	}

	var coordinates []geo.Point
	for _, node := range way.Nodes {
		coordinates = append(coordinates, *geo.NewPointFromLatLng(node.Lat, node.Lon))
	}

	var loc interface{}
	var center interface{}
	if wayType == "LineString" {
		loc = structures.LineString{
			"LineString",
			coordinates,
		}
	} else {
		center = structures.Point{
			Type:        "Point",
			Coordinates: geo.Polylabel(geo.Polygon{Rings: [][]geo.Point{coordinates}}, 1, false),
		}

		loc = structures.Polygon{
			Type:        "Polygon",
			Coordinates: [][]geo.Point{coordinates},
		}
	}

	return structures.Object{
		OsmId:  way.ID,
		Tags:   way.Tags,
		Loc:    loc,
		Center: center,
	}
}

func getPolygonDetection() Detection {
	raw, err := ioutil.ReadFile("./polygon-detection.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []PolygonDetection
	json.Unmarshal(raw, &c)
	return c
}
