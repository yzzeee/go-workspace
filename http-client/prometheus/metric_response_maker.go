package prometheus

import (
	"fmt"
	"strconv"

	"go-practice/common"
)

// MetricResponse 메트릭 응답
type MetricResponse struct {
	Label      string      `json:"label,omitempty"`
	Usage      string      `json:"usage,omitempty"`
	RawUsage   string      `json:"rawUsage,omitempty"`
	Total      string      `json:"total,omitempty"`
	Percentage string      `json:"percentage,omitempty"`
	Unit       string      `json:"unit,omitempty"`
	Values     interface{} `json:"values,omitempty"`
	Error      interface{} `json:"error,omitempty"`
	Queries    []string    `json:"queries,omitempty"`
}

// MakeMetricResponse QueryTemplates의 수와 동일한 resultSet이 인자로 들어오고 해당 resultSet을 이용하여 응답값을 만드는 함수
//nolint:gocyclo
func MakeMetricResponse(metricKey MetricKey, unitTypeKeys []common.UnitTypeKey,
	maxValueUnit string, subLabels []string, isRange bool, resultSets ...interface{}) MetricResponse {
	switch metricKey {
	case
		ContainerCpu, ContainerDiskIORead, ContainerDiskIOWrite, ContainerFileSystem, ContainerMemory,
		ContainerNetworkIn, ContainerNetworkIO, ContainerNetworkOut, ContainerNetworkPacket, ContainerNetworkPacketDrop,
		HaProxyTrafficIn, HaProxyTrafficOut, HaProxyConnectionRate, NodeCpu, NodeCpuLoadAverage,
		NodeDiskIO, NodeFileSystem, NodeMemory, NodeNetworkIn, NodeNetworkIO,
		NodeNetworkOut, NodeNetworkPacket, NodeNetworkPacketDrop, NumberOfContainer, NumberOfDeployment,
		NumberOfIngress, NumberOfNamespace, NumberOfPod, NumberOfService, NumberOfStatefulSet,
		NumberOfVolume, QuotaCountConfigMapHard, QuotaCountConfigMapUsed,
		QuotaCountPersistentVolumeClaimHard, QuotaCountPersistentVolumeClaimUsed, QuotaCountPodHard,
		QuotaCountPodUsed, QuotaCountReplicationControllerHard, QuotaCountReplicationControllerUsed,
		QuotaCountResourceQuotaHard, QuotaCountResourceQuotaUsed, QuotaCountSecretHard,
		QuotaCountSecretUsed, QuotaCountServiceHard, QuotaCountServiceUsed,
		QuotaCountServiceLoadBalancerHard, QuotaCountServiceLoadBalancerUsed, QuotaCountServiceNodePortHard,
		QuotaCountServiceNodePortUsed, QuotaLimitCpuHard, QuotaLimitCpuUsed, QuotaLimitMemoryHard, QuotaLimitMemoryUsed,
		QuotaLimitPodCpu, QuotaLimitPodEphemeralStorage, QuotaLimitPodMemory,
		QuotaRequestCpuHard, QuotaRequestCpuUsed, QuotaRequestMemoryHard, QuotaRequestMemoryUsed, QuotaRequestPodCpu,
		QuotaRequestPodEphemeralStorage, QuotaRequestPodMemory, QuotaRequestStorageHard, QuotaRequestStorageUsed:
		if !isRange {
			var rawUsage = resultSets[0]
			var resultSet0, _ = strconv.ParseFloat(fmt.Sprintf("%s", resultSets[0]), 64)

			if rawUsage == nil {
				rawUsage = ""
			}

			if unitTypeKeys != nil && unitTypeKeys[0] != "" {
				resultSet0 = common.Humanize(resultSet0, unitTypeKeys[0],
					&common.HumanizeOptions{PreferredUnit: maxValueUnit, Precision: 2}).Value
			}
			return MetricResponse{
				Usage:    strconv.FormatFloat(resultSet0, 'f', -1, 64),
				RawUsage: fmt.Sprintf("%s", rawUsage),
			}
		} else {
			var resultSet0 []interface{}
			var maxIdx = 0
			var maxListSize = 0
			for i := 0; i < len(resultSets); i++ {
				if resultSets[i] != nil && len(resultSets[i].([]interface{})) > maxListSize {
					maxListSize = len(resultSets[i].([]interface{}))
					maxIdx = i
				}
			}
			if resultSets[maxIdx] != nil {
				resultSet0 = resultSets[maxIdx].([]interface{})
				for idx, values := range resultSets[maxIdx].([]interface{}) {
					temp := values.(map[string]interface{})
					var value, _ = strconv.ParseFloat(fmt.Sprintf("%s", temp["value"]), 64)
					temp[subLabels[maxIdx]] = common.Humanize(value, unitTypeKeys[maxIdx],
						&common.HumanizeOptions{PreferredUnit: maxValueUnit, Precision: 2}).Value
					for j := 0; j < len(resultSets); j++ {
						if maxIdx != j {
							if resultSets[j] == nil {
								temp[subLabels[j]] = 0

							} else if len(resultSets[j].([]interface{})) > idx {
								var tempJ = resultSets[j].([]interface{})[idx].(map[string]interface{})["value"]
								var valueJ, _ = strconv.ParseFloat(fmt.Sprintf("%s", tempJ), 64)
								temp[subLabels[j]] = common.Humanize(valueJ, unitTypeKeys[j],
									&common.HumanizeOptions{PreferredUnit: maxValueUnit, Precision: 2}).Value
							}
						}
					}
					delete(temp, "value")
				}
			}
			return MetricResponse{
				Values: resultSet0,
			}
		}
	case
		CustomContainerVolume, CustomNodeCpu, CustomNodeFileSystem, CustomNodeMemory, CustomQuotaLimitCpu,
		CustomQuotaLimitMemory, CustomQuotaRequestCpu, CustomQuotaRequestMemory:
		if len(resultSets) != 0 {
			var resultSet0, _ = strconv.ParseFloat(fmt.Sprintf("%s", resultSets[0]), 64)
			var resultSet1, _ = strconv.ParseFloat(fmt.Sprintf("%s", resultSets[1]), 64)
			var resultSet2, _ = strconv.ParseFloat(fmt.Sprintf("%s", resultSets[2]), 64)
			if unitTypeKeys != nil && unitTypeKeys[0] != "" {
				resultSet0 = common.Humanize(resultSet0, unitTypeKeys[0],
					&common.HumanizeOptions{PreferredUnit: maxValueUnit, Precision: 2}).Value
			}
			if unitTypeKeys != nil && unitTypeKeys[1] != "" {
				resultSet1 = common.Humanize(resultSet1, unitTypeKeys[1],
					&common.HumanizeOptions{PreferredUnit: maxValueUnit, Precision: 2}).Value
			}
			if unitTypeKeys != nil && unitTypeKeys[2] != "" {
				resultSet2 = common.Humanize(resultSet2, unitTypeKeys[2],
					&common.HumanizeOptions{PreferredUnit: maxValueUnit, Precision: 2}).Value
			}
			return MetricResponse{
				Usage:      fmt.Sprintf("%v", resultSet0),
				Total:      fmt.Sprintf("%v", resultSet1),
				Percentage: fmt.Sprintf("%v", resultSet2),
			}
		}
	case
		TopNodeCpuByNode, TopNodeFileSystemByNode, TopNodeMemoryByNode,
		TopNodeNetworkInByNode, TopNodeNetworkOutByNode, TopNodePodCountByNode,
		Top5ContainerCpuByNamespace, Top5ContainerCpuByPod, Top5ContainerFileSystemByNamespace,
		Top5ContainerFileSystemByPod, Top5ContainerMemoryByNamespace, Top5ContainerMemoryByPod,
		Top5ContainerNetworkInByNamespace, Top5ContainerNetworkInByPod, Top5ContainerNetworkOutByNamespace,
		Top5ContainerNetworkOutByPod, Top5CountContainerByPod, Top5CountPodByNamespace:
		if len(resultSets) != 0 && resultSets[0] != nil {
			var resultSet0 = resultSets[0].(map[int]interface{})
			if unitTypeKeys != nil && unitTypeKeys[0] != "" {
				for _, values := range resultSet0 {
					temp := values.(map[string]interface{})
					var value, _ = strconv.ParseFloat(fmt.Sprintf("%s", temp["value"]), 64)
					temp["value"] = strconv.FormatFloat(common.Humanize(value, unitTypeKeys[0],
						&common.HumanizeOptions{PreferredUnit: maxValueUnit, Precision: 2}).Value, 'f', -1, 64)
					temp["unit"] = common.Humanize(value, unitTypeKeys[0],
						&common.HumanizeOptions{PreferredUnit: maxValueUnit, Precision: 2}).Unit
				}
			}
			return MetricResponse{
				Values: resultSet0,
			}
		}
	case
		SummaryNodeInfo:
		values := make(map[string]interface{})
		if len(resultSets) != 0 && resultSets[0] != nil {
			resultSet0 := resultSets[0].(map[string]interface{})
			for key, value := range resultSet0 {
				switch MetricKey(key) {
				case CustomNodeCpu, CustomNodeFileSystem, CustomNodeMemory:
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
	case
		SummaryContainerCpuInfo, SummaryContainerMemoryInfo, SummaryCpuQuotaInfo, SummaryMemoryQuotaInfo:
		values := make(map[string]interface{})
		if len(resultSets) != 0 && resultSets[0] != nil {
			resultSet0 := resultSets[0].(map[string]interface{})
			var usagePercentage interface{}
			var rawLimitValue interface{}

			for key, value := range resultSet0 {
				fmt.Println(key)
				switch MetricKey(key) {
				case
					ContainerCpu, ContainerMemory, QuotaRequestCpuHard,
					QuotaRequestMemoryHard, QuotaRequestPodCpu, QuotaRequestPodMemory:
					var label string
					switch MetricKey(key) {
					case ContainerCpu, ContainerMemory:
						label = "used"
					case QuotaRequestCpuHard, QuotaRequestMemoryHard, QuotaRequestPodCpu, QuotaRequestPodMemory:
						label = "request"
					}
					limitVal := make(map[string]interface{})
					val := make(map[string]interface{})
					if rawLimitValue == nil {
						var limitValue interface{}
						switch MetricKey(key) {
						case ContainerCpu, QuotaRequestCpuHard, QuotaRequestPodCpu:
							if metricKey == SummaryCpuQuotaInfo {
								limitValue = resultSet0[string(QuotaLimitCpuHard)]
							} else {
								limitValue = resultSet0[string(QuotaLimitPodCpu)]
							}
							rawLimitValue = common.Get(limitValue, "RawUsage")
						case ContainerMemory, QuotaRequestMemoryHard, QuotaRequestPodMemory:
							if metricKey == SummaryCpuQuotaInfo {
								limitValue = resultSet0[string(QuotaLimitCpuHard)]
							} else {
								limitValue = resultSet0[string(QuotaLimitPodMemory)]
							}
							rawLimitValue = common.Get(limitValue, "RawUsage")
						}
						limitUsage := common.Get(limitValue, "Usage")
						limitUnit := common.Get(limitValue, "Unit")
						limitVal["value"] = limitUsage
						limitVal["unit"] = limitUnit
						var limitPercentage int
						if limitUsage != "0" {
							limitPercentage = 100
						}
						limitVal["percentage"] = limitPercentage
						values["limit"] = limitVal
					}
					rawUsage := common.Get(value, "RawUsage")
					usage := common.Get(value, "Usage")
					unit := common.Get(value, "Unit")
					val["value"] = usage
					val["unit"] = unit

					var percentage interface{}
					if rawLimitValue == "0" || rawLimitValue == nil || rawLimitValue == "" {
						if usage != "0" {
							percentage = 100
						}
					} else {
						if rawUsage != "0" && rawUsage != "" && rawUsage != nil {
							floatUsage, err := strconv.ParseFloat(rawUsage.(string), 64)
							if err != nil {
								fmt.Println("failed to read response body, err=%s\n", err)
							}
							limitFloat, err := strconv.ParseFloat(rawLimitValue.(string), 64)
							if err != nil {
								fmt.Println("failed to read response body, err=%s\n", err)
							}
							percentage = common.RoundFloat(floatUsage/limitFloat*100, 2)

							switch MetricKey(key) {
							case ContainerCpu, ContainerMemory:
								usagePercentage = percentage
							}
							if percentage.(float64) > 100 {
								percentage = 100
							}
						}
					}
					val["percentage"] = percentage
					values[label] = val
				}
			}
			values["percentage"] = usagePercentage

			return MetricResponse{
				Values: values,
			}
		}
	}

	return MetricResponse{}
}
