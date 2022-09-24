package database

import (
	"strings"
	"time"

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
				CreateDate:   sch["creationDate"].(time.Time),
				Time:         sch["time"].(string),
				Duration:     sch["duration"].(string),
				StartDate:    sch["startTime"].(time.Time),
				EndDate:      sch["endTime"].(time.Time),
				Size:         sch["size"].(int64),
				Frequency:    entity.Frequency(sch["frequency"].(string)),
				MatchingType: entity.MatchingType(sch["matchingType"].(string)),
				Active:       sch["active"].(bool),
				Week:         entity.Week(sch["week"].(string)),
				Days:         sch["days"].([]string),
				LastRun:      sch["lastRun"].(time.Time),
				NextRun:      sch["nextRun"].(time.Time),
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
	schedule.LastRun = time.Now().UTC()
	schedule.UdapteNextRun()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (o: Organization{lower_name: $lower_org_name})"+
			"MERGE (n: Schedule{uid: $uid})"+
			"ON CREATE SET n += {uid: $uid, name:$name,active: $active,"+
			"creationDate: datetime({timezone: 'Z'}), "+
			"startTime: $start,"+
			"endTime: $end,"+
			"size:$size,frequency:$frequency,matchingType:$matchingType,"+
			"week:$week, days:$days, lastRun:datetime({timezone: 'Z'}), nextRun:$nextRun,"+
			"time:$time, duration:$duration}"+
			"MERGE (n)-[r:SCHEDULE_FOR]->(o)",

			map[string]interface{}{"lastRun": schedule.LastRun, "start": schedule.StartDate.UTC(), "end": schedule.EndDate.UTC(), "time": schedule.Time, "duration": schedule.Duration, "nextRun": schedule.NextRun, "week": schedule.Week, "days": schedule.Days, "size": schedule.Size, "frequency": schedule.Frequency, "uid": schedule.Id, "name": schedule.Name, "active": schedule.Active, "matchingType": schedule.MatchingType, "lower_org_name": strings.ToLower(orga)})
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

func UpdateSchedule(schedule entity.Schedule, orga string) error {
	driver, err := Driver()
	if err != nil {
		return err
	}
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	schedule.LastRun = time.Now().UTC()
	schedule.UdapteNextRun()

	defer session.Close()
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(
			"MATCH (s:Schedule{uid:$uid})-[r:SCHEDULE_FOR]->(:Organization{lower_name:$orga})"+
				"SET s.LastRun = $lastRun"+
				"SET s.NextRun = $nextRun",

			map[string]interface{}{"lastRun": schedule.LastRun, "nextRun": schedule.NextRun, "uid": schedule.Id, "orga": strings.ToLower(orga)})
		return "", err
	})

	return err

}
