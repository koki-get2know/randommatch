package database

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
)

func CreateJobStatus(uid string) error {
	driver, err := Driver()
	if err != nil {
		return err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MERGE (n: Job{uid: $uid}) "+
			"ON CREATE SET n += {uid: $uid, status: $status, "+
			"creation_date: datetime({timezone: 'Z'}), last_update: datetime({timezone: 'Z'})}",
			map[string]interface{}{"uid": uid, "status": Running})

		if err != nil {
			return nil, err
		}
		_, err = result.Consume()

		if err != nil {
			return nil, err
		}
		return nil, nil

	})
	if err != nil {
		return err
	}
	return nil
}

func updateJobStatus(session neo4j.Session, uid string, status jobStatus) error {
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (n: Job{uid: $uid}) "+
			"SET n += {status: $status, "+
			"last_update: datetime({timezone: 'Z'})}",
			map[string]interface{}{"uid": uid, "status": status})

		if err != nil {
			return nil, err
		}
		_, err = result.Consume()

		if err != nil {
			return nil, err
		}
		return nil, nil

	})
	if err != nil {
		return err
	}
	return nil
}

func GetJobStatus(uid string) (string, error) {
	driver, err := Driver()
	if err != nil {
		return "", err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	res, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (n: Job{uid: $uid}) RETURN n",
			map[string]interface{}{"uid": uid})

		if err != nil {
			return nil, err
		}
		if result.Next() {
			return result.Record().Values[0].(dbtype.Node).Props["status"], nil
		}

		return nil, result.Err()

	})
	if err != nil {
		return "", err
	}
	return res.(string), nil
}
