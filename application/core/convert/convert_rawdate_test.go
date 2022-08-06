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
			"id": "",
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
			"role": "super-user"
		},
		{
			"id": "",
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
			"role": "user"
		}
	]`
	require.JSONEq(t, expected, actual)
}
