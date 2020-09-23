package datums

import (
	"context"

	"github.com/elliotcourant/noachis/pkg/types"
	"github.com/pkg/errors"
)

var (
	_ Datum = DNull{}
)

var (
	Null = DNull{}
)

type DNull struct{}

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
