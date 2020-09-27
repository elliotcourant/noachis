package storage

import (
	"sync/atomic"

	"github.com/dgraph-io/badger/v2"
	"github.com/pkg/errors"
)

var (
	_ Sequence = &badgerSequence{}
)

type (
	badgerSequence struct {
		released int32
		key      []byte
		sequence *badger.Sequence
	}
)

func newBadgerSequence(key []byte, sequence *badger.Sequence) *badgerSequence {
	return &badgerSequence{
		key:      key,
		sequence: sequence,
	}
}

func (b *badgerSequence) assertNotReleased() {
	if b == nil {
		panic("sequence is nil")
	}

	if atomic.LoadInt32(&b.released) > 0 {
		panic("sequence is released")
	}
}

func (b *badgerSequence) Key() []byte {
	b.assertNotReleased()

	return b.key
}

func (b *badgerSequence) Next() (uint64, error) {
	b.assertNotReleased()

	return b.sequence.Next()
}

func (b *badgerSequence) Release() error {
	b.assertNotReleased()

	// We want to make sure that we are the only ones releasing the sequence.
	// So CAS the released int from a 0 to a 1. If this succeeds then we are the
	// only ones trying to release the sequence and we can proceed. If this
	// fails then there are multiple threads trying to release the sequence and
	// we need to return an error.
	if !atomic.CompareAndSwapInt32(&b.released, 0, 1) {
		return errors.Errorf("failed to release, concurrent releases")
	}

	return b.sequence.Release()
}
