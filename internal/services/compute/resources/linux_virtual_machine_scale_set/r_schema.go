// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package linux_virtual_machine_scale_set

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimages"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/base64"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceLinuxVirtualMachineScaleSetSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VirtualMachineName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		// Required
		"admin_username": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_interface": helpers.VirtualMachineScaleSetNetworkInterfaceSchema(),

		"os_disk": helpers.VirtualMachineScaleSetOSDiskSchema(),

		"instances": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"sku": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// Optional
		"additional_capabilities": helpers.VirtualMachineScaleSetAdditionalCapabilitiesSchema(),

		"admin_password": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			Sensitive:        true,
			DiffSuppressFunc: helpers.AdminPasswordDiffSuppressFunc,
		},

		"admin_ssh_key": helpers.SSHKeysSchema(false),

		"automatic_os_upgrade_policy": helpers.VirtualMachineScaleSetAutomatedOSUpgradePolicySchema(),

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

		"computer_name_prefix": {
			Type:     pluginsdk.TypeString,
			Optional: true,

			// Computed since we reuse the VM name if one's not specified
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,

			ValidateFunc: validate.LinuxComputerNamePrefix,
		},

		"custom_data": base64.OptionalSchema(false),

		"data_disk": helpers.VirtualMachineScaleSetDataDiskSchema(),

		"disable_password_authentication": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"do_not_run_extensions_on_overprovisioned_machines": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

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

		"extension": helpers.VirtualMachineScaleSetExtensionsSchema(),

		"extensions_time_budget": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT1H30M",
			ValidateFunc: validation.ISO8601DurationBetween("PT15M", "PT2H"),
		},

		"gallery_application": helpers.VirtualMachineScaleSetGalleryApplicationSchema(),

		"health_probe_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"host_group_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE
			// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validation.AsGeneratedID(commonids.ParseDedicatedHostGroupIDInsensitively),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"max_bid_price": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      -1,
			ValidateFunc: validate.SpotMaxPrice,
		},

		"overprovision": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"plan": helpers.PlanSchema(),

		"platform_fault_domain_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			ForceNew: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
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

		"provision_vm_agent": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
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

		"resilient_vm_creation_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"resilient_vm_deletion_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"rolling_upgrade_policy": helpers.VirtualMachineScaleSetRollingUpgradePolicySchema(),

		"secret": helpers.LinuxSecretSchema(),

		"secure_boot_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"single_placement_group": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
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
			ExactlyOneOf: []string{
				"source_image_id",
				"source_image_reference",
			},
		},

		"source_image_reference": helpers.SourceImageReferenceSchema(false),

		"tags": commonschema.Tags(),

		"upgrade_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(virtualmachinescalesets.UpgradeModeManual),
			ValidateFunc: validation.StringInSlice(virtualmachinescalesets.PossibleValuesForUpgradeMode(), false),
		},

		"user_data": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsBase64,
		},

		"vtpm_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"zone_balance": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"scale_in": helpers.VirtualMachineScaleSetScaleInPolicySchema(),

		"spot_restore": helpers.VirtualMachineScaleSetSpotRestorePolicySchema(),

		"termination_notification": helpers.VirtualMachineScaleSetTerminationNotificationSchema(),

		"zones": commonschema.ZonesMultipleOptional(),

		// Computed
		"unique_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
