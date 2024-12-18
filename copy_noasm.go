//go:build !amd64

package lowlevelfunctions

import "unsafe"

func memcopy_avx2_32(src unsafe.Pointer, src2 unsafe.Pointer) int {
	panic("memcopy_avx2_32 not implemented in your system")

}

func memcopy_avx2_64(src unsafe.Pointer, src2 unsafe.Pointer) int {
	panic("memcopy_avx2_32 not implemented in your system")
}

func copy_AVX2_64(src []byte, src2 []byte) int {
	panic("copy_AVX2_64 not implemented in your system")
}

func copy_AMD_AVX2_32(src []byte, src2 []byte) int {
	panic("copy_AMD_AVX2_32 not implemented in your system")
}
