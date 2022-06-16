package prometheus

import (
	"fmt"
	"net/url"
)

func queryGenerator(paramKeys []interface{}, isRange bool) func(string, url.Values) (string, string) {
	params := make([]interface{}, len(paramKeys))
	return func(queryTemplate string, queryParams url.Values) (string, string) {
		var rangeParams string
		for i, paramKey := range paramKeys {
			param := queryParams.Get(paramKey.(string))
			if param == "" {
				param = ".*"
			}
			params[i] = param
		}

		if isRange {
			start := queryParams.Get("start")
			end := queryParams.Get("end")
			step := queryParams.Get("step")
			rangeParams = fmt.Sprintf("&start=%s&end=%s&step=%s", start, end, step)
		}
		return fmt.Sprintf(queryTemplate, params...), rangeParams
	}
}
