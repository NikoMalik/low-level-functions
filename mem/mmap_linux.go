package mem

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func SysAlloc(n int) (unsafe.Pointer, int) {

	p, err := unix.Mmap(-1, 0, n, unix.PROT_NONE, unix.MAP_ANON|unix.MAP_PRIVATE)
	if err != nil {
		return nil, -1
	}

	return unsafe.Pointer(&p), len(p)

}

func SysFree(m []byte) error {
	return unix.Munmap(m)
}
