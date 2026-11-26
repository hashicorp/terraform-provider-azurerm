// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_start_stop_schedule

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/startstopmanagedinstanceschedules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_instance_id": {
			Description:  "The ID of the Managed Instance.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSqlManagedInstanceID,
		},

		"description": {
			Description: "The description of the schedule.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"schedule": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"start_day": {
						Description:  "Start day of the schedule.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(startstopmanagedinstanceschedules.PossibleValuesForDayOfWeek(), false),
					},

					"start_time": {
						Description:  "Start time of the schedule in 24-hour format (e.g., `08:00`).",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"stop_day": {
						Description:  "Stop day of the schedule.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(startstopmanagedinstanceschedules.PossibleValuesForDayOfWeek(), false),
					},

					"stop_time": {
						Description:  "Stop time of the schedule in 24-hour format (e.g., `17:00`).",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"timezone_id": {
			Description:  "The time zone of the schedule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "UTC",
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"next_execution_time": {
			Description: "Timestamp when the next action will be executed in the corresponding schedule time zone.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"next_run_action": {
			Description: "Next action to be executed (Start or Stop).",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
