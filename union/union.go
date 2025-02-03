package union

import (
	"errors"
	"unsafe"
)

var (
	ErrTypeTooLarge = errors.New("type size exceeds union capacity")
	ErrInvalidType  = errors.New("unsupported type")
)

type Union[T any] struct {
	data [8]byte
}

func (u *Union[T]) Reset()         { u.data = [8]byte{} }
func (u *Union[T]) Data() *[8]byte { return &u.data }

func NewUnion[T any]() *Union[T] { return &Union[T]{} }

func (u *Union[T]) Set(val T) error {
	size := unsafe.Sizeof(val)
	if size > uintptr(len(u.data)) {
		return ErrTypeTooLarge
	}

	src := unsafe.Slice((*byte)(unsafe.Pointer(&val)), size)
	copy(u.data[:], src)
	return nil
}

func (u *Union[T]) Get() (T, error) {
	var zero T
	size := unsafe.Sizeof(zero)
	if size > uintptr(len(u.data)) {
		return zero, ErrTypeTooLarge
	}
	if u.data == [8]byte{} {
		return zero, ErrInvalidType

	}

	return *(*T)(unsafe.Pointer(&u.data[0])), nil
}

func (u *Union[T]) SetInt64(val int64) { u.data = *(*[8]byte)(unsafe.Pointer(&val)) }
func (u *Union[T]) GetInt64() int64    { return *(*int64)(unsafe.Pointer(&u.data)) }

func (u *Union[T]) SetFloat64(val float64) { u.data = *(*[8]byte)(unsafe.Pointer(&val)) }
func (u *Union[T]) GetFloat64() float64    { return *(*float64)(unsafe.Pointer(&u.data)) }
