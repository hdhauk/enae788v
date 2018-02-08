package main

type Edge struct {
	tail   int
	head   int
	weight float64
}

type Vertex struct {
	id        int
	x, y      float64
	parent    int
	distance  float64
	finite    bool
	neighbors map[int]float64
	index     int // for use in heap
}

func (v *Vertex) NextNeighborNextNeighbor() string {
	return "not implemented"
}
