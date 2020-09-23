package types

//go:generate stringer -type=Family -output=family.strings.go
type Family uint8

const (
	UnknownFamily Family = iota
	IntegerFamily
	BooleanFamily
	TextFamily
	ArrayFamily
	OIDFamily
	ObjectFamily
)
