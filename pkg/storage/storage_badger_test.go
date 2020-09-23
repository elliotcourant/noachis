package storage

import (
	"testing"

	"github.com/elliotcourant/noachis/pkg/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewBadgerStorage(t *testing.T) {
	t.Run("with valid directory", func(t *testing.T) {
		dir, cleanup := testutils.NewTempDirectory(t)
		defer cleanup()

		log := testutils.NewTestLogger(t)

		config := Configuration{
			Directory: dir,
			InMemory:  false,
			Logger:    log,
		}

		db, err := newBadgerStorage(config)
		assert.NoError(t, err, "should have created successfully")
		assert.NotNil(t, db, "db should not be nil")

		err = db.Close()
		assert.NoError(t, err, "should close successfully")
	})

	t.Run("with invalid directory", func(t *testing.T) {
		log := testutils.NewTestLogger(t)

		config := Configuration{
			Directory: "/idonotexist",
			InMemory:  false,
			Logger:    log,
		}

		db, err := newBadgerStorage(config)
		assert.Error(t, err)
		assert.Nil(t, db)
	})

	t.Run("in memory", func(t *testing.T) {
		log := testutils.NewTestLogger(t)

		config := Configuration{
			InMemory: true,
			Logger:   log,
		}

		db, err := newBadgerStorage(config)
		assert.NoError(t, err, "should have created successfully")
		assert.NotNil(t, db, "db should not be nil")

		err = db.Close()
		assert.NoError(t, err, "should close successfully")
	})
}

func TestBadgerStorage_NewTransaction(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		db := NewTestStorage(t)

		transaction, err := db.NewTransaction()
		assert.NoError(t, err, "transaction should be created")
		assert.NotNil(t, transaction, "transaction should not be nil")

		err = transaction.Discard()
		assert.NoError(t, err, "transaction should discard")

		err = db.Close()
		assert.NoError(t, err, "storage should close")
	})
}
