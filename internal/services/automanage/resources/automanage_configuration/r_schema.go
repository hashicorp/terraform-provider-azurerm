// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automanage_configuration

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
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

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
