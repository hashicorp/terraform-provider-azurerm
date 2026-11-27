// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_trigger_recurrence

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func logicAppTriggerRecurrenceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Recurrence Triggers to be created within the Logic App Workflow.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"logic_app_id": {
			Description:  "The ID of the Logic App Workflow.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workflows.ValidateWorkflowID,
		},

		"frequency": {
			Description: "The Frequency at which this Trigger should be run.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"Month",
				"Week",
				"Day",
				"Hour",
				"Minute",
				"Second",
			}, false),
		},

		"interval": {
			Description: "The interval used for the Frequency, for example a value of `4` for `interval` and `hour` for `frequency` would run the Trigger every 4 hours.",
			Type:        pluginsdk.TypeInt,
			Required:    true,
		},

		"start_time": {
			Description:  "The start date and time for this trigger in RFC3339 format: `2000-01-02T03:04:05Z`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"schedule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"at_these_hours": {
						Description:  "A list of hours when the trigger should run.",
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"schedule.0.at_these_hours", "schedule.0.at_these_minutes", "schedule.0.on_these_days"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validation.IntBetween(0, 23),
						},
					},
					"at_these_minutes": {
						Description:  "A list of minutes when the trigger should run.",
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"schedule.0.at_these_hours", "schedule.0.at_these_minutes", "schedule.0.on_these_days"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validation.IntBetween(0, 59),
						},
					},
					"on_these_days": {
						Description:  "A list of days when the trigger should run.",
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"schedule.0.at_these_hours", "schedule.0.at_these_minutes", "schedule.0.on_these_days"},
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"Monday",
								"Tuesday",
								"Wednesday",
								"Thursday",
								"Friday",
								"Saturday",
								"Sunday",
							}, false),
						},
					},
				},
			},
		},

		"time_zone": {
			Description:  "The time zone for this trigger. Supported time zone options are listed [here](https://support.microsoft.com/en-us/help/973627/microsoft-time-zone-index-values).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.TriggerRecurrenceTimeZone(),
		},
	}
}
