// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package twofish

import (
	"bytes"
	"testing"
)

var qbox = [2][4][16]byte{
	{
		{0x8, 0x1, 0x7, 0xD, 0x6, 0xF, 0x3, 0x2, 0x0, 0xB, 0x5, 0x9, 0xE, 0xC, 0xA, 0x4},
		{0xE, 0xC, 0xB, 0x8, 0x1, 0x2, 0x3, 0x5, 0xF, 0x4, 0xA, 0x6, 0x7, 0x0, 0x9, 0xD},
		{0xB, 0xA, 0x5, 0xE, 0x6, 0xD, 0x9, 0x0, 0xC, 0x8, 0xF, 0x3, 0x2, 0x4, 0x7, 0x1},
		{0xD, 0x7, 0xF, 0x4, 0x1, 0x2, 0x6, 0xE, 0x9, 0xB, 0x3, 0x0, 0x8, 0x5, 0xC, 0xA},
	},
	{
		{0x2, 0x8, 0xB, 0xD, 0xF, 0x7, 0x6, 0xE, 0x3, 0x1, 0x9, 0x4, 0x0, 0xA, 0xC, 0x5},
		{0x1, 0xE, 0x2, 0xB, 0x4, 0xC, 0x3, 0x7, 0x6, 0xD, 0xA, 0x5, 0xF, 0x9, 0x0, 0x8},
		{0x4, 0xC, 0x7, 0x5, 0x1, 0x6, 0x9, 0xA, 0x0, 0xE, 0xD, 0x8, 0x2, 0xB, 0x3, 0xF},
		{0xB, 0x9, 0x5, 0x1, 0xC, 0x3, 0xD, 0xE, 0x6, 0x4, 0x7, 0xF, 0x2, 0x0, 0x8, 0xA},
	},
}

// genSbox generates the variable sbox
func genSbox(qi int, x byte) byte {
	a0, b0 := x/16, x%16
	for i := 0; i < 2; i++ {
		a1 := a0 ^ b0
		b1 := (a0 ^ ((b0 << 3) | (b0 >> 1)) ^ (a0 << 3)) & 15
		a0 = qbox[qi][2*i][a1]
		b0 = qbox[qi][2*i+1][b1]
	}
	return (b0 << 4) + a0
}

func TestSbox(t *testing.T) {
	for n := range sbox {
		for m := range sbox[n] {
			if genSbox(n, byte(m)) != sbox[n][m] {
				t.Errorf("#%d|%d: sbox value = %d want %d", n, m, sbox[n][m], genSbox(n, byte(m)))
			}
		}
	}
}

var testVectors = []struct {
	key []byte
	dec []byte
	enc []byte
}{
	// These tests are extracted from LibTom
	{
		[]byte{0x9F, 0x58, 0x9F, 0x5C, 0xF6, 0x12, 0x2C, 0x32, 0xB6, 0xBF, 0xEC, 0x2F, 0x2A, 0xE8, 0xC3, 0x5A},
		[]byte{0xD4, 0x91, 0xDB, 0x16, 0xE7, 0xB1, 0xC3, 0x9E, 0x86, 0xCB, 0x08, 0x6B, 0x78, 0x9F, 0x54, 0x19},
		[]byte{0x01, 0x9F, 0x98, 0x09, 0xDE, 0x17, 0x11, 0x85, 0x8F, 0xAA, 0xC3, 0xA3, 0xBA, 0x20, 0xFB, 0xC3},
	},
	{
		[]byte{0x88, 0xB2, 0xB2, 0x70, 0x6B, 0x10, 0x5E, 0x36, 0xB4, 0x46, 0xBB, 0x6D, 0x73, 0x1A, 0x1E, 0x88,
			0xEF, 0xA7, 0x1F, 0x78, 0x89, 0x65, 0xBD, 0x44},
		[]byte{0x39, 0xDA, 0x69, 0xD6, 0xBA, 0x49, 0x97, 0xD5, 0x85, 0xB6, 0xDC, 0x07, 0x3C, 0xA3, 0x41, 0xB2},
		[]byte{0x18, 0x2B, 0x02, 0xD8, 0x14, 0x97, 0xEA, 0x45, 0xF9, 0xDA, 0xAC, 0xDC, 0x29, 0x19, 0x3A, 0x65},
	},
	{
		[]byte{0xD4, 0x3B, 0xB7, 0x55, 0x6E, 0xA3, 0x2E, 0x46, 0xF2, 0xA2, 0x82, 0xB7, 0xD4, 0x5B, 0x4E, 0x0D,
			0x57, 0xFF, 0x73, 0x9D, 0x4D, 0xC9, 0x2C, 0x1B, 0xD7, 0xFC, 0x01, 0x70, 0x0C, 0xC8, 0x21, 0x6F},
		[]byte{0x90, 0xAF, 0xE9, 0x1B, 0xB2, 0x88, 0x54, 0x4F, 0x2C, 0x32, 0xDC, 0x23, 0x9B, 0x26, 0x35, 0xE6},
		[]byte{0x6C, 0xB4, 0x56, 0x1C, 0x40, 0xBF, 0x0A, 0x97, 0x05, 0x93, 0x1C, 0xB6, 0xD4, 0x08, 0xE7, 0xFA},
	},
	// These tests are derived from https://www.schneier.com/code/ecb_ival.txt
	{
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		[]byte{0x9F, 0x58, 0x9F, 0x5C, 0xF6, 0x12, 0x2C, 0x32, 0xB6, 0xBF, 0xEC, 0x2F, 0x2A, 0xE8, 0xC3, 0x5A},
	},
	{
		[]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF, 0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10,
			0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
		},
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		[]byte{0xCF, 0xD1, 0xD2, 0xE5, 0xA9, 0xBE, 0x9C, 0xDF, 0x50, 0x1F, 0x13, 0xB8, 0x92, 0xBD, 0x22, 0x48},
	},
	{
		[]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF, 0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10,
			0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
		},
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		[]byte{0x37, 0x52, 0x7B, 0xE0, 0x05, 0x23, 0x34, 0xB8, 0x9F, 0x0C, 0xFC, 0xCA, 0xE8, 0x7C, 0xFA, 0x20},
	},
}

func TestCipher(t *testing.T) {
	for n, tt := range testVectors {
		// Test if the plaintext (dec) is encrypts to the given
		// ciphertext (enc) using the given certs. Test also if enc can
		// be decrypted again into dec.
		c, err := NewCipher(tt.key)
		if err != nil {
			t.Errorf("#%d: NewCipher: %v", n, err)
			return
		}

		buf := make([]byte, 16)
		c.Encrypt(buf, tt.dec)
		if !bytes.Equal(buf, tt.enc) {
			t.Errorf("#%d: encrypt = %x want %x", n, buf, tt.enc)
		}
		c.Decrypt(buf, tt.enc)
		if !bytes.Equal(buf, tt.dec) {
			t.Errorf("#%d: decrypt = %x want %x", n, buf, tt.dec)
		}

		// Test that 16 zero bytes, encrypted 1000 times then decrypted
		// 1000 times results in zero bytes again.
		zero := make([]byte, 16)
		buf = make([]byte, 16)
		for i := 0; i < 1000; i++ {
			c.Encrypt(buf, buf)
		}
		for i := 0; i < 1000; i++ {
			c.Decrypt(buf, buf)
		}
		if !bytes.Equal(buf, zero) {
			t.Errorf("#%d: encrypt/decrypt 1000: have %x want %x", n, buf, zero)
		}
	}
}
