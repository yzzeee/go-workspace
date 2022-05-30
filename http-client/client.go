package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go-practice/http-client/prometheus"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	queryParams := url.Values{}
	//queryParams.Add("node", "aaa|bbb|ccc")

	// 대시보드 래디얼 차트에서 필요로하는 쿼리 메트릭 키 목록
	var dashboardRadialChartMetricKeys = []string{"radial_chart_node_cpu_usage", "node_total_cpu_core_count"}

	// 반환값
	var dashboardRadialChartMetricResponses = make(map[string]*prometheus.NewMetricResponse)

	//use http
	// 클라이언트 생성(TLS insecure 옵션, 인증 정보, 헤더 설정)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	for _, metricKey := range dashboardRadialChartMetricKeys {
		// 클라이언트에서 요청한 key 에 따른 쿼리 생성
		metricDefinition, isMetric := prometheus.MetricDefinitions[metricKey]
		if isMetric {
			metricQueryTemplates := metricDefinition.QueryTemplates
			metricFilters := metricDefinition.MetricFilter
			queries := make([]string, len(metricQueryTemplates))
			for i, queryTemplate := range metricQueryTemplates {
				queryGenerator := metricDefinition.QueryGenerators[i]
				if queryGenerator != nil {
					queries[i] = queryGenerator(queryTemplate, queryParams)
				} else {
					queries[i] = queryTemplate
				}
			}

			// 프로메테우스 모니터링 API 호출
			responses := make([]interface{}, len(queries))
			for i, query := range queries {
				requestURL := "https://xxx.xxx.xxx.xxx/api/v1/query?query=" + url.QueryEscape(query)

				request, err := http.NewRequest("GET", requestURL, nil)
				if err != nil {
					fmt.Println("failed to create http request, err=%s", err)
				}
				request.Header.Add("Content-Type", "application/json; charset=UTF-8")
				request.Header.Add("Access-Control-Allow-Origin", "*")
				request.Header.Add("Access-Control-Allow-Methods", "*")
				request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ""))

				// 응답 요청
				response, err := client.Do(request)
				if err != nil {
					fmt.Println("failed to call http request, err=%s", err)
				}

				// 응답값 파싱
				// TODO 예외 처리
				responseBytes, _ := ioutil.ReadAll(response.Body)
				a := make(map[string]map[string][]map[string][]interface{})
				_ = json.Unmarshal(responseBytes, &a)
				for _, ele := range a["data"]["result"] {
					responses[i] = ele["value"][1]
				}
			}

			res := &prometheus.NewMetricResponse{}
			if metricFilters != nil {
				fmt.Println("111")
				res = metricFilters(responses...)
			} else {
				fmt.Println("222")
				res.Usage = fmt.Sprintf("%s", responses)
			}
			res.Title = metricDefinition.Title
			res.Unit = string(metricDefinition.Unit)
			dashboardRadialChartMetricResponses[metricKey] = res
		}
	}

	client.CloseIdleConnections()

	fmt.Println(json.Marshal(dashboardRadialChartMetricResponses))
	fmt.Printf("%v", dashboardRadialChartMetricResponses)
}
