package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go-practice/common"
	"go-practice/http-client/config"
	"go-practice/http-client/kubernetes"
	"go-practice/http-client/prometheus"
	"io/ioutil"
	"net/http"
	"net/url"
	_ "strconv"
	_ "time"
)

const (
	queryAPIEndpoint      = "/api/v1/query"
	queryRangeAPIEndpoint = "/api/v1/query_range"
)

func init() {
	config.Init()
}

func main() {
	//now := time.Now()
	//now := time.Date(2022, 6, 27, 10, 15, 30, 0, time.Local)

	bodyParams := map[string]interface{}{
		"metricKeys": []string{"container_cpu"},
		//"start": strconv.Itoa(int(now.Add(-time.Minute*5).Unix())),
		//"end": strconv.Itoa(int(now.Unix())),
		//"step": "120",
		//"node": "master3.ocp4.inno.com|worker2.ocp4.inno.com",
		//"instance": "master3.ocp4.inno.com|worker2.ocp4.inno.com",
		//"namespace": "admin-workspace",
	}

	// -- single 모니터링 하려는 메트릭 키의 목록
	//var metricKeys = []string{"quota"}
	//var metricKeys = []string{"container_info"}
	//var metricKeys = []string{"container_cpu"}
	//var metricKeys = []string{"container_memory"}
	//var metricKeys = []string{"container_disk_io_reads"}
	//var metricKeys = []string{"container_disk_io_writes"}
	//var metricKeys = []string{"container_file_system"}
	//var metricKeys = []string{"container_network_in"}
	//var metricKeys = []string{"container_network_out"}
	//var metricKeys = []string{"container_network_packets_receive"}
	//var metricKeys = []string{"container_network_packets_transmit"}
	//var metricKeys = []string{"container_network_packets_receive_drop"}
	//var metricKeys = []string{"container_network_packets_transmit_drop"}
	//var metricKeys = []string{"container_pod_count"}
	//var metricKeys = []string{"number_of_container"}
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
	//var metricKeys = []string{"quota_object_count_configmaps"}
	//var metricKeys = []string{"quota_object_count_pods"}
	//var metricKeys = []string{"quota_object_count_secrets"}
	//var metricKeys = []string{"quota_object_count_replication_controllers"}
	//var metricKeys = []string{"quota_object_count_services"}
	//var metricKeys = []string{"quota_object_count_services_load_balancers"}
	//var metricKeys = []string{"quota_object_count_services_node_ports"}
	//var metricKeys = []string{"quota_object_count_resource_quotas"}
	//var metricKeys = []string{"quota_object_count_persistent_volume_claims"}
	//var metricKeys = []string{"quota_limit_range"}
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
	//var metricKeys = []string{"number_of_pipeline"}
	//var metricKeys = []string{"node_info"}

	var result = make(map[string]interface{})

	// 반환값
	for _, metricKey := range bodyParams["metricKeys"].([]string) {
		var final []byte

		// 클라이언트에서 요청한 key 에 따른 쿼리 생성
		metricDefinition, isMetric := prometheus.MetricDefinitions[prometheus.MetricKey(metricKey)]

		// 정의된 메트릭 여부 확인
		if isMetric {
			innerMetricKeys := metricDefinition.MetricKeys
			if innerMetricKeys != nil { // 다른 메트릭의 값을 활용하는 메트릭 처리
				innerResult := make(map[string]interface{})
				for _, innerMetricKey := range innerMetricKeys {
					innerResult = common.MergeJSONMaps(innerResult, getQueryResult(innerMetricKey, bodyParams))
				}
				metricResponse := prometheus.MakeMetricResponse(prometheus.MetricKey(metricKey), nil, "", nil, false, innerResult)
				result = common.MergeJSONMaps(result, map[string]interface{}{metricKey: metricResponse.Values})
			} else {
				queryResult := getQueryResult(prometheus.MetricKey(metricKey), bodyParams)
				result = common.MergeJSONMaps(result, queryResult)
			}
		}

		// 최종 결과 확인
		var err error
		if metricKey == "quota_limit_range" {
			final, err = kubernetes.GetLimitRange()
			if err != nil {
				fmt.Printf("failed to get limit range by Kubernetes API, err=%s\n", err)
			}
		} else if metricKey == "number_of_pipeline" {
			final, err = kubernetes.GetPipelines()
			if err != nil {
				fmt.Printf("failed to get pipeline by Kubernetes API, err=%s\n", err)
			}
		} else {
			final, _ = json.Marshal(result)
		}
		fmt.Println("[   FINAL   ]", string(final))
	}
}

func getQueryResult(metricKey prometheus.MetricKey, bodyParams map[string]interface{}) map[string]interface{} {
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

	metricDefinition := prometheus.MetricDefinitions[metricKey]
	label := metricDefinition.Label
	subLabels := metricDefinition.SubLabels
	if subLabels == nil {
		subLabels = []string{label}
	}
	queryTemplates := metricDefinition.QueryTemplates
	unitTypeKeys := metricDefinition.UnitTypeKeys
	primaryUnit := metricDefinition.PrimaryUnit

	queries := make([]string, len(queryTemplates))
	rangeParams := make([]string, len(queryTemplates))

	for i, queryTemplate := range queryTemplates {
		queryGenerator := metricDefinition.QueryGenerators[i]
		if queryGenerator != nil {
			queries[i], rangeParams[i] = queryGenerator(queryTemplate, bodyParams)
		} else {
			queries[i] = queryTemplate
		}
	}

	// 프로메테우스 모니터링 API 호출
	responses := make([]interface{}, len(queries))

	// 조회된 데이터 중 최대값을 통한 단위 저장을 위함
	var maxValue float64
	var maxUnit string
	var isRange bool
	for queryIdx, query := range queries {
		// vector 쿼리와 range 쿼리에 따른 requestURL
		var escapedQuery = url.QueryEscape(query)
		var requestURL = config.ClientConfig.PrometheusRequestURL + queryAPIEndpoint + "?query=" + escapedQuery
		isRange = rangeParams[queryIdx] != ""
		if isRange {
			requestURL = config.ClientConfig.PrometheusRequestURL + queryRangeAPIEndpoint + "?query=" + escapedQuery + rangeParams[queryIdx]
		}
		fmt.Println("[   QUERY    ]", query, requestURL)
		request, err := http.NewRequest("GET", requestURL, nil)
		if err != nil {
			fmt.Printf("failed to create http request, err=%s", err)
		}
		request.Header.Add("Content-Type", "application/json; charset=UTF-8")
		request.Header.Add("Access-Control-Allow-Origin", "*")
		request.Header.Add("Access-Control-Allow-Methods", "*")
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.ClientConfig.PrometheusRequestURL))

		// 응답 요청
		response, err := client.Do(request)
		if err != nil {
			fmt.Printf(
				"failed to call http request, err=%s\n",
				err,
			)
		}

		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("failed to read response body")
		}
		fmt.Println("[  RESPONSE  ]", string(responseBytes))

		// Primary 단위를 기준으로 컨버팅하는 값인지 확인
		isPrimaryUnit := common.Exists(common.UnitTypes[unitTypeKeys[queryIdx]].Units, primaryUnit)

		// 응답값 파싱
		var tempMaxValue float64
		responses[queryIdx], tempMaxValue = prometheus.ParseQueryResult(metricKey, isPrimaryUnit, responseBytes, isRange)
		fmt.Println("[   PARSED    ]", responses[queryIdx], tempMaxValue)
		if tempMaxValue > maxValue {
			maxValue = tempMaxValue
			// 최대값 단위 찾기
			if isPrimaryUnit && unitTypeKeys[queryIdx] != "" {
				maxUnit = common.FindMaxUnitByValues(unitTypeKeys[queryIdx], maxValue)
			}
		}
	}

	metricResponse := prometheus.MakeMetricResponse(metricKey, unitTypeKeys, maxUnit, subLabels, isRange, responses...)
	fmt.Println("[   RESULT   ]", metricResponse)

	metricResponse.Label = metricDefinition.Label
	if maxUnit == "" {
		metricResponse.Unit = metricDefinition.PrimaryUnit
	} else {
		metricResponse.Unit = maxUnit
	}
	metricResponse.Queries = queries
	result[string(metricKey)] = metricResponse

	return result
}
