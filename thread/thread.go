package thread

import (
	"fmt"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

//go:noescape
func rawClone(flags uintptr, stack unsafe.Pointer) int

func thread() {
	// stack for thread 2mb
	runtime.GOMAXPROCS(1)
	const stackSize = 2 << 20
	stack := [stackSize]byte{}
	stackPtr := unsafe.Pointer(&stack[0])
	stackPtr = unsafe.Pointer(uintptr(stackPtr) & ^uintptr(0xF))
	// flags for clone: CLONE_VM | CLONE_FS | CLONE_FILES | CLONE_SIGHAND | CLONE_THREAD
	flags := uintptr(syscall.CLONE_VM | syscall.CLONE_FS | syscall.CLONE_FILES | syscall.CLONE_SIGHAND | syscall.CLONE_THREAD)
	//create thread
	fmt.Printf("Current PID: %d\n", syscall.Getpid())
	time.Sleep(time.Second * 10)

	tid := rawClone(flags, stackPtr)

	if tid < 0 {
		fmt.Println("Error:", syscall.Errno(-tid))
		return
	}
	fmt.Printf("Thread created with TID: %d\n", tid)

	for {
		time.Sleep(1 * time.Second)
		fmt.Printf("Thread running in TID: %d\n", syscall.Gettid())
		tid2 := rawClone(flags, stackPtr)
		fmt.Printf("Thread created with TID: %d\n", tid2)
	}
}
