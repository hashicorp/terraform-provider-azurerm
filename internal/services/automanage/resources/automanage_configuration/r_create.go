// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automanage_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationProfilesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := configurationprofiles.NewConfigurationProfileID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := configurationprofiles.ConfigurationProfile{
				Location:   location.Normalize(model.Location),
				Properties: &configurationprofiles.ConfigurationProfileProperties{},
				Tags:       pointer.To(model.Tags),
			}

			properties.Properties.Configuration = expandConfigurationProfile(model)

			// NOTE: ordering
			if _, err := client.CreateOrUpdate(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			return nil
		},
	}
}

func expandConfigurationProfile(model ConfigurationModel) *any {
	// building configuration profile in json format
	jsonConfig := make(map[string]any)

	if len(model.Antimalware) > 0 {
		antimalwareConfig := model.Antimalware[0]
		jsonConfig["Antimalware/Enable"] = true
		jsonConfig["Antimalware/EnableRealTimeProtection"] = antimalwareConfig.RealTimeProtectionEnabled
		jsonConfig["Antimalware/RunScheduledScan"] = antimalwareConfig.ScheduledScanEnabled
		jsonConfig["Antimalware/ScanType"] = antimalwareConfig.ScanType
		jsonConfig["Antimalware/ScanDay"] = antimalwareConfig.ScanDay
		jsonConfig["Antimalware/ScanTimeInMinutes"] = antimalwareConfig.ScanTimeInMinutes
		if len(antimalwareConfig.Exclusions) > 0 {
			jsonConfig["Antimalware/Exclusions/Extensions"] = antimalwareConfig.Exclusions[0].Extensions
			jsonConfig["Antimalware/Exclusions/Paths"] = antimalwareConfig.Exclusions[0].Paths
			jsonConfig["Antimalware/Exclusions/Processes"] = antimalwareConfig.Exclusions[0].Processes
		}
	}

	if len(model.AzureSecurityBaseline) > 0 {
		azureSecurityBaselineConfig := model.AzureSecurityBaseline[0]
		jsonConfig["AzureSecurityBaseline/Enable"] = true
		jsonConfig["AzureSecurityBaseline/AssignmentType"] = azureSecurityBaselineConfig.AssignmentType
	}

	if len(model.Backup) > 0 {
		backupConfig := model.Backup[0]
		jsonConfig["Backup/Enable"] = true
		if backupConfig.PolicyName != "" {
			jsonConfig["Backup/PolicyName"] = backupConfig.PolicyName
		}
		jsonConfig["Backup/TimeZone"] = backupConfig.TimeZone
		jsonConfig["Backup/InstantRpRetentionRangeInDays"] = backupConfig.InstantRpRetentionRangeInDays
		if len(backupConfig.SchedulePolicy) > 0 {
			schedulePolicyConfig := backupConfig.SchedulePolicy[0]
			jsonConfig["Backup/SchedulePolicy/ScheduleRunFrequency"] = schedulePolicyConfig.ScheduleRunFrequency
			if len(schedulePolicyConfig.ScheduleRunTimes) > 0 {
				jsonConfig["Backup/SchedulePolicy/ScheduleRunTimes"] = schedulePolicyConfig.ScheduleRunTimes
			}
			if len(schedulePolicyConfig.ScheduleRunDays) > 0 {
				jsonConfig["Backup/SchedulePolicy/ScheduleRunDays"] = schedulePolicyConfig.ScheduleRunDays
			}
			jsonConfig["Backup/SchedulePolicy/SchedulePolicyType"] = schedulePolicyConfig.SchedulePolicyType
		}

		if len(backupConfig.RetentionPolicy) > 0 {
			retentionPolicyConfig := backupConfig.RetentionPolicy[0]
			jsonConfig["Backup/RetentionPolicy/RetentionPolicyType"] = retentionPolicyConfig.RetentionPolicyType
			if len(retentionPolicyConfig.DailySchedule) > 0 {
				dailyScheduleConfig := retentionPolicyConfig.DailySchedule[0]
				if len(dailyScheduleConfig.RetentionTimes) > 0 {
					jsonConfig["Backup/RetentionPolicy/DailySchedule/RetentionTimes"] = dailyScheduleConfig.RetentionTimes
				}

				if len(dailyScheduleConfig.RetentionDuration) > 0 {
					jsonConfig["Backup/RetentionPolicy/DailySchedule/RetentionDuration/Count"] = dailyScheduleConfig.RetentionDuration[0].Count
					jsonConfig["Backup/RetentionPolicy/DailySchedule/RetentionDuration/DurationType"] = dailyScheduleConfig.RetentionDuration[0].DurationType
				}
			}

			if len(retentionPolicyConfig.WeeklySchedule) > 0 {
				weeklyScheduleConfig := retentionPolicyConfig.WeeklySchedule[0]
				if len(weeklyScheduleConfig.RetentionTimes) > 0 {
					jsonConfig["Backup/RetentionPolicy/WeeklySchedule/RetentionTimes"] = weeklyScheduleConfig.RetentionTimes
				}

				if len(weeklyScheduleConfig.RetentionDuration) > 0 {
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

	var out any = jsonConfig
	return &out
}
