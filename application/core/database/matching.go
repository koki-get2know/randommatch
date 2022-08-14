package database

import (
	"github.com/google/uuid"
	"github.com/koki/randommatch/entity"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
)

func CreateMatchingStat(MatchingStats entity.Matching) (string, error) {
	driver, err := Driver()
	if err != nil {
		return "", err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	uid, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MERGE (n:MatchingStat {uid: $id}) "+
			"ON CREATE SET n += {uid: $id, numgroups: $numgroups, numConversations: $numconvs, numFaileds: $numfailed, "+
			"creation_date: datetime({timezone: 'Z'}), last_update: datetime({timezone: 'Z'})} "+
			"RETURN n.uid",
			map[string]interface{}{
				"id":        uuid.New().String(),
				"numgroups": MatchingStats.NumGroups,
				"numconvs":  MatchingStats.NumConversations,
				"numfailed": MatchingStats.NumFailed,
			})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return "", err
	}

	return uid.(string), err
}

func GetMatchings() ([]entity.Matching, error) {
	driver, err := Driver()
	if err != nil {
		return []entity.Matching{}, err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	matchings, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (n:MatchingStat) RETURN n", map[string]interface{}{})
		var matchings []entity.Matching

		if err != nil {
			return matchings, err
		}
		for result.Next() {
			matching := result.Record().Values[0].(dbtype.Node).Props

			matchings = append(matchings,
				entity.Matching{
					Id:               matching["uid"].(string),
					NumGroups:        int(matching["numgroups"].(int64)),
					NumConversations: int(matching["numConversations"].(int64)),
					NumFailed:        int(matching["numFaileds"].(int64)),
				})
		}

		if result.Err() != nil {
			return matchings, result.Err()
		}

		return matchings, nil

	})
	if err != nil {
		return []entity.Matching{}, err
	}
	return matchings.([]entity.Matching), nil
}
