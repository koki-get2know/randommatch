package database

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/koki/randommatch/entity"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
)

func CreateMatchingStat(MatchingStats entity.MatchingStat, orgaName string) (string, error) {
	driver, err := Driver()
	if err != nil {
		return "", err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	uid, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		//orgId parameter
		result, err := tx.Run("MATCH (o: Organization{lower_name: $orgname }) "+
			"CREATE (n:MatchingStat {uid: $id, num_groups: $numgroups, num_conversations: $numconvs, num_failures: $numfailed, "+
			"creation_date: datetime({timezone: 'Z'}), last_update: datetime({timezone: 'Z'})}) "+
			"MERGE (n)-[ruo:BELONGS_TO]->(o) "+
			"RETURN n.uid",
			map[string]interface{}{
				"id":        uuid.New().String(),
				"numgroups": MatchingStats.NumGroups,
				"numconvs":  MatchingStats.NumConversations,
				"numfailed": MatchingStats.NumFailed,
				"orgname":   strings.ToLower(orgaName),
			})

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

func GetMatchingStats(organization string) ([]entity.MatchingStat, error) {
	driver, err := Driver()
	if err != nil {
		return []entity.MatchingStat{}, err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	matchings, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (n:MatchingStat) "+
			"MATCH (n)-[BELONGS_TO]->(o:Organization{lower_name: $lower_name}) "+
			"RETURN n",
			map[string]interface{}{"lower_name": strings.ToLower(organization)})
		var matchings []entity.MatchingStat

		if err != nil {
			return matchings, err
		}
		for result.Next() {
			matching := result.Record().Values[0].(dbtype.Node).Props

			matchings = append(matchings,
				entity.MatchingStat{
					Id:               matching["uid"].(string),
					NumGroups:        int(matching["num_groups"].(int64)),
					NumConversations: int(matching["num_conversations"].(int64)),
					NumFailed:        int(matching["num_failures"].(int64)),
					CreatedAt:        matching["creation_date"].(time.Time),
				})
		}
		if result.Err() != nil {
			return matchings, result.Err()
		}

		return matchings, nil

	})

	return matchings.([]entity.MatchingStat), err
}
