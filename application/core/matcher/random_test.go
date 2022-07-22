package matcher

import (
    "fmt"
    "testing"
    
)


func TestRandomChoices(t *testing.T) {
   //id := []string{"2", "5", "6", "8", "10", "12", "24", "25"}
    g.String()
    constraint := []Constraints{dejavu}
    matching := RandomChoices(&g,2, constraint)
    fmt.Printf("Match of %d: [" ,len(matching.Matched))
    for i := 0; i < len(matching.Matched); i++{
       fmt.Printf("%s,",matching.Matched[i].String())
    }
    
    fmt.Printf("]")
}


func TestMatcher(t *testing.T) {
   //id := []string{"2", "5", "6", "8", "10", "12", "24", "25"}
    g.String()
    constraint := []Constraints{dejavu}
    matching := Matcher(&g,2,constraint)
    
    for j := 0; j < len(matching); j++{ 
      fmt.Printf("Match : [" )
      for i := 0; i < len(matching[j].Matched); i++{
         fmt.Printf("%s,",matching[j].Matched[i].String())
      }
    
      fmt.Printf("]")
      fmt.Println("")
      
    }
    g.String()
    WriteToJson(matching)
}



