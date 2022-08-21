package helper

import (
	"log"
	"strings"
	"time"
)

func Track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func Duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}

func Contains(s []any, e string) bool {
	for _, a := range s {
			if a.(string) == e {
					return true
			}
	}
	return false
}

func ContainsString(s []string, e string) bool {
	for _, a := range s {
			if a == e {
					return true
			}
	}
	return false
}

func ItemsWithPrefixInRole(s []any, prefix string) []string {
	orgs := []string{}
	for _, a := range s {
			if strings.HasPrefix(a.(string), prefix) {
					orgs = append(orgs, strings.ToLower(strings.TrimPrefix(a.(string),prefix)) )
			}
	}
	return orgs
}