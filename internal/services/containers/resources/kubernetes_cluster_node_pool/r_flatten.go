// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_node_pool

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceKubernetesClusterNodePoolFlatten(d *pluginsdk.ResourceData, id *agentpools.AgentPoolId, model *agentpools.AgentPool) error {
	clusterId := commonids.NewKubernetesClusterID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName)

	d.Set("name", id.AgentPoolName)
	d.Set("kubernetes_cluster_id", clusterId.ID())

	if model != nil && model.Properties != nil {
		props := model.Properties
		d.Set("zones", zones.FlattenUntyped(props.AvailabilityZones))

		d.Set("auto_scaling_enabled", props.EnableAutoScaling)
		d.Set("node_public_ip_enabled", props.EnableNodePublicIP)
		d.Set("host_encryption_enabled", props.EnableEncryptionAtHost)
		d.Set("fips_enabled", props.EnableFIPS)
		d.Set("ultra_ssd_enabled", props.EnableUltraSSD)

		if v := props.KubeletDiskType; v != nil {
			d.Set("kubelet_disk_type", string(*v))
		}

		if v := props.GpuInstanceProfile; v != nil {
			d.Set("gpu_instance", string(*v))
		}

		if v := props.GpuProfile; v != nil {
			d.Set("gpu_driver", pointer.FromEnum(v.Driver))
		}

		if props.CreationData != nil {
			d.Set("snapshot_id", props.CreationData.SourceResourceId)
		}

		scaleDownMode := string(managedclusters.ScaleDownModeDelete)
		if v := props.ScaleDownMode; v != nil {
			scaleDownMode = string(*v)
		}
		d.Set("scale_down_mode", scaleDownMode)

		if v := props.WorkloadRuntime; v != nil {
			d.Set("workload_runtime", string(*v))
		}

		evictionPolicy := ""
		if v := props.ScaleSetEvictionPolicy; v != nil && *v != "" {
			evictionPolicy = string(*v)
		}
		d.Set("eviction_policy", evictionPolicy)

		if err := d.Set("kubelet_config", helpers.FlattenAgentPoolKubeletConfig(props.KubeletConfig)); err != nil {
			return fmt.Errorf("setting `kubelet_config`: %+v", err)
		}

		linuxOSConfig, err := flattenAgentPoolLinuxOSConfig(props.LinuxOSConfig)
		if err != nil {
			return err
		}
		if err := d.Set("linux_os_config", linuxOSConfig); err != nil {
			return fmt.Errorf("setting `linux_os_config`: %+v", err)
		}

		maxCount := 0
		if props.MaxCount != nil {
			maxCount = int(*props.MaxCount)
		}
		d.Set("max_count", maxCount)

		maxPods := 0
		if props.MaxPods != nil {
			maxPods = int(*props.MaxPods)
		}
		d.Set("max_pods", maxPods)

		minCount := 0
		if props.MinCount != nil {
			minCount = int(*props.MinCount)
		}
		d.Set("min_count", minCount)

		mode := string(managedclusters.AgentPoolModeUser)
		if v := props.Mode; v != nil && *v != "" {
			mode = string(*props.Mode)
		}
		d.Set("mode", mode)

		count := 0
		if props.Count != nil {
			count = int(*props.Count)
		}
		d.Set("node_count", count)

		if err := d.Set("node_labels", props.NodeLabels); err != nil {
			return fmt.Errorf("setting `node_labels`: %+v", err)
		}

		d.Set("node_public_ip_prefix_id", props.NodePublicIPPrefixID)

		if err := d.Set("node_taints", pluginsdk.FlattenSlice(props.NodeTaints)); err != nil {
			return fmt.Errorf("setting `node_taints`: %+v", err)
		}

		d.Set("node_image_version", props.NodeImageVersion)

		// NOTE: workaround for migration from 2022-01-02-preview (<3.12.0) to 2022-03-02-preview (>=3.12.0). Before terraform apply is run against the new API, Azure will respond only with currentOrchestratorVersion, orchestratorVersion will be absent. More details: https://github.com/hashicorp/terraform-provider-azurerm/issues/17833#issuecomment-1227583353
		if props.OrchestratorVersion != nil {
			d.Set("orchestrator_version", props.OrchestratorVersion)
		} else {
			d.Set("orchestrator_version", props.CurrentOrchestratorVersion)
		}

		osDiskSizeGB := 0
		if props.OsDiskSizeGB != nil {
			osDiskSizeGB = int(*props.OsDiskSizeGB)
		}
		d.Set("os_disk_size_gb", osDiskSizeGB)

		osDiskType := agentpools.OSDiskTypeManaged
		if v := props.OsDiskType; v != nil && *v != "" {
			osDiskType = *v
		}
		d.Set("os_disk_type", osDiskType)

		if v := props.OsType; v != nil {
			d.Set("os_type", string(*v))
		}
		if v := props.OsSKU; v != nil {
			d.Set("os_sku", string(*v))
		}
		d.Set("pod_subnet_id", props.PodSubnetID)

		// not returned from the API if not Spot
		priority := string(managedclusters.ScaleSetPriorityRegular)
		if v := props.ScaleSetPriority; v != nil && *v != "" {
			priority = string(*v)
		}
		d.Set("priority", priority)

		d.Set("proximity_placement_group_id", props.ProximityPlacementGroupID)

		spotMaxPrice := -1.0
		if props.SpotMaxPrice != nil {
			spotMaxPrice = *props.SpotMaxPrice
		}
		d.Set("spot_max_price", spotMaxPrice)

		d.Set("vnet_subnet_id", props.VnetSubnetID)
		d.Set("vm_size", props.VMSize)
		d.Set("host_group_id", props.HostGroupID)
		d.Set("capacity_reservation_group_id", props.CapacityReservationGroupID)

		if err := d.Set("upgrade_settings", flattenAgentPoolUpgradeSettings(props.UpgradeSettings)); err != nil {
			return fmt.Errorf("setting `upgrade_settings`: %+v", err)
		}

		if err := d.Set("windows_profile", flattenAgentPoolWindowsProfile(props.WindowsProfile)); err != nil {
			return fmt.Errorf("setting `windows_profile`: %+v", err)
		}

		if err := d.Set("node_network_profile", flattenAgentPoolNetworkProfile(props.NetworkProfile)); err != nil {
			return fmt.Errorf("setting `node_network_profile`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, props.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenAgentPoolUpgradeSettings(input *agentpools.AgentPoolUpgradeSettings) []any {
	// The API returns an empty upgrade settings object for spot node pools, so we need to explicitly check whether there's anything in it
	if input == nil || (input.MaxSurge == nil && input.MaxUnavailable == nil && input.DrainTimeoutInMinutes == nil && input.NodeSoakDurationInMinutes == nil && input.UndrainableNodeBehavior == nil) {
		return []any{}
	}

	values := make(map[string]any)

	if input.MaxSurge != nil && *input.MaxSurge != "" && *input.MaxSurge != "0" && *input.MaxSurge != "0%" {
		values["max_surge"] = *input.MaxSurge
	}
	if input.MaxUnavailable != nil && *input.MaxUnavailable != "" && *input.MaxUnavailable != "0" && *input.MaxUnavailable != "0%" {
		values["max_unavailable"] = *input.MaxUnavailable
	}

	if input.DrainTimeoutInMinutes != nil {
		values["drain_timeout_in_minutes"] = *input.DrainTimeoutInMinutes
	}

	if input.NodeSoakDurationInMinutes != nil {
		values["node_soak_duration_in_minutes"] = *input.NodeSoakDurationInMinutes
	}

	if input.UndrainableNodeBehavior != nil && *input.UndrainableNodeBehavior != "" {
		values["undrainable_node_behavior"] = string(*input.UndrainableNodeBehavior)
	}

	return []any{values}
}

func flattenAgentPoolLinuxOSConfig(input *agentpools.LinuxOSConfig) ([]any, error) {
	if input == nil {
		return make([]any, 0), nil
	}

	var swapFileSizeMB int
	if input.SwapFileSizeMB != nil {
		swapFileSizeMB = int(*input.SwapFileSizeMB)
	}
	transparentHugePageDefrag := pointer.From(input.TransparentHugePageDefrag)
	transparentHugePageEnabled := pointer.From(input.TransparentHugePageEnabled)
	sysctlConfig, err := flattenAgentPoolSysctlConfig(input.Sysctls)
	if err != nil {
		return nil, err
	}
	config := []any{
		map[string]any{
			"swap_file_size_mb":            swapFileSizeMB,
			"sysctl_config":                sysctlConfig,
			"transparent_huge_page_defrag": transparentHugePageDefrag,
			"transparent_huge_page":        transparentHugePageEnabled,
		},
	}

	return config, nil
}

func flattenAgentPoolSysctlConfig(input *agentpools.SysctlConfig) ([]any, error) {
	if input == nil {
		return make([]any, 0), nil
	}

	var fsAioMaxNr int
	if input.FsAioMaxNr != nil {
		fsAioMaxNr = int(*input.FsAioMaxNr)
	}
	var fsFileMax int
	if input.FsFileMax != nil {
		fsFileMax = int(*input.FsFileMax)
	}
	var fsInotifyMaxUserWatches int
	if input.FsInotifyMaxUserWatches != nil {
		fsInotifyMaxUserWatches = int(*input.FsInotifyMaxUserWatches)
	}
	var fsNrOpen int
	if input.FsNrOpen != nil {
		fsNrOpen = int(*input.FsNrOpen)
	}
	var kernelThreadsMax int
	if input.KernelThreadsMax != nil {
		kernelThreadsMax = int(*input.KernelThreadsMax)
	}
	var netCoreNetdevMaxBacklog int
	if input.NetCoreNetdevMaxBacklog != nil {
		netCoreNetdevMaxBacklog = int(*input.NetCoreNetdevMaxBacklog)
	}
	var netCoreOptmemMax int
	if input.NetCoreOptmemMax != nil {
		netCoreOptmemMax = int(*input.NetCoreOptmemMax)
	}
	var netCoreRmemDefault int
	if input.NetCoreRmemDefault != nil {
		netCoreRmemDefault = int(*input.NetCoreRmemDefault)
	}
	var netCoreRmemMax int
	if input.NetCoreRmemMax != nil {
		netCoreRmemMax = int(*input.NetCoreRmemMax)
	}
	var netCoreSomaxconn int
	if input.NetCoreSomaxconn != nil {
		netCoreSomaxconn = int(*input.NetCoreSomaxconn)
	}
	var netCoreWmemDefault int
	if input.NetCoreWmemDefault != nil {
		netCoreWmemDefault = int(*input.NetCoreWmemDefault)
	}
	var netCoreWmemMax int
	if input.NetCoreWmemMax != nil {
		netCoreWmemMax = int(*input.NetCoreWmemMax)
	}
	var netIpv4IpLocalPortRangeMin, netIpv4IpLocalPortRangeMax int
	if input.NetIPv4IPLocalPortRange != nil {
		arr := regexp.MustCompile("[ \t]+").Split(*input.NetIPv4IPLocalPortRange, -1)
		if len(arr) != 2 {
			return nil, fmt.Errorf("parsing `NetIPv4IPLocalPortRange` %s", *input.NetIPv4IPLocalPortRange)
		}
		var err error
		netIpv4IpLocalPortRangeMin, err = strconv.Atoi(arr[0])
		if err != nil {
			return nil, err
		}
		netIpv4IpLocalPortRangeMax, err = strconv.Atoi(arr[1])
		if err != nil {
			return nil, err
		}
	}
	var netIpv4NeighDefaultGcThresh1 int
	if input.NetIPv4NeighDefaultGcThresh1 != nil {
		netIpv4NeighDefaultGcThresh1 = int(*input.NetIPv4NeighDefaultGcThresh1)
	}
	var netIpv4NeighDefaultGcThresh2 int
	if input.NetIPv4NeighDefaultGcThresh2 != nil {
		netIpv4NeighDefaultGcThresh2 = int(*input.NetIPv4NeighDefaultGcThresh2)
	}
	var netIpv4NeighDefaultGcThresh3 int
	if input.NetIPv4NeighDefaultGcThresh3 != nil {
		netIpv4NeighDefaultGcThresh3 = int(*input.NetIPv4NeighDefaultGcThresh3)
	}
	var netIpv4TcpFinTimeout int
	if input.NetIPv4TcpFinTimeout != nil {
		netIpv4TcpFinTimeout = int(*input.NetIPv4TcpFinTimeout)
	}
	var netIpv4TcpkeepaliveIntvl int
	if input.NetIPv4TcpkeepaliveIntvl != nil {
		netIpv4TcpkeepaliveIntvl = int(*input.NetIPv4TcpkeepaliveIntvl)
	}
	var netIpv4TcpKeepaliveProbes int
	if input.NetIPv4TcpKeepaliveProbes != nil {
		netIpv4TcpKeepaliveProbes = int(*input.NetIPv4TcpKeepaliveProbes)
	}
	var netIpv4TcpKeepaliveTime int
	if input.NetIPv4TcpKeepaliveTime != nil {
		netIpv4TcpKeepaliveTime = int(*input.NetIPv4TcpKeepaliveTime)
	}
	var netIpv4TcpMaxSynBacklog int
	if input.NetIPv4TcpMaxSynBacklog != nil {
		netIpv4TcpMaxSynBacklog = int(*input.NetIPv4TcpMaxSynBacklog)
	}
	var netIpv4TcpMaxTwBuckets int
	if input.NetIPv4TcpMaxTwBuckets != nil {
		netIpv4TcpMaxTwBuckets = int(*input.NetIPv4TcpMaxTwBuckets)
	}
	netIpv4TcpTwReuse := pointer.From(input.NetIPv4TcpTwReuse)
	var netNetfilterNfConntrackBuckets int
	if input.NetNetfilterNfConntrackBuckets != nil {
		netNetfilterNfConntrackBuckets = int(*input.NetNetfilterNfConntrackBuckets)
	}
	var netNetfilterNfConntrackMax int
	if input.NetNetfilterNfConntrackMax != nil {
		netNetfilterNfConntrackMax = int(*input.NetNetfilterNfConntrackMax)
	}
	var vmMaxMapCount int
	if input.VMMaxMapCount != nil {
		vmMaxMapCount = int(*input.VMMaxMapCount)
	}
	var vmSwappiness int
	if input.VMSwappiness != nil {
		vmSwappiness = int(*input.VMSwappiness)
	}
	var vmVfsCachePressure int
	if input.VMVfsCachePressure != nil {
		vmVfsCachePressure = int(*input.VMVfsCachePressure)
	}
	return []any{
		map[string]any{
			"fs_aio_max_nr":                      fsAioMaxNr,
			"fs_file_max":                        fsFileMax,
			"fs_inotify_max_user_watches":        fsInotifyMaxUserWatches,
			"fs_nr_open":                         fsNrOpen,
			"kernel_threads_max":                 kernelThreadsMax,
			"net_core_netdev_max_backlog":        netCoreNetdevMaxBacklog,
			"net_core_optmem_max":                netCoreOptmemMax,
			"net_core_rmem_default":              netCoreRmemDefault,
			"net_core_rmem_max":                  netCoreRmemMax,
			"net_core_somaxconn":                 netCoreSomaxconn,
			"net_core_wmem_default":              netCoreWmemDefault,
			"net_core_wmem_max":                  netCoreWmemMax,
			"net_ipv4_ip_local_port_range_min":   netIpv4IpLocalPortRangeMin,
			"net_ipv4_ip_local_port_range_max":   netIpv4IpLocalPortRangeMax,
			"net_ipv4_neigh_default_gc_thresh1":  netIpv4NeighDefaultGcThresh1,
			"net_ipv4_neigh_default_gc_thresh2":  netIpv4NeighDefaultGcThresh2,
			"net_ipv4_neigh_default_gc_thresh3":  netIpv4NeighDefaultGcThresh3,
			"net_ipv4_tcp_fin_timeout":           netIpv4TcpFinTimeout,
			"net_ipv4_tcp_keepalive_intvl":       netIpv4TcpkeepaliveIntvl,
			"net_ipv4_tcp_keepalive_probes":      netIpv4TcpKeepaliveProbes,
			"net_ipv4_tcp_keepalive_time":        netIpv4TcpKeepaliveTime,
			"net_ipv4_tcp_max_syn_backlog":       netIpv4TcpMaxSynBacklog,
			"net_ipv4_tcp_max_tw_buckets":        netIpv4TcpMaxTwBuckets,
			"net_ipv4_tcp_tw_reuse":              netIpv4TcpTwReuse,
			"net_netfilter_nf_conntrack_buckets": netNetfilterNfConntrackBuckets,
			"net_netfilter_nf_conntrack_max":     netNetfilterNfConntrackMax,
			"vm_max_map_count":                   vmMaxMapCount,
			"vm_swappiness":                      vmSwappiness,
			"vm_vfs_cache_pressure":              vmVfsCachePressure,
		},
	}, nil
}

func flattenAgentPoolWindowsProfile(input *agentpools.AgentPoolWindowsProfile) []any {
	if input == nil || input.DisableOutboundNat == nil {
		return []any{}
	}

	outboundNatEnabled := true
	if input.DisableOutboundNat != nil {
		outboundNatEnabled = !*input.DisableOutboundNat
	}

	return []any{
		map[string]any{
			"outbound_nat_enabled": outboundNatEnabled,
		},
	}
}

func flattenAgentPoolNetworkProfile(input *agentpools.AgentPoolNetworkProfile) []any {
	if input == nil || input.NodePublicIPTags == nil && input.AllowedHostPorts == nil && input.ApplicationSecurityGroups == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"allowed_host_ports":             flattenAgentPoolNetworkProfileAllowedHostPorts(input.AllowedHostPorts),
			"application_security_group_ids": pluginsdk.FlattenSlice(input.ApplicationSecurityGroups),
			"node_public_ip_tags":            flattenAgentPoolNetworkProfileNodePublicIPTags(input.NodePublicIPTags),
		},
	}
}

func flattenAgentPoolNetworkProfileAllowedHostPorts(input *[]agentpools.PortRange) []any {
	if input == nil {
		return []any{}
	}
	out := make([]any, 0)
	for _, portRange := range *input {
		out = append(out, map[string]any{
			"port_end":   pointer.From(portRange.PortEnd),
			"port_start": pointer.From(portRange.PortStart),
			"protocol":   pointer.From(portRange.Protocol),
		})
	}
	return out
}

func flattenAgentPoolNetworkProfileNodePublicIPTags(input *[]agentpools.IPTag) map[string]any {
	if input == nil {
		return map[string]any{}
	}
	out := make(map[string]any)

	for _, tag := range *input {
		if tag.IPTagType != nil {
			out[*tag.IPTagType] = tag.Tag
		}
	}

	return out
}
