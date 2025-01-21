package strconv

import (
	"strconv"
	"testing"
)

func TestStrconv16(t *testing.T) {
	got := FormatUint16(8080)
	want := "8080"
	if got != want {
		t.Errorf("ParseUint16(8080) = %q; want %q", got, want)
	}
}

func TestStrconv16_2(t *testing.T) {
	got := FormatUint16(65535)
	want := "65535"
	if got != want {
		t.Errorf("ParseUint16(8080) = %q; want %q", got, want)
	}
}

func TestStrconv16_3(t *testing.T) {
	got := FormatUint16(0)
	want := "0"
	if got != want {
		t.Errorf("ParseUint16(8080) = %q; want %q", got, want)
	}
}

func FormatUint16Default(u uint16) string {
	return strconv.FormatUint(uint64(u), 10)
}

func BenchmarkParseUint16_my(b *testing.B) {

	for i := 0; i < b.N; i++ {
		FormatUint16(8080)
	}
}

func BenchmarkFormatUint16_default(b *testing.B) {

	for i := 0; i < b.N; i++ {
		FormatUint16Default(8080)
	}
}
