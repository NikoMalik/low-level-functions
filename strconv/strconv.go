package strconv

import (
	lowlevelfunctions "github.com/NikoMalik/low-level-functions"
)

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

	return lowlevelfunctions.String(result[i+1:])
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

// TODO: ADD BINARY.LITTLEENDIAN && BINARY.BIGENDIAN FOR LOW-LEVEL FUNCTIONS
type littleEndian struct{}

func (littleEndian) put_uint16(b *[2]byte, value uint16) {
	b[0] = byte(value)
	b[1] = byte(value >> 8)
}
