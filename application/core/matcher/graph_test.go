package matcher

import (
	"testing"

	"github.com/koki/randommatch/entity"
)

var g UserGraph

func fillGraph() {
	nA := entity.User{Id: "1"}
	nB := entity.User{Id: "2"}
	nC := entity.User{Id: "3"}
	nD := entity.User{Id: "4"}
	nE := entity.User{Id: "5"}
	nF := entity.User{Id: "6"}
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
	n1 := entity.User{Id: "4"}
	n2 := entity.User{Id: "5"}
	A := []*entity.User{&n1, &n2}
	sub := g.Subgraph(A)
	sub.String()

}
