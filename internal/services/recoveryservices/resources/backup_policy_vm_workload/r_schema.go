// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_vm_workload

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

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.BackupPolicyName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"recovery_vault_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		"protection_policy": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"policy_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(protectionpolicies.PolicyTypeDifferential),
							string(protectionpolicies.PolicyTypeFull),
							string(protectionpolicies.PolicyTypeIncremental),
							string(protectionpolicies.PolicyTypeLog),
						}, false),
					},

					"backup": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"frequency": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(protectionpolicies.ScheduleRunTypeDaily),
										string(protectionpolicies.ScheduleRunTypeWeekly),
									}, false),
								},

								"frequency_in_minutes": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									ValidateFunc: validation.IntInSlice([]int{
										15,
										30,
										60,
										120,
										240,
										480,
										720,
										1440,
									}),
								},

								"time": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"),
										"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
									),
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
								},
							},
						},
					},

					"retention_daily": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"count": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(7, 9999),
								},
							},
						},
					},

					"retention_weekly": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"count": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 5163),
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
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"count": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 1188),
								},

								"format_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(protectionpolicies.RetentionScheduleFormatDaily),
										string(protectionpolicies.RetentionScheduleFormatWeekly),
									}, false),
								},

								"monthdays": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeInt,
										ValidateFunc: validation.IntBetween(0, 28),
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
								},
							},
						},
					},

					"retention_yearly": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"count": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 99),
								},

								"format_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(protectionpolicies.RetentionScheduleFormatDaily),
										string(protectionpolicies.RetentionScheduleFormatWeekly),
									}, false),
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

								"monthdays": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeInt,
										ValidateFunc: validation.IntBetween(0, 28),
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
								},
							},
						},
					},

					"simple_retention": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"count": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(7, 35),
								},
							},
						},
					},
				},
			},
		},

		"settings": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"time_zone": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"compression_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"workload_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(protectionpolicies.WorkloadTypeSQLDataBase),
				string(protectionpolicies.WorkloadTypeSAPHanaDatabase),
			}, false),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
