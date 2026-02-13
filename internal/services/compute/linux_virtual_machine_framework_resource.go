// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	fwcommonschema "github.com/hashicorp/go-azure-helpers/framework/commonschema"
	fwidentity "github.com/hashicorp/go-azure-helpers/framework/identity"
	fwlocation "github.com/hashicorp/go-azure-helpers/framework/location"
	"github.com/hashicorp/go-azure-helpers/framework/planmodifiers"
	"github.com/hashicorp/go-azure-helpers/framework/planmodifiers/casing"
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/framework/values"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/custompoller"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.FrameworkWrappedResourceWithUpdate = &linuxVirtualMachineResource{}
var _ sdk.FrameworkWrappedResourceWithConfigValidators = &linuxVirtualMachineResource{}
var _ sdk.FrameworkWrappedResourceWithPlanModifier = &linuxVirtualMachineResource{}

type linuxVirtualMachineResource struct{}

func (l *linuxVirtualMachineResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, _ sdk.ResourceMetadata, decodedPlan any, decodedConfig any) {
	plan := sdk.AssertResourceModelType[linuxVirtualMachineResourceModel](decodedPlan, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	config := sdk.AssertResourceModelType[linuxVirtualMachineResourceModel](decodedConfig, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.AllowExtensionOperations.ValueBool() && !plan.ProvisionVMAgent.ValueBool() {
		sdk.SetResponseErrorDiagnostic(resp, "incompatible config set", "`allow_extension_operations` cannot be set to `true` when `provision_vm_agent` is set to `false`")
		return
	}

	// TODO - Suppress diff/Replace in os_disk.0.storage_account_type when using `managed_disk_id`
	if !config.OSManagedDiskID.IsNull() {
		saTypePath := path.Root("os_disk").AtListIndex(0).AtMapKey("storage_account_type")
		if resp.RequiresReplace.Contains(saTypePath) {
			var newRR path.Paths
			for _, v := range resp.RequiresReplace {
				if !v.Equal(saTypePath) {
					newRR = append(newRR, v)
				}
			}
			resp.RequiresReplace = newRR
		}

		// TODO - Is it enough to remove it from RequiresReplace or do we need to also set plan values?
		// osdiskP, diags := decodeLinuxVirtualMachineOSDiskModel(ctx, plan.OSDisk)
		// if diags.HasError() {
		// 	sdk.AppendResponseErrorDiagnostic(response, diags)
		// 	return
		// }
		// osdiskS, diags := decodeLinuxVirtualMachineOSDiskModel(ctx, state.OSDisk)
		// if diags.HasError() {
		// 	sdk.AppendResponseErrorDiagnostic(response, diags)
		// 	return
		// }
		//
		// osdiskP.StorageAccountType = osdiskS.StorageAccountType
		// // set back to response.Plan
	}

}

func (l *linuxVirtualMachineResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		// TODO - "at least one `admin_ssh_key` must be specified when `disable_password_authentication` is set to `true`"
		// TODO - "an `admin_password` must be specified if `disable_password_authentication` is set to `false`"
		// TODO - "patch_mode cannot be set to "AutomaticByPlatform" when "provision_vm_agent" is set to false"
		// TODO - "`patch_mode` must be set to `AutomaticByPlatform` when `bypass_platform_safety_checks_on_user_schedule_enabled` is set to `true`"
		// TODO - "`patch_mode` must be set to `AutomaticByPlatform` when `reboot_setting` is specified"
		// TODO - "`encryption_at_host_enabled` cannot be set to `true` when `os_disk.0.security_encryption_type` is set to `DiskWithVMGuestState`"
		// TODO - "`secure_boot_enabled` must be set to `true` when `os_disk.0.security_encryption_type` is set to `DiskWithVMGuestState`"
		// TODO - "`vtpm_enabled` must be set to `true` when `os_disk.0.security_encryption_type` is specified"
		// TODO - "an `eviction_policy` can only be specified when `priority` is set to `Spot`"
		// TODO - "an `eviction_policy` must be specified when `priority` is set to `Spot`"
		// TODO - "`max_bid_price` can only be configured when `priority` is set to `Spot`"
		// TODO - "once a customer-managed key is used, you can’t change the selection back to a platform-managed key" - Conditional Requires Replace, so should be PlanModifier on the value?
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("admin_username"),
			path.MatchRoot("os_managed_disk_id"),
		),
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("os_managed_disk_id"),
			path.MatchRoot("source_image_id"),
			path.MatchRoot("source_image_reference"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("os_managed_disk_id"),
			path.MatchRoot("os_disk").AtAnyListIndex().AtName("storage_account_type"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("os_managed_disk_id"),
			path.MatchRoot("computer_name"),
		),
		resourcevalidator.RequiredTogether(
			path.MatchRoot("platform_fault_domain"),
			path.MatchRoot("virtual_machine_scale_set_id"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("os_managed_disk_id"),
			path.MatchRoot("admin_password"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("os_managed_disk_id"),
			path.MatchRoot("admin_ssh_key"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("availability_set_id"),
			path.MatchRoot("zone"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("bypass_platform_safety_checks_on_user_schedule_enabled"),
			path.MatchRoot("os_managed_disk_id"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("os_managed_disk_id"),
			path.MatchRoot("custom_data"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("os_managed_disk_id"),
			path.MatchRoot("reboot_setting"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("os_managed_disk_id"),
			path.MatchRoot("patch_assessment_mode"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("os_managed_disk_id"),
			path.MatchRoot("provision_vm_agent"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("os_managed_disk_id"),
			path.MatchRoot("patch_mode"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("capacity_reservation_group_id"),
			path.MatchRoot("proximity_placement_group_id"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("availability_set_id"),
			path.MatchRoot("capacity_reservation_group_id"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("availability_set_id"),
			path.MatchRoot("virtual_machine_scale_set_id"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("dedicated_host_id"),
			path.MatchRoot("dedicated_host_group_id"),
		),
		// Nested
		resourcevalidator.Conflicting(
			path.MatchRoot("os_disk").AtAnyListIndex().AtName("disk_encryption_set_id"),
			path.MatchRoot("os_disk").AtAnyListIndex().AtName("secure_vm_disk_encryption_set_id"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("os_managed_disk_id"),
			path.MatchRoot("os_disk").AtAnyListIndex().AtName("name"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("os_managed_disk_id"),
			path.MatchRoot("os_disk").AtAnyListIndex().AtName("diff_disk_settings"),
		),
	}
}

func (l *linuxVirtualMachineResource) ModelObject() any {
	return &linuxVirtualMachineResourceModel{}
}

func (l *linuxVirtualMachineResource) ResourceType() string {
	return "azurerm_linux_virtual_machine"
}

func (l *linuxVirtualMachineResource) Schema(ctx context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			fwcommonschema.Name: schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: computeValidate.VirtualMachineName,
					},
				},
			},

			fwcommonschema.ResourceGroupName: fwcommonschema.ResourceGroupNameAttribute(),

			fwcommonschema.Location: fwlocation.LocationAttribute(),

			"admin_username": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: computeValidate.LinuxAdminUsername,
					},
				},
			},

			"network_interface_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						typehelpers.WrappedStringValidator{
							Func: commonids.ValidateNetworkInterfaceID,
						},
					),
				},
			},

			"os_managed_disk_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateManagedDiskID,
					},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},

			"size": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},

			"admin_password": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					ignoreAdminPasswordDiffSuppressModifier{},
					stringplanmodifier.RequiresReplace(),
				},
				Sensitive: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: computeValidate.LinuxAdminPassword,
					},
				},
			},

			"allow_extension_operations": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  typehelpers.NewWrappedBoolDefault(true),
			},

			"availability_set_id": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					casing.IgnoreCaseStringPlanModifier(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateAvailabilitySetID,
					},
				},
			},

			"bypass_platform_safety_checks_on_user_schedule_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  typehelpers.NewWrappedBoolDefault(false),
			},

			"capacity_reservation_group_id": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					casing.IgnoreCaseStringPlanModifier(),
				},
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: capacityreservationgroups.ValidateCapacityReservationGroupID,
					},
				},
			},

			"computer_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: computeValidate.LinuxComputerNameFull,
					},
				},
			},

			"custom_data": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validation.StringIsBase64,
					},
				},
			},

			"dedicated_host_id": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					casing.IgnoreCaseStringPlanModifier(),
				},
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateDedicatedHostID,
					},
				},
			},

			"dedicated_host_group_id": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					casing.IgnoreCaseStringPlanModifier(),
				},
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateDedicatedHostGroupID,
					},
				},
			},

			"disable_password_authentication": schema.BoolAttribute{ // Note - not changing this to match `_enabled` pattern at this time to reduce change on resource
				Optional: true,
				Computed: true,
				Default:  typehelpers.NewWrappedBoolDefault(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},

			"disk_controller_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(virtualmachines.DiskControllerTypesSCSI),
						string(virtualmachines.DiskControllerTypesNVMe),
					),
				},
			},

			"edge_zone": fwcommonschema.EdgeZoneOptionalRequiresReplaceAttribute(),

			"encryption_at_host_enabled": schema.BoolAttribute{
				Optional: true,
			},

			"eviction_policy": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(virtualmachines.VirtualMachineEvictionPolicyTypesDeallocate),
						string(virtualmachines.VirtualMachineEvictionPolicyTypesDelete),
					),
				},
			},

			"extensions_time_budget": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  typehelpers.NewWrappedStringDefault("PT1H30M"),
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: azValidate.ISO8601DurationBetween("PT15M", "PT2H"),
					},
				},
			},

			"license_type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"RHEL_BYOS",
						"RHEL_BASE",
						"RHEL_EUS",
						"RHEL_SAPAPPS",
						"RHEL_SAPHA",
						"RHEL_BASESAPAPPS",
						"RHEL_BASESAPHA",
						"SLES_BYOS",
						"SLES_SAP",
						"SLES_HPC",
						"UBUNTU_PRO",
					),
				},
			},

			"max_bid_price": schema.Float64Attribute{
				Optional: true,
				// Computed: true,
				// Default:  typehelpers.NewWrappedFloat64Default(-1.0),
				Validators: []validator.Float64{
					float64validator.AtLeast(-1.0),
				},
			},

			"priority": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: typehelpers.NewWrappedStringDefault(virtualmachines.VirtualMachinePriorityTypesRegular),
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(virtualmachines.VirtualMachinePriorityTypesRegular),
						string(virtualmachines.VirtualMachinePriorityTypesSpot),
					),
				},
			},

			"provision_vm_agent": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  typehelpers.NewWrappedBoolDefault(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},

			"patch_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  typehelpers.NewWrappedStringDefault(virtualmachines.LinuxVMGuestPatchModeImageDefault),
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(virtualmachines.LinuxVMGuestPatchModeAutomaticByPlatform),
						string(virtualmachines.LinuxVMGuestPatchModeImageDefault),
					),
				},
			},

			"patch_assessment_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(virtualmachines.LinuxPatchAssessmentModeAutomaticByPlatform),
						string(virtualmachines.LinuxPatchAssessmentModeImageDefault),
					),
				},
			},

			"proximity_placement_group_id": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					casing.IgnoreCaseStringPlanModifier(),
				},
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: proximityplacementgroups.ValidateProximityPlacementGroupID,
					},
				},
			},

			"reboot_setting": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(virtualmachines.LinuxVMGuestPatchAutomaticByPlatformRebootSettingAlways),
						string(virtualmachines.LinuxVMGuestPatchAutomaticByPlatformRebootSettingIfRequired),
						string(virtualmachines.LinuxVMGuestPatchAutomaticByPlatformRebootSettingNever),
					),
				},
			},

			"secure_boot_enabled": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},

			"source_image_id": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.Any(
						typehelpers.WrappedStringValidator{Func: images.ValidateImageID},
						typehelpers.WrappedStringValidator{Func: computeValidate.SharedImageID},
						typehelpers.WrappedStringValidator{Func: computeValidate.SharedImageVersionID},
						typehelpers.WrappedStringValidator{Func: computeValidate.CommunityGalleryImageID},
						typehelpers.WrappedStringValidator{Func: computeValidate.CommunityGalleryImageVersionID},
						typehelpers.WrappedStringValidator{Func: computeValidate.SharedGalleryImageID},
						typehelpers.WrappedStringValidator{Func: computeValidate.SharedGalleryImageVersionID},
					),
				},
			},

			"virtual_machine_scale_set_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateVirtualMachineScaleSetID,
					},
				},
			},

			"vtpm_enabled": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},

			"platform_fault_domain": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},

			fwcommonschema.Tags: fwcommonschema.TagsResourceAttribute(ctx),

			"user_data": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validation.StringIsBase64,
					},
				},
			},

			// computed only
			"private_ip_address": schema.StringAttribute{
				Computed: true,
			},

			"private_ip_addresses": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},

			"public_ip_address": schema.StringAttribute{
				Computed: true,
			},

			"public_ip_addresses": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},

			"virtual_machine_id": schema.StringAttribute{
				Computed: true,
			},

			"vm_agent_platform_updates_enabled": schema.BoolAttribute{
				Computed: true,
			},

			"zone": fwcommonschema.ZoneSingleOptionalForceNew(),
		},
		Blocks: map[string]schema.Block{
			"os_disk": schema.ListNestedBlock{
				CustomType: typehelpers.NewListNestedObjectTypeOf[linuxVirtualMachineOSDiskModel](ctx),
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"caching": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									string(virtualmachines.CachingTypesNone),
									string(virtualmachines.CachingTypesReadOnly),
									string(virtualmachines.CachingTypesReadWrite),
								),
							},
						},

						"storage_account_type": schema.StringAttribute{
							Optional: true,
							Computed: true,
							// whilst this appears in the Update block the API returns this when changing:
							// Changing property 'osDisk.managedDisk.storageAccountType' is not allowed
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplaceIfConfigured(),
							},
							Validators: []validator.String{
								stringvalidator.OneOf(
									string(virtualmachines.StorageAccountTypesPremiumLRS),
									string(virtualmachines.StorageAccountTypesStandardLRS),
									string(virtualmachines.StorageAccountTypesStandardSSDLRS),
									string(virtualmachines.StorageAccountTypesStandardSSDZRS),
									string(virtualmachines.StorageAccountTypesPremiumZRS),
								),
							},
						},

						// Optional

						"disk_encryption_set_id": schema.StringAttribute{
							Optional: true,
							// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE
							PlanModifiers: []planmodifier.String{
								casing.IgnoreCaseStringPlanModifier(),
							},
							Validators: []validator.String{
								typehelpers.WrappedStringValidator{
									Func: computeValidate.DiskEncryptionSetID,
								},
							},
						},

						"disk_size_gb": schema.Int64Attribute{
							Optional: true,
							Computed: true,
							Validators: []validator.Int64{
								int64validator.Between(0, 4095),
							},
						},

						"name": schema.StringAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplaceIfConfigured(),
							},
						},

						"secure_vm_disk_encryption_set_id": schema.StringAttribute{
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								typehelpers.WrappedStringValidator{
									Func: computeValidate.DiskEncryptionSetID,
								},
							},
						},

						"security_encryption_type": schema.StringAttribute{
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplaceIfConfigured(),
							},
							Validators: []validator.String{
								stringvalidator.OneOf(
									string(virtualmachines.SecurityEncryptionTypesVMGuestStateOnly),
									string(virtualmachines.SecurityEncryptionTypesDiskWithVMGuestState),
								),
							},
						},

						"write_accelerator_enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  typehelpers.NewWrappedBoolDefault(false),
						},

						"id": schema.StringAttribute{
							Computed: true,
						},
					},
					Blocks: map[string]schema.Block{
						"diff_disk_settings": schema.ListNestedBlock{
							CustomType: typehelpers.NewListNestedObjectTypeOf[linuxVirtualMachineDiffDiskSettings](ctx),
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"option": schema.StringAttribute{
										Required: true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.RequiresReplace(),
										},
										Validators: []validator.String{
											stringvalidator.OneOf(
												string(virtualmachines.DiffDiskOptionsLocal),
											),
										},
									},
									"placement": schema.StringAttribute{
										Optional: true,
										Computed: true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.RequiresReplace(),
										},
										Default: typehelpers.NewWrappedStringDefault(virtualmachines.DiffDiskPlacementCacheDisk),
										Validators: []validator.String{
											stringvalidator.OneOf(
												string(virtualmachines.DiffDiskPlacementCacheDisk),
												string(virtualmachines.DiffDiskPlacementResourceDisk),
												string(virtualmachines.DiffDiskPlacementNVMeDisk),
											),
										},
									},
								},
							},
							PlanModifiers: []planmodifier.List{
								listplanmodifier.RequiresReplace(),
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
					},
				},
			},

			"additional_capabilities": schema.ListNestedBlock{
				CustomType: typehelpers.NewListNestedObjectTypeOf[virtualMachineAdditionalCapabilitiesModel](ctx),
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"ultra_ssd_enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  typehelpers.NewWrappedBoolDefault(false),
						},

						"hibernation_enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  typehelpers.NewWrappedBoolDefault(false),
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},

			"admin_ssh_key": schema.SetNestedBlock{
				CustomType: typehelpers.NewSetNestedObjectTypeOf[virtualMachineSSHKeyModel](ctx),
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"public_key": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								planmodifiers.SSHKeyStringPlanModifier(),
							},
							Validators: []validator.String{
								typehelpers.WrappedStringValidator{
									Func: computeValidate.SSHKey,
								},
							},
						},

						"username": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},

			"boot_diagnostics": schema.ListNestedBlock{
				CustomType: typehelpers.NewListNestedObjectTypeOf[virtualMachineBootDiagnosticsModel](ctx),
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"storage_account_uri": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								typehelpers.WrappedStringValidator{
									Func: validation.IsURLWithHTTPorHTTPS,
								},
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},

			"gallery_application": schema.ListNestedBlock{
				CustomType: typehelpers.NewListNestedObjectTypeOf[virtualMachineGalleryApplicationModel](ctx),
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"version_id": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								typehelpers.WrappedStringValidator{
									Func: galleryapplicationversions.ValidateApplicationVersionID,
								},
							},
						},

						"automatic_upgrade_enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  typehelpers.NewWrappedBoolDefault(false),
						},

						"configuration_blob_uri": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								typehelpers.WrappedStringValidator{
									Func: validation.IsURLWithHTTPorHTTPS,
								},
							},
						},

						"order": schema.Int64Attribute{
							Optional: true,
							Computed: true,
							Default:  typehelpers.NewWrappedInt64Default(0),
							Validators: []validator.Int64{
								int64validator.Between(0, math.MaxInt32),
							},
						},

						"tag": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},

						"treat_failure_as_deployment_failure_enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  typehelpers.NewWrappedBoolDefault(false),
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(100),
				},
			},

			"identity": fwidentity.IdentityResourceBlockSchema(ctx),

			"plan": schema.ListNestedBlock{
				CustomType: typehelpers.NewListNestedObjectTypeOf[virtualMachinePlanModel](ctx),
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},

						"product": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},

						"publisher": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},

			"secret": schema.ListNestedBlock{
				CustomType: typehelpers.NewListNestedObjectTypeOf[virtualMachineSecretModel](ctx),
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"key_vault_id": fwcommonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),
					},

					Blocks: map[string]schema.Block{
						"certificate": schema.SetNestedBlock{
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"url": schema.StringAttribute{
										Required: true,
										Validators: []validator.String{
											typehelpers.WrappedStringValidator{
												Func: keyVaultValidate.NestedItemId,
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"source_image_reference": schema.ListNestedBlock{
				CustomType: typehelpers.NewListNestedObjectTypeOf[virtualMachineSourceImageReference](ctx),
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"publisher": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},

						"offer": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},

						"sku": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},

						"version": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
				},
			},

			"os_image_notification": schema.ListNestedBlock{
				CustomType: typehelpers.NewListNestedObjectTypeOf[virtualMachineOSImageNotificationModel](ctx),
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"timeout": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  typehelpers.NewWrappedStringDefault("PT15M"),
							Validators: []validator.String{
								stringvalidator.OneOf("PT15M"), // TODO - is this right? Why is this even exposed?
							},
						},
					},
				},
			},

			"termination_notification": schema.ListNestedBlock{
				CustomType: typehelpers.NewListNestedObjectTypeOf[virtualMachineTerminationNotificationModel](ctx),
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Required: true,
						},

						"timeout": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  typehelpers.NewWrappedStringDefault("PT15M"),
							Validators: []validator.String{
								typehelpers.WrappedStringValidator{
									Func: azValidate.ISO8601DurationBetween("PT5M", "PT15M"),
								},
							},
						},
					},
				},
			},
		},
	}

	if !features.FivePointOh() {
		s.Attributes["vm_agent_platform_updates_enabled"] = schema.BoolAttribute{
			Optional:           true,
			Computed:           false,
			DeprecationMessage: "this property has been deprecated due to a breaking change introduced by the Service team, which redefined it as a read-only field within the API",
		}
	}

	response.Schema = s
}

func (l *linuxVirtualMachineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, metadata sdk.ResourceMetadata, plan any) {
	client := metadata.Client.Compute.VirtualMachinesClient
	data := sdk.AssertResourceModelType[linuxVirtualMachineResourceModel](plan, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	// We need to check explicit config values, not just the planned result
	config := linuxVirtualMachineResourceModel{}
	metadata.DecodeCreate(ctx, req, resp, &config)
	if resp.Diagnostics.HasError() {
		return
	}

	id := virtualmachines.NewVirtualMachineID(metadata.SubscriptionId, data.ResourceGroupName.ValueString(), data.Name.ValueString())

	existing, err := client.Get(ctx, id, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("checking for existing: %+v", err), err.Error())
			return
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		metadata.ResourceRequiresImport(l.ResourceType(), id, resp)
		return
	}

	osDiskIsImported := config.OSManagedDiskID.ValueString() != ""

	vmp := expandPlanModel(ctx, data.Plan, &resp.Diagnostics)

	params := virtualmachines.VirtualMachine{
		ExtendedLocation: expandEdgeZone(data.EdgeZone.ValueString()),
		Location:         fwlocation.Normalize(data.Location.ValueString()),
		Plan:             vmp,
		Properties: &virtualmachines.VirtualMachineProperties{
			AdditionalCapabilities: expandVirtualMachineAdditionalCapabilitiesModel(ctx, data.AdditionalCapabilities, &resp.Diagnostics),
			ApplicationProfile: &virtualmachines.ApplicationProfile{
				GalleryApplications: expandVirtualMachineGalleryApplicationModel(ctx, data.GalleryApplications, &resp.Diagnostics),
			},
			DiagnosticsProfile: &virtualmachines.DiagnosticsProfile{
				BootDiagnostics: expandBootDiagnosticsModel(ctx, data.BootDiagnostics, &resp.Diagnostics),
			},
			ExtensionsTimeBudget: data.ExtensionsTimeBudget.ValueStringPointer(),
			HardwareProfile: &virtualmachines.HardwareProfile{
				VMSize: pointer.ToEnum[virtualmachines.VirtualMachineSizeTypes](data.Size.ValueString()),
			},
			NetworkProfile: &virtualmachines.NetworkProfile{
				NetworkInterfaces: expandVirtualMachineNetworkInterfaceIDsModel(ctx, data.NetworkInterfaceIDs, &resp.Diagnostics),
			},
			Priority: pointer.ToEnum[virtualmachines.VirtualMachinePriorityTypes](data.Priority.ValueString()),
			SecurityProfile: &virtualmachines.SecurityProfile{
				EncryptionAtHost: data.EncryptionAtHost.ValueBoolPointer(),
			},
			StorageProfile: &virtualmachines.StorageProfile{
				DataDisks:      &[]virtualmachines.DataDisk{},
				ImageReference: expandVirtualMachineSourceImageReference(ctx, data.SourceImageReference, &resp.Diagnostics),
			},
			UserData: data.UserData.ValueStringPointer(),
		},
		Tags: fwcommonschema.ExpandTags(ctx, data.Tags, &resp.Diagnostics),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.PlatformFaultDomain.IsUnknown() {
		params.Properties.PlatformFaultDomain = data.PlatformFaultDomain.ValueInt64Pointer()
	}

	if !data.Identity.IsUnknown() {
		ident := &identity.SystemAndUserAssignedMap{}
		fwidentity.ExpandToSystemAndUserAssignedMap(ctx, data.Identity, ident, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		params.Identity = ident
	}

	osDisk := expandVirtualMachineOSDiskModel(ctx, data.OSDisk, virtualmachines.OperatingSystemTypesLinux, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if !osDiskIsImported {
		osProfile := &virtualmachines.OSProfile{
			AdminUsername:            data.AdminUserName.ValueStringPointer(),
			AdminPassword:            values.ValueStringPointer(data.AdminPassword),
			ComputerName:             values.ValueStringPointer(data.ComputerName),
			CustomData:               values.ValueStringPointer(data.CustomData),
			AllowExtensionOperations: data.AllowExtensionOperations.ValueBoolPointer(),
			LinuxConfiguration: &virtualmachines.LinuxConfiguration{
				DisablePasswordAuthentication: values.ValueBoolPointer(data.DisablePasswordAuthentication),
				PatchSettings: &virtualmachines.LinuxPatchSettings{
					PatchMode: pointer.ToEnum[virtualmachines.LinuxVMGuestPatchMode](data.PatchMode.ValueString()), // TODO - Might need to break this out
				},
				ProvisionVMAgent: data.ProvisionVMAgent.ValueBoolPointer(),
			},
		}

		if v := values.ValueStringPointer(data.PatchAssessmentMode); v != nil {
			osProfile.LinuxConfiguration.PatchSettings.AssessmentMode = pointer.ToEnum[virtualmachines.LinuxPatchAssessmentMode](*v)
			if strings.EqualFold(*v, string(virtualmachines.LinuxVMGuestPatchModeAutomaticByPlatform)) {
				if !data.RebootSetting.IsUnknown() {
					osProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings.RebootSetting = pointer.ToEnum[virtualmachines.LinuxVMGuestPatchAutomaticByPlatformRebootSetting](data.RebootSetting.ValueString())
				}

				if v := values.ValueBoolPointer(data.BypassPlatformSecurityChecks); v != nil && *v {
					osProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings = &virtualmachines.LinuxVMGuestPatchAutomaticByPlatformSettings{
						BypassPlatformSafetyChecksOnUserSchedule: v,
					}
				}
			}
		}

		if data.ComputerName.IsUnknown() {
			_, errs := computeValidate.LinuxComputerNameFull(data.Name.ValueString(), "computer_name")
			if len(errs) > 0 {
				resp.Diagnostics.AddError(fmt.Sprintf("unable to assume default computer name %s", id.VirtualMachineName), fmt.Sprintf("Please adjust the `name`, or specify an explicit `computer_name`: +%v", errs[0]))
				return
			}

			osProfile.ComputerName = data.Name.ValueStringPointer()
			data.ComputerName = data.Name
		}

		if !data.AdminSSHKey.IsUnknown() {
			osProfile.LinuxConfiguration.Ssh = &virtualmachines.SshConfiguration{
				PublicKeys: expandVirtualMachineSSHKeyModel(ctx, data.AdminSSHKey, &resp.Diagnostics),
			}
		}

		if !data.Secret.IsUnknown() {
			osProfile.Secrets = expandLinuxVirtualMachineSecretModel(ctx, data.Secret, &resp.Diagnostics)
			if resp.Diagnostics.HasError() {
				return
			}

		}

		if v := values.ValueStringPointer(data.SourceImageID); v != nil {
			skip := false
			if _, errors := validation.Any(validate.CommunityGalleryImageID, validate.CommunityGalleryImageVersionID)(*v, "source_image_id"); len(errors) == 0 {
				params.Properties.StorageProfile.ImageReference = &virtualmachines.ImageReference{
					CommunityGalleryImageId: v,
				}
				skip = true
			}

			if _, errors := validation.Any(validate.SharedGalleryImageID, validate.SharedGalleryImageVersionID)(*v, "source_image_id"); len(errors) == 0 && !skip {
				params.Properties.StorageProfile.ImageReference = &virtualmachines.ImageReference{
					SharedGalleryImageId: v,
				}
				skip = true
			}

			if !skip {
				params.Properties.StorageProfile.ImageReference = &virtualmachines.ImageReference{
					Id: v,
				}
			}
		}

		params.Properties.OsProfile = osProfile
	} else {
		diskId, err := commonids.ParseManagedDiskID(data.OSManagedDiskID.ValueString())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(resp, "parsing ID", err.Error())
			return
		}

		osDisk.ManagedDisk.Id = pointer.To(diskId.ID())
		osDisk.CreateOption = virtualmachines.DiskCreateOptionTypesAttach
	}

	params.Properties.StorageProfile.OsDisk = osDisk

	if !data.DiskControllerType.IsUnknown() {
		params.Properties.StorageProfile.DiskControllerType = pointer.ToEnum[virtualmachines.DiskControllerTypes](data.DiskControllerType.ValueString())
	}

	if !data.LicenseType.IsNull() {
		params.Properties.LicenseType = data.LicenseType.ValueStringPointer()
	}

	// Note: Safe to discard the diags here as this has been validated above.
	osDiskRaw, _ := typehelpers.DecodeObjectListOfOne[linuxVirtualMachineOSDiskModel](ctx, data.OSDisk)
	if osDiskRaw.SecurityEncryptionType.ValueString() != "" {
		params.Properties.SecurityProfile.SecurityType = pointer.To(virtualmachines.SecurityTypesConfidentialVM)
		params.Properties.SecurityProfile.UefiSettings = &virtualmachines.UefiSettings{
			SecureBootEnabled: data.SecureBootEnabled.ValueBoolPointer(),
			VTpmEnabled:       data.VTPMEnabled.ValueBoolPointer(),
		}
	} else {
		if data.SecureBootEnabled.ValueBool() || data.VTPMEnabled.ValueBool() {
			params.Properties.SecurityProfile.SecurityType = pointer.To(virtualmachines.SecurityTypesTrustedLaunch)
			params.Properties.SecurityProfile.UefiSettings = &virtualmachines.UefiSettings{
				SecureBootEnabled: data.SecureBootEnabled.ValueBoolPointer(),
				VTpmEnabled:       data.VTPMEnabled.ValueBoolPointer(),
			}
		}
	}

	if !data.OSImageNotification.IsNull() || !data.TerminationNotification.IsNull() {
		params.Properties.ScheduledEventsProfile = &virtualmachines.ScheduledEventsProfile{
			OsImageNotificationProfile:   expandVirtualMachineOSImageNotificationModel(ctx, data.OSImageNotification, &resp.Diagnostics),
			TerminateNotificationProfile: expandVirtualMachineTerminationNotificationModel(ctx, data.TerminationNotification, &resp.Diagnostics),
		}
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !data.AvailabilitySetID.IsNull() {
		params.Properties.AvailabilitySet = &virtualmachines.SubResource{
			Id: data.AvailabilitySetID.ValueStringPointer(),
		}
	}

	if !data.CapacityReservationsGroupID.IsNull() {
		params.Properties.CapacityReservation = &virtualmachines.CapacityReservationProfile{
			CapacityReservationGroup: &virtualmachines.SubResource{
				Id: data.CapacityReservationsGroupID.ValueStringPointer(),
			},
		}
	}

	if !data.DedicatedHostID.IsNull() {
		params.Properties.Host = &virtualmachines.SubResource{
			Id: data.DedicatedHostID.ValueStringPointer(),
		}
	}

	if !data.DedicatedHostGroupID.IsNull() {
		params.Properties.HostGroup = &virtualmachines.SubResource{
			Id: data.DedicatedHostGroupID.ValueStringPointer(),
		}
	}

	if !data.EvictionPolicy.IsNull() {
		params.Properties.EvictionPolicy = pointer.ToEnum[virtualmachines.VirtualMachineEvictionPolicyTypes](data.EvictionPolicy.ValueString())
	}

	if v := values.ValueFloat64Pointer(data.MaxBidPrice); v != nil && *v > 0 {
		params.Properties.BillingProfile = &virtualmachines.BillingProfile{
			MaxPrice: v,
		}
	}

	if !data.ProximityPlacementGroupID.IsNull() {
		params.Properties.ProximityPlacementGroup = &virtualmachines.SubResource{
			Id: data.ProximityPlacementGroupID.ValueStringPointer(),
		}
	}

	if !data.VMSSID.IsNull() {
		params.Properties.VirtualMachineScaleSet = &virtualmachines.SubResource{
			Id: data.VMSSID.ValueStringPointer(),
		}
	}

	if !data.Zone.IsNull() {
		params.Zones = &zones.Schema{data.Zone.ValueString()}
	}

	if err = client.CreateOrUpdateThenPoll(ctx, id, params, virtualmachines.DefaultCreateOrUpdateOperationOptions()); err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("creating Linux %s", id), err)
		return
	}

	data.ID = types.StringValue(id.ID())
	sdk.SetIdentityOnCreate(ctx, l, id.ID(), resp)

	vmReadBack, err := client.Get(ctx, id, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("retrieving %s", id), err.Error())
		return
	}

	if vmReadBack.Model == nil {
		resp.Diagnostics.AddError(fmt.Sprintf("retrieving %s", id), "model is nil")
		return
	}

	flattenLinuxVirtualMachineModel(ctx, &id, vmReadBack.Model, data, metadata, &resp.Diagnostics)
}

// func resourceLinuxVirtualMachine() *pluginsdk.Resource {
// 	resource := &pluginsdk.Resource{
// 		Create: resourceLinuxVirtualMachineCreate,
// 		Read:   resourceLinuxVirtualMachineRead,
// 		Update: resourceLinuxVirtualMachineUpdate,
// 		Delete: resourceLinuxVirtualMachineDelete,
// 		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
// 			_, err := commonids.ParseVirtualMachineID(id)
// 			return err
// 		}, importVirtualMachine(virtualmachines.OperatingSystemTypesLinux, "azurerm_linux_virtual_machine")),
//
// 		Timeouts: &pluginsdk.ResourceTimeout{
// 			Create: pluginsdk.DefaultTimeout(45 * time.Minute),
// 			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
// 			Update: pluginsdk.DefaultTimeout(45 * time.Minute),
// 			Delete: pluginsdk.DefaultTimeout(45 * time.Minute),
// 		},
//
// 		Schema: map[string]*pluginsdk.Schema{
// 			"name": {
// 				Type:         pluginsdk.TypeString,
// 				Required:     true,
// 				ForceNew:     true,
// 				ValidateFunc: computeValidate.VirtualMachineName,
// 			},
//
// 			"resource_group_name": commonschema.ResourceGroupName(),
//
// 			"location": commonschema.Location(),
//
// 			"admin_username": {
// 				Type:     pluginsdk.TypeString,
// 				Optional: true,
// 				ExactlyOneOf: []string{
// 					"admin_username",
// 					"os_managed_disk_id",
// 				},
// 				ForceNew:     true,
// 				ValidateFunc: computeValidate.LinuxAdminUsername,
// 			},
//
// 			"network_interface_ids": {
// 				Type:     pluginsdk.TypeList,
// 				Required: true,
// 				MinItems: 1,
// 				Elem: &pluginsdk.Schema{
// 					Type:         pluginsdk.TypeString,
// 					ValidateFunc: commonids.ValidateNetworkInterfaceID,
// 				},
// 			},
//
// 			"os_disk": virtualMachineOSDiskSchema(),
//
// 			"os_managed_disk_id": {
// 				// Note: O+C as this is the same value as `os_disk.0.id` - which gains a value from implicit
// 				// disk creation with a VM when an existing disk is not specified here. This is a top-level property
// 				// to enable schema validation to guard against any values for `OsProfile` being set, as these are
// 				// incompatible with specifying an existing disk. i.e. the OsProfile becomes unmanageable.
// 				Type:         pluginsdk.TypeString,
// 				Optional:     true,
// 				Computed:     true,
// 				ForceNew:     true,
// 				ValidateFunc: commonids.ValidateManagedDiskID,
// 				ExactlyOneOf: []string{
// 					"os_managed_disk_id",
// 					"source_image_id",
// 					"source_image_reference",
// 				},
// 			},
//
// 			"size": {
// 				Type:         pluginsdk.TypeString,
// 				Required:     true,
// 				ValidateFunc: validation.StringIsNotEmpty,
// 			},
//
// 			"additional_capabilities": virtualMachineAdditionalCapabilitiesSchema(),
//
// 			"admin_password": {
// 				Type:             pluginsdk.TypeString,
// 				Optional:         true,
// 				ForceNew:         true,
// 				Sensitive:        true,
// 				DiffSuppressFunc: adminPasswordDiffSuppressFunc,
// 				ValidateFunc:     computeValidate.LinuxAdminPassword,
// 				ConflictsWith: []string{
// 					"os_managed_disk_id",
// 				},
// 			},
//
// 			"admin_ssh_key": SSHKeysSchemaVM(),
//
// 			"allow_extension_operations": {
// 				Type:     pluginsdk.TypeBool,
// 				Optional: true,
// 				Computed: true,
// 			},
//
// 			"availability_set_id": {
// 				Type:         pluginsdk.TypeString,
// 				Optional:     true,
// 				ForceNew:     true,
// 				ValidateFunc: commonids.ValidateAvailabilitySetID,
// 				// the Compute/VM API is broken and returns the Availability Set name in UPPERCASE :shrug:
// 				// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
// 				DiffSuppressFunc: suppress.CaseDifference,
// 				ConflictsWith: []string{
// 					"capacity_reservation_group_id",
// 					"virtual_machine_scale_set_id",
// 					"zone",
// 				},
// 			},
//
// 			"boot_diagnostics": bootDiagnosticsSchema(),
//
// 			"bypass_platform_safety_checks_on_user_schedule_enabled": {
// 				Type:     pluginsdk.TypeBool,
// 				Optional: true,
// 				Default:  false,
// 				ConflictsWith: []string{
// 					"os_managed_disk_id",
// 				},
// 			},
//
// 			"capacity_reservation_group_id": {
// 				Type:     pluginsdk.TypeString,
// 				Optional: true,
// 				// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE
// 				// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
// 				DiffSuppressFunc: suppress.CaseDifference,
// 				ValidateFunc:     capacityreservationgroups.ValidateCapacityReservationGroupID,
// 				ConflictsWith: []string{
// 					"availability_set_id",
// 					"proximity_placement_group_id",
// 				},
// 			},
//
// 			"computer_name": {
// 				Type:     pluginsdk.TypeString,
// 				Optional: true,
//
// 				// Computed since we reuse the VM name if one's not specified
// 				Computed: true,
// 				ForceNew: true,
//
// 				ValidateFunc: computeValidate.LinuxComputerNameFull,
// 				ConflictsWith: []string{
// 					"os_managed_disk_id",
// 				},
// 			},
//
// 			"custom_data": {
// 				Type:         pluginsdk.TypeString,
// 				Optional:     true,
// 				ForceNew:     true,
// 				Sensitive:    true,
// 				ValidateFunc: validation.StringIsBase64,
// 				ConflictsWith: []string{
// 					"os_managed_disk_id",
// 				},
// 			},
//
// 			"dedicated_host_id": {
// 				Type:         pluginsdk.TypeString,
// 				Optional:     true,
// 				ValidateFunc: commonids.ValidateDedicatedHostID,
// 				// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE :shrug:
// 				// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
// 				DiffSuppressFunc: suppress.CaseDifference,
// 				ConflictsWith: []string{
// 					"dedicated_host_group_id",
// 				},
// 			},
//
// 			"dedicated_host_group_id": {
// 				Type:         pluginsdk.TypeString,
// 				Optional:     true,
// 				ValidateFunc: commonids.ValidateDedicatedHostGroupID,
// 				// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE
// 				// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
// 				DiffSuppressFunc: suppress.CaseDifference,
// 				ConflictsWith: []string{
// 					"dedicated_host_id",
// 				},
// 			},
//
// 			"disable_password_authentication": {
// 				// O+C OsProfile vs os_managed_disk_id
// 				Type:     pluginsdk.TypeBool,
// 				Optional: true,
// 				ForceNew: true,
// 				Computed: true,
// 			},
//
// 			"disk_controller_type": {
// 				Type:     pluginsdk.TypeString,
// 				Optional: true,
// 				Computed: true,
// 				ValidateFunc: validation.StringInSlice([]string{
// 					string(virtualmachines.DiskControllerTypesNVMe),
// 					string(virtualmachines.DiskControllerTypesSCSI),
// 				}, false),
// 			},
//
// 			"edge_zone": commonschema.EdgeZoneOptionalForceNew(),
//
// 			"encryption_at_host_enabled": {
// 				Type:     pluginsdk.TypeBool,
// 				Optional: true,
// 			},
//
// 			"eviction_policy": {
// 				// only applicable when `priority` is set to `Spot`
// 				Type:     pluginsdk.TypeString,
// 				Optional: true,
// 				ForceNew: true,
// 				ValidateFunc: validation.StringInSlice([]string{
// 					string(virtualmachines.VirtualMachineEvictionPolicyTypesDeallocate),
// 					string(virtualmachines.VirtualMachineEvictionPolicyTypesDelete),
// 				}, false),
// 			},
//
// 			"extensions_time_budget": {
// 				Type:         pluginsdk.TypeString,
// 				Optional:     true,
// 				Default:      "PT1H30M",
// 				ValidateFunc: azValidate.ISO8601DurationBetween("PT15M", "PT2H"),
// 			},
//
// 			"gallery_application": VirtualMachineGalleryApplicationSchema(),
//
// 			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),
//
// 			"license_type": {
// 				Type:     pluginsdk.TypeString,
// 				Optional: true,
// 				ValidateFunc: validation.StringInSlice([]string{
// 					"RHEL_BYOS",
// 					"RHEL_BASE",
// 					"RHEL_EUS",
// 					"RHEL_SAPAPPS",
// 					"RHEL_SAPHA",
// 					"RHEL_BASESAPAPPS",
// 					"RHEL_BASESAPHA",
// 					"SLES_BYOS",
// 					"SLES_SAP",
// 					"SLES_HPC",
// 					"UBUNTU_PRO",
// 				}, false),
// 			},
//
// 			"max_bid_price": {
// 				Type:         pluginsdk.TypeFloat,
// 				Optional:     true,
// 				Default:      -1,
// 				ValidateFunc: validation.FloatAtLeast(-1.0),
// 			},
//
// 			"plan": planSchema(),
//
// 			"priority": {
// 				Type:     pluginsdk.TypeString,
// 				Optional: true,
// 				ForceNew: true,
// 				Default:  string(virtualmachines.VirtualMachinePriorityTypesRegular),
// 				ValidateFunc: validation.StringInSlice([]string{
// 					string(virtualmachines.VirtualMachinePriorityTypesRegular),
// 					string(virtualmachines.VirtualMachinePriorityTypesSpot),
// 				}, false),
// 			},
//
// 			"provision_vm_agent": {
// 				// O+C due to incompatibility between specifying `os_managed_disk_id` and sending OsProfile in the create/update requests
// 				Type:     pluginsdk.TypeBool,
// 				Optional: true,
// 				Computed: true,
// 				ForceNew: true,
// 				ConflictsWith: []string{
// 					"os_managed_disk_id",
// 				},
// 			},
//
// 			"patch_mode": {
// 				Type:     pluginsdk.TypeString,
// 				Optional: true,
// 				Computed: true,
// 				ValidateFunc: validation.StringInSlice([]string{
// 					string(virtualmachines.LinuxVMGuestPatchModeAutomaticByPlatform),
// 					string(virtualmachines.LinuxVMGuestPatchModeImageDefault),
// 				}, false),
// 				ConflictsWith: []string{
// 					"os_managed_disk_id",
// 				},
// 			},
//
// 			"patch_assessment_mode": {
// 				// O+C due to incompatibility between `os_managed_disk_id` and sending `OsProfile` in create/update
// 				Type:     pluginsdk.TypeString,
// 				Optional: true,
// 				Computed: true,
// 				ValidateFunc: validation.StringInSlice([]string{
// 					string(virtualmachines.LinuxPatchAssessmentModeAutomaticByPlatform),
// 					string(virtualmachines.LinuxPatchAssessmentModeImageDefault),
// 				}, false),
// 				ConflictsWith: []string{
// 					"os_managed_disk_id",
// 				},
// 			},
//
// 			"proximity_placement_group_id": {
// 				Type:         pluginsdk.TypeString,
// 				Optional:     true,
// 				ValidateFunc: proximityplacementgroups.ValidateProximityPlacementGroupID,
// 				// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE :shrug:
// 				// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
// 				DiffSuppressFunc: suppress.CaseDifference,
// 				ConflictsWith: []string{
// 					"capacity_reservation_group_id",
// 				},
// 			},
//
// 			"reboot_setting": {
// 				Type:     pluginsdk.TypeString,
// 				Optional: true,
// 				ValidateFunc: validation.StringInSlice([]string{
// 					string(virtualmachines.LinuxVMGuestPatchAutomaticByPlatformRebootSettingAlways),
// 					string(virtualmachines.LinuxVMGuestPatchAutomaticByPlatformRebootSettingIfRequired),
// 					string(virtualmachines.LinuxVMGuestPatchAutomaticByPlatformRebootSettingNever),
// 				}, false),
// 				ConflictsWith: []string{
// 					"os_managed_disk_id",
// 				},
// 			},
//
// 			"secret": linuxSecretSchema(),
//
// 			"secure_boot_enabled": {
// 				Type:     pluginsdk.TypeBool,
// 				Optional: true,
// 				ForceNew: true,
// 			},
//
// 			"source_image_id": {
// 				Type:     pluginsdk.TypeString,
// 				Optional: true,
// 				ForceNew: true,
// 				ValidateFunc: validation.Any(
// 					images.ValidateImageID,
// 					computeValidate.SharedImageID,
// 					computeValidate.SharedImageVersionID,
// 					computeValidate.CommunityGalleryImageID,
// 					computeValidate.CommunityGalleryImageVersionID,
// 					computeValidate.SharedGalleryImageID,
// 					computeValidate.SharedGalleryImageVersionID,
// 				),
// 				ExactlyOneOf: []string{
// 					"os_managed_disk_id",
// 					"source_image_id",
// 					"source_image_reference",
// 				},
// 			},
//
// 			"source_image_reference": sourceImageReferenceSchemaVM(),
//
// 			"virtual_machine_scale_set_id": {
// 				Type:     pluginsdk.TypeString,
// 				Optional: true,
// 				ConflictsWith: []string{
// 					"availability_set_id",
// 				},
// 				ValidateFunc: commonids.ValidateVirtualMachineScaleSetID,
// 			},
//
// 			"vtpm_enabled": {
// 				Type:     pluginsdk.TypeBool,
// 				Optional: true,
// 				ForceNew: true,
// 			},
//
// 			"platform_fault_domain": {
// 				Type:         pluginsdk.TypeInt,
// 				Optional:     true,
// 				Default:      -1,
// 				ForceNew:     true,
// 				RequiredWith: []string{"virtual_machine_scale_set_id"},
// 				ValidateFunc: validation.IntAtLeast(-1),
// 			},
//
// 			"tags": commonschema.Tags(),
//
// 			"os_image_notification": virtualMachineOsImageNotificationSchema(),
//
// 			"termination_notification": virtualMachineTerminationNotificationSchema(),
//
// 			"user_data": {
// 				Type:         pluginsdk.TypeString,
// 				Optional:     true,
// 				ValidateFunc: validation.StringIsBase64,
// 			},
//
// 			"zone": commonschema.ZoneSingleOptionalForceNew(),
//
// 			// Computed
// 			"private_ip_address": {
// 				Type:     pluginsdk.TypeString,
// 				Computed: true,
// 			},
// 			"private_ip_addresses": {
// 				Type:     pluginsdk.TypeList,
// 				Computed: true,
// 				Elem: &pluginsdk.Schema{
// 					Type: pluginsdk.TypeString,
// 				},
// 			},
// 			"public_ip_address": {
// 				Type:     pluginsdk.TypeString,
// 				Computed: true,
// 			},
// 			"public_ip_addresses": {
// 				Type:     pluginsdk.TypeList,
// 				Computed: true,
// 				Elem: &pluginsdk.Schema{
// 					Type: pluginsdk.TypeString,
// 				},
// 			},
// 			"virtual_machine_id": {
// 				Type:     pluginsdk.TypeString,
// 				Computed: true,
// 			},
// 			"vm_agent_platform_updates_enabled": {
// 				Type:     pluginsdk.TypeBool,
// 				Computed: true,
// 			},
// 		},
// 	}
//
// 	if !features.FivePointOh() {
// 		resource.Schema["vm_agent_platform_updates_enabled"] = &pluginsdk.Schema{
// 			Type:       pluginsdk.TypeBool,
// 			Optional:   true,
// 			Computed:   true,
// 			Deprecated: "this property has been deprecated due to a breaking change introduced by the Service team, which redefined it as a read-only field within the API",
// 		}
// 	}
//
// 	return resource
// }

func (l *linuxVirtualMachineResource) Read(ctx context.Context, _ resource.ReadRequest, resp *resource.ReadResponse, metadata sdk.ResourceMetadata, decodedState any) {
	client := metadata.Client.Compute.VirtualMachinesClient

	state := sdk.AssertResourceModelType[linuxVirtualMachineResourceModel](decodedState, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := virtualmachines.ParseVirtualMachineID(state.ID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "parsing ID", err)
		return
	}

	getOpts := virtualmachines.GetOperationOptions{
		Expand: pointer.To(virtualmachines.InstanceViewTypesUserData),
	}

	existing, err := client.Get(ctx, *id, getOpts)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			metadata.MarkAsGone(id, &resp.State, &resp.Diagnostics)
			return
		}

		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("retrieving %s:", id), err)
		return
	}

	flattenLinuxVirtualMachineModel(ctx, id, existing.Model, state, metadata, &resp.Diagnostics)
}

func (l *linuxVirtualMachineResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse, metadata sdk.ResourceMetadata, decodedState any) {
	client := metadata.Client.Compute.VirtualMachinesClient

	state := sdk.AssertResourceModelType[linuxVirtualMachineResourceModel](decodedState, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := virtualmachines.ParseVirtualMachineID(state.ID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "parsing ID", err)
		return
	}

	// TODO - Why are we locking here?
	locks.ByName(id.VirtualMachineName, VirtualMachineResourceName)
	defer locks.UnlockByName(id.VirtualMachineName, VirtualMachineResourceName)

	getOpts := virtualmachines.GetOperationOptions{
		Expand: pointer.To(virtualmachines.InstanceViewTypesUserData),
	}

	existing, err := client.Get(ctx, *id, getOpts)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return
		}

		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("retrieving %s:", id), err)
		return
	}

	deleteOptions := virtualmachines.DefaultDeleteOperationOptions()
	if metadata.Features.VirtualMachine.SkipShutdownAndForceDelete {
		deleteOptions.ForceDeletion = pointer.To(true)
	}

	if err = client.DeleteThenPoll(ctx, *id, deleteOptions); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("deleting Linux %s", id), err.Error())
		return
	}

	if deleteOSDisk := metadata.Features.VirtualMachine.DeleteOSDiskOnDeletion; deleteOSDisk {
		disksClient := metadata.Client.Compute.DisksClient
		managedDiskId := ""
		if props := existing.Model.Properties; props != nil && props.StorageProfile != nil && props.StorageProfile.OsDisk != nil {
			if disk := props.StorageProfile.OsDisk.ManagedDisk; disk != nil && disk.Id != nil {
				managedDiskId = *disk.Id
			}
		}

		if managedDiskId != "" {
			diskId, err := commonids.ParseManagedDiskID(managedDiskId)
			if err != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("parsing disk ID for %s", id), err.Error())
				return
			}

			if err := disksClient.DeleteThenPoll(ctx, *diskId); err != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("deleting %s for Linux %s", diskId, id), err.Error())
				return

			}
		}
	} else {
		disksClient := metadata.Client.Compute.DisksClient
		managedDiskId := ""
		if props := existing.Model.Properties; props != nil && props.StorageProfile != nil && props.StorageProfile.OsDisk != nil {
			if disk := props.StorageProfile.OsDisk.ManagedDisk; disk != nil && disk.Id != nil {
				managedDiskId = *disk.Id
			}
		}
		if managedDiskId != "" {
			diskId, err := commonids.ParseManagedDiskID(managedDiskId)
			if err != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("parsing disk ID for %s", id), err.Error())
				return
			}

			pollerType := custompoller.NewManagedDiskDetachedPoller(disksClient, *diskId, *id)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				resp.Diagnostics.AddError("waiting for OSdis", "")
				return
			}
		}
	}
}

func (l *linuxVirtualMachineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse, _ sdk.ResourceMetadata) {
	if req.ID == "" {
		resourceIdentity := &linuxVirtualMachineResourceIdentityModel{}
		req.Identity.Get(ctx, resourceIdentity)
		id := pointer.To(virtualmachines.NewVirtualMachineID(resourceIdentity.SubscriptionId, resourceIdentity.ResourceGroupName, resourceIdentity.Name))
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id.ID())...)
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (l *linuxVirtualMachineResource) Identity() (id resourceids.ResourceId, idType sdk.ResourceTypeForIdentity) {
	return &commonids.VirtualMachineId{}, sdk.ResourceTypeForIdentityDefault
}

func (l *linuxVirtualMachineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, metadata sdk.ResourceMetadata, decodedPlan any, decodedState any) {
	client := metadata.Client.Compute.VirtualMachinesClient
	diskClient := metadata.Client.Compute.DisksClient

	plan := sdk.AssertResourceModelType[linuxVirtualMachineResourceModel](decodedPlan, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	state := sdk.AssertResourceModelType[linuxVirtualMachineResourceModel](decodedState, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	var config linuxVirtualMachineResourceModel // Raw config values required
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := virtualmachines.ParseVirtualMachineID(state.ID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "parsing ID", err)
		return
	}

	locks.ByName(id.VirtualMachineName, VirtualMachineResourceName)
	defer locks.UnlockByName(id.VirtualMachineName, VirtualMachineResourceName)

	getOpts := virtualmachines.GetOperationOptions{
		Expand: pointer.To(virtualmachines.InstanceViewTypesUserData),
	}

	existing, err := client.Get(ctx, *id, getOpts)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			metadata.MarkAsGone(id, &resp.State, &resp.Diagnostics)
			return
		}

		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("retrieving %s:", id), err)
		return
	}

	if existing.Model == nil {
		resp.Diagnostics.AddError("reading model", fmt.Sprintf("model for Linux %s was nil", *id))
		return
	}

	model := *existing.Model
	if model.Properties == nil {
		resp.Diagnostics.AddError("reading properties", fmt.Sprintf("properties for Linux %s was nil", *id))
		return
	}

	props := *model.Properties

	instanceView, err := client.InstanceView(ctx, *id)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("retrieving InstanceView for Linux %s", id), err.Error())
		return
	}

	shouldTurnBackOn := virtualMachineShouldBeStarted(instanceView.Model)
	hasEphemeralOSDisk := false
	if storage := props.StorageProfile; storage != nil {
		if disk := storage.OsDisk; disk != nil {
			if settings := disk.DiffDiskSettings; settings != nil && settings.Option != nil {
				hasEphemeralOSDisk = *settings.Option == virtualmachines.DiffDiskOptionsLocal
			}
		}
	}

	shouldUpdate := false
	shouldShutDown := false
	shouldDeallocate := false

	update := virtualmachines.VirtualMachineUpdate{
		Properties: &virtualmachines.VirtualMachineProperties{},
	}

	if sdk.HasChange(plan.BootDiagnostics, state.BootDiagnostics) {
		shouldUpdate = true

		update.Properties.DiagnosticsProfile = &virtualmachines.DiagnosticsProfile{
			BootDiagnostics: expandBootDiagnosticsModel(ctx, plan.BootDiagnostics, &resp.Diagnostics),
		}
	}

	if sdk.HasChange(plan.Secret, state.Secret) {
		shouldUpdate = true

		update.Properties.OsProfile = &virtualmachines.OSProfile{
			Secrets: expandLinuxVirtualMachineSecretModel(ctx, plan.Secret, &resp.Diagnostics),
		}
	}

	if sdk.HasChange(plan.Identity, state.Identity) {
		// Due to the behaviour of identity, we need to send this separately to the API via the CreateUpdate, rather than the patch on change.
		// Additionally, removing a UserAssigned identity, requires sending an explicit null for the property which is not currently supported in go-azure-sdk
		ident := &identity.SystemAndUserAssignedMap{}
		fwidentity.ExpandToSystemAndUserAssignedMap(ctx, plan.Identity, ident, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		model.Identity = ident
		if err = client.CreateOrUpdateThenPoll(ctx, *id, model, virtualmachines.DefaultCreateOrUpdateOperationOptions()); err != nil {
			sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("updating identity for Linux %s", id), err)
			return
		}
	}

	if sdk.HasChange(plan.LicenseType, state.LicenseType) {
		shouldUpdate = true

		update.Properties.LicenseType = plan.LicenseType.ValueStringPointer()
	}

	if sdk.HasChange(plan.CapacityReservationsGroupID, state.CapacityReservationsGroupID) {
		shouldUpdate = true
		shouldDeallocate = true

		if !plan.CapacityReservationsGroupID.IsNull() {
			update.Properties.CapacityReservation = &virtualmachines.CapacityReservationProfile{
				CapacityReservationGroup: &virtualmachines.SubResource{
					Id: plan.CapacityReservationsGroupID.ValueStringPointer(),
				},
			}
		} else {
			update.Properties.CapacityReservation = &virtualmachines.CapacityReservationProfile{
				CapacityReservationGroup: &virtualmachines.SubResource{},
			}
		}
	}

	if sdk.HasChange(plan.DedicatedHostID, state.DedicatedHostID) {
		shouldUpdate = true
		shouldDeallocate = true

		if !plan.DedicatedHostID.IsNull() {
			update.Properties.Host = &virtualmachines.SubResource{
				Id: plan.DedicatedHostID.ValueStringPointer(),
			}
		} else {
			update.Properties.Host = &virtualmachines.SubResource{}
		}
	}

	if sdk.HasChange(plan.DedicatedHostGroupID, state.DedicatedHostGroupID) {
		shouldUpdate = true
		shouldDeallocate = true

		if !plan.DedicatedHostGroupID.IsNull() {
			update.Properties.HostGroup = &virtualmachines.SubResource{
				Id: plan.DedicatedHostGroupID.ValueStringPointer(),
			}
		} else {
			update.Properties.HostGroup = &virtualmachines.SubResource{}
		}
	}

	if sdk.HasChange(plan.ExtensionsTimeBudget, state.ExtensionsTimeBudget) {
		shouldUpdate = true
		update.Properties.ExtensionsTimeBudget = plan.ExtensionsTimeBudget.ValueStringPointer()
	}

	if sdk.HasChange(plan.ExtensionsTimeBudget, state.ExtensionsTimeBudget) {
		shouldUpdate = true
		update.Properties.ExtensionsTimeBudget = plan.ExtensionsTimeBudget.ValueStringPointer()
	}

	if sdk.HasChange(plan.GalleryApplications, state.GalleryApplications) {
		shouldUpdate = true
		update.Properties.ApplicationProfile = &virtualmachines.ApplicationProfile{
			GalleryApplications: expandVirtualMachineGalleryApplicationModel(ctx, plan.GalleryApplications, &resp.Diagnostics),
		}
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if sdk.HasChange(plan.MaxBidPrice, state.MaxBidPrice) {
		shouldUpdate = true
		shouldShutDown = true
		shouldDeallocate = true

		update.Properties.BillingProfile = &virtualmachines.BillingProfile{
			MaxPrice: plan.MaxBidPrice.ValueFloat64Pointer(),
		}
	}

	if sdk.HasChange(plan.NetworkInterfaceIDs, state.NetworkInterfaceIDs) {
		shouldUpdate = true
		shouldShutDown = true
		shouldDeallocate = true

		update.Properties.NetworkProfile = &virtualmachines.NetworkProfile{
			NetworkInterfaces: expandVirtualMachineNetworkInterfaceIDsModel(ctx, plan.NetworkInterfaceIDs, &resp.Diagnostics),
		}
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if sdk.HasChange(plan.DiskControllerType, state.DiskControllerType) {
		shouldUpdate = true
		shouldDeallocate = true

		if update.Properties.StorageProfile == nil {
			update.Properties.StorageProfile = &virtualmachines.StorageProfile{}
		}

		update.Properties.StorageProfile.DiskControllerType = pointer.ToEnum[virtualmachines.DiskControllerTypes](plan.DiskControllerType.ValueString())
	}

	if sdk.HasChange(plan.OSDisk, state.OSDisk) {
		shouldUpdate = true
		shouldShutDown = true
		shouldDeallocate = true

		if update.Properties.StorageProfile == nil {
			update.Properties.StorageProfile = &virtualmachines.StorageProfile{}
		}

		update.Properties.StorageProfile.OsDisk = expandVirtualMachineOSDiskModel(ctx, plan.OSDisk, virtualmachines.OperatingSystemTypesLinux, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		if !config.OSManagedDiskID.IsNull() { // If an imported disk was used to create the VM, this needs to be set to match or the API will reject it
			update.Properties.StorageProfile.OsDisk.CreateOption = virtualmachines.DiskCreateOptionTypesAttach
		}
	}

	if sdk.HasChange(plan.VMSSID, state.VMSSID) {
		shouldUpdate = true

		if !plan.VMSSID.IsNull() {
			update.Properties.VirtualMachineScaleSet = &virtualmachines.SubResource{
				Id: plan.VMSSID.ValueStringPointer(),
			}
		} else {
			update.Properties.VirtualMachineScaleSet = &virtualmachines.SubResource{}
		}
	}

	if sdk.HasChange(plan.ProximityPlacementGroupID, state.ProximityPlacementGroupID) {
		shouldUpdate = true
		shouldShutDown = true
		shouldDeallocate = true

		if !plan.ProximityPlacementGroupID.IsNull() {
			update.Properties.ProximityPlacementGroup = &virtualmachines.SubResource{
				Id: plan.ProximityPlacementGroupID.ValueStringPointer(),
			}
		} else {
			update.Properties.ProximityPlacementGroup = &virtualmachines.SubResource{}
		}
	}

	if sdk.HasChange(plan.Size, state.Size) {
		shouldUpdate = true
		shouldShutDown = true

		if !vmSizeAvailableOnHost(ctx, id, plan.Size.ValueString(), metadata, &resp.Diagnostics) {
			shouldDeallocate = true
		}

		if resp.Diagnostics.HasError() {
			return
		}

		update.Properties.HardwareProfile = &virtualmachines.HardwareProfile{
			VMSize: pointer.ToEnum[virtualmachines.VirtualMachineSizeTypes](plan.Size.ValueString()),
		}
	}

	if sdk.HasChange(plan.PatchMode, state.PatchMode) {
		shouldUpdate = true

		if update.Properties.OsProfile == nil {
			update.Properties.OsProfile = &virtualmachines.OSProfile{}
		}

		if update.Properties.OsProfile.LinuxConfiguration == nil {
			update.Properties.OsProfile.LinuxConfiguration = &virtualmachines.LinuxConfiguration{}
		}

		if !plan.PatchMode.IsNull() {
			update.Properties.OsProfile.LinuxConfiguration.PatchSettings = &virtualmachines.LinuxPatchSettings{
				PatchMode: pointer.ToEnum[virtualmachines.LinuxVMGuestPatchMode](plan.PatchMode.ValueString()),
			}
		} else {
			update.Properties.OsProfile.LinuxConfiguration.PatchSettings = &virtualmachines.LinuxPatchSettings{
				PatchMode: pointer.To(virtualmachines.LinuxVMGuestPatchModeImageDefault),
			}
		}
	}

	if sdk.HasChange(plan.PatchAssessmentMode, state.PatchAssessmentMode) {
		shouldUpdate = true

		if update.Properties.OsProfile == nil {
			update.Properties.OsProfile = &virtualmachines.OSProfile{}
		}

		if update.Properties.OsProfile.LinuxConfiguration == nil {
			update.Properties.OsProfile.LinuxConfiguration = &virtualmachines.LinuxConfiguration{}
		}

		update.Properties.OsProfile.LinuxConfiguration.PatchSettings = &virtualmachines.LinuxPatchSettings{
			AssessmentMode: pointer.ToEnum[virtualmachines.LinuxPatchAssessmentMode](plan.PatchAssessmentMode.ValueString()),
		}
	}

	if sdk.HasChange(plan.BypassPlatformSecurityChecks, state.BypassPlatformSecurityChecks) {
		shouldUpdate = true

		if update.Properties.OsProfile == nil {
			update.Properties.OsProfile = &virtualmachines.OSProfile{}
		}

		if update.Properties.OsProfile.LinuxConfiguration == nil {
			update.Properties.OsProfile.LinuxConfiguration = &virtualmachines.LinuxConfiguration{}
		}

		if update.Properties.OsProfile.LinuxConfiguration.PatchSettings == nil {
			update.Properties.OsProfile.LinuxConfiguration.PatchSettings = &virtualmachines.LinuxPatchSettings{}
		}

		if plan.PatchMode.ValueString() == string(virtualmachines.LinuxPatchAssessmentModeAutomaticByPlatform) {
			if update.Properties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings == nil {
				update.Properties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings = &virtualmachines.LinuxVMGuestPatchAutomaticByPlatformSettings{}
			}
			update.Properties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings.BypassPlatformSafetyChecksOnUserSchedule = plan.BypassPlatformSecurityChecks.ValueBoolPointer()
		}
	}

	if sdk.HasChange(plan.RebootSetting, state.RebootSetting) {
		shouldUpdate = true

		if update.Properties.OsProfile == nil {
			update.Properties.OsProfile = &virtualmachines.OSProfile{}
		}

		if update.Properties.OsProfile.LinuxConfiguration == nil {
			update.Properties.OsProfile.LinuxConfiguration = &virtualmachines.LinuxConfiguration{}
		}

		if update.Properties.OsProfile.LinuxConfiguration.PatchSettings == nil {
			update.Properties.OsProfile.LinuxConfiguration.PatchSettings = &virtualmachines.LinuxPatchSettings{}
		}

		if update.Properties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings == nil {
			update.Properties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings = &virtualmachines.LinuxVMGuestPatchAutomaticByPlatformSettings{}
		}

		update.Properties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings.RebootSetting = pointer.ToEnum[virtualmachines.LinuxVMGuestPatchAutomaticByPlatformRebootSetting](plan.RebootSetting.ValueString())
	}

	if sdk.HasChange(plan.AllowExtensionOperations, state.AllowExtensionOperations) {
		shouldUpdate = true
		if update.Properties.OsProfile == nil {
			update.Properties.OsProfile = &virtualmachines.OSProfile{}
		}

		update.Properties.OsProfile.AllowExtensionOperations = plan.AllowExtensionOperations.ValueBoolPointer()
	}

	if sdk.HasChange(plan.OSImageNotification, state.OSImageNotification) || sdk.HasChange(plan.TerminationNotification, state.TerminationNotification) {
		shouldUpdate = true
		update.Properties.ScheduledEventsProfile = &virtualmachines.ScheduledEventsProfile{}

		if sdk.HasChange(plan.OSImageNotification, state.OSImageNotification) {
			update.Properties.ScheduledEventsProfile.OsImageNotificationProfile = expandVirtualMachineOSImageNotificationModel(ctx, plan.OSImageNotification, &resp.Diagnostics)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		if sdk.HasChange(plan.TerminationNotification, state.TerminationNotification) {
			update.Properties.ScheduledEventsProfile.TerminateNotificationProfile = expandVirtualMachineTerminationNotificationModel(ctx, plan.TerminationNotification, &resp.Diagnostics)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}

	if sdk.HasChange(plan.Tags, state.Tags) {
		shouldUpdate = true

		update.Tags = fwcommonschema.ExpandTags(ctx, plan.Tags, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if sdk.HasChange(plan.AdditionalCapabilities, state.AdditionalCapabilities) {
		shouldUpdate = true
		additionalProps := expandVirtualMachineAdditionalCapabilitiesModel(ctx, plan.AdditionalCapabilities, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		// TODO - There must be an easier way of doing this?
		pac, _ := typehelpers.DecodeObjectListOfOne(ctx, plan.AdditionalCapabilities)
		sac, _ := typehelpers.DecodeObjectListOfOne(ctx, state.AdditionalCapabilities)
		if sdk.HasChange(pac.UltraSSDEnabled, sac.UltraSSDEnabled) {
			shouldShutDown = true
			shouldDeallocate = true
		}

		update.Properties.AdditionalCapabilities = additionalProps
	}

	if sdk.HasChange(plan.EncryptionAtHost, state.EncryptionAtHost) {
		update.Properties.SecurityProfile = &virtualmachines.SecurityProfile{
			EncryptionAtHost: plan.EncryptionAtHost.ValueBoolPointer(),
		}
	}

	if sdk.HasChange(plan.UserData, state.UserData) {
		shouldUpdate = true
		update.Properties.UserData = plan.UserData.ValueStringPointer()
	}

	if instanceView.Model != nil && instanceView.Model.Statuses != nil {
		for _, status := range *instanceView.Model.Statuses {
			if status.Code == nil {
				continue
			}

			// could also be the provisioning state which we're not bothered with here
			powerState := strings.ToLower(*status.Code)
			if !strings.HasPrefix(powerState, "powerstate/") {
				continue
			}

			powerState = strings.TrimPrefix(powerState, "powerstate/")
			switch strings.ToLower(powerState) {
			case "deallocated":
				// VM already deallocated, no shutdown or deallocation needed
				shouldShutDown = false
				shouldDeallocate = false
			case "deallocating":
				// VM is deallocating
				// To make sure we do not start updating before this action has finished,
				// only skip the shutdown and send another deallocation request if shouldDeallocate == true
				shouldShutDown = false
			case "stopped":
				// VM already stopped, no shutdown needed
				shouldShutDown = false
			}
		}
	}

	if shouldShutDown {
		if err = client.PowerOffThenPoll(ctx, *id, virtualmachines.DefaultPowerOffOperationOptions()); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("sending Power Off to Linux %s", id), err.Error())
		}
	}

	if shouldDeallocate {
		if hasEphemeralOSDisk {
			resp.Diagnostics.AddWarning("skipping deallocate", fmt.Sprintf("[DEBUG] Skipping deallocation for Linux %s since cannot deallocate a Virtual Machine with an Ephemeral OS Disk", *id))
		} else {
			if err = client.DeallocateThenPoll(ctx, *id, virtualmachines.DefaultDeallocateOperationOptions()); err != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("deallocating Linux %s", id), err.Error())
			}
		}

	}

	po, _ := typehelpers.DecodeObjectListOfOne(ctx, plan.OSDisk)
	so, _ := typehelpers.DecodeObjectListOfOne(ctx, state.OSDisk)
	mID := commonids.NewManagedDiskID(id.SubscriptionId, id.ResourceGroupName, so.Name.ValueString())
	if sdk.HasChange(po.DiskSizeGB, so.DiskSizeGB) {
		diskUpdate := disks.DiskUpdate{
			Properties: &disks.DiskUpdateProperties{
				DiskSizeGB: po.DiskSizeGB.ValueInt64Pointer(),
			},
		}

		if err = diskClient.UpdateThenPoll(ctx, mID, diskUpdate); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("resizing OS Disk %q for Linux %s", mID.DiskName, id), err.Error())
			return
		}
	}

	if sdk.HasChange(po.DiskEncryptionSetID, so.DiskEncryptionSetID) {
		encryptionType, err := retrieveDiskEncryptionSetEncryptionType(ctx, metadata.Client.Compute.DiskEncryptionSetsClient, po.DiskEncryptionSetID.ValueString()) // TODO - Make this native to FW
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("retrieving encryption type for Linux %s", id), err.Error())
			return
		}

		diskUpdate := disks.DiskUpdate{
			Properties: &disks.DiskUpdateProperties{
				Encryption: &disks.Encryption{
					Type:                encryptionType,
					DiskEncryptionSetId: po.DiskEncryptionSetID.ValueStringPointer(),
				},
			},
		}

		if err = diskClient.UpdateThenPoll(ctx, mID, diskUpdate); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("updating encryption settings of OS Disk %q for Linux %s", mID.DiskName, id), err.Error())
			return
		}
	}

	if shouldUpdate {
		if err = client.UpdateThenPoll(ctx, *id, update, virtualmachines.DefaultUpdateOperationOptions()); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("updating Linux %s", id), err.Error())
			return
		}
	}

	if shouldTurnBackOn && (shouldShutDown || shouldDeallocate) {
		if err := client.StartThenPoll(ctx, *id); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("starting Linux %s", id), err.Error())
			return
		}

	}

	vmReadBack, err := client.Get(ctx, *id, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("retrieving %s", id), err.Error())
		return
	}

	if vmReadBack.Model == nil {
		resp.Diagnostics.AddError(fmt.Sprintf("retrieving %s", id), "model is nil")
		return
	}

	flattenLinuxVirtualMachineModel(ctx, id, vmReadBack.Model, plan, metadata, &resp.Diagnostics)

}

// 	return resourceLinuxVirtualMachineRead(d, meta)

// vmSizeAvailableOnHost queries the VM to see if the size requested is allowed. If false, the VM needs to be de-allocated before sending the update.
func vmSizeAvailableOnHost(ctx context.Context, id *virtualmachines.VirtualMachineId, requestedSize string, metadata sdk.ResourceMetadata, diags *diag.Diagnostics) bool {
	client := metadata.Client.Compute.VirtualMachinesClient
	sizes, err := client.ListAvailableSizes(ctx, *id)
	if err != nil {
		diags.AddError(fmt.Sprintf("retrieving available sizes for Linux %s", id), err.Error())
		return false
	}

	if sizes.Model != nil && sizes.Model.Value != nil {
		for _, size := range *sizes.Model.Value {
			if size.Name == nil {
				continue
			}

			if strings.EqualFold(*size.Name, requestedSize) {
				return true
			}
		}
	}

	return false
}
