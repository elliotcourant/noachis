package datums

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullMap_FieldIsNull(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		values := Datums{
			Text("test 1"),
			Null,
			Text("test 3"),
			Text("test 4"),
			Null,
			Null,
			Text("test 7"),
			Null,
			Text("test 9"),
			Null,
		}

		nullMap := NewNullMap(values)
		assert.Equal(t,
			2, nullMap.Length(), "null map length should be 2",
		)

		for i, value := range values {
			switch value {
			case Null:
				assert.Truef(t, nullMap.FieldIsNull(i), "field %d should be null", i)
			default:
				assert.Falsef(t, nullMap.FieldIsNull(i), "field %d should not be null", i)
			}
		}
	})
}
