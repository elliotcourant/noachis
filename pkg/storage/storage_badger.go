package storage

import (
	"encoding/hex"
	"sync"
	"sync/atomic"

	"github.com/dgraph-io/badger/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	_ Storage = &badgerStorage{}
)

type (
	badgerStorage struct {
		closed        int32
		config        Configuration
		db            *badger.DB
		log           *logrus.Entry
		sequenceCache map[string]*badgerSequence
		sequenceLock  sync.RWMutex
	}
)

func (b *badgerStorage) GetSequence(key []byte) (Sequence, error) {
	b.assertNotClosed()

	b.sequenceLock.RLock()
	sequence, ok := b.sequenceCache[string(key)]
	b.sequenceLock.RUnlock()
	if ok {
		return sequence, nil
	}

	b.sequenceLock.Lock()
	defer b.sequenceLock.Unlock()
	badgerSeq, err := b.db.GetSequence(key, 10)
	if err != nil {
		return nil, errors.Wrap(err, "failed to allocate badger sequence")
	}

	sequence = newBadgerSequence(key, badgerSeq)
	b.sequenceCache[string(key)] = sequence

	return sequence, nil
}

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

	if atomic.LoadInt32(&b.closed) > 0 {
		panic("badger storage is closed")
	}
}

func (b *badgerStorage) NewTransaction() (Transaction, error) {
	b.assertNotClosed()

	return &badgerTransaction{
		badgerTxn: b.db.NewTransaction(true),
		log:       b.log,
		storage:   b,
	}, nil
}

func (b *badgerStorage) Close() error {
	if !atomic.CompareAndSwapInt32(&b.closed, 0, 1) {
		return errors.Errorf("cannot close storage, may already be closed")
	}

	b.sequenceLock.Lock()
	defer func() {
		b.sequenceLock.Unlock()
	}()

	for key, sequence := range b.sequenceCache {
		if err := sequence.Release(); err != nil {
			b.log.
				WithField("sequence", hex.EncodeToString([]byte(key))).
				WithError(err).
				Warn("failed to release sequence")
		}

		b.sequenceCache[key] = nil
	}

	// TODO (elliotcourant) Make this a sync.Once.
	return b.db.Close()
}
