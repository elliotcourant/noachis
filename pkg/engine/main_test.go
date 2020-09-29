package engine

import (
	"context"
	"testing"

	"github.com/elliotcourant/noachis/pkg/internal/testutils"
	"github.com/stretchr/testify/require"
)

func RunWithEngine(t *testing.T, testFunc func(t *testing.T, engine Engine, ctx context.Context)) {
	testutils.RunWithContext(t, func(t *testing.T, ctx context.Context) {
		config := Configuration{
			Directory:          "",
			InMemory:           true,
			SequenceAllocation: 10,
			Logger:             testutils.NewTestLogger(t),
		}

		engine, err := NewEngine(config)
		require.NoError(t, err)
		require.NotNil(t, engine)

		testFunc(t, engine, ctx)

		err = engine.Close(ctx)
		require.NoError(t, err)
	})
}
