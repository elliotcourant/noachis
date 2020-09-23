package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadgerTransaction_SetGetHas(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		txn, cleanup := NewTestTransaction(t)
		defer cleanup()

		key, value := []byte("key"), []byte("value")

		{ // Do a read before the key exists to test key not found.
			item, err := txn.Get(key)
			assert.EqualError(t, err, ErrKeyNotFound.Error())
			assert.Nil(t, item, "item should be nil")

			ok, err := txn.Has(key)
			assert.NoError(t, err, "has should succeed")
			assert.False(t, ok, "key should not exist")
		}

		// Then we will set the key.
		err := txn.Set(key, value)
		assert.NoError(t, err, "should set successfully")

		item, err := txn.Get(key)
		assert.NoError(t, err, "read should succeed")
		assert.NotNil(t, item, "item should not be nil")

		ok, err := txn.Has(key)
		assert.NoError(t, err, "has should succeed")
		assert.True(t, ok, "key should exist")
	})
}
