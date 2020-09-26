package descriptors

import (
	"github.com/elliotcourant/noachis/pkg/datums"
)

type IndexDescriptor struct {
	Oid            datums.DOid
	Name           string
	RelationOid    datums.DOid
	IsUnique       bool
	KeyColumns     []ColumnDescriptor
	StoringColumns []ColumnDescriptor
}
