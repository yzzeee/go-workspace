package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go-practice/common"
	"go-practice/http-client/config"
	"go-practice/http-client/prometheus"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	queryAPIEndpoint      = "/api/v1/query"
	queryRangeAPIEndpoint = "/api/v1/query_range"
)

func init() {
	config.Init()
}

func main() {
	queryParams := url.Values{}

	//now := time.Now()
	//now := time.Date(2022, 6, 27, 10, 15, 30, 0, time.Local)
	//queryParams.Add("start", strconv.Itoa(int(now.Add(-time.Minute*5).Unix())))
	//queryParams.Add("end", strconv.Itoa(int(now.Unix())))
	//queryParams.Add("step", "20") // 초 단위 간격
	//queryParams.Add("node", "master3.ocp4.inno.com|worker2.ocp4.inno.com")
	//queryParams.Add("instance", "master3.ocp4.inno.com|worker2.ocp4.inno.com")
	//queryParams.Add("namespace", "admin-workspace")

	//queryParams.Add("start", "1440398866.521")
	//queryParams.Add("end", "1656398866.521")
	//queryParams.Add("step", "120") // 초 단위 간격
	// 타임 스탬프에 대한 값이 조회되지 않는 매트릭 구간도 있어서 항상 개수가 동일하지 않다.
	//queryParams.Add("start", "1656324056.684")
	//queryParams.Add("end", "1656324056.684")
	//queryParams.Add("step", "120") // 초 단위 간격

	//queryParams.Add("start", "1656298468.184")
	//queryParams.Add("end", "1656309268.184")
	//queryParams.Add("step", "600") // 초 단위 간격

	//start=1655801906.107&end=1655802206.107&step=10

	// -- single 모니터링 하려는 메트릭 키의 목록
	//var metricKeys = []string{"quota"}
	//var metricKeys = []string{"container_cpu"}
	var metricKeys = []string{"container_memory"}
	//var metricKeys = []string{"container_file_system"}
	//var metricKeys = []string{"container_network_in"}
	//var metricKeys = []string{"container_network_out"}
	//var metricKeys = []string{"node_info"}
	//var metricKeys = []string{"node_cpu"}
	//var metricKeys = []string{"node_cpu_top"}
	//var metricKeys = []string{"node_cpu_top5_projects"}
	//var metricKeys = []string{"node_cpu_top5_pods"}
	//var metricKeys = []string{"node_memory"}
	//var metricKeys = []string{"node_memory_top"}
	//var metricKeys = []string{"node_memory_top5_projects"}
	//var metricKeys = []string{"node_memory_top5_pods"}
	//var metricKeys = []string{"node_file_system"}
	//var metricKeys = []string{"node_file_system_top"}
	//var metricKeys = []string{"node_file_system_top5_projects"}
	//var metricKeys = []string{"node_file_system_top5_pods"}
	//var metricKeys = []string{"node_network_in"}
	//var metricKeys = []string{"node_network_in_top"}
	//var metricKeys = []string{"node_network_in_top5_projects"}
	//var metricKeys = []string{"node_network_in_top5_pods"}
	//var metricKeys = []string{"node_network_out"}
	//var metricKeys = []string{"node_network_out_top"}
	//var metricKeys = []string{"node_network_out_top5_projects"}
	//var metricKeys = []string{"node_network_out_top5_pods"}
	//var metricKeys = []string{"node_pod_count"}
	//var metricKeys = []string{"node_pod_count_top"}
	//var metricKeys = []string{"node_pod_count_top5_projects"}
	//var metricKeys = []string{"quota_cpu_request"}
	//var metricKeys = []string{"quota_cpu_limit"}
	//var metricKeys = []string{"quota_memory_request"}
	//var metricKeys = []string{"quota_memory_limit"}
	//var metricKeys = []string{"range_node_cpu_usage"}
	//var metricKeys = []string{"range_container_cpu_usage"}
	//var metricKeys = []string{"range_cpu_load_average"}
	//var metricKeys = []string{"range_memory_usage"}
	//var metricKeys = []string{"range_memory_swap"}
	//var metricKeys = []string{"range_network_io"}
	//var metricKeys = []string{"range_network_packet"}
	//var metricKeys = []string{"range_network_bandwidth"}
	//var metricKeys = []string{"range_network_packet_receive_transmit"}
	//var metricKeys = []string{"range_network_packet_receive_transmit_drop"}
	//var metricKeys = []string{"range_file_system"}
	//var metricKeys = []string{"range_disk_io"}
	//var metricKeys = []string{"node_info"}

	// --- mixed
	//var metricKeys = []string{"node_cpu_top", "node_info"}

	var result = make(map[string]interface{})

	// 반환값
	for _, metricKey := range metricKeys {
		// 클라이언트에서 요청한 key 에 따른 쿼리 생성
		metricDefinition, isMetric := prometheus.MetricDefinitions[prometheus.MetricKey(metricKey)]

		// 정의된 메트릭 여부 확인
		if isMetric {
			innerMetricKeys := metricDefinition.MetricKeys
			if innerMetricKeys != nil { // 다른 메트릭의 값을 활용하는 메트릭 처리
				innerResult := make(map[string]interface{})
				for _, innerMetricKey := range innerMetricKeys {
					innerResult = common.MergeJSONMaps(innerResult, getQueryResult(innerMetricKey, queryParams))
				}
				metricResponse := prometheus.MakeMetricResponse(prometheus.MetricKey(metricKey), nil, "", nil, innerResult)
				result = common.MergeJSONMaps(result, map[string]interface{}{metricKey: metricResponse.Values})
			} else {
				queryResult := getQueryResult(prometheus.MetricKey(metricKey), queryParams)
				result = common.MergeJSONMaps(result, queryResult)
			}
		}

		// 최종 결과 확인
		final, _ := json.Marshal(result)
		fmt.Println("[   FINAL   ]", string(final))
	}
}

func getQueryResult(metricKey prometheus.MetricKey, queryParams url.Values) map[string]interface{} {
	var result = make(map[string]interface{})
	//use http
	// 클라이언트 생성(TLS insecure 옵션, 인증 정보, 헤더 설정)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	// 클라이언트 연결 종료 함수 등록
	defer client.CloseIdleConnections()

	metricDefinition, _ := prometheus.MetricDefinitions[metricKey]
	subLabels := metricDefinition.SubLabels
	queryTemplates := metricDefinition.QueryTemplates
	unitTypeKeys := metricDefinition.UnitTypeKeys
	primaryUnit := metricDefinition.PrimaryUnit

	queries := make([]string, len(queryTemplates))
	rangeParams := make([]string, len(queryTemplates))

	for i, queryTemplate := range queryTemplates {
		queryGenerator := metricDefinition.QueryGenerators[i]
		if queryGenerator != nil {
			queries[i], rangeParams[i] = queryGenerator(queryTemplate, queryParams)
		} else {
			queries[i] = queryTemplate
		}
	}

	// 프로메테우스 모니터링 API 호출
	responses := make([]interface{}, len(queries))

	// 조회된 데이터 중 최대값을 통한 단위 저장을 위함
	var maxValue float64
	var maxUnit string
	for queryIdx, query := range queries {
		// vector 쿼리와 range 쿼리에 따른 requestURL
		var escapedQuery = url.QueryEscape(query)
		var requestURL = config.ClientConfig.PrometheusRequestURL + queryAPIEndpoint + "?query=" + escapedQuery
		if rangeParams[queryIdx] != "" {
			requestURL = config.ClientConfig.PrometheusRequestURL + queryRangeAPIEndpoint + "?query=" + escapedQuery + rangeParams[queryIdx]
		}
		fmt.Println("[   QUERY    ]", query)
		request, err := http.NewRequest("GET", requestURL, nil)
		if err != nil {
			fmt.Println("failed to create http request, err=%s", err)
		}
		request.Header.Add("Content-Type", "application/json; charset=UTF-8")
		request.Header.Add("Access-Control-Allow-Origin", "*")
		request.Header.Add("Access-Control-Allow-Methods", "*")
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.ClientConfig.PrometheusToken))

		// 응답 요청
		response, err := client.Do(request)
		if err != nil {
			fmt.Println(
				"failed to call http request, err=%s",
				err,
			)
		}

		if response.StatusCode != 200 {
			break
		}

		// TODO 예외 처리
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("failed to read response body")
		}
		fmt.Println("[  RESPONSE  ]", string(responseBytes))

		// Primary 단위를 기준으로 컨버팅하는 값인지 확인
		isPrimaryUnit := common.Exists(common.UnitTypes[unitTypeKeys[queryIdx]].Units, primaryUnit)

		// 응답값 파싱
		var tempMaxValue float64
		responses[queryIdx], tempMaxValue = prometheus.ParseQueryResult(metricKey, isPrimaryUnit, responseBytes)
		fmt.Println("[   PARSED    ]", responses[queryIdx], tempMaxValue)
		if tempMaxValue > maxValue {
			maxValue = tempMaxValue
			// 최대값 단위 찾기
			if isPrimaryUnit && unitTypeKeys[queryIdx] != "" {
				maxUnit = common.FindMaxUnitByValues(unitTypeKeys[queryIdx], maxValue)
			}
		}
	}

	metricResponse := prometheus.MakeMetricResponse(metricKey, unitTypeKeys, maxUnit, subLabels, responses...)

	metricResponse.Label = metricDefinition.Label
	if maxUnit == "" {
		metricResponse.Unit = metricDefinition.PrimaryUnit
	} else {
		metricResponse.Unit = maxUnit
	}
	result[string(metricKey)] = metricResponse
	fmt.Println("[   RESULT   ]", metricResponse)

	return result
}
