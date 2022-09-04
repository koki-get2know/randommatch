package entity

import (
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
