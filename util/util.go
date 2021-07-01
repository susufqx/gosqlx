package util

import "strings"

func CamelCaseToSnackCase(a string) string {
	s := []rune{}
	for _, v := range a {
		if v >= 65 && v <= 90 {
			s = append(s, 95, v)
		} else {
			s = append(s, v)
		}
	}

	if s[0] == 95 {
		s = s[1:]
	}

	return strings.ToLower(string(s))
}
