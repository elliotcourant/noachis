package exec

import (
	"context"

	"github.com/elliotcourant/noachis/pkg/datums"
	"github.com/elliotcourant/noachis/pkg/engine"
)

type (
	IndexReader interface {
		// Seek will move the IndexReader's cursor to the first record less than
		// or equal to the provided index key. The index key must match this
		// index's key. It cannot be an entire row.
		Seek(ctx context.Context, indexKey datums.Datums) error

		Next(ctx context.Context)

		Item(ctx context.Context) engine.Item

		// Read will retrieve a precise record from the current index.
		Read(ctx context.Context, indexKey datums.Datums) (datums.Datums, error)

		Close(ctx context.Context) error
	}

	IndexWriter interface {
		// StoreRow will take the entire row for the index's table and store the
		// data needed for this specific index.
		StoreRow(ctx context.Context, row datums.Datums) error

		// ValidateRow will take an entire row for the index's table and verify
		// that the row does not disobey this particular index's constraints (if
		// any).
		ValidateRow(ctx context.Context, row datums.Datums) error

		// Close will deallocate the IndexWriter.
		Close(ctx context.Context) error
	}
)
