package prometheus

import "go-practice/common"

// QueryTemplateParserGenerators 쿼리 템플릿과 쿼리 파라미터를 인자로 받아 쿼리를 반환하는 함수 목록 타입
type QueryTemplateParserGenerators []func(queryTemplate string, bodyParams map[string]interface{}) (string, string)

// MetricKey 메트릭 키
type PrometheusVersion string

const (
	v2_20_0 = PrometheusVersion("2.20.0")
	v2_26_0 = PrometheusVersion("2.26.0")
)

type QueryInfo struct {
	ReferenceVersion              PrometheusVersion             // 쿼리를 참조하기 위한 버전, 참조 버전의 쿼리를 사용
	QueryTemplates                []string                      // 쿼리 템플릿
	QueryTemplateParserGenerators QueryTemplateParserGenerators // 쿼리 템플릿에 조건절 추가하여 쿼리를 반환하는 함수 목록(쿼리 템플릿과 맵핑)
}

// MetricDefinition 메트릭 정의 구조체
type MetricDefinition struct {
	Label        string                          // 메트릭의 라벨
	SubLabels    []string                        // 쿼리 템플릿의 라벨
	QueryInfos   map[PrometheusVersion]QueryInfo // 버전별 쿼리 모음
	UnitTypeKeys []common.UnitTypeKey            // 쿼리 결과값의 단위 타입의 키 목록(쿼리 템플릿과 맵핑)
	PrimaryUnit  string                          // 쿼리 결과값의 단위 중 주단위
	MetricKeys   []MetricKey                     // 다른 메트릭 정의를 활용하는 메트릭(다른 메트릭 활용 시 해당 값만 작성)
}

// MetricDefinitions 메트릭 키에 따른 메트릭 정의 상수
var (
	MetricDefinitions = map[MetricKey]MetricDefinition{
		ContainerCpu: {
			Label: "CPU",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 컨테이너의 CPU Core 사용량(Core)
						"sum(rate(container_cpu_usage_seconds_total{container!=\"\",pod!=\"\",namespace=~\"%s\",pod=~\"%s\"}[3m]))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace", "pod"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		ContainerDiskIORead: {
			Label: "DISK READS",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 컨테이너의 읽기 DISK IO
						"sum(irate(container_fs_reads_bytes_total{device!=\"\",node=~\"%s\",namespace=~\"%s\",pod=~\"%s\"}[3m]))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node", "namespace", "pod"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		ContainerDiskIOWrite: {
			Label: "DISK WRITES",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 컨테이너의 쓰기 DISK IO
						"sum(irate(container_fs_writes_bytes_total{device!=\"\",node=~\"%s\",namespace=~\"%s\",pod=~\"%s\"}[3m]))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node", "namespace", "pod"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		ContainerFileSystem: {
			Label: "FILE SYSTEM",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 컨테이너의 파일 시스템 사용량(byte)
						"sum(container_fs_usage_bytes{namespace=~\"%s\",pod=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace", "pod"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		ContainerMemory: {
			Label: "MEMORY",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 컨테이너의 메모리 사용량(byte)
						"sum(container_memory_working_set_bytes{cluster=\"\",container!=\"\",namespace=~\"%s\",pod=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace", "pod"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		ContainerNetworkIn: {
			Label: "NETWORK IN",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 컨테이너의 NETWORK IN(bps)
						"sum(rate(container_network_receive_bytes_total{container=\"POD\",pod!=\"\",namespace=~\"%s\",pod=~\"%s\"}[3m]))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace", "pod"}),
					},
				},
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
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 컨테이너의 NETWORK IN(bps)
						"sum(rate(container_network_receive_bytes_total{container=\"POD\",pod!=\"\",namespace=~\"%s\",pod=~\"%s\"}[3m]))",
						// 컨테이너의 NETWORK OUT(bps)
						"sum(rate(container_network_transmit_bytes_total{container=\"POD\",pod!=\"\",namespace=~\"%s\",pod=~\"%s\"}[3m]))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace", "pod"}),
						queryTemplateParserGenerator([]interface{}{"namespace", "pod"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		ContainerNetworkOut: {
			Label: "NETWORK OUT",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 컨테이너의 NETWORK OUT(bps)
						"sum(rate(container_network_transmit_bytes_total{container=\"POD\",pod!=\"\",namespace=~\"%s\",pod=~\"%s\"}[3m]))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace", "pod"}),
					},
				},
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
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(rate(container_network_receive_packets_total{container=\"POD\",pod!=\"\",namespace=~\"%s\",pod=~\"%s\"}[3m]))",
						"sum(rate(container_network_transmit_packets_total{container=\"POD\",pod!=\"\",namespace=~\"%s\",pod=~\"%s\"}[3m]))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace", "pod"}),
						queryTemplateParserGenerator([]interface{}{"namespace", "pod"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.PacketsPerSec,
				common.PacketsPerSec,
			},
			PrimaryUnit: "pps",
		},
		ContainerNetworkPacketDrop: {
			Label: "NETWORK PACKET DROP",
			SubLabels: []string{
				"NETWORK RECEIVE DROP",
				"NETWORK TRANSMIT DROP",
			},
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 컨테이너의 드롭된 수신 패킷
						"sum(rate(container_network_receive_packets_dropped_total{node=~\"%s\",namespace=~\"%s\",pod=~\"%s\"}[3m]))",
						// 컨테이너의 드롭된 전송 패킷
						"sum(rate(container_network_transmit_packets_dropped_total{node=~\"%s\",namespace=~\"%s\",pod=~\"%s\"}[3m]))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node", "namespace", "pod"}),
						queryTemplateParserGenerator([]interface{}{"node", "namespace", "pod"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Numeric,
				common.Numeric,
			},
			PrimaryUnit: "rps",
		},
		CustomContainerVolume: {
			Label: "PERSISTENT VOLUME",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kubelet_volume_stats_used_bytes{node=~\"%s\"})",
						"sum(kubelet_volume_stats_capacity_bytes{node=~\"%s\"})",
						"sum(kubelet_volume_stats_used_bytes{node=~\"%s\"})/sum(kubelet_volume_stats_capacity_bytes{node=~\"%s\"})*100",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node"}),
						queryTemplateParserGenerator([]interface{}{"node"}),
						queryTemplateParserGenerator([]interface{}{"node", "node"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
				common.BinaryBytes,
				common.Percentage,
			},
			PrimaryUnit: "B",
		},
		CustomQuotaLimitCpu: {
			Label: "CPU LIMIT",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 할당된 CPU LIMIT 쿼터
						"sum(kube_resourcequota{resource=\"limits.cpu\"})",
						// 노드의 CPU Core 수
						"sum(kube_node_status_capacity{resource=\"cpu\",unit=\"core\"})",
						// 노드에 할당된 CPU LIMIT 쿼터 할당량(%)
						"sum(kube_resourcequota{resource=\"limits.cpu\"})/sum(kube_node_status_capacity{resource=\"cpu\",unit=\"core\"})*100",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						nil,
						nil,
						nil,
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
				"",
				common.Percentage,
			},
			PrimaryUnit: "Core",
		},
		CustomQuotaLimitMemory: {
			Label: "MEMORY LIMIT",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 할당된 MEMORY LIMIT 쿼터
						"sum(kube_resourcequota{resource=\"limits.memory\"})",
						// 노드의 총 메모리 크기
						"sum(node_memory_MemTotal_bytes)",
						// 노드에 할당된 MEMORY LIMIT 쿼터 할당량(%)
						"sum(kube_resourcequota{resource=\"limits.memory\"})/sum(node_memory_MemTotal_bytes)*100",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						nil,
						nil,
						nil,
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
				common.BinaryBytes,
				common.Percentage,
			},
			PrimaryUnit: "B",
		},
		CustomQuotaRequestCpu: {
			Label: "CPU REQUEST",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 할당된 CPU REQUEST 쿼터
						"sum(kube_resourcequota{resource=\"requests.cpu\"})",
						// 노드의 CPU Core 수
						"sum(kube_node_status_capacity{resource=\"cpu\",unit=\"core\"})",
						// 할당된 CPU REQUEST 쿼터 할당량(%)
						"sum(kube_resourcequota{resource=\"requests.cpu\"})/sum(kube_node_status_capacity{resource=\"cpu\",unit=\"core\"})*100",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						nil,
						nil,
						nil,
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
				"",
				common.Percentage,
			},
			PrimaryUnit: "Core",
		},
		CustomQuotaRequestMemory: {
			Label: "MEMORY REQUEST",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 할당된 MEMORY REQUEST 쿼터
						"sum(kube_resourcequota{resource=\"requests.memory\"})",
						// 노드의 총 메모리 크기
						"sum(node_memory_MemTotal_bytes)",
						// 노드에 할당된 MEMORY REQUEST 쿼터 할당량(%)
						"sum(kube_resourcequota{resource=\"requests.memory\"})/sum(node_memory_MemTotal_bytes)*100",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						nil,
						nil,
						nil,
					},
				},
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
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드의 CPU Core 사용량(Core)
						"sum(rate(node_cpu_seconds_total{mode!=\"idle\",mode!=\"iowait\",instance=~\"%s\"}[3m]))",
						// 노드의 CPU Core 수
						"sum(kube_node_status_capacity{resource=\"cpu\",unit=\"core\",node=~\"%s\"})",
						// 총 CPU Core 사용량(%)
						"sum(rate(node_cpu_seconds_total{mode!=\"idle\",mode!=\"iowait\",instance=~\"%s\"}[3m]))/sum(kube_node_status_capacity{resource=\"cpu\",unit=\"core\",node=~\"%s\"})*100",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance"}),
						queryTemplateParserGenerator([]interface{}{"node"}),
						queryTemplateParserGenerator([]interface{}{"instance", "node"}),
					},
				},
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
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드의 파일 시스템 사용량(byte)
						"sum(node_filesystem_size_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"}-node_filesystem_avail_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"})",
						// 노드의 총 파일 시스템 크기
						"sum(node_filesystem_size_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"})",
						// 노드의 파일 시스템 사용량(%)
						"sum(node_filesystem_size_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"}-node_filesystem_avail_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"})/sum(node_filesystem_size_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"})*100",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance", "instance"}),
						queryTemplateParserGenerator([]interface{}{"instance"}),
						queryTemplateParserGenerator([]interface{}{"instance", "instance", "instance"}),
					},
				},
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
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드의 메모리 사용량(byte)
						"sum(node_memory_MemTotal_bytes{instance=~\"%s\"}-node_memory_MemAvailable_bytes{instance=~\"%s\"})",
						// 노드의 총 메모리 크기
						"sum(node_memory_MemTotal_bytes{instance=~\"%s\"})",
						// 노드의 메모리 사용량(%)
						"sum(node_memory_MemTotal_bytes{instance=~\"%s\"}-node_memory_MemAvailable_bytes{instance=~\"%s\"})/sum(node_memory_MemTotal_bytes{instance=~\"%s\"})*100",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance", "instance"}),
						queryTemplateParserGenerator([]interface{}{"instance"}),
						queryTemplateParserGenerator([]interface{}{"instance", "instance", "instance"}),
					},
				},
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
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 파드의 컨테이너 수
						"sum(kube_pod_container_info{pod=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"pod"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NumberOfDeployment: {
			Label: "DEPLOYMENT",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"count(kube_deployment_labels{namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NumberOfIngress: {
			Label: "INGRESS",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"count(kube_ingress_labels{namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
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
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드, 네임스페이스 파드 수
						"count(kube_pod_info{node=~\"%s\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node", "namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NumberOfNamespace: {
			Label: "PROJECT",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"count(kube_namespace_status_phase{phase=\"Active\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NumberOfService: {
			Label: "SERVICE",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"count(kube_service_labels{namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NumberOfStatefulSet: {
			Label: "STATEFULSET",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"count(kube_statefulset_labels{namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NumberOfVolume: {
			Label: "VOLUME",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"count(kube_persistentvolume_labels{namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		NodeCpu: {
			Label: "CPU",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(rate(node_cpu_seconds_total{mode!=\"idle\",mode!=\"iowait\",instance=~\"%s\"}[3m]))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		NodeCpuLoadAverage: {
			Label: "CPU LOAD AVERAGE",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(node_load1{job=\"node-exporter\",instance=~\"%s\"})",
						"sum(node_load5{job=\"node-exporter\",instance=~\"%s\"})",
						"sum(node_load15{job=\"node-exporter\",instance=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance"}),
						queryTemplateParserGenerator([]interface{}{"instance"}),
						queryTemplateParserGenerator([]interface{}{"instance"}),
					},
				},
			},
			SubLabels: []string{
				"LOAD AVERAGE 1",
				"LOAD AVERAGE 5",
				"LOAD AVERAGE 15",
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
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(rate(node_disk_io_time_weighted_seconds_total{device=~\"nvme.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\",job=\"node-exporter\",instance=~\"%s\"}[1m]))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		NodeFileSystem: {
			Label: "FILE SYSTEM",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드의 파일 시스템 사용량(byte)
						"sum(node_filesystem_size_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"}-node_filesystem_avail_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance", "instance"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		NodeMemory: {
			Label: "MEMORY",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(node_memory_MemTotal_bytes{instance=~\"%s\"}-node_memory_MemAvailable_bytes{instance=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance", "instance"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		NodeNetworkIO: {
			Label: "NETWORK IO",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 컨테이너의 NETWORK IN(bps)
						"sum(rate(node_network_receive_bytes_total{instance=~\"%s\"}[3m]))",
						// 컨테이너의 NETWORK OUT(bps)
						"sum(rate(node_network_transmit_bytes_total{instance=~\"%s\"}[3m]))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance"}),
						queryTemplateParserGenerator([]interface{}{"instance"}),
					},
				},
			},
			SubLabels: []string{
				"NETWORK IN",
				"NETWORK OUT",
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		NodeNetworkIn: {
			Label: "NETWORK IN",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드의 NETWORK IN(bps)
						"sum(rate(node_network_receive_bytes_total{instance=~\"%s\"}[3m]))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		NodeNetworkOut: {
			Label: "NETWORK OUT",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드의 NETWORK OUT(bps)
						"sum(rate(node_network_transmit_bytes_total{instance=~\"%s\"}[3m]))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance"}),
					},
				},
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
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드의 NETWORK IN(bps)
						"sum(rate(node_network_receive_packets_total{instance=~\"%s\"}[3m]))",
						// 노드의 NETWORK OUT(bps)
						"sum(rate(node_network_transmit_packets_total{instance=~\"%s\"}[3m]))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance"}),
						queryTemplateParserGenerator([]interface{}{"instance"}),
					},
				},
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
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드의 드롭된 수신 패킷
						"sum(rate(node_network_receive_drop_total{device!=\"lo\",job=\"node-exporter\",instance=~\"%s\"}[1m]))",
						// 노드의 드롭된 전송 패킷
						"sum(rate(node_network_transmit_drop_excluding_lo{device!=\"lo\",job=\"node-exporter\",instance=~\"%s\"}[1m]))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance"}),
						queryTemplateParserGenerator([]interface{}{"instance"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Numeric,
				common.Numeric,
			},
			PrimaryUnit: "rps",
		},
		QuotaCountConfigMapHard: {
			Label: "OBJECT COUNT CONFIGMAPS HARD",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"hard\",resource=~\".*configmaps\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountConfigMapUsed: {
			Label: "OBJECT COUNT CONFIGMAPS USED",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"used\",resource=~\".*configmaps\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountPersistentVolumeClaimHard: {
			Label: "OBJECT COUNT PERSISTENT VOLUME CLAIMS HARD",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"hard\",resource=~\"persistentvolumeclaims\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountPersistentVolumeClaimUsed: {
			Label: "OBJECT COUNT PERSISTENT VOLUME CLAIMS USED",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"used\",resource=~\"persistentvolumeclaims\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountPodHard: {
			Label: "OBJECT COUNT PODS HARD",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"hard\",resource=~\".*pods\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountPodUsed: {
			Label: "OBJECT COUNT PODS USED",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"used\",resource=~\".*pods\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountReplicationControllerHard: {
			Label: "OBJECT COUNT REPLICATION CONTROLLERS HARD",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"hard\",resource=~\".*replicationcontrollers\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountReplicationControllerUsed: {
			Label: "OBJECT COUNT REPLICATION CONTROLLERS USED",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"used\",resource=~\".*replicationcontrollers\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountResourceQuotaHard: {
			Label: "OBJECT COUNT RESOURCE QUOTAS HARD",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"hard\",resource=~\".*resourcequotas\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountResourceQuotaUsed: {
			Label: "OBJECT COUNT RESOURCE QUOTAS USED",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"used\",resource=~\".*resourcequotas\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountSecretHard: {
			Label: "OBJECT COUNT SECRETS HARD",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"hard\",resource=~\".*secrets\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountSecretUsed: {
			Label: "OBJECT COUNT SECRETS USED",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"used\",resource=~\".*secrets\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountServiceHard: {
			Label: "OBJECT COUNT SERVICES HARD",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"hard\",resource=~\".*services\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountServiceUsed: {
			Label: "OBJECT COUNT SERVICES USED",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"used\",resource=~\".*services\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountServiceLoadBalancerHard: {
			Label: "OBJECT COUNT SERVICES LOAD BALANCERS HARD",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"hard\",resource=~\".*services.loadbalancers\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountServiceLoadBalancerUsed: {
			Label: "OBJECT COUNT SERVICES LOAD BALANCERS USED",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"used\",resource=~\".*services.loadbalancers\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountServiceNodePortHard: {
			Label: "OBJECT COUNT SERVICES NODE PORTS HARD",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"hard\",resource=~\".*services.nodeports\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaCountServiceNodePortUsed: {
			Label: "OBJECT COUNT SERVICES NODE PORTS USED",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"used\",resource=~\".*services.nodeports\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "",
		},
		QuotaLimitCpuHard: {
			Label: "CPU LIMIT HARD",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"hard\",resource=\"limits.cpu\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "Core",
		},
		QuotaLimitCpuUsed: {
			Label: "CPU LIMIT USED",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"used\",resource=\"limits.cpu\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "Core",
		},
		QuotaLimitMemoryHard: {
			Label: "MEMORY LIMIT HARD",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"hard\",resource=\"limits.memory\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		QuotaLimitMemoryUsed: {
			Label: "MEMORY LIMIT USED",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"used\",resource=\"limits.memory\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		QuotaLimitPodCpu: {
			Label: "POD CPU LIMIT",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_pod_container_resource_limits{resource=\"cpu\",namespace=~\"%s\",pod=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace", "pod"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "Core",
		},
		QuotaLimitPodEphemeralStorage: {
			Label: "POD EPHEMERAL STORAGE LIMIT",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_pod_container_resource_limits{resource=\"ephemeral_storage\",namespace=~\"%s\",pod=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace", "pod"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		QuotaLimitPodMemory: {
			Label: "POD MEMORY LIMIT",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_pod_container_resource_limits{resource=\"memory\",namespace=~\"%s\",pod=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace", "pod"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		QuotaRequestCpuHard: {
			Label: "CPU REQUEST HARD",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"hard\",resource=\"requests.cpu\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "Core",
		},
		QuotaRequestCpuUsed: {
			Label: "CPU REQUEST USED",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"used\",resource=\"requests.cpu\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				"",
			},
			PrimaryUnit: "Core",
		},
		QuotaRequestMemoryHard: {
			Label: "MEMORY REQUEST HARD",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"hard\",resource=\"requests.memory\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		QuotaRequestMemoryUsed: {
			Label: "MEMORY REQUEST USED",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"hard\",resource=\"requests.memory\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		QuotaRequestPodCpu: {
			Label: "POD CPU REQUEST",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_pod_container_resource_requests{resource=\"cpu\",namespace=~\"%s\",pod=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace", "pod"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Numeric,
			},
			PrimaryUnit: "Core",
		},
		QuotaRequestPodEphemeralStorage: {
			Label: "POD EPHEMERAL STORAGE REQUEST",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_pod_container_resource_requests{resource=\"ephemeral_storage\",namespace=~\"%s\",pod=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace", "pod"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		QuotaRequestPodMemory: {
			Label: "POD MEMORY REQUEST",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_pod_container_resource_requests{resource=\"memory\",namespace=~\"%s\",pod=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace", "pod"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		QuotaRequestStorageHard: {
			Label: "STORAGE REQUEST HARD",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"hard\",resource=\"requests.storage\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		QuotaRequestStorageUsed: {
			Label: "STORAGE REQUEST USED",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						"sum(kube_resourcequota{type=\"used\",resource=\"requests.storage\",namespace=~\"%s\"})",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		SummaryNodeInfo: {
			MetricKeys: []MetricKey{CustomNodeCpu, CustomNodeFileSystem, CustomNodeMemory, NodeNetworkIn, NodeNetworkOut, NumberOfPod},
		},
		SummaryContainerCpuInfo: {
			MetricKeys: []MetricKey{ContainerCpu, QuotaRequestPodCpu, QuotaLimitPodCpu},
		},
		SummaryContainerMemoryInfo: {
			MetricKeys: []MetricKey{ContainerMemory, QuotaRequestPodMemory, QuotaLimitPodMemory},
		},
		TopNodeCpuByNode: {
			Label: "CPU",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드의 CPU 사용량에 따른 노드 내림차순 목록
						"sort_desc(sum(rate(node_cpu_seconds_total{mode!=\"idle\",mode!=\"iowait\",instance=~\"%s\"}[3m]))by(instance))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		TopNodeFileSystemByNode: {
			Label: "FILE SYSTEM",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드의 FILE SYSTEM 사용량에 따른 노드 내림차순 목록
						"sort_desc(sum(node_filesystem_size_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"}-node_filesystem_avail_bytes{mountpoint=\"/\",fstype!=\"rootfs\",instance=~\"%s\"})by(instance))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance", "instance"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		TopNodeMemoryByNode: {
			Label: "MEMORY",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드의 MEMORY 사용량에 따른 노드 내림차순 목록
						"sort_desc(sum(node_memory_MemTotal_bytes-node_memory_MemAvailable_bytes{instance=~\"%s\"})by(instance))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		TopNodeNetworkInByNode: {
			Label: "NETWORK IN",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드의 NETWORK IN 에 따른 노드 내림차순 목록
						"sort_desc(sum(rate(node_network_receive_bytes_total{instance=~\"%s\"}[3m]))by(instance))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		TopNodeNetworkOutByNode: {
			Label: "NETWORK OUT",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드의 NETWORK OUT 에 따른 노드 내림차순 목록
						"sort_desc(sum(rate(node_network_receive_bytes_total{instance=~\"%s\"}[3m]))by(instance))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"instance"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		TopNodePodCountByNode: {
			Label: "POD COUNT",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드별 파드 수에 따른 내림차순 목록
						"sort_desc(count(kube_pod_info{node!=\"\",node=~\"%s\"})by(node))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		Top5ContainerCpuByNamespace: {
			Label: "CPU(TOP5 OF PROJECTS)",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 네임스페이스별 CPU 사용량에 따른 내림차순 목록
						"topk(5,sort_desc(sum(rate(container_cpu_usage_seconds_total{container!=\"\",pod!=\"\",node=~\"%s\"}[3m]))by(namespace)))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		Top5ContainerCpuByPod: {
			Label: "CPU(TOP5 OF PODS)",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 파드별 CPU 사용량에 따른 내림차순 목록
						"topk(5,sort_desc(sum(rate(container_cpu_usage_seconds_total{container!=\"\",pod!=\"\",node=~\"%s\",namespace=~\"%s\"}[3m]))by(pod)))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node", "namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Core,
			},
			PrimaryUnit: "Core",
		},
		Top5ContainerFileSystemByNamespace: {
			Label: "FILE SYSTEM(TOP5 OF PROJECTS)",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드의 네임스페이스 중 FILE SYSTEM 사용량에 따른 내림차순 목록
						"topk(5,sort_desc(sum(container_fs_usage_bytes{container!=\"\",pod!=\"\",node=~\"%s\"})by(namespace)))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		Top5ContainerFileSystemByPod: {
			Label: "FILE SYSTEM(TOP5 OF PODS)",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 파드별 FILE SYSTEM 사용량에 따른 내림차순 목록
						"topk(5,sort_desc(sum(container_fs_usage_bytes{container!=\"\",pod!=\"\",node=~\"%s\",namespace=~\"%s\"})by(pod)))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node", "namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		Top5ContainerMemoryByNamespace: {
			Label: "MEMORY(TOP5 OF PROJECTS)",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 네임스페이스별 MEMORY 사용량에 따른 내림차순 목록
						"topk(5,sort_desc(sum(container_memory_working_set_bytes{container!=\"\",pod!=\"\",node=~\"%s\"})by(namespace)))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		Top5ContainerMemoryByPod: {
			Label: "MEMORY(TOP5 OF PODS)",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 파드별 MEMORY 사용량에 따른 내림차순 목록
						"topk(5,sort_desc(sum(container_memory_working_set_bytes{container!=\"\",pod!=\"\",node=~\"%s\",namespace=~\"%s\"})by(pod)))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node", "namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.BinaryBytes,
			},
			PrimaryUnit: "B",
		},
		Top5ContainerNetworkInByNamespace: {
			Label: "NETWORK IN(TOP5 OF PROJECTS)",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// NETWORK IN 에 따른 Top 5 내림차순 목록
						"topk(5,sort_desc(sum(rate(container_network_receive_bytes_total{container=\"POD\",pod!=\"\",node=~\"%s\",namespace=~\"%s\"}[3m]))by(namespace)))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node", "namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		Top5ContainerNetworkInByPod: {
			Label: "NETWORK IN(TOP5 OF PODS)",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// NETWORK IN 에 따른 Top 5 내림차순 목록
						"topk(5,sort_desc(sum(rate(container_network_receive_bytes_total{container=\"POD\",pod!=\"\",node=~\"%s\",namespace=~\"%s\"}[3m]))by(pod)))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node", "namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		Top5ContainerNetworkOutByNamespace: {
			Label: "NETWORK OUT(TOP5 OF PROJECTS)",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드의 네임스페이스 중 NETWORK OUT 에 따른 내림차순 목록
						"topk(5,sort_desc(sum(rate(container_network_receive_bytes_total{namespace!=\"\",node=~\"%s\"}[3m]))by(namespace)))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		Top5ContainerNetworkOutByPod: {
			Label: "NETWORK OUT(TOP5 OF PODS)",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 파드 중 NETWORK OUT 에 따른 내림차순 목록
						"topk(5,sort_desc(sum(rate(container_network_receive_bytes_total{pod!= \"\",node=~\"%s\",namespace=~\"%s\"}[3m]))by(pod)))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node", "namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.DecimalBytesPerSec,
			},
			PrimaryUnit: "Bps",
		},
		Top5CountContainerByPod: {
			Label: "CONTAINER COUNT(TOP5 OF PODS)",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 컨테이너 수에 따른 파드 내림차순 목록
						"topk(5,sort_desc(count(kube_pod_container_info{pod=~\"%s\"})by(pod)))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"pod"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
		Top5CountPodByNamespace: {
			Label: "POD COUNT(TOP5 OF PROJECTS)",
			QueryInfos: map[PrometheusVersion]QueryInfo{
				v2_20_0: {
					QueryTemplates: []string{
						// 노드의 네임스페이스 중 파드 수에 따른 내림차순 목록
						"topk(5,sort_desc(count(kube_pod_info{node=~\"%s\",namespace=~\"%s\"})by(namespace)))",
					},
					QueryTemplateParserGenerators: QueryTemplateParserGenerators{
						queryTemplateParserGenerator([]interface{}{"node", "namespace"}),
					},
				},
			},
			UnitTypeKeys: []common.UnitTypeKey{
				common.Count,
			},
			PrimaryUnit: "",
		},
	}
)
