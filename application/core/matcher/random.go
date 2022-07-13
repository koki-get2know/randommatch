package matcher

import (
    "time"
    "math/rand"
)

/* random choice without constraint
      input : table of id
      output : 2 random id 
      purpose : match 2 persons
      
*/
    
func RandomChoices(ids []int) (int, int) {
     rand.Seed(time.Now().UnixNano()) // initialize the seed to get 
     match_1 :=  rand.Intn(len(ids))
     match_2 := rand.Intn(len(ids))
     return ids[match_1], ids[match_2]
}


