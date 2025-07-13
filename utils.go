package lowlevelfunctions

import (
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"
	"unicode/utf8"
	"unsafe"

	"github.com/NikoMalik/low-level-functions/constants"
)

const (
	PtrSize       = 4 << (^uintptr(0) >> 63)
	StrSize       = unsafe.Sizeof("")
	SliceSize     = int(unsafe.Sizeof([]byte{}))
	CacheLineSize = constants.CacheLinePadSize
	MaxInt32      = 1<<31 - 1
	MaxUintptr    = ^uintptr(0)
)

func MulUintptr(a, b uintptr) (uintptr, bool) {
	if a|b < 1<<(4*PtrSize) || a == 0 {
		return a * b, false
	}
	overflow := b > MaxUintptr/a
	return a * b, overflow
}

func isLittleEndian() bool {
	var i uint16 = 0x0102
	b := (*[2]byte)(unsafe.Pointer(&i))
	return b[0] == 0x02
}
func Uint32_littleEndian(b *[4]byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func Uint32_bigEndian(b *[4]byte) uint32 {
	return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
}

func Uint64_littleEndian(b *[8]byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func Uint64_bigEndian(b *[8]byte) uint64 {
	return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
}

func ReadUnaligned32(p unsafe.Pointer) uint32 {
	b := (*[4]byte)(p)
	if !isLittleEndian() {
		return Uint32_bigEndian(b)
	}
	return Uint32_littleEndian(b)
}

// //go:linkname readUnaligned64 runtime.readUnaligned64
// func readUnaligned64(p unsafe.Pointer) uint64
func ReadUnaligned64(p unsafe.Pointer) uint64 {
	b := (*[8]byte)(p)
	if !isLittleEndian() {
		return Uint64_bigEndian(b)
	}
	return Uint64_littleEndian(b)
}

//go:noescape
func GetG() unsafe.Pointer

//go:noescape
//go:linkname runtime_procPin runtime.procPin
func runtime_procPin() int

//go:noescape
//go:linkname runtime_procUnpin runtime.procUnpin
func runtime_procUnpin()

// Pin pins current p, return pid.
func Pin() int {
	return runtime_procPin()
}

// Unpin unpins current p.
func Unpin() {
	runtime_procUnpin()
}

// Pid returns the id of current p.
func Pid() (id int) {
	id = runtime_procPin()
	runtime_procUnpin()
	return
}

func MallocSlice[T any](len, cap int) []T {
	var t T
	mem, overflow := MulUintptr(unsafe.Sizeof(t), uintptr(cap))
	if overflow || len < 0 || len > cap {
		panic("invalid slice length or capacity")
	}
	return *(*[]T)(unsafe.Pointer(&struct {
		Data uintptr
		Len  int
		Cap  int
	}{uintptr(mallocgc(mem, Pointer(reflect.TypeOf(t)), false)), len, cap}))
}

func GetPrivateField[T any, V any](ptr *T, fieldName string) *V {
	t := reflect.TypeOf(*ptr)

	field, _ := t.FieldByName(fieldName)
	fieldOffset := field.Offset

	return (*V)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + fieldOffset))
}

func Malloc[T any](t T) *T {
	// Allocate memory and then copy the value of t into the allocated memory
	ptr := (*T)(mallocgc(unsafe.Sizeof(t), Pointer(reflect.TypeOf(t)), false))
	*ptr = t // Copy the value of t into the allocated memory
	return ptr
}

//go:linkname newarray runtime.newarray
func newarray(t unsafe.Pointer, n int) unsafe.Pointer

func MakeSlice[T any](len, cap int) []T {
	var typ T
	return *(*[]T)(unsafe.Pointer(&struct {
		Data uintptr
		Len  int
		Cap  int
	}{uintptr(newarray(Pointer(typ), cap)), len, cap}))

}

type Iface struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}

func Inspect(v interface{}) (reflect.Type, unsafe.Pointer) {
	return reflect.TypeOf(v), Pointer(v)
}

func Pointer(v interface{}) unsafe.Pointer {
	return (*Iface)(unsafe.Pointer(&v)).ptr
}

type MutableString []byte

func (m *MutableString) String() string {
	if len(*m) == 0 {
		return ""
	}
	return String(*m)
}

// ONLY FOR SHOW EXAMPLE
// USE DEFAULT [3:5] IN PROD
func BytesFromRange(start, end uintptr) []byte {
	length := end - start
	return *(*[]byte)(unsafe.Pointer(&struct {
		uintptr
		int
		i int
	}{start, int(length), int(length)}))
}

func (m *MutableString) StringNoZero() string {
	if len(*m) == 0 {
		return ""
	}
	result := String(*m)
	return strings.TrimRight(result, "\x00")
}
func (m *MutableString) AppendString(s string) {
	if len(*m) == 0 {
		*m = MakeNoZeroCap(0, len(s))
	}
	*m = append(*m, s...)
}

func (m *MutableString) AppendByte(c byte) {
	if len(*m) == 0 {
		*m = MakeNoZeroCap(0, 1)
	}
	*m = append(*m, c)
}

func (m *MutableString) Get(i int) byte {
	return (*m)[i]
}

func (m *MutableString) Append(data []byte) {
	if len(*m) == 0 {
		*m = MakeNoZeroCap(0, len(data))
	}
	*m = append(*m, data...)
}

func (m *MutableString) Len() int {
	return len(*m)
}

func (m *MutableString) Clear() {
	*m = (*m)[:0]
}

func (m *MutableString) SetString(s string) {

	if cap(*m) >= len(s) {
		*m = (*m)[:len(s)]
		copy(*m, s)
	} else {
		*m = []byte(s)
	}
}

func (m *MutableString) IsEmpty() bool {
	return len(*m) == 0
}

func (m *MutableString) ToUpper() {
	*m = []byte(strings.ToUpper(m.String()))
}

func (m *MutableString) ToLower() {
	*m = []byte(strings.ToLower(m.String()))
}

func (m *MutableString) Trim() {
	*m = []byte(strings.TrimSpace(m.String()))
}

func (m *MutableString) Equals(other string) bool {
	return m.String() == other
}

func UnsafePointer[T any](b T) unsafe.Pointer {
	return unsafe.Pointer(&b)
}

func ConvertUnsafePointer[T any](p unsafe.Pointer) T {
	return *(*T)(p)
}

type String_t struct {
	Data unsafe.Pointer
	Len  int
}

// Slice internals from reflect
type Slice_t struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

//go:linkname memmove runtime.memmove
func memmove(dst, src unsafe.Pointer, n uintptr)

// constants in make for  standart copy is more faster,but if no constant we can use CopyUnsafe()
// make([]byte, len(src)) not constant make([]byte, 0) constant
//
//go:nocheckptr
func CopyUnsafe(dst []byte, src []byte) int {
	memmove(unsafe.Pointer(&dst[0]), unsafe.Pointer(&src[0]), uintptr(len(src)))
	return len(src)
}

//go:nosplit
//go:nocheckptr
func Noescape(up unsafe.Pointer) unsafe.Pointer {
	x := uintptr(up)
	return unsafe.Pointer(x ^ 0)
}

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

func string2(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func string1(b []byte) string {
	h := *(*reflect.StringHeader)(unsafe.Pointer(&b))
	h.Data = uintptr(unsafe.Pointer(&b[0]))
	h.Len = len(b)
	return *(*string)(unsafe.Pointer(&h))
}

func string_b(b []byte) string {
	h := *(*String_t)(unsafe.Pointer(&b))
	h.Data = unsafe.Pointer(&b[0])
	h.Len = len(b)
	return *(*string)(unsafe.Pointer(&h))
}

func String(b []byte) string {
	return *(*string)(unsafe.Pointer(&struct {
		uintptr
		int
	}{*(*uintptr)(unsafe.Pointer(&b)), len(b)}))
}

func string3(b []byte) string {

	return unsafe.String(unsafe.SliceData(b), len(b))
}

func string4(b []byte) string {
	(*reflect.StringHeader)(unsafe.Pointer(&b)).Data = uintptr(unsafe.Pointer(&b[0]))
	(*reflect.StringHeader)(unsafe.Pointer(&b)).Len = len(b)
	return *(*string)(unsafe.Pointer(&b))
}

func unsafeGetBytes(s string) (b []byte) {
	(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	(*reflect.SliceHeader)(unsafe.Pointer(&b)).Cap = len(s)
	(*reflect.SliceHeader)(unsafe.Pointer(&b)).Len = len(s)
	return
}

func unsafeGetBytes_2(s string) []byte {
	const MaxInt32 = 1<<31 - 1
	return (*[MaxInt32]byte)(unsafe.Pointer((*reflect.StringHeader)(
		unsafe.Pointer(&s)).Data))[: len(s)&MaxInt32 : len(s)&MaxInt32]
}

func _stringToBytes_(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func stringTobytes(s string) []byte {
	h := (*(*reflect.SliceHeader)(unsafe.Pointer(&s)))
	h.Len = len(s)
	h.Cap = len(s)
	h.Data = uintptr(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&h))
}

func stringBytes(s string) []byte {
	h := *(*Slice_t)(unsafe.Pointer(&s))
	h.Len = len(s)
	h.Cap = len(s)
	h.Data = unsafe.Pointer(&s)
	return *(*[]byte)(unsafe.Pointer(&h))
}

func stringToBytes_(s string) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&s)), len(s))
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		uintptr
		int
		i int
	}{*(*uintptr)(unsafe.Pointer(&s)), len(s), len(s)}))
}

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
		Swap(&minSize, &maxSize)
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
func Swap[T any](a, b *T) {
	tmp := *a
	*a = *b
	*b = tmp
}

//go:linkname mallocgc runtime.mallocgc
func mallocgc(size uintptr, typ unsafe.Pointer, needzero bool) unsafe.Pointer

func MakeNoZero(l int) []byte {
	return unsafe.Slice((*byte)(mallocgc(uintptr(l), nil, false)), l) //  standart

}

func MakeNoZeroCap(l int, c int) []byte {
	return MakeNoZero(c)[:l]
}

type StringBuffer struct {
	buf []byte
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

func MakeZero[T any](l int) []T { // for now better works with big size
	return unsafe.Slice((*T)(mallocgc(uintptr(l), nil, true)), l)
}

// in future i'll try to replace interface

func MakeZeroCap[T any](l int, c int) []T { //  // for now better works with big size
	return MakeZero[T](c)[:l]
}

func MakeNoZeroString(l int) []string {
	return unsafe.Slice((*string)(mallocgc(uintptr(l), nil, false)), l)
}

func MakeNoZeroCapString(l int, c int) []string {
	return MakeNoZeroString(c)[:l]
}

//go:linkname memequal runtime.memequal
func memequal(a, b unsafe.Pointer, size uintptr) bool

func Equal(a, b []byte, length uintptr) bool {
	return memequal(unsafe.Pointer(&a[0]), unsafe.Pointer(&b[0]), length)

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
func IsEqual[T any](v1, v2 T) bool {
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
	_     CacheLinePadding // 64 or 32
	value atomic.Int32
	_     [constants.CacheLinePadSize - unsafe.Sizeof(atomic.Int32{})]byte
}

func (a *AtomicCounter) Increment(int) {

	a.value.Add(1)
}

func (a *AtomicCounter) Get() int32 {
	return a.value.Load()

}

func GetItem[T any](slice []T, idx int) T { // experimental same performance as original

	if len(slice) == 0 || idx < 0 || idx >= len(slice) {
		panic("index out of range")
	}

	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(&slice[0])) + uintptr(idx)*unsafe.Sizeof(slice[0]))

	return *(*T)(ptr)
}

//go:nocheckptr
func GetItemWithoutCheck[T any](slice []T, idx int) T { // clears the checks for idx and make it faster but not safe

	ptr := (*T)(unsafe.Add(unsafe.Pointer(&slice[0]), uintptr(idx)*unsafe.Sizeof(slice[0])))
	return *ptr
}

var tab64 = [64]uintptr{
	63, 0, 58, 1, 59, 47, 53, 2,
	60, 39, 48, 27, 54, 33, 42, 3,
	61, 51, 37, 40, 49, 18, 28, 20,
	55, 30, 34, 11, 43, 14, 22, 4,
	62, 57, 46, 52, 38, 26, 32, 41,
	50, 36, 17, 19, 29, 10, 13, 21,
	56, 45, 25, 31, 35, 16, 9, 12,
	44, 24, 15, 8, 23, 7, 6, 5,
}

// log2 computes the binary logarithm of x, rounded up to the next integer
func Log2(i uintptr) (n uintptr) {
	if i == 0 {
		return 0
	}

	i |= i >> 1
	i |= i >> 2
	i |= i >> 4
	i |= i >> 8
	i |= i >> 16
	i |= i >> 32

	// Use the lookup table to determine the position of the highest bit.
	return uintptr(tab64[((i-(i>>1))*0x07EDD5E59A4E28C2)>>58])

}

func NextPowerOfTwo(i uintptr) uintptr {
	i--
	i |= i >> 1
	i |= i >> 2
	i |= i >> 4
	i |= i >> 8
	i |= i >> 16
	i |= i >> 32
	i++
	return i
}
