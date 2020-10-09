package schema

import (
	"math"

	"github.com/elliotcourant/noachis/pkg/datums"
)

const (
	OIDSequenceId datums.DOid = math.MaxInt32 - iota
	RelationsTableId
	RelationsPrimaryKeyIndexId
	RelationsByNameUniqueIndexId
)
