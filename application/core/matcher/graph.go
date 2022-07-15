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
    return fmt.Sprintf("%s", n.userId)
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
    g.users = append(g.users, n)
    g.lock.Unlock()
}

// AddEdge adds an edge to the graph
func (g *UserGraph) AddEdge(n1, n2 *User) {
    g.lock.Lock()
    if g.edges == nil {
        g.edges = make(map[User][]*User)
    }
    g.edges[*n1] = append(g.edges[*n1], n2)
    g.edges[*n2] = append(g.edges[*n2], n1)
    g.lock.Unlock()
}

// search a user in a list of user 
func Search(s []*User, n *User)(bool,int){
    index := -1 
    find := false
    for i := 0; i < len(s); i++ {
            if(g.users[i].userId == n.userId){
                  find = true
                  index = i
                  break
            }
      }
      
   return find, index
}

// SearchNode findout a specifique node in a graph

func (g *UserGraph) SearchUser(n *User) (bool, int){
     g.lock.RLock()
     index := -1 
     find := false
     fmt.Printf("Je commence la recherche")
     if g.users != nil {
        fmt.Printf("Je commence la recherche")
        find, index = Search(g.users,n)
      }
      g.lock.RUnlock()
      return find, index

}

// remove a user in a list of users
   
func Remove(s []*User, i int) []*User {
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
}



// TODO RemoveEdge remove an edge from the graph
/*func (g *UserGraph) RemoveEdge(n *User){
      g.lock.Lock()
      
      g.lock.Unlock()
         
}
*/

// TODO RemoveUser remove a user from the graph
/*func (g *UserGraph) RemoveUser(n *User){
      g.lock.Lock()
      find, index := g.SearchUser(n) // find out the index of this node
      fmt.Printf("index:%d ", index)
      if find{
       g.users[index] = g.users[len(g.users) -1] // remplace the contains by the last constains' element
       g.users = g.users[:len(g.users)-1] // slice the last
       delete(g.edges, *n)// delete the corresponding edges
      }
      g.lock.Unlock()

}
*/

// print the graph
func (g *UserGraph) String() {
    g.lock.RLock()
    s := ""
    for i := 0; i < len(g.users); i++ {
        s += g.users[i].String() + " -> "
        near := g.edges[*g.users[i]]
        for j := 0; j < len(near); j++ {
            s += near[j].String() + " "
        }
        s += "\n"
    }
    fmt.Println(s)
    g.lock.RUnlock()
}



