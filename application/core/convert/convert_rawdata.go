package convert

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
)

type User struct {
	UserId                 string   `json:"userId"`
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
	NumberOfMatching       int      `json:"numberOfMatching"`
	NumberMatchingAccepted int      `json:"numberMatchingAccepted"`
	NumberMatchingDeclined int      `json:"numberMatchingDeclined"`
	AverageMatchingRate    int      `json:"averageMatchingRate"`
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

func csvReaderToUsers(r io.Reader) ([]User, error) {
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var header []string

	if len(records) > 0 {
		// skip the header
		header = records[0]
		records = records[1:]
	}
	if len(header) < 15 {
		return nil, fmt.Errorf("header and content not matching")
	}
	var users []User
	for _, record := range records {
		user := User{
			UserId:              "", //randStringRunes(32)
			Name:                record[0],
			Email:               record[1],
			Groups:              strings.Split(record[2], "-"),
			Gender:              record[3],
			Birthday:            record[4],
			Hobbies:             strings.Split(record[5], "-"),
			MatchPreference:     strings.Split(record[6], "-"),
			MatchPreferenceTime: strings.Split(record[7], "-"),
			PositionHeld:        record[8],
			PhoneNumber:         record[10],
			Department:          record[11],
			Location:            record[12],
			Seniority:           record[13],
			Role:                record[14],
		}

		user.MultiMatch, err = strconv.ParseBool(record[9])
		if err != nil {
			fmt.Printf("Warning wrong boolean string value passed for user %v, value passed: %v\n", user.Name, user.MultiMatch)
			user.MultiMatch = false
		}

		users = append(users, user)
	}
	return users, nil
}

func CsvToUsers(csvFile *multipart.FileHeader) ([]User, error) {
	openedFile, err := csvFile.Open()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return csvReaderToUsers(openedFile)
}

func ConvertRawDataToJson(filepath string) []byte {

	csvFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	defer csvFile.Close()
	// Read data
	users, err := csvReaderToUsers(csvFile)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	// Convert to JSON
	jsonData, err := json.Marshal(users)

	if err != nil {
		fmt.Println(err)
		return []byte{}
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
}
