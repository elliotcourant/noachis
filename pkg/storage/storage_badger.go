package storage

import (
	"sync"

	"github.com/dgraph-io/badger/v2"
)

type (
	badgerStorage struct {
		closedSync sync.Mutex
		closed     bool
		db         *badger.DB
	}
)

func newBadgerStorage() {

}
