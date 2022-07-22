package matcher


import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

//TODO should be another the one of parameters' config
const PATHDATA = "users.json"
const PATHRESULT = "output.json"

// Users struct which contains
// an array of users


type UserLink struct {
    Users []Users `json:"users"`
    Links []Link  `json:"links"`
}

type Users struct{
     UserId string `json:"userId"` 
}

type Link struct {
    UserId_1  string `json:"id_1"`
    UserId_2  string `json:"id_2"`
}





// BuildFromJson Make user Graph from Json file
func  (g *UserGraph) BuildFromJson() {
    // Open our jsonFile
    jsonFile, err := ioutil.ReadFile("users.json")
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("Successfully Opened users.json")
    // defer the closing of our jsonFile so that we can parse it later on
    //defer jsonFile.Close()

    // read our opened xmlFile as a byte array.
    //byteValue, _ := ioutil.ReadAll(jsonFile)

    // we initialize our Users array
    var userLink UserLink

    // we unmarshal our byteArray which contains our
    // jsonFile's content into 'users' which we defined above
    json.Unmarshal(jsonFile, &userLink)
    // we iterate through every user within our users array and
    // print out the user Type, their name, and their facebook url
    // as just an example
    fmt.Println("number of User : ", len(userLink.Users))
    fmt.Println("number of link : ", len(userLink.Links))
    var nodes []User
    for i := 0; i < len(userLink.Users); i++ {
        nodes = append(nodes,User{userLink.Users[i].UserId})
        g.AddUser(&nodes[i])

    }
    
    var node1, node2 []User
    for i := 0; i < len(userLink.Links); i++ {
        node1 = append(node1,User{userLink.Links[i].UserId_1})
        node2 = append(node2,User{userLink.Links[i].UserId_2})
        
        g.AddEdge(&node1[i],&node2[i])
    }
    
     
    

}

func WriteToJson(matching map[int]*Matching) {
    /*
       input : all matching
       purpose : serialize matching into Json file
    */
    //var matchingJson []byte
    jsonFile, err := os.Create(PATHRESULT)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer jsonFile.Close()

    for _, value := range matching {
        matchingJson, err := json.Marshal(value)
        if err != nil {
            fmt.Println(err)
            return
        }
        _, err = jsonFile.Write(matchingJson)
        if err != nil {
            fmt.Println(err)
            return
        }

    }
}




