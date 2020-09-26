package descriptors

import (
	"github.com/elliotcourant/noachis/pkg/datums"
)

type SequenceDescriptor struct {
	Oid    datums.DOid
	Name   string
	Parent *datums.DOid
}
