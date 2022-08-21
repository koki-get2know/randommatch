package database

import (
	"strings"

	"github.com/google/uuid"
	"github.com/koki/randommatch/entity"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
)

func CreateOrganization(organization entity.Organization) (string, error) {
	driver, err := Driver()
	if err != nil {
		return "", err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	uid, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MERGE (n:Organization {lower_name: toLower($name)}) "+
			"ON CREATE SET n += {uid: $id, name: $name, description: $desc, "+
			"creation_date: datetime({timezone: 'Z'}), last_update: datetime({timezone: 'Z'})} "+
			"RETURN n.uid",
			map[string]interface{}{"name": organization.Name, "id": uuid.New().String(), "desc": organization.Description})

		if err != nil {
			return "", err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return "", result.Err()
	})

	return uid.(string), err
}

func GetOrganizationById(uid string) (entity.Organization, error) {
	driver, err := Driver()
	if err != nil {
		return entity.Organization{}, err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	res, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (n: Organization{uid: $uid}) RETURN n LIMIT 1",
			map[string]interface{}{"uid": uid})

		if err != nil {
			return entity.Organization{}, err
		}
		if result.Next() {
			return entity.Organization{Id: result.Record().Values[0].(dbtype.Node).Props["uid"].(string),
			 Name: result.Record().Values[0].(dbtype.Node).Props["name"].(string),
			 Description: result.Record().Values[0].(dbtype.Node).Props["description"].(string),
			 }, nil
		}

		return entity.Organization{}, result.Err()

	})

	return res.(entity.Organization), err
}

func GetOrganizationByName(name string) (string, error) {
	driver, err := Driver()
	if err != nil {
		return "", err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	res, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (n: Organization{lower_name: $name}) RETURN n",
			map[string]interface{}{"name": strings.ToLower(name)})

		if err != nil {
			return "", err
		}
		if result.Next() {
			return result.Record().Values[0].(dbtype.Node).Props["uid"], nil
		}

		return "", result.Err()

	})

	return res.(string), err
}