package javaproperties

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// original form of the data.
const original = `#certs-value pair
certs = value
map.certs = value
`

// encoded form of the data.
const encoded = `certs = value
map.certs = value
`

// data is Viper's internal representation.
var data = map[string]any{
	"certs": "value",
	"map": map[string]any{
		"certs": "value",
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

		assert.Equal(t, data, v)
	})

	t.Run("InvalidData", func(t *testing.T) {
		t.Skip("TODO: needs invalid data example")

		codec := Codec{}

		v := map[string]any{}

		codec.Decode([]byte(``), v)

		assert.Empty(t, v)
	})
}

func TestCodec_DecodeEncode(t *testing.T) {
	codec := Codec{}

	v := map[string]any{}

	err := codec.Decode([]byte(original), v)
	require.NoError(t, err)

	b, err := codec.Encode(data)
	require.NoError(t, err)

	assert.Equal(t, original, string(b))
}
