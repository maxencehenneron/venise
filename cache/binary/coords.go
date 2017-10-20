package binary

import (
	"encoding/binary"

	"errors"

	"fmt"

	"github.com/dernise/venise/structures"
)

const COORD_FACTOR float64 = 11930464.7083 // ((2<<31)-1)/360.0

func CoordToInt(coord float64) uint32 {
	return uint32((coord + 180.0) * COORD_FACTOR)
}

func IntToCoord(coord uint32) float64 {
	return float64((float64(coord) / COORD_FACTOR) - 180.0)
}

func MarshalCoords(node structures.Node) ([]byte, error) {
	buf := make([]byte, 2*binary.MaxVarintLen64) // Long and lat are both 10 bytes

	pos := 0
	pos += binary.PutUvarint(buf[pos:], uint64(CoordToInt(node.Lat)))
	pos += binary.PutUvarint(buf[pos:], uint64(CoordToInt(node.Lon)))
	if pos != binary.MaxVarintLen64 {
		return nil, errors.New(fmt.Sprintf("incorrect position in the byte buffer %d", pos))
	}

	return buf, nil
}

func UnmarshallCoord(buf []byte, node structures.Node) structures.Node {
	offset := 0

	lat, n := binary.Uvarint(buf[offset:])
	offset += n
	node.Lat = IntToCoord(uint32(lat))

	long, n := binary.Uvarint(buf[offset:])
	offset += n
	node.Lon = IntToCoord(uint32(long))

	return node
}
