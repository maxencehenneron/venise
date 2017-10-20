package binary

import (
	"testing"

	"github.com/dernise/venise/structures"
	"github.com/stretchr/testify/assert"
)

func TestMarshalCoords(t *testing.T) {
	node := structures.Node{
		123,
		3.0572560,
		50.6292050,
		nil,
		structures.Info{},
		nil,
	}

	bin, err := MarshalCoords(node)
	if err != nil {
		t.Fatal(err.Error())
	}

	unMarshalledCoord := structures.Node{
		123,
		0,
		0,
		nil,
		structures.Info{},
		nil,
	}

	unMarshalledCoord = UnmarshallCoord(bin, unMarshalledCoord)
	assert.Equal(t, 3.057255974331383, unMarshalledCoord.Lat) // close value
	assert.Equal(t, 50.629204919886945, unMarshalledCoord.Lon) // close value
}
