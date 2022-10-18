package containers

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerservice/mgmt/2022-03-02-preview/containerservice"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-08-02-preview/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-08-02-preview/managedclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func SchemaDefaultNodePool() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: func() map[string]*pluginsdk.Schema {
				s := map[string]*pluginsdk.Schema{
					// Required
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.KubernetesAgentPoolName,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  string(containerservice.AgentPoolTypeVirtualMachineScaleSets),
						ValidateFunc: validation.StringInSlice([]string{
							string(containerservice.AgentPoolTypeAvailabilitySet),
							string(containerservice.AgentPoolTypeVirtualMachineScaleSets),
						}, false),
					},

					"vm_size": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"capacity_reservation_group_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: computeValidate.CapacityReservationGroupID,
					},

					// TODO 4.0: change this from enable_* to *_enabled
					"enable_auto_scaling": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					// TODO 4.0: change this from enable_* to *_enabled
					"enable_node_public_ip": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
					},

					// TODO 4.0: change this from enable_* to *_enabled
					"enable_host_encryption": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
					},

					"kubelet_config": schemaNodePoolKubeletConfig(),

					"linux_os_config": schemaNodePoolLinuxOSConfig(),

					"fips_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
					},

					"kubelet_disk_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(containerservice.KubeletDiskTypeOS),
							string(containerservice.KubeletDiskTypeTemporary),
						}, false),
					},

					"max_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						// NOTE: rather than setting `0` users should instead pass `null` here
						ValidateFunc: validation.IntBetween(1, 1000),
					},

					"max_pods": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},

					"message_of_the_day": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"min_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						// NOTE: rather than setting `0` users should instead pass `null` here
						ValidateFunc: validation.IntBetween(1, 1000),
					},

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
						ValidateFunc: azure.ValidateResourceID,
						RequiredWith: []string{"default_node_pool.0.enable_node_public_ip"},
					},

					"node_taints": {
						Type:     pluginsdk.TypeList,
						ForceNew: true,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"tags": commonschema.Tags(),

					"os_disk_size_gb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ForceNew:     true,
						Computed:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"os_disk_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  containerservice.OSDiskTypeManaged,
						ValidateFunc: validation.StringInSlice([]string{
							string(containerservice.OSDiskTypeEphemeral),
							string(containerservice.OSDiskTypeManaged),
						}, false),
					},

					"os_sku": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Computed: true, // defaults to Ubuntu if using Linux
						ValidateFunc: validation.StringInSlice([]string{
							string(containerservice.OSSKUUbuntu),
							string(containerservice.OSSKUCBLMariner),
						}, false),
					},

					"ultra_ssd_enabled": {
						Type:     pluginsdk.TypeBool,
						ForceNew: true,
						Default:  false,
						Optional: true,
					},

					"vnet_subnet_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: azure.ValidateResourceID,
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
						ForceNew:     true,
						ValidateFunc: networkValidate.SubnetID,
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
						ForceNew: true,
					},

					"scale_down_mode": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  string(containerservice.ScaleDownModeDelete),
						ValidateFunc: validation.StringInSlice([]string{
							string(containerservice.ScaleDownModeDeallocate),
							string(containerservice.ScaleDownModeDelete),
						}, false),
					},

					"host_group_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: computeValidate.HostGroupID,
					},

					"upgrade_settings": upgradeSettingsSchema(),

					"workload_runtime": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(containerservice.WorkloadRuntimeOCIContainer),
						}, false),
					},
				}

				s["zones"] = commonschema.ZonesMultipleOptionalForceNew()

				return s
			}(),
		},
	}
}

func schemaNodePoolKubeletConfig() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"cpu_manager_policy": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						"none",
						"static",
					}, false),
				},

				"cpu_cfs_quota_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
				},

				"cpu_cfs_quota_period": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
				},

				"image_gc_high_threshold": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(0, 100),
				},

				"image_gc_low_threshold": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(0, 100),
				},

				"topology_manager_policy": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
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
					ForceNew: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"container_log_max_size_mb": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					ForceNew: true,
				},

				"container_log_max_line": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntAtLeast(2),
				},

				"pod_max_pid": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					ForceNew: true,
				},
			},
		},
	}
}

func schemaNodePoolLinuxOSConfig() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"sysctl_config": schemaNodePoolSysctlConfig(),

				"transparent_huge_page_enabled": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						"always",
						"madvise",
						"never",
					}, false),
				},

				"transparent_huge_page_defrag": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
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
					ForceNew: true,
				},
			},
		},
	}
}

func schemaNodePoolSysctlConfig() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"fs_aio_max_nr": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(65536, 6553500),
				},

				"fs_file_max": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(8192, 12000500),
				},

				"fs_inotify_max_user_watches": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(781250, 2097152),
				},

				"fs_nr_open": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(8192, 20000500),
				},

				"kernel_threads_max": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(20, 513785),
				},

				"net_core_netdev_max_backlog": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(1000, 3240000),
				},

				"net_core_optmem_max": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(20480, 4194304),
				},

				"net_core_rmem_default": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(212992, 134217728),
				},

				"net_core_rmem_max": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(212992, 134217728),
				},

				"net_core_somaxconn": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(4096, 3240000),
				},

				"net_core_wmem_default": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(212992, 134217728),
				},

				"net_core_wmem_max": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(212992, 134217728),
				},

				"net_ipv4_ip_local_port_range_min": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(1024, 60999),
				},

				"net_ipv4_ip_local_port_range_max": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(32768, 65000),
				},

				"net_ipv4_neigh_default_gc_thresh1": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(128, 80000),
				},

				"net_ipv4_neigh_default_gc_thresh2": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(512, 90000),
				},

				"net_ipv4_neigh_default_gc_thresh3": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(1024, 100000),
				},

				"net_ipv4_tcp_fin_timeout": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(5, 120),
				},

				"net_ipv4_tcp_keepalive_intvl": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(10, 75),
				},

				"net_ipv4_tcp_keepalive_probes": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(1, 15),
				},

				"net_ipv4_tcp_keepalive_time": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(30, 432000),
				},

				"net_ipv4_tcp_max_syn_backlog": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(128, 3240000),
				},

				"net_ipv4_tcp_max_tw_buckets": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(8000, 1440000),
				},

				"net_ipv4_tcp_tw_reuse": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
				},

				"net_netfilter_nf_conntrack_buckets": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(65536, 147456),
				},

				"net_netfilter_nf_conntrack_max": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(131072, 1048576),
				},

				"vm_max_map_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(65530, 262144),
				},

				"vm_swappiness": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(0, 100),
				},

				"vm_vfs_cache_pressure": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(0, 100),
				},
			},
		},
	}
}

func ConvertDefaultNodePoolToAgentPool(input *[]managedclusters.ManagedClusterAgentPoolProfile) agentpools.AgentPool {
	defaultCluster := (*input)[0]
	osDiskType := agentpools.OSDiskType(*defaultCluster.OsDiskType)
	osType := agentpools.OSType(*defaultCluster.OsType)
	kubeletDiskType := agentpools.KubeletDiskType(*defaultCluster.KubeletDiskType)
	poolType := agentpools.AgentPoolType(*defaultCluster.Type)
	scaleSetPriority := agentpools.ScaleSetPriority(*defaultCluster.ScaleSetPriority)
	scaleSetEvictionPolicy := agentpools.ScaleSetEvictionPolicy(*defaultCluster.ScaleSetEvictionPolicy)
	mode := agentpools.AgentPoolMode(*defaultCluster.Mode)
	scaleDownMode := agentpools.ScaleDownMode(*defaultCluster.ScaleDownMode)
	upgradeSettings := agentpools.AgentPoolUpgradeSettings{
		MaxSurge: defaultCluster.UpgradeSettings.MaxSurge,
	}
	workloadRuntime := agentpools.WorkloadRuntime(*defaultCluster.WorkloadRuntime)

	kubeletConfig := agentpools.KubeletConfig{
		AllowedUnsafeSysctls:  defaultCluster.KubeletConfig.AllowedUnsafeSysctls,
		ContainerLogMaxFiles:  defaultCluster.KubeletConfig.ContainerLogMaxFiles,
		ContainerLogMaxSizeMB: defaultCluster.KubeletConfig.ContainerLogMaxSizeMB,
		CpuCfsQuota:           defaultCluster.KubeletConfig.CpuCfsQuota,
		CpuCfsQuotaPeriod:     defaultCluster.KubeletConfig.CpuCfsQuotaPeriod,
		CpuManagerPolicy:      defaultCluster.KubeletConfig.CpuManagerPolicy,
		FailSwapOn:            defaultCluster.KubeletConfig.FailSwapOn,
		ImageGcHighThreshold:  defaultCluster.KubeletConfig.ImageGcHighThreshold,
		ImageGcLowThreshold:   defaultCluster.KubeletConfig.ImageGcLowThreshold,
		PodMaxPids:            defaultCluster.KubeletConfig.PodMaxPids,
		TopologyManagerPolicy: defaultCluster.KubeletConfig.TopologyManagerPolicy,
	}

	linuxOSConfig := agentpools.LinuxOSConfig{
		SwapFileSizeMB: defaultCluster.LinuxOSConfig.SwapFileSizeMB,
		Sysctls: &agentpools.SysctlConfig{
			FsAioMaxNr:                     defaultCluster.LinuxOSConfig.Sysctls.FsAioMaxNr,
			FsFileMax:                      defaultCluster.LinuxOSConfig.Sysctls.FsFileMax,
			FsInotifyMaxUserWatches:        defaultCluster.LinuxOSConfig.Sysctls.FsInotifyMaxUserWatches,
			FsNrOpen:                       defaultCluster.LinuxOSConfig.Sysctls.FsNrOpen,
			KernelThreadsMax:               defaultCluster.LinuxOSConfig.Sysctls.KernelThreadsMax,
			NetCoreNetdevMaxBacklog:        defaultCluster.LinuxOSConfig.Sysctls.NetCoreNetdevMaxBacklog,
			NetCoreOptmemMax:               defaultCluster.LinuxOSConfig.Sysctls.NetCoreSomaxconn,
			NetCoreRmemDefault:             defaultCluster.LinuxOSConfig.Sysctls.NetCoreRmemDefault,
			NetCoreRmemMax:                 defaultCluster.LinuxOSConfig.Sysctls.NetCoreSomaxconn,
			NetCoreSomaxconn:               defaultCluster.LinuxOSConfig.Sysctls.NetCoreNetdevMaxBacklog,
			NetCoreWmemDefault:             defaultCluster.LinuxOSConfig.Sysctls.NetCoreWmemDefault,
			NetCoreWmemMax:                 defaultCluster.LinuxOSConfig.Sysctls.NetCoreNetdevMaxBacklog,
			NetIPv4IPLocalPortRange:        defaultCluster.LinuxOSConfig.Sysctls.NetIPv4IPLocalPortRange,
			NetIPv4NeighDefaultGcThresh1:   defaultCluster.LinuxOSConfig.Sysctls.NetIPv4NeighDefaultGcThresh1,
			NetIPv4NeighDefaultGcThresh2:   defaultCluster.LinuxOSConfig.Sysctls.NetIPv4NeighDefaultGcThresh2,
			NetIPv4NeighDefaultGcThresh3:   defaultCluster.LinuxOSConfig.Sysctls.NetIPv4NeighDefaultGcThresh3,
			NetIPv4TcpFinTimeout:           defaultCluster.LinuxOSConfig.Sysctls.NetIPv4TcpFinTimeout,
			NetIPv4TcpKeepaliveProbes:      defaultCluster.LinuxOSConfig.Sysctls.NetIPv4TcpKeepaliveProbes,
			NetIPv4TcpKeepaliveTime:        defaultCluster.LinuxOSConfig.Sysctls.NetIPv4TcpKeepaliveTime,
			NetIPv4TcpMaxSynBacklog:        defaultCluster.LinuxOSConfig.Sysctls.NetIPv4TcpMaxSynBacklog,
			NetIPv4TcpMaxTwBuckets:         defaultCluster.LinuxOSConfig.Sysctls.NetIPv4TcpMaxTwBuckets,
			NetIPv4TcpTwReuse:              defaultCluster.LinuxOSConfig.Sysctls.NetIPv4TcpTwReuse,
			NetIPv4TcpkeepaliveIntvl:       defaultCluster.LinuxOSConfig.Sysctls.NetIPv4TcpkeepaliveIntvl,
			NetNetfilterNfConntrackBuckets: defaultCluster.LinuxOSConfig.Sysctls.NetNetfilterNfConntrackBuckets,
			NetNetfilterNfConntrackMax:     defaultCluster.LinuxOSConfig.Sysctls.NetNetfilterNfConntrackMax,
			VmMaxMapCount:                  defaultCluster.LinuxOSConfig.Sysctls.VmMaxMapCount,
			VmSwappiness:                   defaultCluster.LinuxOSConfig.Sysctls.VmSwappiness,
			VmVfsCachePressure:             defaultCluster.LinuxOSConfig.Sysctls.VmVfsCachePressure,
		},
		TransparentHugePageDefrag:  defaultCluster.LinuxOSConfig.TransparentHugePageDefrag,
		TransparentHugePageEnabled: defaultCluster.LinuxOSConfig.TransparentHugePageEnabled,
	}
	return agentpools.AgentPool{
		Name: &defaultCluster.Name,
		Properties: &agentpools.ManagedClusterAgentPoolProfileProperties{
			Count:                     defaultCluster.Count,
			VmSize:                    defaultCluster.VmSize,
			OsDiskSizeGB:              defaultCluster.OsDiskSizeGB,
			OsDiskType:                &osDiskType,
			VnetSubnetID:              defaultCluster.VnetSubnetID,
			KubeletConfig:             &kubeletConfig,
			LinuxOSConfig:             &linuxOSConfig,
			MaxPods:                   defaultCluster.MaxPods,
			OsType:                    &osType,
			MaxCount:                  defaultCluster.MaxCount,
			MessageOfTheDay:           defaultCluster.MessageOfTheDay,
			MinCount:                  defaultCluster.MinCount,
			EnableAutoScaling:         defaultCluster.EnableAutoScaling,
			EnableFIPS:                defaultCluster.EnableFIPS,
			KubeletDiskType:           &kubeletDiskType,
			Type:                      &poolType,
			OrchestratorVersion:       defaultCluster.OrchestratorVersion,
			ProximityPlacementGroupID: defaultCluster.ProximityPlacementGroupID,
			AvailabilityZones:         defaultCluster.AvailabilityZones,
			EnableNodePublicIP:        defaultCluster.EnableNodePublicIP,
			NodePublicIPPrefixID:      defaultCluster.NodePublicIPPrefixID,
			ScaleSetPriority:          &scaleSetPriority,
			ScaleSetEvictionPolicy:    &scaleSetEvictionPolicy,
			SpotMaxPrice:              defaultCluster.SpotMaxPrice,
			Mode:                      &mode,
			NodeLabels:                defaultCluster.NodeLabels,
			NodeTaints:                defaultCluster.NodeTaints,
			PodSubnetID:               defaultCluster.PodSubnetID,
			ScaleDownMode:             &scaleDownMode,
			Tags:                      defaultCluster.Tags,
			UpgradeSettings:           &upgradeSettings,
			WorkloadRuntime:           &workloadRuntime,
		},
	}
}

func ExpandDefaultNodePool(d *pluginsdk.ResourceData) (*[]managedclusters.ManagedClusterAgentPoolProfile, error) {
	input := d.Get("default_node_pool").([]interface{})

	raw := input[0].(map[string]interface{})
	enableAutoScaling := raw["enable_auto_scaling"].(bool)
	nodeLabelsRaw := raw["node_labels"].(map[string]interface{})
	nodeLabels := make(map[string]string)
	for k, v := range nodeLabelsRaw {
		nodeLabels[k] = v.(string)
	}
	nodeTaintsRaw := raw["node_taints"].([]interface{})
	nodeTaints := utils.ExpandStringSlice(nodeTaintsRaw)

	if len(*nodeTaints) != 0 {
		return nil, fmt.Errorf("The AKS API has removed support for tainting all nodes in the default node pool and it is no longer possible to configure this. To taint a node pool, create a separate one.")
	}

	criticalAddonsEnabled := raw["only_critical_addons_enabled"].(bool)
	if criticalAddonsEnabled {
		*nodeTaints = append(*nodeTaints, "CriticalAddonsOnly=true:NoSchedule")
	}

	t := raw["tags"].(map[string]interface{})

	apType := managedclusters.AgentPoolType(raw["type"].(string))
	kubeletDiskType := raw["kubelet_disk_type"].(managedclusters.KubeletDiskType)
	osType := managedclusters.OSTypeLinux
	mode := managedclusters.AgentPoolModeSystem

	agentPool := managedclusters.ManagedClusterAgentPoolProfile{
		EnableAutoScaling:      utils.Bool(enableAutoScaling),
		EnableFIPS:             utils.Bool(raw["fips_enabled"].(bool)),
		EnableNodePublicIP:     utils.Bool(raw["enable_node_public_ip"].(bool)),
		EnableEncryptionAtHost: utils.Bool(raw["enable_host_encryption"].(bool)),
		KubeletDiskType:        &kubeletDiskType,

		NodeLabels: &nodeLabels,
		NodeTaints: nodeTaints,
		Tags:       tags.Expand(t),
		Type:       &apType,
		VmSize:     utils.String(raw["vm_size"].(string)),

		// at this time the default node pool has to be Linux or the AKS cluster fails to provision with:
		// Pods not in Running status: coredns-7fc597cc45-v5z7x,coredns-autoscaler-7ccc76bfbd-djl7j,metrics-server-cbd95f966-5rl97,tunnelfront-7d9884977b-wpbvn
		// Windows agents can be configured via the separate node pool resource
		OsType: &osType,

		// without this set the API returns:
		// Code="MustDefineAtLeastOneSystemPool" Message="Must define at least one system pool."
		// since this is the "default" node pool we can assume this is a system node pool
		Mode: &mode,

		UpgradeSettings: expandUpgradeSettingsForCluster(raw["upgrade_settings"].([]interface{})),

		// // TODO: support these in time
		// ScaleSetEvictionPolicy: "",
		// ScaleSetPriority:       "",
	}

	zones := zones.Expand(raw["zones"].(*schema.Set).List())
	if len(zones) > 0 {
		agentPool.AvailabilityZones = &zones
	}

	if maxPods := int64(raw["max_pods"].(int)); maxPods > 0 {
		agentPool.MaxPods = utils.Int64(maxPods)
	}

	if v := raw["message_of_the_day"].(string); v != "" {
		messageOfTheDayEncoded := base64.StdEncoding.EncodeToString([]byte(v))
		agentPool.MessageOfTheDay = &messageOfTheDayEncoded
	}

	if prefixID := raw["node_public_ip_prefix_id"].(string); prefixID != "" {
		agentPool.NodePublicIPPrefixID = utils.String(prefixID)
	}

	if osDiskSizeGB := int64(raw["os_disk_size_gb"].(int)); osDiskSizeGB > 0 {
		agentPool.OsDiskSizeGB = utils.Int64(osDiskSizeGB)
	}

	*agentPool.OsDiskType = managedclusters.OSDiskTypeManaged
	if osDiskType := raw["os_disk_type"].(string); osDiskType != "" {
		*agentPool.OsDiskType = managedclusters.OSDiskType(raw["os_disk_type"].(string))
	}

	if osSku := raw["os_sku"].(string); osSku != "" {
		*agentPool.OsSKU = managedclusters.OSSKU(osSku)
	}

	if podSubnetID := raw["pod_subnet_id"].(string); podSubnetID != "" {
		agentPool.PodSubnetID = utils.String(podSubnetID)
	}

	*agentPool.ScaleDownMode = managedclusters.ScaleDownModeDelete
	if scaleDownMode := raw["scale_down_mode"].(string); scaleDownMode != "" {
		*agentPool.ScaleDownMode = managedclusters.ScaleDownMode(scaleDownMode)
	}

	if ultraSSDEnabled, ok := raw["ultra_ssd_enabled"]; ok {
		agentPool.EnableUltraSSD = utils.Bool(ultraSSDEnabled.(bool))
	}

	if vnetSubnetID := raw["vnet_subnet_id"].(string); vnetSubnetID != "" {
		agentPool.VnetSubnetID = utils.String(vnetSubnetID)
	}

	if hostGroupID := raw["host_group_id"].(string); hostGroupID != "" {
		agentPool.HostGroupID = utils.String(hostGroupID)
	}

	if orchestratorVersion := raw["orchestrator_version"].(string); orchestratorVersion != "" {
		agentPool.OrchestratorVersion = utils.String(orchestratorVersion)
	}

	if proximityPlacementGroupId := raw["proximity_placement_group_id"].(string); proximityPlacementGroupId != "" {
		agentPool.ProximityPlacementGroupID = utils.String(proximityPlacementGroupId)
	}

	if workloadRunTime := raw["workload_runtime"].(string); workloadRunTime != "" {
		*agentPool.WorkloadRuntime = managedclusters.WorkloadRuntime(workloadRunTime)
	}

	if capacityReservationGroupId := raw["capacity_reservation_group_id"].(string); capacityReservationGroupId != "" {
		agentPool.CapacityReservationGroupID = utils.String(capacityReservationGroupId)
	}

	count := raw["node_count"].(int)
	maxCount := raw["max_count"].(int)
	minCount := raw["min_count"].(int)

	// Count must always be set (see #6094), RP behaviour has changed
	// since the API version upgrade in v2.1.0 making Count required
	// for all create/update requests
	agentPool.Count = utils.Int64(int64(count))

	if enableAutoScaling {
		// if Count has not been set use min count
		if count == 0 {
			count = minCount
			agentPool.Count = utils.Int64(int64(count))
		}

		// Count must be set for the initial creation when using AutoScaling but cannot be updated
		if d.HasChange("default_node_pool.0.node_count") && !d.IsNewResource() {
			return nil, fmt.Errorf("cannot change `node_count` when `enable_auto_scaling` is set to `true`")
		}

		if maxCount > 0 {
			agentPool.MaxCount = utils.Int64(int64(maxCount))
			if maxCount < count {
				return nil, fmt.Errorf("`node_count`(%d) must be equal to or less than `max_count`(%d) when `enable_auto_scaling` is set to `true`", count, maxCount)
			}
		} else {
			return nil, fmt.Errorf("`max_count` must be configured when `enable_auto_scaling` is set to `true`")
		}

		if minCount > 0 {
			agentPool.MinCount = utils.Int64(int64(minCount))

			if minCount > count && d.IsNewResource() {
				return nil, fmt.Errorf("`node_count`(%d) must be equal to or greater than `min_count`(%d) when `enable_auto_scaling` is set to `true`", count, minCount)
			}
		} else {
			return nil, fmt.Errorf("`min_count` must be configured when `enable_auto_scaling` is set to `true`")
		}

		if minCount > maxCount {
			return nil, fmt.Errorf("`max_count` must be >= `min_count`")
		}
	} else if minCount > 0 || maxCount > 0 {
		return nil, fmt.Errorf("`max_count`(%d) and `min_count`(%d) must be set to `null` when `enable_auto_scaling` is set to `false`", maxCount, minCount)
	}

	if kubeletConfig := raw["kubelet_config"].([]interface{}); len(kubeletConfig) > 0 {
		agentPool.KubeletConfig = expandAgentPoolKubeletConfigForCluster(kubeletConfig)
	}

	if linuxOSConfig := raw["linux_os_config"].([]interface{}); len(linuxOSConfig) > 0 {
		linuxOSConfig, err := expandAgentPoolLinuxOSConfigForCluster(linuxOSConfig)
		if err != nil {
			return nil, err
		}
		agentPool.LinuxOSConfig = linuxOSConfig
	}

	return &[]managedclusters.ManagedClusterAgentPoolProfile{
		agentPool,
	}, nil
}

func expandAgentPoolKubeletConfig(input []interface{}) *agentpools.KubeletConfig {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	result := &agentpools.KubeletConfig{
		CpuCfsQuota: utils.Bool(raw["cpu_cfs_quota_enabled"].(bool)),
		// must be false, otherwise the backend will report error: CustomKubeletConfig.FailSwapOn must be set to false to enable swap file on nodes.
		FailSwapOn:           utils.Bool(false),
		AllowedUnsafeSysctls: utils.ExpandStringSlice(raw["allowed_unsafe_sysctls"].(*pluginsdk.Set).List()),
	}

	if v := raw["cpu_manager_policy"].(string); v != "" {
		result.CpuManagerPolicy = utils.String(v)
	}
	if v := raw["cpu_cfs_quota_period"].(string); v != "" {
		result.CpuCfsQuotaPeriod = utils.String(v)
	}
	if v := raw["image_gc_high_threshold"].(int); v != 0 {
		result.ImageGcHighThreshold = utils.Int64(int64(v))
	}
	if v := raw["image_gc_low_threshold"].(int); v != 0 {
		result.ImageGcLowThreshold = utils.Int64(int64(v))
	}
	if v := raw["topology_manager_policy"].(string); v != "" {
		result.TopologyManagerPolicy = utils.String(v)
	}
	if v := raw["container_log_max_size_mb"].(int); v != 0 {
		result.ContainerLogMaxSizeMB = utils.Int64(int64(v))
	}
	if v := raw["container_log_max_line"].(int); v != 0 {
		result.ContainerLogMaxFiles = utils.Int64(int64(v))
	}
	if v := raw["pod_max_pid"].(int); v != 0 {
		result.PodMaxPids = utils.Int64(int64(v))
	}

	return result
}

func expandAgentPoolKubeletConfigForCluster(input []interface{}) *managedclusters.KubeletConfig {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	result := &managedclusters.KubeletConfig{
		CpuCfsQuota: utils.Bool(raw["cpu_cfs_quota_enabled"].(bool)),
		// must be false, otherwise the backend will report error: CustomKubeletConfig.FailSwapOn must be set to false to enable swap file on nodes.
		FailSwapOn:           utils.Bool(false),
		AllowedUnsafeSysctls: utils.ExpandStringSlice(raw["allowed_unsafe_sysctls"].(*pluginsdk.Set).List()),
	}

	if v := raw["cpu_manager_policy"].(string); v != "" {
		result.CpuManagerPolicy = utils.String(v)
	}
	if v := raw["cpu_cfs_quota_period"].(string); v != "" {
		result.CpuCfsQuotaPeriod = utils.String(v)
	}
	if v := raw["image_gc_high_threshold"].(int); v != 0 {
		result.ImageGcHighThreshold = utils.Int64(int64(v))
	}
	if v := raw["image_gc_low_threshold"].(int); v != 0 {
		result.ImageGcLowThreshold = utils.Int64(int64(v))
	}
	if v := raw["topology_manager_policy"].(string); v != "" {
		result.TopologyManagerPolicy = utils.String(v)
	}
	if v := raw["container_log_max_size_mb"].(int); v != 0 {
		result.ContainerLogMaxSizeMB = utils.Int64(int64(v))
	}
	if v := raw["container_log_max_line"].(int); v != 0 {
		result.ContainerLogMaxFiles = utils.Int64(int64(v))
	}
	if v := raw["pod_max_pid"].(int); v != 0 {
		result.PodMaxPids = utils.Int64(int64(v))
	}

	return result
}

func expandAgentPoolLinuxOSConfig(input []interface{}) (*agentpools.LinuxOSConfig, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	raw := input[0].(map[string]interface{})
	sysctlConfig, err := expandAgentPoolSysctlConfig(raw["sysctl_config"].([]interface{}))
	if err != nil {
		return nil, err
	}

	result := &agentpools.LinuxOSConfig{
		Sysctls: sysctlConfig,
	}
	if v := raw["transparent_huge_page_enabled"].(string); v != "" {
		result.TransparentHugePageEnabled = utils.String(v)
	}
	if v := raw["transparent_huge_page_defrag"].(string); v != "" {
		result.TransparentHugePageDefrag = utils.String(v)
	}
	if v := raw["swap_file_size_mb"].(int); v != 0 {
		result.SwapFileSizeMB = utils.Int64(int64(v))
	}
	return result, nil
}

func expandAgentPoolSysctlConfig(input []interface{}) (*agentpools.SysctlConfig, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	raw := input[0].(map[string]interface{})
	result := &agentpools.SysctlConfig{
		NetIPv4TcpTwReuse: utils.Bool(raw["net_ipv4_tcp_tw_reuse"].(bool)),
	}
	if v := raw["net_core_somaxconn"].(int); v != 0 {
		result.NetCoreSomaxconn = utils.Int64(int64(v))
	}
	if v := raw["net_core_netdev_max_backlog"].(int); v != 0 {
		result.NetCoreNetdevMaxBacklog = utils.Int64(int64(v))
	}
	if v := raw["net_core_rmem_default"].(int); v != 0 {
		result.NetCoreRmemDefault = utils.Int64(int64(v))
	}
	if v := raw["net_core_rmem_max"].(int); v != 0 {
		result.NetCoreRmemMax = utils.Int64(int64(v))
	}
	if v := raw["net_core_wmem_default"].(int); v != 0 {
		result.NetCoreWmemDefault = utils.Int64(int64(v))
	}
	if v := raw["net_core_wmem_max"].(int); v != 0 {
		result.NetCoreWmemMax = utils.Int64(int64(v))
	}
	if v := raw["net_core_optmem_max"].(int); v != 0 {
		result.NetCoreOptmemMax = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_max_syn_backlog"].(int); v != 0 {
		result.NetIPv4TcpMaxSynBacklog = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_max_tw_buckets"].(int); v != 0 {
		result.NetIPv4TcpMaxTwBuckets = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_fin_timeout"].(int); v != 0 {
		result.NetIPv4TcpFinTimeout = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_keepalive_time"].(int); v != 0 {
		result.NetIPv4TcpKeepaliveTime = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_keepalive_probes"].(int); v != 0 {
		result.NetIPv4TcpKeepaliveProbes = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_keepalive_intvl"].(int); v != 0 {
		result.NetIPv4TcpkeepaliveIntvl = utils.Int64(int64(v))
	}
	netIPv4IPLocalPortRangeMin := raw["net_ipv4_ip_local_port_range_min"].(int)
	netIPv4IPLocalPortRangeMax := raw["net_ipv4_ip_local_port_range_max"].(int)
	if (netIPv4IPLocalPortRangeMin != 0 && netIPv4IPLocalPortRangeMax == 0) || (netIPv4IPLocalPortRangeMin == 0 && netIPv4IPLocalPortRangeMax != 0) {
		return nil, fmt.Errorf("`net_ipv4_ip_local_port_range_min` and `net_ipv4_ip_local_port_range_max` should both be set or unset")
	}
	if netIPv4IPLocalPortRangeMin > netIPv4IPLocalPortRangeMax {
		return nil, fmt.Errorf("`net_ipv4_ip_local_port_range_min` should be no larger than `net_ipv4_ip_local_port_range_max`")
	}
	if netIPv4IPLocalPortRangeMin != 0 && netIPv4IPLocalPortRangeMax != 0 {
		result.NetIPv4IPLocalPortRange = utils.String(fmt.Sprintf("%d %d", netIPv4IPLocalPortRangeMin, netIPv4IPLocalPortRangeMax))
	}
	if v := raw["net_ipv4_neigh_default_gc_thresh1"].(int); v != 0 {
		result.NetIPv4NeighDefaultGcThresh1 = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_neigh_default_gc_thresh2"].(int); v != 0 {
		result.NetIPv4NeighDefaultGcThresh2 = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_neigh_default_gc_thresh3"].(int); v != 0 {
		result.NetIPv4NeighDefaultGcThresh3 = utils.Int64(int64(v))
	}
	if v := raw["net_netfilter_nf_conntrack_max"].(int); v != 0 {
		result.NetNetfilterNfConntrackMax = utils.Int64(int64(v))
	}
	if v := raw["net_netfilter_nf_conntrack_buckets"].(int); v != 0 {
		result.NetNetfilterNfConntrackBuckets = utils.Int64(int64(v))
	}
	if v := raw["fs_aio_max_nr"].(int); v != 0 {
		result.FsAioMaxNr = utils.Int64(int64(v))
	}
	if v := raw["fs_inotify_max_user_watches"].(int); v != 0 {
		result.FsInotifyMaxUserWatches = utils.Int64(int64(v))
	}
	if v := raw["fs_file_max"].(int); v != 0 {
		result.FsFileMax = utils.Int64(int64(v))
	}
	if v := raw["fs_nr_open"].(int); v != 0 {
		result.FsNrOpen = utils.Int64(int64(v))
	}
	if v := raw["kernel_threads_max"].(int); v != 0 {
		result.KernelThreadsMax = utils.Int64(int64(v))
	}
	if v := raw["vm_max_map_count"].(int); v != 0 {
		result.VmMaxMapCount = utils.Int64(int64(v))
	}
	if v := raw["vm_swappiness"].(int); v != 0 {
		result.VmSwappiness = utils.Int64(int64(v))
	}
	if v := raw["vm_vfs_cache_pressure"].(int); v != 0 {
		result.VmVfsCachePressure = utils.Int64(int64(v))
	}
	return result, nil
}

func expandAgentPoolLinuxOSConfigForCluster(input []interface{}) (*managedclusters.LinuxOSConfig, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	raw := input[0].(map[string]interface{})
	sysctlConfig, err := expandAgentPoolSysctlConfigForCluster(raw["sysctl_config"].([]interface{}))
	if err != nil {
		return nil, err
	}

	result := &managedclusters.LinuxOSConfig{
		Sysctls: sysctlConfig,
	}
	if v := raw["transparent_huge_page_enabled"].(string); v != "" {
		result.TransparentHugePageEnabled = utils.String(v)
	}
	if v := raw["transparent_huge_page_defrag"].(string); v != "" {
		result.TransparentHugePageDefrag = utils.String(v)
	}
	if v := raw["swap_file_size_mb"].(int); v != 0 {
		result.SwapFileSizeMB = utils.Int64(int64(v))
	}
	return result, nil
}

func expandAgentPoolSysctlConfigForCluster(input []interface{}) (*managedclusters.SysctlConfig, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	raw := input[0].(map[string]interface{})
	result := &managedclusters.SysctlConfig{
		NetIPv4TcpTwReuse: utils.Bool(raw["net_ipv4_tcp_tw_reuse"].(bool)),
	}
	if v := raw["net_core_somaxconn"].(int); v != 0 {
		result.NetCoreSomaxconn = utils.Int64(int64(v))
	}
	if v := raw["net_core_netdev_max_backlog"].(int); v != 0 {
		result.NetCoreNetdevMaxBacklog = utils.Int64(int64(v))
	}
	if v := raw["net_core_rmem_default"].(int); v != 0 {
		result.NetCoreRmemDefault = utils.Int64(int64(v))
	}
	if v := raw["net_core_rmem_max"].(int); v != 0 {
		result.NetCoreRmemMax = utils.Int64(int64(v))
	}
	if v := raw["net_core_wmem_default"].(int); v != 0 {
		result.NetCoreWmemDefault = utils.Int64(int64(v))
	}
	if v := raw["net_core_wmem_max"].(int); v != 0 {
		result.NetCoreWmemMax = utils.Int64(int64(v))
	}
	if v := raw["net_core_optmem_max"].(int); v != 0 {
		result.NetCoreOptmemMax = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_max_syn_backlog"].(int); v != 0 {
		result.NetIPv4TcpMaxSynBacklog = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_max_tw_buckets"].(int); v != 0 {
		result.NetIPv4TcpMaxTwBuckets = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_fin_timeout"].(int); v != 0 {
		result.NetIPv4TcpFinTimeout = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_keepalive_time"].(int); v != 0 {
		result.NetIPv4TcpKeepaliveTime = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_keepalive_probes"].(int); v != 0 {
		result.NetIPv4TcpKeepaliveProbes = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_tcp_keepalive_intvl"].(int); v != 0 {
		result.NetIPv4TcpkeepaliveIntvl = utils.Int64(int64(v))
	}
	netIPv4IPLocalPortRangeMin := raw["net_ipv4_ip_local_port_range_min"].(int)
	netIPv4IPLocalPortRangeMax := raw["net_ipv4_ip_local_port_range_max"].(int)
	if (netIPv4IPLocalPortRangeMin != 0 && netIPv4IPLocalPortRangeMax == 0) || (netIPv4IPLocalPortRangeMin == 0 && netIPv4IPLocalPortRangeMax != 0) {
		return nil, fmt.Errorf("`net_ipv4_ip_local_port_range_min` and `net_ipv4_ip_local_port_range_max` should both be set or unset")
	}
	if netIPv4IPLocalPortRangeMin > netIPv4IPLocalPortRangeMax {
		return nil, fmt.Errorf("`net_ipv4_ip_local_port_range_min` should be no larger than `net_ipv4_ip_local_port_range_max`")
	}
	if netIPv4IPLocalPortRangeMin != 0 && netIPv4IPLocalPortRangeMax != 0 {
		result.NetIPv4IPLocalPortRange = utils.String(fmt.Sprintf("%d %d", netIPv4IPLocalPortRangeMin, netIPv4IPLocalPortRangeMax))
	}
	if v := raw["net_ipv4_neigh_default_gc_thresh1"].(int); v != 0 {
		result.NetIPv4NeighDefaultGcThresh1 = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_neigh_default_gc_thresh2"].(int); v != 0 {
		result.NetIPv4NeighDefaultGcThresh2 = utils.Int64(int64(v))
	}
	if v := raw["net_ipv4_neigh_default_gc_thresh3"].(int); v != 0 {
		result.NetIPv4NeighDefaultGcThresh3 = utils.Int64(int64(v))
	}
	if v := raw["net_netfilter_nf_conntrack_max"].(int); v != 0 {
		result.NetNetfilterNfConntrackMax = utils.Int64(int64(v))
	}
	if v := raw["net_netfilter_nf_conntrack_buckets"].(int); v != 0 {
		result.NetNetfilterNfConntrackBuckets = utils.Int64(int64(v))
	}
	if v := raw["fs_aio_max_nr"].(int); v != 0 {
		result.FsAioMaxNr = utils.Int64(int64(v))
	}
	if v := raw["fs_inotify_max_user_watches"].(int); v != 0 {
		result.FsInotifyMaxUserWatches = utils.Int64(int64(v))
	}
	if v := raw["fs_file_max"].(int); v != 0 {
		result.FsFileMax = utils.Int64(int64(v))
	}
	if v := raw["fs_nr_open"].(int); v != 0 {
		result.FsNrOpen = utils.Int64(int64(v))
	}
	if v := raw["kernel_threads_max"].(int); v != 0 {
		result.KernelThreadsMax = utils.Int64(int64(v))
	}
	if v := raw["vm_max_map_count"].(int); v != 0 {
		result.VmMaxMapCount = utils.Int64(int64(v))
	}
	if v := raw["vm_swappiness"].(int); v != 0 {
		result.VmSwappiness = utils.Int64(int64(v))
	}
	if v := raw["vm_vfs_cache_pressure"].(int); v != 0 {
		result.VmVfsCachePressure = utils.Int64(int64(v))
	}
	return result, nil
}

func FlattenDefaultNodePool(input *[]managedclusters.ManagedClusterAgentPoolProfile, d *pluginsdk.ResourceData) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	agentPool, err := findDefaultNodePool(input, d)
	if err != nil {
		return nil, err
	}

	count := 0
	if agentPool.Count != nil {
		count = int(*agentPool.Count)
	}

	enableUltraSSD := false
	if agentPool.EnableUltraSSD != nil {
		enableUltraSSD = *agentPool.EnableUltraSSD
	}

	enableAutoScaling := false
	if agentPool.EnableAutoScaling != nil {
		enableAutoScaling = *agentPool.EnableAutoScaling
	}

	enableFIPS := false
	if agentPool.EnableFIPS != nil {
		enableFIPS = *agentPool.EnableFIPS
	}

	enableNodePublicIP := false
	if agentPool.EnableNodePublicIP != nil {
		enableNodePublicIP = *agentPool.EnableNodePublicIP
	}

	enableHostEncryption := false
	if agentPool.EnableEncryptionAtHost != nil {
		enableHostEncryption = *agentPool.EnableEncryptionAtHost
	}

	maxCount := 0
	if agentPool.MaxCount != nil {
		maxCount = int(*agentPool.MaxCount)
	}

	maxPods := 0
	if agentPool.MaxPods != nil {
		maxPods = int(*agentPool.MaxPods)
	}

	messageOfTheDay := ""
	if agentPool.MessageOfTheDay != nil {
		messageOfTheDayDecoded, err := base64.StdEncoding.DecodeString(*agentPool.MessageOfTheDay)
		if err != nil {
			return nil, err
		}
		messageOfTheDay = string(messageOfTheDayDecoded)
	}

	minCount := 0
	if agentPool.MinCount != nil {
		minCount = int(*agentPool.MinCount)
	}

	name := ""
	if agentPool.Name != "" {
		name = agentPool.Name
	}

	var nodeLabels map[string]string
	if agentPool.NodeLabels != nil {
		nodeLabels = make(map[string]string)
		for k, v := range *agentPool.NodeLabels {
			nodeLabels[k] = v
		}
	}

	nodePublicIPPrefixID := ""
	if agentPool.NodePublicIPPrefixID != nil {
		nodePublicIPPrefixID = *agentPool.NodePublicIPPrefixID
	}

	criticalAddonsEnabled := false
	if agentPool.NodeTaints != nil {
		for _, taint := range *agentPool.NodeTaints {
			if strings.EqualFold(taint, "CriticalAddonsOnly=true:NoSchedule") {
				criticalAddonsEnabled = true
			}
		}
	}

	osDiskSizeGB := 0
	if agentPool.OsDiskSizeGB != nil {
		osDiskSizeGB = int(*agentPool.OsDiskSizeGB)
	}

	osDiskType := managedclusters.OSDiskTypeManaged
	if *agentPool.OsDiskType != "" {
		osDiskType = *agentPool.OsDiskType
	}

	podSubnetId := ""
	if agentPool.PodSubnetID != nil {
		podSubnetId = *agentPool.PodSubnetID
	}

	vnetSubnetId := ""
	if agentPool.VnetSubnetID != nil {
		vnetSubnetId = *agentPool.VnetSubnetID
	}

	hostGroupID := ""
	if agentPool.HostGroupID != nil {
		hostGroupID = *agentPool.HostGroupID
	}

	orchestratorVersion := ""
	// NOTE: workaround for migration from 2022-01-02-preview (<3.12.0) to 2022-03-02-preview (>=3.12.0). Before terraform apply is run against the new API, Azure will respond only with currentOrchestratorVersion, orchestratorVersion will be absent. More details: https://github.com/hashicorp/terraform-provider-azurerm/issues/17833#issuecomment-1227583353
	if agentPool.OrchestratorVersion != nil {
		orchestratorVersion = *agentPool.OrchestratorVersion
	} else if agentPool.CurrentOrchestratorVersion != nil {
		orchestratorVersion = *agentPool.CurrentOrchestratorVersion
	}

	proximityPlacementGroupId := ""
	if agentPool.ProximityPlacementGroupID != nil {
		proximityPlacementGroupId = *agentPool.ProximityPlacementGroupID
	}

	scaleDownMode := managedclusters.ScaleDownModeDelete
	if *agentPool.ScaleDownMode != "" {
		scaleDownMode = *agentPool.ScaleDownMode
	}

	vmSize := ""
	if agentPool.VmSize != nil {
		vmSize = *agentPool.VmSize
	}
	capacityReservationGroupId := ""
	if agentPool.CapacityReservationGroupID != nil {
		capacityReservationGroupId = *agentPool.CapacityReservationGroupID
	}

	workloadRunTime := ""
	if *agentPool.WorkloadRuntime != "" {
		workloadRunTime = string(*agentPool.WorkloadRuntime)
	}

	upgradeSettings := flattenUpgradeSettingsFromCluster(agentPool.UpgradeSettings)
	linuxOSConfig, err := flattenAgentPoolLinuxOSConfigForCluster(agentPool.LinuxOSConfig)
	if err != nil {
		return nil, err
	}

	out := map[string]interface{}{
		"enable_auto_scaling":           enableAutoScaling,
		"enable_node_public_ip":         enableNodePublicIP,
		"enable_host_encryption":        enableHostEncryption,
		"fips_enabled":                  enableFIPS,
		"host_group_id":                 hostGroupID,
		"kubelet_disk_type":             string(*agentPool.KubeletDiskType),
		"max_count":                     maxCount,
		"max_pods":                      maxPods,
		"message_of_the_day":            messageOfTheDay,
		"min_count":                     minCount,
		"name":                          name,
		"node_count":                    count,
		"node_labels":                   nodeLabels,
		"node_public_ip_prefix_id":      nodePublicIPPrefixID,
		"node_taints":                   []string{},
		"os_disk_size_gb":               osDiskSizeGB,
		"os_disk_type":                  string(osDiskType),
		"os_sku":                        string(*agentPool.OsSKU),
		"scale_down_mode":               string(scaleDownMode),
		"tags":                          agentPool.Tags,
		"type":                          agentPool.Type,
		"ultra_ssd_enabled":             enableUltraSSD,
		"vm_size":                       vmSize,
		"workload_runtime":              workloadRunTime,
		"pod_subnet_id":                 podSubnetId,
		"orchestrator_version":          orchestratorVersion,
		"proximity_placement_group_id":  proximityPlacementGroupId,
		"upgrade_settings":              upgradeSettings,
		"vnet_subnet_id":                vnetSubnetId,
		"only_critical_addons_enabled":  criticalAddonsEnabled,
		"kubelet_config":                flattenAgentPoolKubeletConfigForCluster(agentPool.KubeletConfig),
		"linux_os_config":               linuxOSConfig,
		"zones":                         zones.Flatten(agentPool.AvailabilityZones),
		"capacity_reservation_group_id": capacityReservationGroupId,
	}

	return &[]interface{}{
		out,
	}, nil
}

func flattenAgentPoolKubeletConfigForCluster(input *managedclusters.KubeletConfig) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var cpuManagerPolicy, cpuCfsQuotaPeriod, topologyManagerPolicy string
	var cpuCfsQuotaEnabled bool
	var imageGcHighThreshold, imageGcLowThreshold, containerLogMaxSizeMB, containerLogMaxLines, podMaxPids int

	if input.CpuManagerPolicy != nil {
		cpuManagerPolicy = *input.CpuManagerPolicy
	}
	if input.CpuCfsQuota != nil {
		cpuCfsQuotaEnabled = *input.CpuCfsQuota
	}
	if input.CpuCfsQuotaPeriod != nil {
		cpuCfsQuotaPeriod = *input.CpuCfsQuotaPeriod
	}
	if input.ImageGcHighThreshold != nil {
		imageGcHighThreshold = int(*input.ImageGcHighThreshold)
	}
	if input.ImageGcLowThreshold != nil {
		imageGcLowThreshold = int(*input.ImageGcLowThreshold)
	}
	if input.TopologyManagerPolicy != nil {
		topologyManagerPolicy = *input.TopologyManagerPolicy
	}
	if input.ContainerLogMaxSizeMB != nil {
		containerLogMaxSizeMB = int(*input.ContainerLogMaxSizeMB)
	}
	if input.ContainerLogMaxFiles != nil {
		containerLogMaxLines = int(*input.ContainerLogMaxFiles)
	}
	if input.PodMaxPids != nil {
		podMaxPids = int(*input.PodMaxPids)
	}

	return []interface{}{
		map[string]interface{}{
			"cpu_manager_policy":        cpuManagerPolicy,
			"cpu_cfs_quota_enabled":     cpuCfsQuotaEnabled,
			"cpu_cfs_quota_period":      cpuCfsQuotaPeriod,
			"image_gc_high_threshold":   imageGcHighThreshold,
			"image_gc_low_threshold":    imageGcLowThreshold,
			"topology_manager_policy":   topologyManagerPolicy,
			"allowed_unsafe_sysctls":    utils.FlattenStringSlice(input.AllowedUnsafeSysctls),
			"container_log_max_size_mb": containerLogMaxSizeMB,
			"container_log_max_line":    containerLogMaxLines,
			"pod_max_pid":               podMaxPids,
		},
	}
}

func flattenAgentPoolKubeletConfig(input *agentpools.KubeletConfig) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var cpuManagerPolicy, cpuCfsQuotaPeriod, topologyManagerPolicy string
	var cpuCfsQuotaEnabled bool
	var imageGcHighThreshold, imageGcLowThreshold, containerLogMaxSizeMB, containerLogMaxLines, podMaxPids int

	if input.CpuManagerPolicy != nil {
		cpuManagerPolicy = *input.CpuManagerPolicy
	}
	if input.CpuCfsQuota != nil {
		cpuCfsQuotaEnabled = *input.CpuCfsQuota
	}
	if input.CpuCfsQuotaPeriod != nil {
		cpuCfsQuotaPeriod = *input.CpuCfsQuotaPeriod
	}
	if input.ImageGcHighThreshold != nil {
		imageGcHighThreshold = int(*input.ImageGcHighThreshold)
	}
	if input.ImageGcLowThreshold != nil {
		imageGcLowThreshold = int(*input.ImageGcLowThreshold)
	}
	if input.TopologyManagerPolicy != nil {
		topologyManagerPolicy = *input.TopologyManagerPolicy
	}
	if input.ContainerLogMaxSizeMB != nil {
		containerLogMaxSizeMB = int(*input.ContainerLogMaxSizeMB)
	}
	if input.ContainerLogMaxFiles != nil {
		containerLogMaxLines = int(*input.ContainerLogMaxFiles)
	}
	if input.PodMaxPids != nil {
		podMaxPids = int(*input.PodMaxPids)
	}

	return []interface{}{
		map[string]interface{}{
			"cpu_manager_policy":        cpuManagerPolicy,
			"cpu_cfs_quota_enabled":     cpuCfsQuotaEnabled,
			"cpu_cfs_quota_period":      cpuCfsQuotaPeriod,
			"image_gc_high_threshold":   imageGcHighThreshold,
			"image_gc_low_threshold":    imageGcLowThreshold,
			"topology_manager_policy":   topologyManagerPolicy,
			"allowed_unsafe_sysctls":    utils.FlattenStringSlice(input.AllowedUnsafeSysctls),
			"container_log_max_size_mb": containerLogMaxSizeMB,
			"container_log_max_line":    containerLogMaxLines,
			"pod_max_pid":               podMaxPids,
		},
	}
}

func flattenAgentPoolLinuxOSConfig(input *agentpools.LinuxOSConfig) ([]interface{}, error) {
	if input == nil {
		return make([]interface{}, 0), nil
	}

	var swapFileSizeMB int
	if input.SwapFileSizeMB != nil {
		swapFileSizeMB = int(*input.SwapFileSizeMB)
	}
	var transparentHugePageDefrag string
	if input.TransparentHugePageDefrag != nil {
		transparentHugePageDefrag = *input.TransparentHugePageDefrag
	}
	var transparentHugePageEnabled string
	if input.TransparentHugePageEnabled != nil {
		transparentHugePageEnabled = *input.TransparentHugePageEnabled
	}
	sysctlConfig, err := flattenAgentPoolSysctlConfig(input.Sysctls)
	if err != nil {
		return nil, err
	}
	return []interface{}{
		map[string]interface{}{
			"swap_file_size_mb":             swapFileSizeMB,
			"sysctl_config":                 sysctlConfig,
			"transparent_huge_page_defrag":  transparentHugePageDefrag,
			"transparent_huge_page_enabled": transparentHugePageEnabled,
		},
	}, nil
}

func flattenAgentPoolLinuxOSConfigForCluster(input *managedclusters.LinuxOSConfig) ([]interface{}, error) {
	if input == nil {
		return make([]interface{}, 0), nil
	}

	var swapFileSizeMB int
	if input.SwapFileSizeMB != nil {
		swapFileSizeMB = int(*input.SwapFileSizeMB)
	}
	var transparentHugePageDefrag string
	if input.TransparentHugePageDefrag != nil {
		transparentHugePageDefrag = *input.TransparentHugePageDefrag
	}
	var transparentHugePageEnabled string
	if input.TransparentHugePageEnabled != nil {
		transparentHugePageEnabled = *input.TransparentHugePageEnabled
	}
	sysctlConfig, err := flattenAgentPoolSysctlConfigForCluster(input.Sysctls)
	if err != nil {
		return nil, err
	}
	return []interface{}{
		map[string]interface{}{
			"swap_file_size_mb":             swapFileSizeMB,
			"sysctl_config":                 sysctlConfig,
			"transparent_huge_page_defrag":  transparentHugePageDefrag,
			"transparent_huge_page_enabled": transparentHugePageEnabled,
		},
	}, nil
}

func flattenAgentPoolSysctlConfig(input *agentpools.SysctlConfig) ([]interface{}, error) {
	if input == nil {
		return make([]interface{}, 0), nil
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
	var netIPv4IpLocalPortRangeMin, netIPv4IpLocalPortRangeMax int
	if input.NetIPv4IPLocalPortRange != nil {
		arr := regexp.MustCompile("[ \t]+").Split(*input.NetIPv4IPLocalPortRange, -1)
		if len(arr) != 2 {
			return nil, fmt.Errorf("parsing `NetIPv4IPLocalPortRange` %s", *input.NetIPv4IPLocalPortRange)
		}
		var err error
		netIPv4IpLocalPortRangeMin, err = strconv.Atoi(arr[0])
		if err != nil {
			return nil, err
		}
		netIPv4IpLocalPortRangeMax, err = strconv.Atoi(arr[1])
		if err != nil {
			return nil, err
		}
	}
	var netIPv4NeighDefaultGcThresh1 int
	if input.NetIPv4NeighDefaultGcThresh1 != nil {
		netIPv4NeighDefaultGcThresh1 = int(*input.NetIPv4NeighDefaultGcThresh1)
	}
	var netIPv4NeighDefaultGcThresh2 int
	if input.NetIPv4NeighDefaultGcThresh2 != nil {
		netIPv4NeighDefaultGcThresh2 = int(*input.NetIPv4NeighDefaultGcThresh2)
	}
	var netIPv4NeighDefaultGcThresh3 int
	if input.NetIPv4NeighDefaultGcThresh3 != nil {
		netIPv4NeighDefaultGcThresh3 = int(*input.NetIPv4NeighDefaultGcThresh3)
	}
	var netIPv4TcpFinTimeout int
	if input.NetIPv4TcpFinTimeout != nil {
		netIPv4TcpFinTimeout = int(*input.NetIPv4TcpFinTimeout)
	}
	var netIPv4TcpkeepaliveIntvl int
	if input.NetIPv4TcpkeepaliveIntvl != nil {
		netIPv4TcpkeepaliveIntvl = int(*input.NetIPv4TcpkeepaliveIntvl)
	}
	var netIPv4TcpKeepaliveProbes int
	if input.NetIPv4TcpKeepaliveProbes != nil {
		netIPv4TcpKeepaliveProbes = int(*input.NetIPv4TcpKeepaliveProbes)
	}
	var netIPv4TcpKeepaliveTime int
	if input.NetIPv4TcpKeepaliveTime != nil {
		netIPv4TcpKeepaliveTime = int(*input.NetIPv4TcpKeepaliveTime)
	}
	var netIPv4TcpMaxSynBacklog int
	if input.NetIPv4TcpMaxSynBacklog != nil {
		netIPv4TcpMaxSynBacklog = int(*input.NetIPv4TcpMaxSynBacklog)
	}
	var netIPv4TcpMaxTwBuckets int
	if input.NetIPv4TcpMaxTwBuckets != nil {
		netIPv4TcpMaxTwBuckets = int(*input.NetIPv4TcpMaxTwBuckets)
	}
	var netIPv4TcpTwReuse bool
	if input.NetIPv4TcpTwReuse != nil {
		netIPv4TcpTwReuse = *input.NetIPv4TcpTwReuse
	}
	var netNetfilterNfConntrackBuckets int
	if input.NetNetfilterNfConntrackBuckets != nil {
		netNetfilterNfConntrackBuckets = int(*input.NetNetfilterNfConntrackBuckets)
	}
	var netNetfilterNfConntrackMax int
	if input.NetNetfilterNfConntrackMax != nil {
		netNetfilterNfConntrackMax = int(*input.NetNetfilterNfConntrackMax)
	}
	var vmMaxMapCount int
	if input.VmMaxMapCount != nil {
		vmMaxMapCount = int(*input.VmMaxMapCount)
	}
	var vmSwappiness int
	if input.VmSwappiness != nil {
		vmSwappiness = int(*input.VmSwappiness)
	}
	var vmVfsCachePressure int
	if input.VmVfsCachePressure != nil {
		vmVfsCachePressure = int(*input.VmVfsCachePressure)
	}
	return []interface{}{
		map[string]interface{}{
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
			"net_ipv4_ip_local_port_range_min":   netIPv4IpLocalPortRangeMin,
			"net_ipv4_ip_local_port_range_max":   netIPv4IpLocalPortRangeMax,
			"net_ipv4_neigh_default_gc_thresh1":  netIPv4NeighDefaultGcThresh1,
			"net_ipv4_neigh_default_gc_thresh2":  netIPv4NeighDefaultGcThresh2,
			"net_ipv4_neigh_default_gc_thresh3":  netIPv4NeighDefaultGcThresh3,
			"net_ipv4_tcp_fin_timeout":           netIPv4TcpFinTimeout,
			"net_ipv4_tcp_keepalive_intvl":       netIPv4TcpkeepaliveIntvl,
			"net_ipv4_tcp_keepalive_probes":      netIPv4TcpKeepaliveProbes,
			"net_ipv4_tcp_keepalive_time":        netIPv4TcpKeepaliveTime,
			"net_ipv4_tcp_max_syn_backlog":       netIPv4TcpMaxSynBacklog,
			"net_ipv4_tcp_max_tw_buckets":        netIPv4TcpMaxTwBuckets,
			"net_ipv4_tcp_tw_reuse":              netIPv4TcpTwReuse,
			"net_netfilter_nf_conntrack_buckets": netNetfilterNfConntrackBuckets,
			"net_netfilter_nf_conntrack_max":     netNetfilterNfConntrackMax,
			"vm_max_map_count":                   vmMaxMapCount,
			"vm_swappiness":                      vmSwappiness,
			"vm_vfs_cache_pressure":              vmVfsCachePressure,
		},
	}, nil
}

func flattenAgentPoolSysctlConfigForCluster(input *managedclusters.SysctlConfig) ([]interface{}, error) {
	if input == nil {
		return make([]interface{}, 0), nil
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
	var netIPv4IpLocalPortRangeMin, netIPv4IpLocalPortRangeMax int
	if input.NetIPv4IPLocalPortRange != nil {
		arr := regexp.MustCompile("[ \t]+").Split(*input.NetIPv4IPLocalPortRange, -1)
		if len(arr) != 2 {
			return nil, fmt.Errorf("parsing `NetIPv4IPLocalPortRange` %s", *input.NetIPv4IPLocalPortRange)
		}
		var err error
		netIPv4IpLocalPortRangeMin, err = strconv.Atoi(arr[0])
		if err != nil {
			return nil, err
		}
		netIPv4IpLocalPortRangeMax, err = strconv.Atoi(arr[1])
		if err != nil {
			return nil, err
		}
	}
	var netIPv4NeighDefaultGcThresh1 int
	if input.NetIPv4NeighDefaultGcThresh1 != nil {
		netIPv4NeighDefaultGcThresh1 = int(*input.NetIPv4NeighDefaultGcThresh1)
	}
	var netIPv4NeighDefaultGcThresh2 int
	if input.NetIPv4NeighDefaultGcThresh2 != nil {
		netIPv4NeighDefaultGcThresh2 = int(*input.NetIPv4NeighDefaultGcThresh2)
	}
	var netIPv4NeighDefaultGcThresh3 int
	if input.NetIPv4NeighDefaultGcThresh3 != nil {
		netIPv4NeighDefaultGcThresh3 = int(*input.NetIPv4NeighDefaultGcThresh3)
	}
	var netIPv4TcpFinTimeout int
	if input.NetIPv4TcpFinTimeout != nil {
		netIPv4TcpFinTimeout = int(*input.NetIPv4TcpFinTimeout)
	}
	var netIPv4TcpkeepaliveIntvl int
	if input.NetIPv4TcpkeepaliveIntvl != nil {
		netIPv4TcpkeepaliveIntvl = int(*input.NetIPv4TcpkeepaliveIntvl)
	}
	var netIPv4TcpKeepaliveProbes int
	if input.NetIPv4TcpKeepaliveProbes != nil {
		netIPv4TcpKeepaliveProbes = int(*input.NetIPv4TcpKeepaliveProbes)
	}
	var netIPv4TcpKeepaliveTime int
	if input.NetIPv4TcpKeepaliveTime != nil {
		netIPv4TcpKeepaliveTime = int(*input.NetIPv4TcpKeepaliveTime)
	}
	var netIPv4TcpMaxSynBacklog int
	if input.NetIPv4TcpMaxSynBacklog != nil {
		netIPv4TcpMaxSynBacklog = int(*input.NetIPv4TcpMaxSynBacklog)
	}
	var netIPv4TcpMaxTwBuckets int
	if input.NetIPv4TcpMaxTwBuckets != nil {
		netIPv4TcpMaxTwBuckets = int(*input.NetIPv4TcpMaxTwBuckets)
	}
	var netIPv4TcpTwReuse bool
	if input.NetIPv4TcpTwReuse != nil {
		netIPv4TcpTwReuse = *input.NetIPv4TcpTwReuse
	}
	var netNetfilterNfConntrackBuckets int
	if input.NetNetfilterNfConntrackBuckets != nil {
		netNetfilterNfConntrackBuckets = int(*input.NetNetfilterNfConntrackBuckets)
	}
	var netNetfilterNfConntrackMax int
	if input.NetNetfilterNfConntrackMax != nil {
		netNetfilterNfConntrackMax = int(*input.NetNetfilterNfConntrackMax)
	}
	var vmMaxMapCount int
	if input.VmMaxMapCount != nil {
		vmMaxMapCount = int(*input.VmMaxMapCount)
	}
	var vmSwappiness int
	if input.VmSwappiness != nil {
		vmSwappiness = int(*input.VmSwappiness)
	}
	var vmVfsCachePressure int
	if input.VmVfsCachePressure != nil {
		vmVfsCachePressure = int(*input.VmVfsCachePressure)
	}
	return []interface{}{
		map[string]interface{}{
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
			"net_ipv4_ip_local_port_range_min":   netIPv4IpLocalPortRangeMin,
			"net_ipv4_ip_local_port_range_max":   netIPv4IpLocalPortRangeMax,
			"net_ipv4_neigh_default_gc_thresh1":  netIPv4NeighDefaultGcThresh1,
			"net_ipv4_neigh_default_gc_thresh2":  netIPv4NeighDefaultGcThresh2,
			"net_ipv4_neigh_default_gc_thresh3":  netIPv4NeighDefaultGcThresh3,
			"net_ipv4_tcp_fin_timeout":           netIPv4TcpFinTimeout,
			"net_ipv4_tcp_keepalive_intvl":       netIPv4TcpkeepaliveIntvl,
			"net_ipv4_tcp_keepalive_probes":      netIPv4TcpKeepaliveProbes,
			"net_ipv4_tcp_keepalive_time":        netIPv4TcpKeepaliveTime,
			"net_ipv4_tcp_max_syn_backlog":       netIPv4TcpMaxSynBacklog,
			"net_ipv4_tcp_max_tw_buckets":        netIPv4TcpMaxTwBuckets,
			"net_ipv4_tcp_tw_reuse":              netIPv4TcpTwReuse,
			"net_netfilter_nf_conntrack_buckets": netNetfilterNfConntrackBuckets,
			"net_netfilter_nf_conntrack_max":     netNetfilterNfConntrackMax,
			"vm_max_map_count":                   vmMaxMapCount,
			"vm_swappiness":                      vmSwappiness,
			"vm_vfs_cache_pressure":              vmVfsCachePressure,
		},
	}, nil
}

func findDefaultNodePool(input *[]managedclusters.ManagedClusterAgentPoolProfile, d *pluginsdk.ResourceData) (*managedclusters.ManagedClusterAgentPoolProfile, error) {
	// first try loading this from the Resource Data if possible (e.g. when Created)
	defaultNodePoolName := d.Get("default_node_pool.0.name")

	var agentPool *managedclusters.ManagedClusterAgentPoolProfile
	if defaultNodePoolName != "" {
		// find it
		for _, v := range *input {
			if v.Name != "" && v.Name == defaultNodePoolName {
				agentPool = &v
				break
			}
		}
	}

	if agentPool == nil {
		// otherwise we need to fall back to the name of the first agent pool
		for _, v := range *input {
			if v.Name == "" {
				continue
			}
			if *v.Mode != managedclusters.AgentPoolModeSystem {
				continue
			}

			defaultNodePoolName = v.Name
			agentPool = &v
			break
		}

		if defaultNodePoolName == nil {
			return nil, fmt.Errorf("Unable to Determine Default Agent Pool")
		}
	}

	if agentPool == nil {
		return nil, fmt.Errorf("The Default Agent Pool %q was not found", defaultNodePoolName)
	}

	return agentPool, nil
}
