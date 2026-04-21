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
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipprefixes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/applicationsecuritygroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DefaultNodePoolModel struct {
	Name                       string               `tfschema:"name"`
	TemporaryNameForRotation   string               `tfschema:"temporary_name_for_rotation"`
	Type                       string               `tfschema:"type"`
	VMSize                     string               `tfschema:"vm_size"`
	CapacityReservationGroupID string               `tfschema:"capacity_reservation_group_id"`
	KubeletConfig              []KubeletConfigModel `tfschema:"kubelet_config"`
	LinuxOSConfig              []LinuxOSConfigModel `tfschema:"linux_os_config"`
	FipsEnabled                bool                 `tfschema:"fips_enabled"`
	GPUInstance                string               `tfschema:"gpu_instance"`
	GPUDriver                  string               `tfschema:"gpu_driver"`
	KubeletDiskType            string               `tfschema:"kubelet_disk_type"`
	// MaxCount                   int64                     `tfschema:"max_count"`
	MaxPods int64 `tfschema:"max_pods"`
	// MinCount                   int64                     `tfschema:"min_count"`
	NodeNetworkProfile   []NodeNetworkProfileModel `tfschema:"node_network_profile"`
	NodeCount            int64                     `tfschema:"node_count"`
	NodeLabels           map[string]string         `tfschema:"node_labels"`
	NodePublicIPPrefixID string                    `tfschema:"node_public_ip_prefix_id"`
	Tags                 map[string]interface{}    `tfschema:"tags"`
	OSDiskSizeGB         int64                     `tfschema:"os_disk_size_gb"`
	// OSDiskType                 string                    `tfschema:"os_disk_type"`
	OSSKU                     string `tfschema:"os_sku"`
	UltraSSDEnabled           bool   `tfschema:"ultra_ssd_enabled"`
	VnetSubnetID              string `tfschema:"vnet_subnet_id"`
	OrchestratorVersion       string `tfschema:"orchestrator_version"`
	PodSubnetID               string `tfschema:"pod_subnet_id"`
	ProximityPlacementGroupID string `tfschema:"proximity_placement_group_id"`
	OnlyCriticalAddonsEnabled bool   `tfschema:"only_critical_addons_enabled"`
	// ScaleDownMode             string                 `tfschema:"scale_down_mode"`
	SnapshotID      string                 `tfschema:"snapshot_id"`
	HostGroupID     string                 `tfschema:"host_group_id"`
	UpgradeSettings []UpgradeSettingsModel `tfschema:"upgrade_settings"`
	WorkloadRuntime string                 `tfschema:"workload_runtime"`
	// Zones                     []string               `tfschema:"zones"`
	// AutoScalingEnabled    bool `tfschema:"auto_scaling_enabled"`
	NodePublicIPEnabled   bool `tfschema:"node_public_ip_enabled"`
	HostEncryptionEnabled bool `tfschema:"host_encryption_enabled"`
}

type KubeletConfigModel struct {
	CPUManagerPolicy      string   `tfschema:"cpu_manager_policy"`
	CPUCfsQuotaEnabled    bool     `tfschema:"cpu_cfs_quota_enabled"`
	CPUCfsQuotaPeriod     string   `tfschema:"cpu_cfs_quota_period"`
	ImageGcHighThreshold  int64    `tfschema:"image_gc_high_threshold"`
	ImageGcLowThreshold   int64    `tfschema:"image_gc_low_threshold"`
	TopologyManagerPolicy string   `tfschema:"topology_manager_policy"`
	AllowedUnsafeSysctls  []string `tfschema:"allowed_unsafe_sysctls"`
	ContainerLogMaxSizeMB int64    `tfschema:"container_log_max_size_mb"`
	ContainerLogMaxFiles  int64    `tfschema:"container_log_max_files"`
	PodMaxPid             int64    `tfschema:"pod_max_pid"`
}

type LinuxOSConfigModel struct {
	SysctlConfig              []SysctlConfigModel `tfschema:"sysctl_config"`
	TransparentHugePage       string              `tfschema:"transparent_huge_page"`
	TransparentHugePageDefrag string              `tfschema:"transparent_huge_page_defrag"`
	SwapFileSizeMB            int64               `tfschema:"swap_file_size_mb"`
}

type SysctlConfigModel struct {
	FsAioMaxNr                     int64 `tfschema:"fs_aio_max_nr"`
	FsFileMax                      int64 `tfschema:"fs_file_max"`
	FsInotifyMaxUserWatches        int64 `tfschema:"fs_inotify_max_user_watches"`
	FsNrOpen                       int64 `tfschema:"fs_nr_open"`
	KernelThreadsMax               int64 `tfschema:"kernel_threads_max"`
	NetCoreNetdevMaxBacklog        int64 `tfschema:"net_core_netdev_max_backlog"`
	NetCoreOptmemMax               int64 `tfschema:"net_core_optmem_max"`
	NetCoreRmemDefault             int64 `tfschema:"net_core_rmem_default"`
	NetCoreRmemMax                 int64 `tfschema:"net_core_rmem_max"`
	NetCoreSomaxconn               int64 `tfschema:"net_core_somaxconn"`
	NetCoreWmemDefault             int64 `tfschema:"net_core_wmem_default"`
	NetCoreWmemMax                 int64 `tfschema:"net_core_wmem_max"`
	NetIPv4IPLocalPortRangeMin     int64 `tfschema:"net_ipv4_ip_local_port_range_min"`
	NetIPv4IPLocalPortRangeMax     int64 `tfschema:"net_ipv4_ip_local_port_range_max"`
	NetIPv4NeighDefaultGcThresh1   int64 `tfschema:"net_ipv4_neigh_default_gc_thresh1"`
	NetIPv4NeighDefaultGcThresh2   int64 `tfschema:"net_ipv4_neigh_default_gc_thresh2"`
	NetIPv4NeighDefaultGcThresh3   int64 `tfschema:"net_ipv4_neigh_default_gc_thresh3"`
	NetIPv4TCPFinTimeout           int64 `tfschema:"net_ipv4_tcp_fin_timeout"`
	NetIPv4TCPKeepaliveIntvl       int64 `tfschema:"net_ipv4_tcp_keepalive_intvl"`
	NetIPv4TCPKeepaliveProbes      int64 `tfschema:"net_ipv4_tcp_keepalive_probes"`
	NetIPv4TCPKeepaliveTime        int64 `tfschema:"net_ipv4_tcp_keepalive_time"`
	NetIPv4TCPMaxSynBacklog        int64 `tfschema:"net_ipv4_tcp_max_syn_backlog"`
	NetIPv4TCPMaxTwBuckets         int64 `tfschema:"net_ipv4_tcp_max_tw_buckets"`
	NetIPv4TCPTwReuse              bool  `tfschema:"net_ipv4_tcp_tw_reuse"`
	NetNetfilterNfConntrackBuckets int64 `tfschema:"net_netfilter_nf_conntrack_buckets"`
	NetNetfilterNfConntrackMax     int64 `tfschema:"net_netfilter_nf_conntrack_max"`
	VMMaxMapCount                  int64 `tfschema:"vm_max_map_count"`
	VMSwappiness                   int64 `tfschema:"vm_swappiness"`
	VMVfsCachePressure             int64 `tfschema:"vm_vfs_cache_pressure"`
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
	UndrainableNodeBehavior   string `tfschema:"undrainable_node_behavior"`
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

				"kubelet_config": schemaAutomaticNodePoolKubeletConfig(),

				"linux_os_config": schemaAutomaticNodePoolLinuxOSConfig(),

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

				// "max_count": {
				// 	Type:         pluginsdk.TypeInt,
				// 	Optional:     true,
				// 	Computed:     true,
				// 	ValidateFunc: validation.IntBetween(1, 1000),
				// },

				"max_pods": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Computed: true,
				},

				// "min_count": {
				// 	Type:         pluginsdk.TypeInt,
				// 	Optional:     true,
				// 	Computed:     true,
				// 	ValidateFunc: validation.IntBetween(1, 1000),
				// },

				"node_network_profile": schemaAutomaticNodePoolNetworkProfile(),

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

				// "os_disk_type": {
				// 	Type:     pluginsdk.TypeString,
				// 	Optional: true,
				// 	Default:  agentpools.OSDiskTypeEphemeral,
				//	ValidateFunc: validation.StringInSlice([]string{
				//		string(managedclusters.OSDiskTypeEphemeral),
				//	}, false),
				// },

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

				// "scale_down_mode": {
				// 	Type:     pluginsdk.TypeString,
				// 	Optional: true,
				// 	Default:  string(managedclusters.ScaleDownModeDelete),
				// 	ValidateFunc: validation.StringInSlice([]string{
				//		string(managedclusters.ScaleDownModeDelete),
				//	}, false),
				// },

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

				// "zones": {
				// 	Type:     schema.TypeSet,
				// 	Optional: true,
				// 	Computed: true,
				// 	Elem: &schema.Schema{
				// 		Type:         schema.TypeString,
				// 		ValidateFunc: validation.StringIsNotEmpty,
				// 	},
				// },

				// "auto_scaling_enabled": {
				// 	Type:     pluginsdk.TypeBool,
				// 	Optional: true,
				// 	Default:  false,
				// },

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

func upgradeSettingsSchemaAutomaticClusterDefaultNodePool() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"max_surge": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"drain_timeout_in_minutes": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},
				"node_soak_duration_in_minutes": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 30),
				},
				"undrainable_node_behavior": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForUndrainableNodeBehavior(), true),
				},
			},
		},
	}
}

func schemaAutomaticNodePoolKubeletConfig() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"cpu_manager_policy": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"none",
						"static",
					}, false),
				},

				"cpu_cfs_quota_enabled": {
					Type:     pluginsdk.TypeBool,
					Default:  true,
					Optional: true,
				},

				"cpu_cfs_quota_period": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"image_gc_high_threshold": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 100),
				},

				"image_gc_low_threshold": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 100),
				},

				"topology_manager_policy": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"none",
						"best-effort",
						"restricted",
						"single-numa-node",
					}, false),
				},

				"allowed_unsafe_sysctls": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"container_log_max_size_mb": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},

				"container_log_max_files": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(2),
				},

				"pod_max_pid": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},
			},
		},
	}
}

func schemaAutomaticNodePoolLinuxOSConfig() *pluginsdk.Schema {
	s := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"sysctl_config": schemaNodePoolSysctlConfig(),

				"transparent_huge_page": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"always",
						"madvise",
						"never",
					}, false),
				},

				"transparent_huge_page_defrag": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"always",
						"defer",
						"defer+madvise",
						"madvise",
						"never",
					}, false),
				},

				"swap_file_size_mb": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},
			},
		},
	}
	return s
}

func schemaAutomaticNodePoolNetworkProfile() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"allowed_host_ports": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"port_start": {
								Type:         pluginsdk.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(1, 65535),
							},

							"port_end": {
								Type:         pluginsdk.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(1, 65535),
							},

							"protocol": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(agentpools.ProtocolTCP),
									string(agentpools.ProtocolUDP),
								}, false),
							},
						},
					},
				},

				"application_security_group_ids": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: applicationsecuritygroups.ValidateApplicationSecurityGroupID,
					},
				},

				"node_public_ip_tags": {
					Type:     pluginsdk.TypeMap,
					Optional: true,
					ForceNew: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
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
		result.ImageGcHighThreshold = pointer.To(config.ImageGcHighThreshold)
	}
	if config.ImageGcLowThreshold != 0 {
		result.ImageGcLowThreshold = pointer.To(config.ImageGcLowThreshold)
	}
	if config.TopologyManagerPolicy != "" {
		result.TopologyManagerPolicy = pointer.To(config.TopologyManagerPolicy)
	}
	if config.ContainerLogMaxSizeMB != 0 {
		result.ContainerLogMaxSizeMB = pointer.To(config.ContainerLogMaxSizeMB)
	}
	if config.ContainerLogMaxFiles != 0 {
		result.ContainerLogMaxFiles = pointer.To(config.ContainerLogMaxFiles)
	}
	if config.PodMaxPid != 0 {
		result.PodMaxPids = pointer.To(config.PodMaxPid)
	}

	return result
}

func expandAutomaticClusterNodePoolLinuxOSConfig(input []LinuxOSConfigModel) (*managedclusters.LinuxOSConfig, error) {
	if len(input) == 0 {
		return nil, nil
	}

	config := input[0]
	sysctlConfig, err := expandClusterNodePoolSysctlConfigTyped(config.SysctlConfig)
	if err != nil {
		return nil, err
	}

	result := &managedclusters.LinuxOSConfig{
		Sysctls:                    sysctlConfig,
		TransparentHugePageDefrag:  pointer.To(""),
		TransparentHugePageEnabled: pointer.To(""),
	}

	if config.TransparentHugePage != "" {
		result.TransparentHugePageEnabled = pointer.To(config.TransparentHugePage)
	}
	if config.TransparentHugePageDefrag != "" {
		result.TransparentHugePageDefrag = pointer.To(config.TransparentHugePageDefrag)
	}
	if config.SwapFileSizeMB != 0 {
		result.SwapFileSizeMB = pointer.To(config.SwapFileSizeMB)
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
		result.NetCoreSomaxconn = pointer.To(config.NetCoreSomaxconn)
	}
	if config.NetCoreNetdevMaxBacklog != 0 {
		result.NetCoreNetdevMaxBacklog = pointer.To(config.NetCoreNetdevMaxBacklog)
	}
	if config.NetCoreRmemDefault != 0 {
		result.NetCoreRmemDefault = pointer.To(config.NetCoreRmemDefault)
	}
	if config.NetCoreRmemMax != 0 {
		result.NetCoreRmemMax = pointer.To(config.NetCoreRmemMax)
	}
	if config.NetCoreWmemDefault != 0 {
		result.NetCoreWmemDefault = pointer.To(config.NetCoreWmemDefault)
	}
	if config.NetCoreWmemMax != 0 {
		result.NetCoreWmemMax = pointer.To(config.NetCoreWmemMax)
	}
	if config.NetCoreOptmemMax != 0 {
		result.NetCoreOptmemMax = pointer.To(config.NetCoreOptmemMax)
	}
	if config.NetIPv4TCPMaxSynBacklog != 0 {
		result.NetIPv4TcpMaxSynBacklog = pointer.To(config.NetIPv4TCPMaxSynBacklog)
	}
	if config.NetIPv4TCPMaxTwBuckets != 0 {
		result.NetIPv4TcpMaxTwBuckets = pointer.To(config.NetIPv4TCPMaxTwBuckets)
	}
	if config.NetIPv4TCPFinTimeout != 0 {
		result.NetIPv4TcpFinTimeout = pointer.To(config.NetIPv4TCPFinTimeout)
	}
	if config.NetIPv4TCPKeepaliveTime != 0 {
		result.NetIPv4TcpKeepaliveTime = pointer.To(config.NetIPv4TCPKeepaliveTime)
	}
	if config.NetIPv4TCPKeepaliveProbes != 0 {
		result.NetIPv4TcpKeepaliveProbes = pointer.To(config.NetIPv4TCPKeepaliveProbes)
	}
	if config.NetIPv4TCPKeepaliveIntvl != 0 {
		result.NetIPv4TcpkeepaliveIntvl = pointer.To(config.NetIPv4TCPKeepaliveIntvl)
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
		result.NetIPv4NeighDefaultGcThresh1 = pointer.To(config.NetIPv4NeighDefaultGcThresh1)
	}
	if config.NetIPv4NeighDefaultGcThresh2 != 0 {
		result.NetIPv4NeighDefaultGcThresh2 = pointer.To(config.NetIPv4NeighDefaultGcThresh2)
	}
	if config.NetIPv4NeighDefaultGcThresh3 != 0 {
		result.NetIPv4NeighDefaultGcThresh3 = pointer.To(config.NetIPv4NeighDefaultGcThresh3)
	}
	if config.NetNetfilterNfConntrackMax != 0 {
		result.NetNetfilterNfConntrackMax = pointer.To(config.NetNetfilterNfConntrackMax)
	}
	if config.NetNetfilterNfConntrackBuckets != 0 {
		result.NetNetfilterNfConntrackBuckets = pointer.To(config.NetNetfilterNfConntrackBuckets)
	}
	if config.FsAioMaxNr != 0 {
		result.FsAioMaxNr = pointer.To(config.FsAioMaxNr)
	}
	if config.FsInotifyMaxUserWatches != 0 {
		result.FsInotifyMaxUserWatches = pointer.To(config.FsInotifyMaxUserWatches)
	}
	if config.FsFileMax != 0 {
		result.FsFileMax = pointer.To(config.FsFileMax)
	}
	if config.FsNrOpen != 0 {
		result.FsNrOpen = pointer.To(config.FsNrOpen)
	}
	if config.KernelThreadsMax != 0 {
		result.KernelThreadsMax = pointer.To(config.KernelThreadsMax)
	}
	if config.VMMaxMapCount != 0 {
		result.VMMaxMapCount = pointer.To(config.VMMaxMapCount)
	}
	if config.VMSwappiness != 0 {
		result.VMSwappiness = pointer.To(config.VMSwappiness)
	}
	if config.VMVfsCachePressure != 0 {
		result.VMVfsCachePressure = pointer.To(config.VMVfsCachePressure)
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
		result.DrainTimeoutInMinutes = pointer.To(config.DrainTimeoutInMinutes)
	}
	if config.NodeSoakDurationInMinutes != 0 {
		result.NodeSoakDurationInMinutes = pointer.To(config.NodeSoakDurationInMinutes)
	}

	if config.UndrainableNodeBehavior != "" {
		result.UndrainableNodeBehavior = pointer.To(managedclusters.UndrainableNodeBehavior(config.UndrainableNodeBehavior))
	}

	return result
}

func flattenClusterNodePoolKubeletConfigTyped(input *managedclusters.KubeletConfig) []KubeletConfigModel {
	if input == nil {
		return []KubeletConfigModel{}
	}

	cpuCfsQuotaEnabled := false
	if input.CpuCfsQuota != nil {
		cpuCfsQuotaEnabled = pointer.From(input.CpuCfsQuota)
	}

	allowedUnsafeSysctls := []string{}
	if input.AllowedUnsafeSysctls != nil {
		allowedUnsafeSysctls = pointer.From(input.AllowedUnsafeSysctls)
	}

	result := KubeletConfigModel{
		CPUCfsQuotaEnabled:   cpuCfsQuotaEnabled,
		AllowedUnsafeSysctls: allowedUnsafeSysctls,
	}

	if input.CpuManagerPolicy != nil {
		result.CPUManagerPolicy = pointer.From(input.CpuManagerPolicy)
	}
	if input.CpuCfsQuotaPeriod != nil {
		result.CPUCfsQuotaPeriod = pointer.From(input.CpuCfsQuotaPeriod)
	}
	if input.ImageGcHighThreshold != nil {
		result.ImageGcHighThreshold = pointer.From(input.ImageGcHighThreshold)
	}
	if input.ImageGcLowThreshold != nil {
		result.ImageGcLowThreshold = pointer.From(input.ImageGcLowThreshold)
	}
	if input.TopologyManagerPolicy != nil {
		result.TopologyManagerPolicy = pointer.From(input.TopologyManagerPolicy)
	}
	if input.ContainerLogMaxSizeMB != nil {
		result.ContainerLogMaxSizeMB = pointer.From(input.ContainerLogMaxSizeMB)
	}
	if input.ContainerLogMaxFiles != nil {
		result.ContainerLogMaxFiles = pointer.From(input.ContainerLogMaxFiles)
	}
	if input.PodMaxPids != nil {
		result.PodMaxPid = *input.PodMaxPids
	}

	return []KubeletConfigModel{result}
}

func flattenAutomaticClusterNodePoolLinuxOSConfig(input *managedclusters.LinuxOSConfig) []LinuxOSConfigModel {
	if input == nil {
		return []LinuxOSConfigModel{}
	}

	sysctlConfig := flattenClusterNodePoolSysctlConfigTyped(input.Sysctls)

	result := LinuxOSConfigModel{
		SysctlConfig: sysctlConfig,
	}

	if input.TransparentHugePageEnabled != nil {
		result.TransparentHugePage = pointer.From(input.TransparentHugePageEnabled)
	}
	if input.TransparentHugePageDefrag != nil {
		result.TransparentHugePageDefrag = pointer.From(input.TransparentHugePageDefrag)
	}
	if input.SwapFileSizeMB != nil {
		result.SwapFileSizeMB = pointer.From(input.SwapFileSizeMB)
	}

	return []LinuxOSConfigModel{result}
}

func flattenClusterNodePoolSysctlConfigTyped(input *managedclusters.SysctlConfig) []SysctlConfigModel {
	if input == nil {
		return []SysctlConfigModel{}
	}

	netIPv4TcpTwReuse := false
	if input.NetIPv4TcpTwReuse != nil {
		netIPv4TcpTwReuse = pointer.From(input.NetIPv4TcpTwReuse)
	}

	result := SysctlConfigModel{
		NetIPv4TCPTwReuse: netIPv4TcpTwReuse,
	}

	if input.NetCoreSomaxconn != nil {
		result.NetCoreSomaxconn = pointer.From(input.NetCoreSomaxconn)
	}
	if input.NetCoreNetdevMaxBacklog != nil {
		result.NetCoreNetdevMaxBacklog = pointer.From(input.NetCoreNetdevMaxBacklog)
	}
	if input.NetCoreRmemDefault != nil {
		result.NetCoreRmemDefault = pointer.From(input.NetCoreRmemDefault)
	}
	if input.NetCoreRmemMax != nil {
		result.NetCoreRmemMax = pointer.From(input.NetCoreRmemMax)
	}
	if input.NetCoreWmemDefault != nil {
		result.NetCoreWmemDefault = pointer.From(input.NetCoreWmemDefault)
	}
	if input.NetCoreWmemMax != nil {
		result.NetCoreWmemMax = pointer.From(input.NetCoreWmemMax)
	}
	if input.NetCoreOptmemMax != nil {
		result.NetCoreOptmemMax = pointer.From(input.NetCoreOptmemMax)
	}
	if input.NetIPv4TcpMaxSynBacklog != nil {
		result.NetIPv4TCPMaxSynBacklog = pointer.From(input.NetIPv4TcpMaxSynBacklog)
	}
	if input.NetIPv4TcpMaxTwBuckets != nil {
		result.NetIPv4TCPMaxTwBuckets = pointer.From(input.NetIPv4TcpMaxTwBuckets)
	}
	if input.NetIPv4TcpFinTimeout != nil {
		result.NetIPv4TCPFinTimeout = pointer.From(input.NetIPv4TcpFinTimeout)
	}
	if input.NetIPv4TcpKeepaliveTime != nil {
		result.NetIPv4TCPKeepaliveTime = pointer.From(input.NetIPv4TcpKeepaliveTime)
	}
	if input.NetIPv4TcpKeepaliveProbes != nil {
		result.NetIPv4TCPKeepaliveProbes = pointer.From(input.NetIPv4TcpKeepaliveProbes)
	}
	if input.NetIPv4TcpkeepaliveIntvl != nil {
		result.NetIPv4TCPKeepaliveIntvl = pointer.From(input.NetIPv4TcpkeepaliveIntvl)
	}

	if input.NetIPv4IPLocalPortRange != nil {
		portRange := *input.NetIPv4IPLocalPortRange
		var min, max int64
		if _, err := fmt.Sscanf(portRange, "%d %d", &min, &max); err == nil {
			result.NetIPv4IPLocalPortRangeMin = min
			result.NetIPv4IPLocalPortRangeMax = max
		}
	}

	if input.NetIPv4NeighDefaultGcThresh1 != nil {
		result.NetIPv4NeighDefaultGcThresh1 = pointer.From(input.NetIPv4NeighDefaultGcThresh1)
	}
	if input.NetIPv4NeighDefaultGcThresh2 != nil {
		result.NetIPv4NeighDefaultGcThresh2 = pointer.From(input.NetIPv4NeighDefaultGcThresh2)
	}
	if input.NetIPv4NeighDefaultGcThresh3 != nil {
		result.NetIPv4NeighDefaultGcThresh3 = pointer.From(input.NetIPv4NeighDefaultGcThresh3)
	}
	if input.NetNetfilterNfConntrackMax != nil {
		result.NetNetfilterNfConntrackMax = pointer.From(input.NetNetfilterNfConntrackMax)
	}
	if input.NetNetfilterNfConntrackBuckets != nil {
		result.NetNetfilterNfConntrackBuckets = pointer.From(input.NetNetfilterNfConntrackBuckets)
	}
	if input.FsAioMaxNr != nil {
		result.FsAioMaxNr = pointer.From(input.FsAioMaxNr)
	}
	if input.FsInotifyMaxUserWatches != nil {
		result.FsInotifyMaxUserWatches = pointer.From(input.FsInotifyMaxUserWatches)
	}
	if input.FsFileMax != nil {
		result.FsFileMax = pointer.From(input.FsFileMax)
	}
	if input.FsNrOpen != nil {
		result.FsNrOpen = pointer.From(input.FsNrOpen)
	}
	if input.KernelThreadsMax != nil {
		result.KernelThreadsMax = pointer.From(input.KernelThreadsMax)
	}
	if input.VMMaxMapCount != nil {
		result.VMMaxMapCount = pointer.From(input.VMMaxMapCount)
	}
	if input.VMSwappiness != nil {
		result.VMSwappiness = pointer.From(input.VMSwappiness)
	}
	if input.VMVfsCachePressure != nil {
		result.VMVfsCachePressure = pointer.From(input.VMVfsCachePressure)
	}

	return []SysctlConfigModel{result}
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
		result.DrainTimeoutInMinutes = pointer.From(input.DrainTimeoutInMinutes)
	}
	if input.NodeSoakDurationInMinutes != nil {
		result.NodeSoakDurationInMinutes = pointer.From(input.NodeSoakDurationInMinutes)
	}

	if input.UndrainableNodeBehavior != nil {
		result.UndrainableNodeBehavior = string(pointer.From(input.UndrainableNodeBehavior))
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
	results := make([]NodeNetworkProfileModel, 0)
	result := NodeNetworkProfileModel{
		AllowedHostPorts:            flattenClusterPoolNetworkProfileAllowedHostPortsTyped(input.AllowedHostPorts),
		ApplicationSecurityGroupIDs: []string{},
		NodePublicIPTags:            flattenClusterPoolNetworkProfileNodePublicIPTagsTyped(input.NodePublicIPTags),
	}

	if input.ApplicationSecurityGroups != nil {
		result.ApplicationSecurityGroupIDs = pointer.From(input.ApplicationSecurityGroups)
	}
	results = append(results, result)
	return results
}

func flattenClusterPoolNetworkProfileAllowedHostPortsTyped(input *[]managedclusters.PortRange) []AllowedHostPortsModel {
	if input == nil {
		return []AllowedHostPortsModel{}
	}

	result := make([]AllowedHostPortsModel, 0)
	for _, portRange := range *input {
		model := AllowedHostPortsModel{}
		if portRange.PortEnd != nil {
			model.PortEnd = pointer.From(portRange.PortEnd)
		}
		if portRange.PortStart != nil {
			model.PortStart = pointer.From(portRange.PortStart)
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

func expandClusterPoolNetworkProfileTyped(input []NodeNetworkProfileModel) *managedclusters.AgentPoolNetworkProfile {
	if len(input) == 0 {
		return nil
	}

	profile := input[0]
	result := &managedclusters.AgentPoolNetworkProfile{
		AllowedHostPorts:          expandClusterPoolNetworkProfileAllowedHostPortsTyped(profile.AllowedHostPorts),
		ApplicationSecurityGroups: pointer.To(profile.ApplicationSecurityGroupIDs),
		NodePublicIPTags:          expandClusterPoolNetworkProfileNodePublicIPTagsTyped(profile.NodePublicIPTags),
	}

	return result
}

func expandClusterPoolNetworkProfileAllowedHostPortsTyped(input []AllowedHostPortsModel) *[]managedclusters.PortRange {
	if len(input) == 0 {
		return nil
	}

	out := make([]managedclusters.PortRange, 0, len(input))
	for _, v := range input {
		out = append(out, managedclusters.PortRange{
			PortEnd:   pointer.To(v.PortEnd),
			PortStart: pointer.To(v.PortStart),
			Protocol:  pointer.To(managedclusters.Protocol(v.Protocol)),
		})
	}
	return &out
}

func expandClusterPoolNetworkProfileNodePublicIPTagsTyped(input map[string]string) *[]managedclusters.IPTag {
	if len(input) == 0 {
		return nil
	}

	out := make([]managedclusters.IPTag, 0, len(input))
	for key, val := range input {
		ipTag := managedclusters.IPTag{
			IPTagType: pointer.To(key),
			Tag:       pointer.To(val),
		}
		out = append(out, ipTag)
	}
	return &out
}

// ExpandDefaultNodePoolTyped converts a DefaultNodePoolModel to ManagedClusterAgentPoolProfile
func ExpandDefaultNodePoolTyped(input []DefaultNodePoolModel) (*[]managedclusters.ManagedClusterAgentPoolProfile, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("default_node_pool must be specified")
	}

	raw := input[0]
	// enableAutoScaling := raw.AutoScalingEnabled

	nodeLabels := pointer.To(raw.NodeLabels)
	var nodeTaints *[]string

	if raw.OnlyCriticalAddonsEnabled {
		nodeTaints = pointer.To([]string{"CriticalAddonsOnly=true:NoSchedule"})
	}

	profile := managedclusters.ManagedClusterAgentPoolProfile{
		// EnableAutoScaling:      pointer.To(enableAutoScaling),
		EnableFIPS:             pointer.To(raw.FipsEnabled),
		EnableNodePublicIP:     pointer.To(raw.NodePublicIPEnabled),
		EnableEncryptionAtHost: pointer.To(raw.HostEncryptionEnabled),
		KubeletDiskType:        pointer.To(managedclusters.KubeletDiskType(raw.KubeletDiskType)),
		Name:                   raw.Name,
		NodeLabels:             nodeLabels,
		NodeTaints:             nodeTaints,
		Tags:                   tags.Expand(raw.Tags),
		Type:                   pointer.To(managedclusters.AgentPoolType(raw.Type)),
		VMSize:                 pointer.To(raw.VMSize),

		// at this time the default node pool has to be Linux or the AKS cluster fails to provision with:
		// Pods not in Running status: coredns-7fc597cc45-v5z7x,coredns-autoscaler-7ccc76bfbd-djl7j,metrics-server-cbd95f966-5rl97,tunnelfront-7d9884977b-wpbvn
		// Windows agents can be configured via the separate node pool resource
		OsType: pointer.To(managedclusters.OSTypeLinux),

		// without this set the API returns:
		// Code="MustDefineAtLeastOneSystemPool" Message="Must define at least one system pool."
		// since this is the "default" node pool we can assume this is a system node pool
		Mode: pointer.To(managedclusters.AgentPoolModeSystem),

		UpgradeSettings: expandClusterNodePoolUpgradeSettingsTyped(raw.UpgradeSettings),
	}

	//s := make([]string, 2)
	//s[0] = "1"
	//s[1] = "2"
	//
	//profile.AvailabilityZones = pointer.To(s)

	if raw.MaxPods > 0 {
		profile.MaxPods = pointer.To(raw.MaxPods)
	}

	if raw.NodePublicIPPrefixID != "" {
		profile.NodePublicIPPrefixID = pointer.To(raw.NodePublicIPPrefixID)
	}

	if raw.OSDiskSizeGB > 0 {
		profile.OsDiskSizeGB = pointer.To(raw.OSDiskSizeGB)
	}

	// profile.OsDiskType = pointer.To(managedclusters.OSDiskTypeManaged)
	// if raw.OSDiskType != "" {
	// 	profile.OsDiskType = pointer.To(managedclusters.OSDiskType(raw.OSDiskType))
	//}

	if raw.OSSKU != "" {
		profile.OsSKU = pointer.To(managedclusters.OSSKU(raw.OSSKU))
	}

	if raw.PodSubnetID != "" {
		profile.PodSubnetID = pointer.To(raw.PodSubnetID)
	}

	// scaleDownModeDelete := managedclusters.ScaleDownModeDelete
	// profile.ScaleDownMode = &scaleDownModeDelete
	// if raw.ScaleDownMode != "" {
	profile.ScaleDownMode = pointer.To(managedclusters.ScaleDownModeDelete)
	//}

	if raw.SnapshotID != "" {
		profile.CreationData = &managedclusters.CreationData{
			SourceResourceId: pointer.To(raw.SnapshotID),
		}
	}

	if raw.UltraSSDEnabled {
		profile.EnableUltraSSD = pointer.To(raw.UltraSSDEnabled)
	}

	if raw.VnetSubnetID != "" {
		profile.VnetSubnetID = pointer.To(raw.VnetSubnetID)
	}

	if raw.HostGroupID != "" {
		profile.HostGroupID = pointer.To(raw.HostGroupID)
	}

	if raw.OrchestratorVersion != "" {
		profile.OrchestratorVersion = pointer.To(raw.OrchestratorVersion)
	}

	if raw.ProximityPlacementGroupID != "" {
		profile.ProximityPlacementGroupID = pointer.To(raw.ProximityPlacementGroupID)
	}

	if raw.WorkloadRuntime != "" {
		profile.WorkloadRuntime = pointer.To(managedclusters.WorkloadRuntime(raw.WorkloadRuntime))
	}

	if raw.CapacityReservationGroupID != "" {
		profile.CapacityReservationGroupID = pointer.To(raw.CapacityReservationGroupID)
	}

	if raw.GPUInstance != "" {
		profile.GpuInstanceProfile = pointer.To(managedclusters.GPUInstanceProfile(raw.GPUInstance))
	}

	if raw.GPUDriver != "" {
		profile.GpuProfile = &managedclusters.GPUProfile{
			Driver: pointer.To(managedclusters.GPUDriver(raw.GPUDriver)),
		}
	}

	count := raw.NodeCount
	// maxCount := raw.MaxCount
	// minCount := raw.MinCount

	// Count must always be set (see #6094), RP behaviour has changed
	// since the API version upgrade in v2.1.0 making Count required
	// for all create/update requests
	profile.Count = pointer.To(count)

	// if enableAutoScaling {
	//	// if Count has not been set use min count
	//	if count == 0 {
	//		count = minCount
	//		profile.Count = pointer.To(count)
	//	}

	// // Count must be set for the initial creation when using AutoScaling but cannot be updated
	// if hasNodeCountChange && !isNewResource {
	// 	return nil, fmt.Errorf("cannot change `node_count` when `auto_scaling_enabled` is set to `true`")
	// }

	// if maxCount > 0 {
	// 	profile.MaxCount = pointer.To(maxCount)

	// 	if maxCount < count && isNewResource {
	// 		return nil, fmt.Errorf("`node_count`(%d) must be equal to or less than `max_count`(%d) when `auto_scaling_enabled` is set to `true`", count, maxCount)
	// 	}
	// } else {
	// 	return nil, fmt.Errorf("`max_count` must be configured when `auto_scaling_enabled` is set to `true`")
	// }

	// if minCount > 0 {
	// 	profile.MinCount = pointer.To(minCount)

	// 	if minCount > count && isNewResource {
	// 		return nil, fmt.Errorf("`node_count`(%d) must be equal to or greater than `min_count`(%d) when `auto_scaling_enabled` is set to `true`", count, minCount)
	// 	}
	// } else {
	// 	return nil, fmt.Errorf("`min_count` must be configured when `auto_scaling_enabled` is set to `true`")
	// }
	//
	//	if minCount > maxCount {
	//		return nil, fmt.Errorf("`max_count` must be >= `min_count`")
	//	}
	// } else if minCount > 0 || maxCount > 0 {
	// 	return nil, fmt.Errorf("`max_count`(%d) and `min_count`(%d) must be set to `null` when `auto_scaling_enabled` is set to `false`", maxCount, minCount)
	// }

	if len(raw.KubeletConfig) > 0 {
		profile.KubeletConfig = expandClusterNodePoolKubeletConfigTyped(raw.KubeletConfig)
	}

	if len(raw.LinuxOSConfig) > 0 {
		linuxOSConfig, err := expandAutomaticClusterNodePoolLinuxOSConfig(raw.LinuxOSConfig)
		if err != nil {
			return nil, err
		}
		profile.LinuxOSConfig = linuxOSConfig
	}

	if len(raw.NodeNetworkProfile) > 0 {
		profile.NetworkProfile = expandClusterPoolNetworkProfileTyped(raw.NodeNetworkProfile)
	}

	return &[]managedclusters.ManagedClusterAgentPoolProfile{
		profile,
	}, nil
}

func FlattenDefaultNodePoolTyped(input *[]managedclusters.ManagedClusterAgentPoolProfile, metadata *sdk.ResourceMetaData) ([]DefaultNodePoolModel, error) {
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

	// Preserve temporary_name_for_rotation from existing state since it's not returned by the API
	if metadata != nil {
		var existingModel KubernetesAutomaticClusterModel
		if err := metadata.Decode(&existingModel); err == nil {
			if len(existingModel.DefaultNodePool) > 0 {
				result.TemporaryNameForRotation = existingModel.DefaultNodePool[0].TemporaryNameForRotation
			}
		}
	}

	if agentPool.Count != nil {
		result.NodeCount = pointer.From(agentPool.Count)
	}

	if agentPool.EnableUltraSSD != nil {
		result.UltraSSDEnabled = pointer.From(agentPool.EnableUltraSSD)
	}
	//
	// if agentPool.EnableAutoScaling != nil {
	// 	result.AutoScalingEnabled = pointer.From(agentPool.EnableAutoScaling)
	// }

	if agentPool.EnableFIPS != nil {
		result.FipsEnabled = pointer.From(agentPool.EnableFIPS)
	}

	if agentPool.EnableNodePublicIP != nil {
		result.NodePublicIPEnabled = pointer.From(agentPool.EnableNodePublicIP)
	}

	if agentPool.EnableEncryptionAtHost != nil {
		result.HostEncryptionEnabled = pointer.From(agentPool.EnableEncryptionAtHost)
	}

	if agentPool.GpuInstanceProfile != nil {
		result.GPUInstance = string(pointer.From(agentPool.GpuInstanceProfile))
	}

	if agentPool.GpuProfile != nil && agentPool.GpuProfile.Driver != nil {
		result.GPUDriver = string(pointer.From(agentPool.GpuProfile.Driver))
	}

	// if agentPool.MaxCount != nil {
	// 	result.MaxCount = pointer.From(agentPool.MaxCount)
	// }

	if agentPool.MaxPods != nil {
		result.MaxPods = pointer.From(agentPool.MaxPods)
	}

	// if agentPool.MinCount != nil {
	// 	result.MinCount = pointer.From(agentPool.MinCount)
	// }

	if agentPool.NodeLabels != nil {
		result.NodeLabels = make(map[string]string)
		for k, v := range pointer.From(agentPool.NodeLabels) {
			result.NodeLabels[k] = v
		}
	}

	if agentPool.NodePublicIPPrefixID != nil {
		result.NodePublicIPPrefixID = pointer.From(agentPool.NodePublicIPPrefixID)
	}

	// Check for CriticalAddonsOnly taint
	if agentPool.NodeTaints != nil {
		for _, taint := range pointer.From(agentPool.NodeTaints) {
			if taint == "CriticalAddonsOnly=true:NoSchedule" {
				result.OnlyCriticalAddonsEnabled = true
				break
			}
		}
	}

	if agentPool.OsDiskSizeGB != nil {
		result.OSDiskSizeGB = pointer.From(agentPool.OsDiskSizeGB)
	}

	// if agentPool.OsDiskType != nil {
	// 	result.OSDiskType = string(pointer.From(agentPool.OsDiskType))
	// } else {
	// 	result.OSDiskType = string(managedclusters.OSDiskTypeManaged)
	//}

	if agentPool.PodSubnetID != nil {
		result.PodSubnetID = pointer.From(agentPool.PodSubnetID)
	}

	if agentPool.VnetSubnetID != nil {
		result.VnetSubnetID = pointer.From(agentPool.VnetSubnetID)
	}

	if agentPool.HostGroupID != nil {
		result.HostGroupID = pointer.From(agentPool.HostGroupID)
	}

	// NOTE: workaround for migration from 2022-01-02-preview (<3.12.0) to 2022-03-02-preview (>=3.12.0)
	// Before terraform apply is run against the new API, Azure will respond only with currentOrchestratorVersion
	if agentPool.OrchestratorVersion != nil {
		result.OrchestratorVersion = pointer.From(agentPool.OrchestratorVersion)
	} else if agentPool.CurrentOrchestratorVersion != nil {
		result.OrchestratorVersion = pointer.From(agentPool.CurrentOrchestratorVersion)
	}

	if agentPool.ProximityPlacementGroupID != nil {
		result.ProximityPlacementGroupID = pointer.From(agentPool.ProximityPlacementGroupID)
	}

	// if agentPool.ScaleDownMode != nil {
	// 	result.ScaleDownMode = string(pointer.From(agentPool.ScaleDownMode))
	// } else {
	// 	result.ScaleDownMode = string(managedclusters.ScaleDownModeDelete)
	//}

	if agentPool.CreationData != nil && agentPool.CreationData.SourceResourceId != nil {
		id, err := snapshots.ParseSnapshotIDInsensitively(pointer.From(agentPool.CreationData.SourceResourceId))
		if err != nil {
			return nil, err
		}
		result.SnapshotID = id.ID()
	}

	if agentPool.VMSize != nil {
		result.VMSize = pointer.From(agentPool.VMSize)
	}

	if agentPool.CapacityReservationGroupID != nil {
		result.CapacityReservationGroupID = pointer.From(agentPool.CapacityReservationGroupID)
	}

	if agentPool.WorkloadRuntime != nil {
		result.WorkloadRuntime = string(pointer.From(agentPool.WorkloadRuntime))
	}

	if agentPool.KubeletDiskType != nil {
		result.KubeletDiskType = string(pointer.From(agentPool.KubeletDiskType))
	}

	if agentPool.OsSKU != nil {
		result.OSSKU = string(pointer.From(agentPool.OsSKU))
	}

	if agentPool.Type != nil {
		result.Type = string(pointer.From(agentPool.Type))
	}

	result.UpgradeSettings = flattenClusterNodePoolUpgradeSettingsTyped(agentPool.UpgradeSettings)

	result.LinuxOSConfig = flattenAutomaticClusterNodePoolLinuxOSConfig(agentPool.LinuxOSConfig)

	result.KubeletConfig = flattenClusterNodePoolKubeletConfigTyped(agentPool.KubeletConfig)
	result.NodeNetworkProfile = flattenClusterPoolNetworkProfileTyped(agentPool.NetworkProfile)

	// if agentPool.AvailabilityZones != nil {
	// 	result.Zones = pointer.From(agentPool.AvailabilityZones)
	// }

	if agentPool.Tags != nil {
		result.Tags = tags.Flatten(agentPool.Tags)
	}

	return []DefaultNodePoolModel{result}, nil
}
