package prometheus

import "fmt"

func MetricFilter(args ...interface{}) *NewMetricResponse {
	return &NewMetricResponse{
		"",
		fmt.Sprintf("%v", args[0]),
		fmt.Sprintf("%v", args[1]),
		fmt.Sprintf("%v", args[2]),
		"",
	}
}
