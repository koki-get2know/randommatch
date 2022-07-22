package matcher

import (
	"fmt"
	"testing"
)

func TestRandomChoices(t *testing.T) {
	//id := []string{"2", "5", "6", "8", "10", "12", "24", "25"}
	g.String()
	constraint := []Constraint{dejavu}
	matching := RandomChoices(&g, 2, constraint)
	fmt.Printf("Match of %d: [", len(matching.matched))
	for _, user := range matching.matched {
		fmt.Printf("%s,", user.String())
	}

	fmt.Printf("]")
}

func TestMatcher(t *testing.T) {
	//id := []string{"2", "5", "6", "8", "10", "12", "24", "25"}
	g.String()
	constraint := []Constraint{dejavu}
	matching := Matcher(&g, 2, constraint)

	for _, match := range matching {
		fmt.Printf("Match : [")
		for _, user := range match.matched {
			fmt.Printf("%s,", user.String())
		}

		fmt.Printf("]")
		fmt.Println("")

	}
	g.String()
}
