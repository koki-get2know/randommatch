package database

import (
	"github.com/koki/randommatch/entity"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
)

func GetTags() ([]entity.Tag, error) {
	driver, err := Driver()
	if err != nil {
		return []entity.Tag{}, err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	tags, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (n: Tag) RETURN n",
			map[string]interface{}{})
		var tags []entity.Tag

		if err != nil {
			return tags, err
		}

		for result.Next() {
			tag := result.Record().Values[0].(dbtype.Node).Props

			tags = append(tags,
				entity.Tag{
					Name: tag["name"].(string),
				})
		}

		if result.Err() != nil {
			return tags, result.Err()
		}

		return tags, nil

	})

	return tags.([]entity.Tag), err
}
