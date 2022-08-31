package matcher

import (
	"fmt"

	"math/rand"
	"strconv"
	"time"

	"github.com/jinzhu/copier"

	"github.com/koki/randommatch/entity"
)

var randomChoices = randomChoicesSeed()

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
						if find1, _ := n.UserIn(usersNotToMatch); find1 {
							if find2, _ := user.UserIn(usersNotToMatch); find2 {
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

func minimum(a uint, b uint) uint {
	if a < b {
		return a
	}
	return b
}

func remove[T comparable](l []T, item T) []T {
	for i, elem := range l {
		if elem == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func randomChoicesSeed() func(g *UserGraph, k uint, constraints []Constraint, forbiddenConnections [][]entity.User) *Match {
	rand.Seed(time.Now().UnixNano()) // initialize the seed to get

	return func(g *UserGraph, k uint, constraints []Constraint, forbiddenConnections [][]entity.User) *Match {

		/* random choice without constraint

		   input : a graph of users, k length of tuple match
		   output : k random id
		   purpose : match k user from the graph


		*/

		var matching = &Match{}
		var matchedUsers []entity.User
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
	if uint(len(groupeA.users)) >= matchSizeA && uint(len(groupeB.users)) >= matchSizeB {

		matchA = randomChoices(groupeA, matchSizeA, innerGroupConstraints, forbiddenConnections)

		users := []entity.User{}
		gb := &UserGraph{}
		copier.Copy(&users, matchA.Users)
		copier.Copy(&gb.edges, groupeB.edges)
		copier.Copy(&gb.users, groupeB.users)
		match := false

		for !match && uint(len(matchA.Users)) < (matchSizeA+matchSizeB) && uint(len(gb.users)) >= matchSizeB {
			matchB := randomChoices(gb, matchSizeB, innerGroupConstraints, forbiddenConnections)

			match = true

			for _, u := range matchB.Users {
				u := u
				if Filter(gb, users, &u, interGroupConstraints, forbiddenConnections) && Filter(gb, matchA.Users, &u, innerGroupConstraints, forbiddenConnections) {

					matchA.Users = append(matchA.Users, u)

				} else {
					match = false
					gb.RemoveUser(&u)
					break
				}
				gb.RemoveUser(&u)
			}
		}

	}

	return matchA

}

func doubleCheck(A []*entity.User, B []*entity.User) ([]*entity.User, []*entity.User) {
	APrime := []*entity.User{}
	for _, u1 := range A {
		u1 := u1
		if find, _ := Search(B, u1); !find {
			APrime = append(APrime, u1)
		}
	}
	return APrime, B
}

func Matcher(g *UserGraph, k uint,
	constraints []Constraint, SELECTOR Selector, forbidenconections [][]entity.User,
	A []*entity.User, B []*entity.User,
	interGroupConstraints []Constraint, innerGroupConstraints []Constraint) map[int]*Match {

	/* Matcher without constraint

	   input : g User's graph, k length of tuple match
	   matchSizeA, matchSizeB: size of matching for A and B respectivily; variable for group selector

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
			matched := randomChoices(g, k, constraints, forbidenconections)
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
		if uint(len(A)) > uint(len(B)) {
			A, B = doubleCheck(A, B)
		} else {
			B, A = doubleCheck(B, A)
		}
		if k < 2 {
			break
		}
		var matchSizeA uint
		var matchSizeB uint
		maxMatchSizeA := minimum(uint(len(A)), k-1)
		maxMatchSizeB := minimum(uint(len(B)), k-1)
		if maxMatchSizeA < maxMatchSizeB {
			matchSizeA = minimum(maxMatchSizeA, k/2)
			matchSizeB = minimum(k-matchSizeA, maxMatchSizeB)
		} else {
			matchSizeB = minimum(maxMatchSizeB, k/2)
			matchSizeA = minimum(k-matchSizeB, maxMatchSizeA)
		}

		if matchSizeB == 0 || matchSizeA == 0 || matchSizeB+matchSizeA != k {
			break
		}
		i := 0

		groupA := g.Subgraph(A)
		groupB := g.Subgraph(B)
		fmt.Println(groupA.users)
		fmt.Println(groupB.users)
		fmt.Println("")
		for uint(len(groupA.users))/matchSizeA > 0 && uint(len(groupB.users))/matchSizeB > 0 {

			matched := RandSubGroup(groupA, groupB, matchSizeA, matchSizeB,
				interGroupConstraints, innerGroupConstraints,
				forbidenconections)
			fmt.Println(groupA.users)
			fmt.Println(groupB.users)
			fmt.Println("")
			if matched != nil {
				for _, match := range matched.Users {
					match := match
					groupB.RemoveUser(&match)
					groupA.RemoveUser(&match)

				}

			}
			fmt.Println(groupA.users)
			fmt.Println(groupB.users)
			fmt.Println("")
			matching[i] = matched
			i++

		}

	}

	return matching

}

func GenerateTuple(users []entity.User, connections [][]entity.User, s Selector,
	forbiddenConnections [][]entity.User, size uint,
	A []entity.User, B []entity.User) []Match {
	/*
		         Input :
				    - general
					  connections: Connections plan in the graph;
					  s   : type of selector : basic, group .....
					- specific
					   users: Users for matching for group selector basic
					  *constraint
					   forbiddenConnecitons: variable for forbideenconnections constraints
					  *selector
					    size : size of matching; variable for basic selector
					    A,B : groupe A and B of  user; variable for group selector
	*/
	var results []Match
	var tuples map[int]*Match

	switch s {
	case Basic:
		graph := UsersToGraph(users, connections)
		graph.String()
		tuples = Matcher(graph, size, []Constraint{Unique, ForbiddenConnections}, Basic,
			forbiddenConnections, []*entity.User{}, []*entity.User{},
			[]Constraint{}, []Constraint{})
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
		graph := UsersToGraph(append(A, B...), connections)
		tuples = Matcher(graph, size, []Constraint{}, Group,
			forbiddenConnections, gA, gB,

			[]Constraint{Unique, ForbiddenConnections}, []Constraint{})

	}
	for index, matching := range tuples {
		if len(matching.Users) > 1 {
			results = append(results, Match{Id: strconv.Itoa(index), Users: matching.Users})
		}
	}
	return results
}
