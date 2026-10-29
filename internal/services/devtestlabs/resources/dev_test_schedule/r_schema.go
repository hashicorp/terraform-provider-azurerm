// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_schedule

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/schedules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	devTestValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func devTestScheduleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		// There's a bug in the Azure API where this is returned in lower-case
		// BUG: https://github.com/Azure/azure-rest-api-specs/issues/3964
		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"lab_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: devTestValidate.DevTestLabName(),
		},

		"status": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      schedules.EnableStatusDisabled,
			ValidateFunc: validation.StringInSlice(schedules.PossibleValuesForEnableStatus(), false),
		},

		"task_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"weekly_recurrence": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"time": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]|[0-9])[0-5][0-9]$"),
							"Time of day must match the format HHmm where HH is 00-23 and mm is 00-59",
						),
					},

					"week_days": {
						Type:     pluginsdk.TypeList,
						Optional: true,
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

		"daily_recurrence": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"time": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]|[0-9])[0-5][0-9]$"),
							"Time of day must match the format HHmm where HH is 00-23 and mm is 00-59",
						),
					},
				},
			},
		},

		"hourly_recurrence": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"minute": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 59),
					},
				},
			},
		},

		"time_zone_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: computeValidate.VirtualMachineTimeZoneCaseInsensitive(),
		},

		"notification_settings": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"status": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      schedules.EnableStatusDisabled,
						ValidateFunc: validation.StringInSlice(schedules.PossibleValuesForEnableStatus(), false),
					},
					"time_in_minutes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},
					"webhook_url": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
