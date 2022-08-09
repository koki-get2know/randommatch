package database

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/koki/randommatch/entity"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
)

type jobStatus string

const (
	Pending   jobStatus = "Pending"
	Running   jobStatus = "Running"
	Done      jobStatus = "Done"
	Failed    jobStatus = "Failed"
	Suspended jobStatus = "Suspended"
	Cancelled jobStatus = "Cancelled"
)

func CreateUser(user entity.User) (string, error) {
	driver, err := Driver()
	if err != nil {
		return "", err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	uid, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MERGE (n: User{name: $name, email: $email}) "+
			"ON CREATE SET n += {uid: $uid, "+
			"creation_date: datetime({timezone: 'Z'}), last_update: datetime({timezone: 'Z'})} "+
			"RETURN n.uid",
			map[string]interface{}{"name": user.Name, "uid": uuid.New().String(), "email": user.Email})

		if err != nil {
			return "", err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return "", result.Err()
	})
	if err != nil {
		return "", err
	}

	return uid.(string), err
}

func chunkSlice(slice []entity.User, chunkSize int) [][]entity.User {
	size := len(slice)
	if size <= chunkSize {
		return [][]entity.User{slice}
	}

	var chunks [][]entity.User
	for i := 0; i < size; i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > size {
			end = size
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func mapUsers(users []entity.User) []map[string]interface{} {
	var result = make([]map[string]interface{}, len(users))

	for index, item := range users {
		result[index] = structs.Map(item)
	}

	return result
}

func CreateUsers(users []entity.User, jobId string) error {

	driver, err := Driver()
	if err != nil {
		return err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	if len(users) == 0 {
		updateJobStatus(session, jobId, Done)
		return nil
	}

	if err = updateJobStatus(session, jobId, Running); err != nil {
		return err
	}
	const chunkSize = 1000
	userschunks := chunkSlice(users, chunkSize)

	for _, chunk := range userschunks {
		_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			result, err := tx.Run("UNWIND $users AS user "+
				"MERGE (n: User{name: user.Name, email: user.Email}) "+
				"ON CREATE SET n += {uid: $uid, "+
				"creation_date: datetime({timezone: 'Z'}), last_update: datetime({timezone: 'Z'})} "+
				"RETURN n.uid",
				map[string]interface{}{"users": mapUsers(chunk), "uid": uuid.New().String()})

			if err != nil {
				return "", err
			}
			var rows []string
			for result.Next() {
				rows = append(rows, (result.Record().Values[0]).(string))
			}

			if result.Err() != nil {
				return "", result.Err()
			}

			return strings.Join(rows, "|"), nil

		})
		if err != nil {
			updateJobStatus(session, jobId, Failed)
			return err
		}
	}
	return updateJobStatus(session, jobId, Done)
}

func GetUsers() ([]entity.User, error) {
	driver, err := Driver()
	if err != nil {
		return []entity.User{}, err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	users, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (n: User) RETURN n", map[string]interface{}{})
		var users []entity.User

		if err != nil {
			return users, err
		}
		for result.Next() {
			user := result.Record().Values[0].(dbtype.Node).Props
			users = append(users,
				entity.User{
					Id:   user["uid"].(string),
					Name: user["name"].(string),
				})
		}

		if result.Err() != nil {
			return users, result.Err()
		}

		return users, nil

	})
	if err != nil {
		return []entity.User{}, err
	}
	return users.([]entity.User), nil
}
