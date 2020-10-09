package exec

import (
	"context"
	"fmt"
	"strings"
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
	_ IndexWriter = &indexWriter{}
)

type (
	indexWriter struct {
		state        int32
		relation     descriptors.RelationDescriptor
		index        descriptors.IndexDescriptor
		keyIndexes   []int
		valueIndexes []int
		valueTypes   []types.Type
		log          *logrus.Entry
		txn          engine.Transaction
	}
)

func NewIndexWriter(
	ctx context.Context,
	log *logrus.Entry,
	txn engine.Transaction,
	relation descriptors.RelationDescriptor,
	index descriptors.IndexDescriptor,
) (IndexWriter, error) {
	writer := indexWriter{
		relation:     relation,
		index:        index,
		keyIndexes:   make([]int, len(index.KeyColumns), len(index.KeyColumns)),
		valueIndexes: make([]int, len(index.StoringColumns), len(index.StoringColumns)),
		valueTypes:   make([]types.Type, len(index.StoringColumns), len(index.StoringColumns)),
		log:          log,
		txn:          txn,
	}

	for i, key := range index.KeyColumns {
		writer.keyIndexes[i] = relation.MustGetColumnIndex(key)
	}

	for i, value := range index.StoringColumns {
		writer.valueIndexes[i] = relation.MustGetColumnIndex(value)
		writer.valueTypes[i] = value.Type
	}

	return &writer, nil
}

func (i *indexWriter) assertNotClosed() {
	if atomic.LoadInt32(&i.state) > 0 {
		panic("index writer is closed")
	}
}

func (i *indexWriter) extractKey(ctx context.Context, row datums.Datums) (key datums.Datums, _ error) {
	key = make(datums.Datums, len(i.keyIndexes), len(i.keyIndexes))

	for x, keyIndex := range i.keyIndexes {
		key[x] = row[keyIndex]
	}

	return key, nil
}

func (i *indexWriter) convertRowToKv(ctx context.Context, row datums.Datums) (key, value datums.Datums, _ error) {
	key = make(datums.Datums, len(i.keyIndexes), len(i.keyIndexes))
	value = make(datums.Datums, len(i.valueIndexes), len(i.valueIndexes))

	for x, keyIndex := range i.keyIndexes {
		key[x] = row[keyIndex]
	}

	for x, valueIndex := range i.valueIndexes {
		value[x] = row[valueIndex]
	}

	return key, value, nil
}

func (i *indexWriter) StoreRow(ctx context.Context, row datums.Datums) error {
	i.assertNotClosed()

	key, value, err := i.convertRowToKv(ctx, row)
	if err != nil {
		return err
	}

	encodedKey, err := kv.NewIndexKey(ctx, i.index, key)
	if err != nil {
		return err
	}

	return i.txn.Set(ctx, encodedKey, value, i.valueTypes)
}

func (i *indexWriter) ValidateRow(ctx context.Context, row datums.Datums) error {
	i.assertNotClosed()

	// Current we are only doing unique indexes, if the index is not unique then
	// there is nothing to do.
	if !i.index.IsUnique {
		return nil
	}

	key, err := i.extractKey(ctx, row)
	if err != nil {
		return err
	}

	encodedKey, err := kv.NewIndexKey(ctx, i.index, key)
	if err != nil {
		return err
	}

	ok, err := i.txn.Has(ctx, encodedKey)
	if err != nil {
		return err
	}

	if ok {
		return i.uniqueConstraintViolationError(key)
	}

	return nil
}

func (i *indexWriter) Close(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&i.state, 0, 1) {
		return errors.Errorf("index writer is closed")
	}

	return nil
}

func (i *indexWriter) uniqueConstraintViolationError(key datums.Datums) error {
	columns := make([]string, len(i.index.KeyColumns), len(i.index.KeyColumns))
	for x, col := range i.index.KeyColumns {
		columns[x] = fmt.Sprintf("%s:%s", col.Name, key[x].String())
	}

	return errors.Errorf("key %s violates unique index %s", strings.Join(columns, ", "), i.index.Name)
}
