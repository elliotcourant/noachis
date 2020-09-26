package descriptors

import (
	"testing"

	"github.com/elliotcourant/noachis/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewRelation(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		table := NewRelation("users").
			WithColumn("id", types.Int8, false).
			WithColumn("email", types.Text, false).
			WithColumn("name", types.Text, true).
			WithColumn("age", types.Int2, false).
			WithPrimaryKeyColumns("id").
			WithUniqueIndex("email").
			WithNonUniqueIndex("age")

		assert.NotNil(t, table, "table must not be nil")
	})
}
