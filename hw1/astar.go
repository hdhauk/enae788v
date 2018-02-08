package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type heuristic func(u, goal *Vertex) float64

func aStar(vertices map[int]*Vertex, start, goal int, h heuristic) ([]*Vertex, error) {
	unvisited := make(map[int]bool)
	for k := range vertices {
		unvisited[k] = true
	}
	delete(unvisited, start)

	Q := NewQueue()
	Q.PushVertex(vertices[start])

	goalVertex := vertices[goal]
	var success bool
mainLoop:
	for Q.Peek() != nil {
		v := Q.PopVertex()

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
					fmt.Printf("updating %v in queue \n", u.id)
				} else {
					Q.PushVertex(u)
					fmt.Printf("pushing %v onto queue \n", u.id)
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

	startToFinish := reconstructPath(vertices[start], goalVertex)
	return startToFinish, nil
}

func reconstructPath(start, goal *Vertex) []*Vertex {
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
	return path
}
