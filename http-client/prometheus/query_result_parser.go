package prometheus

import (
	"encoding/json"
	"fmt"
	"go-practice/common"
	"strconv"
)

// ParseQueryResult responseBytes 의 에서 필요한 값을 파싱하고 결과값과,최대값을 반환하는 함수
/* (번호) | responseBytes 포멧(파싱 전) | 결과값 포멧(파싱 후)
 * (1) | { } | { }
 * (2) | { } | { }
 * (3) | { } | { }
 */
func ParseQueryResult(metricKey MetricKey, isPrimaryUnit bool, responseBytes []byte) (interface{}, float64) {
	var result1 interface{}                 // (1)
	var result2 = make(map[int]interface{}) // (2)
	var result3 []interface{}               // (3)
	var maxValue float64
	var response = make(map[string]interface{})
	switch metricKey {
	case // (1)
		Quota, ContainerCpu, ContainerMemory,
		ContainerFileSystem, ContainerNetworkIn, ContainerNetworkOut,
		NodeCpu, NodeMemory, NodeFileSystem,
		NodeNetworkIn, NodeNetworkOut, NodePodCount,
		QuotaCpuRequest, QuotaCpuLimit, QuotaMemoryRequest,
		QuotaMemoryLimit:
		response = make(map[string]interface{})
		_ = json.Unmarshal(responseBytes, &response)
		for _, ele := range response["data"].(map[string]interface{})["result"].([]interface{}) {
			result1 = common.Get(ele, "value").([]interface{})[1]
			maxValue, _ = strconv.ParseFloat(fmt.Sprintf("%s", result1), 64)
		}
	case // (2)
		NodeCpuTop, NodeCpuTop5Projects, NodeCpuTop5Pods,
		NodeMemoryTop, NodeMemoryTop5Projects, NodeMemoryTop5Pods,
		NodeFileSystemTop, NodeFileSystemTop5Projects, NodeFileSystemTop5Pods,
		NodeNetworkInTop, NodeNetworkInTop5Projects, NodeNetworkInTop5Pods,
		NodeNetworkOutTop, NodeNetworkOutTop5Projects, NodeNetworkOutTop5Pods,
		NodePodCountTop, NodePodCountTop5Projects:
		response = make(map[string]interface{})
		_ = json.Unmarshal(responseBytes, &response)

		for i, ele := range response["data"].(map[string]interface{})["result"].([]interface{}) {
			temp := make(map[string]interface{})

			switch metricKey {
			case
				NodePodCountTop:
				temp["id"] = common.Get(ele, "metric.node")
			case
				NodeCpuTop, NodeMemoryTop, NodeFileSystemTop,
				NodeNetworkInTop, NodeNetworkOutTop:
				temp["id"] = common.Get(ele, "metric.instance")
			case
				NodeCpuTop5Projects, NodeMemoryTop5Projects, NodeFileSystemTop5Projects,
				NodeNetworkInTop5Projects, NodeNetworkOutTop5Projects, NodePodCountTop5Projects:
				temp["id"] = common.Get(ele, "metric.namespace")
			case
				NodeCpuTop5Pods, NodeMemoryTop5Pods, NodeFileSystemTop5Pods,
				NodeNetworkInTop5Pods, NodeNetworkOutTop5Pods:
				temp["id"] = common.Get(ele, "metric.pod")
			}
			temp["timestamp"] = common.Get(ele, "value").([]interface{})[0] // value 의 첫 번째 원소는 timestamp
			temp["value"] = common.Get(ele, "value").([]interface{})[1]     // value 의 두 번째 원소는 메트릭 값
			temp["order"] = i                                               // 순서 보장 안되므로 정렬을 위한 인덱스를 넣어줌
			result2[i] = temp

			// 다중 값 중 최대값 저장 후 반환
			if isPrimaryUnit {
				float, _ := strconv.ParseFloat(fmt.Sprintf("%s", temp["value"]), 64)
				if maxValue < float {
					maxValue = float
				}
			}
		}
	case // (3)
		RangeNodeCpuUsage, RangeContainerCpuUsage, RangeCpuLoadAverage,
		RangeNodeMemoryUsage, RangeContainerMemoryUsage, RangeMemorySwap,
		RangeNodeNetworkIO, RangeContainerNetworkIO, RangeNodeNetworkPacket,
		RangeContainerNetworkPacket, RangeNetworkBandwidth, RangeNetworkPacketReceiveTransmit,
		RangeNetworkPacketReceiveTransmitDrop, RangeFileSystem, RangeDiskIO:
		_ = json.Unmarshal(responseBytes, &response)
		fmt.Println(len(common.Get(response, "data.result").([]interface{})))
		if len(common.Get(response, "data.result").([]interface{})) != 0 {
			for _, ele := range response["data"].(map[string]interface{})["result"].([]interface{})[0].(map[string]interface{})["values"].([]interface{}) {
				temp := make(map[string]interface{})
				temp["timestamp"] = common.Get(ele, "0")
				temp["value"] = common.Get(ele, "1")
				result3 = append(result3, temp)

				// 다중 값 중 최대값 저장 후 반환
				if isPrimaryUnit {
					float, _ := strconv.ParseFloat(fmt.Sprintf("%s", temp["value"]), 64)
					if maxValue < float {
						maxValue = float
					}
				}
			}
		}
	}

	if result1 != nil {
		return result1, maxValue
	}
	if len(result2) != 0 {
		return result2, maxValue
	}
	if result3 != nil {
		return result3, maxValue
	}
	return nil, 0
}
