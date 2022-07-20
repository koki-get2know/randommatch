package convert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertRawDataToJson(t *testing.T) {
	jsonData := ConvertRawDataToJson("./data.csv")
	actual := string(jsonData)
	expected := `[
		{
			"UserId": "",
			"Name": "Pins Prestilien",
			"Email": "pinsdev24@gmail.com",
			"Groups": [
				""
			],
			"Genre": "Male",
			"Birthday": "10/01",
			"Hobbies": [
				"Data science",
				"Space",
				"Télévison"
			],
			"MatchPreference": [
				"girls"
			],
			"MatchPreferenceTime": [
				"14:00PM"
			],
			"PositionHeld": "CEO",
			"MultiMatch": false,
			"PhoneNumber": "699999999",
			"Departement": "Informatique",
			"Location": "Kao",
			"Seniority": "",
			"Role": "super-user",
			"NumberOfMatching": 0,
			"NumberMatchingAccepted": 0,
			"NumberMatchingDeclined": 0,
			"AverageMatchingRate": 0
		},
		{
			"UserId": "",
			"Name": "Pins Prestilien",
			"Email": "pinsdev24@gmail.com",
			"Groups": [
				"DS",
				"IA",
				"SPACE"
			],
			"Genre": "Male",
			"Birthday": "10/01",
			"Hobbies": [
				"Jeux vidéos",
				"Musique"
			],
			"MatchPreference": [
				"same groups"
			],
			"MatchPreferenceTime": [
				"14:00PM"
			],
			"PositionHeld": "Admin",
			"MultiMatch": false,
			"PhoneNumber": "699999999",
			"Departement": "Math",
			"Location": "Fpol",
			"Seniority": "",
			"Role": "user",
			"NumberOfMatching": 0,
			"NumberMatchingAccepted": 0,
			"NumberMatchingDeclined": 0,
			"AverageMatchingRate": 0
		}
	]`
	require.JSONEq(t, expected, actual)
}
