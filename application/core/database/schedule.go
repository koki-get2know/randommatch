package database

import (
	"log"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/koki/randommatch/entity"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
)

func ScheduleLinkTTags(tag string, scheduleCode string) error {
	driver, err := Driver()
	if err != nil {
		return err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(
			"MATCH (t: TechTag{lower_name: $lower_tag_name})"+
				"MATCH (n: Schedule{uid: $uid})"+
				"MERGE (n)-[rut:HAS_TECH_TAG]->(t)",
			map[string]interface{}{"uid": scheduleCode, "lower_tag_name": strings.ToLower(tag)})
		return "", err
	})

	return err
}
func ScheduleLinkTags(tag string, scheduleCode string) error {
	driver, err := Driver()
	if err != nil {
		return err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(
			"MATCH (t: Tag{lower_name: $lower_tag_name})"+
				"MATCH (n: Schedule{uid: $uid})"+
				"MERGE (n)-[rut:HAS_TECH_TAG]->(t)",
			map[string]interface{}{"uid": scheduleCode, "lower_tag_name": strings.ToLower(tag)})
		return "", err
	})

	return err
}
func GetSchedule(uid string, orga string) (entity.Schedule, error) {
	driver, err := Driver()
	if err != nil {
		return entity.Schedule{}, err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	schedule, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (s:Schedule{uid:$uid})-[r:SCHEDULE_FOR]->(:Organization{lower_name:$orga})"+
			"RETURN s",
			map[string]interface{}{"uid": uid, "orga": strings.ToLower(orga)})
		var schedule entity.Schedule

		if err != nil {
			return schedule, err
		}
		if result.Next() {
			sch := result.Record().Values[0].(dbtype.Node).Props

			schedule = entity.Schedule{
				Id:           sch["uid"].(string),
				Name:         sch["name"].(string),
				CreateDate:   sch["creation_date"].(time.Time),
				LastUpdate:   sch["last_update"].(time.Time),
				StartDate:    sch["start_time"].(time.Time),
				EndDate:      sch["end_time"].(time.Time),
				Size:         sch["size"].(int64),
				MatchingType: sch["matching_type"].(string),
				Active:       sch["active"].(bool),
			}
		}

		if result.Err() != nil {
			return schedule, result.Err()
		}

		return schedule, nil

	})
	if err != nil {
		return entity.Schedule{}, err
	}
	return schedule.(entity.Schedule), nil
}

func CreateSchedule(schedule entity.Schedule, orga string) error {
	driver, err := Driver()
	if err != nil {
		return err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	log.Println(structs.Map(schedule))
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (o: Organization{lower_name: $lower_org_name})"+
			"MERGE (n: Schedule{uid: $uid})"+
			"ON CREATE SET n += {uid: $uid, name:$name,active: $active, "+
			"creation_date: datetime({timezone: 'Z'}), last_update: datetime({timezone: 'Z'}), "+
			"start_time: datetime({timezone: 'Z'}),"+
			"end_time: datetime({timezone: 'Z'}),"+
			"size:$size,frequency:$frequency,matching_type: $matchingType}"+
			"MERGE (n)-[r:SCHEDULE_FOR]->(o)",

			map[string]interface{}{"size": schedule.Size, "frequency": schedule.Frequency, "uid": schedule.Id, "name": schedule.Name, "active": schedule.Active, "matchingType": schedule.MatchingType, "lower_org_name": strings.ToLower(orga)})
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
