package exec

import (
	"context"
	"testing"

	"github.com/elliotcourant/noachis/pkg/engine"
	"github.com/elliotcourant/noachis/pkg/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewSchemaReader(t *testing.T) {
	t.Run("relation does not exist", func(t *testing.T) {
		RunWithTransaction(t, func(t *testing.T, txn engine.Transaction, ctx context.Context) {
			log := testutils.NewTestLogger(t)
			reader, err := NewSchemaReader(ctx, log, txn)
			assert.NoError(t, err)

			table, err := reader.GetTable(ctx, "users")
			assert.EqualError(t, err, "relation 'users' does not exist")
			assert.Empty(t, table)
		})
	})

	t.Run("list tables", func(t *testing.T) {
		RunWithTransaction(t, func(t *testing.T, txn engine.Transaction, ctx context.Context) {
			log := testutils.NewTestLogger(t)
			reader, err := NewSchemaReader(ctx, log, txn)
			assert.NoError(t, err)

			tables, err := reader.ListTables(ctx)
			assert.NoError(t, err)
			assert.NotEmpty(t, tables)

		})
	})
}
