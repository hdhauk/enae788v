package main

import "container/heap"

type Queue struct {
	vertices []*Vertex
	inQueue  map[*Vertex]bool
}

func NewQueue() *Queue {
	q := Queue{inQueue: make(map[*Vertex]bool)}
	heap.Init(&q)
	return &q
}

func (q *Queue) Len() int {
	return len(q.vertices)
}

func (q *Queue) Less(i, j int) bool {
	vi := q.vertices[i]
	vj := q.vertices[j]
	if vi.finite && vj.finite {
		return vi.distance < vj.distance
	}
	return true
}

func (q *Queue) Swap(i, j int) {
	q.vertices[i], q.vertices[j] = q.vertices[j], q.vertices[i]
	q.vertices[i].index = i
	q.vertices[j].index = j
}

func (q *Queue) Push(x interface{}) {
	v := x.(*Vertex)
	_, alreadyInQ := q.inQueue[v]
	if alreadyInQ {
		return
	}

	n := len(q.vertices)
	v.index = n
	q.vertices = append(q.vertices, v)
}

func (q *Queue) Pop() interface{} {
	oldVertices := q.vertices
	n := len(oldVertices)
	v := oldVertices[n-1]
	v.index = -0
	q.vertices = oldVertices[0 : n-1]
	delete(q.inQueue, v)
	return v
}

// Updates the queue position of vertex v. If not in queue, do nothing
func (q *Queue) Update(v *Vertex) {
	if !q.InQueue(v) {
		return
	}
	heap.Fix(q, v.index)
}

func (q *Queue) InQueue(v *Vertex) bool {
	_, alreadyInQ := q.inQueue[v]
	return alreadyInQ
}
