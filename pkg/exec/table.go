package exec

import (
	"context"

	"github.com/elliotcourant/noachis/pkg/datums"
)

type (
	TableWriter interface {
		InsertRows(ctx context.Context, rows []datums.Datums) error

		DeleteRows(ctx context.Context, primaryKeys []datums.Datums) error

		UpdateRows(ctx context.Context, primaryKeys []datums.Datums, newValues []datums.Datums) error
	}
)
