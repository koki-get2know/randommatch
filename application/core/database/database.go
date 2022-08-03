package database

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

var once sync.Once
var driver *neo4j.Driver

func Driver() (*neo4j.Driver, error) {
	var err error
	once.Do(func() {
		creds := strings.Split(os.Getenv("NEO4J_AUTH"), "/")
		if len(creds) < 2 {
			err = fmt.Errorf("NEO4J_AUTH env variable missing or not set correctly")
			return
		}
		var dr neo4j.Driver
		dr, err = neo4j.NewDriver("bolt://match-db:7687", neo4j.BasicAuth(creds[0], creds[1], ""))
		if err != nil {
			return
		}
		driver = &dr
	})

	return driver, err
}
