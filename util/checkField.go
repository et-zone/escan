package util

import (
	"strings"
)

func IsNil(param interface{}) bool {
	return param == nil
}

func IsEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}
