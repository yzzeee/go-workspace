package prometheus

import (
	"go-practice/common"
	"net/url"
)

// QueryGenerators 쿼리 템플릿과 쿼리 파라미터를 인자로 받아 쿼리를 반환하는 함수 목록 타입
type QueryGenerators []func(queryTemplate string, queryParams url.Values) (string, string)

// MetricDefinition 메트릭 정의 구조체
type MetricDefinition struct {
	Label           string               // 메트릭의 라벨
	SubLabels       []string             // 쿼리 템플릿의 라벨(미필수)
	QueryTemplates  []string             // 쿼리 템플릿
	QueryGenerators QueryGenerators      // 쿼리 템플릿에 조건절 추가하여 쿼리를 반환하는 함수 목록(쿼리 템플릿과 맵핑)
	UnitTypeKeys    []common.UnitTypeKey // 쿼리 결과값의 단위 티입의 키 목록(쿼리 템플릿과 맵핑)
	PrimaryUnit     string               // 쿼리 결과값의 단위 중 주단위

	MetricKeys []MetricKey // 다른 메트릭 정의를 활용하는 메트릭(다른 메트릭 활용 시 해당 값만 작성)
}

// MetricDefinitions 메트릭 키에 따른 메트릭 정의 상수
var (
	MetricDefinitions = map[MetricKey]MetricDefinition{
		Quota: {
			Label: "CPU",
			QueryTemplates: []string{
				// 노드의 CPU Core 사용량(Core)
				"sum(rate(node_cpu_seconds_total{mode!=\"idle\",mode!=\"iowait\",instance=~\"%s\"}[3m]))",
				// 노드의 CPU Core 수
				"sum(kube_node_status_capacity{resource=\"cpu\",unit=\"core\",node=~\"%s\"})",
				// 총 CPU Core 사용량(%)
				"sum(rate(node_cpu_seconds_total{mode!=\"idle\",mode!=\"iowait\",instance=~\"%s\"}[3m]))/sum(kube_node_status_capacity{resource=\"cpu\",unit=\"core\",node=~\"%s\"})*100",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}, false),
				queryGenerator([]interface{}{"node"}, false),
				queryGenerator([]interface{}{"instance", "node"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
				"",
				common.Percentage,
			},
			PrimaryUnit: "Core",
		},
		ContainerCpu: {
			Label: "CPU",
			QueryTemplates: []string{
				// 컨테이너의 CPU Core 사용량(Core)
				"sum(rate(container_cpu_usage_seconds_total{container!=\"\",pod!=\"\",namespace=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		ContainerMemory: {
			Label: "MEMORY",
			QueryTemplates: []string{
				// 컨테이너의 메모리 사용량(byte)
				"sum(container_memory_working_set_bytes{cluster=\"\",container!=\"\",namespace=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		ContainerFileSystem: { // TODO
			Label: "FILE SYSTEM",
			QueryTemplates: []string{
				// 노드의 파일 시스템 사용량(byte)
				"sum(node_filesystem_size_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"}-node_filesystem_avail_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance", "instance"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		ContainerNetworkIn: { // TODO
			Label: "NETWORK IN",
			QueryTemplates: []string{
				// 노드의 NETWORK IN(bps)
				"sum(rate(node_network_receive_bytes_total{instance=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		ContainerNetworkOut: { // TODO
			Label: "NETWORK OUT",
			QueryTemplates: []string{
				// 노드의 NETWORK OUT(bps)
				"sum(rate(node_network_receive_bytes_total{instance=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		NodeInfo: {
			MetricKeys: []MetricKey{NodeCpu, NodeMemory, NodeFileSystem, NodeNetworkIn, NodeNetworkOut, NodePodCount},
		},
		NodeCpu: {
			Label: "CPU",
			QueryTemplates: []string{
				// 노드의 CPU Core 사용량(Core)
				"sum(rate(node_cpu_seconds_total{mode!=\"idle\",mode!=\"iowait\",instance=~\"%s\"}[3m]))",
				// 노드의 CPU Core 수
				"sum(kube_node_status_capacity{resource=\"cpu\",unit=\"core\",node=~\"%s\"})",
				// 총 CPU Core 사용량(%)
				"sum(rate(node_cpu_seconds_total{mode!=\"idle\",mode!=\"iowait\",instance=~\"%s\"}[3m]))/sum(kube_node_status_capacity{resource=\"cpu\",unit=\"core\",node=~\"%s\"})*100",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}, false),
				queryGenerator([]interface{}{"node"}, false),
				queryGenerator([]interface{}{"instance", "node"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
				"",
				common.Percentage,
			},
			PrimaryUnit: "Core",
		},
		NodeCpuTop: {
			Label: "CPU",
			QueryTemplates: []string{
				// 노드의 CPU 사용량에 따른 노드 내림차순 목록
				"sort_desc(sum(rate(node_cpu_seconds_total{mode!=\"idle\",mode!=\"iowait\",instance=~\"%s\"}[3m]))by(instance))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		NodeCpuTop5Projects: {
			Label: "CPU(TOP5 OF PROJECTS)",
			QueryTemplates: []string{
				// 노드의 네임스페이스 중 CPU 사용량에 따른 내림차순 목록(TOP5)
				"topk(5,sort_desc(sum(rate(container_cpu_usage_seconds_total{container!=\"\",pod!=\"\",node=~\"%s\"}[3m]))by(namespace)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		NodeCpuTop5Pods: {
			Label: "CPU(TOP5 OF PODS)",
			QueryTemplates: []string{
				// 노드의 파드 중 CPU 사용량에 따른 내림차순 목록(TOP5)
				"topk(5,sort_desc(sum(rate(container_cpu_usage_seconds_total{container!=\"\",pod!=\"\",node=~\"%s\",namespace=~\"%s\"}[3m]))by(pod)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node", "namespace"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		NodeMemory: {
			Label: "MEMORY",
			QueryTemplates: []string{
				// 노드의 메모리 사용량(byte)
				"sum(node_memory_MemTotal_bytes{instance=~\"%s\"}-node_memory_MemAvailable_bytes{instance=~\"%s\"})",
				// 노드의 총 메모리 크기
				"sum(node_memory_MemTotal_bytes{instance=~\"%s\"})",
				// 노드의 메모리 사용량(%)
				"sum(node_memory_MemTotal_bytes{instance=~\"%s\"}-node_memory_MemAvailable_bytes{instance=~\"%s\"})/sum(node_memory_MemTotal_bytes{instance=~\"%s\"})*100",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance", "instance"}, false),
				queryGenerator([]interface{}{"instance"}, false),
				queryGenerator([]interface{}{"instance", "instance", "instance"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
				common.BinaryBytes,
				common.Percentage,
			},
			PrimaryUnit: "B",
		},
		NodeMemoryTop: {
			Label: "MEMORY",
			QueryTemplates: []string{
				// 노드의 MEMORY 사용량에 따른 노드 내림차순 목록
				"sort_desc(sum(node_memory_MemTotal_bytes-node_memory_MemAvailable_bytes{instance=~\"%s\"})by(instance))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		NodeMemoryTop5Projects: {
			Label: "MEMORY(TOP5 OF PROJECTS)",
			QueryTemplates: []string{
				// 노드의 네임스페이스 중 MEMORY 사용량에 따른 내림차순 목록(TOP5)
				"topk(5,sort_desc(sum(container_memory_working_set_bytes{container!=\"\",pod!=\"\",node=~\"%s\"})by(namespace)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		NodeMemoryTop5Pods: {
			Label: "MEMORY(TOP5 OF PODS)",
			QueryTemplates: []string{
				// 노드의 파드 중 MEMORY 사용량에 따른 내림차순 목록(TOP5)
				"topk(5,sort_desc(sum(container_memory_working_set_bytes{container!=\"\",pod!=\"\",node=~\"%s\"})by(pod)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		NodeFileSystem: {
			Label: "FILE SYSTEM",
			QueryTemplates: []string{
				// 노드의 파일 시스템 사용량(byte)
				"sum(node_filesystem_size_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"}-node_filesystem_avail_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"})",
				// 노드의 총 파일 시스템 크기
				"sum(node_filesystem_size_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"})",
				// 노드의 파일 시스템 사용량(%)
				"sum(node_filesystem_size_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"}-node_filesystem_avail_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"})/sum(node_filesystem_size_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"})*100",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance", "instance"}, false),
				queryGenerator([]interface{}{"instance"}, false),
				queryGenerator([]interface{}{"instance", "instance", "instance"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
				common.BinaryBytes,
				common.Percentage,
			},
			PrimaryUnit: "B",
		},
		NodeFileSystemTop: {
			Label: "FILE SYSTEM",
			QueryTemplates: []string{
				// 노드의 FILE SYSTEM 사용량에 따른 노드 내림차순 목록
				"sort_desc(sum(node_filesystem_size_bytes{instance=~\"%s\",mountpoint=\"/\",fstype!=\"rootfs\"}-node_filesystem_avail_bytes{instance=~\"%s\",mountpoint=\"/\",fstype!=\"rootfs\"})by(instance))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance", "instance"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		NodeFileSystemTop5Projects: {
			Label: "FILE SYSTEM(TOP5 OF PROJECTS)",
			QueryTemplates: []string{
				// 노드의 네임스페이스 중 FILE SYSTEM 사용량에 따른 내림차순 목록(TOP5)
				"topk(5,sort_desc(sum(container_fs_usage_bytes{container!=\"\",pod!=\"\",node=~\"%s\"})by(namespace)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		NodeFileSystemTop5Pods: {
			Label: "FILE SYSTEM(TOP5 OF PODS)",
			QueryTemplates: []string{
				// 노드의 파드 중 FILE SYSTEM 사용량에 따른 내림차순 목록(TOP5)
				"topk(5,sort_desc(sum(container_fs_usage_bytes{container!=\"\",pod!=\"\",node=~\"%s\"})by(pod)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		NodeNetworkIn: {
			Label: "NETWORK IN",
			QueryTemplates: []string{
				// 노드의 NETWORK IN(bps)
				"sum(rate(node_network_receive_bytes_total{instance=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		NodeNetworkInTop: {
			Label: "NETWORK IN",
			QueryTemplates: []string{
				// 노드의 NETWORK IN 에 따른 노드 내림차순 목록
				"sort_desc(sum(rate(node_network_receive_bytes_total{instance=~\"%s\"}[3m]))by(instance))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		NodeNetworkInTop5Projects: {
			Label: "NETWORK IN(TOP5 OF PROJECTS)",
			QueryTemplates: []string{
				// 노드의 네임스페이스 중 NETWORK IN 에 따른 내림차순 목록(TOP5)
				"topk(5,sort_desc(sum(rate(container_network_receive_bytes_total{node=~\"%s\",namespace!=\"\"}[3m]))by(namespace)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		NodeNetworkInTop5Pods: {
			Label: "NETWORK IN(TOP5 OF PODS)",
			QueryTemplates: []string{
				// 노드의 파드 중 NETWORK IN 에 따른 내림차순 목록(TOP5)
				"topk(5,sort_desc(sum(rate(container_network_receive_bytes_total{node=~\"%s\",pod!= \"\"}[3m]))by(pod)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		NodeNetworkOut: {
			Label: "NETWORK OUT",
			QueryTemplates: []string{
				// 노드의 NETWORK OUT(bps)
				"sum(rate(node_network_transmit_bytes_total{instance=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		NodeNetworkOutTop: {
			Label: "NETWORK OUT",
			QueryTemplates: []string{
				// 노드의 NETWORK OUT 에 따른 노드 내림차순 목록
				"sort_desc(sum(rate(node_network_receive_bytes_total{instance=~\"%s\"}[3m]))by(instance))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		NodeNetworkOutTop5Projects: {
			Label: "NETWORK OUT(TOP5 OF PROJECTS)",
			QueryTemplates: []string{
				// 노드의 네임스페이스 중 NETWORK OUT 에 따른 내림차순 목록(TOP5)
				"topk(5,sort_desc(sum(rate(container_network_receive_bytes_total{node=~\"%s\",namespace!=\"\"}[3m]))by(namespace)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		NodeNetworkOutTop5Pods: {
			Label: "NETWORK OUT(TOP5 OF PODS)",
			QueryTemplates: []string{
				// 노드의 파드 중 NETWORK OUT 에 따른 내림차순 목록(TOP5)
				"topk(5,sort_desc(sum(rate(container_network_receive_bytes_total{node=~\"%s\",pod!= \"\"}[3m]))by(pod)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		NodePodCount: {
			Label: "POD COUNT",
			QueryTemplates: []string{
				// 노드의 파드 수
				"sum(kube_pod_info{node=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NodePodCountTop: {
			Label: "POD COUNT",
			QueryTemplates: []string{
				// 노드의 파드 수에 따른 노드 내림차순 목록
				"sort_desc(sum(kube_pod_info{node=~\"%s\",node!=\"\"})by(node))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NodePodCountTop5Projects: {
			Label: "POD COUNT(TOP5 OF PROJECTS)",
			QueryTemplates: []string{
				// 노드의 네임스페이스 중 파드 수에 따른 내림차순 목록(TOP5)
				"topk(5,sort_desc(sum(kube_pod_info{node=~\"%s\"})by(namespace)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}, false),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		QuotaCpuRequest: {
			Label: "CPU REQUEST",
			QueryTemplates: []string{
				// 할당된 CPU REQUEST 쿼터
				"sum(kube_resourcequota{resource=\"requests.cpu\"})",
				// 노드의 CPU Core 수
				"sum(kube_node_status_capacity{resource=\"cpu\",unit=\"core\"})",
				// 할당된 CPU REQUEST 쿼터 할당량(%)
				"sum(kube_resourcequota{resource=\"requests.cpu\"})/sum(kube_node_status_capacity{resource=\"cpu\",unit=\"core\"})*100",
			},
			QueryGenerators: QueryGenerators{
				nil,
				nil,
				nil,
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
				"",
				common.Percentage,
			},
			PrimaryUnit: "Core",
		},
		QuotaCpuLimit: {
			Label: "CPU LIMIT",
			QueryTemplates: []string{
				// 할당된 CPU LIMIT 쿼터
				"sum(kube_resourcequota{resource=\"limits.cpu\"})",
				// 노드의 CPU Core 수
				"sum(kube_node_status_capacity{resource=\"cpu\",unit=\"core\"})",
				// 노드에 할당된 CPU LIMIT 쿼터 할당량(%)
				"sum(kube_resourcequota{resource=\"limits.cpu\"})/sum(kube_node_status_capacity{resource=\"cpu\",unit=\"core\"})*100",
			},
			QueryGenerators: QueryGenerators{
				nil,
				nil,
				nil,
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
				"",
				common.Percentage,
			},
			PrimaryUnit: "Core",
		},
		QuotaMemoryRequest: {
			Label: "MEMORY REQUEST",
			QueryTemplates: []string{
				// 할당된 MEMORY REQUEST 쿼터
				"sum(kube_resourcequota{resource=\"requests.memory\"})",
				// 노드의 총 메모리 크기
				"sum(node_memory_MemTotal_bytes)",
				// 노드에 할당된 MEMORY REQUEST 쿼터 할당량(%)
				"sum(kube_resourcequota{resource=\"requests.memory\"})/sum(node_memory_MemTotal_bytes)*100",
			},
			QueryGenerators: QueryGenerators{
				nil,
				nil,
				nil,
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
				common.BinaryBytes,
				common.Percentage,
			},
			PrimaryUnit: "B",
		},
		QuotaMemoryLimit: {
			Label: "MEMORY LIMIT",
			QueryTemplates: []string{
				// 할당된 MEMORY LIMIT 쿼터
				"sum(kube_resourcequota{resource=\"limits.memory\"})",
				// 노드의 총 메모리 크기
				"sum(node_memory_MemTotal_bytes)",
				// 노드에 할당된 MEMORY LIMIT 쿼터 할당량(%)
				"sum(kube_resourcequota{resource=\"limits.memory\"})/sum(node_memory_MemTotal_bytes)*100",
			},
			QueryGenerators: QueryGenerators{
				nil,
				nil,
				nil,
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
				common.BinaryBytes,
				common.Percentage,
			},
			PrimaryUnit: "B",
		},
		RangeNodeCpuUsage: {
			Label: "CPU USAGE",
			SubLabels: []string{
				"CPU USAGE",
			},
			QueryTemplates: []string{
				"sum(rate(node_cpu_seconds_total{mode!=\"idle\",mode!=\"iowait\",instance=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}, true),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		RangeContainerCpuUsage: {
			Label: "CPU USAGE",
			SubLabels: []string{
				"CPU USAGE",
			},
			QueryTemplates: []string{
				"sum(rate(container_cpu_usage_seconds_total{container!=\"\",pod!=\"\",namespace=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}, true),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		RangeCpuLoadAverage: {
			Label: "CPU LOAD AVERAGE",
			SubLabels: []string{
				"LOAD AVERAGE 1",
				"LOAD AVERAGE 5",
				"LOAD AVERAGE 15",
			},
			QueryTemplates: []string{
				"sum(node_load1{job=\"node-exporter\",instance=~\"%s\"})",
				"sum(node_load5{job=\"node-exporter\",instance=~\"%s\"})",
				"sum(node_load15{job=\"node-exporter\",instance=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}, true),
				queryGenerator([]interface{}{"instance"}, true),
				queryGenerator([]interface{}{"instance"}, true),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
				common.Core,
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		RangeNodeMemoryUsage: {
			Label: "MEMORY USAGE",
			SubLabels: []string{
				"MEMORY USAGE",
			},
			QueryTemplates: []string{
				"sum(node_memory_MemTotal_bytes{instance=~\"%s\"}-node_memory_MemAvailable_bytes{instance=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance", "instance"}, true),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		RangeContainerMemoryUsage: {
			Label: "MEMORY USAGE",
			SubLabels: []string{
				"MEMORY USAGE",
			},
			QueryTemplates: []string{
				"sum(container_memory_working_set_bytes{cluster=\"\",container!=\"\",namespace=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}, true),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		RangeMemorySwap: {
			Label: "MEMORY SWAP",
			SubLabels: []string{
				"MEMORY SWAP",
			},
			QueryTemplates:  []string{},
			QueryGenerators: QueryGenerators{},
			UnitTypeKeys:    []common.UnitTypeKey{},
			PrimaryUnit:     "",
		},
		RangeNodeNetworkIO: {
			Label: "NETWORK IO",
			SubLabels: []string{
				"NETWORK IN",
				"NETWORK OUT",
			},
			QueryTemplates: []string{
				// 컨테이너의 NETWORK IN(bps)
				"sum(rate(node_network_receive_bytes_total{instance=~\"%s\"}[3m]))",
				// 컨테이너의 NETWORK OUT(bps)
				"sum(rate(node_network_transmit_bytes_total{instance=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}, true),
				queryGenerator([]interface{}{"instance"}, true),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		RangeContainerNetworkIO: {
			Label: "NETWORK IO",
			SubLabels: []string{
				"NETWORK IN",
				"NETWORK OUT",
			},
			QueryTemplates: []string{
				// 컨테이너의 NETWORK IN(bps)
				"sum(rate(container_network_receive_bytes_total{container=\"POD\",pod!=\"\",namespace=~\"%s\"}[3m]))",
				// 컨테이너의 NETWORK OUT(bps)
				"sum(rate(container_network_transmit_bytes_total{container=\"POD\",pod!=\"\",namespace=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}, true),
				queryGenerator([]interface{}{"namespace"}, true),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		RangeNodeNetworkPacket: {
			Label: "NETWORK PACKET",
			SubLabels: []string{
				"NETWORK RECEIVE",
				"NETWORK TRANSMIT",
			},
			QueryTemplates: []string{
				// 노드의 NETWORK IN(bps)
				"sum(rate(node_network_receive_packets_total{instance=~\"%s\"}[3m]))",
				// 노드의 NETWORK OUT(bps)
				"sum(rate(node_network_transmit_packets_total{instance=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}, true),
				queryGenerator([]interface{}{"instance"}, true),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.PacketsPerSec,
				common.PacketsPerSec,
			},
			PrimaryUnit: "pps",
		},
		RangeContainerNetworkPacket: {
			Label: "NETWORK PACKET",
			SubLabels: []string{
				"NETWORK RECEIVE",
				"NETWORK TRANSMIT",
			},
			QueryTemplates: []string{
				// 컨테이너의 NETWORK IN(bps)
				"sum(rate(container_network_receive_packets_total{container=\"POD\",pod!=\"\",namespace=~\"%s\"}[3m]))",
				// 컨테이너의 NETWORK OUT(bps)
				"sum(rate(container_network_transmit_packets_total{container=\"POD\",pod!=\"\",namespace=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}, true),
				queryGenerator([]interface{}{"namespace"}, true),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.PacketsPerSec,
				common.PacketsPerSec,
			},
			PrimaryUnit: "pps",
		},
		RangeNetworkBandwidth: {
			Label: "NETWORK BANDWIDTH",
			SubLabels: []string{
				"NETWORK BANDWIDTH",
			},
			QueryTemplates:  []string{},
			QueryGenerators: QueryGenerators{},
			UnitTypeKeys:    []common.UnitTypeKey{},
			PrimaryUnit:     "",
		},
		RangeNetworkPacketReceiveTransmit: {
			Label: "NETWORK PACKET RECEIVE/TRANSMIT",
			SubLabels: []string{
				"NETWORK PACKET RECEIVE",
				"NETWORK PACKET TRANSMIT",
			},
			QueryTemplates:  []string{},
			QueryGenerators: QueryGenerators{},
			UnitTypeKeys:    []common.UnitTypeKey{},
			PrimaryUnit:     "",
		},
		RangeNetworkPacketReceiveTransmitDrop: {
			Label: "NETWORK PACKET RECEIVE/TRANSMIT DROP",
			SubLabels: []string{
				"NETWORK PACKET RECEIVE DROP",
				"NETWORK PACKET TRANSMIT DROP",
			},
			QueryTemplates:  []string{},
			QueryGenerators: QueryGenerators{},
			UnitTypeKeys:    []common.UnitTypeKey{},
			PrimaryUnit:     "",
		},
		RangeFileSystem: {
			Label: "FILE SYSTEM",
			SubLabels: []string{
				"FILE SYSTEM",
			},
			QueryTemplates:  []string{},
			QueryGenerators: QueryGenerators{},
			UnitTypeKeys:    []common.UnitTypeKey{},
			PrimaryUnit:     "",
		},
		RangeDiskIO: {
			Label: "DISK IO",
			SubLabels: []string{
				"DISK IN",
				"DISK OUT",
			},
			QueryTemplates:  []string{},
			QueryGenerators: QueryGenerators{},
			UnitTypeKeys:    []common.UnitTypeKey{},
			PrimaryUnit:     "",
		},
	}
)
