// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_node_pool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/agentpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandAgentPoolKubeletConfig(input []any) *agentpools.KubeletConfig {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)
	result := &agentpools.KubeletConfig{
		CpuCfsQuota: pointer.To(raw["cpu_cfs_quota_enabled"].(bool)),
		// must be false, otherwise the backend will report error: CustomKubeletConfig.FailSwapOn must be set to false to enable swap file on nodes.
		FailSwapOn:           pointer.To(false),
		AllowedUnsafeSysctls: pluginsdk.ExpandStringSlice(raw["allowed_unsafe_sysctls"].(*pluginsdk.Set).List()),
	}

	if v := raw["cpu_manager_policy"].(string); v != "" {
		result.CpuManagerPolicy = pointer.To(v)
	}
	if v := raw["cpu_cfs_quota_period"].(string); v != "" {
		result.CpuCfsQuotaPeriod = pointer.To(v)
	}
	if v := raw["image_gc_high_threshold"].(int); v != 0 {
		result.ImageGcHighThreshold = pointer.To(int64(v))
	}
	if v := raw["image_gc_low_threshold"].(int); v != 0 {
		result.ImageGcLowThreshold = pointer.To(int64(v))
	}
	if v := raw["topology_manager_policy"].(string); v != "" {
		result.TopologyManagerPolicy = pointer.To(v)
	}
	if v := raw["container_log_max_size_mb"].(int); v != 0 {
		result.ContainerLogMaxSizeMB = pointer.To(int64(v))
	}
	if v := raw["container_log_max_files"].(int); v != 0 {
		result.ContainerLogMaxFiles = pointer.To(int64(v))
	}
	if v := raw["pod_max_pid"].(int); v != 0 {
		result.PodMaxPids = pointer.To(int64(v))
	}

	return result
}

func expandAgentPoolUpgradeSettings(input []any) *agentpools.AgentPoolUpgradeSettings {
	setting := &agentpools.AgentPoolUpgradeSettings{}
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)
	maxSurgeRaw := v["max_surge"].(string)
	maxUnavailableRaw := v["max_unavailable"].(string)

	if maxSurgeRaw != "" {
		setting.MaxSurge = pointer.To(maxSurgeRaw)
		setting.MaxUnavailable = pointer.To("0")
	}
	if maxUnavailableRaw != "" {
		setting.MaxUnavailable = pointer.To(maxUnavailableRaw)
		setting.MaxSurge = pointer.To("0")
	}

	if drainTimeoutInMinutesRaw, ok := v["drain_timeout_in_minutes"].(int); ok {
		setting.DrainTimeoutInMinutes = pointer.To(int64(drainTimeoutInMinutesRaw))
	}
	if nodeSoakDurationInMinutesRaw, ok := v["node_soak_duration_in_minutes"].(int); ok {
		setting.NodeSoakDurationInMinutes = pointer.To(int64(nodeSoakDurationInMinutesRaw))
	}
	if undrainableNodeBehaviorRaw, ok := v["undrainable_node_behavior"].(string); ok && undrainableNodeBehaviorRaw != "" {
		setting.UndrainableNodeBehavior = pointer.ToEnum[agentpools.UndrainableNodeBehavior](undrainableNodeBehaviorRaw)
	}
	return setting
}

func expandAgentPoolLinuxOSConfig(input []any) (*agentpools.LinuxOSConfig, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	raw := input[0].(map[string]any)
	sysctlConfig, err := expandAgentPoolSysctlConfig(raw["sysctl_config"].([]any))
	if err != nil {
		return nil, err
	}

	result := &agentpools.LinuxOSConfig{
		Sysctls: sysctlConfig,
	}
	if v := raw["transparent_huge_page"].(string); v != "" {
		result.TransparentHugePageEnabled = pointer.To(v)
	}
	if v := raw["transparent_huge_page_defrag"].(string); v != "" {
		result.TransparentHugePageDefrag = pointer.To(v)
	}
	if v := raw["swap_file_size_mb"].(int); v != 0 {
		result.SwapFileSizeMB = pointer.To(int64(v))
	}
	return result, nil
}

func expandAgentPoolSysctlConfig(input []any) (*agentpools.SysctlConfig, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	raw := input[0].(map[string]any)
	result := &agentpools.SysctlConfig{
		NetIPv4TcpTwReuse: pointer.To(raw["net_ipv4_tcp_tw_reuse"].(bool)),
	}
	if v := raw["net_core_somaxconn"].(int); v != 0 {
		result.NetCoreSomaxconn = pointer.To(int64(v))
	}
	if v := raw["net_core_netdev_max_backlog"].(int); v != 0 {
		result.NetCoreNetdevMaxBacklog = pointer.To(int64(v))
	}
	if v := raw["net_core_rmem_default"].(int); v != 0 {
		result.NetCoreRmemDefault = pointer.To(int64(v))
	}
	if v := raw["net_core_rmem_max"].(int); v != 0 {
		result.NetCoreRmemMax = pointer.To(int64(v))
	}
	if v := raw["net_core_wmem_default"].(int); v != 0 {
		result.NetCoreWmemDefault = pointer.To(int64(v))
	}
	if v := raw["net_core_wmem_max"].(int); v != 0 {
		result.NetCoreWmemMax = pointer.To(int64(v))
	}
	if v := raw["net_core_optmem_max"].(int); v != 0 {
		result.NetCoreOptmemMax = pointer.To(int64(v))
	}
	if v := raw["net_ipv4_tcp_max_syn_backlog"].(int); v != 0 {
		result.NetIPv4TcpMaxSynBacklog = pointer.To(int64(v))
	}
	if v := raw["net_ipv4_tcp_max_tw_buckets"].(int); v != 0 {
		result.NetIPv4TcpMaxTwBuckets = pointer.To(int64(v))
	}
	if v := raw["net_ipv4_tcp_fin_timeout"].(int); v != 0 {
		result.NetIPv4TcpFinTimeout = pointer.To(int64(v))
	}
	if v := raw["net_ipv4_tcp_keepalive_time"].(int); v != 0 {
		result.NetIPv4TcpKeepaliveTime = pointer.To(int64(v))
	}
	if v := raw["net_ipv4_tcp_keepalive_probes"].(int); v != 0 {
		result.NetIPv4TcpKeepaliveProbes = pointer.To(int64(v))
	}
	if v := raw["net_ipv4_tcp_keepalive_intvl"].(int); v != 0 {
		result.NetIPv4TcpkeepaliveIntvl = pointer.To(int64(v))
	}
	netIpv4IPLocalPortRangeMin := raw["net_ipv4_ip_local_port_range_min"].(int)
	netIpv4IPLocalPortRangeMax := raw["net_ipv4_ip_local_port_range_max"].(int)
	if (netIpv4IPLocalPortRangeMin != 0 && netIpv4IPLocalPortRangeMax == 0) || (netIpv4IPLocalPortRangeMin == 0 && netIpv4IPLocalPortRangeMax != 0) {
		return nil, fmt.Errorf("`net_ipv4_ip_local_port_range_min` and `net_ipv4_ip_local_port_range_max` should both be set or unset")
	}
	if netIpv4IPLocalPortRangeMin > netIpv4IPLocalPortRangeMax {
		return nil, fmt.Errorf("`net_ipv4_ip_local_port_range_min` should be no larger than `net_ipv4_ip_local_port_range_max`")
	}
	if netIpv4IPLocalPortRangeMin != 0 && netIpv4IPLocalPortRangeMax != 0 {
		result.NetIPv4IPLocalPortRange = pointer.To(fmt.Sprintf("%d %d", netIpv4IPLocalPortRangeMin, netIpv4IPLocalPortRangeMax))
	}
	if v := raw["net_ipv4_neigh_default_gc_thresh1"].(int); v != 0 {
		result.NetIPv4NeighDefaultGcThresh1 = pointer.To(int64(v))
	}
	if v := raw["net_ipv4_neigh_default_gc_thresh2"].(int); v != 0 {
		result.NetIPv4NeighDefaultGcThresh2 = pointer.To(int64(v))
	}
	if v := raw["net_ipv4_neigh_default_gc_thresh3"].(int); v != 0 {
		result.NetIPv4NeighDefaultGcThresh3 = pointer.To(int64(v))
	}
	if v := raw["net_netfilter_nf_conntrack_max"].(int); v != 0 {
		result.NetNetfilterNfConntrackMax = pointer.To(int64(v))
	}
	if v := raw["net_netfilter_nf_conntrack_buckets"].(int); v != 0 {
		result.NetNetfilterNfConntrackBuckets = pointer.To(int64(v))
	}
	if v := raw["fs_aio_max_nr"].(int); v != 0 {
		result.FsAioMaxNr = pointer.To(int64(v))
	}
	if v := raw["fs_inotify_max_user_watches"].(int); v != 0 {
		result.FsInotifyMaxUserWatches = pointer.To(int64(v))
	}
	if v := raw["fs_file_max"].(int); v != 0 {
		result.FsFileMax = pointer.To(int64(v))
	}
	if v := raw["fs_nr_open"].(int); v != 0 {
		result.FsNrOpen = pointer.To(int64(v))
	}
	if v := raw["kernel_threads_max"].(int); v != 0 {
		result.KernelThreadsMax = pointer.To(int64(v))
	}
	if v := raw["vm_max_map_count"].(int); v != 0 {
		result.VMMaxMapCount = pointer.To(int64(v))
	}
	if v := raw["vm_swappiness"].(int); v != 0 {
		result.VMSwappiness = pointer.To(int64(v))
	}
	if v := raw["vm_vfs_cache_pressure"].(int); v != 0 {
		result.VMVfsCachePressure = pointer.To(int64(v))
	}
	return result, nil
}

func expandAgentPoolWindowsProfile(input []any) *agentpools.AgentPoolWindowsProfile {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)
	outboundNatEnabled := v["outbound_nat_enabled"].(bool)
	return &agentpools.AgentPoolWindowsProfile{
		DisableOutboundNat: pointer.To(!outboundNatEnabled),
	}
}

func expandAgentPoolNetworkProfile(input []any) *agentpools.AgentPoolNetworkProfile {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)
	return &agentpools.AgentPoolNetworkProfile{
		AllowedHostPorts:          expandAgentPoolNetworkProfileAllowedHostPorts(v["allowed_host_ports"].([]any)),
		ApplicationSecurityGroups: pluginsdk.ExpandStringSlice(v["application_security_group_ids"].([]any)),
		NodePublicIPTags:          expandAgentPoolNetworkProfileNodePublicIPTags(v["node_public_ip_tags"].(map[string]any)),
	}
}

func expandAgentPoolNetworkProfileAllowedHostPorts(input []any) *[]agentpools.PortRange {
	if len(input) == 0 {
		return nil
	}
	out := make([]agentpools.PortRange, 0)
	for _, v := range input {
		raw := v.(map[string]any)
		var portEnd, portStart int64
		var protocol agentpools.Protocol
		if raw["port_end"] != nil {
			portEnd = int64(raw["port_end"].(int))
		}
		if raw["port_start"] != nil {
			portStart = int64(raw["port_start"].(int))
		}
		if raw["protocol"] != nil {
			protocol = agentpools.Protocol(raw["protocol"].(string))
		}
		out = append(out, agentpools.PortRange{
			PortEnd:   pointer.To(portEnd),
			PortStart: pointer.To(portStart),
			Protocol:  pointer.To(protocol),
		})
	}
	return &out
}

func expandAgentPoolNetworkProfileNodePublicIPTags(input map[string]any) *[]agentpools.IPTag {
	if len(input) == 0 {
		return nil
	}
	out := make([]agentpools.IPTag, 0)

	for key, val := range input {
		ipTag := agentpools.IPTag{
			IPTagType: pointer.To(key),
			Tag:       pointer.To(val.(string)),
		}
		out = append(out, ipTag)
	}
	return &out
}
