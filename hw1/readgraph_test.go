package main

import (
	"testing"
)

func TestReadProblems(t *testing.T) {
	problems, err := readProblems("problems/problems.txt")
	if err != nil {
		t.Error(err)
	}

	if len(problems) != 5 {
		t.Errorf("expected 5 problems, got %d", len(problems))
	}
}

func TestReadNodes(t *testing.T) {
	nodes, err := readVertices("./problems/nodes_1.txt")
	if err != nil {
		t.Error(err)
	}

	if len(nodes) != 100 {
		t.Errorf("expected 100 nodes, got %d", len(nodes))
	}

}

func TestReadEdges(t *testing.T) {
	edges, err := readEdges("./problems/edges_1.txt")
	if err != nil {
		t.Error(err)
	}

	if len(edges) != 1000 {
		t.Errorf("expected 1000 edges, got %d", len(edges))
	}
}
