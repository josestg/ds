package ds

import (
	"fmt"
	"strings"

	"github.com/josestg/ds/seq"
)

type none = struct{}

type HashSetOptions[E comparable] = HashMapOptions[E]

type HashSet[E comparable] struct {
	backend *HashMap[E, none]
}

func NewHashSet[E comparable]() *HashSet[E] {
	return &HashSet[E]{
		backend: NewHashMap[E, none](),
	}
}

func NewHashSetWith[E comparable](opts HashSetOptions[E]) *HashSet[E] {
	return &HashSet[E]{
		backend: NewHashMapWith[E, none](opts),
	}
}

func (s *HashSet[E]) Add(data E) {
	s.backend.Put(data, none{})
}

func (s *HashSet[E]) Del(data E) {
	s.backend.Del(data)
}
func (s *HashSet[E]) Exists(data E) bool {
	return s.backend.Exists(data)
}

func (s *HashSet[E]) Size() int {
	return s.backend.Size()
}

func (s *HashSet[E]) Empty() bool {
	return s.backend.Empty()
}

func (s *HashSet[E]) String() string {
	var buf strings.Builder
	buf.WriteRune('{')
	for i, v := range seq.Enum(s.Iter) {
		if i > 0 {
			buf.WriteRune(' ')
		}
		_, _ = fmt.Fprint(&buf, v)
	}
	buf.WriteRune('}')
	return buf.String()
}

func (s *HashSet[E]) Iter(yield func(E) bool) {
	for e := range s.backend.Iter {
		if !yield(e.Key()) {
			break
		}
	}
}

func (s *HashSet[E]) Union(s2 *HashSet[E]) *HashSet[E] {
	union := NewHashSet[E]()
	for v := range s.Iter {
		union.Add(v)
	}
	for v := range s2.Iter {
		union.Add(v)
	}
	return union
}

func (s *HashSet[E]) Intersection(s2 *HashSet[E]) *HashSet[E] {
	intersection := NewHashSet[E]()
	var left, right *HashSet[E]
	if s.Size() < s2.Size() {
		left, right = s, s2
	} else {
		left, right = s2, s
	}
	for v := range left.Iter {
		if right.Exists(v) {
			intersection.Add(v)
		}
	}
	return intersection
}
