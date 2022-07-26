package prometheus

import (
	"go-practice/common"
)

// QueryGenerators 쿼리 템플릿과 쿼리 파라미터를 인자로 받아 쿼리를 반환하는 함수 목록 타입
type QueryGenerators []func(queryTemplate string, queryParams map[string]interface{}) (string, string)

// QueryTemplateParsers 쿼리 템플릿과 쿼리 파라미터를 인자로 받아 쿼리를 반환하는 함수 목록 타입
type QueryTemplateParsers []func(queryTemplate string, queryParams map[string]interface{}) (string, string)

// MetricDefinition 메트릭 정의 구조체
type MetricDefinition struct {
	Label           string                          // 메트릭의 라벨
	SubLabels       []string                        // 쿼리 템플릿의 라벨(미필수)
	QueryInfos      map[PrometheusVersion]QueryInfo // 버전별 쿼리 모음
	QueryTemplates  []string                        // 쿼리 템플릿
	QueryGenerators QueryGenerators                 // 쿼리 템플릿에 조건절 추가하여 쿼리를 반환하는 함수 목록(쿼리 템플릿과 맵핑)
	UnitTypeKeys    []common.UnitTypeKey            // 쿼리 결과값의 단위 타입의 키 목록(쿼리 템플릿과 맵핑)
	PrimaryUnit     string                          // 쿼리 결과값의 단위 중 주단위

	MetricKeys []MetricKey // 다른 메트릭 정의를 활용하는 메트릭(다른 메트릭 활용 시 해당 값만 작성)
}

// MetricKey 메트릭 키
type PrometheusVersion string

const (
	v2_20_0 = PrometheusVersion("2.20.0")
	v2_26_0 = PrometheusVersion("2.26.0")
	v2_29_0 = PrometheusVersion("2.29.0")
	v2_32_0 = PrometheusVersion("2.32.0")
)

type QueryInfo struct {
	ReferenceVersion     PrometheusVersion
	QueryTemplates       []string
	QueryTemplateParsers QueryTemplateParsers // 쿼리 템플릿에 조건절 추가하여 쿼리를 반환하는 함수 목록(쿼리 템플릿과 맵핑)
}

// MetricDefinitions 메트릭 키에 따른 메트릭 정의 상수
var (
	MetricDefinitions = map[MetricKey]MetricDefinition{
		ContainerCpu: {
			Label: "CPU USAGE",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(rate(container_cpu_usage_seconds_total{container!=\"\",pod!=\"\",namespace=~\"%s\"}[3m]))",
					},
					QueryTemplateParsers: QueryTemplateParsers{
						queryTemplateParser([]interface{}{"namespace"}),
					},
				},
				v2_26_0: {
					ReferenceVersion: v2_20_0,
				},
			},
			QueryTemplates: []string{
				// 컨테이너의 CPU Core 사용량(Core)
				"sum(rate(container_cpu_usage_seconds_total{container!=\"\",pod!=\"\",namespace=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		ContainerDiskIOReads: {
			Label: "DISK IO READS",
			QueryTemplates: []string{
				// 컨테이너의 읽기 DISK IO
				"sum(irate(container_fs_reads_bytes_total{device!=\"\",node=~\"%s\",namespace=~\"%s\",pod=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node", "namespace", "pod"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		ContainerDiskIOWrites: {
			Label: "DISK IO WRITES",
			QueryTemplates: []string{
				// 컨테이너의 쓰기 DISK IO
				"sum(irate(container_fs_writes_bytes_total{device!=\"\",node=~\"%s\",namespace=~\"%s\",pod=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node", "namespace", "pod"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		ContainerFileSystem: {
			Label: "FILE SYSTEM",
			QueryTemplates: []string{
				// 컨테이너의 파일 시스템 사용량(byte)
				"sum(container_fs_usage_bytes{namespace=~\"%s\",pod=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace", "pod"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		ContainerMemory: {
			Label: "MEMORY USAGE",
			QueryTemplates: []string{
				// 컨테이너의 메모리 사용량(byte)
				"sum(container_memory_working_set_bytes{cluster=\"\",container!=\"\",namespace=~\"%s\",pod=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace", "pod"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		ContainerNetworkIn: {
			Label: "NETWORK IN",
			QueryTemplates: []string{
				// 컨테이너의 NETWORK IN(bps)
				"sum(rate(container_network_receive_bytes_total{container=\"POD\",pod!=\"\",namespace=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		ContainerNetworkIO: {
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
				queryGenerator([]interface{}{"namespace"}),
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		ContainerNetworkOut: {
			Label: "NETWORK OUT",
			QueryTemplates: []string{
				// 컨테이너의 NETWORK OUT(bps)
				"sum(rate(container_network_transmit_bytes_total{container=\"POD\",namespace=~\"%s\",pod!=\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		ContainerNetworkPacket: {
			Label: "NETWORK PACKET",
			SubLabels: []string{
				"NETWORK RECEIVE",
				"NETWORK TRANSMIT",
			},
			QueryTemplates: []string{
				"sum(rate(container_network_receive_packets_total{container=\"POD\",pod!=\"\",namespace=~\"%s\"}[3m]))",
				"sum(rate(container_network_transmit_packets_total{container=\"POD\",pod!=\"\",namespace=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Numeric,
				common.Numeric,
			},
			PrimaryUnit: "",
		},
		ContainerNetworkPacketDrop: {
			Label: "NETWORK PACKET DROP",
			SubLabels: []string{
				"NETWORK RECEIVE DROP",
				"NETWORK TRANSMIT DROP",
			},
			QueryTemplates: []string{
				// 컨테이너의 드롭된 수신 패킷
				"sum(rate(container_network_receive_packets_dropped_total{node=~\"%s\",namespace=~\"%s\",pod=~\"%s\"}[3m]))",
				// 컨테이너의 드롭된 전송 패킷
				"sum(rate(container_network_transmit_packets_dropped_total{node=~\"%s\",namespace=~\"%s\",pod=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node", "namespace", "pod"}),
				queryGenerator([]interface{}{"node", "namespace", "pod"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Numeric,
				common.Numeric,
			},
			PrimaryUnit: "rps",
		},
		CustomQuotaCpuLimit: {
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
		CustomQuotaCpuRequest: {
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
		CustomQuotaMemoryLimit: {
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
		CustomQuotaMemoryRequest: {
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
		CustomNodeCpu: {
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
				queryGenerator([]interface{}{"instance"}),
				queryGenerator([]interface{}{"node"}),
				queryGenerator([]interface{}{"instance", "node"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
				"",
				common.Percentage,
			},
			PrimaryUnit: "Core",
		},
		CustomNodeFileSystem: {
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
				queryGenerator([]interface{}{"instance", "instance"}),
				queryGenerator([]interface{}{"instance"}),
				queryGenerator([]interface{}{"instance", "instance", "instance"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
				common.BinaryBytes,
				common.Percentage,
			},
			PrimaryUnit: "B",
		},
		CustomNodeMemory: {
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
				queryGenerator([]interface{}{"instance", "instance"}),
				queryGenerator([]interface{}{"instance"}),
				queryGenerator([]interface{}{"instance", "instance", "instance"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
				common.BinaryBytes,
				common.Percentage,
			},
			PrimaryUnit: "B",
		},
		NumberOfContainer: {
			Label: "CONTAINER",
			QueryTemplates: []string{
				// 파드의 컨테이너 수
				"sum(kube_pod_container_info{pod=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"pod"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NumberOfDeployment: {
			Label: "DEPLOYMENT",
			QueryTemplates: []string{
				"count(kube_deployment_labels{namespace=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NumberOfIngress: {
			Label: "INGRESS",
			QueryTemplates: []string{
				"count(kube_ingress_labels{namespace=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NumberOfPipeline: {
			Label: "PIPELINE",
		},
		NumberOfPod: {
			Label: "POD",
			QueryTemplates: []string{
				// 노드, 네임스페이스 파드 수
				"count(kube_pod_info{node=~\"%s\",namespace=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node", "namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NumberOfNamespace: {
			Label: "PROJECT",
			QueryTemplates: []string{
				"count(kube_namespace_status_phase{phase=\"Active\",namespace=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NumberOfService: {
			Label: "SERVICE",
			QueryTemplates: []string{
				"count(kube_service_labels{namespace=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NumberOfStatefulSet: {
			Label: "STATEFULSET",
			QueryTemplates: []string{
				"count(kube_statefulset_labels{namespace=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NumberOfVolume: {
			Label: "VOLUME",
			QueryTemplates: []string{
				"count(kube_persistentvolume_labels{namespace=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NodeCpu: {
			Label: "CPU USAGE",
			QueryTemplates: []string{
				"sum(rate(node_cpu_seconds_total{mode!=\"idle\",mode!=\"iowait\",instance=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		NodeCpuLoadAverage: {
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
				queryGenerator([]interface{}{"instance"}),
				queryGenerator([]interface{}{"instance"}),
				queryGenerator([]interface{}{"instance"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
				common.Core,
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		NodeDiskIO: {
			Label: "DISK IO",
			QueryTemplates: []string{
				"sum(rate(node_disk_io_time_weighted_seconds_total{device=~\"nvme.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\",job=\"node-exporter\",instance=~\"%s\"}[1m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}),
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
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance", "instance"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		NodeMemory: {
			Label: "MEMORY USAGE",
			QueryTemplates: []string{
				"sum(node_memory_MemTotal_bytes{instance=~\"%s\"}-node_memory_MemAvailable_bytes{instance=~\"%s\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance", "instance"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		NodeNetworkIO: {
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
				queryGenerator([]interface{}{"instance"}),
				queryGenerator([]interface{}{"instance"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		NodeNetworkIn: {
			Label: "NETWORK IN",
			QueryTemplates: []string{
				// 노드의 NETWORK IN(bps)
				"sum(rate(node_network_receive_bytes_total{instance=~\"%s\"}[3m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}),
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
				queryGenerator([]interface{}{"instance"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		NodeNetworkPacket: {
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
				queryGenerator([]interface{}{"instance"}),
				queryGenerator([]interface{}{"instance"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.PacketsPerSec,
				common.PacketsPerSec,
			},
			PrimaryUnit: "pps",
		},
		NodeNetworkPacketDrop: {
			Label: "NETWORK PACKET DROP",
			SubLabels: []string{
				"NETWORK RECEIVE DROP",
				"NETWORK TRANSMIT DROP",
			},
			QueryTemplates: []string{
				// 노드의 드롭된 수신 패킷
				"sum(rate(node_network_receive_drop_total{device!=\"lo\",job=\"node-exporter\",instance=~\"%s\"}[1m]))",
				// 노드의 드롭된 전송 패킷
				"sum(rate(node_network_transmit_drop_excluding_lo{device!=\"lo\",job=\"node-exporter\",instance=~\"%s\"}[1m]))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}),
				queryGenerator([]interface{}{"instance"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Numeric,
				common.Numeric,
			},
			PrimaryUnit: "rps",
		},
		TopNodeCpuByInstance: {
			Label: "CPU",
			QueryTemplates: []string{
				// 노드의 CPU 사용량에 따른 노드 내림차순 목록
				"sort_desc(sum(rate(node_cpu_seconds_total{mode!=\"idle\",mode!=\"iowait\",instance=~\"%s\"}[3m]))by(instance))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		TopNodeFileSystemByInstance: {
			Label: "FILE SYSTEM",
			QueryTemplates: []string{
				// 노드의 FILE SYSTEM 사용량에 따른 노드 내림차순 목록
				"sort_desc(sum(node_filesystem_size_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"}-node_filesystem_avail_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"})by(instance))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance", "instance"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		TopNodeMemoryByInstance: {
			Label: "MEMORY",
			QueryTemplates: []string{
				// 노드의 MEMORY 사용량에 따른 노드 내림차순 목록
				"sort_desc(sum(node_memory_MemTotal_bytes-node_memory_MemAvailable_bytes{instance=~\"%s\"})by(instance))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		TopNodeNetworkInByInstance: {
			Label: "NETWORK IN",
			QueryTemplates: []string{
				// 노드의 NETWORK IN 에 따른 노드 내림차순 목록
				"sort_desc(sum(rate(node_network_receive_bytes_total{instance=~\"%s\"}[3m]))by(instance))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		TopNodeNetworkOutByInstance: {
			Label: "NETWORK OUT",
			QueryTemplates: []string{
				// 노드의 NETWORK OUT 에 따른 노드 내림차순 목록
				"sort_desc(sum(rate(node_network_receive_bytes_total{instance=~\"%s\"}[3m]))by(instance))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"instance"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		TopCountPodByNode: {
			Label: "POD COUNT",
			QueryTemplates: []string{
				// 노드별 파드 수에 따른 내림차순 목록
				"sort_desc(count(kube_pod_info{node!=\"\",node=~\"%s\"})by(node))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		Top5ContainerCpuByNamespace: {
			Label: "CPU(TOP5 OF PROJECTS)",
			QueryTemplates: []string{
				// 네임스페이스별 CPU 사용량에 따른 내림차순 목록
				"topk(5,sort_desc(sum(rate(container_cpu_usage_seconds_total{container!=\"\",pod!=\"\",node=~\"%s\"}[3m]))by(namespace)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		Top5ContainerCpuByPod: {
			Label: "CPU(TOP5 OF PODS)",
			QueryTemplates: []string{
				// 파드별 CPU 사용량에 따른 내림차순 목록
				"topk(5,sort_desc(sum(rate(container_cpu_usage_seconds_total{container!=\"\",pod!=\"\",node=~\"%s\",namespace=~\"%s\"}[3m]))by(pod)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node", "namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		Top5ContainerFileSystemByNamespace: {
			Label: "FILE SYSTEM(TOP5 OF PROJECTS)",
			QueryTemplates: []string{
				// 노드의 네임스페이스 중 FILE SYSTEM 사용량에 따른 내림차순 목록
				"topk(5,sort_desc(sum(container_fs_usage_bytes{container!=\"\",pod!=\"\",node=~\"%s\"})by(namespace)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		Top5ContainerFileSystemByPod: {
			Label: "FILE SYSTEM(TOP5 OF PODS)",
			QueryTemplates: []string{
				// 파드별 FILE SYSTEM 사용량에 따른 내림차순 목록
				"topk(5,sort_desc(sum(container_fs_usage_bytes{container!=\"\",pod!=\"\",node=~\"%s\",namespace=~\"%s\"})by(pod)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node", "namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		Top5ContainerMemoryByNamespace: {
			Label: "MEMORY(TOP5 OF PROJECTS)",
			QueryTemplates: []string{
				// 네임스페이스별 MEMORY 사용량에 따른 내림차순 목록
				"topk(5,sort_desc(sum(container_memory_working_set_bytes{container!=\"\",pod!=\"\",node=~\"%s\"})by(namespace)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		Top5ContainerMemoryByPod: {
			Label: "MEMORY(TOP5 OF PODS)",
			QueryTemplates: []string{
				// 파드별 MEMORY 사용량에 따른 내림차순 목록
				"topk(5,sort_desc(sum(container_memory_working_set_bytes{container!=\"\",pod!=\"\",node=~\"%s\",namespace=~\"%s\"})by(pod)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node", "namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		Top5ContainerNetworkInByNamespace: {
			Label: "NETWORK IN(TOP5 OF PROJECTS)",
			QueryTemplates: []string{
				// NETWORK IN 에 따른 Top 5 내림차순 목록
				"topk(5,sort_desc(sum(rate(container_network_receive_bytes_total{container=\"POD\",pod!=\"\",node=~\"%s\",namespace=~\"%s\"}[3m]))by(namespace)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node", "namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		Top5ContainerNetworkInByPod: {
			Label: "NETWORK IN(TOP5 OF PODS)",
			QueryTemplates: []string{
				// NETWORK IN 에 따른 Top 5 내림차순 목록
				"topk(5,sort_desc(sum(rate(container_network_receive_bytes_total{container=\"POD\",pod!=\"\",node=~\"%s\",namespace=~\"%s\"}[3m]))by(pod)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node", "namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		Top5ContainerNetworkOutByNamespace: {
			Label: "NETWORK OUT(TOP5 OF PROJECTS)",
			QueryTemplates: []string{
				// 노드의 네임스페이스 중 NETWORK OUT 에 따른 내림차순 목록
				"topk(5,sort_desc(sum(rate(container_network_receive_bytes_total{namespace!=\"\",node=~\"%s\"}[3m]))by(namespace)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		Top5ContainerNetworkOutByPod: {
			Label: "NETWORK OUT(TOP5 OF PODS)",
			QueryTemplates: []string{
				// 파드 중 NETWORK OUT 에 따른 내림차순 목록
				"topk(5,sort_desc(sum(rate(container_network_receive_bytes_total{pod!= \"\",node=~\"%s\",namespace=~\"%s\"}[3m]))by(pod)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node", "namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		Top5CountPodByNamespace: {
			Label: "POD COUNT(TOP5 OF PROJECTS)",
			QueryTemplates: []string{
				// 노드의 네임스페이스 중 파드 수에 따른 내림차순 목록
				"topk(5,sort_desc(count(kube_pod_info{node=~\"%s\",namespace=~\"%s\"})by(namespace)))",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"node", "namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		QuotaCpuLimit: {
			Label: "CPU LIMIT",
			QueryTemplates: []string{
				"sum(kube_resourcequota{resource=\"limits.cpu\"})",
			},
			QueryGenerators: QueryGenerators{
				nil,
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCpuRequest: {
			Label: "CPU REQUEST",
			QueryTemplates: []string{
				"sum(kube_resourcequota{resource=\"requests.cpu\"})",
			},
			QueryGenerators: QueryGenerators{
				nil,
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaMemoryRequest: {
			Label: "MEMORY REQUEST",
			QueryTemplates: []string{
				"sum(kube_resourcequota{resource=\"requests.memory\"})",
			},
			QueryGenerators: QueryGenerators{
				nil,
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		QuotaMemoryLimit: {
			Label: "MEMORY LIMIT",
			QueryTemplates: []string{
				"sum(kube_resourcequota{resource=\"limits.memory\"})",
			},
			QueryGenerators: QueryGenerators{
				nil,
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		QuotaObjectCountConfigmaps: {
			Label: "OBJECT COUNT CONFIGMAPS",
			QueryTemplates: []string{
				"sum(kube_resourcequota{resource=~\".*configmaps\",namespace=~\"%s\",type=\"hard\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaObjectCountPods: {
			Label: "OBJECT COUNT PODS",
			QueryTemplates: []string{
				"sum(kube_resourcequota{resource=~\".*pods\",namespace=~\"%s\",type=\"hard\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaObjectCountSecrets: {
			Label: "OBJECT COUNT SECRETS",
			QueryTemplates: []string{
				"sum(kube_resourcequota{resource=~\".*secrets\",namespace=~\"%s\",type=\"hard\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaObjectCountReplicationControllers: {
			Label: "OBJECT COUNT REPLICATION CONTROLLERS",
			QueryTemplates: []string{
				"sum(kube_resourcequota{resource=~\".*replicationcontrollers\",namespace=~\"%s\",type=\"hard\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaObjectCountServices: {
			Label: "OBJECT COUNT SERVICES",
			QueryTemplates: []string{
				"sum(kube_resourcequota{resource=~\".*services\",namespace=~\"%s\",type=\"hard\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaObjectCountServicesLoadBalancers: {
			Label: "OBJECT COUNT SERVICES LOAD BALANCERS",
			QueryTemplates: []string{
				"sum(kube_resourcequota{resource=~\".*services.loadbalancers\",namespace=~\"%s\",type=\"hard\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaObjectCountServicesNodePorts: {
			Label: "OBJECT COUNT SERVICES NODE PORTS",
			QueryTemplates: []string{
				"sum(kube_resourcequota{resource=~\".*services.nodeports\",namespace=~\"%s\",type=\"hard\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaObjectCountResourceQuotas: {
			Label: "OBJECT COUNT RESOURCE QUOTAS",
			QueryTemplates: []string{
				"sum(kube_resourcequota{resource=~\".*resourcequotas\",namespace=~\"%s\",type=\"hard\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaObjectCountPersistentVolumeClaims: {
			Label: "OBJECT COUNT PERSISTENT VOLUME CLAIMS",
			QueryTemplates: []string{
				"sum(kube_resourcequota{resource=~\".*persistentvolumeclaims\",namespace=~\"%s\",type=\"hard\"})",
			},
			QueryGenerators: QueryGenerators{
				queryGenerator([]interface{}{"namespace"}),
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		SummaryNodeInfo: {
			MetricKeys: []MetricKey{CustomNodeCpu, CustomNodeFileSystem, CustomNodeMemory, NodeNetworkIn, NodeNetworkOut, NumberOfPod},
		},
	}
)
