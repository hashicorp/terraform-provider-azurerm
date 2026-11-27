// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_schedule

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/schedules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	devTestValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func devTestScheduleSchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Dev Test Lab schedule. For instance for task_type `LabVmsStartupTask` the name needs to be `LabVmAutoStart`.",
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
			Description:  "The name of the Dev Test Lab.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: devTestValidate.DevTestLabName(),
		},

		"enabled": {
			Description: "Whether this schedule is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"task_type": {
			Description: "The task type of the schedule.",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},

		"weekly_recurrence": {
			Description: "The properties of a weekly schedule. If the schedule occurs only some days of the week, specify the weekly recurrence.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"time": {
						Description: "The time when the schedule takes effect.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]|[0-9])[0-5][0-9]$"),
							"Time of day must match the format HHmm where HH is 00-23 and mm is 00-59",
						),
					},

					"week_days": {
						Description: "A list of days that this schedule takes effect.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
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
			Description: "The properties of a daily schedule. If the schedule occurs once each day of the week, specify the daily recurrence.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"time": {
						Description: "The time each day when the schedule takes effect.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]|[0-9])[0-5][0-9]$"),
							"Time of day must match the format HHmm where HH is 00-23 and mm is 00-59",
						),
					},
				},
			},
		},

		"hourly_recurrence": {
			Description: "The properties of an hourly schedule. If the schedule occurs multiple times a day, specify the hourly recurrence.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"minute": {
						Description:  "Minutes of the hour the schedule will run.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 59),
					},
				},
			},
		},

		"time_zone_id": {
			Description:  "The time zone ID (e.g. Pacific Standard time).",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: computeValidate.VirtualMachineTimeZoneCaseInsensitive(),
		},

		"notification_settings": {
			Description: "The notification setting of a schedule.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Description: "Whether the notification is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"time_in_minutes": {
						Description:  "Time in minutes before event at which notification will be sent.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},
					"webhook_url": {
						Description: "The webhook URL to which the notification will be sent.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}

	if !features.SixPointOh() {
		s["status"] = &pluginsdk.Schema{
			Description:   "The status of this schedule.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ValidateFunc:  validation.StringInSlice(schedules.PossibleValuesForEnableStatus(), false),
			ConflictsWith: []string{"enabled"},
			Deprecated:    "`status` has been deprecated in favour of `enabled` and will be removed in v6.0 of the AzureRM Provider",
		}

		s["enabled"] = &pluginsdk.Schema{
			Description:   "Whether this schedule is enabled.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"status"},
		}

		notification := s["notification_settings"].Elem.(*pluginsdk.Resource).Schema

		notification["status"] = &pluginsdk.Schema{
			Description:   "The status of the notification.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ValidateFunc:  validation.StringInSlice(schedules.PossibleValuesForEnableStatus(), false),
			ConflictsWith: []string{"notification_settings.0.enabled"},
			Deprecated:    "`notification_settings.status` has been deprecated in favour of `notification_settings.enabled` and will be removed in v6.0 of the AzureRM Provider",
		}

		notification["enabled"] = &pluginsdk.Schema{
			Description:   "Whether the notification is enabled.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"notification_settings.0.status"},
		}
	}

	return s
}
