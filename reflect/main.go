package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func indexOf(arr []string, val string) int {
	for pos, v := range arr {
		if v == val {
			return pos
		}
	}
	return -1
}

func main() {
	var testMap = map[string]string{
		"2.29.0": "one",
		"2.0.0":  "two",
		"5.0.0":  "three",
	}

	// 1번 방법이 더 빠르다.
	startTime1 := time.Now()
	// 입력한 키와 동일한 키가 존재하는지 여부 확인
	keys1 := make([]string, 0, len(testMap))
	for key := range testMap {
		keys1 = append(keys1, key)
	}
	elapseTime1 := time.Since(startTime1)
	fmt.Println(elapseTime1, keys1)

	startTime2 := time.Now()
	keys2 := reflect.ValueOf(testMap).MapKeys()
	strkeys := make([]string, len(keys2))
	for i := 0; i < len(keys2); i++ {
		strkeys[i] = keys2[i].String()
	}
	elapseTime2 := time.Since(startTime2)
	fmt.Println(elapseTime2, keys2)

	// -------------------------------------------------------------

	// 제공된 버전 정보(순차적으로 정렬되어 있음)
	var orderedVersionList = []string{"1.0.0", "2.0.3", "5.0.0"}

	// 사용하고자 하는 버전 정보
	//var inputVersion = "2.0.3"
	var inputVersion = "2.0.3-rc1"

	// 정의된 버전이 존재하는지 확인
	index := indexOf(keys1, inputVersion)
	if index == -1 {
		// 정의된 버전이 없는 경우 제공된 버전 목록에서 하위 버전 중 가장 가까운 버전을 찾음
		iMajor, iMinor, iPatch := parseVersion(inputVersion, orderedVersionList)

		//var usedVersion string
		//orderedVersionObj := map[string][]string{}
		//for _, version := range orderedVersionList {
		//	// 메이저 버전 확인
		//	ma, mi, pa := parseVersion(version, orderedVersionList)
		//
		//	// 마이너 버전 확인
		//	fmt.Println(ma, mi, pa, usedVersion)
		//}
		fmt.Println(iMajor, iMinor, iPatch)
	} else {
		// 버전이 존재하는 경우
		fmt.Println(index, testMap[keys1[index]])
	}
}

func parseVersion(inputVersion string, keys1 []string) (int, int, int) {
	// - 이후의 정보는 제거
	i := strings.Index(inputVersion, "-")

	if i > -1 {
		inputVersion = inputVersion[:i]
	}

	// 사용자 입력 버전 정보 처리(숫자 및 . 을 제외한 문자를 제거)
	regex := regexp.MustCompile(`[^0-9.]`)
	inputVersion = regex.ReplaceAllString(inputVersion, "")

	var major int
	var minor int
	var patch int

	inputVersions := strings.Split(inputVersion, ".")
	size := len(inputVersions)

	if size != 0 {
		val, _ := strconv.Atoi(inputVersions[0])
		major = val
	}
	if size > 1 {
		val, _ := strconv.Atoi(inputVersions[1])
		minor = val
	}
	if size > 2 {
		val, _ := strconv.Atoi(inputVersions[2])
		patch = val
	}
	return major, minor, patch
}
