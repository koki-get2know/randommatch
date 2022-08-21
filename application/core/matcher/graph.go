// Package graph creates a ItemGraph data structure for the Item type
package matcher

import (
	"log"

	"github.com/koki/randommatch/entity"
)

// UserGraph the Items graph
type UserGraph struct {
	users []*entity.User
	edges map[string][]*entity.User
}

// AddNode adds a node to the graph

func (g *UserGraph) AddUser(n *entity.User) {

	if find, _ := Search(g.users, n); !find {
		g.users = append(g.users, n)
	}

}

// AddEdge adds an edge to the graph
func (g *UserGraph) AddEdge(n1, n2 *entity.User) {
	if g.edges == nil {
		g.edges = make(map[string][]*entity.User)
	}
	if find, _ := Search(g.edges[(*n1).Id], n2); !find {
		g.edges[(*n1).Id] = append(g.edges[(*n1).Id], n2)
	}
	if find, _ := Search(g.edges[(*n2).Id], n1); !find {
		g.edges[(*n2).Id] = append(g.edges[(*n2).Id], n1)
	}

}

// search a user in a list of user
func Search(users []*entity.User, n *entity.User) (bool, int) {
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

// SearchNode findout a specifique node in a graph

func (g *UserGraph) SearchUser(n *entity.User) (bool, int) {
	index := -1
	find := false
	if g.users != nil {
		find, index = Search(g.users, n)
	}

	return find, index

}

// remove a user in a list of users

func Remove(s []*entity.User, i int) []*entity.User {
	if i != len(s)-1 {
		s[i] = s[len(s)-1]
	}

	return s[:len(s)-1]
}

// RemoveEdge remove an edge from the graph
func (g *UserGraph) RemoveEdge(n *entity.User) {
	for _, user := range g.users {
		find, index := Search(g.edges[(*user).Id], n)
		if find {
			g.edges[(*user).Id] = Remove(g.edges[(*user).Id], index)
		}
		delete(g.edges, (*n).Id)
	}
}

// RemoveUser remove a user from the graph
func (g *UserGraph) RemoveUser(n *entity.User) {

	g.RemoveEdge(n)
	find, index := g.SearchUser(n) // find out the index of this node
	if find {
		g.users = Remove(g.users, index)

	}

}

// print the graph
func (g *UserGraph) String() {
	s := ""
	for _, usernode := range g.users {
		s += usernode.String() + " -> "
		near := g.edges[(*usernode).Id]
		for _, user := range near {
			s += user.String() + " "
		}
		s += "\n"
	}
	log.Println(s)

}

// UsersToGraph create a graph of some users and connections

func UsersToGraph(users []entity.User, connections [][]entity.User) *UserGraph {

	/* input :
	         users: table users for matching
			 connections: Matrix n*2for connection in the graph where each row is a pair of 2 users connected in the graph;
	   output: return the graph
	*/
	var graph UserGraph
	for _, user := range users {
		user := user
		graph.AddUser(&user)
	}
	for _, usersAlreadyMatch := range connections {
		if len(usersAlreadyMatch) > 0 {
			node := usersAlreadyMatch[0]
			for _, user := range usersAlreadyMatch[1:] {

				user := user
				graph.AddEdge(&node, &user)

			}
		}
	}
	return &graph
}

// Subgraph extract a subgraph G' from a graph G

func (g *UserGraph) Subgraph(users []*entity.User) *UserGraph {
	var subG UserGraph
	subG.edges = make(map[string][]*entity.User)
	subG.users = users

	for _, user := range users {
		subG.edges[(*user).Id] = g.edges[(*user).Id]

	}
	return &subG
}
