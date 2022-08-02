package matcher

import (
	"testing"
)

var g UserGraph

func fillGraph() {
	nA := User{"1"}
	nB := User{"2"}
	nC := User{"3"}
	nD := User{"4"}
	nE := User{"5"}
	nF := User{"6"}
	g.AddUser(&nA)
	g.AddUser(&nB)
	g.AddUser(&nC)
	g.AddUser(&nD)
	g.AddUser(&nE)
	g.AddUser(&nF)

	g.AddEdge(&nA, &nB)
	g.AddEdge(&nA, &nC)
	g.AddEdge(&nB, &nE)
	g.AddEdge(&nF, &nE)
	g.AddEdge(&nA, &nE)

}

func TestAdd(t *testing.T) {
	fillGraph()
	g.String()
	n1 := User{"4"}
	n2 := User{"5"}
	A := []*User{&n1, &n2}
	sub := g.Subgraph(A)
	sub.String()

}
