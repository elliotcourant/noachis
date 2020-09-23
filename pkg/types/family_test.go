package types

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFamily_String(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		families := []Family{
			UnknownFamily,
			IntegerFamily,
			BooleanFamily,
			TextFamily,
			ArrayFamily,
			OIDFamily,
			ObjectFamily,
		}

		for _, family := range families {
			str := family.String()
			assert.NotEmpty(t, str, "string should not be empty")
		}
	})

	t.Run("invalid", func(t *testing.T) {
		badFamily := Family(math.MaxUint8)
		str := badFamily.String()
		assert.Equal(t, "Family(255)", str)
	})
}
