// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Schedule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ScheduleName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"automation_account_name": {
			Description:  "The name of the Automation Account in which the Schedule is created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutomationAccount(),
		},

		"frequency": {
			Description:  "The frequency of the schedule. - can be either `OneTime`, `Minute`, `Day`, `Hour`, `Week`, or `Month`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(schedule.PossibleValuesForScheduleFrequency(), false),
		},

		"interval": {
			Description: "The number of `frequency`s between runs. Only valid when frequency is `Minute`, `Day`, `Hour`, `Week`, or `Month`.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// NOTE: O+C this is set to `1` unless `frequency` is `OneTime` (in which case this property is ignored) o+c can remain since this can be updated without issue
			Computed:     true,
			ValidateFunc: validation.IntBetween(1, 100),
		},

		"start_time": {
			Description: "Start time of the schedule.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C We set this to now + 7 minutes if omitted so this should remain Computed
			Computed:         true,
			DiffSuppressFunc: suppress.RFC3339MinuteTime,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"expiry_time": {
			Description: "The end time of the schedule.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C when frequency is OneTime this has ridiculous value when recurring: "9999-12-31T15:59:00-08:00" which can remain as it can be updated without issue
			Computed:         true,
			DiffSuppressFunc: suppress.RFC3339MinuteTime,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"description": {
			Description: "A description for this Schedule.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"timezone": {
			Description:  "The timezone of the start time. For possible values see: <https://learn.microsoft.com/rest/api/maps/timezone/gettimezoneenumwindows>.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "Etc/UTC",
			ValidateFunc: validate.AzureTimeZoneString(),
		},

		"week_days": {
			Description: "List of days of the week that the job should execute on. Only valid when frequency is `Week`.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(schedule.PossibleValuesForScheduleDay(), false),
			},
			Set:           set.HashStringIgnoreCase,
			ConflictsWith: []string{"month_days", "monthly_occurrence"},
		},

		"month_days": {
			Description: "List of days of the month that the job should execute on. `-1` for last day of the month. Only valid when frequency is `Month`.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
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
			Description: "Only valid when frequency is `Month`.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"day": {
						Description:  "Day of the occurrence.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(schedule.PossibleValuesForScheduleDay(), false),
					},
					"occurrence": {
						Description: "Occurrence of the week within the month. `-1` for last week within the month.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
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
