package lowlevelfunctions

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	mathrand "math/rand"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/NikoMalik/low-level-functions/example"

	"github.com/NikoMalik/low-level-functions/union"
)

var block1kb = 1024
var data []byte

var zeroBuf = make([]byte, 4096)

func stdClear(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

func makeRandom(size int) []byte {
	src := make([]byte, size)
	r := mathrand.New(mathrand.NewSource(42))
	_, _ = r.Read(src)
	return src
}

var sk []byte

func bs(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func TestEncodeCorrectness(t *testing.T) {
	cases := [][]byte{
		{},
		{0x00},
		{0x00, 0xff, 0x12, 0x7a},
		makeRandom(64),
		makeRandom(1024),
	}

	for idx, src := range cases {
		dstStd := make([]byte, EncodedLen(len(src)))
		dstB := make([]byte, EncodedLen(len(src)))

		hex.Encode(dstStd, src)

		EncodeUnrolled8(dstB, src)

		if !bytes.Equal(dstStd, dstB) {
			t.Fatalf("case %d: EncodeUnrolled8 produced wrong result\nexpected=%s\ngot     =%s", idx, string(dstStd), string(dstB))
		}
	}
}

func TestEncodeAgainstStdlibManyRandom(t *testing.T) {
	for i := 0; i < 200; i++ {
		n := mathrand.Intn(2048)
		src := makeRandom(n)

		std := make([]byte, EncodedLen(len(src)))
		got := make([]byte, EncodedLen(len(src)))

		hex.Encode(std, src)
		EncodeUnrolled8(got, src)

		if !bytes.Equal(std, got) {
			t.Fatalf("random %d: EncodeUnrolled8 mismatch (len=%d)", i, n)
		}
	}
}

func BenchmarkEncodeVariants(b *testing.B) {
	sizes := []int{16, 256, 4096}

	for _, size := range sizes {
		src := makeRandom(size)
		dst := make([]byte, EncodedLen(len(src)))

		// b.Run(fmt.Sprintf("Encode/%d", size), func(b *testing.B) {
		// 	b.ReportAllocs()
		// 	b.ResetTimer()
		// 	for i := 0; i < b.N; i++ {
		// 		_ = Encode(dst, src)
		// 	}
		// })

		b.Run(fmt.Sprintf("Unrolled8/%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = EncodeUnrolled8(dst, src)
			}
		})

		b.Run(fmt.Sprintf("Stdlib/%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				hex.Encode(dst, src)
			}
		})
	}
}

func BenchmarkStdStringToBytes(b *testing.B) {
	s := "Hello, this is a sample string for benchmarking."
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sk = []byte(s)
	}
}

func BenchmarkUnsafeStringToBytes(b *testing.B) {
	s := "Hello, this is a sample string for benchmarking."
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sk = bs(s)
	}
}

func sumLoop(arr []int) int {
	sum := 0
	if len(arr) < 8 {
		return 0
	}
	for i := 0; i < len(arr); i++ {
		sum += arr[i]
	}
	return sum
}

var sinkint int

func BenchmarkSumLoop(b *testing.B) {
	arr := make([]int, 1<<20)
	for i := 0; i < len(arr); i++ {
		arr[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkint = sumLoop(arr)
	}
}

func BenchmarkSumUnroll8(b *testing.B) {
	arr := make([]int, 1<<20)
	for i := 0; i < len(arr); i++ {
		arr[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkint = SumUnroll8(arr)
	}
}

func TestSumLoop(t *testing.T) {
	arr := make([]int, 20)
	for i := range arr {
		arr[i] = i + 1
	}
	fmt.Println("sumLoop:   ", sumLoop(arr))
	fmt.Println("sumUnroll8:", SumUnroll8(arr))

	if sumLoop(arr) != SumUnroll8(arr) {
		t.Errorf("SumUnroll8(%v) = %d, want %d", arr, SumUnroll8(arr), sumLoop(arr))
	}
}

var sink byte
var sink2 int64

func BenchmarkMemclr(b *testing.B) {
	buf := make([]byte, 64*1024)
	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Memclr(buf)
		sink ^= buf[0]
	}
}

func BenchmarkStdClear(b *testing.B) {
	buf := make([]byte, 64*1024)
	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stdClear(buf)
		sink ^= buf[0]
	}
}

func BenchmarkMemclrInt64(b *testing.B) {
	data := make([]int64, 8192) // 64 KiB
	b.SetBytes(int64(len(data) * int(unsafe.Sizeof(data[0]))))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Memclr(data)
		sink2 ^= data[0]

	}
}

func BenchmarkStdClearInt64(b *testing.B) {
	data := make([]int64, 8192)
	b.SetBytes(int64(len(data) * int(unsafe.Sizeof(data[0]))))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range data {
			data[j] = 0
		}
		sink2 ^= data[0]
	}
}

type SimpleStruct struct {
	A int64
	B float64
	C bool
}

var sinkStruct SimpleStruct

func stdClearStruct(s []SimpleStruct) {
	for i := range s {
		s[i] = SimpleStruct{}
	}
}

func BenchmarkMemclrStruct(b *testing.B) {
	data := make([]SimpleStruct, 8192) // 64 KiB
	b.SetBytes(int64(len(data) * int(unsafe.Sizeof(data[0]))))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Memclr(data)
		sinkStruct = data[0]
	}
}

func TestContains(t *testing.T) {
	t.Run("int32", func(t *testing.T) {
		tests := []struct {
			name   string
			slice  []int32
			value  int32
			expect bool
		}{
			{"Empty", []int32{}, 42, false},
			{"SingleFound", []int32{42}, 42, true},
			{"SingleNotFound", []int32{41}, 42, false},
			{"First", []int32{42, 43, 44}, 42, true},
			{"Middle", []int32{41, 42, 43}, 42, true},
			{"Last", []int32{41, 42, 43}, 43, true},
			{"NotFound", []int32{41, 42, 43}, 44, false},
			{"Negative", []int32{-1, -2, -3}, -2, true},
			{"MinValue", []int32{math.MinInt32, 0}, math.MinInt32, true},
			{"MaxValue", []int32{0, math.MaxInt32}, math.MaxInt32, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := Contains(tt.slice, tt.value, true); got != tt.expect {
					t.Errorf("Contains(%v, %v) = %v, want %v", tt.slice, tt.value, got, tt.expect)
				}
			})
		}
	})

	t.Run("uint32", func(t *testing.T) {
		tests := []struct {
			name   string
			slice  []uint32
			value  uint32
			expect bool
		}{
			{"Zero", []uint32{0, 1, 2}, 0, true},
			{"MaxValue", []uint32{0, math.MaxUint32}, math.MaxUint32, true},
			{"NotFound", []uint32{1, 2, 3}, 4, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := Contains(tt.slice, tt.value, true); got != tt.expect {
					t.Errorf("Contains(%v, %v) = %v, want %v", tt.slice, tt.value, got, tt.expect)
				}
			})
		}
	})

	t.Run("float32", func(t *testing.T) {
		tests := []struct {
			name   string
			slice  []float32
			value  float32
			expect bool
		}{
			{"ExactMatch", []float32{1.1, 2.2, 3.3}, 2.2, true},
			{"Precision", []float32{0.1 + 0.2}, 0.3, true},
			{"NaN", []float32{float32(math.NaN())}, float32(math.NaN()), true},
			{"Inf", []float32{float32(math.Inf(1))}, float32(math.Inf(1)), true},
			{"NotFound", []float32{1.0, 2.0, 3.0}, 4.0, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := Contains(tt.slice, tt.value, true); got != tt.expect {
					t.Errorf("Contains(%v, %v) = %v, want %v", tt.slice, tt.value, got, tt.expect)
				}
			})
		}
	})

	t.Run("int64", func(t *testing.T) {
		tests := []struct {
			name   string
			slice  []int64
			value  int64
			expect bool
		}{
			{"LargeNumbers", []int64{10000000000, 20000000000}, 20000000000, true},
			{"MinValue", []int64{math.MinInt64, 0}, math.MinInt64, true},
			{"MaxValue", []int64{0, math.MaxInt64}, math.MaxInt64, true},
			{"NotFound", []int64{1, 2, 3}, 4, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := Contains(tt.slice, tt.value, true); got != tt.expect {
					t.Errorf("Contains(%v, %v) = %v, want %v", tt.slice, tt.value, got, tt.expect)
				}
			})
		}
	})

	t.Run("uint64", func(t *testing.T) {
		tests := []struct {
			name   string
			slice  []uint64
			value  uint64
			expect bool
		}{
			{"Zero", []uint64{0, 1, 2}, 0, true},
			{"MaxValue", []uint64{0, math.MaxUint64}, math.MaxUint64, true},
			{"NotFound", []uint64{1, 2, 3}, 4, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := Contains(tt.slice, tt.value, true); got != tt.expect {
					t.Errorf("Contains(%v, %v) = %v, want %v", tt.slice, tt.value, got, tt.expect)
				}
			})
		}
	})

	t.Run("float64", func(t *testing.T) {
		tests := []struct {
			name   string
			slice  []float64
			value  float64
			expect bool
		}{
			{"ExactMatch", []float64{1.1, 2.2, 3.3}, 2.2, true},
			{"Precision", []float64{0.1 + 0.2}, 0.3, true},
			{"NaN", []float64{math.NaN()}, math.NaN(), true},
			{"Inf", []float64{math.Inf(1)}, math.Inf(1), true},
			{"NotFound", []float64{1.0, 2.0, 3.0}, 4.0, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := Contains(tt.slice, tt.value, true); got != tt.expect {
					t.Errorf("Contains(%v, %v) = %v, want %v", tt.slice, tt.value, got, tt.expect)
				}
			})
		}
	})

	t.Run("string", func(t *testing.T) {
		tests := []struct {
			name   string
			slice  []string
			value  string
			expect bool
		}{
			{"ExactMatch", []string{"fkahfkajfka", "jfkajfka", "fjfkajfkjakfjakf"}, "jfkajfka", true},
			{"NotFound", []string{"a", "b", "c"}, "d", false},
			{"EmptySlice", []string{}, "any", false},
			{"FirstElement", []string{"match", "other"}, "match", true},
			{"LastElement", []string{"first", "second", "target"}, "target", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := Contains(tt.slice, tt.value, true); got != tt.expect {
					t.Errorf("Contains(%v, %v) = %v, want %v", tt.slice, tt.value, got, tt.expect)
				}
			})
		}
	})

	t.Run("UnsupportedTypes", func(t *testing.T) {
		t.Run("String", func(t *testing.T) {
			if Contains([]string{"a", "b", "c"}, "b", true) != true {
				t.Error("Contains should work for strings via fallback")
			}
		})

	})

	t.Run("EdgeCases", func(t *testing.T) {
		t.Run("LargeSlice", func(t *testing.T) {
			slice := make([]int64, 1_000_000)
			for i := range slice {
				slice[i] = int64(i)
			}

			if !Contains(slice, int64(999_999), true) {
				t.Error("Should find value in large slice")
			}

			if Contains(slice, int64(-1), true) {
				t.Error("Should not find missing value")
			}
		})

		t.Run("Alignment", func(t *testing.T) {
			slice := []int32{1, 2, 3, 4, 5}
			if !Contains(slice[1:], 5, true) {
				t.Error("Should handle unaligned slices")
			}
		})

		t.Run("CrossElement", func(t *testing.T) {
			slice := []int32{0x11223344, 0x55667788}
			if Contains(slice, 0x33445566, true) {
				t.Error("Should not find cross-element values")
			}
		})
	})
}

func TestBytesToUint64Slice_ValidInput(t *testing.T) {
	values := []uint64{0x1122334455667788, 0x99AABBCCDDEEFF00, 0x0123456789ABCDEF}
	buf := new(bytes.Buffer)
	for _, v := range values {
		_ = binary.Write(buf, binary.LittleEndian, v)
	}

	result := BytesToUint64Slice(buf.Bytes())
	if len(result) != len(values) {
		t.Fatalf("expected length %d, got %d", len(values), len(result))
	}

	for i, v := range result {
		if v != values[i] {
			t.Errorf("index %d: expected 0x%X, got 0x%X", i, values[i], v)
		}
	}
}

func BenchmarkStdClearStruct(b *testing.B) {
	data := make([]SimpleStruct, 8192)
	b.SetBytes(int64(len(data) * int(unsafe.Sizeof(data[0]))))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stdClearStruct(data)
		sinkStruct = data[0]
	}
}

func TestUnsafePointerExtraction(t *testing.T) {

	var empty []byte

	defer func() {

		r := recover()
		if r == nil {
			t.Errorf("Expected panic for empty slice, but no panic occurred")
		}
	}()

	_ = uintptr(unsafe.Pointer(&empty[0]))

	fmt.Println(String(data))
	fmt.Println(string3(data))
	fmt.Println(&data[0])

}

func TestUnsafeCompare(t *testing.T) {

	ss := "ss"
	ff := "ff"
	fmt.Println(CompareImpl(UnsafePointer(ss), UnsafePointer(ff), 2))

	ss2 := "ss"

	fmt.Println(CompareImpl(UnsafePointer(ss), UnsafePointer(ss2), 2))

}

func equalSlice[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSliceFill(t *testing.T) {
	slice := make([]int, 5)
	SliceFill(slice, 42, 1, 4)
	expected := []int{0, 42, 42, 42, 0}
	if !equalSlice(slice, expected) {
		t.Errorf("SliceFill failed: expected %v, got %v", expected, slice)
	}

	slice = make([]int, 5)
	SliceFill(slice, 42, 2, 2)
	expected = []int{0, 0, 0, 0, 0}
	if !equalSlice(slice, expected) {
		t.Errorf("SliceFill failed: expected %v, got %v", expected, slice)
	}

	t.Run("invalid range", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("SliceFill did not panic on invalid range")
			}
		}()
		slice := make([]int, 5)
		SliceFill(slice, 42, 4, 6)
	})
}

func BenchmarkSliceFill(b *testing.B) {
	slice := make([]int, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SliceFill(slice, 42, 0, 1000)
	}
}

func TestEscape(t *testing.T) {
	x := 42
	fmt.Printf("Before: %p\n", &x)
	y := Escape(x)
	fmt.Printf("After: %p\n", &y)

}

func BenchmarkSliceFillStandard(b *testing.B) {
	slice := make([]int, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range slice {
			slice[j] = 42
		}
	}
}

func BenchmarkMemsetSlice(b *testing.B) {
	slice := make([]int, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MemsetSlice(slice, 42)
	}
}

func BenchmarkMemsetSliceStandard(b *testing.B) {
	slice := make([]int, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range slice {
			slice[j] = 42
		}
	}
}

func BenchmarkCompareImplAVX2_16(b *testing.B) {
	a := make([]byte, 16)
	bb := make([]byte, 16)
	for i := 0; i < b.N; i++ {
		_ = CompareImpl(unsafe.Pointer(&a[0]), unsafe.Pointer(&bb[0]), 16)
	}
}

func BenchmarkBytesEqual_16(b *testing.B) {
	a := make([]byte, 16)
	bb := make([]byte, 16)
	for i := 0; i < b.N; i++ {
		_ = bytes.Equal(a, bb)
	}
}

func containsStd(slice [][]byte, value []byte) bool {
	for _, v := range slice {
		if bytes.Equal(v, value) {
			return true
		}
	}
	return false
}

func containsStdString(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func TestContainsByteSliceAVX2(t *testing.T) {
	tests := []struct {
		name     string
		slice    [][]byte
		value    []byte
		expected bool
	}{
		{
			name:     "EmptySlice",
			slice:    [][]byte{},
			value:    []byte{1, 2, 3},
			expected: false,
		},
		{
			name:     "SingleElementFound",
			slice:    [][]byte{{1, 2, 3}},
			value:    []byte{1, 2, 3},
			expected: true,
		},
		{
			name:     "SingleElementNotFound",
			slice:    [][]byte{{1, 2, 3}},
			value:    []byte{4, 5, 6},
			expected: false,
		},
		{
			name:     "MultipleElements",
			slice:    [][]byte{{1}, {2}, {3}, {4, 5}, make([]byte, 1024)},
			value:    []byte{4, 5},
			expected: true,
		},
		{
			name:     "LargeElement",
			slice:    [][]byte{bytes.Repeat([]byte{1}, 1000)},
			value:    bytes.Repeat([]byte{1}, 1000),
			expected: true,
		},
		{
			name:     "DifferentLength",
			slice:    [][]byte{{1, 2, 3, 4}},
			value:    []byte{1, 2, 3},
			expected: false,
		},
		{
			name:     "NilElement",
			slice:    [][]byte{nil, {1, 2, 3}},
			value:    []byte{1, 2, 3},
			expected: true,
		},
		{
			name:     "EmptyValue",
			slice:    [][]byte{{}, {1, 2, 3}},
			value:    []byte{},
			expected: true,
		},
		{
			name:     "EmptyValueNotFound",
			slice:    [][]byte{{1}, {2}, {3}},
			value:    []byte{},
			expected: false,
		},
		{
			name:     "FirstElement",
			slice:    [][]byte{{1, 2}, {3, 4}, {5, 6}},
			value:    []byte{1, 2},
			expected: true,
		},
		{
			name:     "LastElement",
			slice:    [][]byte{{1, 2}, {3, 4}, {5, 6}},
			value:    []byte{5, 6},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slicePtr := unsafe.Pointer(unsafe.SliceData(tt.slice))
			valuePtr := unsafe.Pointer(&tt.value)
			result := containsByteSliceAVX2(
				slicePtr,
				len(tt.slice),
				valuePtr,
				len(tt.value),
			)

			if result != tt.expected {
				t.Errorf("AVX2 implementation: expected %v, got %v", tt.expected, result)
			}

			resultGeneral := Contains(tt.slice, tt.value, true)
			if resultGeneral != tt.expected {
				t.Errorf("General implementation: expected %v, got %v", tt.expected, resultGeneral)
			}
		})
	}

	t.Run("Performance", func(t *testing.T) {
		sizes := []struct {
			sliceSize int
			elemSize  int
		}{
			{100, 16},
			{1000, 32},
			{10000, 64},
			{100000, 128},
		}

		for _, size := range sizes {
			slice := make([][]byte, size.sliceSize)
			for i := range slice {
				slice[i] = make([]byte, size.elemSize)
				rand.Read(slice[i])
			}

			target := make([]byte, size.elemSize)
			rand.Read(target)
			slice = append(slice, target)

			startAVX2 := time.Now()
			slicePtr := unsafe.Pointer(unsafe.SliceData(slice))
			valuePtr := unsafe.Pointer(&target)
			containsByteSliceAVX2(
				slicePtr,
				len(slice),
				valuePtr,
				len(target),
			)
			avx2Time := time.Since(startAVX2)

			startStd := time.Now()

			containsStd(slice, target)
			stdTime := time.Since(startStd)

			t.Logf("Size: %6d x %4d bytes | AVX2: %10v | Std: %10v | Speedup: %.2fx",
				size.sliceSize, size.elemSize,
				avx2Time, stdTime,
				float64(stdTime.Nanoseconds())/float64(avx2Time.Nanoseconds()),
			)
		}
	})
}

func BenchmarkContainsStringSlice(b *testing.B) {
	sizes := []struct {
		sliceSize int
		elemSize  int
	}{
		{100, 16},
		{1000, 32},
		{10000, 64},
		{100000, 128},
	}

	for _, size := range sizes {
		slice := make([]string, size.sliceSize)
		for i := range slice {
			bs := make([]byte, size.elemSize)
			rand.Read(bs)
			slice[i] = string(bs)
		}

		targetBytes := make([]byte, size.elemSize)
		rand.Read(targetBytes)
		target := string(targetBytes)
		slice = append(slice, target)

		b.Run(
			benchName(size.sliceSize, size.elemSize),
			func(b *testing.B) {
				b.Run("Std", func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						containsStdString(slice, target)
					}
				})

				b.Run("AVX2", func(b *testing.B) {
					b.ReportAllocs()
					// slicePtr := unsafe.Pointer(unsafe.SliceData(slice))
					// valuePtr := unsafe.Pointer(&target)
					for i := 0; i < b.N; i++ {
						Contains(slice, target, true)
						// containsStringAVX2(
						// 	slicePtr,
						// 	len(slice),
						// 	valuePtr,
						// 	len(target),
						// )
					}
				})

			},
		)
	}
}

func BenchmarkContainsByteSlice(b *testing.B) {
	sizes := []struct {
		sliceSize int
		elemSize  int
	}{
		{100, 16},
		{1000, 32},
		{10000, 64},
		{100000, 128},
	}

	for _, size := range sizes {
		slice := make([][]byte, size.sliceSize)
		for i := range slice {
			slice[i] = make([]byte, size.elemSize)
			rand.Read(slice[i])
		}

		target := make([]byte, size.elemSize)
		rand.Read(target)
		slice = append(slice, target)

		b.Run(
			benchName(size.sliceSize, size.elemSize),
			func(b *testing.B) {
				b.Run("Std", func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						containsStd(slice, target)
					}
				})

				b.Run("AVX2", func(b *testing.B) {
					b.ReportAllocs()
					slicePtr := unsafe.Pointer(unsafe.SliceData(slice))
					valuePtr := unsafe.Pointer(&target)
					for i := 0; i < b.N; i++ {
						containsByteSliceAVX2(
							slicePtr,
							len(slice),
							valuePtr,
							len(target),
						)
					}
				})
			},
		)
	}
}

func benchName(sliceSize, elemSize int) string {
	return fmt.Sprintf("%d_x_%d", sliceSize, elemSize)
}

func BenchmarkCompareImplAVX2_128(b *testing.B) {
	a := make([]byte, 128)
	bb := make([]byte, 128)
	for i := 0; i < b.N; i++ {
		_ = CompareImpl(unsafe.Pointer(&a[0]), unsafe.Pointer(&bb[0]), 128)
	}
}

func BenchmarkBytesEqual_128(b *testing.B) {
	a := make([]byte, 128)
	bb := make([]byte, 128)
	for i := 0; i < b.N; i++ {
		_ = bytes.Equal(a, bb)
	}
}

func BenchmarkCompareImplAVX2_256(b *testing.B) {
	a := make([]byte, 256)
	bb := make([]byte, 256)
	for i := 0; i < b.N; i++ {
		_ = CompareImpl(unsafe.Pointer(&a[0]), unsafe.Pointer(&bb[0]), 256)
	}
}

func BenchmarkBytesEqual_256(b *testing.B) {
	a := make([]byte, 256)
	bb := make([]byte, 256)
	for i := 0; i < b.N; i++ {
		_ = bytes.Equal(a, bb)
	}
}

func BenchmarkCompareImplAVX2_1024(b *testing.B) {
	a := make([]byte, 1024)
	bb := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		_ = CompareImpl(unsafe.Pointer(&a[0]), unsafe.Pointer(&bb[0]), 1024)
	}
}

func BenchmarkBytesEqual_1024(b *testing.B) {
	a := make([]byte, 1024)
	bb := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		_ = bytes.Equal(a, bb)
	}
}

func BenchmarkCompareImplAVX2_10240(b *testing.B) {
	a := make([]byte, 10240)
	bb := make([]byte, 10240)
	for i := 0; i < b.N; i++ {
		_ = CompareImpl(unsafe.Pointer(&a[0]), unsafe.Pointer(&bb[0]), 1024)
	}
}

func BenchmarkBytesEqual_10240(b *testing.B) {
	a := make([]byte, 10240)
	bb := make([]byte, 10240)
	for i := 0; i < b.N; i++ {
		_ = bytes.Equal(a, bb)
	}
}

// func BenchmarkDirtBytes(b *testing.B) {
// 	for size := block1kb; size < block1kb*20; size += block1kb * 2 {
// 		b.Run(fmt.Sprintf("size=%dkb", size/block1kb), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				data = MallocSlice[byte](size, size)
// 			}
// 		})
// 	}
// }
//
// func BenchmarkArrayBytes(b *testing.B) {
// 	for size := block1kb; size < block1kb*20; size += block1kb * 2 {
// 		b.Run(fmt.Sprintf("size=%dkb", size/block1kb), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				data = MakeSlice[byte](size, size)
// 			}
// 		})
// 	}
// }

func must[T any](v T, err int) T {

	return v
}

// func BenchmarkDirtSys(b *testing.B) {
// 	for size := block1kb; size < block1kb*20; size += block1kb * 2 {
// 		b.Run(fmt.Sprintf("size=%dkb", size/block1kb), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				data = *(*[]byte)(must(mem.SysAlloc(size)))
// 				mem.SysFree(data)
// 			}
// 		})
//
// 	}
// }

// var mm *arena.Arena
//
// func BenchmarkArenaBytes(b *testing.B) {
// 	for size := block1kb; size < block1kb*20; size += block1kb * 2 {
// 		b.Run(fmt.Sprintf("size=%dkb", size/block1kb), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				data, mm = MakeArenaSlice[byte](size, size)
// 				mm.Free()
//
// 			}
// 		})
//
// 	}
//
// }

// func BenchmarkDirtBytes_MakeNozero(b *testing.B) {
// 	for size := block1kb; size < block1kb*20; size += block1kb * 2 {
// 		b.Run(fmt.Sprintf("size=%dkb", size/block1kb), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				data = MakeNoZero(size)
// 			}
// 		})
// 	}
// }
//
// func BenchmarkOriginBytes(b *testing.B) {
// 	for size := block1kb; size < block1kb*20; size += block1kb * 2 {
// 		b.Run(fmt.Sprintf("size=%dkb", size/block1kb), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				data = make([]byte, size)
// 			}
// 		})
// 	}
// }

// func BenchmarkMake(b *testing.B) {
//
// 	b.Run("Make", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			data = make([]byte, 1000, 1000)
// 		}
// 	})
// }

func BenchmarkMutableString_SetString(b *testing.B) {
	var m MutableString
	for i := 0; i < b.N; i++ {
		m.SetString("benchmark")

		m.SetString("ben10")
	}
	_ = m
}

func BenchmarkString_SetString(b *testing.B) {
	var m []byte
	for i := 0; i < b.N; i++ {
		m = []byte("benchmark")
		m = []byte("ben10")
	}
	_ = m
}

func BenchmarkString_AppendString(b *testing.B) {
	s := ""
	for i := 0; i < 1000; i++ {
		s += "test"
	}
}
func BenchmarkMutableString_AppendString(b *testing.B) {
	var m MutableString
	for i := 0; i < 1000; i++ {
		m.AppendString("test")
	}
}

func BenchmarkMutableString_Clear(b *testing.B) {
	var m MutableString
	m.SetString("some long text")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Clear()
	}
}

func BenchmarkMutableString_ToUpper(b *testing.B) {
	m := MutableString("benchmark")
	for i := 0; i < b.N; i++ {
		m.ToUpper()
	}
}

func BenchmarkMutableString_ToLower(b *testing.B) {
	m := MutableString("BENCHMARK")
	for i := 0; i < b.N; i++ {
		m.ToLower()
	}
}

func TestMallocSlice(t *testing.T) {
	slice := MallocSlice[int](5, 10)

	if len(slice) != 5 {
		t.Errorf("Expected length 5, got %d", len(slice))
	}
	if cap(slice) != 10 {
		t.Errorf("Expected capacity 10, got %d", cap(slice))
	}
}

func TestMalloc(t *testing.T) {
	slice := MallocSlice[byte](5, 10)

	if len(slice) != 5 {
		t.Errorf("Expected length 5, got %d", len(slice))
	}
	if cap(slice) != 10 {
		t.Errorf("Expected capacity 10, got %d", cap(slice))
	}

	slice = append(slice, 1, 2, 3, 4, 5)

	if len(slice) != 10 {
		t.Errorf("Expected length 10, got %d", len(slice))
	}
	if cap(slice) != 10 {
		t.Errorf("Expected capacity 10, got %d", cap(slice))
	}
}

func TestCopyMakeSlice(t *testing.T) {
	slice := MakeSlice[int](5, 10)
	slice2 := make([]int, 0, 10)
	slice2 = append(slice2, 1, 2, 3, 4, 5)
	copy(slice, slice2)

	if len(slice) != 5 {
		t.Errorf("Expected length 5, got %d", len(slice))
	}
	if cap(slice) != 10 {
		t.Errorf("Expected capacity 10, got %d", cap(slice))
	}

	if slice[0] != 1 || slice[1] != 2 || slice[2] != 3 || slice[3] != 4 || slice[4] != 5 {
		t.Errorf("Incorrect: %v", slice)
	}

	// fmt.Println(slice)
}

type TestStruct struct {
	A int
	B float64
	C byte
}

func TestMallocStruct(t *testing.T) {
	slice := MallocSlice[TestStruct](3, 6)

	if len(slice) != 3 {
		t.Errorf("Expected length 3, got %d", len(slice))
	}
	if cap(slice) != 6 {
		t.Errorf("Expected capacity 6, got %d", cap(slice))
	}
}

type Example struct {
	Value int
}

// TestMalloc tests the Malloc function
func TestMalloc_noslice(t *testing.T) {
	// Test with an int
	intValue := 42
	intPointer := Malloc(intValue)
	if intPointer == nil {
		t.Errorf("Malloc returned nil for int")
	}
	// Test if the value of the int is the same as the passed value
	if *intPointer != intValue {
		t.Errorf("Expected %v for *intPointer, got %v", intValue, *intPointer)
	}

	// Test with a string
	stringValue := "Hello"
	stringPointer := Malloc(stringValue)
	if stringPointer == nil {
		t.Errorf("Malloc returned nil for string")
	}
	// Test if the value of the string is the same as the passed value
	if *stringPointer != stringValue {
		t.Errorf("Expected %v for *stringPointer, got %v", stringValue, *stringPointer)
	}

	// Test with a custom struct
	structValue := Example{Value: 100}
	structPointer := Malloc(structValue)
	if structPointer == nil {
		t.Errorf("Malloc returned nil for Example struct")
	}
	// Test if the struct's field is the same as the passed value
	if structPointer.Value != structValue.Value {
		t.Errorf("Expected %v for structPointer.Value, got %v", structValue.Value, structPointer.Value)
	}
}

func TestMakeSlice_DataManipulation(t *testing.T) {
	slice := MakeSlice[int](4, 6)

	for i := 0; i < len(slice); i++ {
		slice[i] = i * 10
	}

	for i, v := range slice {
		if v != i*10 {
			t.Errorf("Expected %d, got %d", i*10, v)
		}
	}
}

func TestMakeSliceStruct(t *testing.T) {
	slice := MakeSlice[TestStruct](3, 6)

	if len(slice) != 3 {
		t.Errorf("Expected length 3, got %d", len(slice))
	}
	if cap(slice) != 6 {
		t.Errorf("Expected capacity 6, got %d", cap(slice))
	}
}

func TestMakeSlice_FloatSlice(t *testing.T) {
	slice := MakeSlice[float64](3, 6)

	if len(slice) != 3 {
		t.Errorf("Expected: 3, got: %d", len(slice))
	}
	if cap(slice) != 6 {
		t.Errorf("Expected: 6, got: %d", cap(slice))
	}

	slice[0] = 1.1
	slice[1] = 2.2
	slice[2] = 3.3

	if slice[0] != 1.1 || slice[1] != 2.2 || slice[2] != 3.3 {
		t.Errorf("eror with float64: %v", slice)
	}
}

func TestMakeSlice_EmptySlice(t *testing.T) {
	slice := MakeSlice[int](0, 0)

	if len(slice) != 0 {
		t.Errorf("Expected: 0, got: %d", len(slice))
	}
	if cap(slice) != 0 {
		t.Errorf("Expected: 0, got: %d", cap(slice))
	}
}

func TestMakeSlice_StringSlice(t *testing.T) {
	slice := MakeSlice[string](3, 5)

	if len(slice) != 3 {
		t.Errorf("Expected: 3, got: %d", len(slice))
	}
	if cap(slice) != 5 {
		t.Errorf("Expected: 5, got: %d", cap(slice))
	}
	slice[0] = "hello"
	slice[1] = "world"

	if slice[0] != "hello" || slice[1] != "world" {
		t.Errorf("incorrect: %v", slice)
	}
}

func TestMutableString_String(t *testing.T) {
	ms := MutableString("Hello, world!")
	expected := "Hello, world!"

	result := ms.String()

	if result != expected {
		t.Errorf("Expected: %q, got: %q", expected, result)
	}
}

func TestMutableString_Set(t *testing.T) {
	var m MutableString

	m.SetString("hello")
	if m.String() != "hello" {
		t.Errorf("Expected 'hello', got '%s'", m.String())
	}

	m.SetString("world")
	if m.String() != "world" {
		t.Errorf("Expected 'world', got '%s'", m.String())
	}

	m.SetString("")
	if m.String() != "" {
		t.Errorf("Expected empty string, got '%s'", m.String())
	}
}

func TestMutableString_Clear(t *testing.T) {
	m := MutableString("hello")
	m.Clear()

	if len(m) != 0 {
		t.Errorf("Expected length 0, got %d", len(m))
	}

	if cap(m) < 5 {
		t.Errorf("Expected at least cap 5, got %d", cap(m))
	}
}

func TestMutableString_AppendString(t *testing.T) {
	var m MutableString
	m.AppendString("hello")
	m.AppendString(" world")

	if m.String() != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", m.String())
	}
}

func TestMutableString_AppendByte(t *testing.T) {
	var m MutableString
	m.AppendByte('A')
	m.AppendByte('B')

	if m.String() != "AB" {
		t.Errorf("Expected 'AB', got '%s'", m.String())
	}
}

func TestMutable(t *testing.T) {
	var ms MutableString

	ms.SetString("abcdef")
	fmt.Println(ms.String())
	ms.SetString("xyz")

	fmt.Println(ms.String())

	ms.AppendString("zzz")
	fmt.Println(ms.String())

}

func TestUnsafePointer(t *testing.T) {
	var ha = make(MutableString, 12)

	ha.SetString("Hello, world!")

	ptrOfInf := UnsafePointer(ha)

	fmt.Println(ptrOfInf)

	value := ConvertUnsafePointer[MutableString](ptrOfInf)

	fmt.Println(string(value))

}

func TestMutableString_Modify(t *testing.T) {
	ms := MutableString("Hello, world!")

	ms[0] = 'h'

	expected := "hello, world!"
	result := ms.String()

	if result != expected {
		t.Errorf("Expected: %q, got: %q", expected, result)
	}
}

func TestMutableString_Empty(t *testing.T) {
	var ms MutableString

	expected := ""
	result := ms.String()

	if result != expected {
		t.Errorf("Expected empty string, got: %q", result)
	}
}

func TestMutableString_SetString(t *testing.T) {
	var ms MutableString

	ms.SetString("Hello, Go!")
	expected := "Hello, Go!"

	result := ms.String()

	if result != expected {
		t.Errorf("Expected: %q, got: %q", expected, result)
	}
}

func TestMutableStringCreate(t *testing.T) {
	ms := make(MutableString, 12)
	ms.SetString("Hello, Go!")
	expected := "Hello, Go!"

	result := ms.StringNoZero()

	if result != expected {
		t.Errorf("Expected: %q, got: %q", expected, result)
	}
}

func ConvertSliceManual(from []int64) []int32 {
	to := make([]int32, len(from))
	for i, v := range from {
		to[i] = int32(v)
	}
	return to
}

func TestCopyUnsafe(t *testing.T) {
	// Test 1: Standard case
	source := []byte("Hello, World!")
	destination := make([]byte, len(source))
	n := CopyUnsafe(destination, source)
	if n != len(source) {
		t.Errorf("Expected %d bytes copied, got %d", len(source), n)
	}
	if string(destination) != string(source) {
		t.Errorf("Expected destination %q, got %q", string(source), string(destination))
	}

	// Test 2: Different lengths (source longer than destination)
	source = []byte("Hello, Go!")
	destination = make([]byte, len(source)-2) // Smaller destination
	_ = CopyUnsafe(destination, source)

	// Test 3: Empty slices
	source = []byte{}
	destination = []byte{}
	n = CopyUnsafe(destination, source)
	if n != 0 {
		t.Errorf("Expected 0 bytes copied, got %d", n)
	}
}

func TestStrings(t *testing.T) {
	// Test 1: Standard case
	source := []byte("Hello, World!")

	n := String(source)
	if n != string(source) {
		t.Errorf("Expected %s string copied, got %s", string(source), n)
	}

	source1 := "Hello,world!"

	n1 := StringToBytes(source1)
	if string(n1) != source1 {
		t.Errorf("Expected %s string copied, got %s", source, n1)
	}

	n2 := unsafeGetBytes(source1)
	if string(n2) != source1 {
		t.Errorf("Expected %s string copied, got %s", source, n2)
	}
	// n2 := _stringToBytes_(source1)
	// if string(n2) != source1 {
	// 	t.Errorf("Expected %s string copied, got %s", source, n2)
	// }

}

func BenchmarkCopyUnsafe(b *testing.B) {
	source := []byte("Benchmarking Unsafe Copy!")
	destination := make([]byte, len(source))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CopyUnsafe(destination, source)
	}
}

func BenchmarkCopyStandart(b *testing.B) {
	source := []byte("Benchmarking Standart Copy!")
	destination := make([]byte, len(source))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(destination, source)
	}
}

func BenchmarkCopy_currentBytes(b *testing.B) {
	source := []byte("Benchmarking Current Copy!")
	destination := make([]byte, len(source))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(destination[:26], source)
	}
}

func BenchmarkCopy_currentBytes_UNSAFE(b *testing.B) {
	source := []byte("Benchmarking Current Copy!")
	destination := make([]byte, len(source))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CopyUnsafe(destination[:26], source)
	}
}

func generateTestStrings(count, minLength, maxLength int) []string {
	var data []string
	for i := 0; i < count; i++ {
		length := minLength + i%(maxLength-minLength)
		str := strings.Repeat("a", length)
		data = append(data, str)
	}
	return data
}

func BenchmarkStringBuffer(b *testing.B) {
	for _, n := range []int{10, 100, 1000} {
		b.Run("Custom StringBuffer_"+fmt.Sprint(n), func(b *testing.B) {
			sb := NewStringBuffer(n + 10) // or sb := NewStringBuffer(0) // if you not sure what kind of capasity you need
			// n + 10 is to account for some extra space
			for i := 0; i < b.N; i++ {
				sb.Reset()
				for j := 0; j < n; j++ {
					sb.WriteString("a")
				}
				_ = sb.String()
			}
		})
	}
}

func BenchmarkStringsBuilder(b *testing.B) {
	for _, n := range []int{10, 100, 1000} {
		b.Run("Standard strings.Builder_"+fmt.Sprint(n), func(b *testing.B) {
			var sb strings.Builder
			for i := 0; i < b.N; i++ {
				sb.Reset()
				for j := 0; j < n; j++ {
					sb.WriteString("a")
				}
				_ = sb.String()
			}
		})
	}
}

func BenchmarkString(b *testing.B) {
	data := []byte("This is a benchmark test for String conversion.")

	b.Run("*(*string)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = String(data)
		}
	})

	b.Run("*(*string)(&struct)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = string2(data)
		}

	})

	b.Run(" unsafe.String", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = string3(data)
		}

	})

	b.Run("headerstring", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = string1(data)
		}

	})
	b.Run(" with custom struct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = string_b(data)
		}
	})

	b.Run("multiply headers converting", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = string4(data)
		}
	})
	b.Run("Standard String", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = string(data)
		}
	})
}

func BenchmarkStringToBytesSmallString(b *testing.B) {
	data := "I'm looking forward to season 5 of the boys "

	b.Run("convert with &struct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = StringToBytes(data)
		}
	})

	b.Run("Standard StringToBytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = []byte(data)
		}
	})

	b.Run("unsafeSlice StringToBytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = stringToBytes_(data)
		}

	})
	b.Run("simple converting with *(*) StringToBytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = _stringToBytes_(data)
		}

	})
	b.Run("convert to string with multiply headers", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = unsafeGetBytes(data)
		}
	})

	b.Run("convert to string with custom &struct with unsafe.Pointer data", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = stringBytes(data)
		}

	})

	b.Run("convert to string with array", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = unsafeGetBytes_2(data)
		}

	})

}

// Benchmark Memory Allocation and Copying
func BenchmarkMakeNoZero(b *testing.B) {
	size := 1024

	b.Run("Custom MakeNoZero", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := MakeNoZero(size)
			_ = buf

		}
	})

	b.Run("Standard make([]byte)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := make([]byte, size)
			_ = buf
		}
	})
}

func BenchmarkMakeNoZeroString(b *testing.B) {
	size := 1024

	b.Run("Custom MakeNoZero", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := MakeNoZeroString(size)
			_ = buf

		}
	})

	b.Run("Standard make([]string)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := make([]string, size)
			_ = buf
		}
	})
}

func BenchmarkMakeNoZeroStringSmall(b *testing.B) {
	size := 5

	b.Run("Custom MakeNoZero", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := MakeNoZeroString(size)
			_ = buf

		}
	})

	b.Run("Standard make([]string)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := make([]string, size)
			_ = buf
		}
	})
}

type fny struct {
	a0, a1, a2, a3, a4, a5, a6, a7, a8, a9 string
	b0, b1, b2, b3, b4, b5, b6, b7, b8, b9 int
	c0, c1, c2, c3, c4, c5, c6, c7, c8, c9 float64
	d0, d1, d2, d3, d4, d5, d6, d7, d8, d9 bool
	e0, e1, e2, e3, e4, e5, e6, e7, e8, e9 struct {
		f0, f1, f2, f3, f4, f5, f6, f7, f8, f9 int
	}
}

func BenchmarkMakeNoZeroAny(b *testing.B) {

	b.Run("Custom MakeNoZero", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := MakeZero[fny](1024)
			_ = buf
		}
	})

	b.Run("Standard make([]fny)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := make([]fny, 1024)
			_ = buf
		}
	})
}

func BenchmarkMakeNoZeroCapAny(b *testing.B) {

	b.Run("Custom MakeNoZero", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := MakeZeroCap[fny](0, 1024)
			_ = buf
		}
	})

	b.Run("Standard make([]fny)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := make([]fny, 0, 1024)
			_ = buf
		}
	})
}

func BenchmarkStringToBytes(b *testing.B) {
	// Generate a large number of strings for testing
	testStrings := generateTestStrings(100000, 10, 100) // 100,000 strings with lengths between 10 and 100

	b.Run(" with &struct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testStrings {
				_ = StringToBytes(str)
			}
		}
	})

	b.Run("Standard StringToBytes BigString", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testStrings {
				_ = []byte(str)
			}
		}
	})
	b.Run("just converting *(*string)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testStrings {
				_ = _stringToBytes_(str)
			}
		}

	})
	b.Run("unsafeSlice StringToBytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testStrings {
				_ = unsafeGetBytes(str)
			}
		}

	})
}

func TestEquals(t *testing.T) {
	t.Run("TestEqualsTrue", func(t *testing.T) {
		a := []byte("This is a benchmark test for Equal.....")
		bb := []byte("This is a benchmark test for Equal.....")
		lengthA := uintptr(len(a))
		boole := Equal(a, bb, lengthA)
		if boole == false {
			t.Log("IsEqual: ", Equal(a, bb, lengthA))
		}
	})
	t.Run("TestEqualsFalse", func(t *testing.T) {
		a := []byte("This is a benchmark test for Equal.....")
		bb := []byte("This is a benchmark test for Equal.....1")
		lengthA := uintptr(len(a))
		boole := Equal(a, bb, lengthA)
		if boole == true {
			t.Log("IsEqual: ", Equal(a, bb, lengthA))
		}

	})
}

func BenchmarkEqualTrue(b *testing.B) { // true
	a := []byte("This is a benchmark test for Equal.....")
	bb := []byte("This is a benchmark test for Equal.....")
	lengthA := uintptr(len(a))

	b.Run("Custom Equal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			boole := Equal(a, bb, lengthA)
			if boole == false {
				b.Log("IsEqual: ", Equal(a, bb, lengthA))
			}
		}
	})

	b.Run("Standard bytes.Equal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			boole := bytes.Equal(a, bb)
			if boole == false {
				b.Log("bytes.Equal: ", bytes.Equal(a, bb))
			}

		}
	})

	b.Run("TEST GENERIC EQUAL", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			boole := String(a) == String(bb)
			if boole == false {
				b.Fatal("GenericEqual: ", Equal(a, bb, lengthA))
			}
		}
	})
}

func BenchmarkConvertSlice(b *testing.B) {
	data := []int64{1, 2, 3, 4, 5}

	b.Run("Custom ConvertSlice", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			_, _ = ConvertSlice[int64, int32](data)
		}
	})

	b.Run("Manual ConvertSlice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ConvertSliceManual(data)
		}
	})
}

func BenchmarkGetItem(b *testing.B) {
	intSlice := make([]int, 10000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = GetItem(intSlice, 5000)
	}
}

func BenchmarkStandartIndexing(b *testing.B) {
	intSlice := make([]int, 10000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = intSlice[9000]
	}
}

func BenchmarkGetTimeWithoutChecking(b *testing.B) {
	intSlice := make([]int, 10000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = GetItemWithoutCheck(intSlice, 9000)
	}
}

func BenchmarkNextPower2(b *testing.B) {

	number1 := uintptr(49130)

	b.Run(" NextPowerOfTwo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			number1 = NextPowerOfTwo(number1)
		}
	})

}

func BenchmarkLog2(b *testing.B) {
	const inputUint = uintptr(49130)
	const inputFlt = float64(49130)

	b.Run("math.Log2(float64)", func(b *testing.B) {
		var x float64
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			x = math.Log2(inputFlt)
		}
		_ = x
	})

	b.Run("Log2(uintptr)", func(b *testing.B) {
		var x uintptr
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			x = Log2(inputUint)
		}
		_ = x
	})
}
func TestUnion_SetInt64_GetInt64(t *testing.T) {
	u := union.NewUnion[int64]()
	val := int64(42)

	u.SetInt64(val)
	got := u.GetInt64()

	if got != val {
		t.Errorf("GetInt64() = %v, want %v", got, val)
	}
}

func TestUnion_SetFloat64_GetFloat64(t *testing.T) {
	u := union.NewUnion[float64]()
	val := float64(3.14)

	u.SetFloat64(val)
	got := u.GetFloat64()

	if got != val {
		t.Errorf("GetFloat64() = %v, want %v", got, val)
	}
}

func TestUnion_Set_Get_Int64(t *testing.T) {
	u := union.NewUnion[int64]()
	val := int64(100)

	err := u.Set(val)
	if err != nil {
		t.Fatalf("Set() failed: %v", err)
	}

	got, err := u.Get()
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if got != val {
		t.Errorf("Get() = %v, want %v", got, val)
	}
}

func TestUnion_Set_Get_Float64(t *testing.T) {
	u := union.NewUnion[float64]()
	val := float64(2.71)

	err := u.Set(val)
	if err != nil {
		t.Fatalf("Set() failed: %v", err)
	}

	got, err := u.Get()
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if got != val {
		t.Errorf("Get() = %v, want %v", got, val)
	}
}

func TestUnion_Set_UnsupportedType(t *testing.T) {
	u := union.NewUnion[string]()
	val := "unsupported type"

	err := u.Set(val)
	if err == nil {
		t.Error("Expected error for unsupported type, got nil")
	}
}

func TestUnion_Get_Uninitialized(t *testing.T) {
	u := union.NewUnion[int64]()

	_, err := u.Get()
	if err == nil {
		t.Error("Expected error for uninitialized union, got nil")
	}
}

func TestUnion_Set_TypeTooLarge(t *testing.T) {
	u := union.NewUnion[[16]byte]()
	val := [16]byte{}

	err := u.Set(val)
	if err == nil {
		t.Error("Expected error for type too large, got nil")
	}
}

func TestUnion_SizeOf(t *testing.T) {
	u := union.NewUnion[int64]()
	val := int64(0)

	size := unsafe.Sizeof(val)
	if size > uintptr(len(u.Data())) {
		t.Errorf("Size of int64 (%d) exceeds union capacity (%d)", size, len(u.Data()))
	}
}

func TestGetField(t *testing.T) {
	s := example.NewExmp()

	privateField := GetPrivateField[example.Exmp, int](&s, "private")

	fmt.Println(*privateField)

}

//
// func TestArena(t *testing.T) {
// 	f := MakeShareArray[string](0, 10)
// 	f.Append(("a"))
// 	gg := f.GetElement(0)
// 	fmt.Println(string(*gg))
// 	f.Append(("b"))
// 	fmt.Println(string(*f.GetElement(1)))
// 	fmt.Println(string(*f.GetElement(0)))
//
// 	f.Append("cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc")
// 	fmt.Println(string(*f.GetElement(2)))
//
// }
