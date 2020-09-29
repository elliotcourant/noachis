package schema

import (
	"github.com/elliotcourant/noachis/pkg/descriptors"
	"github.com/elliotcourant/noachis/pkg/types"
)

var (
	RelationsTable = descriptors.NewRelation("relations").
			WithColumn("oid", types.OID, false).
			WithColumn("name", types.Text, false).
			WithPrimaryKeyColumns("oid").
			WithUniqueIndex("name")

	IndexesTable = descriptors.NewRelation("index").
			WithColumn("oid", types.OID, false).
			WithColumn("name", types.Text, false).
			WithColumn("relation_oid", types.OID, false).
			WithColumn("is_unique", types.Bool, false).
			WithPrimaryKeyColumns("oid").
			WithUniqueIndex("name")
)
