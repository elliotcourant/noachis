package engine

import (
	"context"
	"sync/atomic"

	"github.com/elliotcourant/noachis/pkg/datums"
	"github.com/elliotcourant/noachis/pkg/kv"
	"github.com/elliotcourant/noachis/pkg/schema"
	"github.com/elliotcourant/noachis/pkg/storage"
	"github.com/elliotcourant/noachis/pkg/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	_ Transaction = &transactionBase{}
)

type (
	engineTransactionState = int32

	transactionBase struct {
		state  int32
		engine *engineBase
		log    *logrus.Entry
		txn    storage.Transaction
	}
)

const (
	engineTransactionStateActive engineTransactionState = iota
	engineTransactionStateDiscarded
	engineTransactionStateCommitted
)

func (t *transactionBase) assertActive() {
	switch atomic.LoadInt32(&t.state) {
	case engineTransactionStateActive:
		return
	case engineTransactionStateDiscarded:
		panic("transaction is discarded")
	case engineTransactionStateCommitted:
		panic("transaction is committed")
	}
}

func (t *transactionBase) NewObjectId(ctx context.Context) (datums.DOid, error) {
	t.assertActive()

	key, err := kv.NewSequenceKey(schema.OIDSequence)
	if err != nil {
		return 0, err
	}

	sequence, err := t.engine.db.GetSequence(key.Bytes())
	if err != nil {
		return 0, err
	}

	id, err := sequence.Next()
	if err != nil {
		return 0, errors.Wrap(err, "failed to generate new object id")
	}

	return datums.Oid(uint32(id)), nil
}

func (t *transactionBase) Set(ctx context.Context, key kv.Key, datums datums.Datums, datumTypes []types.Type) error {
	t.assertActive()

	t.log.
		WithContext(ctx).
		WithField("key", key.String()).
		Tracef("writing key")

	data, err := kv.EncodeRow(ctx, datums, datumTypes)
	if err != nil {
		return err
	}

	return t.txn.Set(key.Bytes(), data)
}

func (t *transactionBase) Get(ctx context.Context, key kv.Key, datumTypes []types.Type) (result datums.Datums, _ error) {
	t.assertActive()

	t.log.
		WithContext(ctx).
		WithField("key", key.String()).
		Tracef("retrieving key")

	item, err := t.txn.Get(key.Bytes())
	if err != nil {
		return nil, err
	}

	if err := item.Value(func(data []byte) error {
		result, err = kv.DecodeRow(ctx, data, datumTypes)
		return err
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (t *transactionBase) Has(ctx context.Context, key kv.Key) (ok bool, _ error) {
	t.assertActive()

	t.log.
		WithContext(ctx).
		WithField("key", key.String()).
		Tracef("checking if key exists")

	return t.txn.Has(key.Bytes())
}

func (t *transactionBase) NewIterator() (Iterator, error) {
	panic("implement me")
}

func (t *transactionBase) Commit(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(
		&t.state,
		engineTransactionStateActive,
		engineTransactionStateCommitted,
	) {
		return errors.Errorf("cannot commit inactive transaction")
	}

	return t.txn.Commit()
}

func (t *transactionBase) Discard(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(
		&t.state,
		engineTransactionStateActive,
		engineTransactionStateDiscarded,
	) {
		return errors.Errorf("cannot discard inactive transaction")
	}

	return t.txn.Discard()
}
