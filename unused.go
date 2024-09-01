package lowlevelfunctions

import "unsafe"

func new[T any](value T, size int) *T {
	obj := (*T)(mallocgc(uintptr(size), nil, false))
	*obj = value
	return obj
}

//go:linkname lock runtime.lock
func lock(l *mutex)

//go:linkname nanotime runtime.nanotime
func nanotime() int64

//go:linkname unlock runtime.unlock
func unlock(l *mutex)

type mutex struct {
	// Futex-based impl treats it as uint32 key,
	// while sema-based impl as M* waitm.
	// Used to be a union, but unions break precise GC.
	key uintptr
}

//go:linkname sysFree runtime.sysFree
func sysFree(v unsafe.Pointer, n uintptr, sysStat unsafe.Pointer)

//go:linkname sysFreeOS runtime.sysFreeOS
func sysFreeOS(v unsafe.Pointer, n uintptr)

//go:linkname goReady runtime.goready
func goReady(goroutinePtr unsafe.Pointer, traceskip int)

//go:linkname mCall runtime.mcall
func mCall(fn func(unsafe.Pointer))

//go:linkname readGStatus runtime.readgstatus
func readGStatus(gp unsafe.Pointer) uint32

//go:linkname casGStatus runtime.casgstatus
func casGStatus(gp unsafe.Pointer, oldval, newval uint32)

//go:linkname dropG runtime.dropg
func dropG()

//go:linkname schedule runtime.schedule
func schedule()
