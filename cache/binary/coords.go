package binary

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/dernise/venise/structures"
)

func MarshalCoords(node structures.Node) ([]byte, error) {
	buf := make([]byte, 2*binary.MaxVarintLen64)

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
