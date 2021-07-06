package util

import (
	"reflect"
	"strings"
)

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

func MapJoin(m ...map[string]interface{}) map[string]interface{} {
	if len(m) == 0 {
		return nil
	}

	res := map[string]interface{}{}
	for _, mi := range m {
		for k, v := range mi {
			res[k] = v
		}
	}

	return res
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
