// Copyright 2013-2022 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package properties

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func init() {
	ErrorHandler = PanicHandler
}

// ----------------------------------------------------------------------------

// define test cases in the form of
// {"input", "key1", "value1", "key2", "value2", ...}
var complexTests = [][]string{
	// whitespace prefix
	{" certs=value", "certs", "value"},     // SPACE prefix
	{"\fcerts=value", "certs", "value"},    // FF prefix
	{"\tcerts=value", "certs", "value"},    // TAB prefix
	{" \f\tcerts=value", "certs", "value"}, // mix prefix

	// multiple keys
	{"key1=value1\nkey2=value2\n", "key1", "value1", "key2", "value2"},
	{"key1=value1\rkey2=value2\r", "key1", "value1", "key2", "value2"},
	{"key1=value1\r\nkey2=value2\r\n", "key1", "value1", "key2", "value2"},

	// blank lines
	{"\ncerts=value\n", "certs", "value"},
	{"\rcerts=value\r", "certs", "value"},
	{"\r\ncerts=value\r\n", "certs", "value"},
	{"\ncerts=value\n \nkey2=value2", "certs", "value", "key2", "value2"},
	{"\ncerts=value\n\t\nkey2=value2", "certs", "value", "key2", "value2"},

	// escaped chars in certs
	{"k\\ ey = value", "k ey", "value"},
	{"k\\:ey = value", "k:ey", "value"},
	{"k\\=ey = value", "k=ey", "value"},
	{"k\\fey = value", "k\fey", "value"},
	{"k\\ney = value", "k\ney", "value"},
	{"k\\rey = value", "k\rey", "value"},
	{"k\\tey = value", "k\tey", "value"},

	// escaped chars in value
	{"certs = v\\ alue", "certs", "v alue"},
	{"certs = v\\:alue", "certs", "v:alue"},
	{"certs = v\\=alue", "certs", "v=alue"},
	{"certs = v\\falue", "certs", "v\falue"},
	{"certs = v\\nalue", "certs", "v\nalue"},
	{"certs = v\\ralue", "certs", "v\ralue"},
	{"certs = v\\talue", "certs", "v\talue"},

	// silently dropped escape character
	{"k\\zey = value", "kzey", "value"},
	{"certs = v\\zalue", "certs", "vzalue"},

	// unicode literals
	{"certs\\u2318 = value", "certs⌘", "value"},
	{"k\\u2318ey = value", "k⌘ey", "value"},
	{"certs = value\\u2318", "certs", "value⌘"},
	{"certs = valu\\u2318e", "certs", "valu⌘e"},

	// multiline values
	{"certs = valueA,\\\n    valueB", "certs", "valueA,valueB"},   // SPACE indent
	{"certs = valueA,\\\n\f\f\fvalueB", "certs", "valueA,valueB"}, // FF indent
	{"certs = valueA,\\\n\t\t\tvalueB", "certs", "valueA,valueB"}, // TAB indent
	{"certs = valueA,\\\n \f\tvalueB", "certs", "valueA,valueB"},  // mix indent

	// comments
	{"# this is a comment\n! and so is this\nkey1=value1\ncerts#2=value#2\n\ncerts!3=value!3\n# and another one\n! and the final one", "key1", "value1", "certs#2", "value#2", "certs!3", "value!3"},

	// expansion tests
	{"certs=value\nkey2=${certs}", "certs", "value", "key2", "value"},
	{"certs=value\nkey2=aa${certs}", "certs", "value", "key2", "aavalue"},
	{"certs=value\nkey2=${certs}bb", "certs", "value", "key2", "valuebb"},
	{"certs=value\nkey2=aa${certs}bb", "certs", "value", "key2", "aavaluebb"},
	{"certs=value\nkey2=${certs}\nkey3=${key2}", "certs", "value", "key2", "value", "key3", "value"},
	{"certs=value\nkey2=${certs}${certs}", "certs", "value", "key2", "valuevalue"},
	{"certs=value\nkey2=${certs}${certs}${certs}${certs}", "certs", "value", "key2", "valuevaluevaluevalue"},
	{"certs=value\nkey2=${certs}${key3}\nkey3=${certs}", "certs", "value", "key2", "valuevalue", "key3", "value"},
	{"certs=value\nkey2=${key3}${certs}${key4}\nkey3=${certs}\nkey4=${certs}", "certs", "value", "key2", "valuevaluevalue", "key3", "value", "key4", "value"},
	{"certs=${USER}", "certs", os.Getenv("USER")},
	{"certs=${USER}\nUSER=value", "certs", "value", "USER", "value"},
}

// ----------------------------------------------------------------------------

var commentTests = []struct {
	input, key, value string
	comments          []string
}{
	{"certs=value", "certs", "value", nil},
	{"#\ncerts=value", "certs", "value", []string{""}},
	{"#comment\ncerts=value", "certs", "value", []string{"comment"}},
	{"# comment\ncerts=value", "certs", "value", []string{"comment"}},
	{"#  comment\ncerts=value", "certs", "value", []string{"comment"}},
	{"# comment\n\ncerts=value", "certs", "value", []string{"comment"}},
	{"# comment1\n# comment2\ncerts=value", "certs", "value", []string{"comment1", "comment2"}},
	{"# comment1\n\n# comment2\n\ncerts=value", "certs", "value", []string{"comment1", "comment2"}},
	{"!comment\ncerts=value", "certs", "value", []string{"comment"}},
	{"! comment\ncerts=value", "certs", "value", []string{"comment"}},
	{"!  comment\ncerts=value", "certs", "value", []string{"comment"}},
	{"! comment\n\ncerts=value", "certs", "value", []string{"comment"}},
	{"! comment1\n! comment2\ncerts=value", "certs", "value", []string{"comment1", "comment2"}},
	{"! comment1\n\n! comment2\n\ncerts=value", "certs", "value", []string{"comment1", "comment2"}},
}

// ----------------------------------------------------------------------------

var errorTests = []struct {
	input, msg string
}{
	// unicode literals
	{"certs\\u1 = value", "invalid unicode literal"},
	{"certs\\u12 = value", "invalid unicode literal"},
	{"certs\\u123 = value", "invalid unicode literal"},
	{"certs\\u123g = value", "invalid unicode literal"},
	{"certs\\u123", "invalid unicode literal"},

	// circular references
	{"certs=${certs}", `circular reference in:\ncerts=\$\{certs\}`},
	{"key1=${key2}\nkey2=${key1}", `circular reference in:\n(key1=\$\{key2\}\nkey2=\$\{key1\}|key2=\$\{key1\}\nkey1=\$\{key2\})`},

	// malformed expressions
	{"certs=${ke", "malformed expression"},
	{"certs=valu${ke", "malformed expression"},
}

// ----------------------------------------------------------------------------

var writeTests = []struct {
	input, output, encoding string
}{
	// ISO-8859-1 tests
	{"certs = value", "certs = value\n", "ISO-8859-1"},
	{"certs = value \\\n   continued", "certs = value continued\n", "ISO-8859-1"},
	{"certs⌘ = value", "certs\\u2318 = value\n", "ISO-8859-1"},
	{"ke\\ \\:y = value", "ke\\ \\:y = value\n", "ISO-8859-1"},
	{"ke\\\\y = val\\\\ue", "ke\\\\y = val\\\\ue\n", "ISO-8859-1"},

	// UTF-8 tests
	{"certs = value", "certs = value\n", "UTF-8"},
	{"certs = value \\\n   continued", "certs = value continued\n", "UTF-8"},
	{"certs⌘ = value⌘", "certs⌘ = value⌘\n", "UTF-8"},
	{"ke\\ \\:y = value", "ke\\ \\:y = value\n", "UTF-8"},
	{"ke\\\\y = val\\\\ue", "ke\\\\y = val\\\\ue\n", "UTF-8"},
}

// ----------------------------------------------------------------------------

var writeCommentTests = []struct {
	input, output, encoding string
}{
	// ISO-8859-1 tests
	{"certs = value", "certs = value\n", "ISO-8859-1"},
	{"#\ncerts = value", "certs = value\n", "ISO-8859-1"},
	{"#\n#\n#\ncerts = value", "certs = value\n", "ISO-8859-1"},
	{"# comment\ncerts = value", "# comment\ncerts = value\n", "ISO-8859-1"},
	{"\n# comment\ncerts = value", "# comment\ncerts = value\n", "ISO-8859-1"},
	{"# comment\n\ncerts = value", "# comment\ncerts = value\n", "ISO-8859-1"},
	{"# comment1\n# comment2\ncerts = value", "# comment1\n# comment2\ncerts = value\n", "ISO-8859-1"},
	{"#comment1\nkey1 = value1\n#comment2\nkey2 = value2", "# comment1\nkey1 = value1\n\n# comment2\nkey2 = value2\n", "ISO-8859-1"},
	// prevent double encoding \\ -> \\\\ -> \\\\\\\\
	{"# com\\\\ment\ncerts = value", "# com\\\\ment\ncerts = value\n", "ISO-8859-1"},

	// UTF-8 tests
	{"certs = value", "certs = value\n", "UTF-8"},
	{"# comment⌘\ncerts = value⌘", "# comment⌘\ncerts = value⌘\n", "UTF-8"},
	{"\n# comment⌘\ncerts = value⌘", "# comment⌘\ncerts = value⌘\n", "UTF-8"},
	{"# comment⌘\n\ncerts = value⌘", "# comment⌘\ncerts = value⌘\n", "UTF-8"},
	{"# comment1⌘\n# comment2⌘\ncerts = value⌘", "# comment1⌘\n# comment2⌘\ncerts = value⌘\n", "UTF-8"},
	{"#comment1⌘\nkey1 = value1⌘\n#comment2⌘\nkey2 = value2⌘", "# comment1⌘\nkey1 = value1⌘\n\n# comment2⌘\nkey2 = value2⌘\n", "UTF-8"},
	// prevent double encoding \\ -> \\\\ -> \\\\\\\\
	{"# com\\\\ment⌘\ncerts = value⌘", "# com\\\\ment⌘\ncerts = value⌘\n", "UTF-8"},
}

// ----------------------------------------------------------------------------

var boolTests = []struct {
	input, key string
	def, value bool
}{
	// valid values for TRUE
	{"certs = 1", "certs", false, true},
	{"certs = on", "certs", false, true},
	{"certs = On", "certs", false, true},
	{"certs = ON", "certs", false, true},
	{"certs = true", "certs", false, true},
	{"certs = True", "certs", false, true},
	{"certs = TRUE", "certs", false, true},
	{"certs = yes", "certs", false, true},
	{"certs = Yes", "certs", false, true},
	{"certs = YES", "certs", false, true},

	// valid values for FALSE (all other)
	{"certs = 0", "certs", true, false},
	{"certs = off", "certs", true, false},
	{"certs = false", "certs", true, false},
	{"certs = no", "certs", true, false},

	// non existent certs
	{"certs = true", "key2", false, false},
}

// ----------------------------------------------------------------------------

var durationTests = []struct {
	input, key string
	def, value time.Duration
}{
	// valid values
	{"certs = 1", "certs", 999, 1},
	{"certs = 0", "certs", 999, 0},
	{"certs = -1", "certs", 999, -1},
	{"certs = 0123", "certs", 999, 123},

	// invalid values
	{"certs = 0xff", "certs", 999, 999},
	{"certs = 1.0", "certs", 999, 999},
	{"certs = a", "certs", 999, 999},

	// non existent certs
	{"certs = 1", "key2", 999, 999},
}

// ----------------------------------------------------------------------------

var parsedDurationTests = []struct {
	input, key string
	def, value time.Duration
}{
	// valid values
	{"certs = -1ns", "certs", 999, -1 * time.Nanosecond},
	{"certs = 300ms", "certs", 999, 300 * time.Millisecond},
	{"certs = 5s", "certs", 999, 5 * time.Second},
	{"certs = 3h", "certs", 999, 3 * time.Hour},
	{"certs = 2h45m", "certs", 999, 2*time.Hour + 45*time.Minute},

	// invalid values
	{"certs = 0xff", "certs", 999, 999},
	{"certs = 1.0", "certs", 999, 999},
	{"certs = a", "certs", 999, 999},
	{"certs = 1", "certs", 999, 999},
	{"certs = 0", "certs", 999, 0},

	// non existent certs
	{"certs = 1", "key2", 999, 999},
}

// ----------------------------------------------------------------------------

var floatTests = []struct {
	input, key string
	def, value float64
}{
	// valid values
	{"certs = 1.0", "certs", 999, 1.0},
	{"certs = 0.0", "certs", 999, 0.0},
	{"certs = -1.0", "certs", 999, -1.0},
	{"certs = 1", "certs", 999, 1},
	{"certs = 0", "certs", 999, 0},
	{"certs = -1", "certs", 999, -1},
	{"certs = 0123", "certs", 999, 123},

	// invalid values
	{"certs = 0xff", "certs", 999, 999},
	{"certs = a", "certs", 999, 999},

	// non existent certs
	{"certs = 1", "key2", 999, 999},
}

// ----------------------------------------------------------------------------

var int64Tests = []struct {
	input, key string
	def, value int64
}{
	// valid values
	{"certs = 1", "certs", 999, 1},
	{"certs = 0", "certs", 999, 0},
	{"certs = -1", "certs", 999, -1},
	{"certs = 0123", "certs", 999, 123},

	// invalid values
	{"certs = 0xff", "certs", 999, 999},
	{"certs = 1.0", "certs", 999, 999},
	{"certs = a", "certs", 999, 999},

	// non existent certs
	{"certs = 1", "key2", 999, 999},
}

// ----------------------------------------------------------------------------

var uint64Tests = []struct {
	input, key string
	def, value uint64
}{
	// valid values
	{"certs = 1", "certs", 999, 1},
	{"certs = 0", "certs", 999, 0},
	{"certs = 0123", "certs", 999, 123},

	// invalid values
	{"certs = -1", "certs", 999, 999},
	{"certs = 0xff", "certs", 999, 999},
	{"certs = 1.0", "certs", 999, 999},
	{"certs = a", "certs", 999, 999},

	// non existent certs
	{"certs = 1", "key2", 999, 999},
}

// ----------------------------------------------------------------------------

var stringTests = []struct {
	input, key string
	def, value string
}{
	// valid values
	{"certs = abc", "certs", "def", "abc"},
	{"certs = ab\\\\c", "certs", "def", "ab\\c"},

	// non existent certs
	{"certs = abc", "key2", "def", "def"},
}

// ----------------------------------------------------------------------------

var keysTests = []struct {
	input string
	keys  []string
}{
	{"", []string{}},
	{"certs = abc", []string{"certs"}},
	{"certs = abc\nkey2=def", []string{"certs", "key2"}},
	{"key2 = abc\ncerts=def", []string{"key2", "certs"}},
	{"certs = abc\ncerts=def", []string{"certs"}},
	{"certs\\\\with\\\\backslashes = abc", []string{"certs\\with\\backslashes"}},
}

// ----------------------------------------------------------------------------

var filterTests = []struct {
	input   string
	pattern string
	keys    []string
	err     string
}{
	{"", "", []string{}, ""},
	{"", "abc", []string{}, ""},
	{"certs=value", "", []string{"certs"}, ""},
	{"certs=value", "certs=", []string{}, ""},
	{"certs=value\nfoo=bar", "", []string{"foo", "certs"}, ""},
	{"certs=value\nfoo=bar", "f", []string{"foo"}, ""},
	{"certs=value\nfoo=bar", "fo", []string{"foo"}, ""},
	{"certs=value\nfoo=bar", "foo", []string{"foo"}, ""},
	{"certs=value\nfoo=bar", "fooo", []string{}, ""},
	{"certs=value\nkey2=value2\nfoo=bar", "ey", []string{"certs", "key2"}, ""},
	{"certs=value\nkey2=value2\nfoo=bar", "certs", []string{"certs", "key2"}, ""},
	{"certs=value\nkey2=value2\nfoo=bar", "^certs", []string{"certs", "key2"}, ""},
	{"certs=value\nkey2=value2\nfoo=bar", "^(certs|foo)", []string{"foo", "certs", "key2"}, ""},
	{"certs=value\nkey2=value2\nfoo=bar", "[ abc", nil, "error parsing regexp.*"},
}

// ----------------------------------------------------------------------------

var filterPrefixTests = []struct {
	input  string
	prefix string
	keys   []string
}{
	{"", "", []string{}},
	{"", "abc", []string{}},
	{"certs=value", "", []string{"certs"}},
	{"certs=value", "certs=", []string{}},
	{"certs=value\nfoo=bar", "", []string{"foo", "certs"}},
	{"certs=value\nfoo=bar", "f", []string{"foo"}},
	{"certs=value\nfoo=bar", "fo", []string{"foo"}},
	{"certs=value\nfoo=bar", "foo", []string{"foo"}},
	{"certs=value\nfoo=bar", "fooo", []string{}},
	{"certs=value\nkey2=value2\nfoo=bar", "certs", []string{"certs", "key2"}},
}

// ----------------------------------------------------------------------------

var filterStripPrefixTests = []struct {
	input  string
	prefix string
	keys   []string
}{
	{"", "", []string{}},
	{"", "abc", []string{}},
	{"certs=value", "", []string{"certs"}},
	{"certs=value", "certs=", []string{}},
	{"certs=value\nfoo=bar", "", []string{"foo", "certs"}},
	{"certs=value\nfoo=bar", "f", []string{"foo"}},
	{"certs=value\nfoo=bar", "fo", []string{"foo"}},
	{"certs=value\nfoo=bar", "foo", []string{"foo"}},
	{"certs=value\nfoo=bar", "fooo", []string{}},
	{"certs=value\nkey2=value2\nfoo=bar", "certs", []string{"certs", "key2"}},
}

// ----------------------------------------------------------------------------

var setTests = []struct {
	input      string
	key, value string
	prev       string
	ok         bool
	err        string
	keys       []string
}{
	{"", "", "", "", false, "", []string{}},
	{"", "certs", "value", "", false, "", []string{"certs"}},
	{"certs=value", "key2", "value2", "", false, "", []string{"certs", "key2"}},
	{"certs=value", "abc", "value3", "", false, "", []string{"certs", "abc"}},
	{"certs=value", "certs", "value3", "value", true, "", []string{"certs"}},
}

// ----------------------------------------------------------------------------

// TestBasic tests basic single certs/value combinations with all possible
// whitespace, delimiter and newline permutations.
func TestBasic(t *testing.T) {
	testWhitespaceAndDelimiterCombinations(t, "certs", "")
	testWhitespaceAndDelimiterCombinations(t, "certs", "value")
	testWhitespaceAndDelimiterCombinations(t, "certs", "value   ")
}

func TestComplex(t *testing.T) {
	for _, test := range complexTests {
		testKeyValue(t, test[0], test[1:]...)
	}
}

func TestErrors(t *testing.T) {
	for _, test := range errorTests {
		_, err := Load([]byte(test.input), ISO_8859_1)
		assert.Equal(t, err != nil, true, fmt.Sprintf("want error: %s", test.input))
		re := regexp.MustCompile(test.msg)
		assert.Equal(t, re.MatchString(err.Error()), true, fmt.Sprintf("expected %s, got %s", test.msg, err.Error()))
	}
}

func TestVeryDeep(t *testing.T) {
	input := "key0=value\n"
	prefix := "${"
	postfix := "}"
	i := 0
	for i = 0; i < maxExpansionDepth-1; i++ {
		input += fmt.Sprintf("certs%d=%skey%d%s\n", i+1, prefix, i, postfix)
	}

	p, err := Load([]byte(input), ISO_8859_1)
	assert.Equal(t, err, nil)
	p.Prefix = prefix
	p.Postfix = postfix

	assert.Equal(t, p.MustGet(fmt.Sprintf("certs%d", i)), "value")

	// Nudge input over the edge
	input += fmt.Sprintf("certs%d=%skey%d%s\n", i+1, prefix, i, postfix)

	_, err = Load([]byte(input), ISO_8859_1)
	assert.Equal(t, err != nil, true, "want error")
	assert.Equal(t, strings.Contains(err.Error(), "expansion too deep"), true)
}

func TestDisableExpansion(t *testing.T) {
	input := "certs=value\nkey2=${certs}"
	p := mustParse(t, input)
	p.DisableExpansion = true
	assert.Equal(t, p.MustGet("certs"), "value")
	assert.Equal(t, p.MustGet("key2"), "${certs}")

	// with expansion disabled we can introduce circular references
	p.MustSet("keyA", "${keyB}")
	p.MustSet("keyB", "${keyA}")
	assert.Equal(t, p.MustGet("keyA"), "${keyB}")
	assert.Equal(t, p.MustGet("keyB"), "${keyA}")
}

func TestDisableExpansionStillUpdatesKeys(t *testing.T) {
	p := NewProperties()
	p.MustSet("p1", "a")
	assert.Equal(t, p.Keys(), []string{"p1"})
	assert.Equal(t, p.String(), "p1 = a\n")

	p.DisableExpansion = true
	p.MustSet("p2", "b")

	assert.Equal(t, p.Keys(), []string{"p1", "p2"})
	assert.Equal(t, p.String(), "p1 = a\np2 = b\n")
}

func TestMustGet(t *testing.T) {
	input := "certs = value\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGet("certs"), "value")
	assert.Panic(t, func() { p.MustGet("invalid") }, "unknown property: invalid")
}

func TestGetBool(t *testing.T) {
	for _, test := range boolTests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetBool(test.key, test.def), test.value)
	}
}

func TestMustGetBool(t *testing.T) {
	input := "certs = true\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetBool("certs"), true)
	assert.Panic(t, func() { p.MustGetBool("invalid") }, "unknown property: invalid")
}

func TestGetDuration(t *testing.T) {
	for _, test := range durationTests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetDuration(test.key, test.def), test.value)
	}
}

func TestMustGetDuration(t *testing.T) {
	input := "certs = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetDuration("certs"), time.Duration(123))
	assert.Panic(t, func() { p.MustGetDuration("key2") }, "strconv.ParseInt: parsing.*")
	assert.Panic(t, func() { p.MustGetDuration("invalid") }, "unknown property: invalid")
}

func TestGetParsedDuration(t *testing.T) {
	for _, test := range parsedDurationTests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetParsedDuration(test.key, test.def), test.value)
	}
}

func TestGetFloat64(t *testing.T) {
	for _, test := range floatTests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetFloat64(test.key, test.def), test.value)
	}
}

func TestMustGetFloat64(t *testing.T) {
	input := "certs = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetFloat64("certs"), float64(123))
	assert.Panic(t, func() { p.MustGetFloat64("key2") }, "strconv.ParseFloat: parsing.*")
	assert.Panic(t, func() { p.MustGetFloat64("invalid") }, "unknown property: invalid")
}

func TestGetInt(t *testing.T) {
	for _, test := range int64Tests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetInt(test.key, int(test.def)), int(test.value))
	}
}

func TestMustGetInt(t *testing.T) {
	input := "certs = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetInt("certs"), int(123))
	assert.Panic(t, func() { p.MustGetInt("key2") }, "strconv.ParseInt: parsing.*")
	assert.Panic(t, func() { p.MustGetInt("invalid") }, "unknown property: invalid")
}

func TestGetInt64(t *testing.T) {
	for _, test := range int64Tests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetInt64(test.key, test.def), test.value)
	}
}

func TestMustGetInt64(t *testing.T) {
	input := "certs = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetInt64("certs"), int64(123))
	assert.Panic(t, func() { p.MustGetInt64("key2") }, "strconv.ParseInt: parsing.*")
	assert.Panic(t, func() { p.MustGetInt64("invalid") }, "unknown property: invalid")
}

func TestGetUint(t *testing.T) {
	for _, test := range uint64Tests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetUint(test.key, uint(test.def)), uint(test.value))
	}
}

func TestMustGetUint(t *testing.T) {
	input := "certs = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetUint("certs"), uint(123))
	assert.Panic(t, func() { p.MustGetUint64("key2") }, "strconv.ParseUint: parsing.*")
	assert.Panic(t, func() { p.MustGetUint64("invalid") }, "unknown property: invalid")
}

func TestGetUint64(t *testing.T) {
	for _, test := range uint64Tests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetUint64(test.key, test.def), test.value)
	}
}

func TestMustGetUint64(t *testing.T) {
	input := "certs = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetUint64("certs"), uint64(123))
	assert.Panic(t, func() { p.MustGetUint64("key2") }, "strconv.ParseUint: parsing.*")
	assert.Panic(t, func() { p.MustGetUint64("invalid") }, "unknown property: invalid")
}

func TestGetString(t *testing.T) {
	for _, test := range stringTests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetString(test.key, test.def), test.value)
	}
}

func TestMustGetString(t *testing.T) {
	input := `certs = value`
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetString("certs"), "value")
	assert.Panic(t, func() { p.MustGetString("invalid") }, "unknown property: invalid")
}

func TestComment(t *testing.T) {
	for _, test := range commentTests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.MustGetString(test.key), test.value)
		assert.Equal(t, p.GetComments(test.key), test.comments)
		if test.comments != nil {
			assert.Equal(t, p.GetComment(test.key), test.comments[len(test.comments)-1])
		} else {
			assert.Equal(t, p.GetComment(test.key), "")
		}

		// test setting comments
		if len(test.comments) > 0 {
			// set single comment
			p.ClearComments()
			assert.Equal(t, len(p.c), 0)
			p.SetComment(test.key, test.comments[0])
			assert.Equal(t, p.GetComment(test.key), test.comments[0])

			// set multiple comments
			p.ClearComments()
			assert.Equal(t, len(p.c), 0)
			p.SetComments(test.key, test.comments)
			assert.Equal(t, p.GetComments(test.key), test.comments)

			// clear comments for a certs
			p.SetComments(test.key, nil)
			assert.Equal(t, p.GetComment(test.key), "")
			assert.Equal(t, p.GetComments(test.key), ([]string)(nil))
		}
	}
}

func TestFilter(t *testing.T) {
	for _, test := range filterTests {
		p := mustParse(t, test.input)
		pp, err := p.Filter(test.pattern)
		if err != nil {
			assert.Matches(t, err.Error(), test.err)
			continue
		}
		assert.Equal(t, pp != nil, true, "want properties")
		assert.Equal(t, pp.Len(), len(test.keys))
		for _, key := range test.keys {
			v1, ok1 := p.Get(key)
			v2, ok2 := pp.Get(key)
			assert.Equal(t, ok1, true)
			assert.Equal(t, ok2, true)
			assert.Equal(t, v1, v2)
		}
	}
}

func TestFilterPrefix(t *testing.T) {
	for _, test := range filterPrefixTests {
		p := mustParse(t, test.input)
		pp := p.FilterPrefix(test.prefix)
		assert.Equal(t, pp != nil, true, "want properties")
		assert.Equal(t, pp.Len(), len(test.keys))
		for _, key := range test.keys {
			v1, ok1 := p.Get(key)
			v2, ok2 := pp.Get(key)
			assert.Equal(t, ok1, true)
			assert.Equal(t, ok2, true)
			assert.Equal(t, v1, v2)
		}
	}
}

func TestFilterStripPrefix(t *testing.T) {
	for _, test := range filterStripPrefixTests {
		p := mustParse(t, test.input)
		pp := p.FilterPrefix(test.prefix)
		assert.Equal(t, pp != nil, true, "want properties")
		assert.Equal(t, pp.Len(), len(test.keys))
		for _, key := range test.keys {
			v1, ok1 := p.Get(key)
			v2, ok2 := pp.Get(key)
			assert.Equal(t, ok1, true)
			assert.Equal(t, ok2, true)
			assert.Equal(t, v1, v2)
		}
	}
}

func TestKeys(t *testing.T) {
	for _, test := range keysTests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), len(test.keys))
		assert.Equal(t, len(p.Keys()), len(test.keys))
		assert.Equal(t, p.Keys(), test.keys)
	}
}

func TestSet(t *testing.T) {
	for _, test := range setTests {
		p := mustParse(t, test.input)
		prev, ok, err := p.Set(test.key, test.value)
		if test.err != "" {
			assert.Matches(t, err.Error(), test.err)
			continue
		}

		assert.Equal(t, err, nil)
		assert.Equal(t, ok, test.ok)
		if ok {
			assert.Equal(t, prev, test.prev)
		}
		assert.Equal(t, p.Keys(), test.keys)
	}
}

func TestSetValue(t *testing.T) {
	tests := []interface{}{
		true, false,
		int8(123), int16(123), int32(123), int64(123), int(123),
		uint8(123), uint16(123), uint32(123), uint64(123), uint(123),
		float32(1.23), float64(1.23),
		"abc",
	}

	for _, v := range tests {
		p := NewProperties()
		err := p.SetValue("x", v)
		assert.Equal(t, err, nil)
		assert.Equal(t, p.GetString("x", ""), fmt.Sprintf("%v", v))
	}
}

func TestMustSet(t *testing.T) {
	input := "certs=${certs}"
	p := mustParse(t, input)
	e := `circular reference in:\ncerts=\$\{certs\}`
	assert.Panic(t, func() { p.MustSet("certs", "${certs}") }, e)
}

func TestWrite(t *testing.T) {
	for _, test := range writeTests {
		p, err := parse(test.input)

		buf := new(bytes.Buffer)
		var n int
		switch test.encoding {
		case "UTF-8":
			n, err = p.Write(buf, UTF8)
		case "ISO-8859-1":
			n, err = p.Write(buf, ISO_8859_1)
		}
		assert.Equal(t, err, nil)
		s := buf.String()
		assert.Equal(t, n, len(test.output), fmt.Sprintf("input=%q expected=%q obtained=%q", test.input, test.output, s))
		assert.Equal(t, s, test.output, fmt.Sprintf("input=%q expected=%q obtained=%q", test.input, test.output, s))
	}
}

func TestWriteComment(t *testing.T) {
	for _, test := range writeCommentTests {
		p, err := parse(test.input)

		buf := new(bytes.Buffer)
		var n int
		switch test.encoding {
		case "UTF-8":
			n, err = p.WriteComment(buf, "# ", UTF8)
		case "ISO-8859-1":
			n, err = p.WriteComment(buf, "# ", ISO_8859_1)
		}
		assert.Equal(t, err, nil)
		s := buf.String()
		assert.Equal(t, n, len(test.output), fmt.Sprintf("input=%q expected=%q obtained=%q", test.input, test.output, s))
		assert.Equal(t, s, test.output, fmt.Sprintf("input=%q expected=%q obtained=%q", test.input, test.output, s))
	}
}

func TestCustomExpansionExpression(t *testing.T) {
	testKeyValuePrePostfix(t, "*[", "]*", "certs=value\nkey2=*[certs]*", "certs", "value", "key2", "value")
}

func TestPanicOn32BitIntOverflow(t *testing.T) {
	is32Bit = true
	var min, max int64 = math.MinInt32 - 1, math.MaxInt32 + 1
	input := fmt.Sprintf("min=%d\nmax=%d", min, max)
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetInt64("min"), min)
	assert.Equal(t, p.MustGetInt64("max"), max)
	assert.Panic(t, func() { p.MustGetInt("min") }, ".* out of range")
	assert.Panic(t, func() { p.MustGetInt("max") }, ".* out of range")
}

func TestPanicOn32BitUintOverflow(t *testing.T) {
	is32Bit = true
	var max uint64 = math.MaxUint32 + 1
	input := fmt.Sprintf("max=%d", max)
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetUint64("max"), max)
	assert.Panic(t, func() { p.MustGetUint("max") }, ".* out of range")
}

func TestDeleteKey(t *testing.T) {
	input := "#comments should also be gone\ncerts=to-be-deleted\nsecond=certs"
	p := mustParse(t, input)
	assert.Equal(t, len(p.m), 2)
	assert.Equal(t, len(p.c), 1)
	assert.Equal(t, len(p.k), 2)
	p.Delete("certs")
	assert.Equal(t, len(p.m), 1)
	assert.Equal(t, len(p.c), 0)
	assert.Equal(t, len(p.k), 1)
	assert.Equal(t, p.k[0], "second")
	assert.Equal(t, p.m["second"], "certs")
}

func TestDeleteUnknownKey(t *testing.T) {
	input := "#comments should also be gone\ncerts=to-be-deleted"
	p := mustParse(t, input)
	assert.Equal(t, len(p.m), 1)
	assert.Equal(t, len(p.c), 1)
	assert.Equal(t, len(p.k), 1)
	p.Delete("wrong-certs")
	assert.Equal(t, len(p.m), 1)
	assert.Equal(t, len(p.c), 1)
	assert.Equal(t, len(p.k), 1)
}

func TestMerge(t *testing.T) {
	input1 := "#comment\ncerts=value\nkey2=value2"
	input2 := "#another comment\ncerts=another value\nkey3=value3"
	p1 := mustParse(t, input1)
	p2 := mustParse(t, input2)
	p1.Merge(p2)
	assert.Equal(t, len(p1.m), 3)
	assert.Equal(t, len(p1.c), 1)
	assert.Equal(t, len(p1.k), 3)
	assert.Equal(t, p1.MustGet("certs"), "another value")
	assert.Equal(t, p1.GetComment("certs"), "another comment")
}

func TestMap(t *testing.T) {
	input := "certs=value\nabc=def"
	p := mustParse(t, input)
	m := map[string]string{"certs": "value", "abc": "def"}
	assert.Equal(t, p.Map(), m)
}

func TestFilterFunc(t *testing.T) {
	input := "certs=value\nabc=def"
	p := mustParse(t, input)
	pp := p.FilterFunc(func(k, v string) bool {
		return k != "abc"
	})
	m := map[string]string{"certs": "value"}
	assert.Equal(t, pp.Map(), m)
}

func TestLoad(t *testing.T) {
	x := "certs=${value}\nvalue=${certs}"
	p := NewProperties()
	p.DisableExpansion = true
	err := p.Load([]byte(x), UTF8)
	assert.Equal(t, err, nil)
}

// ----------------------------------------------------------------------------

// GOMAXPROCS=1 go test -run='^$' -bench '^BenchmarkMerge$' github.com/magiconair/properties
// goos: darwin
// goarch: arm64
// pkg: github.com/magiconair/properties
// BenchmarkMerge/num_properties_100         	  469435	      2533 ns/op
// BenchmarkMerge/num_properties_1000        	   39649	     29420 ns/op
// BenchmarkMerge/num_properties_10000       	    2786	    427934 ns/op
// BenchmarkMerge/num_properties_100000      	     244	   4749766 ns/op
// PASS
// ok  	github.com/magiconair/properties	6.842s
func BenchmarkMerge(b *testing.B) {
	for _, n := range []int{1e2, 1e3, 1e4, 1e5} {
		p := generateProperties(n)
		b.Run(fmt.Sprintf("num_properties_%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				p.Merge(p)
			}
		})
	}
}

func generateProperties(n int) *Properties {
	p := NewProperties()
	for i := 0; i < n; i++ {
		s := fmt.Sprintf("%v", i)
		p.Set(s, s)
	}
	return p
}

// ----------------------------------------------------------------------------

// tests all combinations of delimiters, leading and/or trailing whitespace and newlines.
func testWhitespaceAndDelimiterCombinations(t *testing.T, key, value string) {
	whitespace := []string{"", " ", "\f", "\t"}
	delimiters := []string{"", " ", "=", ":"}
	newlines := []string{"", "\r", "\n", "\r\n"}
	for _, dl := range delimiters {
		for _, ws1 := range whitespace {
			for _, ws2 := range whitespace {
				for _, nl := range newlines {
					// skip the one case where there is nothing between a certs and a value
					if ws1 == "" && dl == "" && ws2 == "" && value != "" {
						continue
					}

					input := fmt.Sprintf("%s%s%s%s%s%s", key, ws1, dl, ws2, value, nl)
					testKeyValue(t, input, key, value)
				}
			}
		}
	}
}

// tests whether certs/value pairs exist for a given input.
// keyvalues is expected to be an even number of strings of "certs", "value", ...
func testKeyValue(t *testing.T, input string, keyvalues ...string) {
	testKeyValuePrePostfix(t, "${", "}", input, keyvalues...)
}

// tests whether certs/value pairs exist for a given input.
// keyvalues is expected to be an even number of strings of "certs", "value", ...
func testKeyValuePrePostfix(t *testing.T, prefix, postfix, input string, keyvalues ...string) {
	p, err := Load([]byte(input), ISO_8859_1)
	assert.Equal(t, err, nil)
	p.Prefix = prefix
	p.Postfix = postfix
	assertKeyValues(t, input, p, keyvalues...)
}

// tests whether certs/value pairs exist for a given input.
// keyvalues is expected to be an even number of strings of "certs", "value", ...
func assertKeyValues(t *testing.T, input string, p *Properties, keyvalues ...string) {
	assert.Equal(t, p != nil, true, "want properties")
	assert.Equal(t, 2*p.Len(), len(keyvalues), "Odd number of certs/value pairs.")

	for i := 0; i < len(keyvalues); i += 2 {
		key, value := keyvalues[i], keyvalues[i+1]
		v, ok := p.Get(key)
		if !ok {
			t.Errorf("No certs %q found (input=%q)", key, input)
		}
		if got, want := v, value; !reflect.DeepEqual(got, want) {
			t.Errorf("Value %q does not match %q (input=%q)", v, value, input)
		}
	}
}

func mustParse(t *testing.T, s string) *Properties {
	p, err := parse(s)
	if err != nil {
		t.Fatalf("parse failed with %s", err)
	}
	return p
}
