package prometheus

import (
	"fmt"
	"go-practice/common"
	"strconv"
)

// MetricResponse 메트릭 응답 구조체
type MetricResponse struct {
	Label      string      `json:"label,omitempty"`
	Usage      string      `json:"usage,omitempty"`
	Total      string      `json:"total,omitempty"`
	Percentage string      `json:"percentage,omitempty"`
	Unit       string      `json:"unit,omitempty"`
	Values     interface{} `json:"values,omitempty"`
	Error      interface{} `json:"error,omitempty"`
}

// MakeMetricResponse responseBytes 의 에서 필요한 값을 파싱하고 결과값과,최대값을 반환하는 함수
/* (번호) | responseBytes 포멧(파싱 전) | 결과값 포멧(파싱 후)
 * (1) | { } | { }
 * (2) | { } | { }
 * (3) | { } | { }
 * (4) | { } | { }
 * (5) | { } | { }
 */
func MakeMetricResponse(metricKey MetricKey, unitTypeKeys []common.UnitTypeKey, maxValueUnit string, subLabels []string, resultSets ...interface{}) MetricResponse {
	switch metricKey {
	case // (1)
		NodeCpu, NodeFileSystem, NodeMemory, QuotaCpuLimit, QuotaCpuRequest,
		QuotaMemoryLimit, QuotaMemoryRequest:
		if len(resultSets) != 0 {
			var resultSet0, _ = strconv.ParseFloat(fmt.Sprintf("%s", resultSets[0]), 64)
			var resultSet1, _ = strconv.ParseFloat(fmt.Sprintf("%s", resultSets[1]), 64)
			var resultSet2, _ = strconv.ParseFloat(fmt.Sprintf("%s", resultSets[2]), 64)
			if unitTypeKeys != nil && unitTypeKeys[0] != "" {
				resultSet0 = common.Humanize(resultSet0, unitTypeKeys[0], &common.HumanizeOptions{PreferredUnit: maxValueUnit, Precision: 2}).Value
			}
			if unitTypeKeys != nil && unitTypeKeys[1] != "" {
				resultSet1 = common.Humanize(resultSet1, unitTypeKeys[1], &common.HumanizeOptions{PreferredUnit: maxValueUnit, Precision: 2}).Value
			}
			if unitTypeKeys != nil && unitTypeKeys[2] != "" {
				resultSet2 = common.Humanize(resultSet2, unitTypeKeys[2], &common.HumanizeOptions{PreferredUnit: maxValueUnit, Precision: 2}).Value
			}
			return MetricResponse{
				Usage:      fmt.Sprintf("%v", resultSet0),
				Total:      fmt.Sprintf("%v", resultSet1),
				Percentage: fmt.Sprintf("%v", resultSet2),
			}
		}
	case // (2)
		ContainerCpu, ContainerFileSystem, ContainerMemory, ContainerNetworkIn, ContainerNetworkOut,
		NumberOfDeployment, NumberOfIngress, NumberOfPod, NumberOfNamespace, NumberOfService,
		NumberOfStatefulSet, NumberOfVolume, NodeNetworkIn, NodeNetworkOut:
		var resultSet0, _ = strconv.ParseFloat(fmt.Sprintf("%s", resultSets[0]), 64)

		if unitTypeKeys != nil && unitTypeKeys[0] != "" {
			resultSet0 = common.Humanize(resultSet0, unitTypeKeys[0], &common.HumanizeOptions{PreferredUnit: maxValueUnit, Precision: 2}).Value
		}
		return MetricResponse{
			Usage: fmt.Sprintf("%v", resultSet0),
		}
	case // (3)
		TopNodeCpuByInstance, TopNodeFileSystemByInstance, TopNodeMemoryByInstance, TopNodeNetworkInByInstance, TopNodeNetworkOutByInstance,
		TopCountPodByNode, Top5ContainerCpuByNamespace, Top5ContainerCpuByPod, Top5ContainerFileSystemByNamespace, Top5ContainerFileSystemByPod,
		Top5ContainerMemoryByNamespace, Top5ContainerMemoryByPod, Top5ContainerNetworkInByNamespace, Top5ContainerNetworkInByPod, Top5ContainerNetworkOutByNamespace,
		Top5ContainerNetworkOutByPod, Top5CountPodByNamespace:
		if len(resultSets) != 0 {
			var resultSet0 = resultSets[0].(map[int]interface{})
			if unitTypeKeys != nil && unitTypeKeys[0] != "" {
				for _, values := range resultSet0 {
					temp := values.(map[string]interface{})
					var value, _ = strconv.ParseFloat(fmt.Sprintf("%s", temp["value"]), 64)
					temp["value"] = strconv.FormatFloat(common.Humanize(value, unitTypeKeys[0], &common.HumanizeOptions{PreferredUnit: maxValueUnit, Precision: 2}).Value, 'f', -1, 64)
					temp["unit"] = common.Humanize(value, unitTypeKeys[0], &common.HumanizeOptions{PreferredUnit: maxValueUnit, Precision: 2}).Unit
				}
			}
			return MetricResponse{
				Values: resultSet0,
			}
		}
	case // (4)
		RangeContainerCpu, RangeContainerMemory, RangeContainerNetworkIO, RangeContainerNetworkPacket, RangeDiskIO,
		RangeFileSystem, RangeNodeCpu, RangeNetworkBandwidth, RangeNetworkPacketReceiveTransmit, RangeNetworkPacketReceiveTransmitDrop,
		RangeNodeCpuLoadAverage, RangeNodeMemory, RangeNodeNetworkIO, RangeNodeNetworkPacket:
		var resultSet0 []interface{}
		if len(resultSets) != 0 {
			resultSet0 = resultSets[0].([]interface{})
			if unitTypeKeys != nil && unitTypeKeys[0] != "" {
				for idx, values := range resultSets[0].([]interface{}) {
					temp := values.(map[string]interface{})
					var value, _ = strconv.ParseFloat(fmt.Sprintf("%s", temp["value"]), 64)
					temp[subLabels[0]] = common.Humanize(value, unitTypeKeys[0], &common.HumanizeOptions{PreferredUnit: maxValueUnit, Precision: 2}).Value
					for j := 1; j < len(resultSets); j++ {
						var tempJ = resultSets[j].([]interface{})[idx].(map[string]interface{})["value"]
						var valueJ, _ = strconv.ParseFloat(fmt.Sprintf("%s", tempJ), 64)
						temp[subLabels[j]] = common.Humanize(valueJ, unitTypeKeys[j], &common.HumanizeOptions{PreferredUnit: maxValueUnit, Precision: 2}).Value
					}
					delete(temp, "value")
				}
			}
		}
		return MetricResponse{
			Values: resultSet0,
		}
	case // (5)
		SummaryNodeInfo:
		values := make(map[string]interface{})
		if len(resultSets) != 0 {
			resultSet0 := resultSets[0].(map[string]interface{})
			for key, value := range resultSet0 {
				switch MetricKey(key) {
				case NodeCpu, NodeFileSystem, NodeMemory:
					label := fmt.Sprintf("%s", common.Get(value, "Label"))
					usage := common.Get(value, "Usage")
					percentage := common.Get(value, "Percentage")
					unit := common.Get(value, "Unit")
					values[label] = fmt.Sprintf("%s %s (%s%%)", usage, unit, percentage)
				case NodeNetworkIn, NodeNetworkOut:
					label := fmt.Sprintf("%s", common.Get(value, "Label"))
					usage := common.Get(value, "Usage")
					unit := common.Get(value, "Unit")
					values[label] = fmt.Sprintf("%s %s", usage, unit)
				case NumberOfPod:
					label := fmt.Sprintf("%s", common.Get(value, "Label"))
					usage := common.Get(value, "Usage")
					values[label] = fmt.Sprintf("%s", usage)
				}
			}
			return MetricResponse{
				Values: values,
			}
		}
	}
	return MetricResponse{}
}
