package convert

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type User struct {
	UserId                 string
	Name                   string
	Email                  string
	Password               string
	Groups                 []string
	Genre                  string
	Birthday               string
	Hobbies                []string
	MatchPreference        []string
	MatchPreferenceTime    []string
	PositionHeld           string
	MultiMatch             bool
	PhoneNumber            string
	Departement            string
	Location               string
	Seniority              string
	Role                   string
	NumberOfMatching       int
	NumberMatchingAccepted int
	NumberMatchingDeclined int
	AverageMatchingRate    int
	//SubjectOfInterest    []string
}

func ConvertRawDataToJson(filename string) {

	csvFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var user User
	var users []User

	for _, each := range csvData {
		user.UserId = each[0]
		user.Name = each[1]
		user.Email = each[2]
		user.Password = each[3]
		user.Groups = []string{each[4]}
		user.Genre = each[5]
		user.Birthday = each[6]
		user.Hobbies = []string{each[7]}
		user.MatchPreference = []string{each[8]}
		user.MatchPreferenceTime = []string{each[9]}
		user.PositionHeld = each[10]
		user.MultiMatch, _ = strconv.ParseBool(each[11])
		user.PhoneNumber = each[12]
		user.Departement = each[13]
		user.Location = each[14]
		user.Seniority = each[15]
		user.Role = each[16]

		user.NumberOfMatching, _ = strconv.Atoi(each[17])
		user.NumberMatchingAccepted, _ = strconv.Atoi(each[18])
		user.NumberMatchingDeclined, _ = strconv.Atoi(each[19])
		user.AverageMatchingRate, _ = strconv.Atoi(each[20])

		users = append(users, user)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(users)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))

	jsonFile, err := os.Create("./data.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()
}
