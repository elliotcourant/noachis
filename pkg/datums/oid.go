package datums

import (
	"context"
	"fmt"

	"github.com/elliotcourant/noachis/pkg/types"
)

var (
	_ Datum = DOid(0)
)

type (
	DOid uint32
)

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
