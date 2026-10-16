// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_node_pool

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/snapshots"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/helpers"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceKubernetesClusterNodePoolSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: containerValidate.KubernetesAgentPoolName,
		},

		"kubernetes_cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseKubernetesClusterIDInsensitively),
		},

		"node_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(0, 1000),
		},

		"tags": commonschema.Tags(),

		"vm_size": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C AKS RP provides a new feature that will automatically select an available vm size when it's omitted.
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"host_group_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseDedicatedHostGroupIDInsensitively),
		},

		// Optional
		"capacity_reservation_group_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: capacityreservationgroups.ValidateCapacityReservationGroupID,
		},

		"eviction_policy": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForScaleSetEvictionPolicy(), false),
		},

		"kubelet_config": helpers.SchemaNodePoolKubeletConfig(),

		"linux_os_config": helpers.SchemaNodePoolLinuxOSConfig(),

		"fips_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"gpu_instance": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(agentpools.GPUInstanceProfileMIGOneg),
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
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForKubeletDiskType(), false),
		},

		"max_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 1000),
		},

		"max_pods": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(agentpools.AgentPoolModeUser),
			ValidateFunc: validation.StringInSlice([]string{
				string(agentpools.AgentPoolModeSystem),
				string(agentpools.AgentPoolModeUser),
			}, false),
		},

		"min_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: rather than setting `0` users should instead pass `null` here
			ValidateFunc: validation.IntBetween(0, 1000),
		},

		"node_network_profile": helpers.SchemaNodePoolNetworkProfile(),

		"node_labels": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"node_public_ip_prefix_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"node_public_ip_enabled"},
		},

		// Node Taints control the behaviour of the Node Pool, as such they should not be computed and
		// must be specified/reconciled as required
		"node_taints": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"node_image_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"orchestrator_version": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true, // azignore:AZS007 - pre-existing violation
			ValidateFunc:     validation.StringIsNotEmpty,
			DiffSuppressFunc: helpers.SuppressKubernetesVersionDiff,
		},

		"os_disk_size_gb": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntAtLeast(1),
		},

		"os_disk_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      agentpools.OSDiskTypeManaged,
			ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForOSDiskType(), false),
		},

		"os_sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
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
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(agentpools.OSTypeLinux),
			ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForOSType(), false),
		},

		"pod_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"priority": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(agentpools.ScaleSetPriorityRegular),
			ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForScaleSetPriority(), false),
		},

		"proximity_placement_group_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: proximityplacementgroups.ValidateProximityPlacementGroupID,
		},

		"snapshot_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: snapshots.ValidateSnapshotID,
		},

		"spot_max_price": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			ForceNew:     true,
			Default:      -1.0,
			ValidateFunc: computeValidate.SpotMaxPrice,
		},

		"scale_down_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(agentpools.ScaleDownModeDelete),
			ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForScaleDownMode(), false),
		},

		"temporary_name_for_rotation": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: containerValidate.KubernetesAgentPoolName,
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

		"upgrade_settings": upgradeSettingsSchemaNodePoolResource(),

		"windows_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"outbound_nat_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  true,
					},
				},
			},
		},

		"workload_runtime": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForWorkloadRuntime(), false),
		},

		"zones": commonschema.ZonesMultipleOptional(),

		"auto_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"node_public_ip_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"host_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}
}
