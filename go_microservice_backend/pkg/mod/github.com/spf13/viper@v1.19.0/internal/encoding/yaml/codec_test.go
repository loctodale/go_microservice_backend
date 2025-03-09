package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// original form of the data.
const original = `# certs-value pair
certs: value
list:
    - item1
    - item2
    - item3
map:
    certs: value

# nested
# map
nested_map:
    map:
        certs: value
        list:
            - item1
            - item2
            - item3
`

// encoded form of the data.
const encoded = `certs: value
list:
    - item1
    - item2
    - item3
map:
    certs: value
nested_map:
    map:
        certs: value
        list:
            - item1
            - item2
            - item3
`

// decoded form of the data.
//
// In case of YAML it's slightly different from Viper's internal representation
// (e.g. map is decoded into a map with interface certs).
var decoded = map[string]any{
	"certs": "value",
	"list": []any{
		"item1",
		"item2",
		"item3",
	},
	"map": map[string]any{
		"certs": "value",
	},
	"nested_map": map[string]any{
		"map": map[string]any{
			"certs": "value",
			"list": []any{
				"item1",
				"item2",
				"item3",
			},
		},
	},
}

// data is Viper's internal representation.
var data = map[string]any{
	"certs": "value",
	"list": []any{
		"item1",
		"item2",
		"item3",
	},
	"map": map[string]any{
		"certs": "value",
	},
	"nested_map": map[string]any{
		"map": map[string]any{
			"certs": "value",
			"list": []any{
				"item1",
				"item2",
				"item3",
			},
		},
	},
}

func TestCodec_Encode(t *testing.T) {
	codec := Codec{}

	b, err := codec.Encode(data)
	require.NoError(t, err)

	assert.Equal(t, encoded, string(b))
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
