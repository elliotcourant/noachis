package storage

import (
	"sync/atomic"

	"github.com/dgraph-io/badger/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	_ Transaction = &badgerTransaction{}
)

type (
	badgerTransaction struct {
		state     badgerTransactionState
		badgerTxn *badger.Txn
		log       *logrus.Entry
		storage   *badgerStorage
	}
)

func (b *badgerTransaction) assertNotDone() {
	if b == nil {
		panic("badger transaction is nil")
	}

	switch atomic.LoadInt32(&b.state) {
	case badgerTransactionStateActive:
		// If the transaction is active then we are good to use it.
		return
	case badgerTransactionStateDiscarded:
		panic("badger transaction is already discarded")
	case badgerTransactionStateCommitted:
		panic("badger transaction is already committed")
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
	// Try to CAS the state of the transaction to committed. This will make sure
	// that only a single thread can commit the transaction.
	if !atomic.CompareAndSwapInt32(
		&b.state, badgerTransactionStateActive, badgerTransactionStateCommitted,
	) {
		return errors.Errorf("cannot commit inactive badger transaction")
	}

	return b.badgerTxn.Commit()
}

func (b *badgerTransaction) Discard() error {
	if !atomic.CompareAndSwapInt32(
		&b.state, badgerTransactionStateActive, badgerTransactionStateDiscarded,
	) {
		return errors.Errorf("cannot discard inactive badger transaction")
	}

	b.badgerTxn.Discard()

	return nil
}
