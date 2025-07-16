package lowlevelfunctions

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"unsafe"

	"github.com/NikoMalik/low-level-functions/example"

	"github.com/NikoMalik/low-level-functions/union"
)

var block1kb = 1024
var data []byte

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
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for size mismatch, but none occurred")
		}
	}()
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

	b.Run("Default NextPowerOfTwo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			number1 = NextPowerOfTwo(number1)
		}
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
