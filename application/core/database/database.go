package database

import (
	"fmt"
	"os"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

//var once sync.Once
var driver *neo4j.Driver
var err error

func Driver() (*neo4j.Driver, error) {
	if driver == nil {
	//once.Do(func() {

		creds := strings.Split(os.Getenv("NEO4J_AUTH"), "/")
		if len(creds) < 2 {
			err = fmt.Errorf("NEO4J_AUTH env variable missing or not set correctly")
			return nil, err
		}
		var dr neo4j.Driver

		dbhost, found := os.LookupEnv("DB_HOST")
		connectionstring, ok := os.LookupEnv("NEO4J_CNX_STRING")

		if !ok {
			if !found {
				fmt.Println("neither NEO4J_CNX_STRING nor DB_HOST env variable set, defaulting to localhost")
				dbhost = "localhost"
			}
			connectionstring = fmt.Sprintf("bolt://%v:7687", dbhost)
		}
		dr, err = neo4j.NewDriver(connectionstring, neo4j.BasicAuth(creds[0], creds[1], ""))

		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		driver = &dr
	//})
	}
	return driver, err
}
