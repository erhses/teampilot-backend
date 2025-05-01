package utils

import (
	"reflect"
	"strconv"
)

func Str2Int(value string) int {
	returnVal, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return returnVal
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
