package base2048

import (
	"fmt"
	"reflect"
	"testing"
)

type testset struct {
	decoded, encoded string
}

var testsets = []testset{
	{"\x14\xfb\x9c\x03\xd9\x7e", "\xC7\xAF\xE0\xB8\x8B\xC5\x8F\xE0\xAC\x87\xC5\x96"},
	{"\x14\xfb\x9c\x03\xd9", "\xC7\xAF\xE0\xB8\x8B\xC5\x8F\xC6\xA1"},
	{"\x14\xfb\x9c\x03", "\xC7\xAF\xE0\xB8\x8B\xC5\x8B"},

	{"", ""},
	{"f", "\xC6\xAE"},
	{"fo", "\xD5\x93\xC5\x97"},
	{"foo", "\xD5\x93\xDA\x9D\xE0\xBC\x90"},
	{"foob", "\xD5\x93\xDA\x9D\xD7\x93"},
	{"fooba", "\xD5\x93\xDA\x9D\xE0\xB6\xAA\xC6\xA9"},
	{"foobar", "\xD5\x93\xDA\x9D\xE0\xB6\xAA\xE0\xB0\xA8\xC5\x8A"},
	{"foobarb", "\xD5\x93\xDA\x9D\xE0\xB6\xAA\xE0\xB0\xA8\xC9\xB9\xE0\xBC\x8D"},
	{"foobarba", "\xD5\x93\xDA\x9D\xE0\xB6\xAA\xE0\xB0\xA8\xC9\xB9\xC6\xA9"},
	{"foobarbaz", "\xD5\x93\xDA\x9D\xE0\xB6\xAA\xE0\xB0\xA8\xC9\xB9\xCE\x9C\xC6\x82"},
	{"foobarbazq", "\xD5\x93\xDA\x9D\xE0\xB6\xAA\xE0\xB0\xA8\xC9\xB9\xCE\x9C\xE0\xBD\x85\xE0\xBC\x8E"},
	{"foobarbazqu", "\xD5\x93\xDA\x9D\xE0\xB6\xAA\xE0\xB0\xA8\xC9\xB9\xCE\x9C\xE0\xBD\x85\xCE\x8A"},
	{"foobarbazqux", "\xD5\x93\xDA\x9D\xE0\xB6\xAA\xE0\xB0\xA8\xC9\xB9\xCE\x9C\xE0\xBD\x85\xCE\x8A\xC7\x80"},
}

type testlen struct {
	decoded, encoded int
}

var testlens = []testlen{
	{0, 0},
	{1, 1},
	{2, 2},
	{3, 3},
	{4, 3},
	{5, 4},
	{6, 5},
	{7, 6},
	{8, 6},
}

func testEqual(t *testing.T, msg string, args ...interface{}) bool {
	t.Helper()
	if args[len(args)-2] != args[len(args)-1] {
		t.Errorf(msg, args...)
		return false
	}
	return true
}

func testRange(t *testing.T, msg string, args ...interface{}) bool {
	t.Helper()
	actual, ok1 := args[len(args)-3].(int)
	expectLower, ok2 := args[len(args)-2].(int)
	expectUpper, ok3 := args[len(args)-1].(int)
	if !(ok1 && ok2 && ok3) {
		msg := fmt.Sprintf(msg, args...)
		t.Errorf("type mismatch: testRange(%q)", msg)
	}

	if actual < expectLower || actual > expectUpper {
		t.Errorf(msg, args...)
		return false
	}
	return true
}

func TestEncode(t *testing.T) {
	enc := DefaultEncoding
	for _, p := range testsets {
		dbuf := make([]rune, enc.EncodedLen(len(p.decoded)))
		enc.Encode(dbuf, []byte(p.decoded))
		testEqual(t, "Encode(%q) = [% X], want [% X]", p.decoded, string(dbuf), p.encoded)
	}
}

func TestEncodeString(t *testing.T) {
	enc := DefaultEncoding
	for _, p := range testsets {
		got := enc.EncodeToString([]byte(p.decoded))
		testEqual(t, "Encode(%q) = [% X], want [% X]", p.decoded, got, p.encoded)
	}
}

func TestEncodedLen(t *testing.T) {
	enc := DefaultEncoding
	for _, p := range testlens {
		got := enc.EncodedLen(p.decoded)
		testEqual(t, "EncodedLen(%d) = %d, want %d", p.decoded, got, p.encoded)
	}
}

func TestDecode(t *testing.T) {
	enc := DefaultEncoding
	for _, p := range testsets {
		dbuf := make([]byte, enc.DecodedLen(len(p.encoded)))
		count, err := enc.Decode(dbuf, []rune(p.encoded))
		testEqual(t, "Decode([% X]) = error %v, want %v", p.encoded, err, error(nil))
		testEqual(t, "Decode([% X]) = length %v, want %v", p.encoded, count, len(p.decoded))
		testEqual(t, "Decode([% X]) = %q, want %q", p.encoded, string(dbuf[0:count]), p.decoded)
	}
}

func TestDecodeString(t *testing.T) {
	enc := DefaultEncoding
	for _, p := range testsets {
		dbuf, err := enc.DecodeString(p.encoded)
		testEqual(t, "DecodeString([% X]) = error %v, want %v", p.encoded, err, error(nil))
		testEqual(t, "DecodeString([% X]) = %q, want %q", p.encoded, string(dbuf), p.decoded)
	}
}

func TestDecodedLen(t *testing.T) {
	enc := DefaultEncoding
	for _, p := range testlens {
		got := enc.DecodedLen(p.encoded)
		// DecodedLen may return a length one greater
		testRange(t, "DecodedLen(%d) = %d, want %d between %d", p.encoded, got, p.decoded, p.decoded+1)
	}
}

func TestCorruptInputError(t *testing.T) {
	var testerrors = []struct {
		decoded, encoded string
		pos              int64
	}{
		// illegal character
		{"", "Z", 0},
		// illegal character at the first
		{"", "Z\xD5\x93\xDA\x9D\xE0\xB6\xAA\xE0\xB0\xA8\xC5\x8A", 0},
		// illegal character in the middle
		{"f", "\xD5\x93Z\xDA\x9D\xE0\xB6\xAA\xE0\xB0\xA8\xC5\x8A", 1},
		{"fooba", "\xD5\x93\xDA\x9D\xE0\xB6\xAA\xE0\xB0\xA8Z\xC5\x8A", 4},
		// illegal character at the last
		{"fo", "\xD5\x93\xDA\x9DZ", 2},
		// trailing character (\xE0\xBC\x90) in the middle
		{"fo", "\xD5\x93\xDA\x9D\xE0\xBC\x90\xD5\x93", 2},
		// illegal trailing character (\xE0\xBC\x91) at the last
		{"fo", "\xD5\x93\xDA\x9D\xE0\xBC\x91", 2},
	}
	enc := DefaultEncoding
	for _, p := range testerrors {
		dbuf, err := enc.DecodeString(p.encoded)
		want := CorruptInputError(p.pos)
		if !reflect.DeepEqual(want, err) {
			t.Errorf("DecodeString([% X]) = error %v, want %v", p.encoded, err, want)
		}
		if string(dbuf) != p.decoded {
			t.Errorf("DecodeString([% X]) = %q, want %q", p.encoded, dbuf, p.decoded)
		}
	}
}
