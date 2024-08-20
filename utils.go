package lowlevelfunctions

import (
	"fmt"
	"reflect"
	"sync/atomic"
	"unicode/utf8"
	"unsafe"

	"github.com/NikoMalik/low-level-functions/constants"
)

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

func String(b []byte) string {

	return unsafe.String(unsafe.SliceData(b), len(b))
}

//go:linkname schedule runtime.schedule
func schedule()

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		string
		Cap int
	}{s, len(s)},
	))
}

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

func CopyString(s string) string {
	c := MakeNoZero(len(s))
	copy(c, StringToBytes(s))
	return String(c)
}

func ConvertSlice[TFrom, TTo any](from []TFrom) ([]TTo, error) {
	var (
		zeroValFrom TFrom
		zeroValTo   TTo
	)

	maxSize := unsafe.Sizeof(zeroValFrom)
	minSize := unsafe.Sizeof(zeroValTo)

	if minSize > maxSize {
		swap(&minSize, &maxSize)
	}

	if unsafe.Sizeof(zeroValFrom) == minSize {
		if len(from)*int(minSize)%int(maxSize) != 0 {
			return nil, &ErrorSizeUnmatch{
				fromLength: len(from),
				fromSize:   int64(unsafe.Sizeof(zeroValFrom)),
				toSize:     int64(unsafe.Sizeof(zeroValTo)),
			}
		}

		header := *(*reflect.SliceHeader)(unsafe.Pointer(&from))
		header.Len = header.Len * int(minSize) / int(maxSize)
		header.Cap = header.Cap * int(minSize) / int(maxSize)
		result := *(*[]TTo)(unsafe.Pointer(&header))

		return result, nil
	} else {
		if len(from)*int(maxSize)%int(minSize) != 0 {
			return nil, &ErrorSizeUnmatch{
				fromLength: len(from),
				fromSize:   int64(unsafe.Sizeof(zeroValFrom)),
				toSize:     int64(unsafe.Sizeof(zeroValTo)),
			}
		}

		header := *(*reflect.SliceHeader)(unsafe.Pointer(&from))
		header.Len = header.Len * int(maxSize) / int(minSize)
		header.Cap = header.Cap * int(maxSize) / int(minSize)
		result := *(*[]TTo)(unsafe.Pointer(&header))

		return result, nil
	}
}

//go:noinline
func swap[T any](a, b *T) {
	tmp := *a
	*a = *b
	*b = tmp
}

//go:linkname mallocgc runtime.mallocgc
func mallocgc(size uintptr, typ unsafe.Pointer, needzero bool) unsafe.Pointer

//go:linkname sysFree runtime.sysFree
func sysFree(v unsafe.Pointer, n uintptr, sysStat unsafe.Pointer)

//go:linkname sysFreeOS runtime.sysFreeOS
func sysFreeOS(v unsafe.Pointer, n uintptr)

type mutex struct {
	// Futex-based impl treats it as uint32 key,
	// while sema-based impl as M* waitm.
	// Used to be a union, but unions break precise GC.
	key uintptr
}

//go:linkname lock runtime.lock
func lock(l *mutex)

//go:linkname nanotime runtime.nanotime
func nanotime() int64

//go:linkname unlock runtime.unlock
func unlock(l *mutex)

func MakeNoZero(l int) []byte {
	return unsafe.Slice((*byte)(mallocgc(uintptr(l), nil, false)), l) //  standart

}

func MakeNoZeroCap(l int, c int) []byte {
	return MakeNoZero(c)[:l]
}

type StringBuffer struct {
	buf  []byte
	addr *StringBuffer
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

//go:linkname noescape runtime.noescape
func noescape(p unsafe.Pointer) unsafe.Pointer

func (b *StringBuffer) copyCheck() {
	if b.addr == nil {

		b.addr = (*StringBuffer)(noescape(unsafe.Pointer(b)))
	} else if b.addr != b {
		panic("strings: illegal use of non-zero Builder copied by value")
	}
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

func MakeNoZeroString(l int) []string {
	return unsafe.Slice((*string)(mallocgc(uintptr(l), nil, false)), l)
}

func MakeNoZeroCapString(l int, c int) []string {
	return MakeNoZeroString(c)[:l]
}

//go:linkname memequal runtime.memequal
func memequal(a, b unsafe.Pointer, size uintptr) bool

func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	return memequal(unsafe.Pointer(&a[0]), unsafe.Pointer(&b[0]), uintptr(len(a)))

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
func IsEqual(v1, v2 any) bool {
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
	value atomic.Int32
	_     [constants.CacheLinePadSize - unsafe.Sizeof(atomic.Int32{})]byte
}

func (a *AtomicCounter) Increment(int) {

	a.value.Add(1)
}

func (a *AtomicCounter) Get() int32 {
	return a.value.Load()

}

type AtomicCounterWithoutPad struct {
	value atomic.Int32
}

func (a *AtomicCounterWithoutPad) Increment(int) {

	a.value.Add(1)
}

func (a *AtomicCounterWithoutPad) Get() int32 {
	return a.value.Load()

}

type ShardedAtomicCounter struct {
	shards [10]AtomicCounter
}

func (a *ShardedAtomicCounter) Increment(id int) {
	a.shards[id].value.Add(1)
}

func (a *ShardedAtomicCounter) Get() int32 {
	var value int32
	for i := 0; i < 10; i++ {
		value += a.shards[i].Get()

	}
	return value
}

type ShardedAtomicCounterWithoutPad struct {
	shards [10]AtomicCounter
}

func (a *ShardedAtomicCounterWithoutPad) Increment(id int) {
	a.shards[id].value.Add(1)
}

func (a *ShardedAtomicCounterWithoutPad) Get() int32 {
	var value int32
	for i := 0; i < 10; i++ {
		value += a.shards[i].Get()

	}
	return value
}

func GetItem[T any](slice []T, idx int) T { // experimental same performance as original

	if len(slice) == 0 || idx < 0 || idx >= len(slice) {
		panic("index out of range")
	}

	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(&slice[0])) + uintptr(idx)*unsafe.Sizeof(slice[0]))

	return *(*T)(ptr)
}
