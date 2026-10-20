// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_schedule

import (
	// import time/tzdata to embed timezone information in the program
	// add this to resolve https://github.com/hashicorp/terraform-provider-azurerm/issues/20690
	_ "time/tzdata"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/schedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func automationScheduleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ScheduleName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"automation_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutomationAccount(),
		},

		"frequency": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(schedule.ScheduleFrequencyDay),
				string(schedule.ScheduleFrequencyHour),
				string(schedule.ScheduleFrequencyMonth),
				string(schedule.ScheduleFrequencyOneTime),
				string(schedule.ScheduleFrequencyWeek),
			}, false),
		},

		"interval": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C this is set to `1` unless `frequency` is `OneTime` (in which case this property is ignored) o+c can remain since this can be updated without issue
			Computed:     true,
			ValidateFunc: validation.IntBetween(1, 100),
		},

		"start_time": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C We set this to now + 7 minutes if omitted so this should remain Computed
			Computed:         true,
			DiffSuppressFunc: suppress.RFC3339MinuteTime,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"expiry_time": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C when frequency is OneTime this has ridiculous value when recurring: "9999-12-31T15:59:00-08:00" which can remain as it can be updated without issue
			Computed:         true,
			DiffSuppressFunc: suppress.RFC3339MinuteTime,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"timezone": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "Etc/UTC",
			ValidateFunc: validate.AzureTimeZoneString(),
		},

		"week_days": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(schedule.PossibleValuesForScheduleDay(), false),
			},
			Set:           set.HashStringIgnoreCase,
			ConflictsWith: []string{"month_days", "monthly_occurrence"},
		},

		"month_days": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeInt,
				ValidateFunc: validation.All(
					validation.IntBetween(-1, 31),
					validation.IntNotInSlice([]int{0}),
				),
			},
			Set:           set.HashInt,
			ConflictsWith: []string{"week_days", "monthly_occurrence"},
		},

		"monthly_occurrence": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"day": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(schedule.PossibleValuesForScheduleDay(), false),
					},
					"occurrence": {
						Type:     pluginsdk.TypeInt,
						Required: true,
						ValidateFunc: validation.All(
							validation.IntBetween(-1, 5),
							validation.IntNotInSlice([]int{0}),
						),
					},
				},
			},
			ConflictsWith: []string{"week_days", "month_days"},
		},
	}
}
