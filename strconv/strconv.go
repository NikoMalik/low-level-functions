package strconv

import (
	"errors"
	"math"
	"unsafe"
)

var (
	ErrOverflow         = errors.New("overflow")
	ErrInvalidCharacter = errors.New("invalid character")
	ErrInvalidString    = errors.New("invalid string")
	ErrEmptyString      = errors.New("empty string")
)

func _string(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func Digits10(v uint64) uint32 {
	if v < 10 {
		return 1
	}
	if v < 100 {
		return 2
	}
	if v < 1000 {
		return 3
	}
	if v < 1_000_000_000_000 {
		if v < 100_000_000 {
			if v < 1_000_000 {
				if v < 10_000 {
					return 4
				}
				return 5 + boolToUint32(v >= 100_000)
			}
			return 7 + boolToUint32(v >= 10_000_000)
		}
		if v < 10_000_000_000 {
			return 9 + boolToUint32(v >= 1_000_000_000)
		}
		return 11 + boolToUint32(v >= 100_000_000_000)
	}
	return 12 + Digits10(v/1_000_000_000_000)
}

var digits = [...]byte{
	'0', '0', '0', '1', '0', '2', '0', '3', '0', '4', '0', '5', '0', '6', '0', '7', '0', '8', '0', '9',
	'1', '0', '1', '1', '1', '2', '1', '3', '1', '4', '1', '5', '1', '6', '1', '7', '1', '8', '1', '9',
	'2', '0', '2', '1', '2', '2', '2', '3', '2', '4', '2', '5', '2', '6', '2', '7', '2', '8', '2', '9',
	'3', '0', '3', '1', '3', '2', '3', '3', '3', '4', '3', '5', '3', '6', '3', '7', '3', '8', '3', '9',
	'4', '0', '4', '1', '4', '2', '4', '3', '4', '4', '4', '5', '4', '6', '4', '7', '4', '8', '4', '9',
	'5', '0', '5', '1', '5', '2', '5', '3', '5', '4', '5', '5', '5', '6', '5', '7', '5', '8', '5', '9',
	'6', '0', '6', '1', '6', '2', '6', '3', '6', '4', '6', '5', '6', '6', '6', '7', '6', '8', '6', '9',
	'7', '0', '7', '1', '7', '2', '7', '3', '7', '4', '7', '5', '7', '6', '7', '7', '7', '8', '7', '9',
	'8', '0', '8', '1', '8', '2', '8', '3', '8', '4', '8', '5', '8', '6', '8', '7', '8', '8', '8', '9',
	'9', '0', '9', '1', '9', '2', '9', '3', '9', '4', '9', '5', '9', '6', '9', '7', '9', '8', '9', '9',
}

// =========================REDIS_IMPL=================================
func Ull2String(dst []byte, value uint64) int {
	dstlen := len(dst)

	length := Digits10(value)
	if int(length) >= dstlen {
		if dstlen > 0 {
			dst[0] = 0
		}
		return 0
	}
	next := length - 1
	dst[next+1] = 0

	for value >= 100 {
		i := (value % 100) * 2
		value /= 100
		dst[next] = digits[i+1]
		dst[next-1] = digits[i]
		next -= 2
	}

	if value < 10 {
		dst[next] = '0' + byte(value)
	} else {
		i := value * 2
		dst[next] = digits[i+1]
		dst[next-1] = digits[i]
	}

	return int(length)
}

func Ll2String(dst []byte, svalue int64) int {
	dstlen := len(dst)
	negative := 0
	var value uint64

	if svalue < 0 {
		if svalue != math.MinInt64 {
			value = uint64(-svalue)
		} else {
			value = uint64(math.MaxInt64) + 1
		}
		if dstlen < 2 {
			if dstlen > 0 {
				dst[0] = 0
			}
			return 0
		}
		negative = 1
		dst[0] = '-'
		dst = dst[1:]
		dstlen--
	} else {
		value = uint64(svalue)
	}

	length := Ull2String(dst, value)
	if length == 0 {
		return 0
	}

	return length + negative
}

func ParseUint64(s string) (uint64, error) {
	var v uint64
	if len(s) == 0 {
		return 0, ErrEmptyString
	}

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, ErrInvalidCharacter
		}

		d := uint64(c - '0')
		if v > (math.MaxUint64-d)/10 {
			return 0, ErrOverflow
		}
		v = v*10 + d
	}
	return v, nil
}

func ParseInt64(s string) (int64, error) {
	if len(s) == 0 {
		return 0, ErrEmptyString
	}

	negative := false
	start := 0
	if s[0] == '-' {
		negative = true
		start = 1
		if len(s) == 1 {
			return 0, ErrInvalidString
		}
	}

	var v uint64
	for i := start; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, ErrInvalidCharacter
		}

		d := uint64(c - '0')
		if v > (math.MaxInt64-uint64(d))/10 {
			if negative && v == (math.MaxInt64+1-d)/10 {
				v = v*10 + d
				break
			}
			return 0, ErrOverflow
		}
		v = v*10 + d
	}

	if negative {
		if v > math.MaxInt64+1 {
			return 0, ErrOverflow
		}
		return -int64(v), nil
	}
	if v > math.MaxInt64 {
		return 0, ErrOverflow
	}
	return int64(v), nil
}

func boolToUint32(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

// ------------------------
const SAFETY_BUF_SIZE = 32
const FAST_BUF_SIZE = 22

func U64ToA(x uint64, data []byte) int {
	if x < 10 {
		data[0] = '0' + byte(x)
		return 1
	}

	i := 0
	for {
		data[i] = '0' + byte(x%10)
		i++
		x /= 10
		if x == 0 {
			break
		}
	}
	for j, k := 0, i-1; j < k; j, k = j+1, k-1 {
		data[j], data[k] = data[k], data[j]
	}

	return i
}

func I64ToA(x int64, data []byte) int {
	// handle negative
	if x < 0 {
		if x == math.MaxInt64 {
			copy(data, []byte("-9223372036854775808"))
			return 20
		}
		data[0] = '-'
		return 1 + U64ToA(uint64(-x), data[1:])
	}

	return U64ToA(uint64(x), data)
}

func FormatUint16(u uint16) string {

	if u == 0 {
		return "0"
	}

	var result [5]byte
	i := 4
loop:
	if u > 0 {
		result[i] = byte(u%10) + '0'
		u /= 10
		i--
		goto loop
	}

	return _string(result[i+1:])
}

/*
uint8  : 0 to 255
uint16 : 0 to 65535
uint32 : 0 to 4294967295
uint64 : 0 to 18446744073709551615
int8   : -128 to 127
int16  : -32768 to 32767
int32  : -2147483648 to 2147483647
int64  : -9223372036854775808 to 9223372036854775807
*/
