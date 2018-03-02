package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
type SafeFunc func(v, w Point) bool

// Edge define an edge between two vertices.
type Edge struct {
	head, tail Point
}

// RRT build a tree and find a feasible path using the RRT algorithm.
func RRT(obstacles []Circle, start Point, goal Circle, epsilon float64, safe SafeFunc) (path, tree []Edge, err error) {

	vertices := []Point{start}
	edges := []Edge

	var u,v,w Point
	for {
		u = randomSample()
		v = closestMember(u)
		w = smallDistanceAlong(u,v)

		if safe(u,w){
			vertices = append(vertices, w)
			edges = append(edges, NewEdge(v,w))
		}

		if near(w,goal, epsilon){
			break
		}
	}

	path = backtrack(w, start, edges)



	return nil, nil, nil
}
