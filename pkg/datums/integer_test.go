package datums

import (
	"context"
	"testing"

	"github.com/elliotcourant/noachis/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestDInt_Encode(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		val := DInt(1234)
		buf, err := val.Encode(context.Background(), types.Int8)
		assert.NoError(t, err, "should succeed")
		assert.Len(t, buf, 8)
	})
}
