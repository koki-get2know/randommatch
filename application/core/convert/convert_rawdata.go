package convert

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type User struct {
	UserId                 string
	Name                   string
	Email                  string
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

/* Generate random strings

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890$#!@")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}*/

func ConvertRawDataToJson(filename string) []byte {

	csvFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	// Skip first row data
	row1, err := bufio.NewReader(csvFile).ReadSlice('\n')
	if err != nil {
		fmt.Println(err)
	}
	_, err = csvFile.Seek(int64(len(row1)), io.SeekStart)
	if err != nil {
		fmt.Println(err)
	}

	// Read data
	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var user User
	var users []User

	// Create a json data
	for _, each := range csvData {
		user.UserId = "" //randStringRunes(32)
		user.Name = each[0]
		user.Email = each[1]
		user.Groups = strings.Split(each[2], "-") //[]string{each[4]}
		user.Genre = each[3]
		user.Birthday = each[4]
		user.Hobbies = strings.Split(each[5], "-")             //[]string{each[7]}
		user.MatchPreference = strings.Split(each[6], "-")     //[]string{each[8]}
		user.MatchPreferenceTime = strings.Split(each[7], "-") //[]string{each[9]}
		user.PositionHeld = each[8]
		user.MultiMatch, _ = strconv.ParseBool(each[9])
		user.PhoneNumber = each[10]
		user.Departement = each[11]
		user.Location = each[12]
		user.Seniority = each[13]
		user.Role = each[14]

		user.NumberOfMatching = 0
		user.NumberMatchingAccepted = 0
		user.NumberMatchingDeclined = 0
		user.AverageMatchingRate = 0

		users = append(users, user)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(users)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return jsonData
}

func GenerateJsonFile(filename string) {

	jsonData := ConvertRawDataToJson(filename)

	jsonFile, err := os.Create("./data.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()
}
