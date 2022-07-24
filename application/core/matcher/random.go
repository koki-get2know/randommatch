package matcher

import (
	"math/rand"
	"strconv"
	"time"
)

// TODO the selector paramater should be a config variable
const SELECTOR = "random"

type Constraint uint8

const (
	dejavu Constraint = iota
)

type Matching struct {
	matched []*User
}

type Match struct {
	Id    string `json:"id"`
	Users []User `json:"users"`
}

func Filter(g *UserGraph, matched []*User, n *User, constraints []Constraint) bool {
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
		case dejavu:
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

func remove[T comparable](l []T, item T) []T {
	for i, elem := range l {
		if elem == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func RandomChoices(g *UserGraph, k uint, constraints []Constraint) *Matching {

	/* random choice without constraint

	   input : a graph of users, k length of tuple match
	   output : k random id
	   purpose : match k user from the graph

	*/

	var matching = &Matching{}
	var matchedUsers []*User
	rand.Seed(time.Now().UnixNano()) // initialize the seed to get

	var indices []int
	for i := range g.users {
		indices = append(indices, i)
	}

	for uint(len(matchedUsers)) < k && len(indices) > 0 {
		rand.Shuffle(len(indices), func(i, j int) { indices[i], indices[j] = indices[j], indices[i] })
		index := indices[0]
		// check the constraints
		if Filter(g, matchedUsers, g.users[index], constraints) {
			matchedUsers = append(matchedUsers, g.users[index])
		}
		indices = remove(indices, index)
	}
	matching.matched = matchedUsers
	return matching
}

//TODO
/*func RandSubGroup(A *UserGraph, B *Usergraph, constraints []string) *Matching{


}*/

func Matcher(g *UserGraph, k uint, constraints []Constraint) map[int]*Matching {

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
		for k > 0 && uint(len(g.users))/k > 0 {
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

func GenerateTuple(users []User, forbiddenConnections [][]User, size uint) []Match {
	var results []Match
	graph := UsersToGraph(users, forbiddenConnections)
	tuples := Matcher(graph, size, []Constraint{dejavu})
	for index, matching := range tuples {
		var matches []User
		for _, user := range matching.matched {
			matches = append(matches, *user)
		}
		results = append(results, Match{Id: strconv.Itoa(index), Users: matches})
	}
	return results
}
