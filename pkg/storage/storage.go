package storage

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type (
	Configuration struct {
		Directory string

		InMemory bool

		SequenceAllocation int

		Logger *logrus.Entry
	}

	Storage interface {
		// Close will shut down the storage interface. Any operations that can be
		// stoppped will be, but any operations that are attempted after the close
		// is complete will fail. If an error is returned the storage object has
		// still been closed, the error is only a representation of any problems
		// that have happened during closing.
		Close() error
	}

	Item interface {
		// ValueCopy will take the value data store in the database and copy it to
		// the specified byte array. This is done because the byte array for the
		// value on this item is reused between get and iteration operations. So
		// simply doing a "soft" copy of the value would cause an issue where the
		// value could change after being read. This prevents that be copying the
		// value data to the destination to make it separate.
		ValueCopy(destination []byte) ([]byte, error)
	}

	Transaction interface {
		// Get will retrieve the Item from the storage with the specified key. If
		// that item does not exist then an ErrKeyNotFound will be returned.
		Get(key []byte) (Item, error)

		// Has will check and see if there is an Item in the storage with the
		// specified key. This is slightly more efficient than a Get because it
		// does not try to read any data from the storage. It only checks to see if
		// it is there.
		Has(key []byte) (ok bool, _ error)

		// Set will store the specified key value pair in the database as an item.
		Set(key, value []byte) error

		// Commit will persist all of the changes in the transaction to the
		// database.
		Commit() error

		// Discard will throw out all of the changes made during the transaction.
		Discard() error
	}
)
