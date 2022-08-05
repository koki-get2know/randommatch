package matcher

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	"github.com/koki/randommatch/entity"
)

// TODO the selector paramater should be a config variable
type Selector uint8

const (
	Basic Selector = iota
	Group
)

type Constraint uint8

const (
	Unique Constraint = iota
	ForbiddenConnections
)

type Match struct {
	Id    string        `json:"id"`
	Users []entity.User `json:"users"`
}

//TODO integrate  the constraint structure in code
type ConstraintParams[T any] struct { // Type parameters for constraint
	Params map[Constraint][]T
}

//TODO integrate the selector structure in code
type SelectorParams[T any] struct { // Type parameters for constraint
	Params map[Constraint][]T
}

func search(users []entity.User, n entity.User) (bool, int) {
	index := -1
	find := false
	for i, user := range users {
		if user.Id == n.Id {
			find = true
			index = i
			break
		}
	}

	return find, index
}

func Filter(g *UserGraph, matched []entity.User, n *entity.User,
	constraints []Constraint, forbiddenConnections [][]entity.User) bool {
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

constraintloop:
	for _, constraint := range constraints {
		for _, user := range matched {
			switch constraint {
			case Unique:
				// check if an edges exist between n and any user in matched

				if find, _ := Search(g.edges[(*n).Id], &user); find {
					ok = false
					break constraintloop
				}

			case ForbiddenConnections:
				for _, usersNotToMatch := range forbiddenConnections {
					if len(usersNotToMatch) > 0 {
						if find1, _ := search(usersNotToMatch, *n); find1 {
							if find2, _ := search(usersNotToMatch, user); find2 {
								ok = false
								break constraintloop
							}
						}
					}

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

func RandomChoices(g *UserGraph, k uint, constraints []Constraint, forbiddenConnections [][]entity.User) *Match {

	/* random choice without constraint

	   input : a graph of users, k length of tuple match
	   output : k random id
	   purpose : match k user from the graph

	*/

	var matching = &Match{}
	var matchedUsers []entity.User
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

func RandSubGroup(groupeA *UserGraph, groupeB *UserGraph, matchSizeA uint, matchSizeB uint,
	interGroupConstraints []Constraint, innerGroupConstraints []Constraint, forbiddenConnections [][]entity.User) *Match {

	/*
		   input :
		       g : users graph
			   groupeA : subgraph of g A
			   groupeB : subgraph of g B
			   matchSizeA : number of person to take in A
			   matchSizeB : number of person to take in B
			   interGroupConstraint : constraints for matching between groupe A and B
			   innerGroupConstraint : constraints for matching into a group
		       forbedenConnections : parameter of constraint forbidenconnection; it contains the user who cannot match together

		   purpose : match users' groupeA with users' groupeB according to the innerGroupConstraint and interGroupConstraint
		   output : match of size matchSizeA + matchSizeB
	*/
	matchA := &Match{}
	match := false

	if uint(len(groupeA.users)) >= matchSizeA && uint(len(groupeB.users)) >= matchSizeB {

		matchA = RandomChoices(groupeA, matchSizeA, innerGroupConstraints, forbiddenConnections)

		users := []entity.User{}
		gb := &UserGraph{}
		copier.Copy(&users, &matchA.Users)
		copier.Copy(gb, groupeB)
		for !match && uint(len(matchA.Users)) < (matchSizeA+matchSizeB) && uint(len(gb.users)) >= matchSizeB {
			matchB := RandomChoices(gb, matchSizeB, innerGroupConstraints, forbiddenConnections)

			ok := true

			for _, u := range matchB.Users {
				u := u
				if Filter(gb, users, &u, interGroupConstraints, forbiddenConnections) {
					matchA.Users = append(matchA.Users, u)

				} else {
					ok = false
					gb.RemoveUser(&u)
					break
				}
				gb.RemoveUser(&u)
			}
			match = ok
		}
	}
	return matchA

}

func Matcher(g *UserGraph, k uint,
	constraints []Constraint, SELECTOR Selector, forbidenconections [][]entity.User,
	A []*entity.User, B []*entity.User, matchSizeA uint, matchSizeB uint,
	interGroupConstraints []Constraint, innerGroupConstraints []Constraint) map[int]*Match {

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
				match := match
				g.RemoveUser(&match)

			}
			matching[i] = matched
			i++
		}
	case Group:
		/*Extract the subgraph A and B
		    m1 = random choice  users dans A
		    repeat until good match

			   - m2 = random choice  users dans B
			   - check if m1 + m2 can be match
		*/
		groupA := g.Subgraph(A)
		groupB := g.Subgraph(B)
		groupA.String()
		groupB.String()
		i := 0
		if matchSizeB > 0 && matchSizeA > 0 {
			for uint(len(groupA.users))/matchSizeA > 0 && uint(len(groupB.users))/matchSizeB > 0 {
				matched := RandSubGroup(groupA, groupB, matchSizeA, matchSizeB,
					interGroupConstraints, innerGroupConstraints,
					forbidenconections)
				if matched != nil {
					for _, match := range matched.Users {
						match := match
						groupA.RemoveUser(&match)
						groupB.RemoveUser(&match)

					}
				}

				matching[i] = matched
				i++

			}
		}

	}

	return matching

}

func GenerateTuple(users []entity.User, connections [][]entity.User, s Selector,
	forbiddenConnections [][]entity.User, size uint,
	A []entity.User, B []entity.User, sizeA uint, sizeB uint) []Match {
	/*
				         Input :
						    - general
						      users: Users for matching
							  connections: Connections plan in the graph;
							  s   : type of selector : basic, group .....
							- specific
							  *constraint
							   forbiddenConnecitons: variable for forbideenconnections constraints
							  *selector
							    size : size of matching; variable for basic selector\
							    A,B : groupe A and B of  user; variable for group selector
		                        sizeA, sizeB: size of matching for A and B respectivily; variable for group selector
	*/
	var results []Match
	var tuples map[int]*Match

	if len(users) == 0 && s == Group {
		users = append(A, B...)
	}

	graph := UsersToGraph(users, connections)

	switch s {
	case Basic:
		tuples = Matcher(graph, size, []Constraint{Unique}, Basic,
			forbiddenConnections, []*entity.User{}, []*entity.User{},
			0, 0, []Constraint{}, []Constraint{})
	case Group:
		gA := []*entity.User{}
		gB := []*entity.User{}
		for _, u := range A {
			u := u
			gA = append(gA, &u)

		}
		for _, u := range B {
			u := u
			gB = append(gB, &u)
		}

		tuples = Matcher(graph, size, []Constraint{}, Group,
			forbiddenConnections, gA, gB,
			sizeA, sizeB, []Constraint{Unique}, []Constraint{})
	}
	for index, matching := range tuples {

		results = append(results, Match{Id: strconv.Itoa(index), Users: matching.Users})
	}
	return results
}
