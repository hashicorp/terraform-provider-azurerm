package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ConfigurationV0ToV1{}

type ConfigurationV0ToV1 struct {
}

func (ConfigurationV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

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
					},
					"scheduled_scan_day": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  8,
					},
					"scheduled_scan_time_in_minutes": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  0,
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
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"time_zone": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "UTC",
					},
					"instant_rp_retention_range_in_days": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  5,
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
								},
								"schedule_run_times": {
									Type:     pluginsdk.TypeList,
									Optional: true,

									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"schedule_run_days": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"schedule_policy_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "SimpleSchedulePolicy",
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
												},
											},
											"retention_duration": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"count": {
															Type:     pluginsdk.TypeInt,
															Optional: true,
														},
														"duration_type": {
															Type:     pluginsdk.TypeString,
															Optional: true,
															Default:  "Days",
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
												},
											},
											"retention_duration": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"count": {
															Type:     pluginsdk.TypeInt,
															Optional: true,
														},
														"duration_type": {
															Type:     pluginsdk.TypeString,
															Optional: true,
															Default:  "Weeks",
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

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (c ConfigurationV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)

		parsed, err := configurationprofiles.ParseConfigurationProfileIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}
		newId := parsed.ID()

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId
		return rawState, nil
	}
}
