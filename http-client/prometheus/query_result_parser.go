package prometheus

import (
	"encoding/json"
	"fmt"
	"strconv"

	"go-practice/common"
)

// ParseQueryResult responseBytes 의 에서 필요한 값을 파싱하고 결과값과,최대값을 반환하는 함수
/* (번호) responseBytes(파싱 전) => 반환값 형태(파싱 후)
* (1) {"status":"success","data":{"resultType":"vector","result":[{"metric":{},"value":[1657560872.452,"11.754666666666427"]}]}}
*       => 11.754666666666427
* (2) {"status":"success","data":{"resultType":"vector","result":[{"metric":{"instance":"worker1.ocp4.inno.com"},"value":[1657562191.538,"3.313939393939407"]},
                                                                  {"metric":{"instance":"worker2.ocp4.inno.com"},"value":[1657562191.538,"3.1159393939394797"]}]}}
*     => map[0:map[id:worker1.ocp4.inno.com order:0 timestamp:1.657562191538e+09 value:3.313939393939407] 1:map[id:worker2.ocp4.inno.com order:1 timestamp:1.657562191538e+09 value:3.1159393939394797]]
* (3) {"status":"success","data":{"resultType":"matrix","result":[{"metric":{},"values":[[1657561614,"4.194350475285014"],[1657561634,"4.313346351838768"]]}]}}
*       => [map[timestamp:1.657561614e+09 value:4.194350475285014] map[timestamp:1.657561634e+09 value:4.313346351838768]]
*/
func ParseQueryResult(metricKey MetricKey, isPrimaryUnit bool, responseBytes []byte, isRange bool) (interface{}, float64) {
	var result1 interface{}                 // (1)
	var result2 []interface{}               // (2)
	var result3 = make(map[int]interface{}) // (3)
	var maxValue float64
	var response = make(map[string]interface{})
	switch metricKey {
	case
		ContainerCpu, ContainerDiskIOReads, ContainerDiskIOWrites, ContainerFileSystem, ContainerMemory,
		ContainerNetworkIn, ContainerNetworkIO, ContainerNetworkOut, ContainerNetworkPacket, ContainerNetworkPacketDrop,
		CustomNodeCpu, CustomNodeFileSystem, CustomNodeMemory, CustomQuotaCpuLimit, CustomQuotaCpuRequest,
		CustomQuotaMemoryLimit, CustomQuotaMemoryRequest, NodeCpu, NodeCpuLoadAverage, NodeDiskIO,
		NodeFileSystem, NodeMemory, NodeNetworkIn, NodeNetworkIO, NodeNetworkOut,
		NodeNetworkPacket, NodeNetworkPacketDrop, NumberOfContainer, NumberOfDeployment, NumberOfIngress,
		NumberOfPod, NumberOfNamespace, NumberOfService, NumberOfStatefulSet, NumberOfVolume,
		QuotaCpuLimit, QuotaCpuRequest, QuotaMemoryLimit, QuotaMemoryRequest, QuotaObjectCountConfigmaps,
		QuotaObjectCountPods, QuotaObjectCountSecrets, QuotaObjectCountReplicationControllers,
		QuotaObjectCountServices, QuotaObjectCountServicesLoadBalancers, QuotaObjectCountServicesNodePorts,
		QuotaObjectCountResourceQuotas, QuotaObjectCountPersistentVolumeClaims:
		if !isRange { // (1)
			response = make(map[string]interface{})
			_ = json.Unmarshal(responseBytes, &response)
			if response["data"] != nil {
				for _, ele := range response["data"].(map[string]interface{})["result"].([]interface{}) {
					result1 = common.Get(ele, "value").([]interface{})[1]
					maxValue, _ = strconv.ParseFloat(fmt.Sprintf("%s", result1), 64)
				}
			}
		} else { // (2)
			_ = json.Unmarshal(responseBytes, &response)
			if len(common.Get(response, "data.result").([]interface{})) != 0 {
				for _, ele := range response["data"].(map[string]interface{})["result"].([]interface{})[0].(map[string]interface{})["values"].([]interface{}) {
					temp := make(map[string]interface{})
					temp["timestamp"] = common.Get(ele, "0")
					temp["value"] = common.Get(ele, "1")
					result2 = append(result2, temp)

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
	case // (3)
		TopNodeCpuByInstance, TopNodeFileSystemByInstance, TopNodeMemoryByInstance,
		TopNodeNetworkInByInstance, TopNodeNetworkOutByInstance, TopCountPodByNode,
		Top5ContainerCpuByNamespace, Top5ContainerCpuByPod, Top5ContainerFileSystemByNamespace,
		Top5ContainerFileSystemByPod, Top5ContainerMemoryByNamespace, Top5ContainerMemoryByPod,
		Top5ContainerNetworkInByNamespace, Top5ContainerNetworkInByPod, Top5ContainerNetworkOutByNamespace,
		Top5ContainerNetworkOutByPod, Top5CountPodByNamespace:
		response = make(map[string]interface{})
		_ = json.Unmarshal(responseBytes, &response)

		for i, ele := range response["data"].(map[string]interface{})["result"].([]interface{}) {
			temp := make(map[string]interface{})

			switch metricKey {
			case
				TopCountPodByNode:
				temp["id"] = common.Get(ele, "metric.node")
			case
				TopNodeCpuByInstance, TopNodeMemoryByInstance, TopNodeFileSystemByInstance,
				TopNodeNetworkInByInstance, TopNodeNetworkOutByInstance:
				temp["id"] = common.Get(ele, "metric.instance")
			case
				Top5ContainerCpuByNamespace, Top5ContainerFileSystemByNamespace, Top5ContainerMemoryByNamespace,
				Top5ContainerNetworkInByNamespace, Top5ContainerNetworkOutByNamespace,
				Top5CountPodByNamespace:
				temp["id"] = common.Get(ele, "metric.namespace")
			case
				Top5ContainerCpuByPod, Top5ContainerFileSystemByPod, Top5ContainerMemoryByPod,
				Top5ContainerNetworkInByPod, Top5ContainerNetworkOutByPod:
				temp["id"] = common.Get(ele, "metric.pod")
			}
			temp["timestamp"] = common.Get(ele, "value").([]interface{})[0] // value 의 첫 번째 원소는 timestamp
			temp["value"] = common.Get(ele, "value").([]interface{})[1]     // value 의 두 번째 원소는 메트릭 값
			temp["order"] = i                                               // 순서 보장 안되므로 정렬을 위한 인덱스를 넣어줌
			result3[i] = temp

			// 다중 값 중 최대값 저장 후 반환
			if isPrimaryUnit {
				float, _ := strconv.ParseFloat(fmt.Sprintf("%s", temp["value"]), 64)
				if maxValue < float {
					maxValue = float
				}
			}
		}
	}
	if result1 != nil {
		return result1, maxValue
	}
	if result2 != nil {
		return result2, maxValue
	}
	if len(result3) != 0 {
		return result3, maxValue
	}
	return nil, 0
}
