// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_file_share

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceBackupProtectionPolicyFileShareSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z][-_!a-zA-Z0-9]{2,149}$"),
				"Backup Policy name must be 3 - 150 characters long, start with a letter, contain only letters and numbers.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"recovery_vault_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		"timezone": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "UTC",
		},

		"backup": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"frequency": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc: validation.StringInSlice([]string{
							string(protectionpolicies.ScheduleRunTypeDaily),
							string(protectionpolicies.ScheduleRunTypeHourly),
						}, false),
					},

					"time": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"), // time must be on the hour or half past
							"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
						),
						AtLeastOneOf: []string{
							"backup.0.time",
							"backup.0.hourly",
						},
						ConflictsWith: []string{
							"backup.0.hourly",
						},
					},

					"hourly": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						AtLeastOneOf: []string{
							"backup.0.time",
							"backup.0.hourly",
						},
						ConflictsWith: []string{
							"backup.0.time",
						},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"interval": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntInSlice([]int{1, 2, 4, 6, 8, 12}),
								},

								"start_time": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"), // time must be on the hour or half past
										"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
									),
								},

								"window_duration": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(4, 24),
								},
							},
						},
					},
				},
			},
		},

		"retention_daily": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 200),
					},
				},
			},
		},

		"retention_weekly": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 200),
					},

					"weekdays": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Set:      set.HashStringIgnoreCase,
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.IsDayOfTheWeek(true),
						},
					},
				},
			},
		},

		"retention_monthly": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 120),
					},

					"weeks": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Set:      set.HashStringIgnoreCase,
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
						RequiredWith: []string{
							"retention_monthly.0.weekdays",
						},
					},

					"weekdays": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Set:      set.HashStringIgnoreCase,
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.IsDayOfTheWeek(true),
						},
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
					},

					"days": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validation.IntBetween(1, 31), // days in months
						},
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

					"include_last_days": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
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
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 10),
					},

					"months": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Set:      set.HashStringIgnoreCase,
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.IsMonth(true),
						},
					},

					"weeks": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Set:      set.HashStringIgnoreCase,
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
					},

					"weekdays": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Set:      set.HashStringIgnoreCase,
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.IsDayOfTheWeek(true),
						},
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
					},

					"days": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validation.IntBetween(1, 31), // days in months
						},
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

					"include_last_days": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
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

		"backup_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "snapshot",
			ValidateFunc: validation.StringInSlice([]string{
				"snapshot",
				"vault-standard",
			}, false),
		},

		"snapshot_retention_in_days": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  0,
		},
	}
}
