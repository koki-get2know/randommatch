package matcher

import (
	"math/rand"
	"strconv"
	"time"
)

// TODO the selector paramater should be a config variable
type Selector uint8

const (
	Basic Selector = iota
)

type Constraint uint8

const (
	Dejavu Constraint = iota
	ForbiddenConnections
)

type Match struct {
	Id    string `json:"id"`
	Users []User `json:"users"`
}

type ConstraintParams[T any] struct { // Type parameters for constraint
	Params map[Constraint][]T
}

func search(users []User, n User) (bool, int) {
	index := -1
	find := false
	for i, user := range users {
		if user.UserId == n.UserId {
			find = true
			index = i
			break
		}
	}

	return find, index
}
func Filter(g *UserGraph, matched []User, n *User, constraints []Constraint, forbiddenConnections [][]User) bool {
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
		if ok {
			for _, user := range matched {
				if ok {
					switch constraint {
					case Dejavu:
						// check if an edges exist between n and any user in matched

						if find, _ := Search(g.edges[*n], &user); find {
							ok = false

						}

					case ForbiddenConnections:
						for _, usersNotToMatch := range forbiddenConnections {
							if len(usersNotToMatch) > 0 {
								find1, _ := search(usersNotToMatch, *n)
								find2, _ := search(usersNotToMatch, user)
								if find1 && find2 {
									ok = false
									break
								}
							}

						}
					}
				}
			}
		} else {
			break
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

func RandomChoices(g *UserGraph, k uint, constraints []Constraint, forbiddenConnections [][]User) *Match {

	/* random choice without constraint

	   input : a graph of users, k length of tuple match
	   output : k random id
	   purpose : match k user from the graph

	*/

	var matching = &Match{}
	var matchedUsers []User
	rand.Seed(time.Now().UnixNano()) // initialize the seed to get

	var indices []int
	for i := range g.users {
		indices = append(indices, i)
	}

	for uint(len(matchedUsers)) < k && len(indices) > 0 {
		rand.Shuffle(len(indices), func(i, j int) { indices[i], indices[j] = indices[j], indices[i] })
		index := indices[0]

		ok := Filter(g, matchedUsers, g.users[index], constraints, forbiddenConnections) // check the constraints
		if ok {

			matchedUsers = append(matchedUsers, *g.users[index])
		}
		indices = remove(indices, index)
	}
	matching.Users = matchedUsers
	matching.Id = ""
	return matching
}

//TODO
func RandSubGroup(g *UserGraph, A []*User, B []*User, k uint, constraints []Constraint, forbiddenConnections [][]User) *Match {

	/*Extract the subgraph A and B
	    m1 = random choice k-k%2 users dans A
	    repeat until good match

		   - m2 = random choice k/2 users dans B
		   - check m1 + m2
	*/
	groupeA := g.Subgrapn(A)
	groupeB := g.Subgrapn(B)
	match := false
	matchA := RandomChoices(groupeA, k/2, constraints, forbiddenConnections)
	for !match {
		matchB := RandomChoices(groupeB, k-k/2, constraints, forbiddenConnections)
		ok := true
		for _, u := range matchB.Users {

			if Filter(g, matchA.Users, &u, constraints, forbiddenConnections) {
				matchA.Users = append(matchA.Users, u)
			} else {
				ok = false
				break
			}
		}
		match = ok
	}

	return matchA

}

func Matcher(g *UserGraph, k uint, constraints []Constraint, SELECTOR Selector, forbidenconections [][]User) map[int]*Match {

	/* Matcher without constraint

	   input : g User's graph, k length of tuple match
	   output : list of tuple match
	   purpose: match all user in graph the g

	*/
	matching := make(map[int]*Match)

	switch SELECTOR {
	case Basic:
		/*
		   repeat
		     1 - random choices k users
		     2 - remove previous users to the graph
		   until is possible to take k users in graph

		*/
		i := 0
		for k > 0 && uint(len(g.users))/k > 0 {
			matched := RandomChoices(g, k, constraints, forbidenconections)
			for _, match := range matched.Users {
				g.RemoveUser(&match)
			}
			matching[i] = matched
			i++
		}
	}

	return matching

}

func GenerateTuple(users []User, connections [][]User, forbiddenConnections [][]User, size uint) []Match {
	/*
		         Input :
				      users: Users for matching
					  connections: Connections plan in the graph; first element of the tab
	*/
	var results []Match
	graph := UsersToGraph(users, connections)
	tuples := Matcher(graph, size, []Constraint{Dejavu}, Basic, forbiddenConnections)
	for index, matching := range tuples {

		results = append(results, Match{Id: strconv.Itoa(index), Users: matching.Users})
	}
	return results
}
