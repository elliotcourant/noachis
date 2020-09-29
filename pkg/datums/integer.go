package datums

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/elliotcourant/noachis/pkg/types"
)

var (
	_ Datum = DInt(0)
)

type (
	DInt int64
)

func Int(i int64) DInt {
	return DInt(i)
}

func (d DInt) InferredType() types.Type {
	return types.Int8
}

func (d DInt) Encode(ctx context.Context, datumType types.Type) ([]byte, error) {
	destination := make([]byte, 0)
	switch datumType.Family {
	case types.IntegerFamily:
		data := make([]byte, 8, 8)
		binary.BigEndian.PutUint64(data, uint64(d))
		switch datumType.Width {
		case 0:
			// If there is no width specified then that means we need to write
			// this integer with a 16 bit length prefix.
			destination = append(destination, make([]byte, 2, 2)...)
			binary.BigEndian.PutUint16(destination[0:2], uint16(len(data)))
			fallthrough
		case 8:
			destination = append(destination, data...)
			return destination, nil
		default:
			panic("custom lengths are not yet implemented")
		}
	case types.TextFamily:
		panic("encoding integer as text not yet implemented")
	}
	panic("implement me")
}

func (d DInt) String() string {
	return fmt.Sprint(int64(d))
}

func (d DInt) Raw() interface{} {
	return int64(d)
}
