package types

type Type struct {
	Family      Family
	Width       uint16
	ElementType *Type
}

var (
	Unkown = Type{
		Family:      UnknownFamily,
		Width:       0,
		ElementType: nil,
	}

	Int2 = Type{
		Family:      IntegerFamily,
		Width:       2,
		ElementType: nil,
	}

	Int4 = Type{
		Family:      IntegerFamily,
		Width:       4,
		ElementType: nil,
	}

	Int8 = Type{
		Family:      IntegerFamily,
		Width:       8,
		ElementType: nil,
	}

	Text = Type{
		Family:      TextFamily,
		Width:       0,
		ElementType: nil,
	}

	Bool = Type{
		Family:      BooleanFamily,
		Width:       1,
		ElementType: nil,
	}

	OID = Type{
		Family:      OIDFamily,
		Width:       4,
		ElementType: nil,
	}

	Descriptor = Type{
		Family:      DescriptorFamily,
		Width:       0,
		ElementType: nil,
	}
)
