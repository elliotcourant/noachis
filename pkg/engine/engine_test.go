package engine

import (
	"context"
	"testing"

	"github.com/elliotcourant/noachis/pkg/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewEngine(t *testing.T) {
	t.Run("in memory", func(t *testing.T) {
		testutils.RunWithContext(t, func(t *testing.T, ctx context.Context) {
			config := Configuration{
				Directory:          "",
				InMemory:           true,
				SequenceAllocation: 10,
				Logger:             testutils.NewTestLogger(t),
			}

			engine, err := NewEngine(config)
			assert.NoError(t, err)
			assert.NotNil(t, engine)

			err = engine.Close(ctx)
			assert.NoError(t, err)
		})
	})
}

func TestEngineBase_NewTransaction(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		RunWithEngine(t, func(t *testing.T, engine Engine, ctx context.Context) {
			txn, err := engine.NewTransaction(ctx, "test")
			assert.NoError(t, err)
			assert.NotNil(t, txn)

			err = txn.Discard(ctx)
			assert.NoError(t, err)
		})
	})
}
