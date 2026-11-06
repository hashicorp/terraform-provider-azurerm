// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_fleet

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/computefleet/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *fleets.FleetId, model *fleets.Fleet) error {
	state := ComputeFleetResourceModel{
		Name:              id.FleetName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)

		if model.Identity != nil {
			v, err := identity.FlattenSystemAndUserAssignedMapToModel(pointer.To(identity.SystemAndUserAssignedMap(*model.Identity)))
			if err != nil {
				return err
			}
			state.Identity = pointer.From(v)
		}

		state.Plan = r.flattenPlanModel(model.Plan)

		if props := model.Properties; props != nil {
			state.AdditionalCapabilities = r.flattenAdditionalCapabilities(props.ComputeProfile.AdditionalVirtualMachineCapabilities)

			baseVirtualMachineProfileValue, err := r.flattenVirtualMachineProfileModel(&props.ComputeProfile.BaseVirtualMachineProfile, metadata)
			if err != nil {
				return err
			}
			state.VirtualMachineProfile = baseVirtualMachineProfileValue

			state.ComputeApiVersion = pointer.From(props.ComputeProfile.ComputeApiVersion)
			state.PlatformFaultDomainCount = pointer.From(props.ComputeProfile.PlatformFaultDomainCount)
			state.OnDemandCapacity = r.flattenOnDemandCapacityModel(props.RegularPriorityProfile)
			state.SpotCapacity = r.flattenSpotCapacityModel(props.SpotPriorityProfile)
			state.UniqueId = pointer.From(props.UniqueId)
			state.VirtualMachineSizesProfile = r.flattenVirtualMachineSizeProfileModel(&props.VMSizesProfile)
		}
		state.Tags = pointer.From(model.Tags)
		state.Zones = pointer.From(model.Zones)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r Resource) flattenPlanModel(input *fleets.Plan) []PlanModel {
	outputList := make([]PlanModel, 0)
	if input == nil {
		return outputList
	}
	output := PlanModel{
		Name:          input.Name,
		Product:       input.Product,
		Publisher:     input.Publisher,
		PromotionCode: pointer.From(input.PromotionCode),
	}

	return append(outputList, output)
}

func (r Resource) flattenAdditionalCapabilities(input *fleets.AdditionalCapabilities) []AdditionalCapabilitiesModel {
	outputList := make([]AdditionalCapabilitiesModel, 0)
	if input == nil {
		return outputList
	}
	output := AdditionalCapabilitiesModel{
		UltraSsdEnabled:    pointer.From(input.UltraSSDEnabled),
		HibernationEnabled: pointer.From(input.HibernationEnabled),
	}

	return append(outputList, output)
}

func (r Resource) flattenOnDemandCapacityModel(input *fleets.RegularPriorityProfile) []OnDemandCapacityModel {
	outputList := make([]OnDemandCapacityModel, 0)
	if input == nil {
		return outputList
	}

	output := OnDemandCapacityModel{
		AllocationStrategy:      pointer.FromEnum(input.AllocationStrategy),
		TargetCapacity:          pointer.From(input.Capacity),
		MinimumStartingCapacity: pointer.From(input.MinCapacity),
	}

	return append(outputList, output)
}

func (r Resource) flattenSpotCapacityModel(input *fleets.SpotPriorityProfile) []SpotCapacityModel {
	outputList := make([]SpotCapacityModel, 0)
	if input == nil {
		return outputList
	}

	output := SpotCapacityModel{
		AllocationStrategy:      pointer.FromEnum(input.AllocationStrategy),
		TargetCapacity:          pointer.From(input.Capacity),
		EvictionPolicy:          pointer.FromEnum(input.EvictionPolicy),
		MaintainCapacityEnabled: pointer.From(input.Maintain),
		MinimumCapacity:         pointer.From(input.MinCapacity),
	}

	// defaulted since `MaxPricePerVM isn't returned if it's unset
	maximumHourlyPricePerVirtualMachineUsd := float64(-1.0)
	if input.MaxPricePerVM != nil {
		maximumHourlyPricePerVirtualMachineUsd = pointer.From(input.MaxPricePerVM)
	}
	output.MaximumHourlyPricePerVirtualMachineUsd = maximumHourlyPricePerVirtualMachineUsd

	return append(outputList, output)
}

func (r Resource) flattenVirtualMachineSizeProfileModel(inputList *[]fleets.VMSizeProfile) []VirtualMachineSizeProfileModel {
	outputList := make([]VirtualMachineSizeProfileModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := VirtualMachineSizeProfileModel{
			Name: input.Name,
			Rank: pointer.From(input.Rank),
		}
		outputList = append(outputList, output)
	}
	return outputList
}

func (r Resource) flattenVirtualMachineProfileModel(input *fleets.BaseVirtualMachineProfile, metadata sdk.ResourceMetaData) ([]VirtualMachineProfileModel, error) {
	outputList := make([]VirtualMachineProfileModel, 0)
	if input == nil {
		return outputList, nil
	}
	output := VirtualMachineProfileModel{
		GalleryApplicationProfile: r.flattenApplicationProfileModel(input.ApplicationProfile),
		NetworkInterface:          r.flattenNetworkInterfaceModel(input.NetworkProfile),
		UserDataBase64:            pointer.From(input.UserData),
	}

	output.NetworkApiVersion = string(fleets.NetworkApiVersionTwoZeroTwoZeroNegativeOneOneNegativeZeroOne)
	if v := input.NetworkProfile; v != nil {
		output.NetworkApiVersion = pointer.FromEnum(v.NetworkApiVersion)
	}
	if v := input.SecurityProfile; v != nil {
		output.EncryptionAtHostEnabled = pointer.From(v.EncryptionAtHost)
		if v.UefiSettings != nil {
			output.SecureBootEnabled = pointer.From(v.UefiSettings.SecureBootEnabled)
			output.VTpmEnabled = pointer.From(v.UefiSettings.VTpmEnabled)
		}
	}

	if v := input.OsProfile; v != nil {
		osProfile, err := r.flattenOSProfileModel(v, metadata.ResourceData)
		if err != nil {
			return outputList, err
		}
		output.OsProfile = osProfile
		output.ExtensionOperationsEnabled = pointer.From(v.AllowExtensionOperations)
	}

	if v := input.StorageProfile; v != nil {
		output.DataDisks = r.flattenDataDiskModel(v.DataDisks)
		storageImageId := ""
		if v.ImageReference != nil {
			if v.ImageReference.Id != nil {
				storageImageId = pointer.From(v.ImageReference.Id)
			}
			if v.ImageReference.CommunityGalleryImageId != nil {
				storageImageId = *v.ImageReference.CommunityGalleryImageId
			}
			if v.ImageReference.SharedGalleryImageId != nil {
				storageImageId = *v.ImageReference.SharedGalleryImageId
			}
		}
		output.SourceImageId = storageImageId
		output.SourceImageReference = r.flattenImageReference(v.ImageReference, storageImageId != "")
		output.OsDisk = r.flattenOSDiskModel(v.OsDisk)
	}

	if se := input.ScheduledEventsProfile; se != nil {
		if v := se.TerminateNotificationProfile; v != nil {
			if pointer.From(v.Enable) {
				output.ScheduledEventTerminationTimeoutDuration = pointer.From(v.NotBeforeTimeout)
			}
		}
		if v := se.OsImageNotificationProfile; v != nil {
			if pointer.From(v.Enable) {
				output.ScheduledEventOsImageTimeoutDuration = pointer.From(v.NotBeforeTimeout)
			}
		}
	}

	if cr := input.CapacityReservation; cr != nil {
		if v := cr.CapacityReservationGroup; v != nil {
			output.CapacityReservationGroupId = pointer.From(v.Id)
		}
	}

	if dp := input.DiagnosticsProfile; dp != nil {
		if v := dp.BootDiagnostics; v != nil {
			output.BootDiagnosticEnabled = pointer.From(v.Enabled)
			output.BootDiagnosticStorageAccountEndpoint = pointer.From(v.StorageUri)
		}
	}

	extensionProfileValue, err := r.flattenExtensionModel(input.ExtensionProfile, metadata)
	if err != nil {
		return outputList, err
	}
	output.Extension = extensionProfileValue

	if input.ExtensionProfile != nil {
		output.ExtensionsTimeBudgetDuration = pointer.From(input.ExtensionProfile.ExtensionsTimeBudget)
	}

	licenseType := ""
	if v := pointer.From(input.LicenseType); v != "None" {
		licenseType = v
	}
	output.LicenseType = licenseType

	return append(outputList, output), nil
}

func (r Resource) flattenAdminSshKeyModel(input *fleets.SshConfiguration) ([]string, error) {
	outputList := make([]string, 0)
	if input == nil || input.PublicKeys == nil {
		return outputList, nil
	}

	for _, input := range *input.PublicKeys {
		username := helpers.ParseUsernameFromAuthorizedKeysPath(*input.Path)
		if username == nil {
			return outputList, fmt.Errorf("parsing username from %q", pointer.From(input.Path))
		}
		outputList = append(outputList, pointer.From(input.KeyData))
	}

	return outputList, nil
}

func (r Resource) flattenApplicationProfileModel(input *fleets.ApplicationProfile) []GalleryApplicationModel {
	outputList := make([]GalleryApplicationModel, 0)
	if input == nil {
		return outputList
	}

	for _, input := range *input.GalleryApplications {
		output := GalleryApplicationModel{
			VersionId:                              input.PackageReferenceId,
			ConfigurationBlobUri:                   pointer.From(input.ConfigurationReference),
			AutomaticUpgradeEnabled:                pointer.From(input.EnableAutomaticUpgrade),
			Order:                                  pointer.From(input.Order),
			Tag:                                    pointer.From(input.Tags),
			TreatFailureAsDeploymentFailureEnabled: pointer.From(input.TreatFailureAsDeploymentFailure),
		}
		outputList = append(outputList, output)
	}

	return outputList
}

func (r Resource) flattenNetworkInterfaceModel(input *fleets.VirtualMachineScaleSetNetworkProfile) []NetworkInterfaceModel {
	outputList := make([]NetworkInterfaceModel, 0)
	if input == nil {
		return outputList
	}

	for _, input := range *input.NetworkInterfaceConfigurations {
		output := NetworkInterfaceModel{
			Name: input.Name,
		}

		if props := input.Properties; props != nil {
			if v := props.AuxiliaryMode; v != nil && *v != fleets.NetworkInterfaceAuxiliaryModeNone {
				output.AuxiliaryMode = string(*v)
			}

			if v := props.AuxiliarySku; v != nil && *v != fleets.NetworkInterfaceAuxiliarySkuNone {
				output.AuxiliarySku = string(*v)
			}

			output.DeleteOption = pointer.FromEnum(props.DeleteOption)

			if v := props.DnsSettings; v != nil {
				output.DnsServers = pointer.From(v.DnsServers)
			}

			output.AcceleratedNetworkingEnabled = pointer.From(props.EnableAcceleratedNetworking)

			output.IPForwardingEnabled = pointer.From(props.EnableIPForwarding)

			output.IPConfiguration = r.flattenIPConfigurationModel(props.IPConfigurations)

			if v := props.NetworkSecurityGroup; v != nil {
				output.NetworkSecurityGroupId = pointer.From(v.Id)
			}

			output.PrimaryNetworkInterfaceEnabled = pointer.From(props.Primary)
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func (r Resource) flattenOsProfileLinuxSecretsModel(inputList *[]fleets.VaultSecretGroup) []LinuxSecretModel {
	outputList := make([]LinuxSecretModel, 0)
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := LinuxSecretModel{
			Certificate: r.flattenLinuxVaultCertificateModel(input.VaultCertificates),
		}
		if v := input.SourceVault; v != nil {
			output.KeyVaultId = pointer.From(v.Id)
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func (r Resource) flattenOsProfileWindowsSecretsModel(inputList *[]fleets.VaultSecretGroup) []WindowsSecretModel {
	outputList := make([]WindowsSecretModel, 0)
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := WindowsSecretModel{
			Certificate: r.flattenWindowsVaultCertificateModel(input.VaultCertificates),
		}
		if v := input.SourceVault; v != nil {
			output.KeyVaultId = pointer.From(v.Id)
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func (r Resource) flattenOSProfileModel(input *fleets.VirtualMachineScaleSetOSProfile, d *schema.ResourceData) ([]OSProfileModel, error) {
	outputList := make([]OSProfileModel, 0)
	if input == nil {
		return outputList, nil
	}

	output := OSProfileModel{}
	if input.CustomData != nil {
		output.CustomDataBase64 = pointer.From(input.CustomData)
	} else {
		output.CustomDataBase64 = d.Get("virtual_machine_profile.0.os_profile.0.custom_data_base64").(string)
	}

	windowsConfigs := make([]WindowsConfigurationModel, 0)
	if v := input.WindowsConfiguration; v != nil {
		windowsConfig := WindowsConfigurationModel{
			AdditionalUnattendContent: r.flattenAdditionalUnAttendContentModel(v.AdditionalUnattendContent, d),
			WinRM:                     r.flattenWinRMModel(v.WinRM),
			AdminUsername:             pointer.From(input.AdminUsername),
			AdminPassword:             d.Get("virtual_machine_profile.0.os_profile.0.windows_configuration.0.admin_password").(string),
			AutomaticUpdatesEnabled:   pointer.From(v.EnableAutomaticUpdates),
			ComputerNamePrefix:        pointer.From(input.ComputerNamePrefix),
			VirtualMachineAgentPlatformUpdatesEnabled: pointer.From(v.EnableVMAgentPlatformUpdates),
			ProvisionVirtualMachineAgentEnabled:       pointer.From(v.ProvisionVMAgent),
			TimeZone:                                  pointer.From(v.TimeZone),
			Secret:                                    r.flattenOsProfileWindowsSecretsModel(input.Secrets),
		}

		if p := v.PatchSettings; p != nil {
			windowsConfig.PatchMode = pointer.FromEnum(p.PatchMode)
			if a := p.AutomaticByPlatformSettings; a != nil {
				windowsConfig.BypassPlatformSafetyChecksEnabled = pointer.From(a.BypassPlatformSafetyChecksOnUserSchedule)
				windowsConfig.PatchRebooting = pointer.FromEnum(a.RebootSetting)
			}
			windowsConfig.HotPatchingEnabled = pointer.From(p.EnableHotpatching)
		}
		windowsConfigs = append(windowsConfigs, windowsConfig)
	}
	output.WindowsConfiguration = windowsConfigs

	linuxConfigs := make([]LinuxConfigurationModel, 0)
	if v := input.LinuxConfiguration; v != nil {
		linuxConfig := LinuxConfigurationModel{
			AdminUsername:                             pointer.From(input.AdminUsername),
			AdminPassword:                             d.Get("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password").(string),
			ComputerNamePrefix:                        pointer.From(input.ComputerNamePrefix),
			PasswordAuthenticationEnabled:             !pointer.From(v.DisablePasswordAuthentication),
			VirtualMachineAgentPlatformUpdatesEnabled: pointer.From(v.EnableVMAgentPlatformUpdates),
			ProvisionVirtualMachineAgentEnabled:       pointer.From(v.ProvisionVMAgent),
			Secret:                                    r.flattenOsProfileLinuxSecretsModel(input.Secrets),
		}

		if p := v.PatchSettings; p != nil {
			linuxConfig.PatchMode = pointer.FromEnum(p.PatchMode)
			if a := p.AutomaticByPlatformSettings; a != nil {
				linuxConfig.BypassPlatformSafetyChecksEnabled = pointer.From(a.BypassPlatformSafetyChecksOnUserSchedule)
				linuxConfig.PatchRebooting = pointer.FromEnum(a.RebootSetting)
			}
		}

		flattenedSSHPublicKeys, err := r.flattenAdminSshKeyModel(v.Ssh)
		if err != nil {
			return outputList, fmt.Errorf("flattening `linux_configuration.0.admin_ssh_keys`: %+v", err)
		}
		linuxConfig.AdminSSHKeys = flattenedSSHPublicKeys

		linuxConfigs = append(linuxConfigs, linuxConfig)
	}

	output.LinuxConfiguration = linuxConfigs

	return append(outputList, output), nil
}

func (r Resource) flattenLinuxVaultCertificateModel(inputList *[]fleets.VaultCertificate) []LinuxCertificateModel {
	outputList := make([]LinuxCertificateModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := LinuxCertificateModel{
			Url: pointer.From(input.CertificateURL),
		}
		outputList = append(outputList, output)
	}
	return outputList
}

func (r Resource) flattenWindowsVaultCertificateModel(inputList *[]fleets.VaultCertificate) []WindowsCertificateModel {
	outputList := make([]WindowsCertificateModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := WindowsCertificateModel{
			Store: pointer.From(input.CertificateStore),
			Url:   pointer.From(input.CertificateURL),
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func (r Resource) flattenAdditionalUnAttendContentModel(inputList *[]fleets.AdditionalUnattendContent, d *schema.ResourceData) []AdditionalUnattendContentModel {
	outputList := make([]AdditionalUnattendContentModel, 0)
	if inputList == nil {
		return outputList
	}
	for i, input := range *inputList {
		output := AdditionalUnattendContentModel{
			Setting: pointer.FromEnum(input.SettingName),
		}
		existing := make([]interface{}, 0)
		if v, ok := d.GetOk("virtual_machine_profile.0.os_profile.0.windows_configuration.0.additional_unattend_content"); ok {
			existing = v.([]interface{})
		}

		// content isn't returned by the API since it's sensitive data so we need to pull from the state file.
		content := ""
		if len(existing) > i {
			existingVal := existing[i]
			existingRaw, ok := existingVal.(map[string]interface{})
			if ok {
				contentRaw, ok := existingRaw["xml"]
				if ok {
					content = contentRaw.(string)
				}
			}
		}
		output.Xml = content

		outputList = append(outputList, output)
	}
	return outputList
}

func (r Resource) flattenWinRMModel(input *fleets.WinRMConfiguration) []WinRMModel {
	outputList := make([]WinRMModel, 0)
	if input == nil || input.Listeners == nil {
		return outputList
	}

	for _, input := range *input.Listeners {
		output := WinRMModel{
			CertificateUrl: pointer.From(input.CertificateURL),
			Protocol:       pointer.FromEnum(input.Protocol),
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func (r Resource) flattenDataDiskModel(inputList *[]fleets.VirtualMachineScaleSetDataDisk) []DataDiskModel {
	outputList := make([]DataDiskModel, 0)
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := DataDiskModel{
			CreateOption:            string(input.CreateOption),
			Lun:                     input.Lun,
			DeleteOption:            pointer.FromEnum(input.DeleteOption),
			DiskSizeInGiB:           pointer.From(input.DiskSizeGB),
			WriteAcceleratorEnabled: pointer.From(input.WriteAcceleratorEnabled),
		}

		caching := ""
		if v := input.Caching; v != nil && *v != fleets.CachingTypesNone {
			caching = string(*v)
		}
		output.Caching = caching

		if md := input.ManagedDisk; md != nil {
			if v := md.DiskEncryptionSet; v != nil {
				output.DiskEncryptionSetId = pointer.From(v.Id)
			}
			output.StorageAccountType = pointer.FromEnum(md.StorageAccountType)
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func (r Resource) flattenImageReference(input *fleets.ImageReference, hasImageId bool) []SourceImageReferenceModel {
	outputList := make([]SourceImageReferenceModel, 0)
	// since the image id is pulled out as a separate field, if that's set we should return an empty block here
	if input == nil || hasImageId {
		return outputList
	}
	output := SourceImageReferenceModel{
		Publisher: pointer.From(input.Publisher),
		Offer:     pointer.From(input.Offer),
		Sku:       pointer.From(input.Sku),
		Version:   pointer.From(input.Version),
	}

	return append(outputList, output)
}

func (r Resource) flattenOSDiskModel(input *fleets.VirtualMachineScaleSetOSDisk) []OSDiskModel {
	outputList := make([]OSDiskModel, 0)
	if input == nil {
		return outputList
	}

	output := OSDiskModel{
		DeleteOption:            pointer.FromEnum(input.DeleteOption),
		WriteAcceleratorEnabled: pointer.From(input.WriteAcceleratorEnabled),
	}

	if v := input.DiffDiskSettings; v != nil {
		output.DiffDiskOption = pointer.FromEnum(v.Option)
		output.DiffDiskPlacement = pointer.FromEnum(v.Placement)
	}

	caching := ""
	if v := input.Caching; v != nil && *v != fleets.CachingTypesNone {
		caching = string(*v)
	}
	output.Caching = caching

	if input.DiskSizeGB != nil {
		output.DiskSizeInGiB = pointer.From(input.DiskSizeGB)
	}

	if md := input.ManagedDisk; md != nil {
		if v := md.DiskEncryptionSet; v != nil {
			output.DiskEncryptionSetId = pointer.From(v.Id)
		}
		if sp := md.SecurityProfile; sp != nil {
			output.SecurityEncryptionType = pointer.FromEnum(sp.SecurityEncryptionType)
		}
		output.StorageAccountType = pointer.FromEnum(md.StorageAccountType)
	}

	return append(outputList, output)
}

func (r Resource) flattenExtensionModel(input *fleets.VirtualMachineScaleSetExtensionProfile, metadata sdk.ResourceMetaData) ([]ExtensionModel, error) {
	outputList := make([]ExtensionModel, 0)
	if input == nil || input.Extensions == nil {
		return outputList, nil
	}

	for i, input := range *input.Extensions {
		output := ExtensionModel{
			Name: pointer.From(input.Name),
		}

		if props := input.Properties; props != nil {
			output.Publisher = pointer.From(props.Publisher)
			output.Type = pointer.From(props.Type)
			output.TypeHandlerVersion = pointer.From(props.TypeHandlerVersion)
			output.AutoUpgradeMinorVersionEnabled = pointer.From(props.AutoUpgradeMinorVersion)
			output.FailureSuppressionEnabled = pointer.From(props.SuppressFailures)
			output.AutomaticUpgradeEnabled = pointer.From(props.EnableAutomaticUpgrade)
			output.ForceExtensionExecutionOnChange = pointer.From(props.ForceUpdateTag)
			// Sensitive data isn't returned, so we get it from config
			output.ProtectedSettingsJson = metadata.ResourceData.Get("virtual_machine_profile.0.extension." + strconv.Itoa(i) + ".protected_settings_json").(string)
			output.ProtectedSettingsFromKeyVault = r.flattenProtectedSettingsFromKeyVaultModel(props.ProtectedSettingsFromKeyVault)
			output.ExtensionsToProvisionAfterVirtualMachineCreation = pointer.From(props.ProvisionAfterExtensions)
			extSettings := ""
			if props.Settings != nil {
				settings, err := json.Marshal(props.Settings)
				if err != nil {
					return outputList, fmt.Errorf("unmarshaling `settings`: %+v", err)
				}

				extSettings = string(settings)
			}
			output.SettingsJson = extSettings
		}

		outputList = append(outputList, output)
	}

	return outputList, nil
}

func (r Resource) flattenProtectedSettingsFromKeyVaultModel(input *fleets.KeyVaultSecretReference) []ProtectedSettingsFromKeyVaultModel {
	outputList := make([]ProtectedSettingsFromKeyVaultModel, 0)
	if input == nil {
		return outputList
	}

	output := ProtectedSettingsFromKeyVaultModel{
		SecretUrl:     input.SecretURL,
		SourceVaultId: pointer.From(input.SourceVault.Id),
	}

	return append(outputList, output)
}

func (r Resource) flattenIPConfigurationModel(inputList []fleets.VirtualMachineScaleSetIPConfiguration) []IPConfigurationModel {
	outputList := make([]IPConfigurationModel, 0)
	if len(inputList) == 0 {
		return outputList
	}
	for _, input := range inputList {
		output := IPConfigurationModel{
			Name: input.Name,
		}
		if props := input.Properties; props != nil {
			output.PrimaryIpConfigurationEnabled = pointer.From(props.Primary)
			output.Version = pointer.FromEnum(props.PrivateIPAddressVersion)

			addressPools := make([]string, 0)
			if v := props.ApplicationGatewayBackendAddressPools; v != nil {
				addressPools = r.flattenSubResourceId(*v)
			}
			output.ApplicationGatewayBackendAddressPoolIds = addressPools

			lbAddressPools := make([]string, 0)
			if v := props.LoadBalancerBackendAddressPools; v != nil {
				lbAddressPools = r.flattenSubResourceId(*v)
			}
			output.LoadBalancerBackendAddressPoolIds = lbAddressPools

			groupIds := make([]string, 0)
			if v := props.ApplicationSecurityGroups; v != nil {
				groupIds = r.flattenSubResourceId(*v)
			}
			output.ApplicationSecurityGroupIds = groupIds

			if v := props.PublicIPAddressConfiguration; v != nil {
				output.PublicIPAddress = r.flattenPublicIPAddressModel(v)
			}

			if v := props.Subnet; v != nil {
				output.SubnetId = pointer.From(v.Id)
			}
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func (r Resource) flattenPublicIPAddressModel(input *fleets.VirtualMachineScaleSetPublicIPAddressConfiguration) []PublicIPAddressModel {
	outputList := make([]PublicIPAddressModel, 0)
	if input == nil {
		return outputList
	}
	output := PublicIPAddressModel{
		Name: input.Name,
	}

	if v := input.Sku; v != nil {
		if v.Name != nil && v.Tier != nil {
			output.SkuName = fmt.Sprintf("%s_%s", pointer.From(v.Name), pointer.From(v.Tier))
		}
	}

	if props := input.Properties; props != nil {
		output.DeleteOption = pointer.FromEnum(props.DeleteOption)
		if v := props.DnsSettings; v != nil {
			output.DomainNameLabel = v.DomainNameLabel
			output.DomainNameLabelScope = pointer.FromEnum(v.DomainNameLabelScope)
		}
		output.IdleTimeoutInMinutes = pointer.From(props.IdleTimeoutInMinutes)

		if publicIpPrefix := props.PublicIPPrefix; publicIpPrefix != nil {
			output.PublicIPPrefix = pointer.From(publicIpPrefix.Id)
		}

		output.Version = pointer.FromEnum(props.PublicIPAddressVersion)
	}
	return append(outputList, output)
}

func (r Resource) flattenSubResourceId(inputList []fleets.SubResource) []string {
	outputList := make([]string, 0)
	if len(inputList) == 0 {
		return outputList
	}
	for _, input := range inputList {
		outputList = append(outputList, pointer.From(input.Id))
	}
	return outputList
}
