package datums

import (
	"context"

	"github.com/elliotcourant/noachis/pkg/types"
)

type (
	Datum interface {
		InferredType() types.Type

		Encode(ctx context.Context, destination []byte, datumType types.Type) error

		String() string

		Raw() interface{}
	}
)
