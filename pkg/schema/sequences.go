package schema

import (
	"github.com/elliotcourant/noachis/pkg/descriptors"
)

var (
	OIDSequence = descriptors.SequenceDescriptor{
		Oid:    OIDSequenceId,
		Name:   "oid_sequence",
		Parent: nil,
	}
)
