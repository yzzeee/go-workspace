package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"go-practice/resty-client/prometheus"
	"net/url"
)

type prometheusResponse struct {
	Status    string          `json:"status"`
	Data      json.RawMessage `json:"data"`
	ErrorType v1.ErrorType    `json:"errorType,omitempty"`
	Error     string          `json:"error,omitempty"`
	Warnings  []string        `json:"warnings,omitempty"`
}

type prometheusResponses []prometheusResponse

// https://pkg.go.dev/github.com/go-resty/resty#section-readme
func main() {
	queryParams := url.Values{}
	queryParams.Add("nodes", "aaa|bbb")

	// 클라이언트에서 요청한 key 에 따른 쿼리 생성
	metricDefinition, isMetric := prometheus.MetricDefinition["cpu_usage_node"]
	if isMetric {
		metricQueryTemplates := metricDefinition.QueryTemplates
		queries := make([]string, len(metricQueryTemplates))
		for i, queryTemplate := range metricQueryTemplates {
			queryGenerator := metricDefinition.QueryGenerators[i]
			if queryGenerator != nil {
				queries[i] = queryGenerator(queryTemplate, queryParams)
			} else {
				queries[i] = queryTemplate
			}
		}

		// use resty
		// 클라이언트 생성(TLS insecure 옵션, 인증 정보, 헤더 설정)
		client := resty.New()
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		client.SetHeader("Content-Type", "application/json; charset=UTF-8").
			SetHeader("Access-Control-Allow-Origin", "*").
			SetHeader("Access-Control-Allow-Methods", "*").
			SetAuthToken("")

		// 프로메테우스 모니터링 API 호출
		var prometheusResponses prometheusResponses
		tempPrometheusResponses := make([]interface{}, len(queries))
		for i, query := range queries {
			var prometheusResponse prometheusResponse
			response, err := client.R().
				SetQueryString(fmt.Sprintf("query=%s", query)).
				Get("https://xxx.xxx.xxx.xxx/api/v1/query")
			if err != nil {
				fmt.Println("failed to request prometheus server, err=%s", err)
			}
			err = json.Unmarshal(response.Body(), &prometheusResponse)
			if err != nil {
				fmt.Println("failed to parsing prometheus data, err=%s", err)
			}
			tempPrometheusResponses[i] = prometheusResponse
		}

		// TODO : 결과값 파싱
		marshalResponse, err := json.Marshal(tempPrometheusResponses)
		err = json.Unmarshal(marshalResponse, &prometheusResponses)
		if err != nil {
			fmt.Println("failed to parsing prometheus data, err=%s", err)
		}
		fmt.Println(fmt.Sprintf("%s", prometheusResponses), err)

		client.SetCloseConnection(true)
	}
}
