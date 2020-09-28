package kv

import (
	"context"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/elliotcourant/noachis/pkg/datums"
	"github.com/elliotcourant/noachis/pkg/descriptors"
	"github.com/pkg/errors"
)

type Key interface {
	Bytes() []byte

	String() string
}

var (
	_ Key = IndexKey{}
	_ Key = SequenceKey{}
)

func NewIndexKey(
	ctx context.Context,
	index descriptors.IndexDescriptor,
	keyDatums datums.Datums,
) (Key, error) {
	if index.Oid == 0 {
		return nil, errors.Errorf("index is not initialized")
	}

	inputKeyLength, actualKeyLength := len(keyDatums), len(index.KeyColumns)
	if inputKeyLength != actualKeyLength {
		return nil, errors.Errorf("index key length mismatch, received: %d expected: %d", inputKeyLength, actualKeyLength)
	}

	path := make([]string, len(keyDatums)+1, len(keyDatums)+1)
	path[0] = index.Name

	buf := make([]byte, 5)
	buf[0] = RecordKeyType
	binary.BigEndian.PutUint32(buf[1:5], uint32(index.Oid))

	for i, datum := range keyDatums {
		path[i+1] = datum.String()
		encoded, err := datum.Encode(ctx, index.KeyColumns[i].Type)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to encode datum %d for index %s key", i, index.Name)
		}

		buf = append(buf, encoded...)
	}

	return IndexKey{
		Key:  buf,
		Path: "/" + strings.Join(path, "/"),
	}, nil
}

type IndexKey struct {
	Key  []byte
	Path string
}

func (i IndexKey) Bytes() []byte {
	return i.Key
}

func (i IndexKey) String() string {
	return "/index" + i.Path
}

func NewSequenceKey(sequence descriptors.SequenceDescriptor) (Key, error) {
	if sequence.Oid == 0 {
		return nil, errors.Errorf("sequence not initialized")
	}

	buf := make([]byte, 5)
	buf[0] = SequenceKeyType
	binary.BigEndian.PutUint32(buf[1:5], uint32(sequence.Oid))

	return SequenceKey{
		Key:  buf,
		Path: fmt.Sprintf("/%s<%d>", sequence.Name, sequence.Oid),
	}, nil
}

type SequenceKey struct {
	Key  []byte
	Path string
}

func (s SequenceKey) Bytes() []byte {
	return s.Key
}

func (s SequenceKey) String() string {
	return "/sequence" + s.Path
}
