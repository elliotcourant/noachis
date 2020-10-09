package engine

import (
	"context"

	"github.com/elliotcourant/noachis/pkg/kv"
	"github.com/elliotcourant/noachis/pkg/types"
)

var (
	_ Iterator = &iteratorBase{}
)

type (
	iteratorBase struct {
	}
)

func (i iteratorBase) Seek(ctx context.Context, key kv.Key, datumTypes []types.Type) {
	panic("implement me")
}

func (i iteratorBase) Next(ctx context.Context) {
	panic("implement me")
}

func (i iteratorBase) Previous(ctx context.Context) {
	panic("implement me")
}

func (i iteratorBase) Valid() bool {
	panic("implement me")
}

func (i iteratorBase) Item(ctx context.Context) Item {
	panic("implement me")
}

func (i iteratorBase) Close(ctx context.Context) error {
	panic("implement me")
}
