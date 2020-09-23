package storage

import (
	"github.com/dgraph-io/badger/v2"
)

var (
	_ Item = badgerItem{}
)

type (
	badgerItem struct {
		item *badger.Item
	}
)

func newBadgerItem(item *badger.Item) Item {
	return badgerItem{
		item: item,
	}
}

func (b badgerItem) ValueCopy(destination []byte) ([]byte, error) {
	return b.item.ValueCopy(destination)
}

func (b badgerItem) KeyCopy(destination []byte) ([]byte, error) {
	return b.item.KeyCopy(destination), nil
}
