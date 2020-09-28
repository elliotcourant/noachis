package kv

type KeyType = byte

const (
	UnknownKeyType KeyType = iota
	SequenceKeyType
	RecordKeyType
)
