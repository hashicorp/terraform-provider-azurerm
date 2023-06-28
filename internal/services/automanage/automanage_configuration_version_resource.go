package automanage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/automanage/2022-05-04/automanage"
)

type ConfigurationVersionModel struct {
	Name                     string `tfschema:"name"`
	ResourceGroupName        string `tfschema:"resource_group_name"`
	ConfigurationProfileName string `tfschema:"configuration_profile_name"`

	Antimalware               []AntimalwareConfiguration           `tfschema:"antimalware"`
	AzureSecurityBaseline     []AzureSecurityBaselineConfiguration `tfschema:"azure_security_baseline"`
	Backup                    []BackupConfiguration                `tfschema:"backup"`
	LogAnalyticsEnabled       bool                                 `tfschema:"log_analytics_enabled"`
	AutomationAccountEnabled  bool                                 `tfschema:"automation_account_enabled"`
	BootDiagnosticsEnabled    bool                                 `tfschema:"boot_diagnostics_enabled"`
	DefenderForCloudEnabled   bool                                 `tfschema:"defender_for_cloud_enabled"`
	GuestConfigurationEnabled bool                                 `tfschema:"guest_configuration_enabled"`
	StatusChangeAlertEnabled  bool                                 `tfschema:"status_change_alert_enabled"`

	Location string            `tfschema:"location"`
	Tags     map[string]string `tfschema:"tags"`
}

type AutoManageConfigurationVersionResource struct{}

var _ sdk.ResourceWithUpdate = AutoManageConfigurationVersionResource{}

func (r AutoManageConfigurationVersionResource) ResourceType() string {
	return "azurerm_automanage_configuration_version"
}

func (r AutoManageConfigurationVersionResource) ModelObject() interface{} {
	return &ConfigurationVersionModel{}
}

func (r AutoManageConfigurationVersionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.AutomanageConfigurationVersionID
}

func (r AutoManageConfigurationVersionResource) Arguments() map[string]*pluginsdk.Schema {
	c := AutoManageConfigurationResource{}
	arguments := c.Arguments()
	arguments["configuration_profile_name"] = &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}
	return arguments
}

func (r AutoManageConfigurationVersionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AutoManageConfigurationVersionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ConfigurationVersionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Automanage.ConfigurationVersionClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := parse.NewAutomanageConfigurationVersionID(subscriptionId, model.ResourceGroupName, model.ConfigurationProfileName, model.Name)
			existing, err := client.Get(ctx, id.ConfigurationProfileName, id.VersionName, id.ResourceGroup)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := automanage.ConfigurationProfile{
				Location:   utils.String(location.Normalize(model.Location)),
				Properties: &automanage.ConfigurationProfileProperties{},
				Tags:       tags.FromTypedObject(model.Tags),
			}

			properties.Properties.Configuration = expandAutomanageConfigurationVersionProfile(model)

			if _, err := client.CreateOrUpdate(ctx, id.ConfigurationProfileName, id.VersionName, id.ResourceGroup, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AutoManageConfigurationVersionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationVersionClient

			id, err := parse.AutomanageConfigurationVersionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ConfigurationProfileName, id.VersionName, id.ResourceGroup)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ConfigurationVersionModel{
				Name:                     id.VersionName,
				ConfigurationProfileName: id.ConfigurationProfileName,
				ResourceGroupName:        id.ResourceGroup,
				Location:                 location.NormalizeNilable(resp.Location),
			}

			if resp.Properties != nil && resp.Properties.Configuration != nil {
				configMap := resp.Properties.Configuration.(map[string]interface{})

				state.Antimalware = flattenAntimarewareConfig(configMap)

				state.AzureSecurityBaseline = flattenAzureSecurityBaselineConfig(configMap)

				state.Backup = flattenBackupConfig(configMap)

				if val, ok := configMap["AutomationAccount/Enable"]; ok {
					state.AutomationAccountEnabled = val.(bool)
				}

				if val, ok := configMap["BootDiagnostics/Enable"]; ok {
					state.BootDiagnosticsEnabled = val.(bool)
				}

				if val, ok := configMap["DefenderForCloud/Enable"]; ok {
					state.DefenderForCloudEnabled = val.(bool)
				}

				if val, ok := configMap["GuestConfiguration/Enable"]; ok {
					state.GuestConfigurationEnabled = val.(bool)
				}

				if val, ok := configMap["LogAnalytics/Enable"]; ok {
					state.LogAnalyticsEnabled = val.(bool)
				}

				if val, ok := configMap["Alerts/AutomanageStatusChanges/Enable"]; ok {
					state.StatusChangeAlertEnabled = val.(bool)
				}
			}

			if resp.Tags != nil {
				state.Tags = tags.ToTypedObject(resp.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AutoManageConfigurationVersionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationVersionClient

			id, err := parse.AutomanageConfigurationVersionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.ResourceGroup, id.ConfigurationProfileName, id.VersionName); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r AutoManageConfigurationVersionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationVersionClient

			id, err := parse.AutomanageConfigurationVersionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ConfigurationVersionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := automanage.ConfigurationProfile{
				Location: utils.String(location.Normalize(metadata.ResourceData.Get("location").(string))),
				Properties: &automanage.ConfigurationProfileProperties{
					Configuration: expandAutomanageConfigurationVersionProfile(model),
				},
				Tags: tags.Expand(metadata.ResourceData.Get("tags").(map[string]interface{})),
			}

			if _, err := client.CreateOrUpdate(ctx, id.ConfigurationProfileName, id.VersionName, id.ResourceGroup, properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandAutomanageConfigurationVersionProfile(model ConfigurationVersionModel) *map[string]interface{} {
	// building configuration profile in json format
	jsonConfig := make(map[string]interface{})

	if model.Antimalware != nil && len(model.Antimalware) > 0 {
		antimalwareConfig := model.Antimalware[0]
		jsonConfig["Antimalware/Enable"] = true
		jsonConfig["Antimalware/EnableRealTimeProtection"] = antimalwareConfig.RealTimeProtectionEnabled
		jsonConfig["Antimalware/RunScheduledScan"] = antimalwareConfig.ScheduledScanEnabled
		jsonConfig["Antimalware/ScanType"] = antimalwareConfig.ScanType
		jsonConfig["Antimalware/ScanDay"] = antimalwareConfig.ScanDay
		jsonConfig["Antimalware/ScanTimeInMinutes"] = antimalwareConfig.ScanTimeInMinutes
		if antimalwareConfig.Exclusions != nil && len(antimalwareConfig.Exclusions) > 0 {
			jsonConfig["Antimalware/Exclusions/Extensions"] = antimalwareConfig.Exclusions[0].Extensions
			jsonConfig["Antimalware/Exclusions/Paths"] = antimalwareConfig.Exclusions[0].Paths
			jsonConfig["Antimalware/Exclusions/Processes"] = antimalwareConfig.Exclusions[0].Processes
		}
	}

	if model.AzureSecurityBaseline != nil && len(model.AzureSecurityBaseline) > 0 {
		azureSecurityBaselineConfig := model.AzureSecurityBaseline[0]
		jsonConfig["AzureSecurityBaseline/Enable"] = true
		jsonConfig["AzureSecurityBaseline/AssignmentType"] = azureSecurityBaselineConfig.AssignmentType
	}

	if model.Backup != nil && len(model.Backup) > 0 {
		backupConfig := model.Backup[0]
		jsonConfig["Backup/Enable"] = true
		if backupConfig.PolicyName != "" {
			jsonConfig["Backup/PolicyName"] = backupConfig.PolicyName
		}
		jsonConfig["Backup/TimeZone"] = backupConfig.TimeZone
		jsonConfig["Backup/InstantRpRetentionRangeInDays"] = backupConfig.InstantRpRetentionRangeInDays
		if backupConfig.SchedulePolicy != nil && len(backupConfig.SchedulePolicy) > 0 {
			schedulePolicyConfig := backupConfig.SchedulePolicy[0]
			jsonConfig["Backup/SchedulePolicy/ScheduleRunFrequency"] = schedulePolicyConfig.ScheduleRunFrequency
			if schedulePolicyConfig.ScheduleRunTimes != nil && len(schedulePolicyConfig.ScheduleRunTimes) > 0 {
				jsonConfig["Backup/SchedulePolicy/ScheduleRunTimes"] = schedulePolicyConfig.ScheduleRunTimes
			}
			if schedulePolicyConfig.ScheduleRunDays != nil && len(schedulePolicyConfig.ScheduleRunDays) > 0 {
				jsonConfig["Backup/SchedulePolicy/ScheduleRunDays"] = schedulePolicyConfig.ScheduleRunDays
			}
			jsonConfig["Backup/SchedulePolicy/SchedulePolicyType"] = schedulePolicyConfig.SchedulePolicyType
		}

		if backupConfig.RetentionPolicy != nil && len(backupConfig.RetentionPolicy) > 0 {
			retentionPolicyConfig := backupConfig.RetentionPolicy[0]
			jsonConfig["Backup/RetentionPolicy/RetentionPolicyType"] = retentionPolicyConfig.RetentionPolicyType
			if retentionPolicyConfig.DailySchedule != nil && len(retentionPolicyConfig.DailySchedule) > 0 {
				dailyScheduleConfig := retentionPolicyConfig.DailySchedule[0]
				if dailyScheduleConfig.RetentionTimes != nil && len(dailyScheduleConfig.RetentionTimes) > 0 {
					jsonConfig["Backup/RetentionPolicy/DailySchedule/RetentionTimes"] = dailyScheduleConfig.RetentionTimes
				}

				if dailyScheduleConfig.RetentionDuration != nil && len(dailyScheduleConfig.RetentionDuration) > 0 {
					jsonConfig["Backup/RetentionPolicy/DailySchedule/RetentionDuration/Count"] = dailyScheduleConfig.RetentionDuration[0].Count
					jsonConfig["Backup/RetentionPolicy/DailySchedule/RetentionDuration/DurationType"] = dailyScheduleConfig.RetentionDuration[0].DurationType
				}
			}

			if retentionPolicyConfig.WeeklySchedule != nil && len(retentionPolicyConfig.WeeklySchedule) > 0 {
				weeklyScheduleConfig := retentionPolicyConfig.WeeklySchedule[0]
				if weeklyScheduleConfig.RetentionTimes != nil && len(weeklyScheduleConfig.RetentionTimes) > 0 {
					jsonConfig["Backup/RetentionPolicy/WeeklySchedule/RetentionTimes"] = weeklyScheduleConfig.RetentionTimes
				}

				if weeklyScheduleConfig.RetentionDuration != nil && len(weeklyScheduleConfig.RetentionDuration) > 0 {
					jsonConfig["Backup/RetentionPolicy/WeeklySchedule/RetentionDuration/Count"] = weeklyScheduleConfig.RetentionDuration[0].Count
					jsonConfig["Backup/RetentionPolicy/WeeklySchedule/RetentionDuration/DurationType"] = weeklyScheduleConfig.RetentionDuration[0].DurationType
				}
			}
		}
	}

	if model.AutomationAccountEnabled {
		jsonConfig["AutomationAccount/Enable"] = model.AutomationAccountEnabled
	}

	if model.BootDiagnosticsEnabled {
		jsonConfig["BootDiagnostics/Enable"] = model.BootDiagnosticsEnabled
	}

	if model.DefenderForCloudEnabled {
		jsonConfig["DefenderForCloud/Enable"] = model.DefenderForCloudEnabled
	}

	if model.GuestConfigurationEnabled {
		jsonConfig["GuestConfiguration/Enable"] = model.GuestConfigurationEnabled
	}

	if model.LogAnalyticsEnabled {
		jsonConfig["LogAnalytics/Enable"] = model.LogAnalyticsEnabled
	}

	if model.StatusChangeAlertEnabled {
		jsonConfig["Alerts/AutomanageStatusChanges/Enable"] = model.StatusChangeAlertEnabled
	}
	return &jsonConfig
}
