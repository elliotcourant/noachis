package datums

import (
	"context"
	"encoding/binary"
	"io"

	"github.com/elliotcourant/noachis/pkg/types"
	"github.com/pkg/errors"
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

func DecodeText(ctx context.Context, buf io.Reader, datumType types.Type) (DText, error) {
	switch datumType.Family {
	case types.TextFamily:
		size := datumType.Width
		switch datumType.Width {
		case 0:
			sizeBytes := make([]byte, 2, 2)
			if n, err := buf.Read(sizeBytes); err != nil {
				return "", errors.Wrap(err, "failed to read size prefix")
			} else if n != 2 {
				return "", errors.Errorf("failed to read n bytes from buffer, expected %d received %d", 2, n)
			}

			size = binary.BigEndian.Uint16(sizeBytes)
			fallthrough
		default:
			text := make([]byte, size, size)
			if n, err := buf.Read(text); err != nil {
				return "", errors.Wrap(err, "failed to read text datum")
			} else if uint16(n) != size {
				return "", errors.Errorf("failed to read n bytes from buffer, expected %d received %d", size, n)
			}

			return Text(string(text)), nil
		}
	default:
		return "", errors.Errorf("cannot decode text as family %s", datumType.Family)
	}
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
