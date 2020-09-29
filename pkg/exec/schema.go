package exec

import (
	"context"

	"github.com/elliotcourant/noachis/pkg/datums"
	"github.com/elliotcourant/noachis/pkg/descriptors"
)

type (
	SchemaReader interface {
		GetTable(ctx context.Context, tableName string) (descriptors.RelationDescriptor, error)

		GetTableById(ctx context.Context, oid datums.DOid) (descriptors.RelationDescriptor, error)

		ListTables(ctx context.Context) ([]descriptors.RelationDescriptor, error)
	}

	SchemaWriter interface {
		CreateTable(ctx context.Context, table *descriptors.RelationDescriptor) error

		DropTable(ctx context.Context, tableName string) error
	}
)
