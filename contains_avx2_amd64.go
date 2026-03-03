//go:build !amd64

package lowlevelfunctions

import (
	"bytes"
	"unsafe"
)

func containsStringAVX2(
	slicePtr unsafe.Pointer,
	sliceLen int,
	valuePtr unsafe.Pointer,
	valueLen int,
) bool {
	if sliceLen == 0 || valueLen == 0 {
		return false
	}

	slices := *(*[]string)(unsafe.Pointer(&slicePtr))

	value := *(*stirng)(unsafe.Pointer(&valuePtr))

	for i := 0; i < sliceLen; i++ {
		if slices[i] == value {
			return true
		}
	}
	return false

}

func containsByteSliceAVX2(slicePtr unsafe.Pointer, sliceLen int, valuePtr unsafe.Pointer, valueLen int) bool {
	if sliceLen == 0 || valueLen == 0 {
		return false
	}

	slices := *(*[][]byte)(unsafe.Pointer(&slicePtr))

	value := *(*[]byte)(unsafe.Pointer(&valuePtr))

	for i := 0; i < sliceLen; i++ {
		if bytes.Equal(slices[i], value) {
			return true
		}
	}
	return false
}
