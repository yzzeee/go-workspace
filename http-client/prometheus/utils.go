package prometheus

import (
	"fmt"
	"go-practice/common"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type version struct {
	version string
}

func GetTargetPrometheusVersion(definedVersions []string, clusterPrometheusVersion string) string {
	var targetVersion string

	definedVersions = append(definedVersions, clusterPrometheusVersion)
	definedVersions = sortVersions(definedVersions)
	inputVersionIdx := common.IndexOf(sortVersions(definedVersions), clusterPrometheusVersion)

	if inputVersionIdx == 0 {
		targetVersion = definedVersions[inputVersionIdx+1]
	} else {
		targetVersion = definedVersions[inputVersionIdx-1]
	}

	return targetVersion
}

func ParseVersion(version string) string {
	// - 이후의 정보는 제거
	i := strings.Index(version, "-")

	if i > -1 {
		version = version[:i]
	}

	// 사용자 입력 버전 정보 처리(숫자 및 . 을 제외한 문자를 제거)
	regex := regexp.MustCompile(`[^0-9.]`)
	return regex.ReplaceAllString(version, "")
}

func sortVersions(versions []string) []string {
	versionSlice := make([]version, 0)

	for _, v := range versions {
		versionSlice = append(versionSlice, version{v})
	}

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

	return versionKeys
}

func pad(max float64, version []string) int {
	x := make([]string, int(max))
	for i := range x {
		x[i] = "00000"
	}
	for j := range version {
		x[j] = fmt.Sprintf("%05s", version[j])
	}
	n, _ := strconv.Atoi(strings.Join(x, ""))
	return n
}

// queryTemplateParser 쿼리 템플릿과 쿼리 파라미터를 인자로 받아서 쿼리를 생성하는 클로저를 반환하는 함수
func queryTemplateParser(paramKeys []interface{}) func(string, map[string]interface{}) (string, string) {
	params := make([]interface{}, len(paramKeys))
	return func(queryTemplate string, bodyParams map[string]interface{}) (string, string) {
		var rangeParams string
		for i, paramKey := range paramKeys {
			param := bodyParams[paramKey.(string)]
			if param == nil {
				param = ".*"
			}
			if param == "" {
				param = ""
			}
			params[i] = param
		}
		if bodyParams["start"] != nil && bodyParams["end"] != nil && bodyParams["step"] != nil {
			start := bodyParams["start"]
			end := bodyParams["end"]
			step := bodyParams["step"]
			rangeParams = fmt.Sprintf("&start=%s&end=%s&step=%s", start, end, step)
		}
		return fmt.Sprintf(queryTemplate, params...), rangeParams
	}
}
