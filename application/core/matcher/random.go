package matcher

import (
	"math/rand"
	"strconv"
	"time"
)

// TODO the selector paramater should be a config variable
const SELECTOR = "random"

type Matching struct {
	matched []*User
}

type Match struct {
	Id    string
	Users []User
}

func Filter(g *UserGraph, matched []*User, n *User, constraints []string) bool {
	/* Filter

	   input :
	          g:a graph of users,
	          matched: users who already match
	          n: new users
	          constraints: Constraints that must be respected by the matching
	   output : bool
	   purpose : check if user in matched and user n can match

	*/

	ok := true
	for _, constraint := range constraints {
		switch constraint {
		case "deja vu":
			// check if an edges exist between n and any user in matched
			for _, user := range matched {
				if find, _ := Search(g.edges[*n], user); find {
					ok = false
					break
				}
			}
		}
	}

	return ok
}

func RandomChoices(g *UserGraph, k int, constraints []string) *Matching {

	/* random choice without constraint

	   input : a graph of users, k length of tuple match
	   output : k random id
	   purpose : match k user from the graph

	*/

	var matching = &Matching{}
	var matchedUsers []*User
	rand.Seed(time.Now().UnixNano()) // initialize the seed to get
	for i := 0; i < k; i++ {
		index := rand.Intn(len(g.users))
		find, _ := Search(matchedUsers, g.users[index])            // check if g.users[index] already exist in matchedUsers
		ok := Filter(g, matchedUsers, g.users[index], constraints) // check the constraints
		if !find && ok {

			matchedUsers = append(matchedUsers, g.users[index])

		} else {
			i--
		}
	}
	matching.matched = matchedUsers
	return matching
}

//TODO
/*func RandSubGroup(A *UserGraph, B *Usergraph, constraints []string) *Matching{


}*/

func Matcher(g *UserGraph, k int, constraints []string) map[int]*Matching {

	/* Matcher without constraint

	   input : g User's graph, k length of tuple match
	   output : list of tuple match
	   purpose: match all user in graph the g

	*/
	matching := make(map[int]*Matching)

	switch SELECTOR {
	case "random":
		/*
		   repeat
		     1 - random choices k users
		     2 - remove previous users to the graph
		   until is possible to take k users in graph

		*/
		i := 0
		for len(g.users)/k > 0 {
			matched := RandomChoices(g, k, constraints)
			for _, match := range matched.matched {
				g.RemoveUser(match)
			}
			matching[i] = matched
			i++
		}

	}

	return matching

}

func GenerateTuple(users []User, size int) []Match {
	var results []Match
	graph := UsersToGraph(users)
	tuples := Matcher(graph, size, []string{"deja vu"})
	for index, matching := range tuples {
		var matches []User
		for _, user := range matching.matched {
			matches = append(matches, *user)
		}
		results = append(results, Match{Id: strconv.Itoa(index), Users: matches})
	}
	return results
}
