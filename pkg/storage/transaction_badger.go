package storage

import (
	"sync"

	"github.com/dgraph-io/badger/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	_ Transaction = &badgerTransaction{}
)

type (
	badgerTransaction struct {
		badgerTxn *badger.Txn
		done      bool
		doneLock  sync.RWMutex
		log       *logrus.Entry
		storage   *badgerStorage
	}
)

func (b *badgerTransaction) assertNotDone() {
	if b == nil {
		panic("badger transaction is nil")
	}

	b.doneLock.RLock()
	defer b.doneLock.RUnlock()

	if b.done {
		panic("badger transaction is committed or discarded")
	}
}

func (b *badgerTransaction) Get(key []byte) (Item, error) {
	b.assertNotDone()

	item, err := b.badgerTxn.Get(key)
	switch err {
	case badger.ErrKeyNotFound:
		return nil, errors.WithStack(ErrKeyNotFound)
	case nil:
		break
	default:
		return nil, errors.Wrap(err, "failed to retrieve item")
	}

	return newBadgerItem(item), nil
}

func (b *badgerTransaction) Has(key []byte) (ok bool, _ error) {
	b.assertNotDone()

	_, err := b.badgerTxn.Get(key)
	switch err {
	case badger.ErrKeyNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, errors.Wrap(err, "failed to retrieve item")
	}
}

func (b *badgerTransaction) Set(key, value []byte) error {
	b.assertNotDone()

	return b.badgerTxn.Set(key, value)
}

func (b *badgerTransaction) Commit() error {
	b.assertNotDone()

	b.doneLock.Lock()
	defer func() {
		b.done = true
		b.doneLock.Unlock()
	}()

	return b.badgerTxn.Commit()
}

func (b *badgerTransaction) Discard() error {
	b.assertNotDone()

	b.doneLock.Lock()
	defer func() {
		b.done = true
		b.doneLock.Unlock()
	}()

	b.badgerTxn.Discard()

	return nil
}
