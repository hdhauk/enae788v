package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/ungerik/go3d/float64/vec2"
)

// SafeFunc takes to points and return true if the the edge is safe.
type SafeFunc func(v, w *Vertex) bool

// Edge define an edge between two vertices.
type Edge struct {
	head, tail *Vertex
}

// Vertex define a node in the genrated graph.
type Vertex struct {
	Point
	Parent *Vertex
}

func (v Vertex) String() string {
	return fmt.Sprintf("(%.2f, %.2f)", v.X, v.Y)
}

// LinkParent links the parent pointer to the parent vertex.
func (v *Vertex) LinkParent(parent *Vertex) {
	v.Parent = parent
}

// RRT build a tree and find a feasible path using the RRT algorithm.
func RRT(obstacles []Circle, prob Problem, cSpace ConfigSpace, safe SafeFunc, seed int64) (path, tree []Edge, err error) {
	rand.Seed(seed)
	vertices := []*Vertex{&Vertex{Point: prob.Start, Parent: nil}}
	edges := []Edge{}

	var u, v, w *Vertex
	for {
		u = randomSample(cSpace)
		v = closestMember(vertices, u)
		w = smallDistanceAlong(v, u, prob.Epsilon, prob.AllowSmallSteps)
		w.LinkParent(v)

		// Discard vertex w if edge vw is unsafe.
		if !safe(v, w) {
			continue
		}

		vertices = append(vertices, w)
		edges = append(edges, newEdge(v, w))

		if near(w, prob.Goal) {
			break
		}
	}

	path = backtrack(w, &prob.Start, edges)

	return path, edges, nil
}

// randomSample picks a random point within the configuration space.
func randomSample(c ConfigSpace) *Vertex {
	x := c.XMin + rand.Float64()*(c.XMax-c.XMin)
	y := c.YMin + rand.Float64()*(c.YMax-c.YMin)
	theta := -math.Pi + rand.Float64()*(2*math.Pi)
	return newVertex(x, y, theta, nil)
}

// distance returns the cartesian distance between two vertices.
func distance(a, b *Vertex) float64 {
	return math.Sqrt(math.Pow(b.X-a.X, 2) + math.Pow(b.Y-a.Y, 2))
}

// closestMember naively searches for the member in vertices closest to vertex u.
// Runtime: O(n) where n are number of vertices in list.
func closestMember(vertices []*Vertex, u *Vertex) *Vertex {
	var closest *Vertex
	shortest := math.MaxFloat64
	for _, v := range vertices {
		d := distance(u, v)
		if d < shortest {
			closest = v
			shortest = d
		}
	}
	return closest
}

// smallDistanceAlong generate an edge between u and v, and returns a new vertex w
// epsilon distance from u along the uv-edge.
func smallDistanceAlong(u, v *Vertex, epsilon float64, smallSteps bool) *Vertex {
	u2v := vec2.T{v.X - u.X, v.Y - u.Y}

	if smallSteps && u2v.Length() < epsilon {
		return v
	}

	u2vNorm := u2v.Normalize()
	wVec := u2vNorm.Scale(epsilon)
	x := u.X + wVec[0]
	y := u.Y + wVec[1]
	// theta := wVec.Angle()
	theta := v.Theta

	return newVertex(x, y, theta, u)
}

// newEdge return a new edge from two vertices.
func newEdge(tail, head *Vertex) Edge {
	return Edge{
		tail: tail,
		head: head,
	}
}

// newVertex return a new parentless vertex.
func newVertex(x, y, theta float64, parent *Vertex) *Vertex {
	return &Vertex{Point{X: x, Y: y, Theta: theta}, parent}
}

// near returns true if vertex u is within circle goal.
func near(u *Vertex, goal Circle) bool {
	d := distance(u, newVertex(goal.X, goal.Y, 0, nil))
	return d < goal.R
}

// backtrack generate a slice of edges from a leaf vertex to the root of the tree.
func backtrack(w *Vertex, start *Point, edges []Edge) []Edge {

	current := w
	path := []Edge{}
	for current.Parent != nil {
		path = append(path, newEdge(current, current.Parent))
		current = current.Parent
	}

	return path
}
