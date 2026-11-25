// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_global_vm_shutdown_schedule

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func devTestGlobalVmShutdownScheduleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		"virtual_machine_id": {
			Description:  "The resource ID of the target ARM-based Virtual Machine.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualMachineID,
		},

		"enabled": {
			Description: "Whether to enable the schedule.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"daily_recurrence_time": {
			Description: "The time each day when the schedule takes effect. Must match the format HHmm where HH is 00-23 and mm is 00-59 (e.g. 0930, 2300, etc.)",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]|[0-9])[0-5][0-9]$"),
				"Time of day must match the format HHmm where HH is 00-23 and mm is 00-59",
			),
		},

		"timezone": {
			Description:  "The time zone ID (e.g. Pacific Standard time). Refer to this guide for a [full list of accepted time zone names](https://jackstromberg.com/2017/01/list-of-time-zones-consumed-by-azure/).",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.VirtualMachineTimeZoneCaseInsensitive(),
		},

		"notification_settings": {
			Description: "The notification setting of a schedule.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Description: "Whether to enable pre-shutdown notifications.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},
					"time_in_minutes": {
						Description:  "Time in minutes between 15 and 120 before a shutdown event at which a notification will be sent.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      30,
						ValidateFunc: validation.IntBetween(15, 120),
					},
					"webhook_url": {
						Description: "The webhook URL to which the notification will be sent.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"email": {
						Description: "E-mail address to which the notification will be sent.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
