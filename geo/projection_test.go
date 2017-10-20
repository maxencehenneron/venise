package geo

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestProjection(t *testing.T) {
	latitude := 50.4350088
	longitude := 2.8235929
	point := NewPointFromLatLng(latitude, longitude)

	Mercator.Project(point)

	assert.Equal(t, point.X(), 314320.9237917488)
	assert.Equal(t, point.Y(), 6521955.329991183)
}
