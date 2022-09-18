package entity

import (
	"strings"
	"time"

	"github.com/koki/randommatch/utils/helper"
)

type Tag struct {
	Name string `json:"name"`
}

type Organization struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MatchingCycle struct {
	Id string
}

type Frequency string

const (
	Monthly         Frequency = "monthly"
	Weekly          Frequency = "weekly"
	Daily           Frequency = "daily"
	Hourly          Frequency = "hourly"
	Every_half_hour Frequency = "every_half_hour"
	Every_ten_min   Frequency = "every_ten_min"
)

type Week string

const (
	First  Frequency = "first"
	Second Frequency = "second"
	Third  Frequency = "third"
	Fourth Frequency = "fourth"
	Last   Frequency = "last"
)

type Selector uint8

const (
	Basic Selector = iota
	Group
)

type MatchingType string

const (
	Simple MatchingType = "simple"
	Groups MatchingType = "group"
	Tags   MatchingType = "tag"
)

type Schedule struct {
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	CreateDate   time.Time    `json:"-"`
	Time         string       `json:"time"`
	Duration     string       `json:"duration"`
	StartDate    string       `json:"startDate"`
	EndDate      string       `json:"endDate"`
	Size         int64        `json:"size"`
	MatchingType MatchingType `json:"matchingType"`
	Active       bool         `json:"active"`
	Frequency    Frequency    `json:"frequency"`
	Week         Week         `json:"week"`
	Days         string       `json:"days"` // format day1_day2_day3.....
	LastRun      string       `json:"-"`
	NextRun      string       `json:"-"`
}
type MatchingStat struct {
	Id               string    `json:"id"`
	NumGroups        int       `json:"numGroups"`
	NumConversations int       `json:"numConversations"`
	NumFailed        int       `json:"numFailed"`
	CreatedAt        time.Time `json:"createdAt"`
}

type User struct {
	Id                     string   `json:"id"`
	Name                   string   `json:"name"`
	Email                  string   `json:"email"`
	Tags                   []string `json:"tags"`
	Gender                 string   `json:"gender"`
	Birthday               string   `json:"birthday"`
	Hobbies                []string `json:"hobbies"`
	MatchPreference        []string `json:"matchPreference"`
	MatchPreferenceTime    []string `json:"matchPreferenceTime"`
	PositionHeld           string   `json:"positionHeld"`
	MultiMatch             bool     `json:"multiMatch"`
	PhoneNumber            string   `json:"phoneNumber"`
	Department             string   `json:"department"`
	Location               string   `json:"location"`
	Seniority              string   `json:"seniority"`
	Role                   string   `json:"role"`
	NumberOfMatching       int      `json:"-"`
	NumberMatchingAccepted int      `json:"-"`
	NumberMatchingDeclined int      `json:"-"`
	AverageMatchingRate    int      `json:"-"`
	//SubjectOfInterest    []string
}

func (s *Schedule) UdapteNextRun() {

	/*
	 - LastRun has a  following format:     Month separator week separator day
	   1- split LastRun according to the separator
	   2 - update each field
	   3 - assign  the value to nextRun
	*/
	sep := "_"
	months := []string{"January", "February", "March", "April", "May", "June", "july", "August", "September", "October", "November", "December"}
	week := []string{"First", "Second", "Third", "Fourth", "Last"}
	day := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	//
	switch s.Frequency {
	/*
		 -  frequently= monthly
		    step 2- update each field
			    - find out the day in Days schedule attribute abd month in months of a year
				- according to the next day update month
				- week doesn't change
	*/
	case Monthly:
		lastRun := strings.Split(s.LastRun, sep)
		days := strings.Split(s.Days, sep)
		posDay := 0
		for i, d := range days {
			if d == lastRun[2] {
				posDay = i
				break
			}
		}
		posMonth := 0
		for i, m := range months {
			if m == lastRun[0] {
				posMonth = i
				break
			}
		}
		nextMonth := lastRun[0]

		if posDay == len(days)-1 {
			if posMonth == len(months)-1 {
				nextMonth = months[0]
			} else {
				nextMonth = months[posMonth+1]
			}
			s.NextRun = nextMonth + sep + lastRun[1] + sep + days[0]
		} else {
			s.NextRun = nextMonth + sep + lastRun[1] + sep + days[posDay+1]
		}
	case Weekly:
		/*
			      -frenquently weekly
				     step 2:
					 -find out the day in Days schedule attribute, week and month in week of month and months of year respectively
					 - according to the next day update week and month
		*/
		lastRun := strings.Split(s.LastRun, sep)
		days := strings.Split(s.Days, sep)
		posDay := 0
		for i, d := range days {
			if d == lastRun[2] {
				posDay = i
				break
			}
		}
		posWeek := 0
		for i, m := range months {
			if m == lastRun[1] {
				posWeek = i
				break
			}
		}

		posMonth := 0
		for i, m := range months {
			if m == lastRun[0] {
				posMonth = i
				break
			}
		}
		nextWeek := lastRun[1]
		nextMonth := lastRun[0]
		if posDay == len(days)-1 {
			if posWeek == len(week)-1 {
				nextWeek = week[0]
				if posMonth == len(months)-1 {
					nextMonth = months[0]
				} else {
					nextMonth = months[posMonth+1]
				}

			} else {
				nextWeek = week[posWeek+1]
			}
			s.NextRun = nextMonth + sep + nextWeek + sep + days[0]
		} else {
			s.NextRun = nextWeek + sep + nextWeek + sep + days[posDay+1]
		}
	case Daily:
		/*
			-frenquently Daily
					     step 2:
						 -Find out the last day , week and month in week of month and months of year respectively
						 -According to the next day update week and month
		*/
		lastRun := strings.Split(s.LastRun, sep)

		posDay := 0
		for i, d := range day {
			if d == lastRun[2] {
				posDay = i
				break
			}
		}
		posWeek := 0
		for i, m := range months {
			if m == lastRun[1] {
				posWeek = i
				break
			}
		}

		posMonth := 0
		for i, m := range months {
			if m == lastRun[0] {
				posMonth = i
				break
			}
		}

		nextWeek := lastRun[1]
		nextMonth := lastRun[0]
		if posDay == len(day)-1 {
			if posWeek == len(week)-1 {
				nextWeek = week[0]
				if posMonth == len(months)-1 {
					nextMonth = months[0]
				} else {
					nextMonth = months[posMonth+1]
				}

			} else {
				nextWeek = week[posWeek+1]
			}
			s.NextRun = nextMonth + sep + nextWeek + sep + day[0]
		} else {
			s.NextRun = nextWeek + sep + nextWeek + sep + day[posDay+1]
		}
	}
}
func (s *Schedule) UdapteLastRun(date time.Time, sep string) {
	weekNumber := int(date.Day()) / 7
	week := ""
	switch weekNumber {
	case 0:
		week = "First"
	case 1:
		week = "Second"
	case 3:
		week = "Third"
	case 4:
		week = "Fourth"
	default:
		week = "Last"
	}
	s.LastRun = date.Month().String() + sep + week + sep + date.Weekday().String()
}
func (u *User) String() string {
	return u.Id
}

func (u *User) UserIn(users []User) (bool, int) {
	index := -1
	find := false
	for i, user := range users {
		if user.Id == u.Id {
			find = true
			index = i
			break
		}
	}

	return find, index
}

func (u *User) RmUser(users []User) []User {
	if f, i := u.UserIn(users); f {
		users = helper.RemoveByIndex(users, i)

	}
	return users
}
