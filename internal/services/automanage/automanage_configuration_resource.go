// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automanage

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

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

type ConfigurationModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`

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

type AutoManageConfigurationResource struct{}

var _ sdk.ResourceWithUpdate = AutoManageConfigurationResource{}
var _ sdk.ResourceWithStateMigration = AutoManageConfigurationResource{}

func (r AutoManageConfigurationResource) ResourceType() string {
	return "azurerm_automanage_configuration"
}

func (r AutoManageConfigurationResource) ModelObject() interface{} {
	return &ConfigurationModel{}
}

func (r AutoManageConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return configurationprofiles.ValidateConfigurationProfileID
}

func (r AutoManageConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		// "Antimalware/Enable": boolean, true if block exists
		// "Antimalware/EnableRealTimeProtection": boolean,
		// "Antimalware/RunScheduledScan": boolean,
		// "Antimalware/ScanType": string ("Quick", "Full"),
		// "Antimalware/ScanDay": int (0-8) Ex: 0 - daily, 1 - Sunday, 2 - Monday, .... 7 - Saturday, 8 - Disabled,
		// "Antimalware/ScanTimeInMinutes": int (0 - 1440),
		// "Antimalware/Exclusions/Extensions": string (extensions separated by semicolon. Ex: ".ext1;.ext2"),
		// "Antimalware/Exclusions/Paths": string (Paths separated by semicolon. Ex: "c:\excluded-path-1;c:\excluded-path-2"),
		// "Antimalware/Exclusions/Processes": string (Processes separated by semicolon. Ex: "proc1.exe;proc2.exe"),
		"antimalware": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"real_time_protection_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"scheduled_scan_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"scheduled_scan_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "Quick",
						ValidateFunc: validation.StringInSlice([]string{
							"Quick",
							"Full",
						}, false),
					},
					"scheduled_scan_day": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  8,
						ValidateFunc: validation.IntInSlice([]int{
							0, 1, 2, 3, 4, 5, 6, 7, 8,
						}),
					},
					"scheduled_scan_time_in_minutes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntBetween(0, 1439),
					},
					"exclusions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"extensions": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"paths": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"processes": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},

		// "AzureSecurityBaseline/Enable": boolean, true if block exists
		// "AzureSecurityBaseline/AssignmentType": string ("ApplyAndAutoCorrect", "ApplyAndMonitor", "Audit", "DeployAndAutoCorrect"),
		"azure_security_baseline": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"assignment_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "ApplyAndAutoCorrect",
						ValidateFunc: validation.StringInSlice([]string{
							"ApplyAndAutoCorrect",
							"ApplyAndMonitor",
							"Audit",
							"DeployAndAutoCorrect",
						}, false),
					},
				},
			},
		},

		// "Backup/Enable": boolean, true if block exists
		// "Backup/PolicyName": string (length 3 - 150, begin with alphanumeric char, only contain alphanumeric chars and hyphens),
		// "Backup/TimeZone": timezone,
		// "Backup/InstantRpRetentionRangeInDays": int (1 - 5 if ScheduleRunFrequency is Daily, 5 if ScheduleRunFrequency is Weekly),
		// "Backup/SchedulePolicy/ScheduleRunFrequency": string ("Daily", "Weekly"),
		// "Backup/SchedulePolicy/ScheduleRunTimes": list of DateTime,
		// "Backup/SchedulePolicy/ScheduleRunDays": list of strings (["Sunday", "Monday", "Wednesday", "Thursday", "Friday", "Saturday"]),
		// "Backup/SchedulePolicy/SchedulePolicyType": string ("SimpleSchedulePolicy"),
		// "Backup/RetentionPolicy/RetentionPolicyType": string ("LongTermRetentionPolicy"),
		// "Backup/RetentionPolicy/DailySchedule/RetentionTimes": list of DateTime,
		// "Backup/RetentionPolicy/DailySchedule/RetentionDuration/Count": int (7 - 9999),
		// "Backup/RetentionPolicy/DailySchedule/RetentionDuration/DurationType": string ("Days"),
		// "Backup/RetentionPolicy/WeeklySchedule/RetentionTimes":, list of DateTime
		// "Backup/RetentionPolicy/WeeklySchedule/RetentionDuration/Count":, int (1 - 5163)
		// "Backup/RetentionPolicy/WeeklySchedule/RetentionDuration/DurationType": string ("Weeks"),
		"backup": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"policy_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{2,149}$`), "Policy name must be 3 - 150 characters long, begin with an alphanumeric character, and only contain alphanumeric characters and hyphens."),
					},
					"time_zone": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "UTC",
					},
					"instant_rp_retention_range_in_days": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      5,
						ValidateFunc: validation.IntBetween(1, 5),
					},
					"schedule_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"schedule_run_frequency": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "Daily",
									ValidateFunc: validation.StringInSlice([]string{
										"Daily",
										"Weekly",
									}, false),
								},
								"schedule_run_times": {
									Type:     pluginsdk.TypeList,
									Optional: true,

									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringMatch(
											regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"), // time must be on the hour or half past
											"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
										),
									},
								},
								"schedule_run_days": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"Sunday",
											"Monday",
											"Tuesday",
											"Wednesday",
											"Thursday",
											"Friday",
											"Saturday",
										}, false),
									},
								},
								"schedule_policy_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "SimpleSchedulePolicy",
									ValidateFunc: validation.StringInSlice([]string{
										"SimpleSchedulePolicy",
									}, false),
								},
							},
						},
					},
					"retention_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"retention_policy_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "LongTermRetentionPolicy",
									ValidateFunc: validation.StringInSlice([]string{
										"LongTermRetentionPolicy",
									}, false),
								},
								"daily_schedule": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"retention_times": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
													ValidateFunc: validation.StringMatch(
														regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"), // time must be on the hour or half past
														"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
													),
												},
											},
											"retention_duration": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"count": {
															Type:         pluginsdk.TypeInt,
															Optional:     true,
															ValidateFunc: validation.IntBetween(7, 9999),
														},
														"duration_type": {
															Type:     pluginsdk.TypeString,
															Optional: true,
															Default:  "Days",
															ValidateFunc: validation.StringInSlice([]string{
																"Days",
															}, false),
														},
													},
												},
											},
										},
									},
								},
								"weekly_schedule": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"retention_times": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
													ValidateFunc: validation.StringMatch(
														regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"), // time must be on the hour or half past
														"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
													),
												},
											},
											"retention_duration": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"count": {
															Type:         pluginsdk.TypeInt,
															Optional:     true,
															ValidateFunc: validation.IntBetween(1, 5163),
														},
														"duration_type": {
															Type:     pluginsdk.TypeString,
															Optional: true,
															Default:  "Weeks",
															ValidateFunc: validation.StringInSlice([]string{
																"Weeks",
															}, false),
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},

		// "AutomationAccount/Enable": boolean,
		"automation_account_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		// "BootDiagnostics/Enable": boolean,
		"boot_diagnostics_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		// "DefenderForCloud/Enable": boolean,
		"defender_for_cloud_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
		// "GuestConfiguration/Enable": boolean,
		"guest_configuration_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		// "LogAnalytics/Enable": boolean,
		"log_analytics_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		// "Alerts/AutomanageStatusChanges/Enable": boolean,
		"status_change_alert_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"tags": commonschema.Tags(),
	}
}

func (r AutoManageConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AutoManageConfigurationResource) Create() sdk.ResourceFunc {
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
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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
			return nil
		},
	}
}

func (r AutoManageConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationProfilesClient

			id, err := configurationprofiles.ParseConfigurationProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := configurationprofiles.ConfigurationProfile{
				Location: location.Normalize(model.Location),
				Properties: &configurationprofiles.ConfigurationProfileProperties{
					Configuration: expandConfigurationProfile(model),
				},
				Tags: pointer.To(model.Tags),
			}

			if _, err := client.CreateOrUpdate(ctx, *id, properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r AutoManageConfigurationResource) Read() sdk.ResourceFunc {
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
					configMap := (*props.Configuration).(map[string]interface{})

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

			return metadata.Encode(&state)
		},
	}
}

func (r AutoManageConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationProfilesClient

			id, err := configurationprofiles.ParseConfigurationProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r AutoManageConfigurationResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.ConfigurationV0ToV1{},
		},
	}
}
func expandConfigurationProfile(model ConfigurationModel) *interface{} {
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

	var out interface{} = jsonConfig
	return &out
}

func flattenAntiMalwareConfig(configMap map[string]interface{}) []AntimalwareConfiguration {
	if val, ok := configMap["Antimalware/Enable"]; !ok || (val == nil) {
		return nil
	}

	antimalware := make([]AntimalwareConfiguration, 1)
	antimalware[0] = AntimalwareConfiguration{}

	if val, ok := configMap["Antimalware/EnableRealTimeProtection"]; ok {
		antimalware[0].RealTimeProtectionEnabled = val.(bool)
	}

	if val, ok := configMap["Antimalware/RunScheduledScan"]; ok {
		antimalware[0].ScheduledScanEnabled = val.(bool)
	}

	if val, ok := configMap["Antimalware/ScanType"]; ok {
		antimalware[0].ScanType = val.(string)
	}

	if val, ok := configMap["Antimalware/ScanDay"]; ok {
		antimalware[0].ScanDay = int64(val.(float64))
	}

	if val, ok := configMap["Antimalware/ScanTimeInMinutes"]; ok {
		antimalware[0].ScanTimeInMinutes = int64(val.(float64))
	}

	exclusions := AntimalwareExclusions{}
	exclusionsChanged := false
	if val, ok := configMap["Antimalware/Exclusions/Extensions"]; ok {
		exclusions.Extensions = val.(string)
		exclusionsChanged = true
	}

	if val, ok := configMap["Antimalware/Exclusions/Paths"]; ok {
		exclusions.Paths = val.(string)
		exclusionsChanged = true
	}

	if val, ok := configMap["Antimalware/Exclusions/Processes"]; ok {
		exclusions.Processes = val.(string)
		exclusionsChanged = true
	}

	if exclusionsChanged {
		antimalware[0].Exclusions = make([]AntimalwareExclusions, 1)
		antimalware[0].Exclusions[0] = exclusions
	}

	return antimalware
}

func flattenAzureSecurityBaselineConfig(configMap map[string]interface{}) []AzureSecurityBaselineConfiguration {
	if val, ok := configMap["AzureSecurityBaseline/Enable"]; !ok || (val == nil) {
		return nil
	}

	azureSecurityBaseline := make([]AzureSecurityBaselineConfiguration, 1)
	azureSecurityBaseline[0] = AzureSecurityBaselineConfiguration{}

	if val, ok := configMap["AzureSecurityBaseline/AssignmentType"]; ok {
		azureSecurityBaseline[0].AssignmentType = val.(string)
	}

	return azureSecurityBaseline
}

func flattenBackupConfig(configMap map[string]interface{}) []BackupConfiguration {
	if val, ok := configMap["Backup/Enable"]; !ok || (val == nil) {
		return nil
	}

	backup := make([]BackupConfiguration, 1)
	backup[0] = BackupConfiguration{}

	if val, ok := configMap["Backup/PolicyName"]; ok {
		backup[0].PolicyName = val.(string)
	}

	if val, ok := configMap["Backup/TimeZone"]; ok {
		backup[0].TimeZone = val.(string)
	}

	if val, ok := configMap["Backup/InstantRpRetentionRangeInDays"]; ok {
		backup[0].InstantRpRetentionRangeInDays = int64(val.(float64))
	}

	schedulePolicy := SchedulePolicyConfiguration{}
	schedulePolicyChanged := false
	if val, ok := configMap["Backup/SchedulePolicy/ScheduleRunFrequency"]; ok {
		schedulePolicy.ScheduleRunFrequency = val.(string)
		schedulePolicyChanged = true
	}

	if val, ok := configMap["Backup/SchedulePolicy/ScheduleRunTimes"]; ok {
		schedulePolicy.ScheduleRunTimes = flattenToListOfString(val)
		schedulePolicyChanged = true
	}

	if val, ok := configMap["Backup/SchedulePolicy/ScheduleRunDays"]; ok {
		schedulePolicy.ScheduleRunDays = flattenToListOfString(val)
		schedulePolicyChanged = true
	}

	if val, ok := configMap["Backup/SchedulePolicy/SchedulePolicyType"]; ok {
		schedulePolicy.SchedulePolicyType = val.(string)
		schedulePolicyChanged = true
	}

	if schedulePolicyChanged {
		backup[0].SchedulePolicy = make([]SchedulePolicyConfiguration, 1)
		backup[0].SchedulePolicy[0] = schedulePolicy
	}

	retentionPolicy := RetentionPolicyConfiguration{}
	retentionPolicyChanged := false
	if val, ok := configMap["Backup/RetentionPolicy/RetentionPolicyType"]; ok {
		retentionPolicy.RetentionPolicyType = val.(string)
		retentionPolicyChanged = true
	}

	dailySchedule := ScheduleConfiguration{}
	dailyScheduleChanged := false
	if val, ok := configMap["Backup/RetentionPolicy/DailySchedule/RetentionTimes"]; ok {
		dailySchedule.RetentionTimes = flattenToListOfString(val)
		dailyScheduleChanged = true
	}

	retentionDuration := RetentionDurationConfiguration{}
	retentionDurationChanged := false
	if val, ok := configMap["Backup/RetentionPolicy/DailySchedule/RetentionDuration/Count"]; ok {
		retentionDuration.Count = int64(val.(float64))
		retentionDurationChanged = true
	}

	if val, ok := configMap["Backup/RetentionPolicy/DailySchedule/RetentionDuration/DurationType"]; ok {
		retentionDuration.DurationType = val.(string)
		retentionDurationChanged = true
	}

	if retentionDurationChanged {
		dailySchedule.RetentionDuration = make([]RetentionDurationConfiguration, 1)
		dailySchedule.RetentionDuration[0] = retentionDuration
	}

	if dailyScheduleChanged || retentionDurationChanged {
		retentionPolicy.DailySchedule = make([]ScheduleConfiguration, 1)
		retentionPolicy.DailySchedule[0] = dailySchedule
	}

	weeklySchedule := ScheduleConfiguration{}
	weeklyScheduleChanged := false
	if val, ok := configMap["Backup/RetentionPolicy/WeeklySchedule/RetentionTimes"]; ok {
		weeklySchedule.RetentionTimes = flattenToListOfString(val)
		weeklyScheduleChanged = true
	}

	weeklyRetentionDuration := RetentionDurationConfiguration{}
	weeklyRetentionDurationChanged := false
	if val, ok := configMap["Backup/RetentionPolicy/WeeklySchedule/RetentionDuration/Count"]; ok {
		weeklyRetentionDuration.Count = int64(val.(float64))
		weeklyRetentionDurationChanged = true
	}

	if val, ok := configMap["Backup/RetentionPolicy/WeeklySchedule/RetentionDuration/DurationType"]; ok {
		weeklyRetentionDuration.DurationType = val.(string)
		weeklyRetentionDurationChanged = true
	}

	if weeklyRetentionDurationChanged {
		weeklySchedule.RetentionDuration = make([]RetentionDurationConfiguration, 1)
		weeklySchedule.RetentionDuration[0] = weeklyRetentionDuration
	}

	if weeklyScheduleChanged || weeklyRetentionDurationChanged {
		retentionPolicy.WeeklySchedule = make([]ScheduleConfiguration, 1)
		retentionPolicy.WeeklySchedule[0] = weeklySchedule
	}

	if retentionPolicyChanged || dailyScheduleChanged || retentionDurationChanged || weeklyScheduleChanged || weeklyRetentionDurationChanged {
		backup[0].RetentionPolicy = make([]RetentionPolicyConfiguration, 1)
		backup[0].RetentionPolicy[0] = retentionPolicy
	}

	return backup
}

func flattenToListOfString(val interface{}) []string {
	lis := val.([]interface{})
	strs := make([]string, len(lis))
	for i, v := range lis {
		strs[i] = v.(string)
	}
	return strs
}
