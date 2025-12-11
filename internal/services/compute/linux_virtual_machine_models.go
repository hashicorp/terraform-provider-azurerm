package compute

import (
	"context"
	"strings"

	fwcommonschema "github.com/hashicorp/go-azure-helpers/framework/commonschema"
	"github.com/hashicorp/go-azure-helpers/framework/convert"
	"github.com/hashicorp/go-azure-helpers/framework/identity"
	fwidentity "github.com/hashicorp/go-azure-helpers/framework/identity"
	fwlocation "github.com/hashicorp/go-azure-helpers/framework/location"
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/framework/values"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type linuxVirtualMachineResourceIdentityModel struct {
	SubscriptionId    string `tfsdk:"subscription_id"`
	ResourceGroupName string `tfsdk:"resource_group_name"`
	Name              string `tfsdk:"name"`
}

type linuxVirtualMachineResourceModel struct {
	ID       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`

	Name                          types.String                         `tfsdk:"name"`
	ResourceGroupName             types.String                         `tfsdk:"resource_group_name"`
	Location                      types.String                         `tfsdk:"location"`
	AdminUserName                 types.String                         `tfsdk:"admin_username"`
	NetworkInterfaceIDs           types.List                           `tfsdk:"network_interface_ids"`
	OSManagedDiskID               types.String                         `tfsdk:"os_managed_disk_id"`
	Size                          types.String                         `tfsdk:"size"`
	AdminPassword                 types.String                         `tfsdk:"admin_password"`
	AllowExtensionOperations      types.Bool                           `tfsdk:"allow_extension_operations"`
	AvailabilitySetID             types.String                         `tfsdk:"availability_set_id"`
	BypassPlatformSecurityChecks  types.Bool                           `tfsdk:"bypass_platform_safety_checks_on_user_schedule_enabled"`
	CapacityReservationsGroupID   types.String                         `tfsdk:"capacity_reservation_group_id"`
	ComputerName                  types.String                         `tfsdk:"computer_name"`
	CustomData                    types.String                         `tfsdk:"custom_data"`
	DedicatedHostID               types.String                         `tfsdk:"dedicated_host_id"`
	DedicatedHostGroupID          types.String                         `tfsdk:"dedicated_host_group_id"`
	DisablePasswordAuthentication types.Bool                           `tfsdk:"disable_password_authentication"`
	DiskControllerType            types.String                         `tfsdk:"disk_controller_type"`
	EdgeZone                      types.String                         `tfsdk:"edge_zone"`
	EncryptionAtHost              types.Bool                           `tfsdk:"encryption_at_host_enabled"`
	EvictionPolicy                types.String                         `tfsdk:"eviction_policy"`
	ExtensionsTimeBudget          types.String                         `tfsdk:"extensions_time_budget"`
	LicenseType                   types.String                         `tfsdk:"license_type"`
	MaxBidPrice                   types.Float64                        `tfsdk:"max_bid_price"`
	Priority                      types.String                         `tfsdk:"priority"`
	ProvisionVMAgent              types.Bool                           `tfsdk:"provision_vm_agent"`
	PatchMode                     types.String                         `tfsdk:"patch_mode"`
	PatchAssessmentMode           types.String                         `tfsdk:"patch_assessment_mode"`
	ProximityPlacementGroupID     types.String                         `tfsdk:"proximity_placement_group_id"`
	RebootSetting                 types.String                         `tfsdk:"reboot_setting"`
	SecureBootEnabled             types.Bool                           `tfsdk:"secure_boot_enabled"`
	SourceImageID                 types.String                         `tfsdk:"source_image_id"`
	VMSSID                        types.String                         `tfsdk:"virtual_machine_scale_set_id"`
	VTPMEnabled                   types.Bool                           `tfsdk:"vtpm_enabled"`
	PlatformFaultDomain           types.Int64                          `tfsdk:"platform_fault_domain"`
	Tags                          typehelpers.MapValueOf[types.String] `tfsdk:"tags"`
	UserData                      types.String                         `tfsdk:"user_data"`
	Zone                          types.String                         `tfsdk:"zone"`

	PrivateIPAddress              types.String `tfsdk:"private_ip_address"`
	PrivateIPAddresses            types.List   `tfsdk:"private_ip_addresses"`
	PublicIPAddress               types.String `tfsdk:"public_ip_address"`
	PublicIPAddresses             types.List   `tfsdk:"public_ip_addresses"`
	VirtualMachineID              types.String `tfsdk:"virtual_machine_id"`
	VMAgentPlatformUpdatesEnabled types.Bool   `tfsdk:"vm_agent_platform_updates_enabled"`

	OSDisk                  typehelpers.ListNestedObjectValueOf[linuxVirtualMachineOSDiskModel]             `tfsdk:"os_disk"`
	AdditionalCapabilities  typehelpers.ListNestedObjectValueOf[virtualMachineAdditionalCapabilitiesModel]  `tfsdk:"additional_capabilities"`
	AdminSSHKey             typehelpers.SetNestedObjectValueOf[virtualMachineSSHKeyModel]                   `tfsdk:"admin_ssh_key"`
	BootDiagnostics         typehelpers.ListNestedObjectValueOf[virtualMachineBootDiagnosticsModel]         `tfsdk:"boot_diagnostics"`
	GalleryApplications     typehelpers.ListNestedObjectValueOf[virtualMachineGalleryApplicationModel]      `tfsdk:"gallery_application"`
	Identity                typehelpers.ListNestedObjectValueOf[identity.IdentityModel]                     `tfsdk:"identity"`
	Plan                    typehelpers.ListNestedObjectValueOf[virtualMachinePlanModel]                    `tfsdk:"plan"`
	Secret                  typehelpers.ListNestedObjectValueOf[virtualMachineSecretModel]                  `tfsdk:"secret"`
	SourceImageReference    typehelpers.ListNestedObjectValueOf[virtualMachineSourceImageReference]         `tfsdk:"source_image_reference"`
	OSImageNotification     typehelpers.ListNestedObjectValueOf[virtualMachineOSImageNotificationModel]     `tfsdk:"os_image_notification"`
	TerminationNotification typehelpers.ListNestedObjectValueOf[virtualMachineTerminationNotificationModel] `tfsdk:"termination_notification"`
}

func flattenLinuxVirtualMachineModel(ctx context.Context, id *virtualmachines.VirtualMachineId, model *virtualmachines.VirtualMachine, state *linuxVirtualMachineResourceModel, metadata sdk.ResourceMetadata, diags *diag.Diagnostics) {
	networkInterfacesClient := metadata.Client.Network.NetworkInterfacesClient
	publicIPAddressesClient := metadata.Client.Network.PublicIPAddresses
	state.Name = types.StringValue(id.VirtualMachineName)
	state.ResourceGroupName = types.StringValue(id.ResourceGroupName)

	// Set schema defaults that will be absent during import

	if model != nil {
		state.Location = fwlocation.NormalizeValue(model.Location)
		state.EdgeZone = types.StringValue(flattenEdgeZone(model.ExtendedLocation))
		if z := model.Zones; z != nil && len(*z) > 0 {
			state.Zone = types.StringValue((*z)[0])
		}

		fwidentity.FlattenFromSystemAndUserAssignedMap(ctx, model.Identity, &state.Identity, diags)
		if diags.HasError() {
			return
		}

		plan := flattenPlanModel(ctx, model.Plan)
		if diags.HasError() {
			return
		}

		state.Plan = plan

		if props := model.Properties; props != nil {
			state.AdditionalCapabilities = flattenVirtualMachineAdditionalCapabilitiesModel(ctx, props.AdditionalCapabilities)
			if props.CapacityReservation != nil && props.CapacityReservation.CapacityReservationGroup != nil {
				state.CapacityReservationsGroupID = types.StringPointerValue(props.CapacityReservation.CapacityReservationGroup.Id)
			}

			if props.ApplicationProfile != nil {
				state.GalleryApplications = flattenVirtualMachineGalleryApplicationModel(ctx, props.ApplicationProfile.GalleryApplications)
			}

			if props.LicenseType != nil && !strings.EqualFold(*props.LicenseType, "None") {
				state.LicenseType = types.StringPointerValue(props.LicenseType)
			}

			state.BootDiagnostics = flattenBootDiagnosticsModel(ctx, props.DiagnosticsProfile)
			if props.EvictionPolicy != nil {
				state.EvictionPolicy = types.StringValue(pointer.FromEnum(props.EvictionPolicy))
			}
			if props.HardwareProfile != nil {
				state.Size = types.StringValue(pointer.FromEnum(props.HardwareProfile.VMSize))
			}

			state.ExtensionsTimeBudget = types.StringPointerValue(props.ExtensionsTimeBudget)
			if props.BillingProfile != nil {
				state.MaxBidPrice = types.Float64PointerValue(props.BillingProfile.MaxPrice)
			}

			if profile := props.NetworkProfile; profile != nil {
				if profile.NetworkInterfaces != nil && len(*profile.NetworkInterfaces) > 0 {
					nids := make([]types.String, 0)
					for _, n := range *profile.NetworkInterfaces {
						nids = append(nids, types.StringPointerValue(n.Id))
					}

					n, d := types.ListValueFrom(ctx, types.StringType, nids)
					if d.HasError() {
						diags.Append(d...)
						return
					}
					state.NetworkInterfaceIDs = n
					if diags.HasError() {
						return
					}
				}
			}

			if props.Host != nil {
				state.DedicatedHostID = types.StringPointerValue(props.Host.Id)
			}

			if props.HostGroup != nil {
				state.DedicatedHostGroupID = types.StringPointerValue(props.HostGroup.Id)
			}

			if props.VirtualMachineScaleSet != nil {
				state.VMSSID = types.StringPointerValue(props.VirtualMachineScaleSet.Id)
			}

			state.PlatformFaultDomain = types.Int64PointerValue(props.PlatformFaultDomain)

			if profile := props.OsProfile; profile != nil {
				state.AdminUserName = types.StringPointerValue(profile.AdminUsername)
				state.AllowExtensionOperations = types.BoolPointerValue(profile.AllowExtensionOperations)
				state.ComputerName = types.StringPointerValue(profile.ComputerName)

				if conf := profile.LinuxConfiguration; conf != nil {
					state.DisablePasswordAuthentication = types.BoolPointerValue(conf.DisablePasswordAuthentication)
					state.ProvisionVMAgent = types.BoolPointerValue(conf.ProvisionVMAgent)
					state.VMAgentPlatformUpdatesEnabled = types.BoolPointerValue(conf.EnableVMAgentPlatformUpdates)
					state.AdminSSHKey = flattenSSHKeyModel(ctx, conf.Ssh)
					if patch := conf.PatchSettings; patch != nil {
						// TODO - Do we need to set a fall-back value here if the "default" returns as nil?
						state.PatchMode = types.StringValue(pointer.FromEnum(patch.PatchMode))
						state.PatchAssessmentMode = types.StringValue(pointer.FromEnum(patch.AssessmentMode))
						state.BypassPlatformSecurityChecks = types.BoolValue(false) // Default from schema not populated by import
						if patch.AutomaticByPlatformSettings != nil {
							state.BypassPlatformSecurityChecks = types.BoolValue(pointer.From(patch.AutomaticByPlatformSettings.BypassPlatformSafetyChecksOnUserSchedule) == true)
							state.RebootSetting = types.StringValue(pointer.FromEnum(patch.AutomaticByPlatformSettings.RebootSetting))
						}
					}
				}

				state.Secret = flattenLinuxVirtualMachineSecretModel(ctx, profile.Secrets)
			}
			state.Priority = types.StringValue(pointer.FromEnum(props.Priority)) // TODO - Do we still need the "" == "Regular" ?
			if p := props.ProximityPlacementGroup; p != nil {
				state.ProximityPlacementGroupID = types.StringPointerValue(p.Id)
			}

			if profile := props.StorageProfile; profile != nil {
				state.DiskControllerType = types.StringValue(pointer.FromEnum(profile.DiskControllerType))
				state.OSDisk, state.OSManagedDiskID = flattenVirtualMachineOSDiskModel(ctx, profile.OsDisk, diags)
				if diags.HasError() {
					return
				}

				if v := profile.ImageReference; v != nil {
					switch {
					case v.Id != nil:
						state.SourceImageID = types.StringPointerValue(v.Id)
					case v.CommunityGalleryImageId != nil:
						state.SourceImageID = types.StringPointerValue(v.CommunityGalleryImageId)
					case v.SharedGalleryImageId != nil:
						state.SourceImageID = types.StringPointerValue(v.SharedGalleryImageId)
					default:
						state.SourceImageReference = flattenVirtualMachineSourceImageReferenceModel(ctx, v)
					}
				}
			}

			if s := props.ScheduledEventsProfile; s != nil {
				state.OSImageNotification = flattenVirtualMachineOSImageNotificationModel(ctx, s.OsImageNotificationProfile)
				state.TerminationNotification = flattenVirtualMachineTerminationNotificationModel(ctx, s.TerminateNotificationProfile)
			}

			if v := props.SecurityProfile; v != nil {
				state.EncryptionAtHost = types.BoolPointerValue(v.EncryptionAtHost)
				if u := v.UefiSettings; u != nil {
					state.VTPMEnabled = types.BoolPointerValue(u.VTpmEnabled)
					state.SecureBootEnabled = types.BoolPointerValue(u.SecureBootEnabled)
				}
			}

			state.VirtualMachineID = types.StringPointerValue(props.VMId)
			state.UserData = types.StringPointerValue(props.UserData)

			connInfo := retrieveConnectionInformation(ctx, networkInterfacesClient, publicIPAddressesClient, props)

			state.PrivateIPAddress = types.StringValue(connInfo.primaryPrivateAddress)
			state.PublicIPAddress = types.StringValue(connInfo.primaryPublicAddress)
			convert.Flatten(ctx, connInfo.publicAddresses, &state.PublicIPAddresses, diags)
			if diags.HasError() {
				return
			}
			convert.Flatten(ctx, connInfo.privateAddresses, &state.PrivateIPAddresses, diags)
			if diags.HasError() {
				return
			}
		}

		state.Tags = fwcommonschema.FlattenTags(ctx, model.Tags, diags)
	}
}

type linuxVirtualMachineOSDiskModel struct {
	Caching                     types.String `tfsdk:"caching"`
	StorageAccountType          types.String `tfsdk:"storage_account_type"`
	DiskEncryptionSetID         types.String `tfsdk:"disk_encryption_set_id"`
	DiskSizeGB                  types.Int64  `tfsdk:"disk_size_gb"`
	Name                        types.String `tfsdk:"name"`
	SecureVMDiskEncryptionSetID types.String `tfsdk:"secure_vm_disk_encryption_set_id"`
	SecurityEncryptionType      types.String `tfsdk:"security_encryption_type"`
	WriteAcceleratorEnabled     types.Bool   `tfsdk:"write_accelerator_enabled"`
	ID                          types.String `tfsdk:"id"`

	DiffDiskSettings typehelpers.ListNestedObjectValueOf[linuxVirtualMachineDiffDiskSettings] `tfsdk:"diff_disk_settings"`
}

func expandVirtualMachineOSDiskModel(ctx context.Context, input typehelpers.ListNestedObjectValueOf[linuxVirtualMachineOSDiskModel], osType virtualmachines.OperatingSystemTypes, diags *diag.Diagnostics) *virtualmachines.OSDisk {
	osDisk, d := typehelpers.DecodeObjectListOfOne[linuxVirtualMachineOSDiskModel](ctx, input)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	disk := virtualmachines.OSDisk{
		Name:       values.ValueStringPointer(osDisk.Name),
		Caching:    pointer.ToEnum[virtualmachines.CachingTypes](osDisk.Caching.ValueString()),
		DiskSizeGB: values.ValueInt64Pointer(osDisk.DiskSizeGB),
		ManagedDisk: &virtualmachines.ManagedDiskParameters{
			StorageAccountType: pointer.ToEnum[virtualmachines.StorageAccountTypes](osDisk.StorageAccountType.ValueString()),
		},
		WriteAcceleratorEnabled: values.ValueBoolPointer(osDisk.WriteAcceleratorEnabled),

		// these have to be hard-coded so there's no point exposing them
		// for CreateOption, whilst it's possible for this to be "Attach" for an OS Disk
		// from what we can tell this approach has been superseded by provisioning from
		// an image of the machine (e.g. an Image/Shared Image Gallery)
		CreateOption: virtualmachines.DiskCreateOptionTypesFromImage,
		OsType:       pointer.To(osType),
	}

	if v := values.ValueStringPointer(osDisk.SecurityEncryptionType); v != nil {
		disk.ManagedDisk.SecurityProfile = &virtualmachines.VMDiskSecurityProfile{
			SecurityEncryptionType: pointer.To(virtualmachines.SecurityEncryptionTypes(*v)),
		}
	}

	if v := values.ValueStringPointer(osDisk.SecureVMDiskEncryptionSetID); v != nil {
		// if secureVMDiskEncryptionId := osDisk.SecureVMDiskEncryptionSetID.ValueString(); secureVMDiskEncryptionId != "" {
		if virtualmachines.SecurityEncryptionTypesDiskWithVMGuestState != virtualmachines.SecurityEncryptionTypes(osDisk.SecurityEncryptionType.ValueString()) { // TODO - move this to plan modifiers for plan-time catch
			diags.AddError("", "`secure_vm_disk_encryption_set_id` can only be specified when `security_encryption_type` is set to `DiskWithVMGuestState`")
			return nil
		}
		disk.ManagedDisk.SecurityProfile.DiskEncryptionSet = &virtualmachines.SubResource{
			Id: v,
		}
	}

	disk.DiffDiskSettings = expandLinuxVirtualMachineDiffDiskModel(ctx, osDisk.DiffDiskSettings, diags)

	if v := values.ValueStringPointer(osDisk.DiskEncryptionSetID); v != nil {
		disk.ManagedDisk.DiskEncryptionSet = &virtualmachines.SubResource{
			Id: v,
		}
	}

	return &disk
}

func flattenVirtualMachineOSDiskModel(ctx context.Context, input *virtualmachines.OSDisk, diags *diag.Diagnostics) (typehelpers.ListNestedObjectValueOf[linuxVirtualMachineOSDiskModel], types.String) {
	managedDiskIDValue := types.StringNull()

	if input == nil {
		return typehelpers.NewListNestedObjectValueOfNull[linuxVirtualMachineOSDiskModel](ctx), managedDiskIDValue
	}
	osDisk := *input

	result := linuxVirtualMachineOSDiskModel{
		Name:       types.StringPointerValue(osDisk.Name),
		Caching:    types.StringValue(pointer.FromEnum(osDisk.Caching)),
		DiskSizeGB: types.Int64PointerValue(osDisk.DiskSizeGB),

		WriteAcceleratorEnabled: types.BoolPointerValue(osDisk.WriteAcceleratorEnabled),
	}

	if v := osDisk.ManagedDisk; v != nil {
		osDiskId, err := commonids.ParseManagedDiskID(pointer.From(v.Id))
		if err != nil {
			diags.AddError("parsing ID", err.Error())
			return typehelpers.NewListNestedObjectValueOfNull[linuxVirtualMachineOSDiskModel](ctx), managedDiskIDValue
		}

		managedDiskIDValue = types.StringValue(osDiskId.ID())
		result.ID = managedDiskIDValue
		result.StorageAccountType = types.StringValue(pointer.FromEnum(v.StorageAccountType))
		if v.SecurityProfile != nil && v.SecurityProfile.DiskEncryptionSet != nil {
			result.SecureVMDiskEncryptionSetID = types.StringPointerValue(v.SecurityProfile.DiskEncryptionSet.Id)
		}
		if v.DiskEncryptionSet != nil {
			result.DiskEncryptionSetID = types.StringPointerValue(v.DiskEncryptionSet.Id)
		}
	}

	result.DiffDiskSettings = flattenVirtualMachineDiffDiskSettingsModel(ctx, osDisk.DiffDiskSettings)

	return typehelpers.NewListNestedObjectValueOfValueSliceMust(ctx, []linuxVirtualMachineOSDiskModel{result}), managedDiskIDValue
}

type linuxVirtualMachineDiffDiskSettings struct {
	Option    types.String `tfsdk:"option"`
	Placement types.String `tfsdk:"placement"`
}

func expandLinuxVirtualMachineDiffDiskModel(ctx context.Context, input typehelpers.ListNestedObjectValueOf[linuxVirtualMachineDiffDiskSettings], diags *diag.Diagnostics) *virtualmachines.DiffDiskSettings {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	dd, d := typehelpers.DecodeObjectListOfOne(ctx, input)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	result := &virtualmachines.DiffDiskSettings{}

	if !dd.Option.IsNull() {
		result.Option = pointer.ToEnum[virtualmachines.DiffDiskOptions](dd.Option.ValueString())
	}

	if !dd.Placement.IsNull() {
		result.Placement = pointer.ToEnum[virtualmachines.DiffDiskPlacement](dd.Placement.ValueString())
	}

	return result
}

func flattenVirtualMachineDiffDiskSettingsModel(ctx context.Context, input *virtualmachines.DiffDiskSettings) typehelpers.ListNestedObjectValueOf[linuxVirtualMachineDiffDiskSettings] {
	if input == nil {
		return typehelpers.NewListNestedObjectValueOfNull[linuxVirtualMachineDiffDiskSettings](ctx)
	}

	result := linuxVirtualMachineDiffDiskSettings{
		Option:    types.StringValue(pointer.FromEnum(input.Option)),
		Placement: types.StringValue(pointer.FromEnum(input.Placement)),
	}

	return typehelpers.NewListNestedObjectValueOfValueSliceMust(ctx, []linuxVirtualMachineDiffDiskSettings{result})
}

type virtualMachineAdditionalCapabilitiesModel struct {
	UltraSSDEnabled    types.Bool `tfsdk:"ultra_ssd_enabled"`
	HibernationEnabled types.Bool `tfsdk:"hibernation_enabled"`
}

func expandVirtualMachineAdditionalCapabilitiesModel(ctx context.Context, input typehelpers.ListNestedObjectValueOf[virtualMachineAdditionalCapabilitiesModel], diags *diag.Diagnostics) *virtualmachines.AdditionalCapabilities {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	additionalCapabilities, d := typehelpers.DecodeObjectListOfOne(ctx, input)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	return &virtualmachines.AdditionalCapabilities{
		HibernationEnabled: additionalCapabilities.HibernationEnabled.ValueBoolPointer(),
		UltraSSDEnabled:    additionalCapabilities.UltraSSDEnabled.ValueBoolPointer(),
	}
}

func flattenVirtualMachineAdditionalCapabilitiesModel(ctx context.Context, capabilities *virtualmachines.AdditionalCapabilities) typehelpers.ListNestedObjectValueOf[virtualMachineAdditionalCapabilitiesModel] {
	if capabilities == nil {
		return typehelpers.NewListNestedObjectValueOfNull[virtualMachineAdditionalCapabilitiesModel](ctx)
	}

	c := virtualMachineAdditionalCapabilitiesModel{
		UltraSSDEnabled:    types.BoolPointerValue(capabilities.UltraSSDEnabled),
		HibernationEnabled: types.BoolPointerValue(capabilities.HibernationEnabled),
	}

	return typehelpers.NewListNestedObjectValueOfValueSliceMust(ctx, []virtualMachineAdditionalCapabilitiesModel{c})
}

type virtualMachineSSHKeyModel struct {
	PublicKey types.String `tfsdk:"public_key"`
	UserName  types.String `tfsdk:"username"`
}

func expandVirtualMachineSSHKeyModel(ctx context.Context, input typehelpers.SetNestedObjectValueOf[virtualMachineSSHKeyModel], diags *diag.Diagnostics) *[]virtualmachines.SshPublicKey {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	sshKeyList, d := typehelpers.DecodeObjectSet(ctx, input)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	output := make([]virtualmachines.SshPublicKey, 0)

	for _, v := range sshKeyList {
		username := v.UserName.ValueString()
		output = append(output, virtualmachines.SshPublicKey{
			KeyData: v.PublicKey.ValueStringPointer(),
			Path:    pointer.To(formatUsernameForAuthorizedKeysPath(username)),
		})
	}

	return &output
}

func flattenSSHKeyModel(ctx context.Context, configuration *virtualmachines.SshConfiguration) typehelpers.SetNestedObjectValueOf[virtualMachineSSHKeyModel] {
	if configuration == nil {
		return typehelpers.NewSetNestedObjectValueOfNull[virtualMachineSSHKeyModel](ctx)
	}

	result := make([]virtualMachineSSHKeyModel, 0)
	for _, v := range *configuration.PublicKeys {
		if v.Path == nil || v.KeyData == nil {
			continue
		}

		r := virtualMachineSSHKeyModel{
			PublicKey: types.StringPointerValue(v.KeyData),
			UserName: func() types.String {
				n := parseUsernameFromAuthorizedKeysPath(*v.Path)
				return types.StringPointerValue(n)
			}(),
		}

		result = append(result, r)
	}

	return typehelpers.NewSetNestedObjectValueOfValueSliceMust(ctx, result)
}

type virtualMachineBootDiagnosticsModel struct {
	StorageAccountURI types.String `tfsdk:"storage_account_uri"`
}

func expandBootDiagnosticsModel(ctx context.Context, input typehelpers.ListNestedObjectValueOf[virtualMachineBootDiagnosticsModel], diags *diag.Diagnostics) *virtualmachines.BootDiagnostics {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	bootDiagnostics, d := typehelpers.DecodeObjectListOfOne(ctx, input)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	return &virtualmachines.BootDiagnostics{
		Enabled:    pointer.To(bootDiagnostics.StorageAccountURI.ValueString() != ""),
		StorageUri: bootDiagnostics.StorageAccountURI.ValueStringPointer(),
	}
}

func flattenBootDiagnosticsModel(ctx context.Context, profile *virtualmachines.DiagnosticsProfile) typehelpers.ListNestedObjectValueOf[virtualMachineBootDiagnosticsModel] {
	if profile == nil || profile.BootDiagnostics == nil {
		return typehelpers.NewListNestedObjectValueOfNull[virtualMachineBootDiagnosticsModel](ctx)
	}

	r := virtualMachineBootDiagnosticsModel{
		StorageAccountURI: types.StringPointerValue(profile.BootDiagnostics.StorageUri),
	}

	return typehelpers.NewListNestedObjectValueOfValueSliceMust(ctx, []virtualMachineBootDiagnosticsModel{r})
}

type virtualMachineGalleryApplicationModel struct {
	VersionID                   types.String `tfsdk:"version_id"`
	AutomaticUpgradeEnabled     types.Bool   `tfsdk:"automatic_upgrade_enabled"`
	ConfigurationBlobURI        types.String `tfsdk:"configuration_blob_uri"`
	Order                       types.Int64  `tfsdk:"order"`
	Tag                         types.String `tfsdk:"tag"`
	TreatFailureAsDeployFailure types.Bool   `tfsdk:"treat_failure_as_deploy_failure_enabled"`
}

func expandVirtualMachineGalleryApplicationModel(ctx context.Context, input typehelpers.ListNestedObjectValueOf[virtualMachineGalleryApplicationModel], diags *diag.Diagnostics) *[]virtualmachines.VMGalleryApplication {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	out := make([]virtualmachines.VMGalleryApplication, 0)
	if len(input.Elements()) == 0 {
		return &out
	}

	list, d := typehelpers.DecodeObjectList(ctx, input)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	for _, v := range list {
		app := &virtualmachines.VMGalleryApplication{
			PackageReferenceId:              v.VersionID.ValueString(),
			ConfigurationReference:          v.ConfigurationBlobURI.ValueStringPointer(),
			Order:                           v.Order.ValueInt64Pointer(),
			Tags:                            v.Tag.ValueStringPointer(),
			EnableAutomaticUpgrade:          v.AutomaticUpgradeEnabled.ValueBoolPointer(),
			TreatFailureAsDeploymentFailure: v.TreatFailureAsDeployFailure.ValueBoolPointer(),
		}

		out = append(out, *app)
	}

	return &out
}

func flattenVirtualMachineGalleryApplicationModel(ctx context.Context, input *[]virtualmachines.VMGalleryApplication) typehelpers.ListNestedObjectValueOf[virtualMachineGalleryApplicationModel] {
	if input == nil {
		return typehelpers.NewListNestedObjectValueOfNull[virtualMachineGalleryApplicationModel](ctx)
	}

	result := make([]virtualMachineGalleryApplicationModel, 0)
	for _, v := range *input {
		g := virtualMachineGalleryApplicationModel{
			VersionID:                   types.StringValue(v.PackageReferenceId),
			AutomaticUpgradeEnabled:     types.BoolPointerValue(v.EnableAutomaticUpgrade),
			ConfigurationBlobURI:        types.StringPointerValue(v.ConfigurationReference),
			Tag:                         types.StringPointerValue(v.Tags),
			TreatFailureAsDeployFailure: types.BoolPointerValue(v.TreatFailureAsDeploymentFailure),
		}

		result = append(result, g)
	}

	return typehelpers.NewListNestedObjectValueOfValueSliceMust(ctx, result)
}

type virtualMachinePlanModel struct {
	Name      types.String `tfsdk:"name"`
	Product   types.String `tfsdk:"product"`
	Publisher types.String `tfsdk:"publisher"`
}

func expandPlanModel(ctx context.Context, input typehelpers.ListNestedObjectValueOf[virtualMachinePlanModel], diags *diag.Diagnostics) *virtualmachines.Plan {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	vmPlan, d := typehelpers.DecodeObjectListOfOne(ctx, input)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	vmp := &virtualmachines.Plan{} // TODO - Does this work well enough?
	convert.Expand(ctx, &vmPlan, vmp, diags)
	if diags.HasError() {
		return nil
	}

	return vmp
}

func flattenPlanModel(ctx context.Context, input *virtualmachines.Plan) (result typehelpers.ListNestedObjectValueOf[virtualMachinePlanModel]) {
	if input == nil {
		return typehelpers.NewListNestedObjectValueOfNull[virtualMachinePlanModel](ctx)
	}

	p := virtualMachinePlanModel{
		Name:      types.StringPointerValue(input.Name),
		Product:   types.StringPointerValue(input.Product),
		Publisher: types.StringPointerValue(input.Publisher),
	}

	return typehelpers.NewListNestedObjectValueOfValueSliceMust(ctx, []virtualMachinePlanModel{p})
}

type virtualMachineSecretModel struct {
	KeyVaultID   types.String                                                             `tfsdk:"key_vault_id"`
	Certificates typehelpers.SetNestedObjectValueOf[virtualMachineSecretCertificateModel] `tfsdk:"certificate"`
}

func expandLinuxVirtualMachineSecretModel(ctx context.Context, input typehelpers.ListNestedObjectValueOf[virtualMachineSecretModel], diags *diag.Diagnostics) *[]virtualmachines.VaultSecretGroup {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	secretList, d := typehelpers.DecodeObjectList(ctx, input)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	result := make([]virtualmachines.VaultSecretGroup, 0)

	for _, v := range secretList {
		result = append(result, virtualmachines.VaultSecretGroup{
			SourceVault: &virtualmachines.SubResource{
				Id: v.KeyVaultID.ValueStringPointer(),
			},
			VaultCertificates: expandVirtualMachineSecretCertificateModel(ctx, v.Certificates, diags),
		})
		if diags.HasError() {
			return nil
		}
	}

	return &result
}

func flattenLinuxVirtualMachineSecretModel(ctx context.Context, group *[]virtualmachines.VaultSecretGroup) typehelpers.ListNestedObjectValueOf[virtualMachineSecretModel] {
	if group == nil {
		return typehelpers.NewListNestedObjectValueOfNull[virtualMachineSecretModel](ctx)
	}

	result := make([]virtualMachineSecretModel, 0)
	for _, v := range *group {
		if v.SourceVault == nil || v.VaultCertificates == nil {
			continue
		}

		r := virtualMachineSecretModel{
			KeyVaultID:   types.StringPointerValue(v.SourceVault.Id),
			Certificates: flattenVirtualMachineSecretCertificateModel(ctx, v.VaultCertificates),
		}

		result = append(result, r)
	}

	return typehelpers.NewListNestedObjectValueOfValueSliceMust(ctx, result)
}

type virtualMachineSecretCertificateModel struct {
	URL types.String `tfsdk:"url"`
}

func expandVirtualMachineSecretCertificateModel(ctx context.Context, input typehelpers.SetNestedObjectValueOf[virtualMachineSecretCertificateModel], diags *diag.Diagnostics) *[]virtualmachines.VaultCertificate {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}
	result := make([]virtualmachines.VaultCertificate, 0)

	certificatesRaw, d := typehelpers.DecodeObjectSet(ctx, input)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	for _, c := range certificatesRaw {
		result = append(result, virtualmachines.VaultCertificate{
			CertificateURL: c.URL.ValueStringPointer(),
		})
	}

	return &result
}

func flattenVirtualMachineSecretCertificateModel(ctx context.Context, input *[]virtualmachines.VaultCertificate) typehelpers.SetNestedObjectValueOf[virtualMachineSecretCertificateModel] {
	if input == nil {
		return typehelpers.NewSetNestedObjectValueOfNull[virtualMachineSecretCertificateModel](ctx)
	}
	result := make([]virtualMachineSecretCertificateModel, 0)
	for _, v := range *input {
		r := virtualMachineSecretCertificateModel{
			URL: types.StringPointerValue(v.CertificateURL),
		}

		result = append(result, r)
	}

	return typehelpers.NewSetNestedObjectValueOfValueSliceMust(ctx, result)
}

type virtualMachineSourceImageReference struct {
	Publisher types.String `tfsdk:"publisher"`
	Offer     types.String `tfsdk:"offer"`
	Sku       types.String `tfsdk:"sku"`
	Version   types.String `tfsdk:"version"`
}

func expandVirtualMachineSourceImageReference(ctx context.Context, input typehelpers.ListNestedObjectValueOf[virtualMachineSourceImageReference], diags *diag.Diagnostics) *virtualmachines.ImageReference {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	sir, d := typehelpers.DecodeObjectListOfOne(ctx, input)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	return &virtualmachines.ImageReference{
		Publisher: sir.Publisher.ValueStringPointer(),
		Offer:     sir.Offer.ValueStringPointer(),
		Sku:       sir.Sku.ValueStringPointer(),
		Version:   sir.Version.ValueStringPointer(),
	}
}

func flattenVirtualMachineSourceImageReferenceModel(ctx context.Context, input *virtualmachines.ImageReference) typehelpers.ListNestedObjectValueOf[virtualMachineSourceImageReference] {
	if input == nil {
		return typehelpers.NewListNestedObjectValueOfNull[virtualMachineSourceImageReference](ctx)
	}
	result := virtualMachineSourceImageReference{
		Publisher: types.StringPointerValue(input.Publisher),
		Offer:     types.StringPointerValue(input.Offer),
		Sku:       types.StringPointerValue(input.Sku),
		Version:   types.StringPointerValue(input.Version),
	}

	return typehelpers.NewListNestedObjectValueOfValueSliceMust(ctx, []virtualMachineSourceImageReference{result})
}

type virtualMachineOSImageNotificationModel struct {
	Timeout types.String `tfsdk:"timeout"`
}

func expandVirtualMachineOSImageNotificationModel(ctx context.Context, input typehelpers.ListNestedObjectValueOf[virtualMachineOSImageNotificationModel], diags *diag.Diagnostics) *virtualmachines.OSImageNotificationProfile {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	v, d := typehelpers.DecodeObjectListOfOne(ctx, input)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	return &virtualmachines.OSImageNotificationProfile{
		Enable:           pointer.To(true),
		NotBeforeTimeout: v.Timeout.ValueStringPointer(),
	}
}

func flattenVirtualMachineOSImageNotificationModel(ctx context.Context, input *virtualmachines.OSImageNotificationProfile) typehelpers.ListNestedObjectValueOf[virtualMachineOSImageNotificationModel] {
	if input == nil {
		return typehelpers.NewListNestedObjectValueOfNull[virtualMachineOSImageNotificationModel](ctx)
	}
	result := virtualMachineOSImageNotificationModel{
		Timeout: types.StringPointerValue(input.NotBeforeTimeout),
	}

	return typehelpers.NewListNestedObjectValueOfValueSliceMust(ctx, []virtualMachineOSImageNotificationModel{result})
}

type virtualMachineTerminationNotificationModel struct {
	Enabled types.Bool   `tfsdk:"enabled"`
	Timeout types.String `tfsdk:"timeout"`
}

func expandVirtualMachineTerminationNotificationModel(ctx context.Context, input typehelpers.ListNestedObjectValueOf[virtualMachineTerminationNotificationModel], diags *diag.Diagnostics) *virtualmachines.TerminateNotificationProfile {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	v, d := typehelpers.DecodeObjectListOfOne(ctx, input)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	return &virtualmachines.TerminateNotificationProfile{
		Enable:           v.Enabled.ValueBoolPointer(),
		NotBeforeTimeout: v.Timeout.ValueStringPointer(),
	}
}

func flattenVirtualMachineTerminationNotificationModel(ctx context.Context, input *virtualmachines.TerminateNotificationProfile) typehelpers.ListNestedObjectValueOf[virtualMachineTerminationNotificationModel] {
	if input == nil {
		return typehelpers.NewListNestedObjectValueOfNull[virtualMachineTerminationNotificationModel](ctx)
	}
	result := virtualMachineTerminationNotificationModel{
		Enabled: types.BoolPointerValue(input.Enable),
		Timeout: types.StringPointerValue(input.NotBeforeTimeout),
	}

	return typehelpers.NewListNestedObjectValueOfValueSliceMust(ctx, []virtualMachineTerminationNotificationModel{result})
}

func expandVirtualMachineNetworkInterfaceIDsModel(ctx context.Context, input types.List, diags *diag.Diagnostics) *[]virtualmachines.NetworkInterfaceReference {
	result := make([]virtualmachines.NetworkInterfaceReference, 0)

	nicsRaw := make([]string, 0)
	convert.Expand(ctx, input, &nicsRaw, diags)
	if diags.HasError() {
		return nil
	}

	for i, v := range nicsRaw {
		result = append(result, virtualmachines.NetworkInterfaceReference{
			Id: pointer.To(v),
			Properties: &virtualmachines.NetworkInterfaceReferenceProperties{
				Primary: pointer.To(i == 0),
			},
		})
	}

	return &result
}
