package matcher

import (
    "fmt"
    "testing"
    
)


func TestRandomChoices(t *testing.T) {
   //id := []string{"2", "5", "6", "8", "10", "12", "24", "25"}
    g.String()
    constraint := []string{"deja vu"}
    matching := RandomChoices(&g,2, constraint)
    fmt.Printf("Match of %d: [" ,len(matching.matched))
    for i := 0; i < len(matching.matched); i++{
       fmt.Printf("%s,",matching.matched[i].String())
    }
    
    fmt.Printf("]")
}


func TestMatcher(t *testing.T) {
   //id := []string{"2", "5", "6", "8", "10", "12", "24", "25"}
    g.String()
    constraint := []string{"deja vu"}
    matching := Matcher(&g,2,constraint)
    
    for j := 0; j < len(matching); j++{ 
      fmt.Printf("Match : [" )
      for i := 0; i < len(matching[j].matched); i++{
         fmt.Printf("%s,",matching[j].matched[i].String())
      }
    
      fmt.Printf("]")
      fmt.Println("")
      
    }
    g.String()
}



