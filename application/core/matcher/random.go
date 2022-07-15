package matcher

import (
    "time"
    "math/rand"
)
const SELECTOR = "random"

type Matching struct {
       matched []*User
}
/* random choice without constraint
      input : table of id
      output : k random id 
      purpose : match k user from the graph
      
*/
    
func RandomChoices(g *UserGraph, k int) []*User {
     var matchedUsers []*User
     rand.Seed(time.Now().UnixNano()) // initialize the seed to get 
     for i := 0; i < k; i++{
       matchedUsers = append (matchedUsers,g.users[rand.Intn(len(g.users))])
     
     
     }
     return matchedUsers
}

// TODO the selector paramater should be a config variable
/*func Matcher(g *UserGraph) map[int]Matching {
      matching := make(map[int]Matching) 

        if ( SELECTOR == "random" ){
            
        
        }
              
     return matching 

}
*/

