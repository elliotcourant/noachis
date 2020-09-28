package datums

import (
	"context"

	"github.com/elliotcourant/noachis/pkg/types"
)

type (
	Datum interface {
		InferredType() types.Type

		Encode(ctx context.Context, datumType types.Type) ([]byte, error)

		String() string

		Raw() interface{}
	}

	Datums []Datum
)

func (d Datums) String() string {
	panic("not implemented")
}
