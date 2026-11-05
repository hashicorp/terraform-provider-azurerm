// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VirtualMachineName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"network_api_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(virtualmachinescalesets.PossibleValuesForNetworkApiVersion(), false),
			Default:      virtualmachinescalesets.NetworkApiVersionTwoZeroTwoZeroNegativeOneOneNegativeZeroOne,
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
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(0, 1000),
		},

		// For sku I will create a format like: tier_sku name.
		// NOTE: all of the exposed vm sku tier's are Standard so this will continue to be hardcoded
		// Examples: Standard_HC44rs_4, Standard_D48_v3_6, Standard_M64s_20, Standard_HB120-96rs_v3_8
		"sku_name": {
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
						Type:     pluginsdk.TypeString,
						Required: true,
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
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validate.SkuProfileVMSizeName,
								},
								"rank": {
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

		"boot_diagnostics": helpers.BootDiagnosticsSchema(),

		"capacity_reservation_group_id": {
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
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"eviction_policy": {
			// only applicable when `priority` is set to `Spot`
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(virtualmachinescalesets.PossibleValuesForVirtualMachineEvictionPolicyTypes(), false),
		},

		"extension_operations_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
		},

		// Due to bug in RP extensions cannot currently be supported in Terraform ETA for full support is mid Jan 2022
		"extension": helpers.OrchestratedVirtualMachineScaleSetExtensionsSchema(),

		"extensions_time_budget": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT1H30M",
			ValidateFunc: validation.ISO8601DurationBetween("PT15M", "PT2H"),
		},

		// whilst the Swagger defines multiple at this time only UAI is supported
		"identity": commonschema.UserAssignedIdentityOptional(),

		"license_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
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
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      -1,
			ValidateFunc: validate.SpotMaxPrice,
		},

		"plan": helpers.PlanSchema(),

		"platform_fault_domain_count": {
			Type:     pluginsdk.TypeInt,
			Required: true,
			ForceNew: true,
		},

		"priority": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(virtualmachinescalesets.VirtualMachinePriorityTypesRegular),
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualmachinescalesets.VirtualMachinePriorityTypesRegular),
				string(virtualmachinescalesets.VirtualMachinePriorityTypesSpot),
			}, false),
		},

		"proximity_placement_group_id": {
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

		"rolling_upgrade_policy": helpers.VirtualMachineScaleSetRollingUpgradePolicySchema(),

		// NOTE: single_placement_group is now supported in orchestrated VMSS
		// Since null is now a valid value for this field there is no default
		// for this bool
		"single_placement_group": {
			Type:     pluginsdk.TypeBool,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Optional: true,
		},

		"source_image_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
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
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"termination_notification": helpers.OrchestratedVirtualMachineScaleSetTerminationNotificationSchema(),

		"zones": commonschema.ZonesMultipleOptional(),

		"tags": commonschema.Tags(),

		// Computed
		"unique_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"upgrade_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(virtualmachinescalesets.UpgradeModeManual),
			ValidateFunc: validation.StringInSlice(virtualmachinescalesets.PossibleValuesForUpgradeMode(), false),
		},

		"user_data_base64": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsBase64,
		},

		"priority_mix": helpers.OrchestratedVirtualMachineScaleSetPriorityMixPolicySchema(),
	}
}
