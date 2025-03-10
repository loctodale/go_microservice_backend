//go:build !go1.15
// +build !go1.15

package properties

import (
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

// TestMustGetParsedDuration works with go before go1.15 where the panic
// message was changed slightly. We keep this test (!) here to demonstrate the
// backwards compatibility and to keep the author happy as long as it does not
// affect any real users. Thank you! Frank :)
//
// See https://github.com/magiconair/properties/pull/63
func TestMustGetParsedDuration(t *testing.T) {
	input := "certs = 123ms\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetParsedDuration("certs"), 123*time.Millisecond)
	assert.Panic(t, func() { p.MustGetParsedDuration("key2") }, `time: invalid duration ghi`)
	assert.Panic(t, func() { p.MustGetParsedDuration("invalid") }, "unknown property: invalid")
}
