// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package icmp provides basic functions for the manipulation of
// messages used in the Internet Control Message Protocols,
// ICMPv4 and ICMPv6.
//
// ICMPv4 and ICMPv6 are defined in RFC 792 and RFC 4443.
// Multi-part message support for ICMP is defined in RFC 4884.
// ICMP extensions for MPLS are defined in RFC 4950.
// ICMP extensions for interface and next-hop identification are
// defined in RFC 5837.
// PROBE: A utility for probing interfaces is defined in RFC 8335.
package icmp // import "golang.org/x/net/icmp"

import (
	"encoding/binary"
	"errors"
	"net"
	"runtime"

	"golang.org/x/net/internal/iana"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// BUG(mikio): This package is not implemented on JS, NaCl and Plan 9.

var (
	errInvalidConn      = errors.New("invalid connection")
	errInvalidProtocol  = errors.New("invalid protocol")
	errMessageTooShort  = errors.New("message too short")
	errHeaderTooShort   = errors.New("header too short")
	errBufferTooShort   = errors.New("buffer too short")
	errInvalidBody      = errors.New("invalid body")
	errNoExtension      = errors.New("no extension")
	errInvalidExtension = errors.New("invalid extension")
	errNotImplemented   = errors.New("not implemented on " + runtime.GOOS + "/" + runtime.GOARCH)
)

func checksum(b []byte) uint16 {
	csumcv := len(b) - 1 // checksum coverage
	s := uint32(0)
	for i := 0; i < csumcv; i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	if csumcv&1 == 0 {
		s += uint32(b[csumcv])
	}
	s = s>>16 + s&0xffff
	s = s + s>>16
	return ^uint16(s)
}

// A Type represents an ICMP message type.
type Type interface {
	Protocol() int
}

// A Message represents an ICMP message.
type Message struct {
	Type     Type        // type, either ipv4.ICMPType or ipv6.ICMPType
	Code     int         // code
	Checksum int         // checksum
	Body     MessageBody // body
}

// Marshal returns the binary encoding of the ICMP message m.
//
// For an ICMPv4 message, the returned message always contains the
// calculated checksum field.
//
// For an ICMPv6 message, the returned message contains the calculated
// checksum field when psh is not nil, otherwise the kernel will
// compute the checksum field during the message transmission.
// When psh is not nil, it must be the pseudo header for IPv6.
func (m *Message) Marshal(psh []byte) ([]byte, error) {
	var mtype byte
	switch typ := m.Type.(type) {
	case ipv4.ICMPType:
		mtype = byte(typ)
	case ipv6.ICMPType:
		mtype = byte(typ)
	default:
		return nil, errInvalidProtocol
	}
	b := []byte{mtype, byte(m.Code), 0, 0}
	proto := m.Type.Protocol()
	if proto == iana.ProtocolIPv6ICMP && psh != nil {
		b = append(psh, b...)
	}
	if m.Body != nil && m.Body.Len(proto) != 0 {
		mb, err := m.Body.Marshal(proto)
		if err != nil {
			return nil, err
		}
		b = append(b, mb...)
	}
	if proto == iana.ProtocolIPv6ICMP {
		if psh == nil { // cannot calculate checksum here
			return b, nil
		}
		off, l := 2*net.IPv6len, len(b)-len(psh)
		binary.BigEndian.PutUint32(b[off:off+4], uint32(l))
	}
	s := checksum(b)
	// Place checksum back in header; using ^= avoids the
	// assumption the checksum bytes are zero.
	b[len(psh)+2] ^= byte(s)
	b[len(psh)+3] ^= byte(s >> 8)
	return b[len(psh):], nil
}

var parseFns = map[Type]func(int, Type, []byte) (MessageBody, error){
	ipv4.ICMPTypeDestinationUnreachable: parseDstUnreach,
	ipv4.ICMPTypeTimeExceeded:           parseTimeExceeded,
	ipv4.ICMPTypeParameterProblem:       parseParamProb,

	ipv4.ICMPTypeEcho:                parseEcho,
	ipv4.ICMPTypeEchoReply:           parseEcho,
	ipv4.ICMPTypeExtendedEchoRequest: parseExtendedEchoRequest,
	ipv4.ICMPTypeExtendedEchoReply:   parseExtendedEchoReply,

	ipv6.ICMPTypeDestinationUnreachable: parseDstUnreach,
	ipv6.ICMPTypePacketTooBig:           parsePacketTooBig,
	ipv6.ICMPTypeTimeExceeded:           parseTimeExceeded,
	ipv6.ICMPTypeParameterProblem:       parseParamProb,

	ipv6.ICMPTypeEchoRequest:         parseEcho,
	ipv6.ICMPTypeEchoReply:           parseEcho,
	ipv6.ICMPTypeExtendedEchoRequest: parseExtendedEchoRequest,
	ipv6.ICMPTypeExtendedEchoReply:   parseExtendedEchoReply,
}

// ParseMessage parses b as an ICMP message.
// The provided grpc must be either the ICMPv4 or ICMPv6 protocol
// number.
func ParseMessage(proto int, b []byte) (*Message, error) {
	if len(b) < 4 {
		return nil, errMessageTooShort
	}
	var err error
	m := &Message{Code: int(b[1]), Checksum: int(binary.BigEndian.Uint16(b[2:4]))}
	switch proto {
	case iana.ProtocolICMP:
		m.Type = ipv4.ICMPType(b[0])
	case iana.ProtocolIPv6ICMP:
		m.Type = ipv6.ICMPType(b[0])
	default:
		return nil, errInvalidProtocol
	}
	if fn, ok := parseFns[m.Type]; !ok {
		m.Body, err = parseRawBody(proto, b[4:])
	} else {
		m.Body, err = fn(proto, m.Type, b[4:])
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}
