package manageddevopspools

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2025-01-21/pools"
)

func expandAgentProfileModel(input []AgentProfileModel) (pools.BaseAgentProfileImpl, error) {
	if len(input) == 0 {
		return nil, nil
	}

	agentProfile := input[0]
	resource_predictions := expandResourcePredictionsModel(agentProfile.ResourcePredictions)
	switch agentProfile.Kind {
	case AgentProfileKindStateful:
		stateful := &pools.Stateful{
			GracePeriodTimeSpan: agentProfile.GracePeriodTimeSpan,
			MaxAgentLifetime:    agentProfile.MaxAgentLifetime,
		}

		if resource_predictions != nil {
			stateful.ResourcePredictions = pointer.To(interface{}(pointer.From(resource_predictions)))
		}

		if len(agentProfile.ResourcePredictionsProfile) > 0 {
			resourcePredictionsProfile := agentProfile.ResourcePredictionsProfile[0]
			if resourcePredictionsProfile.Kind == string(pools.ResourcePredictionsProfileTypeAutomatic) {
				automaticPredictionsProfile := &pools.AutomaticResourcePredictionsProfile{
					Kind:                 pools.ResourcePredictionsProfileTypeAutomatic,
					PredictionPreference: (*pools.PredictionPreference)(resourcePredictionsProfile.PredictionPreference),
				}
				stateful.ResourcePredictionsProfile = automaticPredictionsProfile
			}

			if resourcePredictionsProfile.Kind == string(pools.ResourcePredictionsProfileTypeManual) {
				manualPredictionsProfile := &pools.ManualResourcePredictionsProfile{
					Kind: pools.ResourcePredictionsProfileTypeManual,
				}
				stateful.ResourcePredictionsProfile = manualPredictionsProfile
			}
		}

		return stateful.AgentProfile(), nil

	case AgentProfileKindStateless:
		stateless := &pools.StatelessAgentProfile{
			Kind:                agentProfile.Kind,
			ResourcePredictions: pointer.To(interface{}(expandResourcePredictionsModel(agentProfile.ResourcePredictions))),
		}

		if resource_predictions != nil {
			stateless.ResourcePredictions = pointer.To(interface{}(pointer.From(resource_predictions)))
		}

		if len(agentProfile.ResourcePredictionsProfile) > 0 {
			resourcePredictionsProfile := agentProfile.ResourcePredictionsProfile[0]
			if resourcePredictionsProfile.Kind == string(pools.ResourcePredictionsProfileTypeAutomatic) {
				automaticPredictionsProfile := &pools.AutomaticResourcePredictionsProfile{
					Kind:                 pools.ResourcePredictionsProfileTypeAutomatic,
					PredictionPreference: (*pools.PredictionPreference)(resourcePredictionsProfile.PredictionPreference),
				}
				stateless.ResourcePredictionsProfile = automaticPredictionsProfile
			}

			if resourcePredictionsProfile.Kind == string(pools.ResourcePredictionsProfileTypeManual) {
				manualPredictionsProfile := &pools.ManualResourcePredictionsProfile{
					Kind: pools.ResourcePredictionsProfileTypeManual,
				}
				stateless.ResourcePredictionsProfile = manualPredictionsProfile
			}
		}

		return stateless.AgentProfile(), nil

	default:
		return nil, fmt.Errorf("invalid agent_profile kind provided: %s", agentProfile.Kind)
	}
}

func expandResourcePredictionsModel(input []ResourcePredictionsModel) *ResourcePredictionsSdkModel {
	if len(input) == 0 {
		return nil
	}

	resourcePredictions := input[0]
	var parsedDaysData []map[string]interface{}
	if err := json.Unmarshal([]byte(resourcePredictions.DaysData), &parsedDaysData); err != nil {
		return nil
	}

	return &ResourcePredictionsSdkModel{
		DaysData: parsedDaysData,
		TimeZone: resourcePredictions.TimeZone,
	}
}

func expandOrganizationProfileModel(input []OrganizationProfileModel) (pools.BaseOrganizationProfileImpl, error) {
	if len(input) == 0 {
		return nil, nil
	}

	organizationProfile := input[0]
	if organizationProfile.Kind == "AzureDevOps" {
		poolOrganizations := []pools.Organization{}
		for _, org := range organizationProfile.Organizations {
			poolOrganization := pools.Organization{
				Parallelism: org.Parallelism,
				Projects:    org.Projects,
				Url:         org.Url,
			}
			poolOrganizations = append(poolOrganizations, poolOrganization)
		}

		azureDevOpsOrganizationProfile := pools.AzureDevOpsOrganizationProfile{
			Kind:          organizationProfile.Kind,
			Organizations: poolOrganizations,
		}

		poolPermissionProfile := &pools.AzureDevOpsPermissionProfile{}
		if len(organizationProfile.PermissionProfile) > 0 {
			permissionProfile := organizationProfile.PermissionProfile[0]
			poolPermissionProfile.Groups = permissionProfile.Groups
			poolPermissionProfile.Kind = pools.AzureDevOpsPermissionType(permissionProfile.Kind)
			poolPermissionProfile.Users = permissionProfile.Users

			azureDevOpsOrganizationProfile.PermissionProfile = poolPermissionProfile
		}

		return azureDevOpsOrganizationProfile, nil
	} else {
		return nil, fmt.Errorf("invalid organization_profile `Kind` Provided: %s", organizationProfile.Kind)
	}
}

func expandFabricProfileModel(input []FabricProfileModel) (BaseFabricProfileImpl, error) {
	if len(input) == 0 {
		return nil, nil
	}

	fabricProfile := input[0]
	if fabricProfile.Kind == "Vmss" {
		vmssFabricProfile := pools.VMSSFabricProfile{
			Images:         expandImageModel(fabricProfile.Images),
			NetworkProfile: expandNetworkProfileModel(fabricProfile.NetworkProfile),
			OsProfile:      expandOsProfileModel(fabricProfile.OsProfile),
			Sku:            expandDevOpsAzureSkuModel(fabricProfile.Sku),
			StorageProfile: expandStorageProfileModel(fabricProfile.StorageProfile),
			Kind:           fabricProfile.Kind,
		}

		return vmssFabricProfile, nil
	} else {
		return nil, fmt.Errorf("invalid fabric_profile Kind Provided: %s", fabricProfile.Kind)
	}
}

func expandImageModel(input []ImageModel) []pools.PoolImage {
	output := []pools.PoolImage{}

	for _, image := range input {
		poolImage := pools.PoolImage{
			Aliases:            image.Aliases,
			Buffer:             image.Buffer,
			WellKnownImageName: image.WellKnownImageName,
		}

		if image.ResourceId != nil {
			poolImage.ResourceId = image.ResourceId
		}

		output = append(output, poolImage)
	}

	return output
}

func expandNetworkProfileModel(input []NetworkProfileModel) *pools.NetworkProfile {
	if len(input) == 0 {
		return nil
	}

	networkProfile := input[0]
	return &pools.NetworkProfile{
		SubnetId: networkProfile.SubnetId,
	}
}

func expandOsProfileModel(input []OsProfileModel) *pools.OsProfile {
	if len(input) == 0 {
		return nil
	}

	osProfile := input[0]
	logonType := pools.LogonType(osProfile.LogonType)
	return &pools.OsProfile{
		LogonType:                 &logonType,
		SecretsManagementSettings: expandSecretsManagementSettingsModel(osProfile.SecretsManagementSettings),
	}
}

func expandDevOpsAzureSkuModel(input []DevOpsAzureSkuModel) pools.DevOpsAzureSku {
	return pools.DevOpsAzureSku{
		Name: input[0].Name,
	}
}

func expandStorageProfileModel(input []StorageProfileModel) *pools.StorageProfile {
	if len(input) == 0 {
		return nil
	}

	storageProfile := input[0]
	osDiskStorageAccountType := pools.OsDiskStorageAccountType(storageProfile.OsDiskStorageAccountType)
	output := &pools.StorageProfile{
		OsDiskStorageAccountType: &osDiskStorageAccountType,
	}

	if len(storageProfile.DataDisks) > 0 {
		dataDisksOut := []pools.DataDisk{}
		for _, disk := range storageProfile.DataDisks {
			cachingType := pools.CachingType(pointer.From(disk.Caching))
			storageAccountType := pools.StorageAccountType(pointer.From(disk.StorageAccountType))
			diskOut := pools.DataDisk{
				Caching:            pointer.To(cachingType),
				DiskSizeGiB:        disk.DiskSizeGB,
				DriveLetter:        disk.DriveLetter,
				StorageAccountType: pointer.To(storageAccountType),
			}

			dataDisksOut = append(dataDisksOut, diskOut)
		}

		output.DataDisks = &dataDisksOut
	}

	return output
}

func expandSecretsManagementSettingsModel(input []SecretsManagementSettingsModel) *pools.SecretsManagementSettings {
	if len(input) == 0 {
		return nil
	}

	secretsManagementSettings := input[0]
	output := &pools.SecretsManagementSettings{
		CertificateStoreLocation: secretsManagementSettings.CertificateStoreLocation,
		KeyExportable:            secretsManagementSettings.KeyExportable,
		ObservedCertificates:     secretsManagementSettings.ObservedCertificates,
	}

	if secretsManagementSettings.CertificateStoreName != nil {
		output.CertificateStoreName = pointer.To(pools.CertificateStoreNameOption(pointer.From(secretsManagementSettings.CertificateStoreName)))
	}

	return output
}

func flattenAgentProfileToModel(input pools.AgentProfile) []AgentProfileModel {
	if stateful, ok := input.(pools.Stateful); ok {
		return flattenStatefulAgentProfileToModel(stateful)
	}

	if stateless, ok := input.(pools.StatelessAgentProfile); ok {
		return flattenStatelessAgentProfileToModel(stateless)
	}

	return nil
}

func flattenStatefulAgentProfileToModel(input pools.Stateful) []AgentProfileModel {
	agentProfileModel := AgentProfileModel{
		GracePeriodTimeSpan:        input.GracePeriodTimeSpan,
		Kind:                       AgentProfileKindStateful,
		MaxAgentLifetime:           input.MaxAgentLifetime,
		ResourcePredictionsProfile: flattenResourcePredictionsProfileToModel(input.ResourcePredictionsProfile),
	}

	if input.ResourcePredictions != nil {
		if predictions, ok := (pointer.From(input.ResourcePredictions)).(ResourcePredictionsModel); ok {
			agentProfileModel.ResourcePredictions = []ResourcePredictionsModel{predictions}
		}
	}

	return []AgentProfileModel{agentProfileModel}
}

func flattenStatelessAgentProfileToModel(input pools.StatelessAgentProfile) []AgentProfileModel {
	agentProfileModel := AgentProfileModel{
		Kind:                       AgentProfileKindStateless,
		ResourcePredictionsProfile: flattenResourcePredictionsProfileToModel(input.ResourcePredictionsProfile),
	}

	if input.ResourcePredictions != nil {
		if predictions, ok := (pointer.From(input.ResourcePredictions)).(ResourcePredictionsModel); ok {
			agentProfileModel.ResourcePredictions = []ResourcePredictionsModel{predictions}
		}
	}

	return []AgentProfileModel{agentProfileModel}
}

func flattenResourcePredictionsProfileToModel(input pools.ResourcePredictionsProfile) []ResourcePredictionsProfileModel {
	if automatic, ok := input.(pools.AutomaticResourcePredictionsProfile); ok {
		return flattenAutomaticResourcePredictionsProfileToModel(automatic)
	}

	if manual, ok := input.(pools.ManualResourcePredictionsProfile); ok {
		return flattenManualResourcePredictionsProfileToModel(manual)
	}

	return nil
}

func flattenAutomaticResourcePredictionsProfileToModel(input pools.AutomaticResourcePredictionsProfile) []ResourcePredictionsProfileModel {
	resourcePredictionsProfileModel := ResourcePredictionsProfileModel{
		Kind:                 string(input.Kind),
		PredictionPreference: pointer.To(string(pointer.From(input.PredictionPreference))),
	}

	return []ResourcePredictionsProfileModel{resourcePredictionsProfileModel}
}

func flattenManualResourcePredictionsProfileToModel(input pools.ManualResourcePredictionsProfile) []ResourcePredictionsProfileModel {
	resourcePredictionsProfileModel := ResourcePredictionsProfileModel{
		Kind: string(input.Kind),
	}

	return []ResourcePredictionsProfileModel{resourcePredictionsProfileModel}
}

func flattenOrganizationProfileToModel(input pools.OrganizationProfile) []OrganizationProfileModel {
	if azureDevOps, ok := input.(pools.AzureDevOpsOrganizationProfile); ok {
		return flattenAzureDevOpsOrganizationProfileToModel(azureDevOps)
	}

	return nil
}

func flattenAzureDevOpsOrganizationProfileToModel(input pools.AzureDevOpsOrganizationProfile) []OrganizationProfileModel {
	organizationProfileModel := OrganizationProfileModel{
		Kind:              input.Kind,
		Organizations:     flattenOrganizationsToModel(input.Organizations),
		PermissionProfile: flattenAzureDevOpsPermissionProfileToModel(input.PermissionProfile),
	}

	return []OrganizationProfileModel{organizationProfileModel}
}

func flattenOrganizationsToModel(input []pools.Organization) []OrganizationModel {
	output := []OrganizationModel{}

	for _, org := range input {
		organizationModel := OrganizationModel{
			Parallelism: org.Parallelism,
			Projects:    org.Projects,
			Url:         org.Url,
		}
		output = append(output, organizationModel)
	}

	return output
}

func flattenAzureDevOpsPermissionProfileToModel(input *pools.AzureDevOpsPermissionProfile) []PermissionProfileModel {
	if input == nil {
		return nil
	}

	permissionProfileModel := PermissionProfileModel{
		Groups: input.Groups,
		Kind:   string(input.Kind),
		Users:  input.Users,
	}

	return []PermissionProfileModel{permissionProfileModel}
}

func flattenFabricProfileToModel(input pools.FabricProfile) []FabricProfileModel {
	if vmssProfile, ok := input.(pools.VMSSFabricProfile); ok {
		return flattenVmssFabricProfileToModel(vmssProfile)
	}

	return nil
}

func flattenVmssFabricProfileToModel(input pools.VMSSFabricProfile) []FabricProfileModel {
	fabricProfileModel := FabricProfileModel{
		Images:         flattenImagesToModel(input.Images),
		Kind:           input.Kind,
		NetworkProfile: flattenNetworkProfileToModel(input.NetworkProfile),
		OsProfile:      flattenOsProfileToModel(input.OsProfile),
		Sku:            flattenDevOpsAzureSkuToModel(input.Sku),
		StorageProfile: flattenStorageProfileToModel(input.StorageProfile),
	}

	return []FabricProfileModel{fabricProfileModel}
}

func flattenNetworkProfileToModel(input *pools.NetworkProfile) []NetworkProfileModel {
	if input == nil {
		return nil
	}

	networkProfileModel := NetworkProfileModel{
		SubnetId: input.SubnetId,
	}

	return []NetworkProfileModel{networkProfileModel}
}

func flattenOsProfileToModel(input *pools.OsProfile) []OsProfileModel {
	if input == nil {
		return nil
	}

	osProfileModel := OsProfileModel{
		LogonType:                 string(pointer.From(input.LogonType)),
		SecretsManagementSettings: flattenSecretsManagementSettingsToModel(input.SecretsManagementSettings),
	}

	return []OsProfileModel{osProfileModel}
}

func flattenSecretsManagementSettingsToModel(input *pools.SecretsManagementSettings) []SecretsManagementSettingsModel {
	if input == nil {
		return nil
	}

	secretsManagementSettingsModel := SecretsManagementSettingsModel{
		CertificateStoreLocation: input.CertificateStoreLocation,
		KeyExportable:            input.KeyExportable,
		ObservedCertificates:     input.ObservedCertificates,
	}

	if input.CertificateStoreName != nil {
		secretsManagementSettingsModel.CertificateStoreName = pointer.To(string(pointer.From(input.CertificateStoreName)))
	}

	return []SecretsManagementSettingsModel{secretsManagementSettingsModel}
}

func flattenImagesToModel(input []pools.PoolImage) []ImageModel {
	output := []ImageModel{}

	for _, image := range input {
		imageModel := ImageModel{
			Aliases:            image.Aliases,
			Buffer:             image.Buffer,
			WellKnownImageName: image.WellKnownImageName,
		}

		if image.ResourceId != nil {
			imageModel.ResourceId = image.ResourceId
		}

		output = append(output, imageModel)
	}

	return output
}

func flattenDevOpsAzureSkuToModel(input pools.DevOpsAzureSku) []DevOpsAzureSkuModel {
	return []DevOpsAzureSkuModel{
		{
			Name: input.Name,
		},
	}
}

func flattenStorageProfileToModel(input *pools.StorageProfile) []StorageProfileModel {
	if input == nil {
		return nil
	}

	storageProfileModel := StorageProfileModel{
		OsDiskStorageAccountType: string(pointer.From(input.OsDiskStorageAccountType)),
	}

	if input.DataDisks != nil {
		dataDisksOut := []DataDiskModel{}
		for _, disk := range pointer.From(input.DataDisks) {
			diskOut := DataDiskModel{
				Caching:            pointer.To(string(pointer.From(disk.Caching))),
				DiskSizeGB:         disk.DiskSizeGiB,
				DriveLetter:        disk.DriveLetter,
				StorageAccountType: pointer.To(string(pointer.From(disk.StorageAccountType))),
			}

			dataDisksOut = append(dataDisksOut, diskOut)
		}

		storageProfileModel.DataDisks = dataDisksOut
	}

	return []StorageProfileModel{storageProfileModel}
}
