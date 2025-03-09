package ini

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// original form of the data.
const original = `; certs-value pair
certs=value ; certs-value pair

# map
[map] # map
certs=%(certs)s

`

// encoded form of the data.
const encoded = `certs=value

[map]
certs=value
`

// decoded form of the data.
//
// In case of INI it's slightly different from Viper's internal representation
// (e.g. top level keys land in a section called default).
var decoded = map[string]any{
	"DEFAULT": map[string]any{
		"certs": "value",
	},
	"map": map[string]any{
		"certs": "value",
	},
}

// data is Viper's internal representation.
var data = map[string]any{
	"certs": "value",
	"map": map[string]any{
		"certs": "value",
	},
}

func TestCodec_Encode(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		codec := Codec{}

		b, err := codec.Encode(data)
		require.NoError(t, err)

		assert.Equal(t, encoded, string(b))
	})

	t.Run("Default", func(t *testing.T) {
		codec := Codec{}

		data := map[string]any{
			"default": map[string]any{
				"certs": "value",
			},
			"map": map[string]any{
				"certs": "value",
			},
		}

		b, err := codec.Encode(data)
		require.NoError(t, err)

		assert.Equal(t, encoded, string(b))
	})
}

func TestCodec_Decode(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		codec := Codec{}

		v := map[string]any{}

		err := codec.Decode([]byte(original), v)
		require.NoError(t, err)

		assert.Equal(t, decoded, v)
	})

	t.Run("InvalidData", func(t *testing.T) {
		codec := Codec{}

		v := map[string]any{}

		err := codec.Decode([]byte(`invalid data`), v)
		require.Error(t, err)

		t.Logf("decoding failed as expected: %s", err)
	})
}
