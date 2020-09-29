package exec

import (
	"context"
	"testing"

	"github.com/elliotcourant/noachis/pkg/engine"
	"github.com/elliotcourant/noachis/pkg/internal/testutils"
	"github.com/stretchr/testify/require"
)

func RunWithEngine(t *testing.T, testFunc func(t *testing.T, engine engine.Engine, ctx context.Context)) {
	testutils.RunWithContext(t, func(t *testing.T, ctx context.Context) {
		config := engine.Configuration{
			Directory:          "",
			InMemory:           true,
			SequenceAllocation: 10,
			Logger:             testutils.NewTestLogger(t),
		}

		eng, err := engine.NewEngine(config)
		require.NoError(t, err)
		require.NotNil(t, eng)

		testFunc(t, eng, ctx)

		err = eng.Close(ctx)
		require.NoError(t, err)
	})
}

func RunWithTransaction(t *testing.T, testFunc func(t *testing.T, txn engine.Transaction, ctx context.Context)) {
	RunWithEngine(t, func(t *testing.T, engine engine.Engine, ctx context.Context) {
		txn, err := engine.NewTransaction(ctx, t.Name())
		require.NoError(t, err)
		require.NotNil(t, txn)

		testFunc(t, txn, ctx)

		err = txn.Discard(ctx)
		require.NoError(t, err)
	})
}
