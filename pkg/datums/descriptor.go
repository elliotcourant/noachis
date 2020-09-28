package datums

import (
	"context"
	"encoding/json"

	"github.com/elliotcourant/noachis/pkg/types"
)

var (
	_ Datum = DDescriptor("")
)

type DDescriptor []byte

func (d DDescriptor) InferredType() types.Type {
	return types.Descriptor
}

func (d DDescriptor) Encode(ctx context.Context, datumType types.Type) ([]byte, error) {
	panic("implement me")
}

func (d DDescriptor) String() string {
	return string(d)
}

func (d DDescriptor) Raw() interface{} {
	var obj interface{}
	if err := json.Unmarshal(d, &obj); err != nil {
		return []byte(d)
	}

	return obj
}
