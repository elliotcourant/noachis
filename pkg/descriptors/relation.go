package descriptors

import (
	"github.com/elliotcourant/noachis/pkg/datums"
)

type RelationDescriptor struct {
	Oid          datums.DOid
	Name         string
	NextColumnId uint32
	Columns      []ColumnDescriptor
}
