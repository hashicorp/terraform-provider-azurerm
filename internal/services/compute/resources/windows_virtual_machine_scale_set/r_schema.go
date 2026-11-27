// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package windows_virtual_machine_scale_set

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

func resourceWindowsVirtualMachineScaleSetSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Windows Virtual Machine Scale Set.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VirtualMachineName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"admin_username": {
			Description:  "The username of the local administrator on each Virtual Machine Scale Set instance.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"admin_password": {
			Description:      "The Password which should be used for the local-administrator on this Virtual Machine.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			Sensitive:        true,
			ValidateFunc:     validation.StringIsNotEmpty,
			DiffSuppressFunc: helpers.AdminPasswordDiffSuppressFunc,
		},

		"network_interface": helpers.VirtualMachineScaleSetNetworkInterfaceSchema(),

		"os_disk": helpers.VirtualMachineScaleSetOSDiskSchema(),

		"instances": {
			Description:  "The number of Virtual Machines in the Scale Set.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"sku": {
			Description:  "The Virtual Machine SKU for the Scale Set, such as `Standard_D4_v5`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"additional_capabilities": helpers.VirtualMachineScaleSetAdditionalCapabilitiesSchema(),

		"additional_unattend_content": helpers.AdditionalUnattendContentSchema(),

		"automatic_os_upgrade_policy": helpers.VirtualMachineScaleSetAutomatedOSUpgradePolicySchema(),

		"automatic_instance_repair": helpers.VirtualMachineScaleSetAutomaticRepairsPolicySchema(),

		"automatic_zone_rebalancing_enabled": {
			Description: "Whether to enable automatic zone balance. Requires `zones` to be set. When set to `true`, a `health_probe_id` or a health extension must be configured.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"boot_diagnostics": helpers.BootDiagnosticsSchema(),

		"capacity_reservation_group_id": {
			Description:  "The ID of the Capacity Reservation Group which the Virtual Machine Scale Set should be allocated to.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: capacityreservationgroups.ValidateCapacityReservationGroupID,
			ConflictsWith: []string{
				"proximity_placement_group_id",
			},
		},

		"computer_name_prefix": {
			Description: "The prefix which should be used for the name of the Virtual Machines in this Scale Set. If the value of the `name` field is not a valid `computer_name_prefix`, then you must specify `computer_name_prefix`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C since we reuse the VM name if one's not specified
			Computed: true,
			ForceNew: true,

			ValidateFunc: validate.WindowsComputerNamePrefix,
		},

		"custom_data": base64.OptionalSchema(false),

		"data_disk": helpers.VirtualMachineScaleSetDataDiskSchema(),

		"disk_controller_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because Azure assigns a default disk controller type based on the VM size and image when not specified
			Computed:     true,
			ValidateFunc: validation.StringInSlice(virtualmachinescalesets.PossibleValuesForDiskControllerTypes(), false),
		},

		"do_not_run_extensions_on_overprovisioned_machines": {
			Description: "Whether Virtual Machine Extensions should be run on Overprovisioned Virtual Machines in the Scale Set.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"automatic_updates_enabled": {
			Description: "Whether to enable automatic update for this Virtual Machine.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"encryption_at_host_enabled": {
			Description: "Whether all of the disks (including the temp disk) attached to this Virtual Machine should be encrypted by enabling Encryption at Host.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"eviction_policy": {
			Description: "The eviction policy for Virtual Machines in this Scale Set.",
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

		"extension": helpers.VirtualMachineScaleSetExtensionsSchema(),

		"extensions_time_budget": {
			Description:  "The duration allocated for all extensions to start. The time duration should be between `15` minutes and `120` minutes (inclusive) and should be specified in ISO 8601 format.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT1H30M",
			ValidateFunc: validation.ISO8601DurationBetween("PT15M", "PT2H"),
		},

		"gallery_application": helpers.VirtualMachineScaleSetGalleryApplicationSchema(),

		"health_probe_id": {
			Description:  "The ID of a Load Balancer Probe which should be used to determine the health of an instance.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"host_group_id": {
			Description:  "The ID of the Dedicated Host group that the Virtual Machine Scale Set resides in.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseDedicatedHostGroupIDInsensitively),
			// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE
			// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"license_type": {
			Description: "The type of on-premise license (also known as [Azure Hybrid Use Benefit](https://learn.microsoft.com/windows-server/get-started/azure-hybrid-benefit)) which should be used for this Virtual Machine Scale Set.",
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
			Description:  "The maximum price you're willing to pay for each Virtual Machine in this Scale Set, in US Dollars; which must be greater than the current spot price. If this bid price falls below the current spot price the Virtual Machines in the Scale Set will be evicted using the `eviction_policy`.",
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      -1,
			ValidateFunc: validate.SpotMaxPrice,
		},

		"overprovision": {
			Description: "Whether Azure should over-provision Virtual Machines in this Scale Set. This means that multiple Virtual Machines will be provisioned and Azure will keep the instances which become available first - which improves provisioning success rates and improves deployment time. You're not billed for these over-provisioned VM's and they don't count towards the Subscription Quota.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"plan": helpers.PlanSchema(),

		"platform_fault_domain_count": {
			Description: "The number of fault domains that are used by this Linux Virtual Machine Scale Set.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
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

		"provision_vm_agent": {
			Description: "Whether the Azure VM Agent should be provisioned on each Virtual Machine in the Scale Set.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     true,
		},

		"proximity_placement_group_id": {
			Description:  "The ID of the Proximity Placement Group in which the Virtual Machine Scale Set should be assigned to.",
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

		"secret": helpers.WindowsSecretSchema(),

		"secure_boot_enabled": {
			Description: "Whether Secure Boot and Trusted Launch is enabled for the Virtual Machine.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"single_placement_group": {
			Description: "Whether this Virtual Machine Scale Set should be limited to a Single Placement Group, which means the number of instances will be capped at 100 Virtual Machines.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"source_image_id": {
			Description: "The ID of an Image which each Virtual Machine in this Scale Set should be based on. Possible Image ID types include `Image ID`, `Shared Image ID`, `Shared Image Version ID`, `Community Gallery Image ID`, `Community Gallery Image Version ID`, `Shared Gallery Image ID` and `Shared Gallery Image Version ID`.",
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
			ExactlyOneOf: []string{
				"source_image_id",
				"source_image_reference",
			},
		},

		"source_image_reference": helpers.SourceImageReferenceSchema(false),

		"tags": commonschema.Tags(),

		"timezone": {
			Description:  "The time zone of the Virtual Machine.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.VirtualMachineTimeZone(),
		},

		"upgrade_mode": {
			Description:  "How Upgrades (e.g. changing the Image/SKU) should be performed to Virtual Machine Instances.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(virtualmachinescalesets.UpgradeModeManual),
			ValidateFunc: validation.StringInSlice(virtualmachinescalesets.PossibleValuesForUpgradeMode(), false),
		},

		"user_data": {
			Description:  "The Base64-Encoded User Data which should be used for this Virtual Machine Scale Set.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsBase64,
		},

		"vtpm_enabled": {
			Description: "Whether vTPM (Virtual Trusted Platform Module) and Trusted Launch is enabled for the Virtual Machine.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"winrm_listener": helpers.WinRmListenerSchema(),

		"zone_balance": {
			Description: "Whether the Virtual Machines in this Scale Set should be strictly evenly distributed across Availability Zones.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"scale_in": helpers.VirtualMachineScaleSetScaleInPolicySchema(),

		"spot_restore": helpers.VirtualMachineScaleSetSpotRestorePolicySchema(),

		"resilient_vm_creation_enabled": {
			Description: "Whether `resilient VM creation` is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"resilient_vm_deletion_enabled": {
			Description: "Whether `resilient VM deletion` is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"termination_notification": helpers.VirtualMachineScaleSetTerminationNotificationSchema(),

		"zones": commonschema.ZonesMultipleOptional(),

		"unique_id": {
			Description: "The Unique ID for this Windows Virtual Machine Scale Set.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
