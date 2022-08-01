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
			"userId": "",
			"name": "Pins Prestilien",
			"email": "pinsdev24@gmail.com",
			"groups": [
				""
			],
			"gender": "Male",
			"birthday": "10/01",
			"hobbies": [
				"Data science",
				"Space",
				"Télévison"
			],
			"matchPreference": [
				"girls"
			],
			"matchPreferenceTime": [
				"14:00PM"
			],
			"positionHeld": "CEO",
			"multiMatch": false,
			"phoneNumber": "699999999",
			"department": "Informatique",
			"location": "Kao",
			"seniority": "",
			"role": "super-user",
			"numberOfMatching": 0,
			"numberMatchingAccepted": 0,
			"numberMatchingDeclined": 0,
			"averageMatchingRate": 0
		},
		{
			"userId": "",
			"name": "Pins Prestilien",
			"email": "pinsdev24@gmail.com",
			"groups": [
				"DS",
				"IA",
				"SPACE"
			],
			"gender": "Male",
			"birthday": "10/01",
			"hobbies": [
				"Jeux vidéos",
				"Musique"
			],
			"matchPreference": [
				"same groups"
			],
			"matchPreferenceTime": [
				"14:00PM"
			],
			"positionHeld": "Admin",
			"multiMatch": false,
			"phoneNumber": "699999999",
			"department": "Math",
			"location": "Fpol",
			"seniority": "",
			"role": "user",
			"numberOfMatching": 0,
			"numberMatchingAccepted": 0,
			"numberMatchingDeclined": 0,
			"averageMatchingRate": 0
		}
	]`
	require.JSONEq(t, expected, actual)
}
