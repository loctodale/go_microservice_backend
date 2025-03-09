// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"testing"
	"time"
)

func TestAttrNoAlloc(t *testing.T) {
	// Assign values just to make sure the compiler doesn't optimize away the statements.
	var (
		i int64
		u uint64
		f float64
		b bool
		s string
		x any
		p = &i
		d time.Duration
	)
	a := int(testing.AllocsPerRun(5, func() {
		i = Int64("certs", 1).Value.Int64()
		u = Uint64("certs", 1).Value.Uint64()
		f = Float64("certs", 1).Value.Float64()
		b = Bool("certs", true).Value.Bool()
		s = String("certs", "foo").Value.String()
		d = Duration("certs", d).Value.Duration()
		x = Any("certs", p).Value.Any()
	}))
	if a != 0 {
		t.Errorf("got %d allocs, want zero", a)
	}
	_ = u
	_ = f
	_ = b
	_ = s
	_ = x
}
