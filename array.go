package ds

// #cgo CFLAGS: -I${SRCDIR}/c/include
// #include <josestg/ds/array.h>
import "C"

import (
	"fmt"
	"iter"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

type Array[T any] struct {
	backend   C.array_t
	gcEnabled bool
}

// NewArray creates new array without any GC.
// You should call Free() manually when you done.
func NewArray[T any](length int) *Array[T] {
	return NewGCArray[T](length, false)
}

// NewGCArray creates a new array and register the Free method to GC,
// so no need to call free manually.
func NewGCArray[T any](length int, enabled bool) *Array[T] {
	var zero T
	elemSize := unsafe.Sizeof(zero)
	if elemSize == 0 {
		panic("ds.Array: zero-sized types are not supported")
	}

	// *T currently not supported.
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Pointer {
		panic("ds.Array: pointer types are not supported as array elements")
	}

	var backend C.array_t
	s := C.array_init(
		&backend,
		C.size_t(length),
		C.size_t(elemSize),
	)
	mustOk(s)
	a := &Array[T]{
		backend:   backend,
		gcEnabled: enabled,
	}
	if a.gcEnabled {
		runtime.SetFinalizer(a, func(a *Array[T]) { a.Free() })
	}
	return a
}

// Free frees the backing array and set the length to zero.
func (a *Array[T]) Free() {
	// ensure free idempotency.
	if a.backend.head != nil {
		s := C.array_deinit(&a.backend)
		mustOk(s)
		if a.gcEnabled {
			// remove the finalizer once it freed.
			runtime.SetFinalizer(a, nil)
		}
	}
}

func (a *Array[T]) Len() int {
	var length C.size_t
	mustOk(C.array_len(&a.backend, &length))
	return int(length)
}

func (a *Array[T]) Set(index int, value T) {
	a.boundCheck(index)
	s := C.array_set(
		&a.backend,
		C.size_t(index),
		unsafe.Pointer(&value),
		C.size_t(unsafe.Sizeof(value)),
	)
	mustOk(s)
}

func (a *Array[T]) Get(index int) T {
	a.boundCheck(index)
	var out T
	s := C.array_get(
		&a.backend,
		C.size_t(index),
		unsafe.Pointer(&out),
		C.size_t(unsafe.Sizeof(out)),
	)
	mustOk(s)
	return out
}

func (a *Array[T]) Iter(reversed bool) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		if reversed {
			a.iterBackward(yield)
		} else {
			a.iterForward(yield)
		}
	}
}

func (a *Array[T]) iterForward(yield func(int, T) bool) {
	for i := range a.Len() {
		v := a.Get(i)
		if !yield(i, v) {
			break
		}
	}
}

func (a *Array[T]) iterBackward(yield func(int, T) bool) {
	for i := range a.Len() {
		v := a.Get(i)
		if !yield(i, v) {
			break
		}
	}
}

func (a *Array[T]) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range a.iterForward {
		if i > 0 {
			sb.WriteRune(' ')
		}
		_, err := fmt.Fprintf(&sb, "%v", v)
		if err != nil {
			panic(fmt.Errorf("ds.Array: to string at index %d: %v", i, err))
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func (a *Array[T]) boundCheck(index int) {
	n := a.Len()
	if index < 0 || index >= n {
		panic("index out of range")
	}
}

func mustOk(s C.status_t) {
	if err := errorOf(s); err != nil {
		panic(err)
	}
}

func errorOf(s C.status_t) error {
	switch s {
	default:
		return fmt.Errorf("ds.Array: status_t(%v): unrecognized status", s)
	case C.S_OK:
		return nil
	case C.S_ERR_SELF_IS_NULL:
		return fmt.Errorf("ds.Array: status_t(%v): self is null", s)
	case C.S_ERR_RETURN_PARAMS_IS_NULL:
		return fmt.Errorf("ds.Array: status_t(%v): out params is missing", s)
	case C.S_ERR_OUT_OF_MEMORY:
		return fmt.Errorf("ds.Array: status_t(%v): out of memory", s)
	case C.S_ERR_OUT_OF_RANGE:
		return fmt.Errorf("ds.Array: status_t(%v): index out of range", s)
	case C.S_ERR_INVALID_ARGUMENTS:
		return fmt.Errorf("ds.Array: status_t(%v): invalid arguments", s)
	case C.S_ERR_ELEMENT_SIZE_MISMATCH:
		return fmt.Errorf("ds.Array: status_t(%v): type size mismatched", s)
	}
}
