package matcher

import (
	"fmt"
	"testing"
)

func TestRandomChoices(t *testing.T) {
	//id := []string{"2", "5", "6", "8", "10", "12", "24", "25"}
	g.String()
	var forbiddenConnections [][]User
	constraint := []Constraint{Dejavu}
	matching := RandomChoices(&g, 2, constraint, forbiddenConnections)
	fmt.Printf("Match of %d: [", len(matching.Users))
	for _, user := range matching.Users {
		fmt.Printf("%s,", user.String())
	}

	fmt.Printf("]")
}

func TestRanSubGroup(t *testing.T) {
	g.String()
	var forbiddenConnections [][]User
	constraint := []Constraint{Dejavu}
	A := []*User{&User{"4"}, &User{"5"}, &User{"6"}}
	subA := g.Subgrapn(A)
	fmt.Println("Sous groupe A")
	subA.String()
	B := []*User{&User{"1"}, &User{"2"}, &User{"3"}}
	fmt.Println("Sous groupe B")
	subB := g.Subgrapn(B)
	subB.String()
	matching := RandSubGroup(&g, A, B, 3, constraint, forbiddenConnections)
	fmt.Printf("Match of %d: [", len(matching.Users))
	for _, user := range matching.Users {
		fmt.Printf("%s,", user.String())
	}

	fmt.Printf("]")
}
func TestMatcher(t *testing.T) {

	g.String()
	var forbiddenConnections [][]User
	A := []User{User{"2"}, User{"1"}}
	forbiddenConnections = append(forbiddenConnections, A)
	constraint := []Constraint{Dejavu}
	SELECTOR := Basic
	matching := Matcher(&g, 2, constraint, SELECTOR, forbiddenConnections)

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

func TestGenTuple(t *testing.T) {

	users := []User{User{"1"}, User{"2"}, User{"3"}, User{"4"}, User{"5"}, User{"6"}}
	var connections, forbiddenConnections [][]User

	matching := GenerateTuple(users, connections, forbiddenConnections, 2)

	for _, match := range matching {
		fmt.Printf("Match : [")
		for _, user := range match.Users {
			fmt.Printf("%s,", user.String())
		}

		fmt.Printf("]")
		fmt.Println("")
	}
}
