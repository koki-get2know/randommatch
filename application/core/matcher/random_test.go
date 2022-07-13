package matcher

import (
    "fmt"
    "testing"
    
)

func TestRandomChoices(t *testing.T) {
   id := []int{2, 5, 6, 8, 10, 12, 24, 25}
    match_1, match_2 := RandomChoices(id)
    fmt.Printf("id:%d match id:%d\n", match_1,match_2)
}

