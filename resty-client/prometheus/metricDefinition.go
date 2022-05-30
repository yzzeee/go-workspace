package prometheus

import (
	"net/url"
)

var (
	MetricDefinition = map[string]Metric{
		"node_cpu_usage": {
			QueryTemplates: []string{
				"sum(rate(node_cpu_seconds_total{mode!=\"idle\",mode!=\"iowait\",instance=~\"%s\"}[3m]))",
				"sum(kube_node_status_capacity{resource=\"cpu\",unit=\"core\",node=~\"%s\"})",
				"round(sum(rate(node_cpu_seconds_total{mode!=\"idle\",mode!=\"iowait\",instance=~\"%s\"}[3m]))/sum(kube_node_status_capacity{resource=\"cpu\",unit=\"core\",node=~\"%s\"}) * 100, 0.1)",
			},
			QueryGenerators: QueryGenerators{
				setQueryParams([]interface{}{"instance"}),
				setQueryParams([]interface{}{"node"}),
				setQueryParams([]interface{}{"instance", "node"})},
			MetricFilter: MetricFilter,
			Title:        "CPU",
			Unit:         MilliCore,
		},
	}
)

type QueryGenerators []func(queryTemplate string, queryParams url.Values) string

type Metric struct {
	QueryTemplates  []string
	QueryGenerators QueryGenerators
	MetricFilter    func(args ...interface{}) *NewMetricResponse
	Title           string
	Unit            Unit
}

type NewMetricResponse struct {
	Title      string `json:"title"`
	Usage      string `json:"usage"`
	Total      string `json:"total,omitempty"`
	Percentage string `json:"percentage,omitempty"`
	Unit       string `json:"unit"`
}

// Unit 단위 종류
type Unit string

const (
	Core      = Unit("core")
	MilliCore = Unit("millicore")
	Byte      = Unit("byte")
)
