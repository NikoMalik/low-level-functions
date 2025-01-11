package lowlevelfunctions

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestMutableString_String(t *testing.T) {
	ms := MutableString("Hello, world!")
	expected := "Hello, world!"

	result := ms.String()

	if result != expected {
		t.Errorf("Expected: %q, got: %q", expected, result)
	}
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
