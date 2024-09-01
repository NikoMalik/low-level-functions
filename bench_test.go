package lowlevelfunctions

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func ConvertSliceManual(from []int64) []int32 {
	to := make([]int32, len(from))
	for i, v := range from {
		to[i] = int32(v)
	}
	return to
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

	b.Run("Custom String", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = String(data)
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

	b.Run("Custom StringToBytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = StringToBytes(data)
		}
	})

	b.Run("Standard StringToBytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = []byte(data)
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

	b.Run("Custom StringToBytesBigString", func(b *testing.B) {
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
}

func BenchmarkEqualTrue(b *testing.B) { // true
	a := []byte("This is a benchmark test for Equal.....")
	bb := []byte("This is a benchmark test for Equal.....")

	b.Run("Custom Equal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			boole := Equal(a, bb)
			if boole == false {
				b.Log("IsEqual: ", Equal(a, bb))
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

func BenchmarkStandardIndexing(b *testing.B) {
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
