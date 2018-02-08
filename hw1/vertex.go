package main

type Edge struct {
	tail   int
	head   int
	weight float64
}

type Vertex struct {
	id          int
	x, y        float64 // cartesian position
	parent      *Vertex
	costToStart float64 // current costToStart (dijkstra)
	priority    float64 // based on costToStart and heuristic
	finite      bool
	neighbors   map[int]float64
	index       int // for use in heap
}

func (v *Vertex) NextNeighborNextNeighbor() string {
	return "not implemented"
}
