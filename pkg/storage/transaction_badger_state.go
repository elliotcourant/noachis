package storage

type badgerTransactionState = int32

const (
	badgerTransactionStateActive badgerTransactionState = iota
	badgerTransactionStateDiscarded
	badgerTransactionStateCommitted
)
