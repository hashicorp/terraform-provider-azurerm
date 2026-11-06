// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_fleet

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) expandPlanModel(inputList []PlanModel) *fleets.Plan {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.Plan{
		Name:      input.Name,
		Product:   input.Product,
		Publisher: input.Publisher,
	}

	if input.PromotionCode != "" {
		output.PromotionCode = pointer.To(input.PromotionCode)
	}

	return &output
}

func (r Resource) expandOnDemandCapacityModel(inputList []OnDemandCapacityModel) *fleets.RegularPriorityProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.RegularPriorityProfile{
		AllocationStrategy: pointer.ToEnum[fleets.RegularPriorityAllocationStrategy](input.AllocationStrategy),
		Capacity:           pointer.To(input.TargetCapacity),
		MinCapacity:        pointer.To(input.MinimumStartingCapacity),
	}

	return &output
}

func (r Resource) expandSpotCapacityModel(inputList []SpotCapacityModel) *fleets.SpotPriorityProfile {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.SpotPriorityProfile{
		AllocationStrategy: pointer.ToEnum[fleets.SpotAllocationStrategy](input.AllocationStrategy),
		Capacity:           pointer.To(input.TargetCapacity),
		EvictionPolicy:     pointer.ToEnum[fleets.EvictionPolicy](input.EvictionPolicy),
		Maintain:           pointer.To(input.MaintainCapacityEnabled),
		MinCapacity:        pointer.To(input.MinimumCapacity),
	}

	if input.MaximumHourlyPricePerVirtualMachineUsd > 0 || input.MaximumHourlyPricePerVirtualMachineUsd == -1 {
		output.MaxPricePerVM = pointer.To(input.MaximumHourlyPricePerVirtualMachineUsd)
	}
	return &output
}

func (r Resource) expandVirtualMachineSizeProfileModel(inputList []VirtualMachineSizeProfileModel, metadata sdk.ResourceMetaData) *[]fleets.VMSizeProfile {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make([]fleets.VMSizeProfile, 0)
	for i, v := range inputList {
		input := v
		output := fleets.VMSizeProfile{
			Name: input.Name,
		}
		if !metadata.ResourceData.GetRawConfig().AsValueMap()["virtual_machine_sizes_profile"].AsValueSlice()[i].AsValueMap()["rank"].IsNull() {
			output.Rank = pointer.To(input.Rank)
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func (r Resource) expandAdditionalCapabilities(inputList []AdditionalCapabilitiesModel) *fleets.AdditionalCapabilities {
	if len(inputList) == 0 {
		return nil
	}

	capabilities := fleets.AdditionalCapabilities{
		UltraSSDEnabled:    pointer.To(inputList[0].UltraSsdEnabled),
		HibernationEnabled: pointer.To(inputList[0].HibernationEnabled),
	}

	return &capabilities
}

func (r Resource) expandVirtualMachineProfileModel(inputList []VirtualMachineProfileModel, computeFleetModel *ComputeFleetResourceModel, d *schema.ResourceData) (*fleets.BaseVirtualMachineProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := fleets.BaseVirtualMachineProfile{
		OsProfile:          r.expandOSProfileModel(inputList, computeFleetModel),
		ApplicationProfile: r.expandApplicationProfileModel(input.GalleryApplicationProfile),
		DiagnosticsProfile: &fleets.DiagnosticsProfile{
			BootDiagnostics: &fleets.BootDiagnostics{
				Enabled:    pointer.To(input.BootDiagnosticEnabled),
				StorageUri: pointer.To(input.BootDiagnosticStorageAccountEndpoint),
			},
		},
		NetworkProfile: &fleets.VirtualMachineScaleSetNetworkProfile{
			NetworkInterfaceConfigurations: r.expandNetworkInterfaceModel(input.NetworkInterface),
		},
		StorageProfile: &fleets.VirtualMachineScaleSetStorageProfile{
			ImageReference: r.expandImageReference(input.SourceImageReference, input.SourceImageId),
			OsDisk:         r.expandOSDiskModel(input),
			DataDisks:      r.expandDataDiskModel(input.DataDisks),
		},
	}

	output.NetworkProfile.NetworkApiVersion = pointer.To(fleets.NetworkApiVersionTwoZeroTwoZeroNegativeOneOneNegativeZeroOne)
	if input.NetworkApiVersion != "" {
		output.NetworkProfile.NetworkApiVersion = pointer.ToEnum[fleets.NetworkApiVersion](input.NetworkApiVersion)
	}

	if input.ScheduledEventTerminationTimeoutDuration != "" {
		output.ScheduledEventsProfile = &fleets.ScheduledEventsProfile{
			TerminateNotificationProfile: &fleets.TerminateNotificationProfile{
				Enable:           pointer.To(true),
				NotBeforeTimeout: pointer.To(input.ScheduledEventTerminationTimeoutDuration),
			},
		}
	}

	if input.ScheduledEventOsImageTimeoutDuration != "" {
		if output.ScheduledEventsProfile == nil {
			output.ScheduledEventsProfile = &fleets.ScheduledEventsProfile{}
		}
		output.ScheduledEventsProfile.OsImageNotificationProfile = &fleets.OSImageNotificationProfile{
			Enable:           pointer.To(true),
			NotBeforeTimeout: pointer.To(input.ScheduledEventOsImageTimeoutDuration),
		}
	}

	if input.CapacityReservationGroupId != "" {
		output.CapacityReservation = &fleets.CapacityReservationProfile{
			CapacityReservationGroup: r.expandSubResource(input.CapacityReservationGroupId),
		}
	}

	encryptionAtHostEnabledExist := d.GetRawConfig().AsValueMap()["virtual_machine_profile"].AsValueSlice()[0].AsValueMap()["encryption_at_host_enabled"]
	if !encryptionAtHostEnabledExist.IsNull() {
		output.SecurityProfile = &fleets.SecurityProfile{
			EncryptionAtHost: pointer.To(input.EncryptionAtHostEnabled),
		}
	}

	if v := input.OsDisk; len(v) > 0 && v[0].SecurityEncryptionType != "" {
		if output.SecurityProfile == nil {
			output.SecurityProfile = &fleets.SecurityProfile{}
		}
		output.SecurityProfile.SecurityType = pointer.To(fleets.SecurityTypesConfidentialVM)
		output.SecurityProfile.UefiSettings = &fleets.UefiSettings{
			SecureBootEnabled: pointer.To(input.SecureBootEnabled),
			VTpmEnabled:       pointer.To(input.VTpmEnabled),
		}
	} else {
		if input.SecureBootEnabled {
			if output.SecurityProfile == nil {
				output.SecurityProfile = &fleets.SecurityProfile{}
			}
			output.SecurityProfile.UefiSettings = &fleets.UefiSettings{
				SecureBootEnabled: pointer.To(input.SecureBootEnabled),
			}
			output.SecurityProfile.SecurityType = pointer.To(fleets.SecurityTypesTrustedLaunch)
		}

		if input.VTpmEnabled {
			if output.SecurityProfile == nil {
				output.SecurityProfile = &fleets.SecurityProfile{}
			}
			if output.SecurityProfile.UefiSettings == nil {
				output.SecurityProfile.UefiSettings = &fleets.UefiSettings{}
			}
			output.SecurityProfile.UefiSettings.VTpmEnabled = pointer.To(input.VTpmEnabled)

			output.SecurityProfile.SecurityType = pointer.To(fleets.SecurityTypesTrustedLaunch)
		}
	}

	extensionProfileValue, err := r.expandExtensionModel(input.Extension, input.ExtensionsTimeBudgetDuration)
	if err != nil {
		return nil, err
	}
	output.ExtensionProfile = extensionProfileValue

	output.LicenseType = pointer.To("None")
	if input.LicenseType != "" {
		output.LicenseType = pointer.To(input.LicenseType)
	}

	if input.UserDataBase64 != "" {
		output.UserData = pointer.To(input.UserDataBase64)
	}

	return &output, nil
}

func (r Resource) expandApplicationProfileModel(inputList []GalleryApplicationModel) *fleets.ApplicationProfile {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make([]fleets.VMGalleryApplication, 0)
	for _, input := range inputList {
		output := fleets.VMGalleryApplication{
			EnableAutomaticUpgrade:          pointer.To(input.AutomaticUpgradeEnabled),
			Order:                           pointer.To(input.Order),
			PackageReferenceId:              input.VersionId,
			TreatFailureAsDeploymentFailure: pointer.To(input.TreatFailureAsDeploymentFailureEnabled),
		}

		if input.ConfigurationBlobUri != "" {
			output.ConfigurationReference = pointer.To(input.ConfigurationBlobUri)
		}

		if input.Tag != "" {
			output.Tags = pointer.To(input.Tag)
		}
		outputList = append(outputList, output)
	}

	return &fleets.ApplicationProfile{
		GalleryApplications: &outputList,
	}
}

func (r Resource) expandSubResource(input string) *fleets.SubResource {
	if input == "" {
		return nil
	}

	return &fleets.SubResource{
		Id: pointer.To(input),
	}
}

func (r Resource) expandSubResources(inputList []string) *[]fleets.SubResource {
	if len(inputList) == 0 {
		return nil
	}
	outputList := make([]fleets.SubResource, 0)
	for _, input := range inputList {
		output := r.expandSubResource(input)
		if output != nil {
			outputList = append(outputList, pointer.From(output))
		}
	}

	return &outputList
}

func (r Resource) expandExtensionModel(inputList []ExtensionModel, timeBudget string) (*fleets.VirtualMachineScaleSetExtensionProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	output := fleets.VirtualMachineScaleSetExtensionProfile{}
	extensions := make([]fleets.VirtualMachineScaleSetExtension, 0)
	for _, v := range inputList {
		extension := fleets.VirtualMachineScaleSetExtension{
			Name: pointer.To(v.Name),
			Properties: &fleets.VirtualMachineScaleSetExtensionProperties{
				AutoUpgradeMinorVersion: pointer.To(v.AutoUpgradeMinorVersionEnabled),
				SuppressFailures:        pointer.To(v.FailureSuppressionEnabled),
				EnableAutomaticUpgrade:  pointer.To(v.AutomaticUpgradeEnabled),
			},
		}

		if len(v.ProtectedSettingsFromKeyVault) > 0 {
			extension.Properties.ProtectedSettingsFromKeyVault = &fleets.KeyVaultSecretReference{
				SecretURL:   v.ProtectedSettingsFromKeyVault[0].SecretUrl,
				SourceVault: pointer.From(r.expandSubResource(v.ProtectedSettingsFromKeyVault[0].SourceVaultId)),
			}
		}

		if v.ForceExtensionExecutionOnChange != "" {
			extension.Properties.ForceUpdateTag = pointer.To(v.ForceExtensionExecutionOnChange)
		}

		if v.ProtectedSettingsJson != "" {
			protectedSettingsValue := make(map[string]interface{})
			if err := json.Unmarshal([]byte(v.ProtectedSettingsJson), &protectedSettingsValue); err != nil {
				return nil, fmt.Errorf("unmarshaling `protected_settings_json`: %+v", err)
			}
			extension.Properties.ProtectedSettings = pointer.To(protectedSettingsValue)
		}

		if len(v.ExtensionsToProvisionAfterVirtualMachineCreation) > 0 {
			extension.Properties.ProvisionAfterExtensions = pointer.To(v.ExtensionsToProvisionAfterVirtualMachineCreation)
		}

		if v.Publisher != "" {
			extension.Properties.Publisher = pointer.To(v.Publisher)
		}

		if v.SettingsJson != "" {
			result := make(map[string]interface{})
			if err := json.Unmarshal([]byte(v.SettingsJson), &result); err != nil {
				return nil, fmt.Errorf("unmarshaling `settings_json`: %+v", err)
			}
			extension.Properties.Settings = pointer.To(result)
		}

		if v.Type != "" {
			extension.Properties.Type = pointer.To(v.Type)
		}

		if v.TypeHandlerVersion != "" {
			extension.Properties.TypeHandlerVersion = pointer.To(v.TypeHandlerVersion)
		}

		extensions = append(extensions, extension)
	}

	output.Extensions = &extensions

	if timeBudget != "" {
		output.ExtensionsTimeBudget = pointer.To(timeBudget)
	}
	return &output, nil
}

func (r Resource) expandNetworkInterfaceModel(inputList []NetworkInterfaceModel) *[]fleets.VirtualMachineScaleSetNetworkConfiguration {
	outputList := make([]fleets.VirtualMachineScaleSetNetworkConfiguration, 0)
	for _, v := range inputList {
		output := fleets.VirtualMachineScaleSetNetworkConfiguration{
			Name: v.Name,
			Properties: &fleets.VirtualMachineScaleSetNetworkConfigurationProperties{
				EnableAcceleratedNetworking: pointer.To(v.AcceleratedNetworkingEnabled),
				EnableIPForwarding:          pointer.To(v.IPForwardingEnabled),
				NetworkSecurityGroup:        r.expandSubResource(v.NetworkSecurityGroupId),
				Primary:                     pointer.To(v.PrimaryNetworkInterfaceEnabled),
				IPConfigurations:            pointer.From(r.expandIPConfigurationModel(v.IPConfiguration)),
			},
		}

		if len(v.DnsServers) > 0 {
			output.Properties.DnsSettings = &fleets.VirtualMachineScaleSetNetworkConfigurationDnsSettings{
				DnsServers: pointer.To(v.DnsServers),
			}
		}

		auxiliaryMode := fleets.NetworkInterfaceAuxiliaryModeNone
		if v.AuxiliaryMode != "" {
			auxiliaryMode = fleets.NetworkInterfaceAuxiliaryMode(v.AuxiliaryMode)
		}
		output.Properties.AuxiliaryMode = pointer.To(auxiliaryMode)

		if v.DeleteOption != "" {
			output.Properties.DeleteOption = pointer.ToEnum[fleets.DeleteOptions](v.DeleteOption)
		}

		auxiliarySku := fleets.NetworkInterfaceAuxiliarySkuNone
		if v.AuxiliarySku != "" {
			auxiliarySku = fleets.NetworkInterfaceAuxiliarySku(v.AuxiliarySku)
		}
		output.Properties.AuxiliarySku = pointer.To(auxiliarySku)

		outputList = append(outputList, output)
	}
	return &outputList
}

func (r Resource) expandIPConfigurationModel(inputList []IPConfigurationModel) *[]fleets.VirtualMachineScaleSetIPConfiguration {
	outputList := make([]fleets.VirtualMachineScaleSetIPConfiguration, 0)
	for _, input := range inputList {
		output := fleets.VirtualMachineScaleSetIPConfiguration{
			Name: input.Name,
			Properties: &fleets.VirtualMachineScaleSetIPConfigurationProperties{
				ApplicationGatewayBackendAddressPools: r.expandSubResources(input.ApplicationGatewayBackendAddressPoolIds),
				ApplicationSecurityGroups:             r.expandSubResources(input.ApplicationSecurityGroupIds),
				LoadBalancerBackendAddressPools:       r.expandSubResources(input.LoadBalancerBackendAddressPoolIds),
				Primary:                               pointer.To(input.PrimaryIpConfigurationEnabled),
				PublicIPAddressConfiguration:          r.expandPublicIPAddressModel(input.PublicIPAddress),
			},
		}

		if input.SubnetId != "" {
			output.Properties.Subnet = &fleets.ApiEntityReference{
				Id: pointer.To(input.SubnetId),
			}
		}

		if input.Version != "" {
			output.Properties.PrivateIPAddressVersion = pointer.ToEnum[fleets.IPVersion](input.Version)
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func (r Resource) expandPublicIPAddressModel(inputList []PublicIPAddressModel) *fleets.VirtualMachineScaleSetPublicIPAddressConfiguration {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetPublicIPAddressConfiguration{
		Name:       input.Name,
		Properties: &fleets.VirtualMachineScaleSetPublicIPAddressConfigurationProperties{},
	}

	if input.DomainNameLabel != "" {
		output.Properties.DnsSettings = &fleets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings{
			DomainNameLabel: input.DomainNameLabel,
		}
	}

	if input.DomainNameLabelScope != "" {
		if output.Properties.DnsSettings == nil {
			output.Properties.DnsSettings = &fleets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings{}
		}
		output.Properties.DnsSettings.DomainNameLabelScope = pointer.ToEnum[fleets.DomainNameLabelScopeTypes](input.DomainNameLabelScope)
	}

	if v := input.SkuName; v != "" {
		skuParts := strings.Split(v, "_")
		output.Sku = &fleets.PublicIPAddressSku{
			Name: pointer.ToEnum[fleets.PublicIPAddressSkuName](skuParts[0]),
			Tier: pointer.ToEnum[fleets.PublicIPAddressSkuTier](skuParts[1]),
		}
	}

	if input.IdleTimeoutInMinutes > 0 {
		output.Properties.IdleTimeoutInMinutes = pointer.To(input.IdleTimeoutInMinutes)
	}
	if input.PublicIPPrefix != "" {
		output.Properties.PublicIPPrefix = r.expandSubResource(input.PublicIPPrefix)
	}
	if input.DeleteOption != "" {
		output.Properties.DeleteOption = pointer.ToEnum[fleets.DeleteOptions](input.DeleteOption)
	}
	if input.Version != "" {
		output.Properties.PublicIPAddressVersion = pointer.ToEnum[fleets.IPVersion](input.Version)
	}

	return &output
}

func (r Resource) expandOSProfileModel(inputList []VirtualMachineProfileModel, computeFleetModel *ComputeFleetResourceModel) *fleets.VirtualMachineScaleSetOSProfile {
	osProfile := &inputList[0].OsProfile[0]
	output := fleets.VirtualMachineScaleSetOSProfile{
		AllowExtensionOperations: pointer.To(inputList[0].ExtensionOperationsEnabled),
	}
	if osProfile.CustomDataBase64 != "" {
		output.CustomData = pointer.To(osProfile.CustomDataBase64)
	}

	if lConfig := osProfile.LinuxConfiguration; len(lConfig) > 0 {
		linuxConfig := fleets.LinuxConfiguration{
			DisablePasswordAuthentication: pointer.To(!lConfig[0].PasswordAuthenticationEnabled),
			ProvisionVMAgent:              pointer.To(lConfig[0].ProvisionVirtualMachineAgentEnabled),
			EnableVMAgentPlatformUpdates:  pointer.To(lConfig[0].VirtualMachineAgentPlatformUpdatesEnabled),
			PatchSettings: &fleets.LinuxPatchSettings{
				// 'AssessmentMode' can only be set `ImageDefault` on Virtual Machine Scale Sets.
				AssessmentMode: pointer.To(fleets.LinuxPatchAssessmentModeImageDefault),
			},
		}

		// AutomaticByPlatformSettings cannot be set if the PatchMode is not `AutomaticByPlatform`
		if lConfig[0].PatchMode == string(fleets.LinuxVMGuestPatchModeAutomaticByPlatform) {
			linuxConfig.PatchSettings.AutomaticByPlatformSettings = &fleets.LinuxVMGuestPatchAutomaticByPlatformSettings{
				BypassPlatformSafetyChecksOnUserSchedule: pointer.To(lConfig[0].BypassPlatformSafetyChecksEnabled),
			}

			if lConfig[0].PatchRebooting != "" {
				linuxConfig.PatchSettings.AutomaticByPlatformSettings.RebootSetting = pointer.ToEnum[fleets.LinuxVMGuestPatchAutomaticByPlatformRebootSetting](lConfig[0].PatchRebooting)
			}
		}
		if lConfig[0].PatchMode != "" {
			linuxConfig.PatchSettings.PatchMode = pointer.ToEnum[fleets.LinuxVMGuestPatchMode](lConfig[0].PatchMode)
		}
		if lConfig[0].AdminUsername != "" {
			output.AdminUsername = pointer.To(lConfig[0].AdminUsername)
		}
		if lConfig[0].AdminPassword != "" {
			output.AdminPassword = pointer.To(lConfig[0].AdminPassword)
		}

		// `ComputerNamePrefix` default is the first 9 characters of the compute fleet name according to Azure portal
		defaultComputerNamePrefix := computeFleetModel.Name
		if len(defaultComputerNamePrefix) > 9 {
			defaultComputerNamePrefix = defaultComputerNamePrefix[:9]
		}
		output.ComputerNamePrefix = pointer.To(defaultComputerNamePrefix)
		if lConfig[0].ComputerNamePrefix != "" {
			output.ComputerNamePrefix = pointer.To(lConfig[0].ComputerNamePrefix)
		}
		output.Secrets = r.expandOsProfileLinuxSecretsModel(lConfig[0].Secret)

		if lConfig[0].AdminUsername != "" || len(lConfig[0].AdminSSHKeys) > 0 {
			publicKeys := make([]fleets.SshPublicKey, 0)
			for _, v := range lConfig[0].AdminSSHKeys {
				output := fleets.SshPublicKey{
					Path: pointer.To(fmt.Sprintf("/home/%s/.ssh/authorized_keys", lConfig[0].AdminUsername)),
				}
				if v != "" {
					output.KeyData = pointer.To(v)
				}
				publicKeys = append(publicKeys, output)
			}

			linuxConfig.Ssh = &fleets.SshConfiguration{
				PublicKeys: pointer.To(publicKeys),
			}
		}

		output.LinuxConfiguration = &linuxConfig
	}

	if winConfig := osProfile.WindowsConfiguration; len(winConfig) > 0 {
		windowsConfig := fleets.WindowsConfiguration{
			AdditionalUnattendContent:    r.expandAdditionalUnAttendContentModel(winConfig[0].AdditionalUnattendContent),
			EnableAutomaticUpdates:       pointer.To(winConfig[0].AutomaticUpdatesEnabled),
			EnableVMAgentPlatformUpdates: pointer.To(winConfig[0].VirtualMachineAgentPlatformUpdatesEnabled),
			ProvisionVMAgent:             pointer.To(winConfig[0].ProvisionVirtualMachineAgentEnabled),
			PatchSettings: &fleets.PatchSettings{
				EnableHotpatching: pointer.To(winConfig[0].HotPatchingEnabled),
				// 'AssessmentMode' can only be set `ImageDefault` on Virtual Machine Scale Sets.
				AssessmentMode: pointer.To(fleets.WindowsPatchAssessmentModeImageDefault),
			},
		}

		if winRm := winConfig[0].WinRM; len(winRm) > 0 {
			listenerList := make([]fleets.WinRMListener, 0)
			for _, v := range winRm {
				output := fleets.WinRMListener{
					Protocol: pointer.ToEnum[fleets.ProtocolTypes](v.Protocol),
				}

				if v.CertificateUrl != "" {
					output.CertificateURL = pointer.To(v.CertificateUrl)
				}
				listenerList = append(listenerList, output)
			}
			windowsConfig.WinRM = &fleets.WinRMConfiguration{
				Listeners: pointer.To(listenerList),
			}
		}

		if winConfig[0].AdminUsername != "" {
			output.AdminUsername = pointer.To(winConfig[0].AdminUsername)
		}
		if winConfig[0].AdminPassword != "" {
			output.AdminPassword = pointer.To(winConfig[0].AdminPassword)
		}

		// `ComputerNamePrefix` default is the first 9 characters of the compute fleet name according to Azure portal
		defaultComputerNamePrefix := computeFleetModel.Name
		if len(defaultComputerNamePrefix) > 9 {
			defaultComputerNamePrefix = defaultComputerNamePrefix[:9]
		}
		output.ComputerNamePrefix = pointer.To(defaultComputerNamePrefix)
		if winConfig[0].ComputerNamePrefix != "" {
			output.ComputerNamePrefix = pointer.To(winConfig[0].ComputerNamePrefix)
		}

		// AutomaticByPlatformSettings cannot be set if the PatchMode is not `AutomaticByPlatform`
		if winConfig[0].PatchMode == string(fleets.WindowsVMGuestPatchModeAutomaticByPlatform) {
			windowsConfig.PatchSettings.AutomaticByPlatformSettings = &fleets.WindowsVMGuestPatchAutomaticByPlatformSettings{
				BypassPlatformSafetyChecksOnUserSchedule: pointer.To(winConfig[0].BypassPlatformSafetyChecksEnabled),
			}

			if winConfig[0].PatchRebooting != "" {
				windowsConfig.PatchSettings.AutomaticByPlatformSettings.RebootSetting = pointer.ToEnum[fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSetting](winConfig[0].PatchRebooting)
			}
		}
		if winConfig[0].PatchMode != "" {
			windowsConfig.PatchSettings.PatchMode = pointer.ToEnum[fleets.WindowsVMGuestPatchMode](winConfig[0].PatchMode)
		}
		if winConfig[0].TimeZone != "" {
			windowsConfig.TimeZone = pointer.To(winConfig[0].TimeZone)
		}
		output.WindowsConfiguration = &windowsConfig

		output.Secrets = r.expandOsProfileWindowsSecretsModel(winConfig[0].Secret)
	}

	return &output
}

func (r Resource) expandOsProfileLinuxSecretsModel(inputList []LinuxSecretModel) *[]fleets.VaultSecretGroup {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make([]fleets.VaultSecretGroup, 0)
	for _, v := range inputList {
		output := fleets.VaultSecretGroup{
			SourceVault: r.expandSubResource(v.KeyVaultId),
		}

		if len(v.Certificate) > 0 {
			vcs := make([]fleets.VaultCertificate, 0)
			for _, v := range v.Certificate {
				vc := fleets.VaultCertificate{}
				if v.Url != "" {
					vc.CertificateURL = pointer.To(v.Url)
					vcs = append(vcs, vc)
				}
			}
			output.VaultCertificates = &vcs
		}

		outputList = append(outputList, output)
	}
	return &outputList
}

func (r Resource) expandOsProfileWindowsSecretsModel(inputList []WindowsSecretModel) *[]fleets.VaultSecretGroup {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make([]fleets.VaultSecretGroup, 0)
	for _, v := range inputList {
		output := fleets.VaultSecretGroup{
			SourceVault: r.expandSubResource(v.KeyVaultId),
		}

		if len(v.Certificate) > 0 {
			vcs := make([]fleets.VaultCertificate, 0)
			for _, v := range v.Certificate {
				vc := fleets.VaultCertificate{}
				if v.Store != "" {
					vc.CertificateStore = pointer.To(v.Store)
				}
				if v.Url != "" {
					vc.CertificateURL = pointer.To(v.Url)
				}
				vcs = append(vcs, vc)
			}

			output.VaultCertificates = &vcs
		}

		outputList = append(outputList, output)
	}
	return &outputList
}

func (r Resource) expandAdditionalUnAttendContentModel(inputList []AdditionalUnattendContentModel) *[]fleets.AdditionalUnattendContent {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make([]fleets.AdditionalUnattendContent, 0)
	for _, input := range inputList {
		output := fleets.AdditionalUnattendContent{
			SettingName: pointer.ToEnum[fleets.SettingNames](input.Setting),
			Content:     pointer.To(input.Xml),
			// no other possible values
			ComponentName: pointer.To(fleets.ComponentNameMicrosoftNegativeWindowsNegativeShellNegativeSetup),
			PassName:      pointer.To(fleets.PassNameOobeSystem),
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func (r Resource) expandDataDiskModel(inputList []DataDiskModel) *[]fleets.VirtualMachineScaleSetDataDisk {
	outputList := make([]fleets.VirtualMachineScaleSetDataDisk, 0)
	for _, input := range inputList {
		output := fleets.VirtualMachineScaleSetDataDisk{
			CreateOption:            fleets.DiskCreateOptionTypes(input.CreateOption),
			Lun:                     input.Lun,
			WriteAcceleratorEnabled: pointer.To(input.WriteAcceleratorEnabled),
		}

		if input.DeleteOption != "" {
			output.DeleteOption = pointer.ToEnum[fleets.DiskDeleteOptionTypes](input.DeleteOption)
		}

		if input.DiskSizeInGiB > 0 {
			output.DiskSizeGB = pointer.To(input.DiskSizeInGiB)
		}

		caching := string(fleets.CachingTypesNone)
		if input.Caching != "" {
			caching = input.Caching
		}
		output.Caching = pointer.ToEnum[fleets.CachingTypes](caching)

		if input.StorageAccountType != "" {
			output.ManagedDisk = &fleets.VirtualMachineScaleSetManagedDiskParameters{
				StorageAccountType: pointer.ToEnum[fleets.StorageAccountTypes](input.StorageAccountType),
			}
		}

		if input.DiskEncryptionSetId != "" {
			if output.ManagedDisk == nil {
				output.ManagedDisk = &fleets.VirtualMachineScaleSetManagedDiskParameters{}
			}
			output.ManagedDisk.DiskEncryptionSet = &fleets.DiskEncryptionSetParameters{
				Id: pointer.To(input.DiskEncryptionSetId),
			}
		}

		outputList = append(outputList, output)
	}
	return &outputList
}

func (r Resource) expandImageReference(inputList []SourceImageReferenceModel, imageId string) *fleets.ImageReference {
	if imageId != "" {
		// With Version            : "/communityGalleries/publicGalleryName/images/myGalleryImageName/versions/(major.minor.patch | latest)"
		// Versionless(e.g. latest): "/communityGalleries/publicGalleryName/images/myGalleryImageName"
		if _, errs := validation.Any(validate.CommunityGalleryImageID, validate.CommunityGalleryImageVersionID)(imageId, "source_image_id"); len(errs) == 0 {
			return &fleets.ImageReference{
				CommunityGalleryImageId: pointer.To(imageId),
			}
		}

		// Shared Image Gallery with Cross-Tenant Sharing
		// With Version            : "/sharedGalleries/galleryUniqueName/images/myGalleryImageName/versions/(major.minor.patch | latest)"
		// Versionless(e.g. latest): "/sharedGalleries/galleryUniqueName/images/myGalleryImageName"
		if _, errs := validation.Any(validate.SharedGalleryImageID, validate.SharedGalleryImageVersionID)(imageId, "source_image_id"); len(errs) == 0 {
			return &fleets.ImageReference{
				SharedGalleryImageId: pointer.To(imageId),
			}
		}

		return &fleets.ImageReference{
			// Standard Shared Image ID: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/images/{imageDefinitionName}/versions/{imageVersion}
			// Standard Image ID: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/Microsoft.Compute/images/{imageName}
			Id: pointer.To(imageId),
		}
	}

	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	return &fleets.ImageReference{
		Publisher: pointer.To(input.Publisher),
		Offer:     pointer.To(input.Offer),
		Sku:       pointer.To(input.Sku),
		Version:   pointer.To(input.Version),
	}
}

func (r Resource) expandOSDiskModel(input *VirtualMachineProfileModel) *fleets.VirtualMachineScaleSetOSDisk {
	osType := fleets.OperatingSystemTypesLinux
	if len(input.OsProfile) > 0 && len(input.OsProfile[0].WindowsConfiguration) > 0 {
		osType = fleets.OperatingSystemTypesWindows
	}

	if input == nil || len(input.OsDisk) == 0 {
		return nil
	}

	inputOsDisk := &input.OsDisk[0]
	output := fleets.VirtualMachineScaleSetOSDisk{
		DeleteOption:            pointer.ToEnum[fleets.DiskDeleteOptionTypes](inputOsDisk.DeleteOption),
		OsType:                  pointer.To(osType),
		WriteAcceleratorEnabled: pointer.To(inputOsDisk.WriteAcceleratorEnabled),
		ManagedDisk: &fleets.VirtualMachineScaleSetManagedDiskParameters{
			StorageAccountType: pointer.ToEnum[fleets.StorageAccountTypes](inputOsDisk.StorageAccountType),
		},
		// these have to be hard-coded so there's no point exposing them
		CreateOption: fleets.DiskCreateOptionTypesFromImage,
	}

	if inputOsDisk.DiffDiskOption != "" {
		output.DiffDiskSettings = &fleets.DiffDiskSettings{
			Option: pointer.ToEnum[fleets.DiffDiskOptions](inputOsDisk.DiffDiskOption),
		}
	}

	if inputOsDisk.DiffDiskPlacement != "" {
		if output.DiffDiskSettings == nil {
			output.DiffDiskSettings = &fleets.DiffDiskSettings{}
		}
		output.DiffDiskSettings.Placement = pointer.ToEnum[fleets.DiffDiskPlacement](inputOsDisk.DiffDiskPlacement)
	}

	if inputOsDisk.DiskSizeInGiB > 0 {
		output.DiskSizeGB = pointer.To(inputOsDisk.DiskSizeInGiB)
	}

	caching := fleets.CachingTypesNone
	if v := inputOsDisk.Caching; v != "" {
		caching = fleets.CachingTypes(v)
	}
	output.Caching = pointer.To(caching)

	if inputOsDisk.DiskEncryptionSetId != "" {
		output.ManagedDisk.DiskEncryptionSet = &fleets.DiskEncryptionSetParameters{
			Id: pointer.To(inputOsDisk.DiskEncryptionSetId),
		}
	}

	if inputOsDisk.SecurityEncryptionType != "" {
		output.ManagedDisk.SecurityProfile = &fleets.VMDiskSecurityProfile{
			SecurityEncryptionType: pointer.ToEnum[fleets.SecurityEncryptionTypes](inputOsDisk.SecurityEncryptionType),
		}
	}

	return &output
}
