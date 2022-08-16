package database

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/koki/randommatch/entity"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
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

func DeleteUser(id string) error{
	driver, err := Driver()
	if err != nil {
		return err
	}

	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (u:User{uid:$uid}) "+
			"DETACH DELETE u ",
			map[string]interface{}{"uid": id})

		if err != nil {
			return nil, err
		}

		return result.Consume()
	})
	return err
}

func DeleteUsers() error {
	driver, err := Driver()
	if err != nil {
		return err
	}

	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (u:User) "+
			"DETACH DELETE u ",
			map[string]interface{}{})

		if err != nil {
			return nil, err
		}

		return result.Consume()
	})
	return err
}

func mapUsers(users []entity.User) []map[string]interface{} {
	var result = make([]map[string]interface{}, len(users))

	for index, item := range users {
		item.Id = uuid.New().String()
		result[index] = structs.Map(item)
	}

	return result
}

func mapMatches(tuples [][]entity.User) [][]map[string]interface{} {

	var result = make([][]map[string]interface{}, len(tuples))
	for index, item := range tuples {
		users := []map[string]interface{}{}
		for _, user := range item {
			users = append(users, structs.Map(user))
		}
		result[index] = users
	}
	fmt.Println(result)
	return result
}
func CreateUsers(users []entity.User, orgaUid string) (string, error) {
	jobId := uuid.New().String()
	if err := CreateJobStatus(jobId); err != nil {
		return "", err
	}
	status := make(chan JobStatus)
	go func() {
		if err := createUsers(users, orgaUid, status); err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		driver, err := Driver()
		if err != nil {
			fmt.Print("Driver error", err)
			return
		}
		session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
		defer session.Close()
		for st := range status {
			if err := updateJobStatus(session, jobId, st); err != nil {
				fmt.Println("Error while updating job", jobId, err)
			}
		}
	}()
	return jobId, nil
}

func createUsers(users []entity.User, orgaUid string, out chan JobStatus) error {
	defer close(out)
	driver, err := Driver()
	if err != nil {
		out <- Failed
		return err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	if len(users) == 0 {
		out <- Done
		return nil
	}

	out <- Running
	const chunkSize = 1000
	userschunks := chunkSlice(users, chunkSize)

	for _, chunk := range userschunks {
		_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			result, err := tx.Run("MATCH (o: Organization{uid: $orguid }) "+
				"UNWIND $users AS user "+
				"MERGE (u: User{name: user.Name, email: user.Email}) "+
				"ON CREATE SET u += {uid: user.Id, "+
				"creation_date: datetime({timezone: 'Z'}), last_update: datetime({timezone: 'Z'})} "+
				"MERGE (u)-[ruo:WORKS_FOR]->(o) "+
				"ON CREATE SET ruo.since = datetime({timezone: 'Z'}) "+
				"WITH user.Groups AS tags, u AS u "+
				"UNWIND tags AS tag "+
				"MERGE (t: Tag {name: tag}) "+
				"MERGE (u)-[rut:HAS_TAG]->(t) "+
				"RETURN u.uid",
				map[string]interface{}{"users": mapUsers(chunk), "orguid": orgaUid})

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
			out <- Failed
			return err
		}
	}
	out <- Done
	return nil
}

// Getlink  get all relationship from DB
func GetLink() ([][]entity.User, error) {
	driver, err := Driver()
	if err != nil {
		return [][]entity.User{}, err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	links, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (n:User)-[r:MET]->(ou:User) RETURN n, ou",
			map[string]interface{}{})

		link := [][]entity.User{}
		if err != nil {
			return link, err
		}
		for result.Next() {
			user := result.Record().Values[0].(dbtype.Node).Props
			ou := result.Record().Values[1].(dbtype.Node).Props
			var users []entity.User
			users = append(users,
				entity.User{
					Id:   user["uid"].(string),
					Name: user["name"].(string),
					//Groups: tags,
				}, entity.User{
					Id:   ou["uid"].(string),
					Name: ou["name"].(string),
					//Groups: tags,
				})

			link = append(link, users)
		}

		if result.Err() != nil {
			return link, result.Err()
		}

		return link, nil

	})
	if err != nil {
		return [][]entity.User{}, err
	}
	return links.([][]entity.User), nil
}

// CreateLink create relationship (known) between 2 users in BD
func CreateLink(tuples [][]entity.User) error {
	driver, err := Driver()
	if err != nil {
		return err
	}

	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("UNWIND $tuples AS tuple "+
			" MATCH (a:User{uid:tuple[0].Id}),(b:User{uid:tuple[1].Id}) "+
			"MERGE (a)-[r:MET]-(b) "+
			"ON CREATE SET r.on = datetime({timezone: 'Z'}) ",
			map[string]interface{}{"tuples": mapMatches(tuples)})

		if err != nil {
			return nil, err
		}

		return result.Consume()
	})

	return err
}

func GetUsers() ([]entity.User, error) {
	driver, err := Driver()
	if err != nil {
		return []entity.User{}, err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	users, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (n: User) OPTIONAL MATCH (n)-[r:HAS_TAG]->(t: Tag) RETURN  n, COLLECT(t.name)",
			map[string]interface{}{})
		var users []entity.User

		if err != nil {
			return users, err
		}
		for result.Next() {
			var tags []string
			for _, tag := range result.Record().Values[1].([]interface{}) {
				tags = append(tags, tag.(string))
			}
			user := result.Record().Values[0].(dbtype.Node).Props

			users = append(users,
				entity.User{
					Id:     user["uid"].(string),
					Name:   user["name"].(string),
					Groups: tags,
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

func GetEmailsFromUIds(uids []string) (map[string]string, error) {
	driver, err := Driver()
	if err != nil {
		return map[string]string{}, err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	mapUidEmail, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (n: User) WHERE n.uid IN $uids RETURN  n.uid, n.email",
			map[string]interface{}{"uids": uids})

		if err != nil {
			return map[string]string{}, err
		}
		var uidEmails = make(map[string]string)
		for result.Next() {
			uidEmails[result.Record().Values[0].(string)] = result.Record().Values[1].(string)
		}

		if result.Err() != nil {
			return map[string]string{}, result.Err()
		}

		return uidEmails, nil

	})

	if err != nil {
		return map[string]string{}, err
	}
	return mapUidEmail.(map[string]string), nil
}
