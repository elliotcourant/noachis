package datums

import (
	"context"
	"encoding/binary"

	"github.com/elliotcourant/noachis/pkg/types"
)

var (
	_ Datum = DText("")
)

type (
	DText string
)

func Text(str string) DText {
	return DText(str)
}

func (d DText) InferredType() types.Type {
	return types.Text
}

func (d DText) Encode(ctx context.Context, datumType types.Type) ([]byte, error) {
	switch datumType.Family {
	case types.TextFamily:
		switch datumType.Width {
		case 0:
			buf := make([]byte, 2+len(d), 2+len(d))
			binary.BigEndian.PutUint16(buf[0:2], uint16(len(d)))
			copy(buf[2:], d)
			return buf, nil
		default:
			panic("sized text not implemented")
		}
	default:
		panic("cannot encode text as another datum type")
	}
}

func (d DText) String() string {
	return string(d)
}

func (d DText) Raw() interface{} {
	return string(d)
}
