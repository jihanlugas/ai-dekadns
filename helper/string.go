package helper

import (
	"errors"
	"strings"
)

func AllowedKey(data map[string]interface{}, allowedList []string) error {
	for field := range data {
		allow := false
		for _, allowField := range allowedList {
			if field == allowField {
				allow = true
				break
			}
		}
		if !allow {
			return errors.New(field)
		}
	}
	return nil
}

func GetTextBetween(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}
