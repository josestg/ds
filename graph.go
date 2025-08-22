package ds

import (
	"iter"

	"github.com/josestg/ds/seq"
)

type Graph[V comparable] struct {
	directed  bool
	adjacency *HashMap[V, *SinglyLinkedList[V]]
}

func NewGraph[V comparable](directed bool) *Graph[V] {
	return &Graph[V]{
		directed:  directed,
		adjacency: NewHashMap[V, *SinglyLinkedList[V]](),
	}
}

func (g *Graph[V]) ensureNode(v V) *SinglyLinkedList[V] {
	neighbors, ok := g.adjacency.Get(v)
	if !ok {
		neighbors = NewSinglyLinkedList[V]()
		g.adjacency.Put(v, neighbors)
	}
	return neighbors
}

func (g *Graph[V]) Size() int {
	return g.adjacency.Size()
}

func (g *Graph[V]) Empty() bool {
	return g.adjacency.Empty()
}

func (g *Graph[V]) AddEdge(from, to V) {
	list := g.ensureNode(from)
	for v := range list.Iter {
		if v == to {
			return
		}
	}
	list.Append(to)
	if !g.directed {
		rev := g.ensureNode(to)
		for v := range rev.Iter {
			if v == from {
				return
			}
		}
		rev.Append(from)
	}
}

func (g *Graph[V]) DelEdge(from, to V) {
	if list, ok := g.adjacency.Get(from); ok {
		for i, v := range seq.Enum(list.Iter) {
			if v == to {
				_ = list.Remove(i)
				break
			}
		}
	}
	if !g.directed {
		if list, ok := g.adjacency.Get(to); ok {
			for i, v := range seq.Enum(list.Iter) {
				if v == from {
					_ = list.Remove(i)
					break
				}
			}
		}
	}
}

func (g *Graph[V]) HasEdge(from, to V) bool {
	if list, ok := g.adjacency.Get(from); ok {
		for v := range list.Iter {
			if v == to {
				return true
			}
		}
	}
	return false
}

func (g *Graph[V]) HasVertex(v V) bool {
	return g.adjacency.Exists(v)
}

func (g *Graph[V]) Vertex(yield func(V) bool) {
	for v := range g.adjacency.Iter {
		if !yield(v.Key()) {
			break
		}
	}
}

func (g *Graph[V]) Neighbors(v V) iter.Seq[V] {
	return func(yield func(V) bool) {
		if list, ok := g.adjacency.Get(v); ok {
			for val := range list.Iter {
				if !yield(val) {
					break
				}
			}
		}
	}
}

type WalkAlgorithm int

const (
	BFS WalkAlgorithm = iota
	DFSPreOrder
	DFSPostOrder
)

type GraphWalker[T comparable] struct {
	alg     WalkAlgorithm
	graph   *Graph[T]
	visited *HashSet[T]
}

func NewGraphWalker[T comparable](g *Graph[T], alg WalkAlgorithm) *GraphWalker[T] {
	return &GraphWalker[T]{
		alg:     alg,
		graph:   g,
		visited: NewHashSet[T](),
	}
}

func (w *GraphWalker[T]) Visited(node T) bool {
	return w.visited.Exists(node)
}

func (w *GraphWalker[T]) Explored() bool {
	for n := range w.graph.Vertex {
		if !w.Visited(n) {
			return false
		}
	}
	return true
}

func (w *GraphWalker[T]) WalkAll(visit func(T)) {
	for n := range w.graph.Vertex {
		if !w.Visited(n) {
			w.Walk(n, visit)
		}
	}
}

func (w *GraphWalker[T]) Walk(start T, visit func(T)) {
	switch w.alg {
	case BFS:
		w.bfs(start, visit)
	case DFSPreOrder, DFSPostOrder:
		w.dfs(start, visit)
	default:
		panic("unknown walk algorithm")
	}
}

func (w *GraphWalker[T]) dfs(start T, visit func(T)) {
	var traverse func(T)
	traverse = func(n T) {
		if w.Visited(n) {
			return
		}
		if w.alg == DFSPreOrder {
			visit(n)
		}
		w.visited.Add(n)
		for adj := range w.graph.Neighbors(n) {
			traverse(adj)
		}
		if w.alg == DFSPostOrder {
			visit(n)
		}
	}
	traverse(start)
}

func (w *GraphWalker[T]) bfs(start T, visit func(T)) {
	q := NewQueue[T]()
	q.Enqueue(start)
	w.visited.Add(start)
	for !q.Empty() {
		node := q.Dequeue()
		visit(node)
		for neighbor := range w.graph.Neighbors(node) {
			if !w.Visited(neighbor) {
				w.visited.Add(neighbor)
				q.Enqueue(neighbor)
			}
		}
	}
}

func (w *GraphWalker[T]) HasCycle() bool {
	stack := NewHashSet[T]()
	visited := NewHashSet[T]()

	var visit func(T) bool
	visit = func(n T) bool {
		if stack.Exists(n) {
			return true
		}
		if visited.Exists(n) {
			return false
		}

		visited.Add(n)
		stack.Add(n)

		for neighbor := range w.graph.Neighbors(n) {
			if visit(neighbor) {
				return true
			}
		}

		stack.Del(n)
		return false
	}

	for n := range w.graph.Vertex {
		if !visited.Exists(n) {
			if visit(n) {
				return true
			}
		}
	}

	return false
}
