package database

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type frequency string

const (
	monthly         frequency = "monthly"
	weekly          frequency = "weekly"
	daily           frequency = "daily"
	hourly          frequency = "hourly"
	every_half_hour frequency = "every_half_hour"
	every_ten_min   frequency = "every_ten_min"
)

func CreateSchedule(uid string) error {
	driver, err := Driver()
	if err != nil {
		return err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MERGE (n: Schedule{uid: $uid}) "+
			"ON CREATE SET n += {uid: $uid, active: $active, "+
			"creation_date: datetime({timezone: 'Z'}), last_update: datetime({timezone: 'Z'})}, "+
			"frequency: $frequency, start_time: datetime({timezone: 'Z'})",
			map[string]interface{}{"uid": uid, "active": true, "frequency": monthly})

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
