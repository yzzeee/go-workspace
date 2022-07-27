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
		//"metricKeys": []string{"container_cpu"},
		"metricKeys": []string{"container_disk_io_read"},
		//"metricKeys": []string{"container_disk_io_write"},
		//"metricKeys": []string{"container_file_system"},
		//"metricKeys": []string{"container_memory"},
		//"metricKeys": []string{"container_network_in"},
		//"metricKeys": []string{"container_network_io"},
		//"metricKeys": []string{"container_network_out"},
		//"metricKeys": []string{"container_network_packet"},
		//"metricKeys": []string{"container_network_packet_drop"},
		//"metricKeys": []string{"custom_node_cpu"},
		//"metricKeys": []string{"custom_node_file_system"},
		//"metricKeys": []string{"custom_node_memory"},
		//"metricKeys": []string{"custom_quota_limit_cpu"},
		//"metricKeys": []string{"custom_quota_limit_memory"},
		//"metricKeys": []string{"custom_quota_request_cpu"},
		//"metricKeys": []string{"custom_quota_request_memory"},
		//"metricKeys": []string{"node_cpu"},
		//"metricKeys": []string{"node_cpu_load_average"},
		//"metricKeys": []string{"node_disk_io"},
		//"metricKeys": []string{"node_file_system"},
		//"metricKeys": []string{"node_memory"},
		//"metricKeys": []string{"node_network_in"},
		//"metricKeys": []string{"node_network_io"},
		//"metricKeys": []string{"node_network_out"},
		//"metricKeys": []string{"node_network_packet"},
		//"metricKeys": []string{"node_network_packet_drop"},
		//"metricKeys": []string{"number_of_container"},
		//"metricKeys": []string{"number_of_deployment"},
		//"metricKeys": []string{"number_of_ingress"},
		//"metricKeys": []string{"number_of_namespace"},
		//"metricKeys": []string{"number_of_pipeline"},
		//"metricKeys": []string{"number_of_pod"},
		//"metricKeys": []string{"number_of_service"},
		//"metricKeys": []string{"number_of_stateful_set"},
		//"metricKeys": []string{"number_of_volume"},
		//"metricKeys": []string{"quota_count_config_map_hard"},
		//"metricKeys": []string{"quota_count_config_map_used"},
		//"metricKeys": []string{"quota_count_persistent_volume_claim_hard"},
		//"metricKeys": []string{"quota_count_persistent_volume_claim_used"},
		//"metricKeys": []string{"quota_count_pod_hard"},
		//"metricKeys": []string{"quota_count_pod_used"},
		//"metricKeys": []string{"quota_count_replication_controller_hard"},
		//"metricKeys": []string{"quota_count_replication_controller_used"},
		//"metricKeys": []string{"quota_count_resource_quota_hard"},
		//"metricKeys": []string{"quota_count_resource_quota_used"},
		//"metricKeys": []string{"quota_count_secret_hard"},
		//"metricKeys": []string{"quota_count_secret_used"},
		//"metricKeys": []string{"quota_count_service_hard"},
		//"metricKeys": []string{"quota_count_service_used"},
		//"metricKeys": []string{"quota_count_service_load_balancer_hard"},
		//"metricKeys": []string{"quota_count_service_load_balancer_used"},
		//"metricKeys": []string{"quota_count_service_node_port_hard"},
		//"metricKeys": []string{"quota_count_service_node_port_used"},
		//"metricKeys": []string{"quota_limit_cpu_hard"},
		//"metricKeys": []string{"quota_limit_cpu_used"},
		//"metricKeys": []string{"quota_limit_memory_hard"},
		//"metricKeys": []string{"quota_limit_memory_used"},
		//"metricKeys": []string{"quota_request_cpu_hard"},
		//"metricKeys": []string{"quota_request_cpu_used"},
		//"metricKeys": []string{"quota_request_memory_hard"},
		//"metricKeys": []string{"quota_request_memory_used"},
		//"metricKeys": []string{"summary_node_info"},
		//"metricKeys": []string{"top_node_cpu_by_node"},
		//"metricKeys": []string{"top_node_file_system_by_node"},
		//"metricKeys": []string{"top_node_memory_by_node"},
		//"metricKeys": []string{"top_node_network_in_by_node"},
		//"metricKeys": []string{"top_node_network_out_by_node"},
		//"metricKeys": []string{"top_node_pod_count_by_node"},
		//"metricKeys": []string{"top5_container_cpu_by_namespace"},
		//"metricKeys": []string{"top5_container_cpu_by_pod"},
		//"metricKeys": []string{"top5_container_file_system_by_namespace"},
		//"metricKeys": []string{"top5_container_file_system_by_pod"},
		//"metricKeys": []string{"top5_container_memory_by_namespace"},
		//"metricKeys": []string{"top5_container_memory_by_pod"},
		//"metricKeys": []string{"top5_container_network_in_by_namespace"},
		//"metricKeys": []string{"top5_container_network_in_by_pod"},
		//"metricKeys": []string{"top5_container_network_out_by_namespace"},
		//"metricKeys": []string{"top5_container_network_out_by_pod"},
		//"metricKeys": []string{"top5_count_pod_by_namespace"},
		//"start": strconv.Itoa(int(now.Add(-time.Minute*5).Unix())),
		//"end": strconv.Itoa(int(now.Unix())),
		//"step": "120",
		//"node": "master3.ocp4.inno.com|worker2.ocp4.inno.com",
		//"instance": "master3.ocp4.inno.com|worker2.ocp4.inno.com",
		//"namespace": ".*",
	}

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
			//final, err = kubernetes.GetPipelines()
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
	//var clusterPrometheusVersion = "2.1.5-rc1"
	var clusterPrometheusVersion = "2.26.4"

	queryInfos := metricDefinition.QueryInfos
	definedVersions := make([]string, 0, len(queryInfos))
	for prometheusVersion := range queryInfos {
		definedVersions = append(definedVersions, string(prometheusVersion))
	}

	var targetVersion prometheus.PrometheusVersion
	clusterPrometheusVersion = prometheus.ParseVersion(clusterPrometheusVersion)

	index := common.IndexOf(definedVersions, clusterPrometheusVersion)

	if index != -1 { // 동일한 버전이 있는 경우
		targetVersion = prometheus.PrometheusVersion(definedVersions[index])
	} else { // 동일한 버전이 없는 경우
		targetVersion = prometheus.PrometheusVersion(prometheus.GetTargetPrometheusVersion(definedVersions, clusterPrometheusVersion))
	}

	referenceVersion := queryInfos[targetVersion].ReferenceVersion
	if referenceVersion != "" {
		targetVersion = referenceVersion
	}
	queryTemplates := queryInfos[targetVersion].QueryTemplates
	queryTemplateParsers := queryInfos[targetVersion].QueryTemplateParsers
	unitTypeKeys := metricDefinition.UnitTypeKeys
	primaryUnit := metricDefinition.PrimaryUnit

	queries := make([]string, len(queryTemplates))
	rangeParams := make([]string, len(queryTemplates))

	for i, queryTemplate := range queryTemplates {
		queryTemplateParser := queryTemplateParsers[i]
		if queryTemplateParser != nil {
			queries[i], rangeParams[i] = queryTemplateParser(queryTemplate, bodyParams)
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
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.ClientConfig.PrometheusToken))

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
