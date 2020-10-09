package exec

import (
	"context"
	"sync/atomic"

	"github.com/elliotcourant/noachis/pkg/datums"
	"github.com/elliotcourant/noachis/pkg/descriptors"
	"github.com/elliotcourant/noachis/pkg/engine"
	"github.com/elliotcourant/noachis/pkg/kv"
	"github.com/elliotcourant/noachis/pkg/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	_ IndexReader = &indexReader{}
)

type (
	indexReader struct {
		state      int32
		index      descriptors.IndexDescriptor
		valueTypes []types.Type
		log        *logrus.Entry
		txn        engine.Transaction
		itr        engine.Iterator
	}
)

func NewIndexReader(
	ctx context.Context,
	log *logrus.Entry,
	txn engine.Transaction,
	index descriptors.IndexDescriptor,
) (IndexReader, error) {
	itr, err := txn.NewIterator()
	if err != nil {
		return nil, err
	}

	reader := indexReader{
		index:      index,
		valueTypes: make([]types.Type, len(index.StoringColumns), len(index.StoringColumns)),
		log:        log,
		txn:        txn,
		itr:        itr,
	}

	for i, col := range index.StoringColumns {
		reader.valueTypes[i] = col.Type
	}

	return &reader, nil
}

func (i *indexReader) assertNotClosed() {
	if atomic.LoadInt32(&i.state) > 0 {
		panic("index reader is closed")
	}
}

func (i *indexReader) Next(ctx context.Context) {
	i.assertNotClosed()

	i.itr.Next(ctx)
}

func (i *indexReader) Item(ctx context.Context) engine.Item {
	return i.itr.Item(ctx)
}

func (i *indexReader) Seek(ctx context.Context, indexKey datums.Datums) error {
	i.assertNotClosed()

	var key kv.Key
	var err error
	if indexKey == nil {
		key, err = kv.NewMinimumIndexKey(ctx, i.index)
	} else {
		key, err = kv.NewIndexKey(ctx, i.index, indexKey)
	}
	if err != nil {
		return err
	}

	i.itr.Seek(ctx, key, i.valueTypes)

	return nil
}

func (i *indexReader) Read(ctx context.Context, indexKey datums.Datums) (datums.Datums, error) {
	i.assertNotClosed()

	key, err := kv.NewIndexKey(ctx, i.index, indexKey)
	if err != nil {
		return nil, err
	}

	return i.txn.Get(ctx, key, i.valueTypes)
}

func (i *indexReader) Close(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&i.state, 0, 1) {
		return errors.Errorf("index reader is already closed")
	}

	return i.itr.Close(ctx)
}
