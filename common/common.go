package common

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Get 두 개의 인자를 받아 두 번째 인자의 경로에 해당하는 값을 반환
func Get(args ...interface{}) interface{} {
	if len(args) < 2 {
		panic("Invalid number of argument. At least 2 arguments are required")
	}
	var (
		path     string
		mapData  interface{}
		fallback interface{}
	)

	if len(args) >= 1 {
		mapData = args[0]
	}
	if len(args) >= 2 {
		path = args[1].(string)
	}
	if len(args) >= 3 {
		fallback = args[2]
	}
	if mapData == nil {
		return fallback
	}
	defer func() interface{} {
		if r := recover(); r != nil {
			return fallback
		}
		return fallback
	}()
	data := mapData
	paths := strings.Split(path, ".")
	for _, key := range paths {
		value := reflect.ValueOf(data)
		dataType := value.Type().Kind()
		if dataType == reflect.Map {
			data = data.(map[string]interface{})[key]
			continue
		}

		if dataType == reflect.Slice {
			indx, err := strconv.Atoi(key)
			if err != nil {
				return fallback
			}
			data = value.Index(indx).Interface()
			continue
		}

		if dataType == reflect.Struct {
			data = value.FieldByName(key).Interface()
			continue
		}

		if dataType == reflect.String {
			indx, err := strconv.Atoi(key)
			if err != nil {
				return fallback
			}
			data = string(data.(string)[indx])
			continue
		}
		return fallback
	}
	return data
}

func Exists(s []string, target string) bool {
	i := sort.SearchStrings(s, target)
	return i < len(s) && s[i] == target
}

func MergeJSONMaps(maps ...map[string]interface{}) (result map[string]interface{}) {
	result = make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
