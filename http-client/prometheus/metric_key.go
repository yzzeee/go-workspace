package prometheus

// MetricKey 메트릭 키
type MetricKey string

const (
	Quota                                 = MetricKey("quota")
	ContainerCpu                          = MetricKey("container_cpu")
	ContainerMemory                       = MetricKey("container_memory")
	ContainerFileSystem                   = MetricKey("container_file_system")
	ContainerNetworkIn                    = MetricKey("container_network_in")
	ContainerNetworkOut                   = MetricKey("container_network_out")
	NodeInfo                              = MetricKey("node_info")
	NodeCpu                               = MetricKey("node_cpu")
	NodeCpuTop                            = MetricKey("node_cpu_top")
	NodeCpuTop5Projects                   = MetricKey("node_cpu_top5_projects")
	NodeCpuTop5Pods                       = MetricKey("node_cpu_top5_pods")
	NodeMemory                            = MetricKey("node_memory")
	NodeMemoryTop                         = MetricKey("node_memory_top")
	NodeMemoryTop5Projects                = MetricKey("node_memory_top5_projects")
	NodeMemoryTop5Pods                    = MetricKey("node_memory_top5_pods")
	NodeFileSystem                        = MetricKey("node_file_system")
	NodeFileSystemTop                     = MetricKey("node_file_system_top")
	NodeFileSystemTop5Projects            = MetricKey("node_file_system_top5_projects")
	NodeFileSystemTop5Pods                = MetricKey("node_file_system_top5_pods")
	NodeNetworkIn                         = MetricKey("node_network_in")
	NodeNetworkInTop                      = MetricKey("node_network_in_top")
	NodeNetworkInTop5Projects             = MetricKey("node_network_in_top5_projects")
	NodeNetworkInTop5Pods                 = MetricKey("node_network_in_top5_pods")
	NodeNetworkOut                        = MetricKey("node_network_out")
	NodeNetworkOutTop                     = MetricKey("node_network_out_top")
	NodeNetworkOutTop5Projects            = MetricKey("node_network_out_top5_projects")
	NodeNetworkOutTop5Pods                = MetricKey("node_network_out_top5_pods")
	NodePodCount                          = MetricKey("node_pod_count")
	NodePodCountTop                       = MetricKey("node_pod_count_top")
	NodePodCountTop5Projects              = MetricKey("node_pod_count_top5_projects")
	QuotaCpuRequest                       = MetricKey("quota_cpu_request")
	QuotaCpuLimit                         = MetricKey("quota_cpu_limit")
	QuotaMemoryRequest                    = MetricKey("quota_memory_request")
	QuotaMemoryLimit                      = MetricKey("quota_memory_limit")
	RangeNodeCpuUsage                     = MetricKey("range_node_cpu_usage")
	RangeContainerCpuUsage                = MetricKey("range_container_cpu_usage")
	RangeCpuLoadAverage                   = MetricKey("range_cpu_load_average")
	RangeNodeMemoryUsage                  = MetricKey("range_node_memory_usage")
	RangeContainerMemoryUsage             = MetricKey("range_container_memory_usage")
	RangeMemorySwap                       = MetricKey("range_memory_swap")
	RangeNodeNetworkIO                    = MetricKey("range_node_network_io")
	RangeContainerNetworkIO               = MetricKey("range_container_network_io")
	RangeNodeNetworkPacket                = MetricKey("range_node_network_packet")
	RangeContainerNetworkPacket           = MetricKey("range_container_network_packet")
	RangeNetworkBandwidth                 = MetricKey("range_network_bandwidth")
	RangeNetworkPacketReceiveTransmit     = MetricKey("range_network_packet_receive_transmit")
	RangeNetworkPacketReceiveTransmitDrop = MetricKey("range_network_packet_receive_transmit_drop")
	RangeFileSystem                       = MetricKey("range_file_system")
	RangeDiskIO                           = MetricKey("range_disk_io")
)
