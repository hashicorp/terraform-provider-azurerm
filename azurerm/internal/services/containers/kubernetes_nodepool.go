package containers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-03-01/containerservice"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	computeValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SchemaDefaultNodePool() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
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

				// Optional
				"availability_zones": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"enable_auto_scaling": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"enable_node_public_ip": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
				},

				"enable_host_encryption": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
				},

				"kubelet_config": schemaNodePoolKubeletConfig(),

				"linux_os_config": schemaNodePoolLinuxOSConfig(),

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
					ForceNew: true,
					Optional: true,
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

				"tags": tags.Schema(),

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
				"proximity_placement_group_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: computeValidate.ProximityPlacementGroupID,
				},
				"only_critical_addons_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
				},

				"upgrade_settings": upgradeSettingsSchema(),
			},
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
					ValidateFunc: validation.IntBetween(1024, 60999),
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
					ValidateFunc: validation.IntBetween(131072, 589824),
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

func ConvertDefaultNodePoolToAgentPool(input *[]containerservice.ManagedClusterAgentPoolProfile) containerservice.AgentPool {
	defaultCluster := (*input)[0]
	return containerservice.AgentPool{
		Name: defaultCluster.Name,
		ManagedClusterAgentPoolProfileProperties: &containerservice.ManagedClusterAgentPoolProfileProperties{
			Count:                     defaultCluster.Count,
			VMSize:                    defaultCluster.VMSize,
			OsDiskSizeGB:              defaultCluster.OsDiskSizeGB,
			OsDiskType:                defaultCluster.OsDiskType,
			VnetSubnetID:              defaultCluster.VnetSubnetID,
			KubeletConfig:             defaultCluster.KubeletConfig,
			LinuxOSConfig:             defaultCluster.LinuxOSConfig,
			MaxPods:                   defaultCluster.MaxPods,
			OsType:                    defaultCluster.OsType,
			MaxCount:                  defaultCluster.MaxCount,
			MinCount:                  defaultCluster.MinCount,
			EnableAutoScaling:         defaultCluster.EnableAutoScaling,
			Type:                      defaultCluster.Type,
			OrchestratorVersion:       defaultCluster.OrchestratorVersion,
			ProximityPlacementGroupID: defaultCluster.ProximityPlacementGroupID,
			AvailabilityZones:         defaultCluster.AvailabilityZones,
			EnableNodePublicIP:        defaultCluster.EnableNodePublicIP,
			NodePublicIPPrefixID:      defaultCluster.NodePublicIPPrefixID,
			ScaleSetPriority:          defaultCluster.ScaleSetPriority,
			ScaleSetEvictionPolicy:    defaultCluster.ScaleSetEvictionPolicy,
			SpotMaxPrice:              defaultCluster.SpotMaxPrice,
			Mode:                      defaultCluster.Mode,
			NodeLabels:                defaultCluster.NodeLabels,
			NodeTaints:                defaultCluster.NodeTaints,
			Tags:                      defaultCluster.Tags,
			UpgradeSettings:           defaultCluster.UpgradeSettings,
		},
	}
}

func ExpandDefaultNodePool(d *pluginsdk.ResourceData) (*[]containerservice.ManagedClusterAgentPoolProfile, error) {
	input := d.Get("default_node_pool").([]interface{})

	raw := input[0].(map[string]interface{})
	enableAutoScaling := raw["enable_auto_scaling"].(bool)
	nodeLabelsRaw := raw["node_labels"].(map[string]interface{})
	nodeLabels := utils.ExpandMapStringPtrString(nodeLabelsRaw)
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

	profile := containerservice.ManagedClusterAgentPoolProfile{
		EnableAutoScaling:      utils.Bool(enableAutoScaling),
		EnableNodePublicIP:     utils.Bool(raw["enable_node_public_ip"].(bool)),
		EnableEncryptionAtHost: utils.Bool(raw["enable_host_encryption"].(bool)),
		Name:                   utils.String(raw["name"].(string)),
		NodeLabels:             nodeLabels,
		NodeTaints:             nodeTaints,
		Tags:                   tags.Expand(t),
		Type:                   containerservice.AgentPoolType(raw["type"].(string)),
		VMSize:                 utils.String(raw["vm_size"].(string)),

		// at this time the default node pool has to be Linux or the AKS cluster fails to provision with:
		// Pods not in Running status: coredns-7fc597cc45-v5z7x,coredns-autoscaler-7ccc76bfbd-djl7j,metrics-server-cbd95f966-5rl97,tunnelfront-7d9884977b-wpbvn
		// Windows agents can be configured via the separate node pool resource
		OsType: containerservice.OSTypeLinux,

		// without this set the API returns:
		// Code="MustDefineAtLeastOneSystemPool" Message="Must define at least one system pool."
		// since this is the "default" node pool we can assume this is a system node pool
		Mode: containerservice.AgentPoolModeSystem,

		UpgradeSettings: expandUpgradeSettings(raw["upgrade_settings"].([]interface{})),

		// // TODO: support these in time
		// ScaleSetEvictionPolicy: "",
		// ScaleSetPriority:       "",
	}

	availabilityZonesRaw := raw["availability_zones"].([]interface{})
	availabilityZones := utils.ExpandStringSlice(availabilityZonesRaw)

	// otherwise: Standard Load Balancer is required for availability zone.
	if len(*availabilityZones) > 0 {
		profile.AvailabilityZones = availabilityZones
	}

	if maxPods := int32(raw["max_pods"].(int)); maxPods > 0 {
		profile.MaxPods = utils.Int32(maxPods)
	}

	if prefixID := raw["node_public_ip_prefix_id"].(string); prefixID != "" {
		profile.NodePublicIPPrefixID = utils.String(prefixID)
	}

	if osDiskSizeGB := int32(raw["os_disk_size_gb"].(int)); osDiskSizeGB > 0 {
		profile.OsDiskSizeGB = utils.Int32(osDiskSizeGB)
	}

	profile.OsDiskType = containerservice.OSDiskTypeManaged
	if osDiskType := raw["os_disk_type"].(string); osDiskType != "" {
		profile.OsDiskType = containerservice.OSDiskType(raw["os_disk_type"].(string))
	}

	if vnetSubnetID := raw["vnet_subnet_id"].(string); vnetSubnetID != "" {
		profile.VnetSubnetID = utils.String(vnetSubnetID)
	}

	if orchestratorVersion := raw["orchestrator_version"].(string); orchestratorVersion != "" {
		profile.OrchestratorVersion = utils.String(orchestratorVersion)
	}

	if proximityPlacementGroupId := raw["proximity_placement_group_id"].(string); proximityPlacementGroupId != "" {
		profile.ProximityPlacementGroupID = utils.String(proximityPlacementGroupId)
	}

	count := raw["node_count"].(int)
	maxCount := raw["max_count"].(int)
	minCount := raw["min_count"].(int)

	// Count must always be set (see #6094), RP behaviour has changed
	// since the API version upgrade in v2.1.0 making Count required
	// for all create/update requests
	profile.Count = utils.Int32(int32(count))

	if enableAutoScaling {
		// if Count has not been set use min count
		if count == 0 {
			count = minCount
			profile.Count = utils.Int32(int32(count))
		}

		// Count must be set for the initial creation when using AutoScaling but cannot be updated
		if d.HasChange("default_node_pool.0.node_count") && !d.IsNewResource() {
			return nil, fmt.Errorf("cannot change `node_count` when `enable_auto_scaling` is set to `true`")
		}

		if maxCount > 0 {
			profile.MaxCount = utils.Int32(int32(maxCount))
			if maxCount < count {
				return nil, fmt.Errorf("`node_count`(%d) must be equal to or less than `max_count`(%d) when `enable_auto_scaling` is set to `true`", count, maxCount)
			}
		} else {
			return nil, fmt.Errorf("`max_count` must be configured when `enable_auto_scaling` is set to `true`")
		}

		if minCount > 0 {
			profile.MinCount = utils.Int32(int32(minCount))

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
		profile.KubeletConfig = expandAgentPoolKubeletConfig(kubeletConfig)
	}

	if linuxOSConfig := raw["linux_os_config"].([]interface{}); len(linuxOSConfig) > 0 {
		linuxOSConfig, err := expandAgentPoolLinuxOSConfig(linuxOSConfig)
		if err != nil {
			return nil, err
		}
		profile.LinuxOSConfig = linuxOSConfig
	}

	return &[]containerservice.ManagedClusterAgentPoolProfile{
		profile,
	}, nil
}

func expandAgentPoolKubeletConfig(input []interface{}) *containerservice.KubeletConfig {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	result := &containerservice.KubeletConfig{
		CPUCfsQuota: utils.Bool(raw["cpu_cfs_quota_enabled"].(bool)),
		// must be false, otherwise the backend will report error: CustomKubeletConfig.FailSwapOn must be set to false to enable swap file on nodes.
		FailSwapOn:           utils.Bool(false),
		AllowedUnsafeSysctls: utils.ExpandStringSlice(raw["allowed_unsafe_sysctls"].(*pluginsdk.Set).List()),
	}

	if v := raw["cpu_manager_policy"].(string); v != "" {
		result.CPUManagerPolicy = utils.String(v)
	}
	if v := raw["cpu_cfs_quota_period"].(string); v != "" {
		result.CPUCfsQuotaPeriod = utils.String(v)
	}
	if v := raw["image_gc_high_threshold"].(int); v != 0 {
		result.ImageGcHighThreshold = utils.Int32(int32(v))
	}
	if v := raw["image_gc_low_threshold"].(int); v != 0 {
		result.ImageGcLowThreshold = utils.Int32(int32(v))
	}
	if v := raw["topology_manager_policy"].(string); v != "" {
		result.TopologyManagerPolicy = utils.String(v)
	}
	if v := raw["container_log_max_size_mb"].(int); v != 0 {
		result.ContainerLogMaxSizeMB = utils.Int32(int32(v))
	}
	if v := raw["container_log_max_line"].(int); v != 0 {
		result.ContainerLogMaxFiles = utils.Int32(int32(v))
	}
	if v := raw["pod_max_pid"].(int); v != 0 {
		result.PodMaxPids = utils.Int32(int32(v))
	}

	return result
}

func expandAgentPoolLinuxOSConfig(input []interface{}) (*containerservice.LinuxOSConfig, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	raw := input[0].(map[string]interface{})
	sysctlConfig, err := expandAgentPoolSysctlConfig(raw["sysctl_config"].([]interface{}))
	if err != nil {
		return nil, err
	}

	result := &containerservice.LinuxOSConfig{
		Sysctls: sysctlConfig,
	}
	if v := raw["transparent_huge_page_enabled"].(string); v != "" {
		result.TransparentHugePageEnabled = utils.String(v)
	}
	if v := raw["transparent_huge_page_defrag"].(string); v != "" {
		result.TransparentHugePageDefrag = utils.String(v)
	}
	if v := raw["swap_file_size_mb"].(int); v != 0 {
		result.SwapFileSizeMB = utils.Int32(int32(v))
	}
	return result, nil
}

func expandAgentPoolSysctlConfig(input []interface{}) (*containerservice.SysctlConfig, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	raw := input[0].(map[string]interface{})
	result := &containerservice.SysctlConfig{
		NetIpv4TCPTwReuse: utils.Bool(raw["net_ipv4_tcp_tw_reuse"].(bool)),
	}
	if v := raw["net_core_somaxconn"].(int); v != 0 {
		result.NetCoreSomaxconn = utils.Int32(int32(v))
	}
	if v := raw["net_core_netdev_max_backlog"].(int); v != 0 {
		result.NetCoreNetdevMaxBacklog = utils.Int32(int32(v))
	}
	if v := raw["net_core_rmem_default"].(int); v != 0 {
		result.NetCoreRmemDefault = utils.Int32(int32(v))
	}
	if v := raw["net_core_rmem_max"].(int); v != 0 {
		result.NetCoreRmemMax = utils.Int32(int32(v))
	}
	if v := raw["net_core_wmem_default"].(int); v != 0 {
		result.NetCoreWmemDefault = utils.Int32(int32(v))
	}
	if v := raw["net_core_wmem_max"].(int); v != 0 {
		result.NetCoreWmemMax = utils.Int32(int32(v))
	}
	if v := raw["net_core_optmem_max"].(int); v != 0 {
		result.NetCoreOptmemMax = utils.Int32(int32(v))
	}
	if v := raw["net_ipv4_tcp_max_syn_backlog"].(int); v != 0 {
		result.NetIpv4TCPMaxSynBacklog = utils.Int32(int32(v))
	}
	if v := raw["net_ipv4_tcp_max_tw_buckets"].(int); v != 0 {
		result.NetIpv4TCPMaxTwBuckets = utils.Int32(int32(v))
	}
	if v := raw["net_ipv4_tcp_fin_timeout"].(int); v != 0 {
		result.NetIpv4TCPFinTimeout = utils.Int32(int32(v))
	}
	if v := raw["net_ipv4_tcp_keepalive_time"].(int); v != 0 {
		result.NetIpv4TCPKeepaliveTime = utils.Int32(int32(v))
	}
	if v := raw["net_ipv4_tcp_keepalive_probes"].(int); v != 0 {
		result.NetIpv4TCPKeepaliveProbes = utils.Int32(int32(v))
	}
	if v := raw["net_ipv4_tcp_keepalive_intvl"].(int); v != 0 {
		result.NetIpv4TcpkeepaliveIntvl = utils.Int32(int32(v))
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
		result.NetIpv4IPLocalPortRange = utils.String(fmt.Sprintf("%d %d", netIpv4IPLocalPortRangeMin, netIpv4IPLocalPortRangeMax))
	}
	if v := raw["net_ipv4_neigh_default_gc_thresh1"].(int); v != 0 {
		result.NetIpv4NeighDefaultGcThresh1 = utils.Int32(int32(v))
	}
	if v := raw["net_ipv4_neigh_default_gc_thresh2"].(int); v != 0 {
		result.NetIpv4NeighDefaultGcThresh2 = utils.Int32(int32(v))
	}
	if v := raw["net_ipv4_neigh_default_gc_thresh3"].(int); v != 0 {
		result.NetIpv4NeighDefaultGcThresh3 = utils.Int32(int32(v))
	}
	if v := raw["net_netfilter_nf_conntrack_max"].(int); v != 0 {
		result.NetNetfilterNfConntrackMax = utils.Int32(int32(v))
	}
	if v := raw["net_netfilter_nf_conntrack_buckets"].(int); v != 0 {
		result.NetNetfilterNfConntrackBuckets = utils.Int32(int32(v))
	}
	if v := raw["fs_aio_max_nr"].(int); v != 0 {
		result.FsAioMaxNr = utils.Int32(int32(v))
	}
	if v := raw["fs_inotify_max_user_watches"].(int); v != 0 {
		result.FsInotifyMaxUserWatches = utils.Int32(int32(v))
	}
	if v := raw["fs_file_max"].(int); v != 0 {
		result.FsFileMax = utils.Int32(int32(v))
	}
	if v := raw["fs_nr_open"].(int); v != 0 {
		result.FsNrOpen = utils.Int32(int32(v))
	}
	if v := raw["kernel_threads_max"].(int); v != 0 {
		result.KernelThreadsMax = utils.Int32(int32(v))
	}
	if v := raw["vm_max_map_count"].(int); v != 0 {
		result.VMMaxMapCount = utils.Int32(int32(v))
	}
	if v := raw["vm_swappiness"].(int); v != 0 {
		result.VMSwappiness = utils.Int32(int32(v))
	}
	if v := raw["vm_vfs_cache_pressure"].(int); v != 0 {
		result.VMVfsCachePressure = utils.Int32(int32(v))
	}
	return result, nil
}

func FlattenDefaultNodePool(input *[]containerservice.ManagedClusterAgentPoolProfile, d *pluginsdk.ResourceData) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	agentPool, err := findDefaultNodePool(input, d)
	if err != nil {
		return nil, err
	}

	var availabilityZones []string
	if agentPool.AvailabilityZones != nil {
		availabilityZones = *agentPool.AvailabilityZones
	}

	count := 0
	if agentPool.Count != nil {
		count = int(*agentPool.Count)
	}

	enableAutoScaling := false
	if agentPool.EnableAutoScaling != nil {
		enableAutoScaling = *agentPool.EnableAutoScaling
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

	minCount := 0
	if agentPool.MinCount != nil {
		minCount = int(*agentPool.MinCount)
	}

	name := ""
	if agentPool.Name != nil {
		name = *agentPool.Name
	}

	var nodeLabels map[string]string
	if agentPool.NodeLabels != nil {
		nodeLabels = make(map[string]string)
		for k, v := range agentPool.NodeLabels {
			nodeLabels[k] = *v
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

	osDiskType := containerservice.OSDiskTypeManaged
	if agentPool.OsDiskType != "" {
		osDiskType = agentPool.OsDiskType
	}

	vnetSubnetId := ""
	if agentPool.VnetSubnetID != nil {
		vnetSubnetId = *agentPool.VnetSubnetID
	}

	orchestratorVersion := ""
	if agentPool.OrchestratorVersion != nil {
		orchestratorVersion = *agentPool.OrchestratorVersion
	}

	proximityPlacementGroupId := ""
	if agentPool.ProximityPlacementGroupID != nil {
		proximityPlacementGroupId = *agentPool.ProximityPlacementGroupID
	}

	vmSize := ""
	if agentPool.VMSize != nil {
		vmSize = *agentPool.VMSize
	}

	upgradeSettings := flattenUpgradeSettings(agentPool.UpgradeSettings)
	linuxOSConfig, err := flattenAgentPoolLinuxOSConfig(agentPool.LinuxOSConfig)
	if err != nil {
		return nil, err
	}
	return &[]interface{}{
		map[string]interface{}{
			"availability_zones":           availabilityZones,
			"enable_auto_scaling":          enableAutoScaling,
			"enable_node_public_ip":        enableNodePublicIP,
			"enable_host_encryption":       enableHostEncryption,
			"max_count":                    maxCount,
			"max_pods":                     maxPods,
			"min_count":                    minCount,
			"name":                         name,
			"node_count":                   count,
			"node_labels":                  nodeLabels,
			"node_public_ip_prefix_id":     nodePublicIPPrefixID,
			"node_taints":                  []string{},
			"os_disk_size_gb":              osDiskSizeGB,
			"os_disk_type":                 string(osDiskType),
			"tags":                         tags.Flatten(agentPool.Tags),
			"type":                         string(agentPool.Type),
			"vm_size":                      vmSize,
			"orchestrator_version":         orchestratorVersion,
			"proximity_placement_group_id": proximityPlacementGroupId,
			"upgrade_settings":             upgradeSettings,
			"vnet_subnet_id":               vnetSubnetId,
			"only_critical_addons_enabled": criticalAddonsEnabled,
			"kubelet_config":               flattenAgentPoolKubeletConfig(agentPool.KubeletConfig),
			"linux_os_config":              linuxOSConfig,
		},
	}, nil
}

func flattenAgentPoolKubeletConfig(input *containerservice.KubeletConfig) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var cpuManagerPolicy, cpuCfsQuotaPeriod, topologyManagerPolicy string
	var cpuCfsQuotaEnabled bool
	var imageGcHighThreshold, imageGcLowThreshold, containerLogMaxSizeMB, containerLogMaxLines, podMaxPids int

	if input.CPUManagerPolicy != nil {
		cpuManagerPolicy = *input.CPUManagerPolicy
	}
	if input.CPUCfsQuota != nil {
		cpuCfsQuotaEnabled = *input.CPUCfsQuota
	}
	if input.CPUCfsQuotaPeriod != nil {
		cpuCfsQuotaPeriod = *input.CPUCfsQuotaPeriod
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

func flattenAgentPoolLinuxOSConfig(input *containerservice.LinuxOSConfig) ([]interface{}, error) {
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

func flattenAgentPoolSysctlConfig(input *containerservice.SysctlConfig) ([]interface{}, error) {
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
	var netIpv4IpLocalPortRangeMin, netIpv4IpLocalPortRangeMax int
	if input.NetIpv4IPLocalPortRange != nil {
		arr := regexp.MustCompile("[ \t]+").Split(*input.NetIpv4IPLocalPortRange, -1)
		if len(arr) != 2 {
			return nil, fmt.Errorf("parsing `NetIpv4IPLocalPortRange` %s", *input.NetIpv4IPLocalPortRange)
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
	if input.NetIpv4NeighDefaultGcThresh1 != nil {
		netIpv4NeighDefaultGcThresh1 = int(*input.NetIpv4NeighDefaultGcThresh1)
	}
	var netIpv4NeighDefaultGcThresh2 int
	if input.NetIpv4NeighDefaultGcThresh2 != nil {
		netIpv4NeighDefaultGcThresh2 = int(*input.NetIpv4NeighDefaultGcThresh2)
	}
	var netIpv4NeighDefaultGcThresh3 int
	if input.NetIpv4NeighDefaultGcThresh3 != nil {
		netIpv4NeighDefaultGcThresh3 = int(*input.NetIpv4NeighDefaultGcThresh3)
	}
	var netIpv4TcpFinTimeout int
	if input.NetIpv4TCPFinTimeout != nil {
		netIpv4TcpFinTimeout = int(*input.NetIpv4TCPFinTimeout)
	}
	var netIpv4TcpkeepaliveIntvl int
	if input.NetIpv4TcpkeepaliveIntvl != nil {
		netIpv4TcpkeepaliveIntvl = int(*input.NetIpv4TcpkeepaliveIntvl)
	}
	var netIpv4TcpKeepaliveProbes int
	if input.NetIpv4TCPKeepaliveProbes != nil {
		netIpv4TcpKeepaliveProbes = int(*input.NetIpv4TCPKeepaliveProbes)
	}
	var netIpv4TcpKeepaliveTime int
	if input.NetIpv4TCPKeepaliveTime != nil {
		netIpv4TcpKeepaliveTime = int(*input.NetIpv4TCPKeepaliveTime)
	}
	var netIpv4TcpMaxSynBacklog int
	if input.NetIpv4TCPMaxSynBacklog != nil {
		netIpv4TcpMaxSynBacklog = int(*input.NetIpv4TCPMaxSynBacklog)
	}
	var netIpv4TcpMaxTwBuckets int
	if input.NetIpv4TCPMaxTwBuckets != nil {
		netIpv4TcpMaxTwBuckets = int(*input.NetIpv4TCPMaxTwBuckets)
	}
	var netIpv4TcpTwReuse bool
	if input.NetIpv4TCPTwReuse != nil {
		netIpv4TcpTwReuse = *input.NetIpv4TCPTwReuse
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

func findDefaultNodePool(input *[]containerservice.ManagedClusterAgentPoolProfile, d *pluginsdk.ResourceData) (*containerservice.ManagedClusterAgentPoolProfile, error) {
	// first try loading this from the Resource Data if possible (e.g. when Created)
	defaultNodePoolName := d.Get("default_node_pool.0.name")

	var agentPool *containerservice.ManagedClusterAgentPoolProfile
	if defaultNodePoolName != "" {
		// find it
		for _, v := range *input {
			if v.Name != nil && *v.Name == defaultNodePoolName {
				agentPool = &v
				break
			}
		}
	}

	if agentPool == nil {
		// otherwise we need to fall back to the name of the first agent pool
		for _, v := range *input {
			if v.Name == nil {
				continue
			}
			if v.Mode != containerservice.AgentPoolModeSystem {
				continue
			}

			defaultNodePoolName = *v.Name
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
