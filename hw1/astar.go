package main

import (
	"math"

	"github.com/pkg/errors"
)

type heuristic func(u, goal *Vertex) float64

type searchResult struct {
	path       []*Vertex
	searchTree []*Vertex
	pathCost   float64
}

func aStar(vertices map[int]*Vertex, start, goal int, h heuristic) (*searchResult, error) {
	unvisited := make(map[int]bool)
	for k := range vertices {
		unvisited[k] = true
	}
	delete(unvisited, start)

	Q := NewQueue()
	Q.PushVertex(vertices[start])

	goalVertex := vertices[goal]

	searchTree := []*Vertex{}
	var success bool
mainLoop:
	for Q.Peek() != nil {
		v := Q.PopVertex()
		searchTree = append(searchTree, v)

		for neighID, d := range v.neighbors {
			u := vertices[neighID]
			_, notSeen := unvisited[u.id]
			newShorterDistance := u.costToStart > v.costToStart+d
			if notSeen || newShorterDistance {
				delete(unvisited, neighID)
				u.finite = true
				u.parent = v
				u.costToStart = v.costToStart + d
				u.priority = u.costToStart + h(u, goalVertex)
				if Q.InQueue(u) {
					Q.UpdateVertex(u)
					// fmt.Printf("updating %v in queue \n", u.id)
				} else {
					Q.PushVertex(u)
					// fmt.Printf("pushing %v onto queue \n", u.id)
				}
			}
			if v.id == goal {
				goalVertex = v
				success = true
				break mainLoop
			}
		}

	}
	if !success {
		return nil, errors.New("could not find goal")
	}

	startToFinish, pathCost := reconstructPath(vertices[start], goalVertex)

	results := &searchResult{
		path:       startToFinish,
		pathCost:   pathCost,
		searchTree: searchTree,
	}

	return results, nil
}

func reconstructPath(start, goal *Vertex) ([]*Vertex, float64) {
	// Walk backword along parent pointers
	path := []*Vertex{goal}
	next := goal.parent
	current := goal
	for next.id != start.id {
		path = append(path, next)
		current = next
		next = current.parent
	}

	// Reverse slice
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path, goal.costToStart
}

func cartesianDistance(u, goal *Vertex) float64 {
	dist := math.Sqrt(math.Pow(goal.x-u.x, 2) + math.Pow(goal.y-u.y, 2))
	return dist
}
