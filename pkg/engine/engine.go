package engine

import (
	"context"
	"sync/atomic"

	"github.com/elliotcourant/noachis/pkg/datums"
	"github.com/elliotcourant/noachis/pkg/kv"
	"github.com/elliotcourant/noachis/pkg/storage"
	"github.com/elliotcourant/noachis/pkg/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type (
	Configuration struct {
		Directory          string
		InMemory           bool
		SequenceAllocation int
		Logger             *logrus.Entry
	}

	Engine interface {
		NewTransaction(ctx context.Context, sessionId string) (Transaction, error)

		Close(ctx context.Context) error
	}

	Transaction interface {
		NewObjectId(ctx context.Context) (datums.DOid, error)

		Set(
			ctx context.Context,
			key kv.Key, datums datums.Datums, datumTypes []types.Type,
		) error

		Get(
			ctx context.Context,
			key kv.Key, datumTypes []types.Type,
		) (datums.Datums, error)

		Has(
			ctx context.Context, key kv.Key,
		) (ok bool, _ error)

		NewIterator() (Iterator, error)

		Commit(ctx context.Context) error

		Discard(ctx context.Context) error
	}

	Iterator interface {
		Seek(
			ctx context.Context,
			key kv.Key, datumTypes []types.Type,
		)

		Next(ctx context.Context)

		Previous(ctx context.Context)

		Valid() bool

		Item(ctx context.Context) Item

		Close(ctx context.Context) error
	}

	Item interface {
		Key() datums.Datums

		Value() datums.Datums
	}
)

type (
	engineBase struct {
		state  int32
		config Configuration
		log    *logrus.Entry
		db     storage.Storage
	}
)

func (e *engineBase) NewTransaction(ctx context.Context, sessionId string) (Transaction, error) {
	txn, err := e.db.NewTransaction()
	if err != nil {
		return nil, err
	}

	return &transactionBase{
		engine: e,
		log:    e.log,
		txn:    txn,
	}, nil
}

func (e *engineBase) Close(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&e.state, 0, 1) {
		return errors.Errorf("engine already closed")
	}

	return e.db.Close()
}

func NewEngine(config Configuration) (Engine, error) {
	db, err := storage.NewStorage(storage.Configuration{
		Directory:          config.Directory,
		InMemory:           config.InMemory,
		SequenceAllocation: config.SequenceAllocation,
		Logger:             config.Logger,
	})
	if err != nil {
		return nil, err
	}

	return &engineBase{
		config: config,
		log:    config.Logger,
		db:     db,
	}, nil
}
