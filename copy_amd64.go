//go:build amd64 && !purego

package lowlevelfunctions

import "unsafe"

//go:noescape
func copy_AVX2_32(src []byte, src2 []byte) int

//go:noescape
func copy_AVX2_64(src []byte, src2 []byte) int

//go:noescape
func memcopy_avx2_32(src unsafe.Pointer, src2 unsafe.Pointer) int

//go:noescape
func memcopy_avx2_64(src unsafe.Pointer, src2 unsafe.Pointer) int
