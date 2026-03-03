//go:build !amd64

package lowlevelfunctions

import "unsafe"

func SumUnroll8(arr []int) int {
	n := len(arr)
	if n < 8 {
		return 0
	}
	sum := 0
	size := unsafe.Sizeof(arr[0])
	base := unsafe.Pointer(&arr[0])
	for i := 0; i+8 <= n; i += 8 {
		sum += *(*int)(unsafe.Add(base, uintptr(i)*size)) +
			*(*int)(unsafe.Add(base, uintptr(i+1)*size)) +
			*(*int)(unsafe.Add(base, uintptr(i+2)*size)) +
			*(*int)(unsafe.Add(base, uintptr(i+3)*size)) +
			*(*int)(unsafe.Add(base, uintptr(i+4)*size)) +
			*(*int)(unsafe.Add(base, uintptr(i+5)*size)) +
			*(*int)(unsafe.Add(base, uintptr(i+6)*size)) +
			*(*int)(unsafe.Add(base, uintptr(i+7)*size))
	}
	for i := (n / 8) * 8; i < n; i++ {
		sum += *(*int)(unsafe.Add(base, uintptr(i)*size))
	}
	return sum
}
