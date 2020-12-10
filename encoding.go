// Package base2048 implements base2048 encoding of binary data
package base2048

const (
	bitsPerChar  = 11
	bitsPerByte  = 8
	encoderSize  = 2048
	trailingSize = 8
)

// Encoding is a radix 2048 encoding/decoding scheme, defined by
// a 2048 unicode characters and a trailing 8 unicode characters.
// It has no standard (RFC, etc...) specifications.
type Encoding struct {
	encode    [encoderSize]rune
	decodeMap map[rune]uint16
	tail      [trailingSize]rune
	tailMap   map[rune]uint16
}

// NewEncoding returns a new Encoding defined by the given unicode characters,
// which should be a 2048-characters slice for encoder and a 8-characters slice
// for trailing.
func NewEncoding(encoder []rune, trailing []rune) *Encoding {
	if len(encoder) != encoderSize {
		panic("encoder is not 2048 characters")
	}

	if len(trailing) != trailingSize {
		panic("trailing is not 8 characters")
	}

	for i := 0; i < len(encoder); i++ {
		if encoder[i] == '\n' || encoder[i] == '\r' {
			panic("encoder contains newline character")
		}
	}

	for i := 0; i < len(trailing); i++ {
		if trailing[i] == '\n' || trailing[i] == '\r' {
			panic("trailing contains newline character")
		}
	}

	enc := new(Encoding)
	copy(enc.encode[:], encoder)
	copy(enc.tail[:], trailing)
	enc.decodeMap = make(map[rune]uint16, len(encoder))
	enc.tailMap = make(map[rune]uint16, len(trailing))

	for i := 0; i < len(encoder); i++ {
		enc.decodeMap[encoder[i]] = uint16(i)
	}

	for i := 0; i < len(trailing); i++ {
		enc.tailMap[trailing[i]] = uint16(i)
	}

	return enc
}

// DefaultEncoding is the default base2048 encoding defined in this module.
var DefaultEncoding = NewEncoding(DefaultEncodeChars, DefaultTrailingChars) //nolint:gochecknoglobals

// Encode encodes src using the encoding enc, writing
// EncodedLen(len(src)) characters to dst.
func (enc *Encoding) Encode(dst []rune, src []byte) {
	if len(src) == 0 {
		return
	}

	// enc is a pointer receiver, so the use of enc.encode within the hot
	// loop below means a nil check at every operation. Lift that nil check
	// outside of the loop to speed up the encoder.
	_ = enc.encode

	var (
		stage     uint16
		remaining uint8
		di        int
	)

	se := len(src)
	for si := 0; si < se; si++ {
		b := uint16(src[si])

		need := bitsPerChar - remaining
		if need <= bitsPerByte {
			remaining = bitsPerByte - need
			index := (stage << need) | (b >> remaining)
			stage = b & ((1 << remaining) - 1)
			dst[di] = enc.encode[index]
			di++
		} else {
			remaining += bitsPerByte
			stage = (stage << bitsPerByte) | b
		}
	}

	if remaining == 0 {
		return
	}

	// Add the remaining small block
	if remaining <= (bitsPerChar - bitsPerByte) {
		dst[di] = enc.tail[stage]
	} else {
		dst[di] = enc.encode[stage]
	}
}

// EncodeToString returns the base2048 encoding of src.
func (enc *Encoding) EncodeToString(src []byte) string {
	buf := make([]rune, enc.EncodedLen(len(src)))
	enc.Encode(buf, src)

	return string(buf)
}

// EncodedLen returns the length in characters of the base2048 encoding
// of an input buffer of bytes length n.
func (enc *Encoding) EncodedLen(n int) int {
	return (n*bitsPerByte + bitsPerChar - 1) / bitsPerChar
}

// Decode decodes src using the encoding enc. It writes at most
// DecodedLen(len(src)) bytes to dst and returns the number of bytes
// written. If src contains invalid base2048 data, it will return
// the number of bytes successfully written and CorruptInputError.
// New line characters (\r and \n) are ignored.
func (enc *Encoding) Decode(dst []byte, src []rune) (n int, err error) {
	if len(src) == 0 {
		return 0, nil
	}

	// Lift the nil check outside of the loop. enc.decodeMap is directly
	// used later in this function, to let the compiler know that the
	// receiver can't be nil.
	_ = enc.decodeMap

	var (
		stage     uint32
		remaining uint8
		residue   uint8
	)

	// Truncate trailing newline characters
	se := len(src) - 1
	for se >= 0 && (src[se] == '\r' || src[se] == '\n') {
		se--
	}

	for si := 0; si <= se; si++ {
		if src[si] == '\r' || src[si] == '\n' {
			continue
		}

		residue = (residue + bitsPerChar) % bitsPerByte

		var (
			newBits      uint16
			newBitsCount uint8
			ok           bool
		)

		if newBits, ok = enc.decodeMap[src[si]]; ok {
			if si == se {
				newBitsCount = bitsPerChar - residue
			} else {
				newBitsCount = bitsPerChar
			}
		} else {
			newBitsCount = bitsPerByte - remaining
			newBits, ok = enc.tailMap[src[si]]
			if !ok || si < se || newBits >= (1<<newBitsCount) {
				return n, CorruptInputError(si)
			}
		}

		stage = (stage << newBitsCount) | uint32(newBits)
		remaining += newBitsCount

		for remaining >= bitsPerByte {
			remaining -= bitsPerByte
			dst[n] = byte(stage >> remaining)
			stage &= (1 << remaining) - 1
			n++
		}
	}

	return n, nil
}

// DecodeString returns the bytes represented by the base64 string s.
func (enc *Encoding) DecodeString(s string) ([]byte, error) {
	sbuf := []rune(s)
	dbuf := make([]byte, enc.DecodedLen(len(sbuf)))
	n, err := enc.Decode(dbuf, sbuf)

	return dbuf[:n], err
}

// DecodedLen returns the maximum length in bytes of the decoded data
// corresponding to n characters of base2048-encoded data.
func (enc *Encoding) DecodedLen(n int) int {
	return n * bitsPerChar / bitsPerByte
}
