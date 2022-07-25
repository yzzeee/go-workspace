package prometheus

import (
	"fmt"
)

// queryGenerator 쿼리 템플릿과 쿼리 파라미터를 인자로 받아서 쿼리를 생성하는 클로저를 반환하는 함수
func queryGenerator(paramKeys []interface{}) func(string, map[string]interface{}) (string, string) {
	params := make([]interface{}, len(paramKeys))
	return func(queryTemplate string, bodyParams map[string]interface{}) (string, string) {
		var rangeParams string
		for i, paramKey := range paramKeys {
			param := bodyParams[paramKey.(string)]
			if param == "" || params == nil {
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
