package exec

import (
	"context"
	"testing"

	"github.com/elliotcourant/noachis/pkg/datums"
	"github.com/elliotcourant/noachis/pkg/descriptors"
	"github.com/elliotcourant/noachis/pkg/engine"
	"github.com/elliotcourant/noachis/pkg/internal/testutils"
	"github.com/elliotcourant/noachis/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewIndexWriter(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		RunWithTransaction(t, func(t *testing.T, txn engine.Transaction, ctx context.Context) {
			table := descriptors.NewRelation("users").
				WithColumn("id", types.Int8, false).
				WithColumn("email", types.Text, false).
				WithColumn("name", types.Text, true).
				WithColumn("password", types.Text, false).
				WithPrimaryKeyColumns("id").
				WithUniqueIndex("email").
				WithNonUniqueIndex("email", "password")
			table.Oid = 1
			table.PrimaryKeyIndex.Oid = 2
			table.Indexes[0].Oid = 3
			table.Indexes[1].Oid = 4

			writer, err := NewIndexWriter(
				ctx,
				testutils.NewTestLogger(t),
				txn,
				*table,
				table.PrimaryKeyIndex,
			)
			assert.NoError(t, err)

			row := datums.Datums{
				datums.Int(1),
				datums.Text("email@email.com"),
				datums.Null,
				datums.Text("password"),
			}

			// Make sure that there is no error since the row does not exist
			// yet.
			err = writer.ValidateRow(ctx, row)
			assert.NoError(t, err)

			// Store the row in the index successfully.
			err = writer.StoreRow(ctx, row)
			assert.NoError(t, err)

			// Make sure that there is an error if we try to insert the same row
			// again.
			err = writer.ValidateRow(ctx, row)
			assert.EqualError(t, err, "key id:1 violates unique index pk_users")

			// Get rid of the index writer.
			err = writer.Close(ctx)
			assert.NoError(t, err)
		})
	})
}
