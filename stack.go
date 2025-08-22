package ds

import "github.com/josestg/ds/adt"

type stackable[E any] interface {
	adt.Sizer
	adt.Emptier
	adt.Tailer[E]
	adt.Popper[E]
	adt.Appender[E]
	adt.Stringer
}

type Stack[E any] struct {
	b stackable[E]
}

func NewStack[E any]() *Stack[E] {
	return NewStackWith[E](NewDoublyLinkedList[E]())
}

func NewStackWith[E any](b stackable[E]) *Stack[E] {
	return &Stack[E]{b: b}
}

func (s *Stack[E]) Empty() bool {
	s.ensureBackend()
	return s.b.Empty()
}

func (s *Stack[E]) Size() int {
	s.ensureBackend()
	return s.b.Size()
}

func (s *Stack[E]) Peek() E {
	s.ensureBackend()
	if s.Empty() {
		panic("stack.Peek: stack is empty")
	}
	return s.b.Tail()
}

func (s *Stack[E]) Push(data E) {
	s.ensureBackend()
	s.b.Append(data)
}

func (s *Stack[E]) Pop() E {
	s.ensureBackend()
	if s.Empty() {
		panic("stack.Pop: stack underflow")
	}
	return s.b.Pop()
}

func (s *Stack[E]) String() string {
	s.ensureBackend()
	return s.b.String()
}

func (s *Stack[E]) ensureBackend() {
	if s.b == nil {
		s.b = NewSinglyLinkedList[E]()
	}
}
