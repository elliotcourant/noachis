package storage

import (
	"sync"

	"github.com/dgraph-io/badger/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	_ Storage = &badgerStorage{}
)

type (
	badgerStorage struct {
		closed     bool
		closedSync sync.RWMutex
		config     Configuration
		db         *badger.DB
		log        *logrus.Entry
	}
)

func newBadgerStorage(configuration Configuration) (*badgerStorage, error) {
	options := badger.DefaultOptions(configuration.Directory).
		WithInMemory(configuration.InMemory).
		WithLogger(configuration.Logger)

	db, err := badger.Open(options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open badger store")
	}

	return &badgerStorage{
		config: configuration,
		log:    configuration.Logger,
		db:     db,
	}, nil
}

func (b *badgerStorage) assertNotClosed() {
	if b == nil {
		panic("badger storage is nil")
	}
	b.closedSync.RLock()
	defer b.closedSync.RUnlock()

	if b.closed {
		panic("badger storage is closed")
	}
}

func (b *badgerStorage) Close() error {
	b.assertNotClosed()
	b.closedSync.Lock()
	defer func() {
		b.closed = true
		b.closedSync.Unlock()
	}()

	// TODO (elliotcourant) Make this a sync.Once.
	return b.db.Close()
}
