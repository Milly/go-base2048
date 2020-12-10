package base2048

import (
	"testing"
)

func TestCorruptInputError(t *testing.T) {
	{
		err := CorruptInputError(0)
		got := err.Error()
		want := "illegal base2048 data at input 0"
		testEqual(t, "CorruptInputError(%d) = %q, want %q", 0, got, want)
	}

	{
		err := CorruptInputError(4239)
		got := err.Error()
		want := "illegal base2048 data at input 4239"
		testEqual(t, "CorruptInputError(%d) = %q, want %q", 0, got, want)
	}
}
