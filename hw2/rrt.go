package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/ungerik/go3d/float64/vec2"
)

func main() {
	configPath := flag.String("-c", "problems.json", "config file")
	problem := flag.Int("-p", 0, "which problem in config file to solve (0-indexed)")

	configFile, err := os.Open(*configPath)
	defer configFile.Close()
	if err != nil {
		log.Fatalf("could not open config file: %v", err)
	}

	config, err := parseConfig(configFile)
	if err != nil {
		log.Fatalf("could not parse config file: %v", err)
	}

	obstacleFile, err := os.Open(config.ObstaclesPath)
	if err != nil {
		log.Fatalf("could not open obstacle file: %v", err)
	}
	defer obstacleFile.Close()
	obstacles, err := readObstacles(obstacleFile)
	if err != nil {
		log.Fatalf("could not read obstacles from file: %v", err)
	}

	fmt.Println(obstacles)
	fmt.Println(*problem)
	fmt.Printf("%+v\n", config)

}

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

// RRT build a tree and find a feasible path using the RRT algorithm.
func RRT(obstacles []Circle, prob Problem, cSpace ConfigSpace, safe SafeFunc, seed int64) (path, tree []Edge, err error) {
	rand.Seed(seed)
	vertices := []*Vertex{&Vertex{Point: prob.Start, Parent: nil}}
	edges := []Edge{}

	var u, v, w *Vertex
	for {
		u = randomSample(cSpace)
		v = closestMember(vertices, u)
		w = smallDistanceAlong(u, v, prob.Epsilon)

		if safe(u, w) {
			vertices = append(vertices, w)
			edges = append(edges, newEdge(v, w))
		}

		if near(w, prob.Goal) {
			break
		}
	}

	path = backtrack(edges[len(edges)-1], &prob.Start, edges)

	return path, edges, nil
}

func randomSample(c ConfigSpace) *Vertex {
	x := c.XMin + rand.Float64()*(c.XMax-c.XMin)
	y := c.YMin + rand.Float64()*(c.YMin-c.YMin)
	return newVertex(x, y, nil)
}

func distance(a, b *Vertex) float64 {
	return math.Sqrt(math.Pow(b.X-a.X, 2) + math.Pow(b.Y-a.Y, 2))
}

func closestMember(vertices []*Vertex, u *Vertex) *Vertex {
	bestVertex := vertices[0]
	bestDistance := distance(u, vertices[0])
	for _, v := range vertices {
		if distance(u, v) < bestDistance {
			bestVertex = v
			bestDistance = distance(u, v)
		}
	}
	return bestVertex
}

func smallDistanceAlong(u, v *Vertex, epsilon float64) *Vertex {
	u2v := vec2.T{v.X - u.X, v.Y - u.Y}
	u2vNorm := u2v.Normalize()
	wVec := u2vNorm.Scale(epsilon)
	return newVertex(wVec[0], wVec[1], u)
}

func newEdge(u, v *Vertex) Edge {
	return Edge{
		tail: u,
		head: v,
	}
}

func newVertex(x, y float64, parent *Vertex) *Vertex {
	return &Vertex{Point{x, y}, parent}
}

func near(u *Vertex, goal Circle) bool {
	d := distance(u, newVertex(goal.X, goal.Y, nil))
	return d < goal.R
}

func backtrack(w Edge, start *Point, edges []Edge) []Edge {
	// path := []Edge{}

	// current := Edge{}
	// for current.parent != start{
	// 	path = append(path, )
	// }
	return nil
}
