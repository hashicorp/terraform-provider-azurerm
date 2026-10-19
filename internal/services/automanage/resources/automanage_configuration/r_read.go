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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/resources/automanage_configuration/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationProfilesClient

			id, err := configurationprofiles.ParseConfigurationProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ConfigurationModel{
				Name:              id.ConfigurationProfileName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if props := model.Properties; props != nil && props.Configuration != nil {
					configMap := (*props.Configuration).(map[string]any)

					state.Antimalware = flattenAntiMalwareConfig(configMap)

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
				state.Tags = pointer.From(model.Tags)
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}
			return metadata.Encode(&state)
		},
	}
}

type AzureSecurityBaselineConfiguration struct {
	AssignmentType string `tfschema:"assignment_type"`
}

type AntimalwareConfiguration struct {
	Exclusions                []AntimalwareExclusions `tfschema:"exclusions"`
	RealTimeProtectionEnabled bool                    `tfschema:"real_time_protection_enabled"`
	ScheduledScanEnabled      bool                    `tfschema:"scheduled_scan_enabled"`
	ScanType                  string                  `tfschema:"scheduled_scan_type"`
	ScanDay                   int64                   `tfschema:"scheduled_scan_day"`
	ScanTimeInMinutes         int64                   `tfschema:"scheduled_scan_time_in_minutes"`
}

type AntimalwareExclusions struct {
	Extensions string `tfschema:"extensions"`
	Paths      string `tfschema:"paths"`
	Processes  string `tfschema:"processes"`
}

type BackupConfiguration struct {
	PolicyName                    string                         `tfschema:"policy_name"`
	TimeZone                      string                         `tfschema:"time_zone"`
	InstantRpRetentionRangeInDays int64                          `tfschema:"instant_rp_retention_range_in_days"`
	SchedulePolicy                []SchedulePolicyConfiguration  `tfschema:"schedule_policy"`
	RetentionPolicy               []RetentionPolicyConfiguration `tfschema:"retention_policy"`
}

type ScheduleConfiguration struct {
	RetentionTimes    []string                         `tfschema:"retention_times"`
	RetentionDuration []RetentionDurationConfiguration `tfschema:"retention_duration"`
}

type RetentionPolicyConfiguration struct {
	RetentionPolicyType string                  `tfschema:"retention_policy_type"`
	DailySchedule       []ScheduleConfiguration `tfschema:"daily_schedule"`
	WeeklySchedule      []ScheduleConfiguration `tfschema:"weekly_schedule"`
}

type RetentionDurationConfiguration struct {
	Count        int64  `tfschema:"count"`
	DurationType string `tfschema:"duration_type"`
}

type SchedulePolicyConfiguration struct {
	ScheduleRunFrequency string   `tfschema:"schedule_run_frequency"`
	ScheduleRunTimes     []string `tfschema:"schedule_run_times"`
	ScheduleRunDays      []string `tfschema:"schedule_run_days"`
	SchedulePolicyType   string   `tfschema:"schedule_policy_type"`
}

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.ConfigurationV0ToV1{},
		},
	}
}
