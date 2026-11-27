// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_vm

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceBackupProtectionPolicyVMSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Backup Policy.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z][-_!a-zA-Z0-9]{2,149}$"),
				"Backup Policy name must be 3 - 150 characters long, start with a letter, contain only letters and numbers.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"instant_restore_retention_days": {
			Description:  "The instant restore retention range in days.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(1, 30),
		},

		"recovery_vault_name": {
			Description:  "The name of the Recovery Services Vault to use.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		// `tiering_policy` is defined as the map type in Swagger and currently the key only supports `ArchivedRP`. Service team confirmed that they would support other keys in the future.
		"tiering_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"archived_restore_point": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"mode": {
									Description: "The tiering mode to control automatic tiering of recovery points.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ValidateFunc: validation.StringInSlice([]string{
										string(protectionpolicies.TieringModeTierAfter),
										string(protectionpolicies.TieringModeTierRecommended),
									}, false),
								},

								"duration": {
									Description:  "The number of days/weeks/months/years to retain backups in current tier before tiering.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntAtLeast(3),
								},

								"duration_type": {
									Description: "The retention duration type.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
									ValidateFunc: validation.StringInSlice([]string{
										string(protectionpolicies.RetentionDurationTypeDays),
										string(protectionpolicies.RetentionDurationTypeWeeks),
										string(protectionpolicies.RetentionDurationTypeMonths),
										string(protectionpolicies.RetentionDurationTypeYears),
									}, false),
								},
							},
						},
					},
				},
			},
		},

		"timezone": {
			Description: "The timezone.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "UTC",
		},

		"instant_restore_resource_group": {
			Description: "The instant restore Resource Group name.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"prefix": {
						Description:  "The prefix for the `instant_restore_resource_group` name.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"suffix": {
						Description:  "The suffix for the `instant_restore_resource_group` name.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"backup": {
			Description: "Configures the Policy backup frequency, times & days.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"frequency": {
						Description: "Sets the backup frequency.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(protectionpolicies.ScheduleRunTypeHourly),
							string(protectionpolicies.ScheduleRunTypeDaily),
							string(protectionpolicies.ScheduleRunTypeWeekly),
						}, false),
					},

					"time": {
						Description: "The time of day to perform the backup in 24hour format.",
						// applies to all backup schedules & retention times (they all must be the same)
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"), // time must be on the hour or half past
							"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
						),
					},

					"weekdays": {
						Description: "The days of the week to perform backups on. This is used when `frequency` is `Weekly`.",
						// only for weekly
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Set:      set.HashStringIgnoreCase,
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							ValidateFunc:     validation.IsDayOfTheWeek(true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},

					"hour_interval": {
						Description: "Interval in hour at which backup is triggered. This is used when `frequency` is `Hourly`.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						ValidateFunc: validation.IntInSlice([]int{
							4,
							6,
							8,
							12,
						}),
					},

					"hour_duration": {
						Description:  "Duration of the backup window in hours.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(4, 24),
					},
				},
			},
		},

		"consistency_type": {
			Description:  "The consistency type for the backup policy.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(protectionpolicies.PossibleValuesForIaasVMSnapshotConsistencyType(), false),
		},

		"policy_type": {
			Description: "Type of the Backup Policy.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     string(protectionpolicies.IAASVMPolicyTypeVOne),
			ValidateFunc: validation.StringInSlice([]string{
				string(protectionpolicies.IAASVMPolicyTypeVOne),
				string(protectionpolicies.IAASVMPolicyTypeVTwo),
			}, false),
		},

		"retention_daily": {
			Description: "Configures the policy daily retention.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Description:  "The number of daily backups to keep.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 9999), // Azure no longer supports less than 7 daily backups. This should be updated in 3.0 provider

					},
				},
			},
		},

		"retention_weekly": {
			Description: "Configures the policy weekly retention.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Description:  "The number of weekly backups to keep.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 9999),
					},

					"weekdays": {
						Description: "The weekday backups to retain.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						Set:         set.HashStringIgnoreCase,
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							ValidateFunc:     validation.IsDayOfTheWeek(true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},
				},
			},
		},

		"retention_monthly": {
			Description: "Configures the policy monthly retention.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Description:  "The number of monthly backups to keep.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 9999),
					},

					"weeks": {
						Description: "The weeks of the month to retain backups of.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Set:         set.HashStringIgnoreCase,
						RequiredWith: []string{
							"retention_monthly.0.weekdays",
						},
						ConflictsWith: []string{
							"retention_monthly.0.days",
							"retention_monthly.0.include_last_days",
						},
						AtLeastOneOf: []string{
							"retention_monthly.0.weekdays",
							"retention_monthly.0.weeks",
							"retention_monthly.0.days",
							"retention_monthly.0.include_last_days",
						},
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								string(protectionpolicies.WeekOfMonthFirst),
								string(protectionpolicies.WeekOfMonthSecond),
								string(protectionpolicies.WeekOfMonthThird),
								string(protectionpolicies.WeekOfMonthFourth),
								string(protectionpolicies.WeekOfMonthLast),
							}, false),
						},
					},

					"weekdays": {
						Description: "The weekday backups to retain.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Set:         set.HashStringIgnoreCase,
						RequiredWith: []string{
							"retention_monthly.0.weeks",
						},
						ConflictsWith: []string{
							"retention_monthly.0.days",
							"retention_monthly.0.include_last_days",
						},
						AtLeastOneOf: []string{
							"retention_monthly.0.weekdays",
							"retention_monthly.0.weeks",
							"retention_monthly.0.days",
							"retention_monthly.0.include_last_days",
						},
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							ValidateFunc:     validation.IsDayOfTheWeek(true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},

					"days": {
						Description: "The days of the month to retain backups of.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						ConflictsWith: []string{
							"retention_monthly.0.weeks",
							"retention_monthly.0.weekdays",
						},
						AtLeastOneOf: []string{
							"retention_monthly.0.weekdays",
							"retention_monthly.0.weeks",
							"retention_monthly.0.days",
							"retention_monthly.0.include_last_days",
						},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validation.IntBetween(1, 31), // days in months
						},
					},

					"include_last_days": {
						Description: "Including the last day of the month.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
						ConflictsWith: []string{
							"retention_monthly.0.weeks",
							"retention_monthly.0.weekdays",
						},
						AtLeastOneOf: []string{
							"retention_monthly.0.weekdays",
							"retention_monthly.0.weeks",
							"retention_monthly.0.days",
							"retention_monthly.0.include_last_days",
						},
					},
				},
			},
		},

		"retention_yearly": {
			Description: "Configures the policy yearly retention.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Description:  "The number of yearly backups to keep.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 9999),
					},

					"months": {
						Description: "The months of the year to retain backups of.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						Set:         set.HashStringIgnoreCase,
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							ValidateFunc:     validation.IsMonth(true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},

					"weeks": {
						Description: "The weeks of the month to retain backups of.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Set:         set.HashStringIgnoreCase,
						RequiredWith: []string{
							"retention_yearly.0.weekdays",
						},
						ConflictsWith: []string{
							"retention_yearly.0.days",
							"retention_yearly.0.include_last_days",
						},
						AtLeastOneOf: []string{
							"retention_yearly.0.weeks",
							"retention_yearly.0.weekdays",
							"retention_yearly.0.days",
							"retention_yearly.0.include_last_days",
						},
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								string(protectionpolicies.WeekOfMonthFirst),
								string(protectionpolicies.WeekOfMonthSecond),
								string(protectionpolicies.WeekOfMonthThird),
								string(protectionpolicies.WeekOfMonthFourth),
								string(protectionpolicies.WeekOfMonthLast),
							}, false),
						},
					},

					"weekdays": {
						Description: "The weekday backups to retain.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Set:         set.HashStringIgnoreCase,
						RequiredWith: []string{
							"retention_yearly.0.weeks",
						},
						ConflictsWith: []string{
							"retention_yearly.0.days",
							"retention_yearly.0.include_last_days",
						},
						AtLeastOneOf: []string{
							"retention_yearly.0.weeks",
							"retention_yearly.0.weekdays",
							"retention_yearly.0.days",
							"retention_yearly.0.include_last_days",
						},
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							ValidateFunc:     validation.IsDayOfTheWeek(true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},

					"days": {
						Description: "The days of the month to retain backups of.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						ConflictsWith: []string{
							"retention_yearly.0.weeks",
							"retention_yearly.0.weekdays",
						},
						AtLeastOneOf: []string{
							"retention_yearly.0.weeks",
							"retention_yearly.0.weekdays",
							"retention_yearly.0.days",
							"retention_yearly.0.include_last_days",
						},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validation.IntBetween(1, 31), // days in months
						},
					},

					"include_last_days": {
						Description: "Including the last day of the month.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
						ConflictsWith: []string{
							"retention_yearly.0.weeks",
							"retention_yearly.0.weekdays",
						},
						AtLeastOneOf: []string{
							"retention_yearly.0.weeks",
							"retention_yearly.0.weekdays",
							"retention_yearly.0.days",
							"retention_yearly.0.include_last_days",
						},
					},
				},
			},
		},
	}
}
