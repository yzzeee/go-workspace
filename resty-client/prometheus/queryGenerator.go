package prometheus

import (
	"fmt"
	"net/url"
)

func setQueryParams(test []interface{}) func(string, url.Values) string {
	vv := make([]interface{}, len(test))
	return func(queryTemplate string, queryParams url.Values) string {
		for i, t := range test {
			v := queryParams.Get(t.(string))
			if v == "" {
				v = ".*"
			}
			vv[i] = v
		}
		return fmt.Sprintf(queryTemplate, vv...)
	}
}
