package main

import (
	"fmt"
	"math"
	"math/rand"
)

var timestep = 0.1

// SafeFunc takes to points and return true if the the edge is safe.
type SafeFunc func(*PathPoint) bool

// Edge define an edge between two vertices.
type Edge struct {
	head, tail *Vertex
	path       []*PathPoint
}

func (e *Edge) String() string {
	return fmt.Sprintf("edge(head:0x%x, tail:0x%x, len(path)=%d)", &e.head, &e.tail, len(e.path))
}

type PathPoint struct {
	x, y, θ, v, w, a, γ float64
}

// Point is a point in 2D space with an optional direction angle associated with it.
type Point struct {
	X, Y, Theta float64
	V           float64 // linear velocity
	W           float64 // angular velocity
}

// Vertex define a node in the genrated graph.
type Vertex struct {
	Point
	Parent      *Vertex
	Edge2Parent *Edge
}

func (v *Vertex) String() string {
	hasPath := "False"
	if v.Edge2Parent != nil && len(v.Edge2Parent.path) > 1 {
		hasPath = "True"
	}
	return fmt.Sprintf("(x:%.2f, y:%.2f, θ:%.2f, v:%.2f, w:%.2f, parent:0x%x, path2parent: %s)",
		v.X, v.Y, v.Theta, v.V, v.W, &v.Parent, hasPath)
}

// LinkParent links the parent pointer to the parent vertex.
func (v *Vertex) LinkParent(parent *Vertex) {
	v.Parent = parent
}

// RRT build a tree and find a feasible path using the RRT algorithm.
func RRT(obstacles []Circle, prob Problem, cSpace *ConfigSpace, safe SafeFunc, seed int64) (path []*PathPoint, tree []*Edge, err error) {
	rand.Seed(seed)
	vertices := []*Vertex{&Vertex{Point: prob.Start, Parent: nil}}
	edges := []*Edge{}

	var u, v, w *Vertex
	var ok bool
	var edge *Edge
	for {
		u = randomSample(cSpace)
		v = closestMember(vertices, u)
		w, edge, ok = forwardSim(v, u, prob.Epsilon, prob.Delta, safe, cSpace)
		if !ok {
			// fmt.Println("found unsafe path...")
			continue
		}
		w.Edge2Parent = edge

		vertices = append(vertices, w)
		edges = append(edges, edge)

		if near(w, prob.Goal) {
			break
		}
	}
	vertexInGoal := w

	path = backtrack(vertexInGoal, &prob.Start, edges)

	return path, edges, nil
}

// randomSample picks a random point within the configuration space.
func randomSample(c *ConfigSpace) *Vertex {
	x := c.XMin + rand.Float64()*(c.XMax-c.XMin)
	y := c.YMin + rand.Float64()*(c.YMax-c.YMin)
	theta := -math.Pi + rand.Float64()*(2*math.Pi)

	v := c.VMin + rand.Float64()*(c.VMax-c.VMin)
	w := c.WMin + rand.Float64()*(c.WMax-c.WMin)

	return newVertex(x, y, theta, v, w, nil)
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

// newEdge return a new edge from two vertices.
func newEdge(tail, head *Vertex) Edge {
	return Edge{
		tail: tail,
		head: head,
	}
}

// newVertex return a new parentless vertex.
func newVertex(x, y, theta, v, w float64, parent *Vertex) *Vertex {
	return &Vertex{Point{X: x, Y: y, Theta: theta, V: v, W: w}, parent, nil}
}

// near returns true if vertex u is within circle goal.
func near(u *Vertex, goal Circle) bool {
	d := distance(u, newVertex(goal.X, goal.Y, 0, 0, 0, nil))
	return d < goal.R
}

// backtrack generate a slice of edges from a leaf vertex to the root of the tree.
func backtrack(inGoal *Vertex, start *Point, edges []*Edge) []*PathPoint {
	current := inGoal
	path := []*PathPoint{}
	i := 0
	for current.Parent != nil {
		path = append(current.Edge2Parent.path, path...)
		current = current.Parent
		i++
	}
	return path
}

// forwardSim forward simulate a trajectory toward u using ½-car like model
func forwardSim(v, u *Vertex, epsilon, delta float64, safe SafeFunc, cspace *ConfigSpace) (*Vertex, *Edge, bool) {

	changeInLinVelocity := u.V - v.V
	changeInAngVelocity := u.W - v.W
	avgSpeed := v.V + changeInLinVelocity/2
	timeToTravelEpsilon := clamp(epsilon/avgSpeed, 1, 10)

	// determine acceleration
	a := clamp(changeInLinVelocity/timeToTravelEpsilon, cspace.AMin, cspace.AMax)
	gamma := clamp(changeInAngVelocity/timeToTravelEpsilon, cspace.GammaMin, cspace.GammaMax)

	h := timestep // global variable. also used for printing...
	X := Point{v.X, v.Y, v.Theta, v.V, v.W}
	path := []*PathPoint{}
	var i float64
	for i = 0.0; i*h < timeToTravelEpsilon; i++ {
		next := euler(X, h, a, gamma)

		if !safe(&next) {
			return nil, nil, false
		}
		path = append(path, &next)
		X = Point{next.x, next.y, next.θ, next.v, next.w}
	}

	head := Vertex{Point: X, Parent: v}
	edge := Edge{head: &head, tail: v, path: path}
	return &head, &edge, true
}

type stateSpace struct {
	x, y, theta, v, w float64
}

func euler(X Point, h, a, gamma float64) PathPoint {
	x := X.X + h*(X.V*math.Cos(X.Theta))
	y := X.Y + h*(X.V*math.Sin(X.Theta))
	theta := X.Theta + h*X.W
	v := X.V + h*a
	w := X.W + h*gamma
	return PathPoint{x, y, theta, v, w, a, gamma}

}

func clamp(a, min, max float64) float64 {
	if a > max {
		return max
	} else if a < min {
		return min
	}
	return a
}
