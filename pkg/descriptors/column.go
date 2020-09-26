package descriptors

import (
	"github.com/elliotcourant/noachis/pkg/types"
)

type ColumnDescriptor struct {
	Id         uint32
	Name       string
	IsNullable bool
	Type       types.Type
}
