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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualMachineID,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"daily_recurrence_time": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]|[0-9])[0-5][0-9]$"),
				"Time of day must match the format HHmm where HH is 00-23 and mm is 00-59",
			),
		},

		"timezone": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.VirtualMachineTimeZoneCaseInsensitive(),
		},

		"notification_settings": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"time_in_minutes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      30,
						ValidateFunc: validation.IntBetween(15, 120),
					},
					"webhook_url": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"email": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
