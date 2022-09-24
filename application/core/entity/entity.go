package entity

import (
	"fmt"
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
	StartDate    time.Time    `json:"startDate"`
	EndDate      time.Time    `json:"endDate"`
	Size         int64        `json:"size"`
	MatchingType MatchingType `json:"matchingType"`
	Active       bool         `json:"active"`
	Frequency    Frequency    `json:"frequency"`
	Week         Week         `json:"week"`
	Days         []string     `json:"days"` // format day1_day2_day3.....
	LastRun      time.Time    `json:"-"`
	NextRun      time.Time    `json:"-"`
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
	   1- find out the index of the day of lastRun
	   2- update the date accoording to the next run day
	         - if frequency == daily then next run day is just +1 from last run day

	*/

	//
	postDay := 0

	for i, d := range s.Days {
		if d == s.LastRun.Weekday().String() {
			postDay = i
			break
		}
	}
	newDay := postDay
	if postDay == len(s.Days)-1 {
		newDay = 0
	}
	switch s.Frequency {
	case Daily:
		s.NextRun = s.LastRun.AddDate(0, 0, 1)
	case Weekly:
		fmt.Println(s.LastRun.Weekday().String())
		s.NextRun = s.LastRun.AddDate(0, 0, helper.HMDays(s.LastRun.Weekday().String(), s.Days[newDay]))
	case Monthly:
		newMonth := 0
		if postDay == len(s.Days)-1 {
			newMonth = 1
		}
		s.NextRun = s.LastRun.AddDate(0, newMonth, helper.HMDays(s.LastRun.Weekday().String(), s.Days[newDay]))
	}
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
