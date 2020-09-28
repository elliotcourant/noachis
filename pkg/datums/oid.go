package datums

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/elliotcourant/noachis/pkg/types"
	"github.com/pkg/errors"
)

var (
	_ Datum = DOid(0)
)

type (
	DOid uint32
)

func Oid(i uint32) DOid {
	return DOid(i)
}

func DecodeOid(ctx context.Context, buf io.Reader, datumType types.Type) (DOid, error) {
	switch datumType.Family {
	case types.OIDFamily:
		if datumType.Width != 4 {
			return 0, errors.Errorf("cannot decode custom width oid's, not implemented")
		}

		data := make([]byte, 4, 4)
		if n, err := buf.Read(data); err != nil {
			return 0, errors.Wrap(err, "failed to read 4 bytes from buffer")
		} else if n != 4 {
			return 0, errors.Errorf("failed to read n bytes from buffer, expected %d received %d", 4, n)
		}

		val := binary.BigEndian.Uint32(data)

		return Oid(val), nil
	default:
		return 0, errors.Errorf("cannot decode oid as family %s", datumType.Family)
	}
}

func (d DOid) InferredType() types.Type {
	return types.OID
}

func (d DOid) Encode(ctx context.Context, datumType types.Type) ([]byte, error) {
	panic("implement me")
}

func (d DOid) String() string {
	return fmt.Sprint(uint32(d))
}

func (d DOid) Raw() interface{} {
	return uint32(d)
}
