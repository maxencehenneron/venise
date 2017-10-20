package geo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/magiconair/properties/assert"
)

type JsonParser struct {
	Rings [][]Point `json:"rings"`
}

func TestPolylabel(t *testing.T) {
	raw, err := ioutil.ReadFile("./water.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var p JsonParser
	json.Unmarshal(raw, &p)

	polygon := Polygon{p.Rings}
	bestPoint := Polylabel(polygon, 1, true)
	assert.Equal(t, bestPoint, Point{3865.85009765625, 2124.87841796875})
}
