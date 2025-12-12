package manageddevopspools

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2025-01-21/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func manualResourcePredictionsProfileSchema(parentPath string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{parentPath + ".automatic_resource_predictions_profile"},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"time_zone": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "UTC",
				},
				"all_week_schedule": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					ConflictsWith: []string{
						parentPath + ".manual_resource_predictions_profile.0.sunday_schedule",
						parentPath + ".manual_resource_predictions_profile.0.monday_schedule",
						parentPath + ".manual_resource_predictions_profile.0.tuesday_schedule",
						parentPath + ".manual_resource_predictions_profile.0.wednesday_schedule",
						parentPath + ".manual_resource_predictions_profile.0.thursday_schedule",
						parentPath + ".manual_resource_predictions_profile.0.friday_schedule",
						parentPath + ".manual_resource_predictions_profile.0.saturday_schedule",
					},
					ValidateFunc: validation.IntAtLeast(1),
				},
				"sunday_schedule": {
					Type:          pluginsdk.TypeMap,
					Optional:      true,
					ConflictsWith: []string{parentPath + ".manual_resource_predictions_profile.0.all_week_schedule"},
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeInt,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
				"monday_schedule": {
					Type:          pluginsdk.TypeMap,
					Optional:      true,
					ConflictsWith: []string{parentPath + ".manual_resource_predictions_profile.0.all_week_schedule"},
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeInt,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
				"tuesday_schedule": {
					Type:          pluginsdk.TypeMap,
					Optional:      true,
					ConflictsWith: []string{parentPath + ".manual_resource_predictions_profile.0.all_week_schedule"},
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeInt,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
				"wednesday_schedule": {
					Type:          pluginsdk.TypeMap,
					Optional:      true,
					ConflictsWith: []string{parentPath + ".manual_resource_predictions_profile.0.all_week_schedule"},
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeInt,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
				"thursday_schedule": {
					Type:          pluginsdk.TypeMap,
					Optional:      true,
					ConflictsWith: []string{parentPath + ".manual_resource_predictions_profile.0.all_week_schedule"},
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeInt,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
				"friday_schedule": {
					Type:          pluginsdk.TypeMap,
					Optional:      true,
					ConflictsWith: []string{parentPath + ".manual_resource_predictions_profile.0.all_week_schedule"},
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeInt,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
				"saturday_schedule": {
					Type:          pluginsdk.TypeMap,
					Optional:      true,
					ConflictsWith: []string{parentPath + ".manual_resource_predictions_profile.0.all_week_schedule"},
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeInt,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},
	}
}

func automaticResourcePredictionsProfileSchema(parentPath string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{parentPath + ".manual_resource_predictions_profile"},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"prediction_preference": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(pools.PredictionPreferenceBalanced),
					ValidateFunc: validation.StringInSlice(pools.PossibleValuesForPredictionPreference(), false),
				},
			},
		},
	}
}

func manualResourcePredictionsProfileSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"time_zone": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"all_week_schedule": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
				"sunday_schedule": {
					Type:     pluginsdk.TypeMap,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeInt,
					},
				},
				"monday_schedule": {
					Type:     pluginsdk.TypeMap,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeInt,
					},
				},
				"tuesday_schedule": {
					Type:     pluginsdk.TypeMap,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeInt,
					},
				},
				"wednesday_schedule": {
					Type:     pluginsdk.TypeMap,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeInt,
					},
				},
				"thursday_schedule": {
					Type:     pluginsdk.TypeMap,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeInt,
					},
				},
				"friday_schedule": {
					Type:     pluginsdk.TypeMap,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeInt,
					},
				},
				"saturday_schedule": {
					Type:     pluginsdk.TypeMap,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeInt,
					},
				},
			},
		},
	}
}

func automaticResourcePredictionsProfileSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"prediction_preference": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func expandStatefulAgentProfileModel(input []StatefulAgentProfileModel) pools.AgentProfile {
	stateful := &pools.Stateful{
		Kind: "Stateful",
	}

	if len(input) == 0 {
		return stateful
	}

	agentProfile := input[0]

	stateful.GracePeriodTimeSpan = pointer.To(agentProfile.GracePeriodTimeSpan)
	stateful.MaxAgentLifetime = pointer.To(agentProfile.MaxAgentLifetime)

	if len(agentProfile.ManualResourcePredictionsProfile) > 0 {
		resourcePredictionsProfile := agentProfile.ManualResourcePredictionsProfile[0]

		resourcePredictions := expandResourcePredictionsModel(resourcePredictionsProfile)
		if resourcePredictions != nil {
			stateful.ResourcePredictions = pointer.To(interface{}(*resourcePredictions))
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
		if resourcePredictionsProfile.PredictionPreference != "" {
			predictionPreference := pools.PredictionPreference(resourcePredictionsProfile.PredictionPreference)
			automaticPredictionsProfile.PredictionPreference = &predictionPreference
		}

		stateful.ResourcePredictionsProfile = automaticPredictionsProfile
	}

	return stateful
}

func expandStatelessAgentProfileModel(input []StatelessAgentProfileModel) pools.AgentProfile {
	stateless := &pools.StatelessAgentProfile{
		Kind: "Stateless",
	}

	if len(input) == 0 {
		return stateless
	}

	agentProfile := input[0]

	if len(agentProfile.ManualResourcePredictionsProfile) > 0 {
		resourcePredictionsProfile := agentProfile.ManualResourcePredictionsProfile[0]

		resourcePredictions := expandResourcePredictionsModel(resourcePredictionsProfile)
		if resourcePredictions != nil {
			stateless.ResourcePredictions = pointer.To(interface{}(*resourcePredictions))
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
		if resourcePredictionsProfile.PredictionPreference != "" {
			predictionPreference := pools.PredictionPreference(resourcePredictionsProfile.PredictionPreference)
			automaticPredictionsProfile.PredictionPreference = &predictionPreference
		}

		stateless.ResourcePredictionsProfile = automaticPredictionsProfile
	}

	return stateless
}

func expandResourcePredictionsModel(input ManualResourcePredictionsProfileModel) *ResourcePredictionsSdkModel {
	var daysData []map[string]int64

	if input.AllWeekSchedule > 0 {
		allWeekMap := map[string]int64{
			"00:00:00": input.AllWeekSchedule,
		}
		daysData = append(daysData, allWeekMap)
	} else {
		// Per-day schedule - create 7 map entries (one per day)
		daysData = []map[string]int64{
			input.SundaySchedule,    // 0 = Sunday
			input.MondaySchedule,    // 1 = Monday
			input.TuesdaySchedule,   // 2 = Tuesday
			input.WednesdaySchedule, // 3 = Wednesday
			input.ThursdaySchedule,  // 4 = Thursday
			input.FridaySchedule,    // 5 = Friday
			input.SaturdaySchedule,  // 6 = Saturday
		}
	}

	return &ResourcePredictionsSdkModel{
		DaysData: daysData,
		TimeZone: input.TimeZone,
	}
}

func expandAzureDevOpsOrganizationProfileModel(input []AzureDevOpsOrganizationProfileModel) pools.OrganizationProfile {
	if len(input) == 0 {
		return nil
	}

	organizationProfile := input[0]
	poolOrganizations := []pools.Organization{}
	for _, org := range organizationProfile.Organizations {
		poolOrganization := pools.Organization{
			Parallelism: pointer.To(org.Parallelism),
			Projects:    pointer.To(org.Projects),
			Url:         org.Url,
		}
		poolOrganizations = append(poolOrganizations, poolOrganization)
	}

	azureDevOpsOrganizationProfile := pools.AzureDevOpsOrganizationProfile{
		Organizations: poolOrganizations,
	}

	if len(organizationProfile.PermissionProfile) > 0 {
		permissionProfile := organizationProfile.PermissionProfile[0]
		poolPermissionProfile := &pools.AzureDevOpsPermissionProfile{
			Kind: pools.AzureDevOpsPermissionType(permissionProfile.Kind),
		}

		if poolPermissionProfile.Kind == pools.AzureDevOpsPermissionTypeSpecificAccounts && len(permissionProfile.AdministratorAccounts) > 0 {
			adminAccounts := permissionProfile.AdministratorAccounts[0]
			poolPermissionProfile.Groups = pointer.To(adminAccounts.Groups)
			poolPermissionProfile.Users = pointer.To(adminAccounts.Users)
		}

		azureDevOpsOrganizationProfile.PermissionProfile = poolPermissionProfile
	}

	return azureDevOpsOrganizationProfile
}

func expandVmssFabricProfileModel(input []VmssFabricProfileModel) pools.FabricProfile {
	if len(input) == 0 {
		return nil
	}

	fabricProfile := input[0]
	vmssFabricProfile := pools.VMSSFabricProfile{
		Images:         expandImageModel(fabricProfile.Images),
		OsProfile:      expandOsProfileModel(fabricProfile.OsProfile),
		Sku:            pools.DevOpsAzureSku{Name: fabricProfile.SkuName},
		StorageProfile: expandStorageProfileModel(fabricProfile.StorageProfile),
	}

	if fabricProfile.SubnetId != "" {
		vmssFabricProfile.NetworkProfile = &pools.NetworkProfile{
			SubnetId: fabricProfile.SubnetId,
		}
	}

	return vmssFabricProfile
}

func expandImageModel(input []ImageModel) []pools.PoolImage {
	output := []pools.PoolImage{}

	for _, image := range input {
		poolImage := pools.PoolImage{}

		if len(image.Aliases) > 0 {
			poolImage.Aliases = pointer.To(image.Aliases)
		}

		if image.Buffer != "" {
			poolImage.Buffer = pointer.To(image.Buffer)
		}

		// Only apply well_known_image_name or resource_id if they are set, otherwise SDK may throw error
		if image.WellKnownImageName != "" {
			poolImage.WellKnownImageName = pointer.To(image.WellKnownImageName)
		}

		if image.ResourceId != "" {
			poolImage.ResourceId = pointer.To(image.ResourceId)
		}

		output = append(output, poolImage)
	}

	return output
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
			cachingType := pools.CachingType(disk.Caching)
			storageAccountType := pools.StorageAccountType(disk.StorageAccountType)
			diskOut := pools.DataDisk{
				Caching:            pointer.To(cachingType),
				DiskSizeGiB:        pointer.To(disk.DiskSizeGB),
				DriveLetter:        pointer.To(disk.DriveLetter),
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
		KeyExportable:        secretsManagementSettings.KeyExportable,
		ObservedCertificates: secretsManagementSettings.ObservedCertificates,
	}

	if secretsManagementSettings.CertificateStoreLocation != "" {
		output.CertificateStoreLocation = pointer.To(secretsManagementSettings.CertificateStoreLocation)
	}

	if secretsManagementSettings.CertificateStoreName != "" {
		certificateStoreName := pools.CertificateStoreNameOption(secretsManagementSettings.CertificateStoreName)
		output.CertificateStoreName = pointer.To(certificateStoreName)
	}

	return output
}

func flattenStatefulAgentProfileToModel(input pools.Stateful) []StatefulAgentProfileModel {
	statefulAgentProfileModel := StatefulAgentProfileModel{
		GracePeriodTimeSpan: pointer.From(input.GracePeriodTimeSpan),
		MaxAgentLifetime:    pointer.From(input.MaxAgentLifetime),
	}

	if input.ResourcePredictionsProfile != nil {
		if automatic, ok := input.ResourcePredictionsProfile.(pools.AutomaticResourcePredictionsProfile); ok {
			statefulAgentProfileModel.AutomaticResourcePredictionsProfile = []AutomaticResourcePredictionsProfileModel{
				{
					PredictionPreference: string(pointer.From(automatic.PredictionPreference)),
				},
			}
		} else if _, ok := input.ResourcePredictionsProfile.(pools.ManualResourcePredictionsProfile); ok {
			manualProfile := ManualResourcePredictionsProfileModel{}

			if input.ResourcePredictions != nil {
				manualProfile = flattenManualResourcePredictionsModel(pointer.From(input.ResourcePredictions))
			}

			statefulAgentProfileModel.ManualResourcePredictionsProfile = []ManualResourcePredictionsProfileModel{manualProfile}
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
					PredictionPreference: string(pointer.From(automatic.PredictionPreference)),
				},
			}
		} else if _, ok := input.ResourcePredictionsProfile.(pools.ManualResourcePredictionsProfile); ok {
			manualProfile := ManualResourcePredictionsProfileModel{}

			if input.ResourcePredictions != nil {
				manualProfile = flattenManualResourcePredictionsModel(pointer.From(input.ResourcePredictions))
			}

			statelessAgentProfileModel.ManualResourcePredictionsProfile = []ManualResourcePredictionsProfileModel{manualProfile}
		}
	}
	return []StatelessAgentProfileModel{statelessAgentProfileModel}
}

func flattenManualResourcePredictionsModel(input interface{}) ManualResourcePredictionsProfileModel {
	manualProfile := ManualResourcePredictionsProfileModel{}

	if input == nil {
		return manualProfile
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return manualProfile
	}

	var sdkModel ResourcePredictionsSdkModel
	if err := json.Unmarshal(jsonBytes, &sdkModel); err != nil {
		return manualProfile
	}

	manualProfile.TimeZone = sdkModel.TimeZone

	if len(sdkModel.DaysData) == 1 {
		if agentCount, exists := sdkModel.DaysData[0]["00:00:00"]; exists {
			manualProfile.AllWeekSchedule = agentCount
		}
	} else if len(sdkModel.DaysData) == 7 {
		if len(sdkModel.DaysData[0]) > 0 {
			manualProfile.SundaySchedule = sdkModel.DaysData[0]
		}
		if len(sdkModel.DaysData[1]) > 0 {
			manualProfile.MondaySchedule = sdkModel.DaysData[1]
		}
		if len(sdkModel.DaysData[2]) > 0 {
			manualProfile.TuesdaySchedule = sdkModel.DaysData[2]
		}
		if len(sdkModel.DaysData[3]) > 0 {
			manualProfile.WednesdaySchedule = sdkModel.DaysData[3]
		}
		if len(sdkModel.DaysData[4]) > 0 {
			manualProfile.ThursdaySchedule = sdkModel.DaysData[4]
		}
		if len(sdkModel.DaysData[5]) > 0 {
			manualProfile.FridaySchedule = sdkModel.DaysData[5]
		}
		if len(sdkModel.DaysData[6]) > 0 {
			manualProfile.SaturdaySchedule = sdkModel.DaysData[6]
		}
	}

	return manualProfile
}

func flattenAzureDevOpsOrganizationProfileToModel(input pools.AzureDevOpsOrganizationProfile) []AzureDevOpsOrganizationProfileModel {
	organizationProfileModel := AzureDevOpsOrganizationProfileModel{
		Organizations: flattenOrganizationsToModel(input.Organizations),
	}

	if input.PermissionProfile != nil {
		permissionProfile := AzureDevOpsPermissionProfileModel{
			Kind: string(input.PermissionProfile.Kind),
		}

		if input.PermissionProfile.Kind == pools.AzureDevOpsPermissionTypeSpecificAccounts {
			adminAccounts := AzureDevOpsAdministratorAccountsModel{
				Groups: pointer.From(input.PermissionProfile.Groups),
				Users:  pointer.From(input.PermissionProfile.Users),
			}
			permissionProfile.AdministratorAccounts = []AzureDevOpsAdministratorAccountsModel{adminAccounts}
		}

		organizationProfileModel.PermissionProfile = []AzureDevOpsPermissionProfileModel{permissionProfile}
	}

	return []AzureDevOpsOrganizationProfileModel{organizationProfileModel}
}

func flattenOrganizationsToModel(input []pools.Organization) []OrganizationModel {
	output := []OrganizationModel{}

	for _, org := range input {
		organizationModel := OrganizationModel{
			Parallelism: pointer.From(org.Parallelism),
			Projects:    pointer.From(org.Projects),
			Url:         org.Url,
		}
		output = append(output, organizationModel)
	}

	return output
}

func flattenVmssFabricProfileToModel(input pools.VMSSFabricProfile) []VmssFabricProfileModel {
	vmssFabricProfileModel := VmssFabricProfileModel{
		Images:         flattenImagesToModel(input.Images),
		OsProfile:      flattenOsProfileToModel(input.OsProfile),
		SkuName:        input.Sku.Name,
		StorageProfile: flattenStorageProfileToModel(input.StorageProfile),
	}

	if input.NetworkProfile != nil {
		vmssFabricProfileModel.SubnetId = input.NetworkProfile.SubnetId
	}

	return []VmssFabricProfileModel{vmssFabricProfileModel}
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
		CertificateStoreLocation: pointer.From(input.CertificateStoreLocation),
		KeyExportable:            input.KeyExportable,
		ObservedCertificates:     input.ObservedCertificates,
	}

	if input.CertificateStoreName != nil {
		secretsManagementSettingsModel.CertificateStoreName = string(pointer.From(input.CertificateStoreName))
	}

	return []SecretsManagementSettingsModel{secretsManagementSettingsModel}
}

func flattenImagesToModel(input []pools.PoolImage) []ImageModel {
	output := []ImageModel{}

	for _, image := range input {
		imageModel := ImageModel{
			Aliases:            pointer.From(image.Aliases),
			Buffer:             pointer.From(image.Buffer),
			WellKnownImageName: pointer.From(image.WellKnownImageName),
		}

		if image.ResourceId != nil {
			imageModel.ResourceId = pointer.From(image.ResourceId)
		}

		output = append(output, imageModel)
	}

	return output
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
				Caching:            string(pointer.From(disk.Caching)),
				DiskSizeGB:         pointer.From(disk.DiskSizeGiB),
				DriveLetter:        pointer.From(disk.DriveLetter),
				StorageAccountType: string(pointer.From(disk.StorageAccountType)),
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
