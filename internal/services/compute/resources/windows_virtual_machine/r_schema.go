// Copyright IBM Corp. 2014, 2026
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func windowsVirtualMachineSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Windows Virtual Machine.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VirtualMachineName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"admin_password": {
			Description:      "The Password which should be used for the local-administrator on this Virtual Machine.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			Sensitive:        true,
			ValidateFunc:     validate.WindowsAdminPassword,
			DiffSuppressFunc: helpers.AdminPasswordDiffSuppressFunc,
			RequiredWith: []string{
				"admin_username",
			},
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"admin_username": {
			Description:  "The username of the local administrator used for the Virtual Machine.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.WindowsAdminUsername,
			RequiredWith: []string{
				"admin_password",
			},
			ExactlyOneOf: []string{
				"admin_username",
				"os_managed_disk_id",
			},
		},

		"network_interface_ids": {
			Description: "A list of Network Interface IDs which should be attached to this Virtual Machine. The first Network Interface ID in this list will be the Primary Network Interface on the Virtual Machine.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: commonids.ValidateNetworkInterfaceID,
			},
		},

		"os_disk": helpers.VirtualMachineOSDiskSchema(),

		"os_managed_disk_id": {
			Description: "The ID of an existing Managed Disk to use as the OS Disk for this Windows Virtual Machine.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
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
			Description:  "The SKU which should be used for this Virtual Machine, such as `Standard_D4_v5`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"additional_capabilities": helpers.VirtualMachineAdditionalCapabilitiesSchema(),

		"additional_unattend_content": helpers.AdditionalUnattendContentSchemaVM(),

		"allow_extension_operations": {
			Description: "Whether Extension Operations should be allowed on this Virtual Machine.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"availability_set_id": {
			Description:  "The ID of the Availability Set in which the Virtual Machine should exist.",
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
			Description: "Whether to skip platform scheduled patching when a user schedule is associated with the VM.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"capacity_reservation_group_id": {
			Description:  "The ID of the Capacity Reservation Group which the Virtual Machine should be allocated to.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: capacityreservationgroups.ValidateCapacityReservationGroupID,
			// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE
			// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
			DiffSuppressFunc: suppress.CaseDifference,
			ConflictsWith: []string{
				"availability_set_id",
				"proximity_placement_group_id",
			},
		},

		"computer_name": {
			Description: "The Hostname which should be used for this Virtual Machine. If the value of the `name` field is not a valid `computer_name`, then you must specify `computer_name`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C since we reuse the VM name if one's not specified
			Computed: true,
			ForceNew: true,

			ValidateFunc: validate.WindowsComputerNameFull,
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"custom_data": {
			Description:  "The Base64-Encoded Custom Data which should be used for this Virtual Machine.",
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
			Description:  "The ID of a Dedicated Host where this machine should be run on.",
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
			Description:  "The ID of a Dedicated Host Group that this Windows Virtual Machine should be run within.",
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
			Description:  "The Disk Controller Type used for this Virtual Machine.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(virtualmachines.PossibleValuesForDiskControllerTypes(), false),
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"automatic_updates_enabled": {
			Description: "Whether Automatic Updates are Enabled for the Windows Virtual Machine.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true, // updating this is not allowed "Changing property 'windowsConfiguration.enableAutomaticUpdates' is not allowed." Target="windowsConfiguration.enableAutomaticUpdates"
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"encryption_at_host_enabled": {
			Description: "Whether all of the disks (including the temp disk) attached to this Virtual Machine should be encrypted by enabling Encryption at Host.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"eviction_policy": {
			Description: "What should happen when the Virtual Machine is evicted for price reasons when using a Spot instance.",
			// only applicable when `priority` is set to `Spot`
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(virtualmachines.PossibleValuesForVirtualMachineEvictionPolicyTypes(), false),
		},

		"extensions_time_budget": {
			Description:  "The duration allocated for all extensions to start. The time duration should be between 15 minutes and 120 minutes (inclusive) and should be specified in ISO 8601 format.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT1H30M",
			ValidateFunc: validation.ISO8601DurationBetween("PT15M", "PT2H"),
		},

		"gallery_application": helpers.VirtualMachineGalleryApplicationSchema(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"license_type": {
			Description: "The type of on-premise license (also known as [Azure Hybrid Use Benefit](https://learn.microsoft.com/windows-server/get-started/azure-hybrid-benefit)) which should be used for this Virtual Machine.",
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
			Description:  "The maximum price you're willing to pay for this Virtual Machine, in US Dollars; which must be greater than the current spot price. If this bid price falls below the current spot price the Virtual Machine will be evicted using the `eviction_policy`.",
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      -1,
			ValidateFunc: validation.FloatAtLeast(-1.0),
		},

		"patch_mode": {
			Description:  "The mode of in-guest patching to this Windows Virtual Machine. For more information on patch modes see the [product documentation](https://learn.microsoft.com/azure/virtual-machines/automatic-vm-guest-patching#patch-orchestration-modes).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(virtualmachines.PossibleValuesForWindowsVMGuestPatchMode(), false),
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"patch_assessment_mode": {
			Description:  "The mode of VM Guest Patching for the Virtual Machine.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(virtualmachines.PossibleValuesForWindowsPatchAssessmentMode(), false),
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"hotpatching_enabled": {
			Description: "Whether the VM should be patched without requiring a reboot. For more information about hot patching see the [product documentation](https://learn.microsoft.com/azure/automanage/automanage-hotpatch).",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"plan": helpers.PlanSchema(),

		"priority": {
			Description: "The priority of this Virtual Machine.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     string(virtualmachines.VirtualMachinePriorityTypesRegular),
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualmachines.VirtualMachinePriorityTypesRegular),
				string(virtualmachines.VirtualMachinePriorityTypesSpot),
			}, false),
		},

		"provision_vm_agent": {
			Description: "Whether the Azure VM Agent should be provisioned on this Virtual Machine.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
			ConflictsWith: []string{
				"os_managed_disk_id",
			},
		},

		"proximity_placement_group_id": {
			Description:  "The ID of the Proximity Placement Group which the Virtual Machine should be assigned to.",
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
			Description: "The reboot setting for platform scheduled patching.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
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
			Description: "Whether Secure Boot and Trusted Launch is enabled for the Virtual Machine.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"source_image_id": {
			Description: "The ID of the Image which this Virtual Machine should be created from. Possible Image ID types include `Image ID`s, `Shared Image ID`s, `Shared Image Version ID`s, `Community Gallery Image ID`s, `Community Gallery Image Version ID`s, `Shared Gallery Image ID`s and `Shared Gallery Image Version ID`s.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
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
			Description:  "The Time Zone which should be used by the Virtual Machine.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.VirtualMachineTimeZone(),
		},

		"virtual_machine_scale_set_id": {
			Description:  "The Orchestrated Virtual Machine Scale Set that this Virtual Machine should be created within.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateVirtualMachineScaleSetID,
			ConflictsWith: []string{
				"availability_set_id",
			},
		},

		"platform_fault_domain": {
			Description:  "The Platform Fault Domain in which this Windows Virtual Machine should be created.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      -1,
			ValidateFunc: validation.IntAtLeast(-1),
			RequiredWith: []string{"virtual_machine_scale_set_id"},
		},

		"user_data": {
			Description:  "The Base64-Encoded User Data which should be used for this Virtual Machine.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsBase64,
		},

		"vtpm_enabled": {
			Description: "Whether vTPM (virtual Trusted Platform Module) and Trusted Launch is enabled for the Virtual Machine.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"winrm_listener": helpers.WinRmListenerSchemaVM(),

		"zone": commonschema.ZoneSingleOptionalForceNew(),

		// Computed
		"private_ip_address": {
			Description: "The Primary Private IP Address assigned to this Virtual Machine.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"private_ip_addresses": {
			Description: "A list of Private IP Addresses assigned to this Virtual Machine.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		"public_ip_address": {
			Description: "The Primary Public IP Address assigned to this Virtual Machine.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"public_ip_addresses": {
			Description: "A list of the Public IP Addresses assigned to this Virtual Machine.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		"virtual_machine_id": {
			Description: "A 128-bit identifier which uniquely identifies this Virtual Machine.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"vm_agent_platform_updates_enabled": {
			Description: "Whether Virtual Machine Agent Platform Updates are `enabled` on this Virtual Machine.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},
	}
}
