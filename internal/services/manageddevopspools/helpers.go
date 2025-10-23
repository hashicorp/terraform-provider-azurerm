package manageddevopspools

import (
	"fmt"
	"encoding/json"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2025-01-21/pools"
)

func expandStatefulAgentProfileModel(input []StatefulAgentProfileModel) (pools.AgentProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	agentProfile := input[0]

	stateful := &pools.Stateful{
		Kind:                "Stateful",
		GracePeriodTimeSpan: agentProfile.GracePeriodTimeSpan,
		MaxAgentLifetime:    agentProfile.MaxAgentLifetime,
	}

	if len(agentProfile.ManualResourcePredictionsProfile) > 0 {
		resourcePredictionsProfile := agentProfile.ManualResourcePredictionsProfile[0]

		resourcePredictions := expandResourcePredictionsModel(resourcePredictionsProfile.ResourcePredictions)
		if resourcePredictions != nil {
			stateful.ResourcePredictions = pointer.To(interface{}(pointer.From(resourcePredictions)))
		}

		manualPredictionsProfile := &pools.ManualResourcePredictionsProfile{
			Kind: pools.ResourcePredictionsProfileTypeManual,
		}
		stateful.ResourcePredictionsProfile = manualPredictionsProfile
	} else if len(agentProfile.AutomaticResourcePredictionsProfile) > 0 {
		automaticPredictionsProfile := &pools.AutomaticResourcePredictionsProfile{
			Kind: pools.ResourcePredictionsProfileTypeAutomatic,
		}

		resourcePredictionsProfile := agentProfile.AutomaticResourcePredictionsProfile[0]
		if resourcePredictionsProfile.PredictionPreference != nil {
			automaticPredictionsProfile.PredictionPreference = (*pools.PredictionPreference)(resourcePredictionsProfile.PredictionPreference)
		}

		stateful.ResourcePredictionsProfile = automaticPredictionsProfile
	}

	return stateful, nil
}

func expandStatelessAgentProfileModel(input []StatelessAgentProfileModel) (pools.AgentProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	agentProfile := input[0]

	stateless := &pools.StatelessAgentProfile{
		Kind: "Stateless",
	}

	if len(agentProfile.ManualResourcePredictionsProfile) > 0 {
		resourcePredictionsProfile := agentProfile.ManualResourcePredictionsProfile[0]

		resourcePredictions := expandResourcePredictionsModel(resourcePredictionsProfile.ResourcePredictions)
		if resourcePredictions != nil {
			stateless.ResourcePredictions = pointer.To(interface{}(pointer.From(resourcePredictions)))
		}

		manualPredictionsProfile := &pools.ManualResourcePredictionsProfile{
			Kind: pools.ResourcePredictionsProfileTypeManual,
		}
		stateless.ResourcePredictionsProfile = manualPredictionsProfile
	} else if len(agentProfile.AutomaticResourcePredictionsProfile) > 0 {
		automaticPredictionsProfile := &pools.AutomaticResourcePredictionsProfile{
			Kind: pools.ResourcePredictionsProfileTypeAutomatic,
		}

		resourcePredictionsProfile := agentProfile.AutomaticResourcePredictionsProfile[0]
		if resourcePredictionsProfile.PredictionPreference != nil {
			automaticPredictionsProfile.PredictionPreference = (*pools.PredictionPreference)(resourcePredictionsProfile.PredictionPreference)
		}

		stateless.ResourcePredictionsProfile = automaticPredictionsProfile
	}

	return stateless, nil
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

func expandAzureDevOpsOrganizationProfileModel(input []AzureDevOpsOrganizationProfileModel) (pools.OrganizationProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	organizationProfile := input[0]
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
		Organizations: poolOrganizations,
	}

	poolPermissionProfile := &pools.AzureDevOpsPermissionProfile{}
	if organizationProfile.PermissionProfileKind != nil {
		poolPermissionProfile.Kind = pools.AzureDevOpsPermissionType(*organizationProfile.PermissionProfileKind)

		// Only set groups and users if the kind is SpecificAccounts
		if *organizationProfile.PermissionProfileKind == string(pools.AzureDevOpsPermissionTypeSpecificAccounts) &&
			len(organizationProfile.AdministratorAccounts) > 0 {
			specificAccounts := organizationProfile.AdministratorAccounts[0]
			poolPermissionProfile.Groups = specificAccounts.Groups
			poolPermissionProfile.Users = specificAccounts.Users
		}

		azureDevOpsOrganizationProfile.PermissionProfile = poolPermissionProfile
	}

	return azureDevOpsOrganizationProfile, nil
}

func expandVmssFabricProfileModel(input []VmssFabricProfileModel) (pools.FabricProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	fabricProfile := input[0]
	vmssFabricProfile := pools.VMSSFabricProfile{
		Images:         expandImageModel(fabricProfile.Images),
		NetworkProfile: expandNetworkProfileModel(fabricProfile.NetworkProfile),
		OsProfile:      expandOsProfileModel(fabricProfile.OsProfile),
		Sku:            expandDevOpsAzureSkuModel(fabricProfile.Sku),
		StorageProfile: expandStorageProfileModel(fabricProfile.StorageProfile),
	}

	return vmssFabricProfile, nil
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

func flattenStatefulAgentProfileToModel(input pools.Stateful) []StatefulAgentProfileModel {
	statefulAgentProfileModel := StatefulAgentProfileModel{
		GracePeriodTimeSpan: input.GracePeriodTimeSpan,
		MaxAgentLifetime:    input.MaxAgentLifetime,
	}

	if input.ResourcePredictionsProfile != nil {
		if automatic, ok := input.ResourcePredictionsProfile.(pools.AutomaticResourcePredictionsProfile); ok {
			statefulAgentProfileModel.AutomaticResourcePredictionsProfile = []AutomaticResourcePredictionsProfileModel{
				{
					PredictionPreference: pointer.To(string(pointer.From(automatic.PredictionPreference))),
				},
			}
		} else if _, ok := input.ResourcePredictionsProfile.(pools.ManualResourcePredictionsProfile); ok {
			statefulAgentProfileModel.ManualResourcePredictionsProfile = []ManualResourcePredictionsProfileModel{
				{
					ResourcePredictions: []ResourcePredictionsModel{},
				},
			}

			if input.ResourcePredictions != nil {
				if predModel := flattenResourcePredictionsModel(pointer.From(input.ResourcePredictions)); predModel != nil {
					if len(statefulAgentProfileModel.ManualResourcePredictionsProfile) > 0 {
						statefulAgentProfileModel.ManualResourcePredictionsProfile[0].ResourcePredictions = []ResourcePredictionsModel{*predModel}
					}
				}
			}
		}
	}

	return []StatefulAgentProfileModel{statefulAgentProfileModel}
}

func flattenStatelessAgentProfileToModel(input pools.StatelessAgentProfile) []StatelessAgentProfileModel {
	statelessAgentProfileModel := StatelessAgentProfileModel{}

	if input.ResourcePredictionsProfile != nil {
		if automatic, ok := input.ResourcePredictionsProfile.(pools.AutomaticResourcePredictionsProfile); ok {
			statelessAgentProfileModel.AutomaticResourcePredictionsProfile = []AutomaticResourcePredictionsProfileModel{
				{
					PredictionPreference: pointer.To(string(pointer.From(automatic.PredictionPreference))),
				},
			}
		} else if _, ok := input.ResourcePredictionsProfile.(pools.ManualResourcePredictionsProfile); ok {
			statelessAgentProfileModel.ManualResourcePredictionsProfile = []ManualResourcePredictionsProfileModel{
				{
					ResourcePredictions: []ResourcePredictionsModel{}, // Will be populated below if ResourcePredictions exists
				},
			}

			if input.ResourcePredictions != nil {
				if predModel := flattenResourcePredictionsModel(pointer.From(input.ResourcePredictions)); predModel != nil {
					if len(statelessAgentProfileModel.ManualResourcePredictionsProfile) > 0 {
						statelessAgentProfileModel.ManualResourcePredictionsProfile[0].ResourcePredictions = []ResourcePredictionsModel{*predModel}
					}
				}
			}
		}
	}
	return []StatelessAgentProfileModel{statelessAgentProfileModel}
}

func flattenResourcePredictionsModel(input interface{}) *ResourcePredictionsModel {
	if input == nil {
		return nil
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return nil
	}

	var sdkModel ResourcePredictionsSdkModel
	if err := json.Unmarshal(jsonBytes, &sdkModel); err != nil {
		return nil
	}

	daysDataBytes, err := json.Marshal(sdkModel.DaysData)
	if err != nil {
		return nil
	}

	return &ResourcePredictionsModel{
		TimeZone: sdkModel.TimeZone,
		DaysData: string(daysDataBytes),
	}
}

func flattenAzureDevOpsOrganizationProfileToModel(input pools.AzureDevOpsOrganizationProfile) []AzureDevOpsOrganizationProfileModel {
	organizationProfileModel := AzureDevOpsOrganizationProfileModel{
		Organizations: flattenOrganizationsToModel(input.Organizations),
	}

	if input.PermissionProfile != nil {
		organizationProfileModel.PermissionProfileKind = (*string)(&input.PermissionProfile.Kind)

		// Only populate specific accounts if it's SpecificAccounts type and has groups/users
		if input.PermissionProfile.Kind == pools.AzureDevOpsPermissionTypeSpecificAccounts {
			specificAccounts := AzureDevOpsAdministratorAccountsModel{
				Groups: input.PermissionProfile.Groups,
				Users:  input.PermissionProfile.Users,
			}
			organizationProfileModel.AdministratorAccounts = []AzureDevOpsAdministratorAccountsModel{specificAccounts}
		}
	}

	return []AzureDevOpsOrganizationProfileModel{organizationProfileModel}
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

func flattenVmssFabricProfileToModel(input pools.VMSSFabricProfile) []VmssFabricProfileModel {
	vmssFabricProfileModel := VmssFabricProfileModel{
		Images:         flattenImagesToModel(input.Images),
		NetworkProfile: flattenNetworkProfileToModel(input.NetworkProfile),
		OsProfile:      flattenOsProfileToModel(input.OsProfile),
		Sku:            flattenDevOpsAzureSkuToModel(input.Sku),
		StorageProfile: flattenStorageProfileToModel(input.StorageProfile),
	}

	return []VmssFabricProfileModel{vmssFabricProfileModel}
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

// identity defined both systemAssigned and userAssigned Identity type in Swagger but only support userAssigned Identity,
// so add a workaround to convert type here.
func expandManagedDevopsToUserAssignedIdentity(input []identity.ModelUserAssigned) (*identity.LegacySystemAndUserAssignedMap, error) {
	if len(input) == 0 {
		return nil, nil
	}

	identityValue, err := identity.ExpandUserAssignedMapFromModel(input)
	if err != nil {
		return nil, fmt.Errorf("expanding `identity`: %+v", err)
	}

	output := identity.LegacySystemAndUserAssignedMap{
		Type:        identityValue.Type,
		IdentityIds: identityValue.IdentityIds,
	}

	return &output, nil
}

func flattenManagedDevopsUserAssignedToLegacyIdentity(input *identity.LegacySystemAndUserAssignedMap) ([]identity.ModelUserAssigned, error) {
	if input == nil {
		return nil, nil
	}

	tmp := identity.UserAssignedMap{
		Type:        input.Type,
		IdentityIds: input.IdentityIds,
	}

	output, err := identity.FlattenUserAssignedMapToModel(&tmp)
	if err != nil {
		return nil, fmt.Errorf("expanding `identity`: %+v", err)
	}

	return *output, nil
}