package datums

import (
	"context"
	"io"
	"math"

	"github.com/elliotcourant/noachis/pkg/types"
	"github.com/pkg/errors"
)

var (
	_ Datum = DNull{}
)

var (
	Null = DNull{}
)

type (
	DNull struct{}

	NullMap []byte
)

func (D DNull) InferredType() types.Type {
	return types.Unkown
}

func (D DNull) Encode(ctx context.Context, datumType types.Type) ([]byte, error) {
	return nil, errors.Errorf("cannot encode null, use a null map")
}

func (D DNull) String() string {
	return "NULL"
}

func (D DNull) Raw() interface{} {
	return nil
}

// NewNullMap builds an array of 8 bit unsigned integers that is used to store
// which fields are null (and therefore not written to the record) and which are
// not. The first byte is used to indicate the number of bytes in the null map.
// The following bytes each store the status for 8 fields. The highest bit being
// the lowest field. So the bits are read from left to right when checking if a
// field is null or not.
func NewNullMap(values Datums) NullMap {
	// Because we are only using a single byte to store the number of 8 bit
	// integers in the null map we are limited to 255 * 8 fields that can be
	// null mapped.
	if len(values) > 255*8 {
		panic("too many datums for null map")
	}

	// Calculate the size of the array that we will need. Each byte can store
	// whether or not up to 8 fields are null. So we need to round up to the
	// nearest whole byte.
	size := int(math.Ceil(float64(len(values)) / 8))

	// The first byte of the null map should indicate how many bytes are used to
	// represent the null map. This means that there can be a maximum of 255
	// bytes in a null map or 255 * 8 number of fields.
	nullMap := make(NullMap, 1+size, 1+size)
	nullMap[0] = uint8(size)

	// Now check which values are null. Values that are null are stored as a 1
	// in the null map. Values that are populated are stored as a 0.
	for i, value := range values {
		if value == Null {
			// Fuck you that's why.
			index := int(math.Floor(float64(i)/8)) + 1
			nullMap[index] = nullMap[index] | (1 << (7 - (i % 8)))
		}
	}

	return nullMap
}

// DecodeNullMap will parse the prefix of a byte array and return a null map
// from it. It will also return the new offset where the rest of the data for
// the record should be stored.
func DecodeNullMap(reader io.Reader) (nullMap NullMap, err error) {
	lengthByte := make([]byte, 1, 1)
	if n, err := reader.Read(lengthByte); err != nil {
		return nil, errors.Wrapf(err, "failed to read null map size")
	} else if n != 1 {
		return nil, errors.Errorf("cannot read null map")
	}

	nullMap = make(NullMap, lengthByte[0]+1, lengthByte[0]+1)
	nullMap[0] = lengthByte[0]
	n, err := reader.Read(nullMap[1:])
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read null map")
	} else if uint8(n) != lengthByte[0] {
		return nil, errors.Errorf("failed to read entire null map")
	}

	return nullMap, nil
}

// Length will return the number of bytes that this null map uses to store
// fields states. It does not return the number of fields stored in the null
// map.
func (n NullMap) Length() int {
	return int(n[0])
}

// FieldIsNull will return true if the field at the specified index is null. It
// does this using magic. It will find the byte within the array that contains
// the bit for the provided index. It will then grab the bit (from left to
// right) for that index. If that bit is 1, then the field is null. If it is 0
// then the field is not null.
func (n NullMap) FieldIsNull(index int) bool {
	// Fuck ya bois.
	return n[int(math.Floor(float64(index)/8))+1]&(1<<(7-(index%8))) != 0
}
