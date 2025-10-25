package lowlevelfunctions

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/bits"
	"reflect"
	"strings"
	"sync/atomic"
	"unicode/utf8"
	"unsafe"

	"github.com/NikoMalik/low-level-functions/constants"
)

const (
	PtrSize                = 4 << (^uintptr(0) >> 63)
	StrSize                = unsafe.Sizeof("")
	SliceSize              = int(unsafe.Sizeof([]byte{}))
	CacheLineSize          = constants.CacheLinePadSize
	MaxInt32               = 1<<31 - 1
	MaxUintptr             = ^uintptr(0)
	minSizeForMallocHeader = PtrSize * ptrBits
	ptrBits                = 8 * PtrSize
	_PageSize              = 1 << _PageShift
	_PageMask              = _PageSize - 1
)

var hexTbl = [16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

var reverseHexTable = []byte(
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\xff\xff\xff\xff\xff\xff" +
		"\xff\x0a\x0b\x0c\x0d\x0e\x0f\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\x0a\x0b\x0c\x0d\x0e\x0f\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff",
)

type InvalidByteError byte

func (e InvalidByteError) Error() string {
	return "invalid byte: " + string(e)
}

var ErrInvalidFormat = errors.New("invalid format")

var ErrLength = errors.New("invalid hex string length")

func EncodedLen(n int) int { return n * 2 }

func DecodeUnrolled16(dst, src []byte) (int, error) { // with checks
	n := len(src)
	if n == 0 {
		return 0, nil
	}
	if n%2 == 1 {
		if reverseHexTable[src[n-1]] > 0x0f {
			return 0, InvalidByteError(src[n-1])
		}
		return 0, ErrLength
	}

	srcPtr := unsafe.Pointer(unsafe.SliceData(src))
	dstPtr := unsafe.Pointer(unsafe.SliceData(dst))
	t := &reverseHexTable
	i := 0
	j := 0

	for ; i+16 <= n; i += 16 {
		p0 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+0)))
		q0 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+1)))
		p1 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+2)))
		q1 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+3)))
		p2 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+4)))
		q2 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+5)))
		p3 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+6)))
		q3 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+7)))
		p4 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+8)))
		q4 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+9)))
		p5 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+10)))
		q5 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+11)))
		p6 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+12)))
		q6 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+13)))
		p7 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+14)))
		q7 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+15)))

		a0 := (*t)[p0]
		b0 := (*t)[q0]
		if a0 > 0x0f {
			return j, InvalidByteError(p0)
		}
		if b0 > 0x0f {
			return j, InvalidByteError(q0)
		}
		a1 := (*t)[p1]
		b1 := (*t)[q1]
		if a1 > 0x0f {
			return j, InvalidByteError(p1)
		}
		if b1 > 0x0f {
			return j, InvalidByteError(q1)
		}
		a2 := (*t)[p2]
		b2 := (*t)[q2]
		if a2 > 0x0f {
			return j, InvalidByteError(p2)
		}
		if b2 > 0x0f {
			return j, InvalidByteError(q2)
		}
		a3 := (*t)[p3]
		b3 := (*t)[q3]
		if a3 > 0x0f {
			return j, InvalidByteError(p3)
		}
		if b3 > 0x0f {
			return j, InvalidByteError(q3)
		}
		a4 := (*t)[p4]
		b4 := (*t)[q4]
		if a4 > 0x0f {
			return j, InvalidByteError(p4)
		}
		if b4 > 0x0f {
			return j, InvalidByteError(q4)
		}
		a5 := (*t)[p5]
		b5 := (*t)[q5]
		if a5 > 0x0f {
			return j, InvalidByteError(p5)
		}
		if b5 > 0x0f {
			return j, InvalidByteError(q5)
		}
		a6 := (*t)[p6]
		b6 := (*t)[q6]
		if a6 > 0x0f {
			return j, InvalidByteError(p6)
		}
		if b6 > 0x0f {
			return j, InvalidByteError(q6)
		}
		a7 := (*t)[p7]
		b7 := (*t)[q7]
		if a7 > 0x0f {
			return j, InvalidByteError(p7)
		}
		if b7 > 0x0f {
			return j, InvalidByteError(q7)
		}

		*(*byte)(unsafe.Add(dstPtr, uintptr(j+0))) = (a0 << 4) | b0
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+1))) = (a1 << 4) | b1
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+2))) = (a2 << 4) | b2
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+3))) = (a3 << 4) | b3
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+4))) = (a4 << 4) | b4
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+5))) = (a5 << 4) | b5
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+6))) = (a6 << 4) | b6
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+7))) = (a7 << 4) | b7

		j += 8
	}

	for ; i < n; i += 2 {
		p := *(*byte)(unsafe.Add(srcPtr, uintptr(i)))
		q := *(*byte)(unsafe.Add(srcPtr, uintptr(i+1)))

		a := (*t)[p]
		b := (*t)[q]
		if a > 0x0f {
			return j, InvalidByteError(p)
		}
		if b > 0x0f {
			return j, InvalidByteError(q)
		}

		*(*byte)(unsafe.Add(dstPtr, uintptr(j))) = (a << 4) | b
		j++
	}

	return j, nil
}

// ValidateHex checks if src is a valid hex string (contains only 0-9, a-f, A-F and has even length).
func ValidateHex(src []byte) error {
	if len(src)%2 == 1 {
		if reverseHexTable[src[len(src)-1]] > 0x0f {
			return InvalidByteError(src[len(src)-1])
		}
		return ErrLength
	}

	srcPtr := unsafe.Pointer(unsafe.SliceData(src))
	t := &reverseHexTable
	i := 0

	for ; i+8 <= len(src); i += 8 {
		p0 := *(*byte)(unsafe.Add(srcPtr, uintptr(i)))
		p1 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+1)))
		p2 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+2)))
		p3 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+3)))
		p4 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+4)))
		p5 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+5)))
		p6 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+6)))
		p7 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+7)))

		a0 := (*t)[p0]
		a1 := (*t)[p1]
		a2 := (*t)[p2]
		a3 := (*t)[p3]
		a4 := (*t)[p4]
		a5 := (*t)[p5]
		a6 := (*t)[p6]
		a7 := (*t)[p7]

		if (a0 | a1 | a2 | a3 | a4 | a5 | a6 | a7) > 0x0f {
			if a0 > 0x0f {
				return InvalidByteError(p0)
			}
			if a1 > 0x0f {
				return InvalidByteError(p1)
			}
			if a2 > 0x0f {
				return InvalidByteError(p2)
			}
			if a3 > 0x0f {
				return InvalidByteError(p3)
			}
			if a4 > 0x0f {
				return InvalidByteError(p4)
			}
			if a5 > 0x0f {
				return InvalidByteError(p5)
			}
			if a6 > 0x0f {
				return InvalidByteError(p6)
			}
			if a7 > 0x0f {
				return InvalidByteError(p7)
			}
		}
	}

	for ; i < len(src); i++ {
		p := *(*byte)(unsafe.Add(srcPtr, uintptr(i)))
		if (*t)[p] > 0x0f {
			return InvalidByteError(p)
		}
	}

	return nil
}

// fast analog for word % p
// https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
func Fastrange(word, p uint64) uint64 {
	hi, _ := bits.Mul64(word, p)
	return hi
}

// safety decode with checks
func DecodeUnrolled8Checks(dst, src []byte) (int, error) {
	n := len(src)
	if n == 0 {
		return 0, nil
	}
	if n%2 == 1 {
		if reverseHexTable[src[n-1]] > 0x0f {
			return 0, InvalidByteError(src[n-1])
		}
		return 0, ErrLength
	}

	srcPtr := unsafe.Pointer(unsafe.SliceData(src))
	dstPtr := unsafe.Pointer(unsafe.SliceData(dst))
	t := &reverseHexTable
	i := 0
	j := 0

	for ; i+8 <= n; i += 8 {
		p0 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+0)))
		q0 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+1)))
		p1 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+2)))
		q1 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+3)))
		p2 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+4)))
		q2 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+5)))
		p3 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+6)))
		q3 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+7)))

		a0 := (*t)[p0]
		b0 := (*t)[q0]
		a1 := (*t)[p1]
		b1 := (*t)[q1]
		a2 := (*t)[p2]
		b2 := (*t)[q2]
		a3 := (*t)[p3]
		b3 := (*t)[q3]

		if (a0 | b0 | a1 | b1 | a2 | b2 | a3 | b3) > 0x0f {
			if a0 > 0x0f {
				return j, InvalidByteError(p0)
			}
			if b0 > 0x0f {
				return j, InvalidByteError(q0)
			}
			if a1 > 0x0f {
				return j, InvalidByteError(p1)
			}
			if b1 > 0x0f {
				return j, InvalidByteError(q1)
			}
			if a2 > 0x0f {
				return j, InvalidByteError(p2)
			}
			if b2 > 0x0f {
				return j, InvalidByteError(q2)
			}
			if a3 > 0x0f {
				return j, InvalidByteError(p3)
			}
			if b3 > 0x0f {
				return j, InvalidByteError(q3)
			}
		}

		*(*byte)(unsafe.Add(dstPtr, uintptr(j+0))) = (a0 << 4) | b0
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+1))) = (a1 << 4) | b1
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+2))) = (a2 << 4) | b2
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+3))) = (a3 << 4) | b3

		j += 4
	}

	for ; i < n; i += 2 {
		p := *(*byte)(unsafe.Add(srcPtr, uintptr(i)))
		q := *(*byte)(unsafe.Add(srcPtr, uintptr(i+1)))

		a := (*t)[p]
		b := (*t)[q]
		if a > 0x0f {
			return j, InvalidByteError(p)
		}
		if b > 0x0f {
			return j, InvalidByteError(q)
		}

		*(*byte)(unsafe.Add(dstPtr, uintptr(j))) = (a << 4) | b
		j++
	}

	return j, nil
}

// based, must be correct hex,no checks
func DecodeUnrolled8(dst, src []byte) int {
	n := len(src)
	if n == 0 {
		return 0
	}
	srcPtr := unsafe.Pointer(unsafe.SliceData(src))
	dstPtr := unsafe.Pointer(unsafe.SliceData(dst))
	t := &reverseHexTable
	i := 0
	j := 0
	for ; i+8 <= n; i += 8 {
		p0 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+0)))
		q0 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+1)))
		p1 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+2)))
		q1 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+3)))
		p2 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+4)))
		q2 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+5)))
		p3 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+6)))
		q3 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+7)))
		a0 := (*t)[p0]
		b0 := (*t)[q0]
		a1 := (*t)[p1]
		b1 := (*t)[q1]
		a2 := (*t)[p2]
		b2 := (*t)[q2]
		a3 := (*t)[p3]
		b3 := (*t)[q3]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+0))) = (a0 << 4) | b0
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+1))) = (a1 << 4) | b1
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+2))) = (a2 << 4) | b2
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+3))) = (a3 << 4) | b3
		j += 4
	}

	for ; i < n; i += 2 {
		p := *(*byte)(unsafe.Add(srcPtr, uintptr(i)))
		q := *(*byte)(unsafe.Add(srcPtr, uintptr(i+1)))
		a := (*t)[p]
		b := (*t)[q]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j))) = (a << 4) | b
		j++
	}

	return j
}

func EncodeUnrolled8(dst, src []byte) int {
	n := len(src)
	if n == 0 {
		return 0
	}
	srcPtr := unsafe.Pointer(unsafe.SliceData(src))
	dstPtr := unsafe.Pointer(unsafe.SliceData(dst))
	j := 0
	i := 0
	t := &hexTbl
	for ; i+8 <= n; i += 8 {
		v0 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+0)))
		v1 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+1)))
		v2 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+2)))
		v3 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+3)))
		v4 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+4)))
		v5 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+5)))
		v6 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+6)))
		v7 := *(*byte)(unsafe.Add(srcPtr, uintptr(i+7)))

		*(*byte)(unsafe.Add(dstPtr, uintptr(j+0))) = t[v0>>4]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+1))) = t[v0&0x0f]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+2))) = t[v1>>4]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+3))) = t[v1&0x0f]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+4))) = t[v2>>4]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+5))) = t[v2&0x0f]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+6))) = t[v3>>4]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+7))) = t[v3&0x0f]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+8))) = t[v4>>4]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+9))) = t[v4&0x0f]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+10))) = t[v5>>4]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+11))) = t[v5&0x0f]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+12))) = t[v6>>4]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+13))) = t[v6&0x0f]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+14))) = t[v7>>4]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+15))) = t[v7&0x0f]

		j += 16
	}

	for ; i < n; i++ {
		v := *(*byte)(unsafe.Add(srcPtr, uintptr(i)))
		*(*byte)(unsafe.Add(dstPtr, uintptr(j))) = t[v>>4]
		*(*byte)(unsafe.Add(dstPtr, uintptr(j+1))) = t[v&0x0f]
		j += 2
	}

	return j
}

// Bool2int returns 0 if x is false or 1 if x is true.
func Bool2int(x bool) int {
	// Avoid branches. In the SSA compiler, this compiles to
	// exactly what you would want it to.
	return int(*(*uint8)(unsafe.Pointer(&x)))
}

func MulUintptr(a, b uintptr) (uintptr, bool) {
	if a|b < 1<<(4*PtrSize) || a == 0 {
		return a * b, false
	}
	overflow := b > MaxUintptr/a
	return a * b, overflow
}

func isLittleEndian() bool {
	var i uint16 = 0x0102
	b := (*[2]byte)(unsafe.Pointer(&i))
	return b[0] == 0x02
}
func Uint32_littleEndian(b *[4]byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func Uint32_bigEndian(b *[4]byte) uint32 {
	return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
}

func Uint64_littleEndian(b *[8]byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func Uint64_bigEndian(b *[8]byte) uint64 {
	return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
}

func ReadUnaligned32(p unsafe.Pointer) uint32 {
	b := (*[4]byte)(p)
	if !isLittleEndian() {
		return Uint32_bigEndian(b)
	}
	return Uint32_littleEndian(b)
}

//go:linkname memclrNoHeapPointers runtime.memclrNoHeapPointers
func memclrNoHeapPointers(p unsafe.Pointer, n uintptr)

// MemclrZero sets memory of slice to zero, assuming T has no heap pointers.
// T MUST NOT contain any references (e.g. pointers, strings, slices, maps, funcs).
func Memclr[T any](s []T) {
	if len(s) == 0 {
		return
	}
	size := unsafe.Sizeof(s[0]) * uintptr(len(s))
	ptr := unsafe.Pointer(&s[0])
	memclrNoHeapPointers(ptr, size)
}

// //go:linkname readUnaligned64 runtime.readUnaligned64
// func readUnaligned64(p unsafe.Pointer) uint64
func ReadUnaligned64(p unsafe.Pointer) uint64 {
	b := (*[8]byte)(p)
	if !isLittleEndian() {
		return Uint64_bigEndian(b)
	}
	return Uint64_littleEndian(b)
}

func MemsetSlice[T any](s []T, value T) {
	if len(s) == 0 {
		return
	}
	s[0] = value
	for i := 1; i < len(s); i *= 2 {
		CopyUnsafe(s[i:], s[:i])
	}
}

type ContainKey interface {
	[]byte | string | int8 | uint8 | int16 | uint16 | int32 | uint32 | float32 |
		int64 | uint64 | float64
}

// only work for predifined types
func Contains[T ContainKey](slice []T, value T, littleEndian bool) bool {
	if len(slice) == 0 {
		return false
	}
	switch any(*new(T)).(type) {
	case int8, uint8:
		data := unsafe.Slice(
			(*byte)(unsafe.Pointer(unsafe.SliceData(slice))),
			len(slice),
		)
		byteVal := *(*byte)(unsafe.Pointer(&value))
		return bytes.Contains(data, []byte{byteVal})

	case int16, uint16:
		data := unsafe.Slice(
			(*byte)(unsafe.Pointer(unsafe.SliceData(slice))),
			len(slice)*2,
		)
		var buf [2]byte
		if littleEndian {
			binary.LittleEndian.PutUint16(buf[:], *(*uint16)(unsafe.Pointer(&value)))
			return bytes.Contains(data, buf[:])
		}
		binary.BigEndian.PutUint16(buf[:], *(*uint16)(unsafe.Pointer(&value)))
		return bytes.Contains(data, buf[:])

	case int32, uint32, float32:
		data := unsafe.Slice(
			(*byte)(unsafe.Pointer(unsafe.SliceData(slice))),
			len(slice)*4,
		)
		var buf [4]byte
		if littleEndian {
			binary.LittleEndian.PutUint32(buf[:], *(*uint32)(unsafe.Pointer(&value)))
			return bytes.Contains(data, buf[:])
		}
		binary.BigEndian.PutUint32(buf[:], *(*uint32)(unsafe.Pointer(&value)))
		return bytes.Contains(data, buf[:])

	case int64, uint64, float64:
		data := unsafe.Slice(
			(*byte)(unsafe.Pointer(unsafe.SliceData(slice))),
			len(slice)*8,
		)
		var buf [8]byte
		if littleEndian {
			binary.LittleEndian.PutUint64(buf[:], *(*uint64)(unsafe.Pointer(&value)))
			return bytes.Contains(data, buf[:])

		}
		binary.BigEndian.PutUint64(buf[:], *(*uint64)(unsafe.Pointer(&value)))

		return bytes.Contains(data, buf[:])

	case string:
		slicePtr := unsafe.Pointer(unsafe.SliceData(slice))
		sliceLen := len(slice)
		valuePtr := Noescape(unsafe.Pointer(&value))
		value := *(*string)(Noescape(valuePtr))
		return containsStringAVX2(
			slicePtr,
			sliceLen,
			valuePtr,
			len(value),
		)

	case []byte:
		slicePtr := unsafe.Pointer(unsafe.SliceData(slice))
		sliceLen := len(slice)
		valuePtr := Noescape(unsafe.Pointer(&value))
		value := *(*[]byte)(Noescape(valuePtr))

		return containsByteSliceAVX2(
			slicePtr,
			sliceLen,
			valuePtr,
			len(value),
		)

	default:
		return false
	}
}

func SliceFill[T any](s []T, value T, start, end int) {
	if start < 0 || end > len(s) || start > end {
		panic(fmt.Sprintf("SliceFill: invalid range [%d:%d] for slice of length %d", start, end, len(s)))
	}
	if start == end {
		return
	}
	MemsetSlice(s[start:end], value)
}

//go:noescape
func GetG() unsafe.Pointer

//go:noescape
//go:linkname runtime_procPin runtime.procPin
func runtime_procPin() int

//go:noescape
//go:linkname runtime_procUnpin runtime.procUnpin
func runtime_procUnpin()

// Pin pins current p, return pid.
func Pin() int {
	return runtime_procPin()
}

// Unpin unpins current p.
func Unpin() {
	runtime_procUnpin()
}

// Pid returns the id of current p.
func Pid() (id int) {
	id = runtime_procPin()
	runtime_procUnpin()
	return
}

func MallocSlice[T any](len, cap int) []T {
	var t T
	mem, overflow := MulUintptr(unsafe.Sizeof(t), uintptr(cap))
	if overflow || len < 0 || len > cap {
		panic("invalid slice length or capacity")
	}
	return *(*[]T)(unsafe.Pointer(&struct {
		Data uintptr
		Len  int
		Cap  int
	}{uintptr(mallocgc(mem, Pointer(reflect.TypeOf(t)), false)), len, cap}))
}

func GetPrivateField[T any, V any](ptr *T, fieldName string) *V {
	t := reflect.TypeOf(*ptr)

	field, _ := t.FieldByName(fieldName)
	fieldOffset := field.Offset

	return (*V)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + fieldOffset))
}

// pointers must be zeroed
func Malloc[T any](t T) *T {
	// Allocate memory and then copy the value of t into the allocated memory
	ptr := (*T)(mallocgc(unsafe.Sizeof(t), Pointer(reflect.TypeOf(t)), true))
	*ptr = t // Copy the value of t into the allocated memory
	return ptr
}

//go:linkname newarray runtime.newarray
func newarray(t unsafe.Pointer, n int) unsafe.Pointer

func MakeSlice[T any](len, cap int) []T {
	var typ T
	return *(*[]T)(unsafe.Pointer(&struct {
		Data uintptr
		Len  int
		Cap  int
	}{uintptr(newarray(Pointer(typ), cap)), len, cap}))

}

type Iface struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}

func Inspect(v interface{}) (reflect.Type, unsafe.Pointer) {
	return reflect.TypeOf(v), Pointer(v)
}

func Pointer(v interface{}) unsafe.Pointer {
	return (*Iface)(unsafe.Pointer(&v)).ptr
}

type MutableString []byte

func (m *MutableString) String() string {
	return String(*m)
}

func (m *MutableString) SubString(start, end int) string {
	if start < 0 || end > len(*m) || start > end {
		panic("invalid substring range")
	}
	return unsafe.String(unsafe.SliceData((*m)[start:end]), end-start)
}

func (m *MutableString) SubSlice(start, end int) MutableString {
	if start < 0 || end > len(*m) || start > end {
		panic("invalid slice range")
	}
	return (*m)[start:end]
}

func (m *MutableString) Reserve(cp int) {
	if cp > cap(*m) {
		newBuf := MakeNoZeroCap(len(*m), cp)
		copy(newBuf, *m)
		*m = newBuf
	}
}

func (m *MutableString) Shrink() {
	if cap(*m) > len(*m) {
		newBuf := MakeNoZero(len(*m))
		copy(newBuf, *m)
		*m = newBuf
	}
}

func (m *MutableString) WriteUnaligned32(i int, v uint32) {
	if i < 0 || i+4 > len(*m) {
		panic("index out of range")
	}
	b := (*[4]byte)(unsafe.Pointer(&(*m)[i]))
	if isLittleEndian() {
		b[0], b[1], b[2], b[3] = byte(v), byte(v>>8), byte(v>>16), byte(v>>24)
	} else {
		b[3], b[2], b[1], b[0] = byte(v), byte(v>>8), byte(v>>16), byte(v>>24)
	}
}

func (m *MutableString) AppendUnaligned32(v uint32) {
	if len(*m) == 0 {
		*m = MakeNoZeroCap(0, 4)
	}
	b := [4]byte{}
	if isLittleEndian() {
		b[0], b[1], b[2], b[3] = byte(v), byte(v>>8), byte(v>>16), byte(v>>24)
	} else {
		b[3], b[2], b[1], b[0] = byte(v), byte(v>>8), byte(v>>16), byte(v>>24)
	}
	*m = append(*m, b[:]...)
}

func (m *MutableString) WriteToBuffer(b io.Writer) (int, error) {
	return b.Write(*m)
}

// ONLY FOR SHOW EXAMPLE
// USE DEFAULT [3:5] IN PROD
func BytesFromRange(start, end uintptr) []byte {
	length := end - start
	return *(*[]byte)(unsafe.Pointer(&struct {
		uintptr
		int
		i int
	}{start, int(length), int(length)}))
}

func (m *MutableString) StringNoZero() string {
	if len(*m) == 0 {
		return ""
	}
	result := String(*m)
	return strings.TrimRight(result, "\x00")
}
func (m *MutableString) AppendString(s string) {
	if len(*m) == 0 {
		*m = MakeNoZeroCap(0, len(s))
	}
	*m = append(*m, s...)
}

func (m *MutableString) AppendByte(c byte) {
	if len(*m) == 0 {
		*m = MakeNoZeroCap(0, 1)
	}
	*m = append(*m, c)
}

func (m *MutableString) Get(i int) byte {
	return (*m)[i]
}

func (m *MutableString) Append(data []byte) {
	if len(*m) == 0 {
		*m = MakeNoZeroCap(0, len(data))
	}
	*m = append(*m, data...)
}

func (m *MutableString) Len() int {
	return len(*m)
}

func (m *MutableString) Clear() {
	*m = (*m)[:0]
}

func (m *MutableString) SetString(s string) {

	if cap(*m) >= len(s) {
		*m = (*m)[:len(s)]
		copy(*m, s)
	} else {
		*m = []byte(s)
	}
}

func (m *MutableString) IsEmpty() bool {
	return len(*m) == 0
}

func (m *MutableString) ToUpper() {
	*m = []byte(strings.ToUpper(m.String()))
}

func (m *MutableString) ToLower() {
	*m = []byte(strings.ToLower(m.String()))
}

func (m *MutableString) Trim() {
	*m = []byte(strings.TrimSpace(m.String()))
}

func (m *MutableString) Equals(other string) bool {
	return m.String() == other
}

func UnsafePointer[T any](b T) unsafe.Pointer {
	return unsafe.Pointer(&b)
}

func ConvertUnsafePointer[T any](p unsafe.Pointer) T {
	return *(*T)(p)
}

type String_t struct {
	Data unsafe.Pointer
	Len  int
}

// Slice internals from reflect
type Slice_t struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// 10
func ParseInt(b []byte) (int64, error) {
	if len(b) == 0 {
		return 0, ErrInvalidFormat
	}

	result := int64(0)
	ptr := unsafe.Pointer(unsafe.SliceData(b))
	n := len(b)
	i := 0

	for ; i+8 <= n; i += 8 {
		c0 := *(*byte)(unsafe.Add(ptr, uintptr(i)))
		c1 := *(*byte)(unsafe.Add(ptr, uintptr(i+1)))
		c2 := *(*byte)(unsafe.Add(ptr, uintptr(i+2)))
		c3 := *(*byte)(unsafe.Add(ptr, uintptr(i+3)))
		c4 := *(*byte)(unsafe.Add(ptr, uintptr(i+4)))
		c5 := *(*byte)(unsafe.Add(ptr, uintptr(i+5)))
		c6 := *(*byte)(unsafe.Add(ptr, uintptr(i+6)))
		c7 := *(*byte)(unsafe.Add(ptr, uintptr(i+7)))

		valid := (c0 - '0') | (c1 - '0') | (c2 - '0') | (c3 - '0') |
			(c4 - '0') | (c5 - '0') | (c6 - '0') | (c7 - '0')
		if valid > 9 {
			if c0 < '0' || c0 > '9' {
				return 0, ErrInvalidFormat
			}
			if c1 < '0' || c1 > '9' {
				return 0, ErrInvalidFormat
			}
			if c2 < '0' || c2 > '9' {
				return 0, ErrInvalidFormat
			}
			if c3 < '0' || c3 > '9' {
				return 0, ErrInvalidFormat
			}
			if c4 < '0' || c4 > '9' {
				return 0, ErrInvalidFormat
			}
			if c5 < '0' || c5 > '9' {
				return 0, ErrInvalidFormat
			}
			if c6 < '0' || c6 > '9' {
				return 0, ErrInvalidFormat
			}
			if c7 < '0' || c7 > '9' {
				return 0, ErrInvalidFormat
			}
		}

		result = result*100000000 +
			int64(c0-'0')*10000000 +
			int64(c1-'0')*1000000 +
			int64(c2-'0')*100000 +
			int64(c3-'0')*10000 +
			int64(c4-'0')*1000 +
			int64(c5-'0')*100 +
			int64(c6-'0')*10 +
			int64(c7-'0')
	}

	for ; i < n; i++ {
		c := *(*byte)(unsafe.Add(ptr, uintptr(i)))
		if c < '0' || c > '9' {
			return 0, ErrInvalidFormat
		}
		result = result*10 + int64(c-'0')
	}

	return result, nil
}

//go:linkname memmove runtime.memmove
func memmove(dst, src unsafe.Pointer, n uintptr)

// constants in make for  standart copy is more faster,but if no constant we can use CopyUnsafe()
// make([]byte, len(src)) not constant make([]byte, 0) constant
//
//go:nocheckptr
func CopyUnsafe[T any](dst []T, src []T) int {
	if len(dst) == 0 || len(src) == 0 {
		return 0
	}
	if len(src) > len(dst) {
		src = src[:len(dst)]
	}
	memmove(
		unsafe.Pointer(unsafe.SliceData(dst)),
		unsafe.Pointer(unsafe.SliceData(src)),
		uintptr(len(src))*unsafe.Sizeof(src[0]),
	)
	return len(src)
}

// Noescape forces any pointerx not escape to the heap
//
//go:nosplit
//go:nocheckptr
func Noescape(up unsafe.Pointer) unsafe.Pointer {
	x := uintptr(up)
	return unsafe.Pointer(x ^ 0)
}

var alwaysFalse bool
var escapeSink any

// Escape forces any pointers in x to escape to the heap.
func Escape[T any](x T) T {
	if alwaysFalse {
		escapeSink = x
	}
	return x
}

type ErrorSizeUnmatch struct {
	fromLength int
	fromSize   int64

	toSize int64
}

func (err *ErrorSizeUnmatch) Error() string {
	return fmt.Sprintf(
		"size mismatch: source length = '%d',"+
			"source size = '%d', destination size = '%d'",
		err.fromLength, err.fromSize, err.toSize)
}

func string2(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func string1(b []byte) string {
	h := *(*reflect.StringHeader)(unsafe.Pointer(&b))
	h.Data = uintptr(unsafe.Pointer(&b[0]))
	h.Len = len(b)
	return *(*string)(unsafe.Pointer(&h))
}

func string_b(b []byte) string {
	h := *(*String_t)(unsafe.Pointer(&b))
	h.Data = unsafe.Pointer(&b[0])
	h.Len = len(b)
	return *(*string)(unsafe.Pointer(&h))
}

func string3(b []byte) string {
	return *(*string)(unsafe.Pointer(&struct {
		uintptr
		int
	}{*(*uintptr)(unsafe.Pointer(&b)), len(b)}))
}

func String(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func string4(b []byte) string {
	(*reflect.StringHeader)(unsafe.Pointer(&b)).Data = uintptr(unsafe.Pointer(&b[0]))
	(*reflect.StringHeader)(unsafe.Pointer(&b)).Len = len(b)
	return *(*string)(unsafe.Pointer(&b))
}

func unsafeGetBytes(s string) (b []byte) {
	(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	(*reflect.SliceHeader)(unsafe.Pointer(&b)).Cap = len(s)
	(*reflect.SliceHeader)(unsafe.Pointer(&b)).Len = len(s)
	return
}

func unsafeGetBytes_2(s string) []byte {
	const MaxInt32 = 1<<31 - 1
	return (*[MaxInt32]byte)(unsafe.Pointer((*reflect.StringHeader)(
		unsafe.Pointer(&s)).Data))[: len(s)&MaxInt32 : len(s)&MaxInt32]
}

func _stringToBytes_(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func stringTobytes(s string) []byte {
	h := (*(*reflect.SliceHeader)(unsafe.Pointer(&s)))
	h.Len = len(s)
	h.Cap = len(s)
	h.Data = uintptr(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&h))
}

func stringBytes(s string) []byte {
	h := *(*Slice_t)(unsafe.Pointer(&s))
	h.Len = len(s)
	h.Cap = len(s)
	h.Data = unsafe.Pointer(&s)
	return *(*[]byte)(unsafe.Pointer(&h))
}

func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func stringToBytes_(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		uintptr
		int
		i int
	}{*(*uintptr)(unsafe.Pointer(&s)), len(s), len(s)}))
}

func CompareSlice(a, b []byte, length int) bool {
	if length > len(a) || length > len(b) {
		return false
	}
	return CompareImpl(unsafe.Pointer(&a[0]), unsafe.Pointer(&b[0]), length)
}

func CopyString(s string) string {
	if len(s) == 0 {
		return ""
	}
	b := make([]byte, len(s))
	copy(b, s)
	return unsafe.String(&b[0], len(b))
}

func Swap[T any](a, b *T) {
	tmp := *a
	*a = *b
	*b = tmp
}

func ConvertSlice[TFrom, TTo any](from []TFrom) ([]TTo, error) {
	var (
		zeroValFrom TFrom
		zeroValTo   TTo
	)

	maxSize := unsafe.Sizeof(zeroValFrom)
	minSize := unsafe.Sizeof(zeroValTo)

	if minSize > maxSize {
		Swap(&minSize, &maxSize)
	}

	if unsafe.Sizeof(zeroValFrom) == minSize {
		if len(from)*int(minSize)%int(maxSize) != 0 {
			return nil, &ErrorSizeUnmatch{
				fromLength: len(from),
				fromSize:   int64(unsafe.Sizeof(zeroValFrom)),
				toSize:     int64(unsafe.Sizeof(zeroValTo)),
			}
		}

		newLen := len(from) * int(minSize) / int(maxSize)
		newCap := cap(from) * int(minSize) / int(maxSize)

		return *(*[]TTo)(unsafe.Pointer(&struct {
			Data uintptr
			Len  int
			Cap  int
		}{*(*uintptr)(unsafe.Pointer(&from)), newLen, newCap})), nil
	} else {
		if len(from)*int(maxSize)%int(minSize) != 0 {
			return nil, &ErrorSizeUnmatch{
				fromLength: len(from),
				fromSize:   int64(unsafe.Sizeof(zeroValFrom)),
				toSize:     int64(unsafe.Sizeof(zeroValTo)),
			}
		}

		newLen := len(from) * int(maxSize) / int(minSize)
		newCap := cap(from) * int(maxSize) / int(minSize)

		return *(*[]TTo)(unsafe.Pointer(&struct {
			Data uintptr
			Len  int
			Cap  int
		}{*(*uintptr)(unsafe.Pointer(&from)), newLen, newCap})), nil

	}
}

//go:linkname mallocgc runtime.mallocgc
func mallocgc(size uintptr, typ unsafe.Pointer, needzero bool) unsafe.Pointer

func MakeNoZero(l int) []byte {
	return unsafe.Slice((*byte)(mallocgc(uintptr(l), nil, false)), l) //  standart

}

func MakeNoZeroCap(l int, c int) []byte {
	return MakeNoZero(c)[:l]
}

func AlignedAlloc(size, align uintptr) unsafe.Pointer {
	addr := mallocgc(Roundupsize(NextPowerOfTwo(size+align), false), nil, false)
	offset := align - uintptr(addr)%align
	return unsafe.Pointer(uintptr(addr) + offset)
}

type StringBuffer struct {
	buf []byte
}

func NewStringBuffer(cap int) *StringBuffer {
	return &StringBuffer{
		buf: MakeNoZeroCap(0, cap),
	}
}

func (b *StringBuffer) String() string {
	return String(b.buf)
}

func (b *StringBuffer) Bytes() []byte {
	return b.buf
}

func (b *StringBuffer) Len() int {
	return len(b.buf)
}

func (b *StringBuffer) Cap() int {
	return cap(b.buf)
}

func (b *StringBuffer) Reset() {
	b.buf = b.buf[:0] // reuse the underlying storage
}

func (b *StringBuffer) grow(n int) {
	buf := MakeNoZero(2*cap(b.buf) + n)[:len(b.buf)]
	copy(buf, b.buf)
	b.buf = buf
}

func (b *StringBuffer) Grow(n int) {
	// Check if n is negative
	if n < 0 {
		// Panic with the message "fast.StringBuffer.Grow: negative count"
		panic("fast.StringBuffer.Grow: negative count")
	}

	// Check if the buffer's available capacity is less than n
	if cap(b.buf)-len(b.buf) < n {
		// Call the grow method to increase the capacity
		b.grow(n)
	}
}

func (b *StringBuffer) Write(p []byte) (int, error) {
	b.buf = append(b.buf, p...)
	return len(p), nil
}

func (b *StringBuffer) WriteByte(c byte) error {

	b.buf = append(b.buf, c)
	return nil
}

func (b *StringBuffer) WriteRune(r rune) (int, error) {

	n := len(b.buf)
	b.buf = utf8.AppendRune(b.buf, r)
	return len(b.buf) - n, nil
}

func (b *StringBuffer) WriteString(s string) (int, error) {

	b.buf = append(b.buf, s...)
	return len(s), nil
}

func ConvertOne[TFrom, TTo any](from TFrom) (TTo, error) {
	var (
		zeroValFrom TFrom
		zeroValTo   TTo
	)

	if unsafe.Sizeof(zeroValFrom) != unsafe.Sizeof(zeroValTo) { // need same size to convert
		return zeroValTo, &ErrorSizeUnmatch{
			fromSize: int64(unsafe.Sizeof(zeroValFrom)),
			toSize:   int64(unsafe.Sizeof(zeroValTo)),
		}
	}

	value := *(*TTo)(unsafe.Pointer(&from))

	return value, nil
}

func MustConvertOne[TFrom, TTo any](from TFrom) TTo {

	return *(*TTo)(unsafe.Pointer(&from))

}

func MakeZero[T any](l int) []T { // for now better works with big size
	return unsafe.Slice((*T)(mallocgc(uintptr(l), nil, true)), l)
}

// in future i'll try to replace interface

func MakeZeroCap[T any](l int, c int) []T { //  // for now better works with big size
	return MakeZero[T](c)[:l]
}

func MakeNoZeroString(l int) []string {
	return unsafe.Slice((*string)(mallocgc(uintptr(l), nil, false)), l)
}

func MakeNoZeroCapString(l int, c int) []string {
	return MakeNoZeroString(c)[:l]
}

//go:linkname memequal runtime.memequal
func memequal(a, b unsafe.Pointer, size uintptr) bool

func Equal(a, b []byte, length uintptr) bool {
	return memequal(unsafe.Pointer(&a[0]), unsafe.Pointer(&b[0]), length)

}

func IsNil(v any) bool {
	/*
		var x *int
		var y any
		fmt.Println(x == nil) // false
		fmt.Println(isNil(x))           // true
		fmt.Println(x == nil) // true
		fmt.Println(isNil(y))           // panic


	*/

	return reflect.ValueOf(v).IsNil()
}

// IsEqual checks if two variables point to the same memory location.
//
// It uses unsafe.Pointer to get the memory address of the variables.
// The equality check is performed by comparing the memory addresses.
//
// Parameters:
// - v1: The first variable.
// - v2: The second variable.
//
// Returns:
// - bool: True if the variables point to the same memory location, false otherwise.
func IsEqual[T any](v1, v2 T) bool {
	// Get the memory address of the variables using unsafe.Pointer.
	// The & operator returns the memory address of a variable.
	// The unsafe.Pointer type is used to store and manipulate untyped memory.
	// It is commonly used in low-level programming to bypass type safety checks.
	//
	// The &v1 and &v2 expressions take the address of v1 and v2 variables respectively.
	// The expressions return pointers to the variables.
	return unsafe.Pointer(&v1) == unsafe.Pointer(&v2)
}

type CacheLinePadding struct {
	_ [constants.CacheLinePadSize]byte
}

// Example of using cache line padding

type AtomicCounter struct {
	_     CacheLinePadding // 64 or 32
	value atomic.Int32
	_     [constants.CacheLinePadSize - unsafe.Sizeof(atomic.Int32{})]byte
}

func (a *AtomicCounter) Increment(int) {

	a.value.Add(1)
}

func (a *AtomicCounter) Get() int32 {
	return a.value.Load()

}

func BytesToUint64Slice(b []byte) []uint64 {
	if len(b) == 0 {
		return nil
	}
	if len(b)%8 != 0 {
		panic("BytesToUint64Slice: length of byte slice must be a multiple of 8")
	}
	return unsafe.Slice((*uint64)(unsafe.Pointer(unsafe.SliceData(b))), len(b)/8)
}

func GetItem[T any](slice []T, idx int) T { // experimental same performance as original

	if len(slice) == 0 || idx < 0 || idx >= len(slice) {
		panic("index out of range")
	}

	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(&slice[0])) + uintptr(idx)*unsafe.Sizeof(slice[0]))

	return *(*T)(ptr)
}

//go:nocheckptr
func GetItemWithoutCheck[T any](slice []T, idx int) T { // clears the checks for idx and make it faster but not safe

	ptr := (*T)(unsafe.Add(unsafe.Pointer(&slice[0]), uintptr(idx)*unsafe.Sizeof(slice[0])))
	return *ptr
}

var tab64 = [64]uintptr{
	63, 0, 58, 1, 59, 47, 53, 2,
	60, 39, 48, 27, 54, 33, 42, 3,
	61, 51, 37, 40, 49, 18, 28, 20,
	55, 30, 34, 11, 43, 14, 22, 4,
	62, 57, 46, 52, 38, 26, 32, 41,
	50, 36, 17, 19, 29, 10, 13, 21,
	56, 45, 25, 31, 35, 16, 9, 12,
	44, 24, 15, 8, 23, 7, 6, 5,
}

// log2 computes the binary logarithm of x, rounded up to the next integer
func Log2(i uintptr) (n uintptr) {
	if i == 0 {
		return 0
	}

	i |= i >> 1
	i |= i >> 2
	i |= i >> 4
	i |= i >> 8
	i |= i >> 16
	i |= i >> 32

	// Use the lookup table to determine the position of the highest bit.
	return uintptr(tab64[((i-(i>>1))*0x07EDD5E59A4E28C2)>>58])

}

func IsPowerOfTwo(x uintptr) bool {
	return x&(x-1) == 0
}

// Bswap64 returns its input with byte order reversed
// 0x0102030405060708 -> 0x0807060504030201
func Bswap64(x uint64) uint64 {
	c8 := uint64(0x00ff00ff00ff00ff)
	a := x >> 8 & c8
	b := (x & c8) << 8
	x = a | b
	c16 := uint64(0x0000ffff0000ffff)
	a = x >> 16 & c16
	b = (x & c16) << 16
	x = a | b
	c32 := uint64(0x00000000ffffffff)
	a = x >> 32 & c32
	b = (x & c32) << 32
	x = a | b
	return x
}

const (
	minHeapAlign    = 8
	_MaxSmallSize   = 32768
	smallSizeDiv    = 8
	smallSizeMax    = 1024
	largeSizeDiv    = 128
	_NumSizeClasses = 68
	_PageShift      = 13
	maxObjsPerSpan  = 1024
)

var class_to_size = [_NumSizeClasses]uint16{0, 8, 16, 24, 32, 48, 64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 224, 240, 256, 288, 320, 352, 384, 416, 448, 480, 512, 576, 640, 704, 768, 896, 1024, 1152, 1280, 1408, 1536, 1792, 2048, 2304, 2688, 3072, 3200, 3456, 4096, 4864, 5376, 6144, 6528, 6784, 6912, 8192, 9472, 9728, 10240, 10880, 12288, 13568, 14336, 16384, 18432, 19072, 20480, 21760, 24576, 27264, 28672, 32768}
var class_to_allocnpages = [_NumSizeClasses]uint8{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 3, 2, 3, 1, 3, 2, 3, 4, 5, 6, 1, 7, 6, 5, 4, 3, 5, 7, 2, 9, 7, 5, 8, 3, 10, 7, 4}
var class_to_divmagic = [_NumSizeClasses]uint32{0, ^uint32(0)/8 + 1, ^uint32(0)/16 + 1, ^uint32(0)/24 + 1, ^uint32(0)/32 + 1, ^uint32(0)/48 + 1, ^uint32(0)/64 + 1, ^uint32(0)/80 + 1, ^uint32(0)/96 + 1, ^uint32(0)/112 + 1, ^uint32(0)/128 + 1, ^uint32(0)/144 + 1, ^uint32(0)/160 + 1, ^uint32(0)/176 + 1, ^uint32(0)/192 + 1, ^uint32(0)/208 + 1, ^uint32(0)/224 + 1, ^uint32(0)/240 + 1, ^uint32(0)/256 + 1, ^uint32(0)/288 + 1, ^uint32(0)/320 + 1, ^uint32(0)/352 + 1, ^uint32(0)/384 + 1, ^uint32(0)/416 + 1, ^uint32(0)/448 + 1, ^uint32(0)/480 + 1, ^uint32(0)/512 + 1, ^uint32(0)/576 + 1, ^uint32(0)/640 + 1, ^uint32(0)/704 + 1, ^uint32(0)/768 + 1, ^uint32(0)/896 + 1, ^uint32(0)/1024 + 1, ^uint32(0)/1152 + 1, ^uint32(0)/1280 + 1, ^uint32(0)/1408 + 1, ^uint32(0)/1536 + 1, ^uint32(0)/1792 + 1, ^uint32(0)/2048 + 1, ^uint32(0)/2304 + 1, ^uint32(0)/2688 + 1, ^uint32(0)/3072 + 1, ^uint32(0)/3200 + 1, ^uint32(0)/3456 + 1, ^uint32(0)/4096 + 1, ^uint32(0)/4864 + 1, ^uint32(0)/5376 + 1, ^uint32(0)/6144 + 1, ^uint32(0)/6528 + 1, ^uint32(0)/6784 + 1, ^uint32(0)/6912 + 1, ^uint32(0)/8192 + 1, ^uint32(0)/9472 + 1, ^uint32(0)/9728 + 1, ^uint32(0)/10240 + 1, ^uint32(0)/10880 + 1, ^uint32(0)/12288 + 1, ^uint32(0)/13568 + 1, ^uint32(0)/14336 + 1, ^uint32(0)/16384 + 1, ^uint32(0)/18432 + 1, ^uint32(0)/19072 + 1, ^uint32(0)/20480 + 1, ^uint32(0)/21760 + 1, ^uint32(0)/24576 + 1, ^uint32(0)/27264 + 1, ^uint32(0)/28672 + 1, ^uint32(0)/32768 + 1}
var size_to_class8 = [smallSizeMax/smallSizeDiv + 1]uint8{0, 1, 2, 3, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16, 16, 17, 17, 18, 18, 19, 19, 19, 19, 20, 20, 20, 20, 21, 21, 21, 21, 22, 22, 22, 22, 23, 23, 23, 23, 24, 24, 24, 24, 25, 25, 25, 25, 26, 26, 26, 26, 27, 27, 27, 27, 27, 27, 27, 27, 28, 28, 28, 28, 28, 28, 28, 28, 29, 29, 29, 29, 29, 29, 29, 29, 30, 30, 30, 30, 30, 30, 30, 30, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32}
var size_to_class128 = [(_MaxSmallSize-smallSizeMax)/largeSizeDiv + 1]uint8{32, 33, 34, 35, 36, 37, 37, 38, 38, 39, 39, 40, 40, 40, 41, 41, 41, 42, 43, 43, 44, 44, 44, 44, 44, 45, 45, 45, 45, 45, 45, 46, 46, 46, 46, 47, 47, 47, 47, 47, 47, 48, 48, 48, 49, 49, 50, 51, 51, 51, 51, 51, 51, 51, 51, 51, 51, 52, 52, 52, 52, 52, 52, 52, 52, 52, 52, 53, 53, 54, 54, 54, 54, 55, 55, 55, 55, 55, 56, 56, 56, 56, 56, 56, 56, 56, 56, 56, 56, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 58, 58, 58, 58, 58, 58, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 61, 61, 61, 61, 61, 62, 62, 62, 62, 62, 62, 62, 62, 62, 62, 62, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67}

func divRoundUp(n, a uintptr) uintptr {
	// a is generally a power of two. This will get inlined and
	// the compiler will optimize the division.
	return (n + a - 1) / a
}

// Returns size of the memory block that mallocgc will allocate if you ask for the size,
// minus any inline space for metadata.
func Roundupsize(size uintptr, noscan bool) (reqSize uintptr) {
	reqSize = size
	if reqSize <= _MaxSmallSize-8 {
		// Small object.
		if !noscan && reqSize > minSizeForMallocHeader { // !noscan && !heapBitsInSpan(reqSize)
			reqSize += 8
		}
		// (reqSize - size) is either mallocHeaderSize or 0. We need to subtract mallocHeaderSize
		// from the result if we have one, since mallocgc will add it back in.
		if reqSize <= smallSizeMax-8 {
			return uintptr(class_to_size[size_to_class8[divRoundUp(reqSize, smallSizeDiv)]]) - (reqSize - size)
		}
		return uintptr(class_to_size[size_to_class128[divRoundUp(reqSize-smallSizeMax, largeSizeDiv)]]) - (reqSize - size)
	}
	// Large object. Align reqSize up to the next page. Check for overflow.
	reqSize += _PageSize - 1
	if reqSize < size {
		return size
	}
	return reqSize &^ (_PageSize - 1)
}

// Bswap32 returns its input with byte order reversed
// 0x01020304 -> 0x04030201
func Bswap32(x uint32) uint32 {
	c8 := uint32(0x00ff00ff)
	a := x >> 8 & c8
	b := (x & c8) << 8
	x = a | b
	c16 := uint32(0x0000ffff)
	a = x >> 16 & c16
	b = (x & c16) << 16
	x = a | b
	return x
}

func NextPowerOfTwo(i uintptr) uintptr {
	i--
	i |= i >> 1
	i |= i >> 2
	i |= i >> 4
	i |= i >> 8
	i |= i >> 16
	i |= i >> 32
	i++
	return i
}
