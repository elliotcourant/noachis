package storage

import (
	"testing"

	"github.com/elliotcourant/noachis/pkg/internal/testutils"
	"github.com/stretchr/testify/require"
)

func NewTestStorage(t *testing.T) Storage {
	log := testutils.NewTestLogger(t)

	config := Configuration{
		InMemory: true,
		Logger:   log,
	}

	db, err := NewStorage(config)
	require.NoError(t, err, "should create storage")
	require.NotNil(t, db)

	return db
}

func NewTestTransaction(t *testing.T) (txn Transaction, cleanup func()) {
	db := NewTestStorage(t)

	tx, err := db.NewTransaction()
	require.NoError(t, err, "should not have error")
	require.NotNil(t, tx, "transaction should not be nil")

	return tx, func() {
		require.NoError(t, tx.Discard(), "should discard")
		require.NoError(t, db.Close(), "should close")
	}
}
