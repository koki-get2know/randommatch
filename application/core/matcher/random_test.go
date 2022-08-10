package matcher

import (
	"fmt"
	"testing"

	"github.com/koki/randommatch/entity"
)

func TestRandomChoices(t *testing.T) {
	//id := []string{"2", "5", "6", "8", "10", "12", "24", "25"}
	g.String()

	var forbiddenConnections [][]entity.User
	constraint := []Constraint{Unique}
	RandomChoices := randomChoicesSeed()
	matching := RandomChoices(&g, 2, constraint, forbiddenConnections)

	fmt.Printf("Match of %d: [", len(matching.Users))
	for _, user := range matching.Users {
		fmt.Printf("%s,", user.String())
	}

	fmt.Printf("]")
}

func TestRanSubGroup(t *testing.T) {
	g.String()

	var forbiddenConnections [][]entity.User
	interConstraint := []Constraint{Unique}
	A := []*entity.User{{Id: "5"}}
	subA := g.Subgraph(A)
	fmt.Println("Sous groupe A")
	subA.String()
	B := []*entity.User{{Id: "3"}}
	fmt.Println("Sous groupe B")
	subB := g.Subgraph(B)
	subB.String()
	matching := RandSubGroup(subA, subB, 1, 1, interConstraint, []Constraint{Unique}, forbiddenConnections)

	fmt.Printf("Match of %d: [", len(matching.Users))
	for _, user := range matching.Users {
		fmt.Printf("%s,", user.String())
	}

	fmt.Printf("]")
}
func TestMatcher1(t *testing.T) {

	g.String()

	var forbiddenConnections [][]entity.User
	A := []entity.User{{Id: "2"}, {Id: "1"}}
	forbiddenConnections = append(forbiddenConnections, A)
	constraint := []Constraint{Unique}
	SELECTOR := Basic
	matching := Matcher(&g, 2, constraint, SELECTOR, forbiddenConnections,
		[]*entity.User{}, []*entity.User{}, []Constraint{}, []Constraint{})

	for _, match := range matching {
		fmt.Printf("Match : [")
		for _, user := range match.Users {
			fmt.Printf("%s,", user.String())

		}

		fmt.Printf("]")
		fmt.Println("")

	}
	g.String()
}
func TestMatcher2(t *testing.T) {
	var G UserGraph

	nA := entity.User{Id: "1"}
	nB := entity.User{Id: "2"}
	nC := entity.User{Id: "3"}
	nD := entity.User{Id: "4"}
	nE := entity.User{Id: "5"}
	nF := entity.User{Id: "6"}

	G.AddUser(&nA)
	G.AddUser(&nB)
	G.AddUser(&nC)
	G.AddUser(&nD)
	G.AddUser(&nE)
	G.AddUser(&nF)

	G.AddEdge(&nA, &nB)
	G.AddEdge(&nA, &nC)
	G.AddEdge(&nB, &nE)
	G.AddEdge(&nF, &nE)
	G.AddEdge(&nA, &nE)

	G.String()

	var forbiddenConnections [][]entity.User
	C := []entity.User{{Id: "4"}, {Id: "1"}}
	forbiddenConnections = append(forbiddenConnections, C)
	A := []*entity.User{{Id: "5"}, {Id: "6"}}
	B := []*entity.User{{Id: "1"}, {Id: "2"}, {Id: "4"}}
	interConstraint := []Constraint{Unique}
	SELECTOR := Group
	matching := Matcher(&G, 0, []Constraint{}, SELECTOR, forbiddenConnections, A, B, interConstraint, []Constraint{Unique})

	for _, match := range matching {
		fmt.Printf("Match : [")
		for _, user := range match.Users {
			fmt.Printf("%s,", user.String())

		}

		fmt.Printf("]")
		fmt.Println("")

	}

}

func TestGenTuple(t *testing.T) {

	users := []entity.User{}
	var connections, forbiddenConnections [][]entity.User
	A := []entity.User{{Id: "tata"}, {Id: "titi"}, {Id: "tato"}}
	B := []entity.User{{Id: "tata"}, {Id: "titi"}, {Id: "tato"}}
	C := []entity.User{{Id: "tete"}, {Id: "titi"}, {Id: "tonton"}}
	D := []entity.User{{Id: "tata"}, {Id: "titi"}}
	forbiddenConnections = append(forbiddenConnections, C)
	forbiddenConnections = append(forbiddenConnections, D)
	matching := GenerateTuple(users, connections, Group, forbiddenConnections, 2, A, B)

	for _, match := range matching {
		fmt.Printf("Match : [")
		for _, user := range match.Users {
			fmt.Printf("%s,", user.String())
		}

		fmt.Printf("]")
		fmt.Println("")
	}
}
