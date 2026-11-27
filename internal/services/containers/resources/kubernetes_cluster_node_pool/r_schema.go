// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_node_pool

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/snapshots"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/helpers"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceKubernetesClusterNodePoolSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Node Pool which should be created within the Kubernetes Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: containerValidate.KubernetesAgentPoolName,
		},

		"kubernetes_cluster_id": {
			Description:  "The ID of the Kubernetes Cluster where this Node Pool should exist.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseKubernetesClusterIDInsensitively),
		},

		"node_count": {
			Description:  "The initial number of nodes which should exist within this Node Pool.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(0, 1000),
		},

		"tags": commonschema.Tags(),

		"vm_size": {
			Description: "The SKU which should be used for the Virtual Machines used in this Node Pool.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C AKS RP provides a new feature that will automatically select an available vm size when it's omitted.
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"host_group_id": {
			Description:  "The fully qualified resource ID of the Dedicated Host Group to provision Virtual Machines from.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseDedicatedHostGroupIDInsensitively),
		},

		// Optional
		"capacity_reservation_group_id": {
			Description:  "The ID of the Capacity Reservation Group where this Node Pool should exist.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: capacityreservationgroups.ValidateCapacityReservationGroupID,
		},

		"eviction_policy": {
			Description:  "The Eviction Policy which should be used for Virtual Machines within the Virtual Machine Scale Set powering this Node Pool.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForScaleSetEvictionPolicy(), false),
		},

		"kubelet_config": helpers.SchemaNodePoolKubeletConfig(),

		"linux_os_config": helpers.SchemaNodePoolLinuxOSConfig(),

		"fips_enabled": {
			Description: "Whether the nodes in this Node Pool should have Federal Information Processing Standard enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"gpu_instance": {
			Description: "The GPU MIG instance profile for supported GPU VM SKU.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(agentpools.GPUInstanceProfileMIGOneg),
				string(managedclusters.GPUInstanceProfileMIGTwog),
				string(managedclusters.GPUInstanceProfileMIGThreeg),
				string(managedclusters.GPUInstanceProfileMIGFourg),
				string(managedclusters.GPUInstanceProfileMIGSeveng),
			}, false),
		},

		"gpu_driver": {
			Description:  "Whether to install the GPU Driver for the nodes.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForGPUDriver(), false),
		},

		"kubelet_disk_type": {
			Description:  "The type of disk used by kubelet.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForKubeletDiskType(), false),
		},

		"max_count": {
			Description:  "The maximum number of nodes which should exist within this Node Pool.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 1000),
		},

		"max_pods": {
			Description: "The maximum number of pods that can run on each agent.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"mode": {
			Description: "Whether this Node Pool should be used for System or User resources.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(agentpools.AgentPoolModeUser),
			ValidateFunc: validation.StringInSlice([]string{
				string(agentpools.AgentPoolModeSystem),
				string(agentpools.AgentPoolModeUser),
			}, false),
		},

		"min_count": {
			Description: "The minimum number of nodes which should exist within this Node Pool.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// NOTE: rather than setting `0` users should instead pass `null` here
			ValidateFunc: validation.IntBetween(0, 1000),
		},

		"node_network_profile": helpers.SchemaNodePoolNetworkProfile(),

		"node_labels": {
			Description: "A map of Kubernetes labels which should be applied to nodes in this Node Pool.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"node_public_ip_prefix_id": {
			Description:  "Resource ID for the Public IP Addresses Prefix for the nodes in this Node Pool. `node_public_ip_enabled` should be `true`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"node_public_ip_enabled"},
		},

		// Node Taints control the behaviour of the Node Pool, as such they should not be computed and
		// must be specified/reconciled as required
		"node_taints": {
			Description: "A list of Kubernetes taints which should be applied to nodes in the agent pool (e.g. `key=value:NoSchedule`).",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"node_image_version": {
			Description: "The current node image version running on this Node Pool.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"orchestrator_version": {
			Description:      "Version of Kubernetes used for the Agents. If not specified, the latest recommended version will be used at provisioning time (but won't auto-upgrade). AKS does not require an exact patch version to be specified, minor version aliases such as `1.22` are also supported. - The minor version's latest GA patch is automatically chosen in that case. More details can be found in [the documentation](https://learn.microsoft.com/azure/aks/supported-kubernetes-versions?tabs=azure-cli#alias-minor-version).",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true, // azignore:AZS007 - pre-existing violation
			ValidateFunc:     validation.StringIsNotEmpty,
			DiffSuppressFunc: helpers.SuppressKubernetesVersionDiff,
		},

		"os_disk_size_gb": {
			Description:  "The Agent Operating System disk size in GB.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntAtLeast(1),
		},

		"os_disk_type": {
			Description:  "The type of disk which should be used for the Operating System.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      agentpools.OSDiskTypeManaged,
			ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForOSDiskType(), false),
		},

		"os_sku": {
			Description: "The OS SKU used by the agent pool. If not specified. Changing between `AzureLinux` and `Ubuntu` does not replace the resource; any other change forces a new resource to be created.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because defaults to Ubuntu if using Linux
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(agentpools.OSSKUAzureLinux),
				string(agentpools.OSSKUAzureLinuxThree),
				string(agentpools.OSSKUUbuntu),
				string(agentpools.OSSKUUbuntuTwoTwoZeroFour),
				string(agentpools.OSSKUUbuntuTwoFourZeroFour),
				string(agentpools.OSSKUWindowsTwoZeroOneNine),
				string(agentpools.OSSKUWindowsTwoZeroTwoTwo),
				string(agentpools.OSSKUWindowsTwoZeroTwoFive),
			}, false),
		},

		"os_type": {
			Description:  "The Operating System which should be used for this Node Pool.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(agentpools.OSTypeLinux),
			ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForOSType(), false),
		},

		"pod_ip_allocation_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"pod_subnet_id"},
			ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForPodIPAllocationMode(), false),
		},

		"pod_subnet_id": {
			Description:  "The ID of the Subnet where the pods in the Node Pool should exist.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"priority": {
			Description:  "The Priority for Virtual Machines within the Virtual Machine Scale Set that powers this Node Pool.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(agentpools.ScaleSetPriorityRegular),
			ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForScaleSetPriority(), false),
		},

		"proximity_placement_group_id": {
			Description:  "The ID of the Proximity Placement Group where the Virtual Machine Scale Set that powers this Node Pool will be placed.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: proximityplacementgroups.ValidateProximityPlacementGroupID,
		},

		"snapshot_id": {
			Description:  "The ID of the Snapshot which should be used to create this Node Pool.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: snapshots.ValidateSnapshotID,
		},

		"spot_max_price": {
			Description:  "The maximum price you're willing to pay in USD per Virtual Machine.",
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			ForceNew:     true,
			Default:      -1.0,
			ValidateFunc: computeValidate.SpotMaxPrice,
		},

		"scale_down_mode": {
			Description:  "How the node pool should deal with scaled-down nodes.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(agentpools.ScaleDownModeDelete),
			ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForScaleDownMode(), false),
		},

		"security": helpers.SchemaNodePoolSecurity("security"),

		"temporary_name_for_rotation": {
			Description:  "The name of the temporary node pool used to cycle the node pool when one of the relevant properties are updated.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: containerValidate.KubernetesAgentPoolName,
		},

		"ultra_ssd_enabled": {
			Description: "Used to specify whether the UltraSSD is enabled in the Node Pool. See [the documentation](https://learn.microsoft.com/azure/aks/use-ultra-disks) for more information.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"vnet_subnet_id": {
			Description:  "The ID of the Subnet where this Node Pool should exist.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"upgrade_settings": upgradeSettingsSchemaNodePoolResource(),

		"windows_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"outbound_nat_enabled": {
						Description: "Whether the Windows nodes in this Node Pool should have outbound NAT enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						ForceNew:    true,
						Default:     true,
					},
				},
			},
		},

		"workload_runtime": {
			Description:  "Used to specify the workload runtime.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForWorkloadRuntime(), false),
		},

		"zones": commonschema.ZonesMultipleOptional(),

		"auto_scaling_enabled": {
			Description: "Whether to enable [auto-scaler](https://learn.microsoft.com/azure/aks/cluster-autoscaler).",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"node_public_ip_enabled": {
			Description: "Whether each node should have a Public IP Address.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"host_encryption_enabled": {
			Description: "Whether the nodes in this Node Pool should have host encryption enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},
	}
}
