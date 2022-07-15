package matcher

import (
    "fmt"
    "testing"
    
)


func TestRandomChoices(t *testing.T) {
   //id := []string{"2", "5", "6", "8", "10", "12", "24", "25"}
    g.String()
    matchedUsers := RandomChoices(&g,3)
    fmt.Printf("Match of %d: [" ,len(matchedUsers))
    for i := 0; i < len(matchedUsers); i++{
       fmt.Printf("%s,",matchedUsers[i].String())
    }
    
    fmt.Printf("]")
}

