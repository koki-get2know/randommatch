package database

import (
	"log"
	"strings"

	"github.com/koki/randommatch/entity"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func CreateTechTags(tag string) error {
	driver, err := Driver()
	if err != nil {
		return err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(
			"MERGE (t: TechTag {lower_name: $lower_name}) "+
				"ON CREATE SET t += {name: $tag, lower_name:$lower_name}",
			map[string]interface{}{"tag": tag, "lower_name": strings.ToLower(tag)})
		return "", err
	})

	return err
}

func UserLinkTags(users []entity.User, tag string) error {
	driver, err := Driver()
	if err != nil {
		return err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	log.Println(MapUsers(users))
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(
			"MATCH (t: TechTag{lower_name:$lower_tag_name})"+
				"UNWIND $users AS user "+
				"MATCH (u:User{uid: user.Id})"+
				"MERGE (u)-[rut:HAS_TECH_TAG]->(t)",
			map[string]interface{}{"users": MapUsers(users), "lower_tag_name": strings.ToLower(tag)})
		return "", err
	})

	return err

}
