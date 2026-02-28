// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package manageddevopspools

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2025-01-21/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/manageddevopspools/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func manualResourcePredictionSchema(parentPath string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{parentPath + ".automatic_resource_prediction"},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"time_zone_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      "UTC",
					ValidateFunc: validate.ResourcePredictionsProfileTimeZone(),
				},

				"all_week_schedule": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					ConflictsWith: []string{
						parentPath + ".manual_resource_prediction.0.sunday_schedule",
						parentPath + ".manual_resource_prediction.0.monday_schedule",
						parentPath + ".manual_resource_prediction.0.tuesday_schedule",
						parentPath + ".manual_resource_prediction.0.wednesday_schedule",
						parentPath + ".manual_resource_prediction.0.thursday_schedule",
						parentPath + ".manual_resource_prediction.0.friday_schedule",
						parentPath + ".manual_resource_prediction.0.saturday_schedule",
					},
					ValidateFunc: validation.IntAtLeast(1),
				},

				"sunday_schedule": dayScheduleSchemaOptional(parentPath + ".manual_resource_prediction.0.all_week_schedule"),

				"monday_schedule": dayScheduleSchemaOptional(parentPath + ".manual_resource_prediction.0.all_week_schedule"),

				"tuesday_schedule": dayScheduleSchemaOptional(parentPath + ".manual_resource_prediction.0.all_week_schedule"),

				"wednesday_schedule": dayScheduleSchemaOptional(parentPath + ".manual_resource_prediction.0.all_week_schedule"),

				"thursday_schedule": dayScheduleSchemaOptional(parentPath + ".manual_resource_prediction.0.all_week_schedule"),

				"friday_schedule": dayScheduleSchemaOptional(parentPath + ".manual_resource_prediction.0.all_week_schedule"),

				"saturday_schedule": dayScheduleSchemaOptional(parentPath + ".manual_resource_prediction.0.all_week_schedule"),
			},
		},
	}
}

func automaticResourcePredictionSchema(parentPath string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{parentPath + ".manual_resource_prediction"},
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

func dayScheduleSchemaOptional(conflictsWith ...string) *pluginsdk.Schema {
	s := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"time": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d:[0-5]\d$`),
						"must be a valid 24-hour time in format HH:MM:SS",
					),
				},

				"count": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(0),
				},
			},
		},
		ConflictsWith: conflictsWith,
	}
	return s
}

func dayScheduleSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"time": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"count": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func manualResourcePredictionSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"time_zone_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"all_week_schedule": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"sunday_schedule": dayScheduleSchemaComputed(),

				"monday_schedule": dayScheduleSchemaComputed(),

				"tuesday_schedule": dayScheduleSchemaComputed(),

				"wednesday_schedule": dayScheduleSchemaComputed(),

				"thursday_schedule": dayScheduleSchemaComputed(),

				"friday_schedule": dayScheduleSchemaComputed(),

				"saturday_schedule": dayScheduleSchemaComputed(),
			},
		},
	}
}

func automaticResourcePredictionSchemaComputed() *pluginsdk.Schema {
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

func expandStatefulAgentModel(input []StatefulAgentModel) pools.AgentProfile {
	stateful := &pools.Stateful{
		Kind: "Stateful",
	}

	if len(input) == 0 {
		return stateful
	}

	agentProfile := input[0]

	stateful.GracePeriodTimeSpan = pointer.To(agentProfile.GracePeriodTimeSpan)
	stateful.MaxAgentLifetime = pointer.To(agentProfile.MaxAgentLifetime)

	if len(agentProfile.ManualResourcePrediction) > 0 {
		manualResourcePrediction := agentProfile.ManualResourcePrediction[0]

		resourcePredictions := expandResourcePredictionsModel(manualResourcePrediction)
		if resourcePredictions != nil {
			stateful.ResourcePredictions = pointer.To(interface{}(*resourcePredictions))
		}

		manualPredictionProfile := &pools.ManualResourcePredictionsProfile{
			Kind: pools.ResourcePredictionsProfileTypeManual,
		}
		stateful.ResourcePredictionsProfile = manualPredictionProfile
	} else if len(agentProfile.AutomaticResourcePrediction) > 0 {
		automaticResourcePrediction := agentProfile.AutomaticResourcePrediction[0]

		automaticPredictionProfile := &pools.AutomaticResourcePredictionsProfile{
			Kind: pools.ResourcePredictionsProfileTypeAutomatic,
		}

		if automaticResourcePrediction.PredictionPreference != "" {
			predictionPreference := pools.PredictionPreference(automaticResourcePrediction.PredictionPreference)
			automaticPredictionProfile.PredictionPreference = &predictionPreference
		}

		stateful.ResourcePredictionsProfile = automaticPredictionProfile
	}

	return stateful
}

func expandStatelessAgentModel(input []StatelessAgentModel) pools.AgentProfile {
	stateless := &pools.StatelessAgentProfile{
		Kind: "Stateless",
	}

	if len(input) == 0 {
		return stateless
	}

	agentProfile := input[0]

	if len(agentProfile.ManualResourcePrediction) > 0 {
		manualResourcePrediction := agentProfile.ManualResourcePrediction[0]

		resourcePredictions := expandResourcePredictionsModel(manualResourcePrediction)
		if resourcePredictions != nil {
			stateless.ResourcePredictions = pointer.To(interface{}(*resourcePredictions))
		}

		manualPredictionProfile := &pools.ManualResourcePredictionsProfile{
			Kind: pools.ResourcePredictionsProfileTypeManual,
		}
		stateless.ResourcePredictionsProfile = manualPredictionProfile
	} else if len(agentProfile.AutomaticResourcePrediction) > 0 {
		automaticPredictionProfile := &pools.AutomaticResourcePredictionsProfile{
			Kind: pools.ResourcePredictionsProfileTypeAutomatic,
		}

		automaticResourcePrediction := agentProfile.AutomaticResourcePrediction[0]
		if automaticResourcePrediction.PredictionPreference != "" {
			predictionPreference := pools.PredictionPreference(automaticResourcePrediction.PredictionPreference)
			automaticPredictionProfile.PredictionPreference = &predictionPreference
		}

		stateless.ResourcePredictionsProfile = automaticPredictionProfile
	}

	return stateless
}

func expandResourcePredictionsModel(input ManualResourcePredictionModel) *ResourcePredictionsSdkModel {
	var daysData []map[string]int64

	if input.AllWeekSchedule > 0 {
		daysData = []map[string]int64{{"00:00:00": input.AllWeekSchedule}}
	} else {
		daysData = []map[string]int64{
			expandDaySchedule(input.SundaySchedule),
			expandDaySchedule(input.MondaySchedule),
			expandDaySchedule(input.TuesdaySchedule),
			expandDaySchedule(input.WednesdaySchedule),
			expandDaySchedule(input.ThursdaySchedule),
			expandDaySchedule(input.FridaySchedule),
			expandDaySchedule(input.SaturdaySchedule),
		}
	}

	return &ResourcePredictionsSdkModel{
		DaysData: daysData,
		TimeZone: input.TimeZoneName,
	}
}

func expandDaySchedule(input []DayScheduleModel) map[string]int64 {
	m := make(map[string]int64, len(input))
	for _, entry := range input {
		m[entry.Time] = entry.Count
	}
	return m
}

func expandAzureDevOpsOrganizationModel(input []AzureDevOpsOrganizationModel) pools.OrganizationProfile {
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

	if len(organizationProfile.Permission) > 0 {
		permissionProfile := organizationProfile.Permission[0]
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

func expandVmssFabricModel(input []VmssFabricModel) pools.FabricProfile {
	if len(input) == 0 {
		return nil
	}

	fabricProfile := input[0]
	vmssFabricProfile := pools.VMSSFabricProfile{
		Images:         expandImageModel(fabricProfile.Images),
		OsProfile:      expandSecurityModel(fabricProfile.Security),
		Sku:            pools.DevOpsAzureSku{Name: fabricProfile.SkuName},
		StorageProfile: expandStorageModel(fabricProfile.OsDiskStorageAccountType, fabricProfile.Storage),
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

		if image.Id != "" {
			poolImage.ResourceId = pointer.To(image.Id)
		}

		output = append(output, poolImage)
	}

	return output
}

func expandSecurityModel(input []SecurityModel) *pools.OsProfile {
	if len(input) == 0 {
		return nil
	}

	security := input[0]
	logonType := pools.LogonTypeService
	if security.InteractiveLogonEnabled {
		logonType = pools.LogonTypeInteractive
	}
	return &pools.OsProfile{
		LogonType:                 &logonType,
		SecretsManagementSettings: expandKeyVaultManagementSettingsModel(security.KeyVaultManagementSettings),
	}
}

func expandStorageModel(osDiskStorageAccountType string, input []StorageModel) *pools.StorageProfile {
	osDiskType := pools.OsDiskStorageAccountType(osDiskStorageAccountType)
	output := &pools.StorageProfile{
		OsDiskStorageAccountType: &osDiskType,
	}

	if len(input) > 0 {
		disk := input[0]
		cachingType := pools.CachingType(disk.Caching)
		storageAccountType := pools.StorageAccountType(disk.StorageAccountType)
		diskOut := pools.DataDisk{
			Caching:            pointer.To(cachingType),
			DiskSizeGiB:        pointer.To(disk.DiskSizeInGB),
			DriveLetter:        pointer.To(disk.DriveLetter),
			StorageAccountType: pointer.To(storageAccountType),
		}
		output.DataDisks = &[]pools.DataDisk{diskOut}
	}

	return output
}

func expandKeyVaultManagementSettingsModel(input []KeyVaultManagementSettingsModel) *pools.SecretsManagementSettings {
	if len(input) == 0 {
		return nil
	}

	keyVaultManagementSettings := input[0]
	output := &pools.SecretsManagementSettings{
		KeyExportable:        keyVaultManagementSettings.KeyExportable,
		ObservedCertificates: keyVaultManagementSettings.KeyVaultCertificateIds,
	}

	if keyVaultManagementSettings.CertificateStoreLocation != "" {
		output.CertificateStoreLocation = pointer.To(keyVaultManagementSettings.CertificateStoreLocation)
	}

	if keyVaultManagementSettings.CertificateStoreName != "" {
		certificateStoreName := pools.CertificateStoreNameOption(keyVaultManagementSettings.CertificateStoreName)
		output.CertificateStoreName = pointer.To(certificateStoreName)
	}

	return output
}

func flattenStatefulAgentToModel(input pools.Stateful) []StatefulAgentModel {
	statefulAgentModel := StatefulAgentModel{
		GracePeriodTimeSpan: pointer.From(input.GracePeriodTimeSpan),
		MaxAgentLifetime:    pointer.From(input.MaxAgentLifetime),
	}

	if input.ResourcePredictionsProfile != nil {
		if automatic, ok := input.ResourcePredictionsProfile.(pools.AutomaticResourcePredictionsProfile); ok {
			statefulAgentModel.AutomaticResourcePrediction = []AutomaticResourcePredictionModel{
				{
					PredictionPreference: string(pointer.From(automatic.PredictionPreference)),
				},
			}
		} else if _, ok := input.ResourcePredictionsProfile.(pools.ManualResourcePredictionsProfile); ok {
			manualProfile := ManualResourcePredictionModel{}

			if input.ResourcePredictions != nil {
				manualProfile = flattenManualResourcePredictionsModel(pointer.From(input.ResourcePredictions))
			}

			statefulAgentModel.ManualResourcePrediction = []ManualResourcePredictionModel{manualProfile}
		}
	}

	return []StatefulAgentModel{statefulAgentModel}
}

func flattenStatelessAgentToModel(input pools.StatelessAgentProfile) []StatelessAgentModel {
	statelessAgentModel := StatelessAgentModel{}

	if input.ResourcePredictionsProfile != nil {
		if automatic, ok := input.ResourcePredictionsProfile.(pools.AutomaticResourcePredictionsProfile); ok {
			statelessAgentModel.AutomaticResourcePrediction = []AutomaticResourcePredictionModel{
				{
					PredictionPreference: string(pointer.From(automatic.PredictionPreference)),
				},
			}
		} else if _, ok := input.ResourcePredictionsProfile.(pools.ManualResourcePredictionsProfile); ok {
			manualProfile := ManualResourcePredictionModel{}

			if input.ResourcePredictions != nil {
				manualProfile = flattenManualResourcePredictionsModel(pointer.From(input.ResourcePredictions))
			}

			statelessAgentModel.ManualResourcePrediction = []ManualResourcePredictionModel{manualProfile}
		}
	}
	return []StatelessAgentModel{statelessAgentModel}
}

func flattenManualResourcePredictionsModel(input interface{}) ManualResourcePredictionModel {
	manualProfile := ManualResourcePredictionModel{}

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

	manualProfile.TimeZoneName = sdkModel.TimeZone

	if len(sdkModel.DaysData) == 1 {
		for _, count := range sdkModel.DaysData[0] {
			manualProfile.AllWeekSchedule = count
			break
		}
	} else if len(sdkModel.DaysData) == 7 {
		manualProfile.SundaySchedule = flattenDaySchedule(sdkModel.DaysData[0])
		manualProfile.MondaySchedule = flattenDaySchedule(sdkModel.DaysData[1])
		manualProfile.TuesdaySchedule = flattenDaySchedule(sdkModel.DaysData[2])
		manualProfile.WednesdaySchedule = flattenDaySchedule(sdkModel.DaysData[3])
		manualProfile.ThursdaySchedule = flattenDaySchedule(sdkModel.DaysData[4])
		manualProfile.FridaySchedule = flattenDaySchedule(sdkModel.DaysData[5])
		manualProfile.SaturdaySchedule = flattenDaySchedule(sdkModel.DaysData[6])
	}

	return manualProfile
}

func flattenDaySchedule(input map[string]int64) []DayScheduleModel {
	output := make([]DayScheduleModel, 0, len(input))
	for t, count := range input {
		output = append(output, DayScheduleModel{Time: t, Count: count})
	}
	sort.Slice(output, func(i, j int) bool {
		return output[i].Time < output[j].Time
	})
	return output
}

func flattenAzureDevOpsOrganizationToModel(input pools.AzureDevOpsOrganizationProfile) []AzureDevOpsOrganizationModel {
	organizationModel := AzureDevOpsOrganizationModel{
		Organizations: flattenOrganizationsToModel(input.Organizations),
	}

	if input.PermissionProfile != nil {
		permission := AzureDevOpsPermissionModel{
			Kind: string(input.PermissionProfile.Kind),
		}

		if input.PermissionProfile.Kind == pools.AzureDevOpsPermissionTypeSpecificAccounts {
			adminAccounts := AzureDevOpsAdministratorAccountsModel{
				Groups: pointer.From(input.PermissionProfile.Groups),
				Users:  pointer.From(input.PermissionProfile.Users),
			}
			permission.AdministratorAccounts = []AzureDevOpsAdministratorAccountsModel{adminAccounts}
		}

		organizationModel.Permission = []AzureDevOpsPermissionModel{permission}
	}

	return []AzureDevOpsOrganizationModel{organizationModel}
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

func flattenVmssFabricToModel(input pools.VMSSFabricProfile) []VmssFabricModel {
	vmssFabricModel := VmssFabricModel{
		Images:                   flattenImagesToModel(input.Images),
		Security:                 flattenSecurityToModel(input.OsProfile),
		SkuName:                  input.Sku.Name,
		OsDiskStorageAccountType: flattenOsDiskStorageAccountType(input.StorageProfile),
		Storage:                  flattenStorageToModel(input.StorageProfile),
	}

	if input.NetworkProfile != nil {
		vmssFabricModel.SubnetId = input.NetworkProfile.SubnetId
	}

	return []VmssFabricModel{vmssFabricModel}
}

func flattenSecurityToModel(input *pools.OsProfile) []SecurityModel {
	if input == nil {
		return []SecurityModel{}
	}

	securityModel := SecurityModel{
		InteractiveLogonEnabled:    pointer.From(input.LogonType) == pools.LogonTypeInteractive,
		KeyVaultManagementSettings: flattenKeyVaultManagementSettingsToModel(input.SecretsManagementSettings),
	}

	return []SecurityModel{securityModel}
}

func flattenKeyVaultManagementSettingsToModel(input *pools.SecretsManagementSettings) []KeyVaultManagementSettingsModel {
	if input == nil {
		return []KeyVaultManagementSettingsModel{}
	}

	keyvaultManagementSettingsModel := KeyVaultManagementSettingsModel{
		CertificateStoreLocation: pointer.From(input.CertificateStoreLocation),
		KeyExportable:            input.KeyExportable,
		KeyVaultCertificateIds:   input.ObservedCertificates,
	}

	if input.CertificateStoreName != nil {
		keyvaultManagementSettingsModel.CertificateStoreName = string(pointer.From(input.CertificateStoreName))
	}

	return []KeyVaultManagementSettingsModel{keyvaultManagementSettingsModel}
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
			imageModel.Id = pointer.From(image.ResourceId)
		}

		output = append(output, imageModel)
	}

	return output
}

func flattenOsDiskStorageAccountType(input *pools.StorageProfile) string {
	if input == nil || input.OsDiskStorageAccountType == nil {
		return string(pools.OsDiskStorageAccountTypeStandard)
	}
	return string(pointer.From(input.OsDiskStorageAccountType))
}

func flattenStorageToModel(input *pools.StorageProfile) []StorageModel {
	if input == nil || input.DataDisks == nil || len(pointer.From(input.DataDisks)) == 0 {
		return []StorageModel{}
	}

	disk := pointer.From(input.DataDisks)[0]
	return []StorageModel{{
		Caching:            string(pointer.From(disk.Caching)),
		DiskSizeInGB:       pointer.From(disk.DiskSizeGiB),
		DriveLetter:        pointer.From(disk.DriveLetter),
		StorageAccountType: string(pointer.From(disk.StorageAccountType)),
	}}
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
		return []identity.ModelUserAssigned{}, nil
	}

	tmp := identity.UserAssignedMap{
		Type:        input.Type,
		IdentityIds: input.IdentityIds,
	}

	output, err := identity.FlattenUserAssignedMapToModel(&tmp)
	if err != nil {
		return []identity.ModelUserAssigned{}, fmt.Errorf("expanding `identity`: %+v", err)
	}

	return *output, nil
}
