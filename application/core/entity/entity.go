package entity

type Organization struct {
	Id          string
	Name        string
	Description string
}

type User struct {
	Id                     string   `json:"id"`
	Name                   string   `json:"name"`
	Email                  string   `json:"email"`
	Groups                 []string `json:"groups"`
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

func (n *User) String() string {
	return n.Id
}
