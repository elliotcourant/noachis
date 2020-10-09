package exec

import (
	"context"

	"github.com/elliotcourant/noachis/pkg/datums"
	"github.com/elliotcourant/noachis/pkg/descriptors"
	"github.com/elliotcourant/noachis/pkg/engine"
	"github.com/elliotcourant/noachis/pkg/schema"
	"github.com/elliotcourant/noachis/pkg/storage"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	_ SchemaReader = &schemaReader{}
)

type (
	schemaReader struct {
		state                 int32
		relationsByNameReader IndexReader
		relationsReader       IndexReader
		log                   *logrus.Entry
		txn                   engine.Transaction
	}
)

func NewSchemaReader(
	ctx context.Context,
	log *logrus.Entry,
	txn engine.Transaction,
) (SchemaReader, error) {
	reader := &schemaReader{
		relationsByNameReader: nil,
		relationsReader:       nil,
		log:                   log,
		txn:                   txn,
	}

	var err error
	reader.relationsByNameReader, err = NewIndexReader(ctx, log, txn, schema.RelationsByNameUniqueIndex)
	if err != nil {
		return nil, err
	}

	reader.relationsReader, err = NewIndexReader(ctx, log, txn, schema.RelationsPrimaryKeyIndex)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (s *schemaReader) GetTable(ctx context.Context, tableName string) (descriptors.RelationDescriptor, error) {
	if table, ok := schema.GetSystemTableMaybe(tableName); ok {
		return *table, nil
	}

	relationIdRecord, err := s.relationsByNameReader.Read(ctx, datums.Datums{
		datums.Text(tableName),
	})
	if errors.Is(err, storage.ErrKeyNotFound) {
		return descriptors.RelationDescriptor{}, errors.Errorf("relation '%s' does not exist", tableName)
	} else if err != nil {
		return descriptors.RelationDescriptor{}, err
	}

	return s.GetTableById(ctx, relationIdRecord[0].(datums.DOid))
}

func (s *schemaReader) GetTableById(ctx context.Context, oid datums.DOid) (descriptors.RelationDescriptor, error) {
	relationDescriptor, err := s.relationsReader.Read(ctx, datums.Datums{
		oid,
	})
	if err != nil {
		return descriptors.RelationDescriptor{}, err
	}

	return descriptors.DecodeRelationDescriptor(relationDescriptor[0].(datums.DDescriptor))
}

func (s *schemaReader) ListTables(ctx context.Context) ([]descriptors.RelationDescriptor, error) {
	if err := s.relationsReader.Seek(ctx, nil); err != nil {
		return nil, err
	}

	tables := make([]descriptors.RelationDescriptor, 0)
	for {
		var descriptorDatum datums.DDescriptor
		if item := s.relationsReader.Item(ctx).Value(); len(item) == 1 && item[0] != nil {
			descriptorDatum = item[0].(datums.DDescriptor)
		} else {
			s.log.WithContext(ctx).Warn("found invalid descriptor datum on relations primary key index")
			break
		}

		descriptor, err := descriptors.DecodeRelationDescriptor(descriptorDatum)
		if err != nil {
			return nil, err
		}

		tables = append(tables, descriptor)

		s.relationsReader.Next(ctx)
	}

	tables = append(tables, schema.SystemTablesList...)

	return tables, nil
}
