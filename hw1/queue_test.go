package main

import (
	"container/heap"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomFloat64(min int, max int) float64 {
	p := rng.Perm(max - min + 1)
	d := rng.Float64()
	return float64(p[min]) + d
}

func TestRandomQueue(t *testing.T) {
	q := NewQueue()

	for i := 0; i < 100; i++ {
		v := &Vertex{
			id:       i + 1,
			distance: RandomFloat64(1, 1000),
			finite:   true,
		}
		heap.Push(q, v)
	}

	largest := 0.0
	for len(q.vertices) > 0 {
		v := heap.Pop(q).(*Vertex)
		if v.distance < largest {
			t.Errorf("v.distance smaller. Got %.4f, largest seen %.4f", v.distance, largest)
		}
		largest = v.distance
	}
}

func TestChangingQueue(t *testing.T) {
	q := NewQueue()

	var u *Vertex
	for i := 0; i < 100; i++ {
		v := &Vertex{
			id:       i + 1,
			distance: RandomFloat64(1, 10),
			finite:   true,
		}
		if i == 50 {
			u = v
			fmt.Println(v)
		}
		heap.Push(q, v)
	}
	magicNumber := 69.69
	fmt.Println(u)
	u.distance = magicNumber
	heap.Fix(q, u.index)

	largest := 0.0
	for len(q.vertices) > 0 {
		v := heap.Pop(q).(*Vertex)
		if v.distance < largest {
			t.Errorf("v.distance smaller. Got %.4f, largest seen %.4f", v.distance, largest)
		}
		largest = v.distance
	}
	if largest != magicNumber {
		t.Errorf("Expected largest distance to be &%.2f, but found %.4f", magicNumber, largest)
	}
}
