package schema

import (
	"github.com/elliotcourant/noachis/pkg/descriptors"
	"github.com/elliotcourant/noachis/pkg/types"
)

var (
	RelationsColumnOid = descriptors.ColumnDescriptor{
		Id:         1,
		Name:       "oid",
		IsNullable: false,
		Type:       types.OID,
	}

	RelationsColumnName = descriptors.ColumnDescriptor{
		Id:         2,
		Name:       "name",
		IsNullable: false,
		Type:       types.Text,
	}

	RelationsPrimaryKeyIndex = descriptors.IndexDescriptor{
		Oid:         RelationsPrimaryKeyIndexId,
		Name:        "pk_relations",
		RelationOid: RelationsTableId,
		IsUnique:    true,
		KeyColumns: []descriptors.ColumnDescriptor{
			RelationsColumnOid,
		},
		StoringColumns: []descriptors.ColumnDescriptor{
			{
				Id:         3,
				Name:       "data",
				IsNullable: false,
				Type:       types.Descriptor,
			},
		},
	}

	RelationsByNameUniqueIndex = descriptors.IndexDescriptor{
		Oid:         RelationsByNameUniqueIndexId,
		Name:        "uq_relations_name",
		RelationOid: RelationsTableId,
		IsUnique:    true,
		KeyColumns: []descriptors.ColumnDescriptor{
			RelationsColumnName,
		},
		StoringColumns: []descriptors.ColumnDescriptor{
			RelationsColumnOid,
		},
	}

	RelationsTable = descriptors.NewRelation("relations").
			WithColumnRaw(RelationsColumnOid).
			WithColumnRaw(RelationsColumnName).
			WithPrimaryKey(func(relation *descriptors.RelationDescriptor) descriptors.IndexDescriptor {
			return RelationsPrimaryKeyIndex
		}).
		WithIndex(func(relation *descriptors.RelationDescriptor) descriptors.IndexDescriptor {
			return RelationsByNameUniqueIndex
		})

	IndexesTable = descriptors.NewRelation("indexes").
			WithColumn("oid", types.OID, false).
			WithColumn("name", types.Text, false).
			WithColumn("relation_oid", types.OID, false).
			WithColumn("is_unique", types.Bool, false).
			WithPrimaryKeyColumns("oid").
			WithUniqueIndex("name")
)

var (
	SystemTables = map[string]*descriptors.RelationDescriptor{
		"relations": RelationsTable,
		"indexes":   IndexesTable,
	}

	SystemTablesList = []descriptors.RelationDescriptor{
		*RelationsTable,
		*IndexesTable,
	}
)

func GetSystemTableMaybe(tableName string) (*descriptors.RelationDescriptor, bool) {
	table, ok := SystemTables[tableName]
	return table, ok
}

func IsSystemTable(tableName string) bool {
	_, ok := SystemTables[tableName]
	return ok
}
