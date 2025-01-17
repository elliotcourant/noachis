package descriptors

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elliotcourant/noachis/pkg/datums"
	"github.com/elliotcourant/noachis/pkg/types"
	"github.com/pkg/errors"
)

var (
	_ datums.Datum = &RelationDescriptor{}
)

type RelationDescriptor struct {
	Oid             datums.DOid
	Name            string
	NextColumnId    uint32
	Columns         []ColumnDescriptor
	PrimaryKeyIndex IndexDescriptor
	Indexes         []IndexDescriptor
}

func DecodeRelationDescriptor(descriptor datums.DDescriptor) (RelationDescriptor, error) {
	relDesc := RelationDescriptor{}
	if err := json.Unmarshal(descriptor, &relDesc); err != nil {
		return relDesc, errors.Wrap(err, "failed to decode relation descriptor")
	}

	return relDesc, nil
}

func (r *RelationDescriptor) InferredType() types.Type {
	return types.Descriptor
}

func (r *RelationDescriptor) Encode(ctx context.Context, datumType types.Type) ([]byte, error) {
	encoded, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode relation descriptor")
	}

	return encoded, nil
}

func (r *RelationDescriptor) String() string {
	if r.Oid > 0 {
		return fmt.Sprintf("%s (%d)", r.Name, r.Oid)
	}

	return fmt.Sprintf("%s", r.Name)
}

func (r *RelationDescriptor) Raw() interface{} {
	return r
}

func NewRelation(name string) *RelationDescriptor {
	return &RelationDescriptor{
		Name:            name,
		NextColumnId:    1,
		Columns:         make([]ColumnDescriptor, 0),
		PrimaryKeyIndex: IndexDescriptor{},
		Indexes:         make([]IndexDescriptor, 0),
	}
}

func (r *RelationDescriptor) WithColumn(
	name string, columnType types.Type, nullable bool,
) *RelationDescriptor {
	columnId := r.NextColumnId
	r.NextColumnId++

	r.Columns = append(r.Columns, ColumnDescriptor{
		Id:         columnId,
		Name:       name,
		IsNullable: nullable,
		Type:       columnType,
	})

	return r
}

func (r *RelationDescriptor) WithColumnRaw(column ColumnDescriptor) *RelationDescriptor {
	r.NextColumnId = column.Id + 1
	r.Columns = append(r.Columns, column)

	return r
}

func (r *RelationDescriptor) WithPrimaryKey(
	creator func(relation *RelationDescriptor) IndexDescriptor,
) *RelationDescriptor {
	r.PrimaryKeyIndex = creator(r)

	return r
}

func (r *RelationDescriptor) WithPrimaryKeyColumns(
	columnNames ...string,
) *RelationDescriptor {
	return r.WithPrimaryKey(func(relation *RelationDescriptor) IndexDescriptor {
		index := IndexDescriptor{
			Name:           fmt.Sprintf("pk_%s", relation.Name),
			IsUnique:       true,
			KeyColumns:     make([]ColumnDescriptor, len(columnNames), len(columnNames)),
			StoringColumns: relation.Columns,
		}

		for i, name := range columnNames {
			index.KeyColumns[i] = relation.MustGetColumnByName(name)
		}

		return index
	})
}

func (r *RelationDescriptor) WithIndex(
	creator func(relation *RelationDescriptor) IndexDescriptor,
) *RelationDescriptor {
	r.Indexes = append(r.Indexes, creator(r))

	return r
}

func (r *RelationDescriptor) WithUniqueIndex(uniqueColumnNames ...string) *RelationDescriptor {
	return r.WithIndex(func(relation *RelationDescriptor) IndexDescriptor {
		index := IndexDescriptor{
			Name:           fmt.Sprintf("uq_%s_%s", relation.Name, strings.Join(uniqueColumnNames, "_")),
			IsUnique:       true,
			KeyColumns:     make([]ColumnDescriptor, len(uniqueColumnNames), len(uniqueColumnNames)),
			StoringColumns: r.PrimaryKeyIndex.KeyColumns,
		}

		for i, name := range uniqueColumnNames {
			index.KeyColumns[i] = relation.MustGetColumnByName(name)
		}

		return index
	})
}

func (r *RelationDescriptor) WithUniqueIndexId(id datums.DOid, uniqueColumnNames ...string) *RelationDescriptor {
	return r.WithIndex(func(relation *RelationDescriptor) IndexDescriptor {
		index := IndexDescriptor{
			Oid:            id,
			RelationOid:    relation.Oid,
			Name:           fmt.Sprintf("uq_%s_%s", relation.Name, strings.Join(uniqueColumnNames, "_")),
			IsUnique:       true,
			KeyColumns:     make([]ColumnDescriptor, len(uniqueColumnNames), len(uniqueColumnNames)),
			StoringColumns: r.PrimaryKeyIndex.KeyColumns,
		}

		for i, name := range uniqueColumnNames {
			index.KeyColumns[i] = relation.MustGetColumnByName(name)
		}

		return index
	})
}

func (r *RelationDescriptor) WithNonUniqueIndex(columnNames ...string) *RelationDescriptor {
	return r.WithIndex(func(relation *RelationDescriptor) IndexDescriptor {
		keySize := len(columnNames) + len(r.PrimaryKeyIndex.KeyColumns)
		index := IndexDescriptor{
			Name:           fmt.Sprintf("ix_%s_%s", relation.Name, strings.Join(columnNames, "_")),
			IsUnique:       true,
			KeyColumns:     make([]ColumnDescriptor, keySize, keySize),
			StoringColumns: r.PrimaryKeyIndex.KeyColumns,
		}

		for i, name := range columnNames {
			index.KeyColumns[i] = relation.MustGetColumnByName(name)
		}

		for i, column := range r.PrimaryKeyIndex.KeyColumns {
			index.KeyColumns[i+len(columnNames)] = column
		}

		return index
	})
}

func (r *RelationDescriptor) MustGetColumnByName(name string) ColumnDescriptor {
	for _, column := range r.Columns {
		if column.Name == name {
			return column
		}
	}

	panic(fmt.Sprintf("cannot get column by name %s", name))
}

func (r *RelationDescriptor) MustGetColumnIndex(input ColumnDescriptor) int {
	for i, column := range r.Columns {
		if column.Id == input.Id {
			return i
		}
	}

	panic(fmt.Sprintf("cannot get column index for %s (%d)", input.Name, input.Id))
}
