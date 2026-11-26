// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package orchestrated_virtual_machine_scale_set

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimages"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func orchestratedVirtualMachineScaleSetSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Virtual Machine Scale Set.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VirtualMachineName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"network_api_version": {
			Description:  "The Microsoft.Network API version used when creating networking resources in the Network Interface Configurations for Virtual Machine Scale Set.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      virtualmachinescalesets.NetworkApiVersionTwoZeroTwoZeroNegativeOneOneNegativeZeroOne,
			ValidateFunc: validation.StringInSlice(virtualmachinescalesets.PossibleValuesForNetworkApiVersion(), false),
			DiffSuppressFunc: func(_, old, new string, d *pluginsdk.ResourceData) bool {
				// This `DiffSuppressFunc` is used to keep compatible with the legacy Orchestrated VMSS and can be removed once the legacy VMSS is removed.
				if _, ok := d.GetOk("sku_name"); !ok {
					if old == "" && new == string(virtualmachinescalesets.NetworkApiVersionTwoZeroTwoZeroNegativeOneOneNegativeZeroOne) {
						return true
					}
				}

				return false
			},
		},

		"network_interface": helpers.OrchestratedVirtualMachineScaleSetNetworkInterfaceSchema(),

		"os_disk": helpers.OrchestratedVirtualMachineScaleSetOSDiskSchema(),

		"instances": {
			Description:  "The number of Virtual Machines in the Virtual Machine Scale Set.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(0, 1000),
		},

		// For sku I will create a format like: tier_sku name.
		// NOTE: all of the exposed vm sku tier's are Standard so this will continue to be hardcoded
		// Examples: Standard_HC44rs_4, Standard_D48_v3_6, Standard_M64s_20, Standard_HB120-96rs_v3_8
		"sku_name": {
			Description:  "The name of the SKU to be used by this Virtual Machine Scale Set.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.OrchestratedVirtualMachineScaleSetSku,
		},

		"sku_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allocation_strategy": {
						Description: "The allocation strategy for the Virtual Machine Scale Set based on which the VMs will be allocated.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice(
							virtualmachinescalesets.PossibleValuesForAllocationStrategy(),
							false,
						),
					},
					"virtual_machine_size": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						MaxItems: 5,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description:  "The name of the VM SKU which should be used for this Virtual Machine Scale Set, such as `Standard_B1ls` or `Standard_B2s`.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validate.SkuProfileVMSizeName,
								},
								"rank": {
									Description:  "The priority rank of the VM size. Lower values represent higher priority.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntBetween(1, 3),
								},
							},
						},
					},
				},
			},
		},

		"os_profile": helpers.OrchestratedVirtualMachineScaleSetOSProfileSchema(),

		// Optional
		// NOTE: The schema for the automatic instance repair has merged so they are
		// identical for both uniform and flex mode VMSS's
		"automatic_instance_repair": helpers.VirtualMachineScaleSetAutomaticRepairsPolicySchema(),

		"automatic_zone_rebalancing_enabled": {
			Description: "Whether to enable automatic zone balance. Requires `zones` to be set. When set to `true`, a health extension must be configured.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"boot_diagnostics": helpers.BootDiagnosticsSchema(),

		"capacity_reservation_group_id": {
			Description:  "The ID of the Capacity Reservation Group which the Virtual Machine Scale Set should be allocated to.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: capacityreservationgroups.ValidateCapacityReservationGroupID,
			ConflictsWith: []string{
				"proximity_placement_group_id",
			},
		},

		"data_disk": helpers.OrchestratedVirtualMachineScaleSetDataDiskSchema(),

		// Optional
		"additional_capabilities": helpers.OrchestratedVirtualMachineScaleSetAdditionalCapabilitiesSchema(),

		"encryption_at_host_enabled": {
			Description: "Whether disks attached to this Virtual Machine Scale Set should be encrypted by enabling Encryption at Host.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"eviction_policy": {
			Description: "The Policy which should be used by Spot Virtual Machines that are Evicted from the Scale Set.",
			// only applicable when `priority` is set to `Spot`
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(virtualmachinescalesets.PossibleValuesForVirtualMachineEvictionPolicyTypes(), false),
		},

		"extension_operations_enabled": {
			Description: "Whether extension operations should be allowed on the Virtual Machine Scale Set.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     true,
		},

		// Due to bug in RP extensions cannot currently be supported in Terraform ETA for full support is mid Jan 2022
		"extension": helpers.OrchestratedVirtualMachineScaleSetExtensionsSchema(),

		"extensions_time_budget": {
			Description:  "The time allotted for all extensions to start. The time duration should be between 15 minutes and 120 minutes (inclusive) and should be specified in ISO 8601 format.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT1H30M",
			ValidateFunc: validation.ISO8601DurationBetween("PT15M", "PT2H"),
		},

		// whilst the Swagger defines multiple at this time only UAI is supported
		"identity": commonschema.UserAssignedIdentityOptional(),

		"license_type": {
			Description: "The type of on-premise license (also known as Azure Hybrid Use Benefit) which should be used for this Virtual Machine Scale Set.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"None",
				"Windows_Client",
				"Windows_Server",
			}, false),
			DiffSuppressFunc: func(_, old, new string, _ *pluginsdk.ResourceData) bool {
				if old == "None" && new == "" || old == "" && new == "None" {
					return true
				}

				return false
			},
		},

		"max_bid_price": {
			Description:  "The maximum price you're willing to pay for each Virtual Machine in this Scale Set, in US Dollars; which must be greater than the current spot price. If this bid price falls below the current spot price the Virtual Machines in the Scale Set will be evicted using the eviction_policy.",
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      -1,
			ValidateFunc: validate.SpotMaxPrice,
		},

		"plan": helpers.PlanSchema(),

		"platform_fault_domain_count": {
			Description: "The number of fault domains that are used by this Virtual Machine Scale Set.",
			Type:        pluginsdk.TypeInt,
			Required:    true,
			ForceNew:    true,
		},

		"priority": {
			Description: "The Priority of this Virtual Machine Scale Set.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     string(virtualmachinescalesets.VirtualMachinePriorityTypesRegular),
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualmachinescalesets.VirtualMachinePriorityTypesRegular),
				string(virtualmachinescalesets.VirtualMachinePriorityTypesSpot),
			}, false),
		},

		"proximity_placement_group_id": {
			Description:  "The ID of the Proximity Placement Group which the Virtual Machine should be assigned to.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: proximityplacementgroups.ValidateProximityPlacementGroupID,
			// the Compute API is broken and returns the Resource Group name in UPPERCASE :shrug:, github issue: https://github.com/Azure/azure-rest-api-specs/issues/10016
			DiffSuppressFunc: suppress.CaseDifference,
			ConflictsWith: []string{
				"capacity_reservation_group_id",
			},
		},

		"resilient_virtual_machine_creation_enabled": {
			Description: "Whether `resilient virtual machine creation` is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"resilient_virtual_machine_deletion_enabled": {
			Description: "Whether `resilient virtual machine deletion` is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"rolling_upgrade_policy": helpers.VirtualMachineScaleSetRollingUpgradePolicySchema(),

		// NOTE: single_placement_group is now supported in orchestrated VMSS
		// Since null is now a valid value for this field there is no default
		// for this bool
		"single_placement_group": {
			Description: "Whether this Virtual Machine Scale Set should be limited to a Single Placement Group, which means the number of instances will be capped at 100 Virtual Machines.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"source_image_id": {
			Description: "The ID of an Image which each Virtual Machine in this Scale Set should be based on. Possible Image ID types include `Image ID`s, `Shared Image ID`s, `Shared Image Version ID`s, `Community Gallery Image ID`s, `Community Gallery Image Version ID`s, `Shared Gallery Image ID`s and `Shared Gallery Image Version ID`s.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.Any(
				images.ValidateImageID,
				validation.AsGeneratedID(galleryimages.ParseGalleryImageIDInsensitively),
				validation.AsGeneratedID(galleryimageversions.ParseImageVersionIDInsensitively),
				validate.CommunityGalleryImageID,
				validate.CommunityGalleryImageVersionID,
				validate.SharedGalleryImageID,
				validate.SharedGalleryImageVersionID,
			),
			ConflictsWith: []string{
				"source_image_reference",
			},
		},

		"source_image_reference": helpers.SourceImageReferenceSchemaOrchestratedVMSS(),

		"zone_balance": {
			Description: "Whether the Virtual Machines in this Scale Set should be strictly evenly distributed across Availability Zones.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"termination_notification": helpers.OrchestratedVirtualMachineScaleSetTerminationNotificationSchema(),

		"zones": commonschema.ZonesMultipleOptional(),

		"tags": commonschema.Tags(),

		// Computed
		"unique_id": {
			Description: "The Unique ID for the Orchestrated Virtual Machine Scale Set.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"upgrade_mode": {
			Description:  "How upgrades (e.g. changing the Image/SKU) should be performed to Virtual Machine Instances.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(virtualmachinescalesets.UpgradeModeManual),
			ValidateFunc: validation.StringInSlice(virtualmachinescalesets.PossibleValuesForUpgradeMode(), false),
		},

		"user_data_base64": {
			Description:  "The Base64-Encoded User Data which should be used for this Virtual Machine Scale Set.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsBase64,
		},

		"priority_mix": helpers.OrchestratedVirtualMachineScaleSetPriorityMixPolicySchema(),
	}
}
