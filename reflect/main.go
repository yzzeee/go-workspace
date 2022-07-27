package main

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

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

	fmt.Println("================================================================================================")

	// 사용하고자 하는 버전 정보
	//var inputVersion = "2.0.4"
	//var inputVersion = "2.29.0"
	var inputVersion = "2.1.5-rc1"

	fmt.Println("입력한 버전: ", inputVersion)
	fmt.Println("================================================================================================")

	versionSlice := make([]version, 0)
	versionSlice = append(versionSlice, version{"1.1.2"})
	versionSlice = append(versionSlice, version{"1.1.5"})
	versionSlice = append(versionSlice, version{"2.0.4"})
	versionSlice = append(versionSlice, version{"2.0.9"})
	versionSlice = append(versionSlice, version{"5.0.3"})
	versionSlice = append(versionSlice, version{parseVersion(inputVersion)}) // 사용자가 입력한 버전

	sort.Slice(versionSlice, func(i, j int) bool {
		a := strings.Split(versionSlice[i].version, ".")
		b := strings.Split(versionSlice[j].version, ".")
		max := math.Max(float64(len(a)), float64(len(b)))
		return pad(max, a) < pad(max, b)
	})

	versionKeys := make([]string, 0, len(versionSlice))
	for key := range versionSlice {
		versionKeys = append(versionKeys, versionSlice[key].version)
	}

	inputVersionIdx := indexOf(versionKeys, parseVersion(inputVersion))
	var usedVersion string
	if inputVersionIdx == 0 {
		usedVersion = versionKeys[inputVersionIdx+1]
	} else {
		usedVersion = versionKeys[inputVersionIdx-1]
	}
	fmt.Println("선택된 버전 1: ", usedVersion)

	fmt.Println("================================================================================================")

	// 제공된 버전 정보(순차적으로 정렬되어 있음)
	var orderedVersionList = []string{"1.1.2", "1.1.5", "2.0.4", "2.0.9", "5.0.3"}

	// 정의된 버전이 존재하는지 확인
	index := indexOf(keys1, inputVersion)
	if index == -1 {
		// 정의된 버전이 없는 경우 제공된 버전 목록에서 하위 버전 중 가장 가까운 버전을 찾음
		foundVersion := parseVersion2(inputVersion, orderedVersionList)

		// 선택된 버전 확인
		fmt.Println("선택된 버전 2: ", foundVersion)
	} else {
		// 버전이 존재하는 경우
		fmt.Println(index, testMap[keys1[index]])
		fmt.Println(inputVersion)
	}
}

func indexOf(arr []string, val string) int {
	for pos, v := range arr {
		if v == val {
			return pos
		}
	}
	return -1
}

type version struct {
	version string
}

func pad(max float64, version []string) int {
	x := make([]string, int(max))
	for i := range x {
		x[i] = "000"
	}
	for j := range version {
		x[j] = fmt.Sprintf("%03s", version[j])
	}
	n, _ := strconv.Atoi(strings.Join(x, ""))
	return n
}

func parseVersion(version string) string {
	// - 이후의 정보는 제거
	i := strings.Index(version, "-")

	if i > -1 {
		version = version[:i]
	}

	// 사용자 입력 버전 정보 처리(숫자 및 . 을 제외한 문자를 제거)
	regex := regexp.MustCompile(`[^0-9.]`)
	return regex.ReplaceAllString(version, "")
}

func splitedVersion(inputVersion string) (int, int, int) {
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

func parseVersion2(inputVersion string, orderedVersionList []string) string {
	var orderedSplitedVersionList [][3]int

	var maxMajor int
	var maxMinor int
	var maxPatch int

	var selectedMajor = -1
	var selectedMinor = -1
	var selectedPatch = -1

	inputMajor, inputMinor, inputPatch := splitedVersion(inputVersion)

	var diffMajor = inputMajor
	var diffMinor = inputMinor
	var diffPatch = inputPatch

	var sameMajorExist = false
	var sameMajorAndMinorExist = false

	for _, version := range orderedVersionList {
		major, minor, patch := splitedVersion(version)
		orderedSplitedVersionList = append(orderedSplitedVersionList, [3]int{major, minor, patch})
	}

	// major 찾기
	// 같은 메이저를 찾거나, 하위 메이저중 가까운 메이저를 선택한다.
	for _, version := range orderedSplitedVersionList {
		major := version[0]

		if major == inputMajor {
			sameMajorExist = true
			selectedMajor = major
			break
		}

		if major < inputMajor {
			diff := inputMajor - major
			if diff < diffMajor {
				diffMajor = diff
				selectedMajor = major
			}
		}

		if major > maxMajor {
			maxMajor = major
		}
	}
	// 하위 메이저가 존재 하지 않는 경우, 상위 중에 가까운 메이저를 사용한다.
	if selectedMajor == -1 {
		diffMajor = maxMajor

		for _, version := range orderedSplitedVersionList {
			major := version[0]

			if major > inputMajor {
				diff := major - inputMajor
				if diff < diffMajor {
					diffMajor = diff
					selectedMajor = major

				}
			}
		}
	}

	// minor 찾기
	// 같은 메이저와 같은 마이너를 찾으면 바로 루프를 빠져나오고, 그 이외의 경우에는 선택한 메이저로 시작하며 하위 마이너중 가까운 마이너를 선택한다.
	for _, version := range orderedSplitedVersionList {
		major := version[0]
		minor := version[1]

		if major != selectedMajor {
			continue
		}

		if sameMajorExist && minor == inputMinor {
			sameMajorAndMinorExist = true
			selectedMinor = minor
			break
		}

		if minor <= inputMinor {
			diff := inputMinor - minor
			if diff < diffMinor {
				diffMinor = diff
				selectedMinor = minor
			}
		}

		if minor > maxMinor {
			maxMinor = minor
		}
	}
	// 하위 마이너가 존재 하지 않는 경우, 가장 가까운 마이너를 사용한다.
	if selectedMinor == -1 {
		diffMinor = maxMinor

		for _, version := range orderedSplitedVersionList {
			major := version[0]
			minor := version[1]

			if major != selectedMajor {
				continue
			}

			diff := minor - inputMinor
			if diff <= diffMinor {
				diffMinor = diff
				selectedMinor = minor
			}
		}
	}

	// patch 찾기
	// 메이저와 마이너가 같은 경우에는 하위 중에 가까운 패치버전을, 그 이외의 경우에는 제일 큰 패치 버전을 사용한다.
	for _, version := range orderedSplitedVersionList {
		major := version[0]
		minor := version[1]
		patch := version[2]

		if major != selectedMajor || minor != selectedMinor {
			continue
		}

		if sameMajorAndMinorExist {
			if patch == inputPatch {
				selectedPatch = patch
				break
			}

			diff := inputPatch - patch
			if patch < inputPatch {
				if diff < diffPatch {
					diffPatch = diff
					selectedPatch = patch
				}
			}
		}

		if patch > maxPatch {
			maxPatch = patch
		}
	}
	if selectedPatch == -1 {
		selectedPatch = maxPatch
	}

	return fmt.Sprintf("%d.%d.%d", selectedMajor, selectedMinor, selectedPatch)
}
