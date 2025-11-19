package strconv

import (
	"math"
	"strconv"
	"testing"
)

var sink string
var _sink string
var sinkUint64 uint64
var sinkInt64 int64
var sinkErr error

func BenchmarkULL2String(b *testing.B) {
	buf := make([]byte, SAFETY_BUF_SIZE)
	for i := 0; i < b.N; i++ {
		n := Ull2String(buf, 1234567890123456789)
		_sink = _string(buf[:n])
	}
}

func BenchmarkU64string(b *testing.B) {
	buf := make([]byte, SAFETY_BUF_SIZE)
	for i := 0; i < b.N; i++ {
		n := U64ToA(1234567890123456789, buf)
		_sink = _string(buf[:n])
	}
}

func BenchmarkLL2String(b *testing.B) {
	buf := make([]byte, SAFETY_BUF_SIZE)
	for i := 0; i < b.N; i++ {
		n := Ll2String(buf, -1234567890123456789)
		_sink = _string(buf[:n])
	}
}
func BenchmarkLtuString(b *testing.B) {
	buf := make([]byte, SAFETY_BUF_SIZE)
	for i := 0; i < b.N; i++ {
		n := I64ToA(-1234567890123456789, buf)
		_sink = _string(buf[:n])
	}
}

func BenchmarkStrconvUint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_sink = strconv.FormatUint(1234567890123456789, 10)
	}
}

func BenchmarkStrconvInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_sink = strconv.FormatInt(-1234567890123456789, 10)
	}
}

func BenchmarkParseUint64(b *testing.B) {
	s := "1234567890123456789"
	for i := 0; i < b.N; i++ {
		sinkUint64, sinkErr = ParseUint64(s)
	}
}

func BenchmarkParseInt64(b *testing.B) {
	s := "-1234567890123456789"
	for i := 0; i < b.N; i++ {
		sinkInt64, sinkErr = ParseInt64(s)
	}
}

func BenchmarkStrconvParseUint64(b *testing.B) {
	s := "1234567890123456789"
	for i := 0; i < b.N; i++ {
		sinkUint64, sinkErr = strconv.ParseUint(s, 10, 64)
	}
}

func BenchmarkStrconvParseInt64(b *testing.B) {
	s := "-1234567890123456789"
	for i := 0; i < b.N; i++ {
		sinkInt64, sinkErr = strconv.ParseInt(s, 10, 64)
	}
}

func TestParseUint64(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
		wantErr  bool
	}{
		{"0", 0, false},
		{"1", 1, false},
		{"1234567890", 1234567890, false},
		{"18446744073709551615", math.MaxUint64, false},
		{"18446744073709551616", 0, true},
		{"-1", 0, true},
		{"abc", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		got, err := ParseUint64(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseUint64(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.expected {
			t.Errorf("ParseUint64(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestParseInt64(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		wantErr  bool
	}{
		{"0", 0, false},
		{"1", 1, false},
		{"-1", -1, false},
		{"1234567890", 1234567890, false},
		{"-1234567890", -1234567890, false},
		{"9223372036854775807", math.MaxInt64, false},
		{"-9223372036854775808", math.MinInt64, false},
		{"9223372036854775808", 0, true},
		{"-9223372036854775809", 0, true},
		{"abc", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		got, err := ParseInt64(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseInt64(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.expected {
			t.Errorf("ParseInt64(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestUll2String(t *testing.T) {
	cases := []uint64{
		0, 1, 9, 10, 11, 99, 100, 9999,
		123456789, 1000000000000, 18446744073709551615,
	}

	for _, v := range cases {
		buf := make([]byte, 32)
		n := Ull2String(buf, v)
		got := _string(buf[:n])
		want := strconv.FormatUint(v, 10)
		if got != want {
			t.Fatalf("Ull2String(%d) = %q, want %q", v, got, want)
		}
	}
}

func BenchmarkU64ToA(b *testing.B) {
	buf := make([]byte, SAFETY_BUF_SIZE)
	var s string

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n := U64ToA(uint64(i), buf)
		s = _string(buf[:n])
	}
	sink = s
}

func TestU64ToA_Size(t *testing.T) {
	tests := []uint64{
		0,
		9,
		10,
		99,
		100,
		12345,
		9999999999999999999,
		math.MaxUint64,
	}

	buf := make([]byte, SAFETY_BUF_SIZE)

	for _, v := range tests {
		n := U64ToA(v, buf)
		if n == 0 {
			t.Fatalf("U64ToA returned 0 for %d", v)
		}
		if n > 20 {
			t.Fatalf("U64ToA produced too long output for %d: %d bytes", v, n)
		}
	}
}

func TestI64ToA_Size(t *testing.T) {
	tests := []int64{
		0,
		-1,
		1,
		9,
		10,
		-10,
		99,
		-99,
		100,
		-100,
		12345,
		-12345,
		math.MaxInt64,
		math.MinInt64,
	}

	buf := make([]byte, SAFETY_BUF_SIZE)

	for _, v := range tests {
		n := I64ToA(v, buf)
		if n == 0 {
			t.Fatalf("I64ToA returned 0 for %d", v)
		}
		if n > 21 {
			t.Fatalf("I64ToA produced too long output for %d: %d bytes", v, n)
		}
	}
}

func BenchmarkI64ToA(b *testing.B) {
	buf := make([]byte, SAFETY_BUF_SIZE)
	var s string

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n := I64ToA(int64(i), buf)
		s = _string(buf[:n])
	}
	sink = s
}

func BenchmarkStrconvU64(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		s = strconv.FormatUint(uint64(i), 10)
	}
	sink = s
}

func BenchmarkStrconvI64(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		s = strconv.Itoa(i)
	}
	sink = s
}

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

func BenchmarkFormatUint16_my(b *testing.B) {
	var u string

	for i := 0; i < b.N; i++ {
		u = FormatUint16(8080)
		_ = u
	}
}

func BenchmarkFormatUint16_default(b *testing.B) {
	var u string

	for i := 0; i < b.N; i++ {
		u = FormatUint16Default(8080)
		_ = u
	}
}
