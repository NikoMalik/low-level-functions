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

// Benchmark String Conversion
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

	b.Run("Standard make([]byte)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := make([]string, size)
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

func BenchmarkAtomicCounter(b *testing.B) {
	counter := AtomicCounter{}
	withoutPad := AtomicCounterWithoutPad{}
	withShar := ShardedAtomicCounter{}
	withSharNoPad := ShardedAtomicCounterWithoutPad{}

	b.Run("With CacheLine Padding", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			go func() {
				counter.Increment(1)
				_ = counter.Get()
			}()
		}
	})

	b.Run("Without CacheLine Padding", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			go func() {
				withoutPad.Increment(1)
				_ = withoutPad.Get()
			}()
		}
	})

	b.Run("With CacheLine Padding and shard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			go func() {
				withShar.Increment(1)
				_ = withShar.Get()
			}()
		}
	})

	b.Run("Without CacheLine Padding and with shard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			go func() {
				withSharNoPad.Increment(1)
				_ = withSharNoPad.Get()
			}()
		}
	})

}
