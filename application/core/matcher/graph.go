// Package graph creates a ItemGraph data structure for the Item type
package matcher

import (
	"fmt"
	"sync"
	//"github.com/cheekybits/genny/generic"
)

// Item the type of the binary search tree
//type Item generic.Type

// Node a single node that composes the tree
type User struct {
	userId string
}

func (n *User) String() string {
	return n.userId
}

// UserGraph the Items graph
type UserGraph struct {
	users []*User
	edges map[User][]*User
	lock  sync.RWMutex
}

// AddNode adds a node to the graph
func (g *UserGraph) AddUser(n *User) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.users = append(g.users, n)

}

// AddEdge adds an edge to the graph
func (g *UserGraph) AddEdge(n1, n2 *User) {
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.edges == nil {
		g.edges = make(map[User][]*User)
	}
	g.edges[*n1] = append(g.edges[*n1], n2)
	g.edges[*n2] = append(g.edges[*n2], n1)

}

// search a user in a list of user
func Search(users []*User, n *User) (bool, int) {
	index := -1
	find := false
	for i, user := range users {
		if user.userId == n.userId {
			find = true
			index = i
			break
		}
	}

	return find, index
}

// SearchNode findout a specifique node in a graph

func (g *UserGraph) SearchUser(n *User) (bool, int) {
	g.lock.RLock()
	defer g.lock.RUnlock()
	index := -1
	find := false
	if g.users != nil {
		find, index = Search(g.users, n)
	}

	return find, index

}

// remove a user in a list of users

func Remove(s []*User, i int) []*User {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// RemoveEdge remove an edge from the graph
func (g *UserGraph) RemoveEdge(n *User) {
	g.lock.Lock()
	defer g.lock.Unlock()
	for _, user := range g.users {
		find, index := Search(g.edges[*user], n)
		if find {
			g.edges[*user] = Remove(g.edges[*user], index)
		}
		delete(g.edges, *n)
	}
}

// RemoveUser remove a user from the graph
func (g *UserGraph) RemoveUser(n *User) {

	g.RemoveEdge(n)
	find, index := g.SearchUser(n) // find out the index of this node
	g.lock.Lock()
	defer g.lock.Unlock()
	if find {
		g.users = Remove(g.users, index)

	}

}

// print the graph
func (g *UserGraph) String() {
	g.lock.RLock()
	defer g.lock.RUnlock()
	s := ""
	for i := 0; i < len(g.users); i++ {
		s += g.users[i].String() + " -> "
		near := g.edges[*g.users[i]]
		for _, user := range near {
			s += user.String() + " "
		}
		s += "\n"
	}
	fmt.Println(s)

}
