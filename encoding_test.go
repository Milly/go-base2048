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

func testEqual(t *testing.T, msg string, args ...interface{}) {
	t.Helper()

	if args[len(args)-2] != args[len(args)-1] {
		t.Errorf(msg, args...)
	}
}

func testRange(t *testing.T, msg string, args ...interface{}) {
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
	}
}

func testPanic(t *testing.T, proc func(), msg string, args ...interface{}) {
	t.Helper()

	defer func() {
		err := recover()
		if err != args[len(args)-1] {
			msg = fmt.Sprintf(msg, args...)
			t.Errorf("panic = %v, %s", err, msg)
		}
	}()

	proc()
}

func TestNewEncodingWithInvalidEncoderLength(t *testing.T) {
	encoder := make([]rune, 2047)
	copy(encoder, DefaultEncodeChars[:2047])
	testPanic(t, func() {
		NewEncoding(encoder, DefaultTrailingChars)
	}, "NewEncoding() = panic want %q", "encoder is not 2048 characters")
}

func TestNewEncodingWithInvalidTrailingLength(t *testing.T) {
	trailing := make([]rune, 7)
	copy(trailing, DefaultTrailingChars[:7])
	testPanic(t, func() {
		NewEncoding(DefaultEncodeChars, trailing)
	}, "NewEncoding() = panic want %q", "trailing is not 8 characters")
}

func TestNewEncodingWithEncoderContainsCLRF(t *testing.T) {
	testsets := []struct {
		pos   int
		value rune
	}{
		{0, '\r'},
		{0, '\n'},
		{1, '\r'},
		{1, '\n'},
		{2047, '\r'},
		{2047, '\n'},
	}
	encoder := make([]rune, 2048)

	for _, p := range testsets {
		copy(encoder, DefaultEncodeChars)
		encoder[p.pos] = p.value

		testPanic(t, func() {
			NewEncoding(encoder, DefaultTrailingChars)
		}, "NewEncoding() = panic want %q", "encoder contains newline character")
	}
}

func TestNewEncodingWithTrailingContainsCLRF(t *testing.T) {
	testsets := []struct {
		pos   int
		value rune
	}{
		{0, '\r'},
		{0, '\n'},
		{1, '\r'},
		{1, '\n'},
		{7, '\r'},
		{7, '\n'},
	}
	trailing := make([]rune, 8)

	for _, p := range testsets {
		copy(trailing, DefaultTrailingChars)
		trailing[p.pos] = p.value

		testPanic(t, func() {
			NewEncoding(DefaultEncodeChars, trailing)
		}, "NewEncoding() = panic want %q", "trailing contains newline character")
	}
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

func TestDecodeWithCRLF(t *testing.T) {
	testsets := []struct {
		decoded, encoded string
	}{
		{"foo", "\r\xD5\x93\xDA\x9D\xE0\xBC\x90"},
		{"foo", "\xD5\x93\n\xDA\x9D\xE0\xBC\x90"},
		{"foo", "\xD5\x93\xDA\x9D\xE0\xBC\x90\r"},
		{"foo", "\xD5\x93\xDA\x9D\xE0\xBC\x90\n"},
		{"foo", "\xD5\x93\xDA\x9D\xE0\xBC\x90\r\n"},
	}
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

func TestDecodeError(t *testing.T) {
	testerrors := []struct {
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

func TestEncodeDecodeMatches(t *testing.T) {
	enc := DefaultEncoding
	testsets := [][]byte{
		{},
		{0},
		{0, 0},
		{0xff},
		{0xff, 0xff},

		{0x80, 0}, // 10000000 00000000
		{0x40, 0}, // 01000000 00000000
		{0x20, 0}, // 00100000 00000000
		{0x10, 0}, // 00010000 00000000
		{0x08, 0}, // 00001000 00000000
		{0x04, 0}, // 00000100 00000000
		{0x02, 0}, // 00000010 00000000
		{0x01, 0}, // 00000001 00000000
		{0, 0x80}, // 00000000 10000000
		{0, 0x40}, // 00000000 01000000
		{0, 0x20}, // 00000000 00100000
		{0, 0x10}, // 00000000 00010000
		{0, 0x08}, // 00000000 00001000
		{0, 0x04}, // 00000000 00000100
		{0, 0x02}, // 00000000 00000010
		{0, 0x01}, // 00000000 00000001

		{0x7f, 0xff}, // 01111111 11111111
		{0xbf, 0xff}, // 10111111 11111111
		{0xdf, 0xff}, // 11011111 11111111
		{0xef, 0xff}, // 11101111 11111111
		{0xf7, 0xff}, // 11110111 11111111
		{0xfb, 0xff}, // 11111011 11111111
		{0xfd, 0xff}, // 11111101 11111111
		{0xfe, 0xff}, // 11111110 11111111
		{0xff, 0x7f}, // 11111111 01111111
		{0xff, 0xbf}, // 11111111 10111111
		{0xff, 0xdf}, // 11111111 11011111
		{0xff, 0xef}, // 11111111 11101111
		{0xff, 0xf7}, // 11111111 11110111
		{0xff, 0xfb}, // 11111111 11111011
		{0xff, 0xfd}, // 11111111 11111101
		{0xff, 0xfe}, // 11111111 11111110

		// 7-octal = 56-bits = 5*(11-bits) + 1-bits
		{0, 0, 0, 0, 0, 0, 0}, // ...0
		{0, 0, 0, 0, 0, 0, 1}, // ...1

		// 3-octal = 24-bits = 2*(11-bits) + 2-bits
		{0, 0, 0}, // ...00
		{0, 0, 1}, // ...01
		{0, 0, 2}, // ...11
		{0, 0, 3}, // ...10

		// 10-octal = 80-bits = 7*(11-bits) + 3-bits
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // ...000
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, // ...001
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 2}, // ...010
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 3}, // ...011
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 4}, // ...100
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 5}, // ...101
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 6}, // ...110
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 7}, // ...111
	}

	for _, in := range testsets {
		encoded := enc.EncodeToString(in)
		decoded, _ := enc.DecodeString(encoded)
		testEqual(t, "Decode(Encode()) = [% X], want [% X]", string(decoded), string(in))
	}
}
