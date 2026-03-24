// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-07-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-07-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-07-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipprefixes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DefaultNodePoolModel struct {
	Name                       string                    `tfschema:"name"`
	TemporaryNameForRotation   string                    `tfschema:"temporary_name_for_rotation"`
	Type                       string                    `tfschema:"type"`
	VMSize                     string                    `tfschema:"vm_size"`
	CapacityReservationGroupID string                    `tfschema:"capacity_reservation_group_id"`
	KubeletConfig              []KubeletConfigModel      `tfschema:"kubelet_config"`
	LinuxOSConfig              []LinuxOSConfigModel      `tfschema:"linux_os_config"`
	FipsEnabled                bool                      `tfschema:"fips_enabled"`
	GPUInstance                string                    `tfschema:"gpu_instance"`
	GPUDriver                  string                    `tfschema:"gpu_driver"`
	KubeletDiskType            string                    `tfschema:"kubelet_disk_type"`
	MaxCount                   int64                     `tfschema:"max_count"`
	MaxPods                    int64                     `tfschema:"max_pods"`
	MinCount                   int64                     `tfschema:"min_count"`
	NodeNetworkProfile         []NodeNetworkProfileModel `tfschema:"node_network_profile"`
	NodeCount                  int64                     `tfschema:"node_count"`
	NodeLabels                 map[string]string         `tfschema:"node_labels"`
	NodePublicIPPrefixID       string                    `tfschema:"node_public_ip_prefix_id"`
	Tags                       map[string]interface{}    `tfschema:"tags"`
	OSDiskSizeGB               int64                     `tfschema:"os_disk_size_gb"`
	OSDiskType                 string                    `tfschema:"os_disk_type"`
	OSSKU                      string                    `tfschema:"os_sku"`
	UltraSSDEnabled            bool                      `tfschema:"ultra_ssd_enabled"`
	VnetSubnetID               string                    `tfschema:"vnet_subnet_id"`
	OrchestratorVersion        string                    `tfschema:"orchestrator_version"`
	PodSubnetID                string                    `tfschema:"pod_subnet_id"`
	ProximityPlacementGroupID  string                    `tfschema:"proximity_placement_group_id"`
	OnlyCriticalAddonsEnabled  bool                      `tfschema:"only_critical_addons_enabled"`
	ScaleDownMode              string                    `tfschema:"scale_down_mode"`
	SnapshotID                 string                    `tfschema:"snapshot_id"`
	HostGroupID                string                    `tfschema:"host_group_id"`
	UpgradeSettings            []UpgradeSettingsModel    `tfschema:"upgrade_settings"`
	WorkloadRuntime            string                    `tfschema:"workload_runtime"`
	Zones                      []string                  `tfschema:"zones"`
	AutoScalingEnabled         bool                      `tfschema:"auto_scaling_enabled"`
	NodePublicIPEnabled        bool                      `tfschema:"node_public_ip_enabled"`
	HostEncryptionEnabled      bool                      `tfschema:"host_encryption_enabled"`
}

type KubeletConfigModel struct {
	CPUManagerPolicy      string   `tfschema:"cpu_manager_policy"`
	CPUCfsQuotaEnabled    bool     `tfschema:"cpu_cfs_quota_enabled"`
	CPUCfsQuotaPeriod     string   `tfschema:"cpu_cfs_quota_period"`
	ImageGcHighThreshold  int      `tfschema:"image_gc_high_threshold"`
	ImageGcLowThreshold   int      `tfschema:"image_gc_low_threshold"`
	TopologyManagerPolicy string   `tfschema:"topology_manager_policy"`
	AllowedUnsafeSysctls  []string `tfschema:"allowed_unsafe_sysctls"`
	ContainerLogMaxSizeMB int      `tfschema:"container_log_max_size_mb"`
	ContainerLogMaxLine   int      `tfschema:"container_log_max_line"`
	PodMaxPid             int      `tfschema:"pod_max_pid"`
}

type LinuxOSConfigModel struct {
	SysctlConfig              []SysctlConfigModel `tfschema:"sysctl_config"`
	TransparentHugePage       string              `tfschema:"transparent_huge_page"`
	TransparentHugePageDefrag string              `tfschema:"transparent_huge_page_defrag"`
	SwapFileSizeMB            int                 `tfschema:"swap_file_size_mb"`
}

type SysctlConfigModel struct {
	FsAioMaxNr                     int  `tfschema:"fs_aio_max_nr"`
	FsFileMax                      int  `tfschema:"fs_file_max"`
	FsInotifyMaxUserWatches        int  `tfschema:"fs_inotify_max_user_watches"`
	FsNrOpen                       int  `tfschema:"fs_nr_open"`
	KernelThreadsMax               int  `tfschema:"kernel_threads_max"`
	NetCoreNetdevMaxBacklog        int  `tfschema:"net_core_netdev_max_backlog"`
	NetCoreOptmemMax               int  `tfschema:"net_core_optmem_max"`
	NetCoreRmemDefault             int  `tfschema:"net_core_rmem_default"`
	NetCoreRmemMax                 int  `tfschema:"net_core_rmem_max"`
	NetCoreSomaxconn               int  `tfschema:"net_core_somaxconn"`
	NetCoreWmemDefault             int  `tfschema:"net_core_wmem_default"`
	NetCoreWmemMax                 int  `tfschema:"net_core_wmem_max"`
	NetIPv4IPLocalPortRangeMin     int  `tfschema:"net_ipv4_ip_local_port_range_min"`
	NetIPv4IPLocalPortRangeMax     int  `tfschema:"net_ipv4_ip_local_port_range_max"`
	NetIPv4NeighDefaultGcThresh1   int  `tfschema:"net_ipv4_neigh_default_gc_thresh1"`
	NetIPv4NeighDefaultGcThresh2   int  `tfschema:"net_ipv4_neigh_default_gc_thresh2"`
	NetIPv4NeighDefaultGcThresh3   int  `tfschema:"net_ipv4_neigh_default_gc_thresh3"`
	NetIPv4TCPFinTimeout           int  `tfschema:"net_ipv4_tcp_fin_timeout"`
	NetIPv4TCPKeepaliveIntvl       int  `tfschema:"net_ipv4_tcp_keepalive_intvl"`
	NetIPv4TCPKeepaliveProbes      int  `tfschema:"net_ipv4_tcp_keepalive_probes"`
	NetIPv4TCPKeepaliveTime        int  `tfschema:"net_ipv4_tcp_keepalive_time"`
	NetIPv4TCPMaxSynBacklog        int  `tfschema:"net_ipv4_tcp_max_syn_backlog"`
	NetIPv4TCPMaxTwBuckets         int  `tfschema:"net_ipv4_tcp_max_tw_buckets"`
	NetIPv4TCPTwReuse              bool `tfschema:"net_ipv4_tcp_tw_reuse"`
	NetNetfilterNfConntrackBuckets int  `tfschema:"net_netfilter_nf_conntrack_buckets"`
	NetNetfilterNfConntrackMax     int  `tfschema:"net_netfilter_nf_conntrack_max"`
	VMMaxMapCount                  int  `tfschema:"vm_max_map_count"`
	VMSwappiness                   int  `tfschema:"vm_swappiness"`
	VMVfsCachePressure             int  `tfschema:"vm_vfs_cache_pressure"`
}

type NodeNetworkProfileModel struct {
	AllowedHostPorts            []AllowedHostPortsModel `tfschema:"allowed_host_ports"`
	ApplicationSecurityGroupIDs []string                `tfschema:"application_security_group_ids"`
	NodePublicIPTags            map[string]string       `tfschema:"node_public_ip_tags"`
}

type AllowedHostPortsModel struct {
	PortStart int64  `tfschema:"port_start"`
	PortEnd   int64  `tfschema:"port_end"`
	Protocol  string `tfschema:"protocol"`
}

type UpgradeSettingsModel struct {
	MaxSurge                  string `tfschema:"max_surge"`
	DrainTimeoutInMinutes     int64  `tfschema:"drain_timeout_in_minutes"`
	NodeSoakDurationInMinutes int64  `tfschema:"node_soak_duration_in_minutes"`
}

func SchemaDefaultAutomaticClusterNodePoolTyped() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validate.KubernetesAgentPoolName,
				},

				"temporary_name_for_rotation": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.KubernetesAgentPoolName,
				},

				"type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					Default:  string(managedclusters.AgentPoolTypeVirtualMachineScaleSets),
					ValidateFunc: validation.StringInSlice([]string{
						string(managedclusters.AgentPoolTypeVirtualMachineScaleSets),
					}, false),
				},

				"vm_size": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"capacity_reservation_group_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: capacityreservationgroups.ValidateCapacityReservationGroupID,
				},

				"kubelet_config": schemaNodePoolKubeletConfig(),

				"linux_os_config": schemaNodePoolLinuxOSConfig(),

				"fips_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"gpu_instance": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(managedclusters.GPUInstanceProfileMIGOneg),
						string(managedclusters.GPUInstanceProfileMIGTwog),
						string(managedclusters.GPUInstanceProfileMIGThreeg),
						string(managedclusters.GPUInstanceProfileMIGFourg),
						string(managedclusters.GPUInstanceProfileMIGSeveng),
					}, false),
				},

				"gpu_driver": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForGPUDriver(), false),
				},

				"kubelet_disk_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(managedclusters.KubeletDiskTypeOS),
						string(managedclusters.KubeletDiskTypeTemporary),
					}, false),
				},

				"max_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 1000),
				},

				"max_pods": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Computed: true,
				},

				"min_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 1000),
				},

				"node_network_profile": schemaNodePoolNetworkProfile(),

				"node_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 1000),
				},

				"node_labels": {
					Type:     pluginsdk.TypeMap,
					Optional: true,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"node_public_ip_prefix_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: publicipprefixes.ValidatePublicIPPrefixID,
					RequiredWith: []string{"default_node_pool.0.node_public_ip_enabled"},
				},

				"tags": commonschema.Tags(),

				"os_disk_size_gb": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntAtLeast(1),
				},

				"os_disk_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  agentpools.OSDiskTypeEphemeral,
					ValidateFunc: validation.StringInSlice([]string{
						string(managedclusters.OSDiskTypeEphemeral),
					}, false),
				},

				"os_sku": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(agentpools.OSSKUAzureLinux),
						string(agentpools.OSSKUAzureLinuxThree),
						string(agentpools.OSSKUUbuntu),
						string(agentpools.OSSKUUbuntuTwoTwoZeroFour),
						string(agentpools.OSSKUWindowsTwoZeroOneNine),
						string(agentpools.OSSKUWindowsTwoZeroTwoTwo),
					}, false),
				},

				"ultra_ssd_enabled": {
					Type:     pluginsdk.TypeBool,
					Default:  false,
					Optional: true,
				},

				"vnet_subnet_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: commonids.ValidateSubnetID,
				},

				"orchestrator_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"pod_subnet_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: commonids.ValidateSubnetID,
				},

				"proximity_placement_group_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: proximityplacementgroups.ValidateProximityPlacementGroupID,
				},

				"only_critical_addons_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				"scale_down_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(managedclusters.ScaleDownModeDelete),
					ValidateFunc: validation.StringInSlice([]string{
						string(managedclusters.ScaleDownModeDeallocate),
						string(managedclusters.ScaleDownModeDelete),
					}, false),
				},

				"snapshot_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: snapshots.ValidateSnapshotID,
				},

				"host_group_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: computeValidate.HostGroupID,
				},

				"upgrade_settings": upgradeSettingsSchemaAutomaticClusterDefaultNodePool(),

				"workload_runtime": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(managedclusters.WorkloadRuntimeOCIContainer),
					}, false),
				},

				"zones": {
					Type:     schema.TypeSet,
					Optional: true,
					Computed: true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"auto_scaling_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"node_public_ip_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"host_encryption_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func expandClusterNodePoolKubeletConfigTyped(input []KubeletConfigModel) *managedclusters.KubeletConfig {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	result := &managedclusters.KubeletConfig{
		CpuCfsQuota:          pointer.To(config.CPUCfsQuotaEnabled),
		FailSwapOn:           pointer.To(false), // must be false to enable swap file on nodes
		AllowedUnsafeSysctls: pointer.To(config.AllowedUnsafeSysctls),
	}

	if config.CPUManagerPolicy != "" {
		result.CpuManagerPolicy = pointer.To(config.CPUManagerPolicy)
	}
	if config.CPUCfsQuotaPeriod != "" {
		result.CpuCfsQuotaPeriod = pointer.To(config.CPUCfsQuotaPeriod)
	}
	if config.ImageGcHighThreshold != 0 {
		result.ImageGcHighThreshold = pointer.To(int64(config.ImageGcHighThreshold))
	}
	if config.ImageGcLowThreshold != 0 {
		result.ImageGcLowThreshold = pointer.To(int64(config.ImageGcLowThreshold))
	}
	if config.TopologyManagerPolicy != "" {
		result.TopologyManagerPolicy = pointer.To(config.TopologyManagerPolicy)
	}
	if config.ContainerLogMaxSizeMB != 0 {
		result.ContainerLogMaxSizeMB = pointer.To(int64(config.ContainerLogMaxSizeMB))
	}
	if config.ContainerLogMaxLine != 0 {
		result.ContainerLogMaxFiles = pointer.To(int64(config.ContainerLogMaxLine))
	}
	if config.PodMaxPid != 0 {
		result.PodMaxPids = pointer.To(int64(config.PodMaxPid))
	}

	return result
}

func expandClusterNodePoolLinuxOSConfigTyped(input []LinuxOSConfigModel) (*managedclusters.LinuxOSConfig, error) {
	if len(input) == 0 {
		return nil, nil
	}

	config := input[0]
	sysctlConfig, err := expandClusterNodePoolSysctlConfigTyped(config.SysctlConfig)
	if err != nil {
		return nil, err
	}

	result := &managedclusters.LinuxOSConfig{
		Sysctls: sysctlConfig,
	}

	if config.TransparentHugePage != "" {
		result.TransparentHugePageEnabled = pointer.To(config.TransparentHugePage)
	}
	if config.TransparentHugePageDefrag != "" {
		result.TransparentHugePageDefrag = pointer.To(config.TransparentHugePageDefrag)
	}
	if config.SwapFileSizeMB != 0 {
		result.SwapFileSizeMB = pointer.To(int64(config.SwapFileSizeMB))
	}

	return result, nil
}

func expandClusterNodePoolSysctlConfigTyped(input []SysctlConfigModel) (*managedclusters.SysctlConfig, error) {
	if len(input) == 0 {
		return nil, nil
	}

	config := input[0]
	result := &managedclusters.SysctlConfig{
		NetIPv4TcpTwReuse: pointer.To(config.NetIPv4TCPTwReuse),
	}

	if config.NetCoreSomaxconn != 0 {
		result.NetCoreSomaxconn = pointer.To(int64(config.NetCoreSomaxconn))
	}
	if config.NetCoreNetdevMaxBacklog != 0 {
		result.NetCoreNetdevMaxBacklog = pointer.To(int64(config.NetCoreNetdevMaxBacklog))
	}
	if config.NetCoreRmemDefault != 0 {
		result.NetCoreRmemDefault = pointer.To(int64(config.NetCoreRmemDefault))
	}
	if config.NetCoreRmemMax != 0 {
		result.NetCoreRmemMax = pointer.To(int64(config.NetCoreRmemMax))
	}
	if config.NetCoreWmemDefault != 0 {
		result.NetCoreWmemDefault = pointer.To(int64(config.NetCoreWmemDefault))
	}
	if config.NetCoreWmemMax != 0 {
		result.NetCoreWmemMax = pointer.To(int64(config.NetCoreWmemMax))
	}
	if config.NetCoreOptmemMax != 0 {
		result.NetCoreOptmemMax = pointer.To(int64(config.NetCoreOptmemMax))
	}
	if config.NetIPv4TCPMaxSynBacklog != 0 {
		result.NetIPv4TcpMaxSynBacklog = pointer.To(int64(config.NetIPv4TCPMaxSynBacklog))
	}
	if config.NetIPv4TCPMaxTwBuckets != 0 {
		result.NetIPv4TcpMaxTwBuckets = pointer.To(int64(config.NetIPv4TCPMaxTwBuckets))
	}
	if config.NetIPv4TCPFinTimeout != 0 {
		result.NetIPv4TcpFinTimeout = pointer.To(int64(config.NetIPv4TCPFinTimeout))
	}
	if config.NetIPv4TCPKeepaliveTime != 0 {
		result.NetIPv4TcpKeepaliveTime = pointer.To(int64(config.NetIPv4TCPKeepaliveTime))
	}
	if config.NetIPv4TCPKeepaliveProbes != 0 {
		result.NetIPv4TcpKeepaliveProbes = pointer.To(int64(config.NetIPv4TCPKeepaliveProbes))
	}
	if config.NetIPv4TCPKeepaliveIntvl != 0 {
		result.NetIPv4TcpkeepaliveIntvl = pointer.To(int64(config.NetIPv4TCPKeepaliveIntvl))
	}

	if (config.NetIPv4IPLocalPortRangeMin != 0 && config.NetIPv4IPLocalPortRangeMax == 0) ||
		(config.NetIPv4IPLocalPortRangeMin == 0 && config.NetIPv4IPLocalPortRangeMax != 0) {
		return nil, fmt.Errorf("`net_ipv4_ip_local_port_range_min` and `net_ipv4_ip_local_port_range_max` should both be set or unset")
	}
	if config.NetIPv4IPLocalPortRangeMin > config.NetIPv4IPLocalPortRangeMax {
		return nil, fmt.Errorf("`net_ipv4_ip_local_port_range_min` should be no larger than `net_ipv4_ip_local_port_range_max`")
	}
	if config.NetIPv4IPLocalPortRangeMin != 0 && config.NetIPv4IPLocalPortRangeMax != 0 {
		result.NetIPv4IPLocalPortRange = pointer.To(fmt.Sprintf("%d %d", config.NetIPv4IPLocalPortRangeMin, config.NetIPv4IPLocalPortRangeMax))
	}

	if config.NetIPv4NeighDefaultGcThresh1 != 0 {
		result.NetIPv4NeighDefaultGcThresh1 = pointer.To(int64(config.NetIPv4NeighDefaultGcThresh1))
	}
	if config.NetIPv4NeighDefaultGcThresh2 != 0 {
		result.NetIPv4NeighDefaultGcThresh2 = pointer.To(int64(config.NetIPv4NeighDefaultGcThresh2))
	}
	if config.NetIPv4NeighDefaultGcThresh3 != 0 {
		result.NetIPv4NeighDefaultGcThresh3 = pointer.To(int64(config.NetIPv4NeighDefaultGcThresh3))
	}
	if config.NetNetfilterNfConntrackMax != 0 {
		result.NetNetfilterNfConntrackMax = pointer.To(int64(config.NetNetfilterNfConntrackMax))
	}
	if config.NetNetfilterNfConntrackBuckets != 0 {
		result.NetNetfilterNfConntrackBuckets = pointer.To(int64(config.NetNetfilterNfConntrackBuckets))
	}
	if config.FsAioMaxNr != 0 {
		result.FsAioMaxNr = pointer.To(int64(config.FsAioMaxNr))
	}
	if config.FsInotifyMaxUserWatches != 0 {
		result.FsInotifyMaxUserWatches = pointer.To(int64(config.FsInotifyMaxUserWatches))
	}
	if config.FsFileMax != 0 {
		result.FsFileMax = pointer.To(int64(config.FsFileMax))
	}
	if config.FsNrOpen != 0 {
		result.FsNrOpen = pointer.To(int64(config.FsNrOpen))
	}
	if config.KernelThreadsMax != 0 {
		result.KernelThreadsMax = pointer.To(int64(config.KernelThreadsMax))
	}
	if config.VMMaxMapCount != 0 {
		result.VMMaxMapCount = pointer.To(int64(config.VMMaxMapCount))
	}
	if config.VMSwappiness != 0 {
		result.VMSwappiness = pointer.To(int64(config.VMSwappiness))
	}
	if config.VMVfsCachePressure != 0 {
		result.VMVfsCachePressure = pointer.To(int64(config.VMVfsCachePressure))
	}

	return result, nil
}

func expandClusterNodePoolUpgradeSettingsTyped(input []UpgradeSettingsModel) *managedclusters.AgentPoolUpgradeSettings {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	result := &managedclusters.AgentPoolUpgradeSettings{}

	if config.MaxSurge != "" {
		result.MaxSurge = pointer.To(config.MaxSurge)
	}
	if config.DrainTimeoutInMinutes != 0 {
		result.DrainTimeoutInMinutes = pointer.To(int64(config.DrainTimeoutInMinutes))
	}
	if config.NodeSoakDurationInMinutes != 0 {
		result.NodeSoakDurationInMinutes = pointer.To(int64(config.NodeSoakDurationInMinutes))
	}

	return result
}

func flattenClusterNodePoolKubeletConfigTyped(input *managedclusters.KubeletConfig) []KubeletConfigModel {
	if input == nil {
		return []KubeletConfigModel{}
	}

	cpuCfsQuotaEnabled := false
	if input.CpuCfsQuota != nil {
		cpuCfsQuotaEnabled = *input.CpuCfsQuota
	}

	allowedUnsafeSysctls := []string{}
	if input.AllowedUnsafeSysctls != nil {
		allowedUnsafeSysctls = *input.AllowedUnsafeSysctls
	}

	result := KubeletConfigModel{
		CPUCfsQuotaEnabled:   cpuCfsQuotaEnabled,
		AllowedUnsafeSysctls: allowedUnsafeSysctls,
	}

	if input.CpuManagerPolicy != nil {
		result.CPUManagerPolicy = *input.CpuManagerPolicy
	}
	if input.CpuCfsQuotaPeriod != nil {
		result.CPUCfsQuotaPeriod = *input.CpuCfsQuotaPeriod
	}
	if input.ImageGcHighThreshold != nil {
		result.ImageGcHighThreshold = int(*input.ImageGcHighThreshold)
	}
	if input.ImageGcLowThreshold != nil {
		result.ImageGcLowThreshold = int(*input.ImageGcLowThreshold)
	}
	if input.TopologyManagerPolicy != nil {
		result.TopologyManagerPolicy = *input.TopologyManagerPolicy
	}
	if input.ContainerLogMaxSizeMB != nil {
		result.ContainerLogMaxSizeMB = int(*input.ContainerLogMaxSizeMB)
	}
	if input.ContainerLogMaxFiles != nil {
		result.ContainerLogMaxLine = int(*input.ContainerLogMaxFiles)
	}
	if input.PodMaxPids != nil {
		result.PodMaxPid = int(*input.PodMaxPids)
	}

	return []KubeletConfigModel{result}
}

func flattenClusterNodePoolLinuxOSConfigTyped(input *managedclusters.LinuxOSConfig) ([]LinuxOSConfigModel, error) {
	if input == nil {
		return []LinuxOSConfigModel{}, nil
	}

	sysctlConfig, err := flattenClusterNodePoolSysctlConfigTyped(input.Sysctls)
	if err != nil {
		return nil, err
	}

	result := LinuxOSConfigModel{
		SysctlConfig: sysctlConfig,
	}

	if input.TransparentHugePageEnabled != nil {
		result.TransparentHugePage = *input.TransparentHugePageEnabled
	}
	if input.TransparentHugePageDefrag != nil {
		result.TransparentHugePageDefrag = *input.TransparentHugePageDefrag
	}
	if input.SwapFileSizeMB != nil {
		result.SwapFileSizeMB = int(*input.SwapFileSizeMB)
	}

	return []LinuxOSConfigModel{result}, nil
}

func flattenClusterNodePoolSysctlConfigTyped(input *managedclusters.SysctlConfig) ([]SysctlConfigModel, error) {
	if input == nil {
		return []SysctlConfigModel{}, nil
	}

	netIPv4TcpTwReuse := false
	if input.NetIPv4TcpTwReuse != nil {
		netIPv4TcpTwReuse = *input.NetIPv4TcpTwReuse
	}

	result := SysctlConfigModel{
		NetIPv4TCPTwReuse: netIPv4TcpTwReuse,
	}

	if input.NetCoreSomaxconn != nil {
		result.NetCoreSomaxconn = int(*input.NetCoreSomaxconn)
	}
	if input.NetCoreNetdevMaxBacklog != nil {
		result.NetCoreNetdevMaxBacklog = int(*input.NetCoreNetdevMaxBacklog)
	}
	if input.NetCoreRmemDefault != nil {
		result.NetCoreRmemDefault = int(*input.NetCoreRmemDefault)
	}
	if input.NetCoreRmemMax != nil {
		result.NetCoreRmemMax = int(*input.NetCoreRmemMax)
	}
	if input.NetCoreWmemDefault != nil {
		result.NetCoreWmemDefault = int(*input.NetCoreWmemDefault)
	}
	if input.NetCoreWmemMax != nil {
		result.NetCoreWmemMax = int(*input.NetCoreWmemMax)
	}
	if input.NetCoreOptmemMax != nil {
		result.NetCoreOptmemMax = int(*input.NetCoreOptmemMax)
	}
	if input.NetIPv4TcpMaxSynBacklog != nil {
		result.NetIPv4TCPMaxSynBacklog = int(*input.NetIPv4TcpMaxSynBacklog)
	}
	if input.NetIPv4TcpMaxTwBuckets != nil {
		result.NetIPv4TCPMaxTwBuckets = int(*input.NetIPv4TcpMaxTwBuckets)
	}
	if input.NetIPv4TcpFinTimeout != nil {
		result.NetIPv4TCPFinTimeout = int(*input.NetIPv4TcpFinTimeout)
	}
	if input.NetIPv4TcpKeepaliveTime != nil {
		result.NetIPv4TCPKeepaliveTime = int(*input.NetIPv4TcpKeepaliveTime)
	}
	if input.NetIPv4TcpKeepaliveProbes != nil {
		result.NetIPv4TCPKeepaliveProbes = int(*input.NetIPv4TcpKeepaliveProbes)
	}
	if input.NetIPv4TcpkeepaliveIntvl != nil {
		result.NetIPv4TCPKeepaliveIntvl = int(*input.NetIPv4TcpkeepaliveIntvl)
	}

	if input.NetIPv4IPLocalPortRange != nil {
		portRange := *input.NetIPv4IPLocalPortRange
		var min, max int
		if _, err := fmt.Sscanf(portRange, "%d %d", &min, &max); err == nil {
			result.NetIPv4IPLocalPortRangeMin = min
			result.NetIPv4IPLocalPortRangeMax = max
		}
	}

	if input.NetIPv4NeighDefaultGcThresh1 != nil {
		result.NetIPv4NeighDefaultGcThresh1 = int(*input.NetIPv4NeighDefaultGcThresh1)
	}
	if input.NetIPv4NeighDefaultGcThresh2 != nil {
		result.NetIPv4NeighDefaultGcThresh2 = int(*input.NetIPv4NeighDefaultGcThresh2)
	}
	if input.NetIPv4NeighDefaultGcThresh3 != nil {
		result.NetIPv4NeighDefaultGcThresh3 = int(*input.NetIPv4NeighDefaultGcThresh3)
	}
	if input.NetNetfilterNfConntrackMax != nil {
		result.NetNetfilterNfConntrackMax = int(*input.NetNetfilterNfConntrackMax)
	}
	if input.NetNetfilterNfConntrackBuckets != nil {
		result.NetNetfilterNfConntrackBuckets = int(*input.NetNetfilterNfConntrackBuckets)
	}
	if input.FsAioMaxNr != nil {
		result.FsAioMaxNr = int(*input.FsAioMaxNr)
	}
	if input.FsInotifyMaxUserWatches != nil {
		result.FsInotifyMaxUserWatches = int(*input.FsInotifyMaxUserWatches)
	}
	if input.FsFileMax != nil {
		result.FsFileMax = int(*input.FsFileMax)
	}
	if input.FsNrOpen != nil {
		result.FsNrOpen = int(*input.FsNrOpen)
	}
	if input.KernelThreadsMax != nil {
		result.KernelThreadsMax = int(*input.KernelThreadsMax)
	}
	if input.VMMaxMapCount != nil {
		result.VMMaxMapCount = int(*input.VMMaxMapCount)
	}
	if input.VMSwappiness != nil {
		result.VMSwappiness = int(*input.VMSwappiness)
	}
	if input.VMVfsCachePressure != nil {
		result.VMVfsCachePressure = int(*input.VMVfsCachePressure)
	}

	return []SysctlConfigModel{result}, nil
}

func flattenClusterNodePoolUpgradeSettingsTyped(input *managedclusters.AgentPoolUpgradeSettings) []UpgradeSettingsModel {
	if input == nil || (input.MaxSurge == nil && input.DrainTimeoutInMinutes == nil && input.NodeSoakDurationInMinutes == nil) {
		return []UpgradeSettingsModel{}
	}

	result := UpgradeSettingsModel{}

	if input.MaxSurge != nil {
		result.MaxSurge = *input.MaxSurge
	}
	if input.DrainTimeoutInMinutes != nil {
		result.DrainTimeoutInMinutes = int64(*input.DrainTimeoutInMinutes)
	}
	if input.NodeSoakDurationInMinutes != nil {
		result.NodeSoakDurationInMinutes = int64(*input.NodeSoakDurationInMinutes)
	}

	return []UpgradeSettingsModel{result}
}

func findDefaultNodePoolTyped(input *[]managedclusters.ManagedClusterAgentPoolProfile) (*managedclusters.ManagedClusterAgentPoolProfile, error) {
	if input == nil {
		return nil, fmt.Errorf("agent pool profiles is nil")
	}

	var agentPool *managedclusters.ManagedClusterAgentPoolProfile
	for _, v := range *input {
		if v.Name == "" {
			continue
		}
		if v.Mode == nil || *v.Mode != managedclusters.AgentPoolModeSystem {
			continue
		}

		agentPool = &v
		break
	}

	if agentPool == nil {
		return nil, fmt.Errorf("unable to determine default agent pool - no System mode pool found")
	}

	return agentPool, nil
}

func flattenClusterPoolNetworkProfileTyped(input *managedclusters.AgentPoolNetworkProfile) []NodeNetworkProfileModel {
	if input == nil || (input.NodePublicIPTags == nil && input.AllowedHostPorts == nil && input.ApplicationSecurityGroups == nil) {
		return []NodeNetworkProfileModel{}
	}

	result := NodeNetworkProfileModel{
		AllowedHostPorts:            flattenClusterPoolNetworkProfileAllowedHostPortsTyped(input.AllowedHostPorts),
		ApplicationSecurityGroupIDs: []string{},
		NodePublicIPTags:            flattenClusterPoolNetworkProfileNodePublicIPTagsTyped(input.NodePublicIPTags),
	}

	if input.ApplicationSecurityGroups != nil {
		result.ApplicationSecurityGroupIDs = *input.ApplicationSecurityGroups
	}

	return []NodeNetworkProfileModel{result}
}

func flattenClusterPoolNetworkProfileAllowedHostPortsTyped(input *[]managedclusters.PortRange) []AllowedHostPortsModel {
	if input == nil {
		return []AllowedHostPortsModel{}
	}

	result := make([]AllowedHostPortsModel, 0)
	for _, portRange := range *input {
		model := AllowedHostPortsModel{}
		if portRange.PortEnd != nil {
			model.PortEnd = int64(*portRange.PortEnd)
		}
		if portRange.PortStart != nil {
			model.PortStart = int64(*portRange.PortStart)
		}
		if portRange.Protocol != nil {
			model.Protocol = string(*portRange.Protocol)
		}
		result = append(result, model)
	}
	return result
}

func flattenClusterPoolNetworkProfileNodePublicIPTagsTyped(input *[]managedclusters.IPTag) map[string]string {
	if input == nil {
		return map[string]string{}
	}

	result := make(map[string]string)
	for _, tag := range *input {
		if tag.IPTagType != nil && tag.Tag != nil {
			result[*tag.IPTagType] = *tag.Tag
		}
	}

	return result
}

func FlattenDefaultNodePoolTyped(input *[]managedclusters.ManagedClusterAgentPoolProfile) ([]DefaultNodePoolModel, error) {
	if input == nil {
		return []DefaultNodePoolModel{}, nil
	}

	agentPool, err := findDefaultNodePoolTyped(input)
	if err != nil {
		return nil, err
	}

	result := DefaultNodePoolModel{
		Name: agentPool.Name,
	}

	if agentPool.Count != nil {
		result.NodeCount = int64(*agentPool.Count)
	}

	if agentPool.EnableUltraSSD != nil {
		result.UltraSSDEnabled = *agentPool.EnableUltraSSD
	}

	if agentPool.EnableAutoScaling != nil {
		result.AutoScalingEnabled = *agentPool.EnableAutoScaling
	}

	if agentPool.EnableFIPS != nil {
		result.FipsEnabled = *agentPool.EnableFIPS
	}

	if agentPool.EnableNodePublicIP != nil {
		result.NodePublicIPEnabled = *agentPool.EnableNodePublicIP
	}

	if agentPool.EnableEncryptionAtHost != nil {
		result.HostEncryptionEnabled = *agentPool.EnableEncryptionAtHost
	}

	if agentPool.GpuInstanceProfile != nil {
		result.GPUInstance = string(*agentPool.GpuInstanceProfile)
	}

	if agentPool.GpuProfile != nil && agentPool.GpuProfile.Driver != nil {
		result.GPUDriver = string(*agentPool.GpuProfile.Driver)
	}

	if agentPool.MaxCount != nil {
		result.MaxCount = int64(*agentPool.MaxCount)
	}

	if agentPool.MaxPods != nil {
		result.MaxPods = int64(*agentPool.MaxPods)
	}

	if agentPool.MinCount != nil {
		result.MinCount = int64(*agentPool.MinCount)
	}

	if agentPool.NodeLabels != nil {
		result.NodeLabels = make(map[string]string)
		for k, v := range *agentPool.NodeLabels {
			result.NodeLabels[k] = v
		}
	}

	if agentPool.NodePublicIPPrefixID != nil {
		result.NodePublicIPPrefixID = *agentPool.NodePublicIPPrefixID
	}

	// Check for CriticalAddonsOnly taint
	if agentPool.NodeTaints != nil {
		for _, taint := range *agentPool.NodeTaints {
			if taint == "CriticalAddonsOnly=true:NoSchedule" {
				result.OnlyCriticalAddonsEnabled = true
				break
			}
		}
	}

	if agentPool.OsDiskSizeGB != nil {
		result.OSDiskSizeGB = int64(*agentPool.OsDiskSizeGB)
	}

	if agentPool.OsDiskType != nil {
		result.OSDiskType = string(*agentPool.OsDiskType)
	} else {
		result.OSDiskType = string(managedclusters.OSDiskTypeManaged)
	}

	if agentPool.PodSubnetID != nil {
		result.PodSubnetID = *agentPool.PodSubnetID
	}

	if agentPool.VnetSubnetID != nil {
		result.VnetSubnetID = *agentPool.VnetSubnetID
	}

	if agentPool.HostGroupID != nil {
		result.HostGroupID = *agentPool.HostGroupID
	}

	// NOTE: workaround for migration from 2022-01-02-preview (<3.12.0) to 2022-03-02-preview (>=3.12.0)
	// Before terraform apply is run against the new API, Azure will respond only with currentOrchestratorVersion
	if agentPool.OrchestratorVersion != nil {
		result.OrchestratorVersion = *agentPool.OrchestratorVersion
	} else if agentPool.CurrentOrchestratorVersion != nil {
		result.OrchestratorVersion = *agentPool.CurrentOrchestratorVersion
	}

	if agentPool.ProximityPlacementGroupID != nil {
		result.ProximityPlacementGroupID = *agentPool.ProximityPlacementGroupID
	}

	if agentPool.ScaleDownMode != nil {
		result.ScaleDownMode = string(*agentPool.ScaleDownMode)
	} else {
		result.ScaleDownMode = string(managedclusters.ScaleDownModeDelete)
	}

	if agentPool.CreationData != nil && agentPool.CreationData.SourceResourceId != nil {
		id, err := snapshots.ParseSnapshotIDInsensitively(*agentPool.CreationData.SourceResourceId)
		if err != nil {
			return nil, err
		}
		result.SnapshotID = id.ID()
	}

	if agentPool.VMSize != nil {
		result.VMSize = *agentPool.VMSize
	}

	if agentPool.CapacityReservationGroupID != nil {
		result.CapacityReservationGroupID = *agentPool.CapacityReservationGroupID
	}

	if agentPool.WorkloadRuntime != nil {
		result.WorkloadRuntime = string(*agentPool.WorkloadRuntime)
	}

	if agentPool.KubeletDiskType != nil {
		result.KubeletDiskType = string(*agentPool.KubeletDiskType)
	}

	if agentPool.OsSKU != nil {
		result.OSSKU = string(*agentPool.OsSKU)
	}

	if agentPool.Type != nil {
		result.Type = string(*agentPool.Type)
	}

	result.UpgradeSettings = flattenClusterNodePoolUpgradeSettingsTyped(agentPool.UpgradeSettings)

	linuxOSConfig, err := flattenClusterNodePoolLinuxOSConfigTyped(agentPool.LinuxOSConfig)
	if err != nil {
		return nil, err
	}
	result.LinuxOSConfig = linuxOSConfig

	result.KubeletConfig = flattenClusterNodePoolKubeletConfigTyped(agentPool.KubeletConfig)
	result.NodeNetworkProfile = flattenClusterPoolNetworkProfileTyped(agentPool.NetworkProfile)

	if agentPool.AvailabilityZones != nil {
		result.Zones = *agentPool.AvailabilityZones
	}

	if agentPool.Tags != nil {
		result.Tags = tags.Flatten(agentPool.Tags)
	}

	return []DefaultNodePoolModel{result}, nil
}
