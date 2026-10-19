// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_virtual_machine

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimages"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func windowsVirtualMachineSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VirtualMachineName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"admin_password": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			Sensitive:        true,
			DiffSuppressFunc: helpers.AdminPasswordDiffSuppressFunc,
			RequiredWith: []string{
				"admin_username",
			},
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
			ValidateFunc: validate.WindowsAdminPassword,
		},

		"admin_username": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			RequiredWith: []string{
				"admin_password",
			},
			ExactlyOneOf: []string{
				"admin_username",
				"os_managed_disk_id",
			},
			ValidateFunc: validate.WindowsAdminUsername,
		},

		"network_interface_ids": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: commonids.ValidateNetworkInterfaceID,
			},
		},

		"os_disk": helpers.VirtualMachineOSDiskSchema(),

		"os_managed_disk_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C as this is the same value as `os_disk.0.id` - which gains a value from implicit
			// disk creation with a VM when an existing disk is not specified here. This is a top-level property
			// to enable schema validation to guard against any values for `OsProfile` being set, as these are
			// incompatible with specifying an existing disk. i.e. the OsProfile becomes unmanageable.
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateManagedDiskID,
			ExactlyOneOf: []string{
				"os_managed_disk_id",
				"source_image_id",
				"source_image_reference",
			},
		},

		"size": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"additional_capabilities": helpers.VirtualMachineAdditionalCapabilitiesSchema(),

		"additional_unattend_content": helpers.AdditionalUnattendContentSchemaVM(),

		"allow_extension_operations": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"availability_set_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateAvailabilitySetID,
			// the Compute/VM API is broken and returns the Availability Set name in UPPERCASE :shrug:
			// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
			DiffSuppressFunc: suppress.CaseDifference,
			ConflictsWith: []string{
				"capacity_reservation_group_id",
				"virtual_machine_scale_set_id",
				"zone",
			},
		},

		"boot_diagnostics": helpers.BootDiagnosticsSchema(),

		"bypass_platform_safety_checks_on_user_schedule_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"capacity_reservation_group_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE
			// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     capacityreservationgroups.ValidateCapacityReservationGroupID,
			ConflictsWith: []string{
				"availability_set_id",
				"proximity_placement_group_id",
			},
		},

		"computer_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C since we reuse the VM name if one's not specified
			Computed: true,
			ForceNew: true,

			ValidateFunc: validate.WindowsComputerNameFull,
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"custom_data": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsBase64,
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"dedicated_host_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateDedicatedHostID,
			// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE :shrug:
			// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
			DiffSuppressFunc: suppress.CaseDifference,
			ConflictsWith: []string{
				"dedicated_host_group_id",
			},
		},

		"dedicated_host_group_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateDedicatedHostGroupID,
			// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE
			// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
			DiffSuppressFunc: suppress.CaseDifference,
			ConflictsWith: []string{
				"dedicated_host_id",
			},
		},

		"disk_controller_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(virtualmachines.PossibleValuesForDiskControllerTypes(), false),
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"automatic_updates_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true, // updating this is not allowed "Changing property 'windowsConfiguration.enableAutomaticUpdates' is not allowed." Target="windowsConfiguration.enableAutomaticUpdates"
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"encryption_at_host_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"eviction_policy": {
			// only applicable when `priority` is set to `Spot`
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(virtualmachines.PossibleValuesForVirtualMachineEvictionPolicyTypes(), false),
		},

		"extensions_time_budget": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT1H30M",
			ValidateFunc: validation.ISO8601DurationBetween("PT15M", "PT2H"),
		},

		"gallery_application": helpers.VirtualMachineGalleryApplicationSchema(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

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
			ValidateFunc: validation.FloatAtLeast(-1.0),
		},

		"patch_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(virtualmachines.PossibleValuesForWindowsVMGuestPatchMode(), false),
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"patch_assessment_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(virtualmachines.PossibleValuesForWindowsPatchAssessmentMode(), false),
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"hotpatching_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"plan": helpers.PlanSchema(),

		"priority": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(virtualmachines.VirtualMachinePriorityTypesRegular),
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualmachines.VirtualMachinePriorityTypesRegular),
				string(virtualmachines.VirtualMachinePriorityTypesSpot),
			}, false),
		},

		"provision_vm_agent": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"proximity_placement_group_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: proximityplacementgroups.ValidateProximityPlacementGroupID,
			// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE :shrug:
			// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
			DiffSuppressFunc: suppress.CaseDifference,
			ConflictsWith: []string{
				"capacity_reservation_group_id",
			},
		},

		"reboot_setting": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualmachines.WindowsVMGuestPatchAutomaticByPlatformRebootSettingAlways),
				string(virtualmachines.WindowsVMGuestPatchAutomaticByPlatformRebootSettingIfRequired),
				string(virtualmachines.WindowsVMGuestPatchAutomaticByPlatformRebootSettingNever),
			}, false),
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"secret": helpers.WindowsSecretSchemaVM(),

		"secure_boot_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"source_image_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
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
				"os_managed_disk_id",
				"source_image_id",
				"source_image_reference",
			},
		},

		"source_image_reference": helpers.SourceImageReferenceSchemaVM(),

		"tags": commonschema.Tags(),

		"os_image_notification": helpers.VirtualMachineOsImageNotificationSchema(),

		"termination_notification": helpers.VirtualMachineTerminationNotificationSchema(),

		"timezone": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.VirtualMachineTimeZone(),
		},

		"virtual_machine_scale_set_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ConflictsWith: []string{
				"availability_set_id",
			},
			ValidateFunc: commonids.ValidateVirtualMachineScaleSetID,
		},

		"platform_fault_domain": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      -1,
			ForceNew:     true,
			RequiredWith: []string{"virtual_machine_scale_set_id"},
			ValidateFunc: validation.IntAtLeast(-1),
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

		"winrm_listener": helpers.WinRmListenerSchemaVM(),

		"zone": commonschema.ZoneSingleOptionalForceNew(),

		// Computed
		"private_ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"private_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		"public_ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"public_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		"virtual_machine_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"vm_agent_platform_updates_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
	}
}
