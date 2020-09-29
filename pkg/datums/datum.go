package datums

import (
	"context"
	"strings"

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
	items := make([]string, len(d), len(d))
	for i, datum := range d {
		items[i] = datum.String()
	}

	return strings.Join(items, "|")
}
