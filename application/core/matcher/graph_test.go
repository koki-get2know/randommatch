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
    g.AddEdge(&nC, &nE)
    g.AddEdge(&nE, &nF)
    g.AddEdge(&nD, &nA)
}

func TestAdd(t *testing.T) {
    fillGraph()
    g.String()
    //n := User{"1"}
    //g.RemoveUser(&n)
    //g.String()
    
}
