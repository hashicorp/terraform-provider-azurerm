// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_scaling_plan

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/scalingplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualDesktopScalingPlanSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"friendly_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 512),
		},

		"time_zone": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"exclusion_tag": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"schedule": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"days_of_week": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForDayOfWeek(), false),
						},
					},

					"ramp_up_start_time": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validateTime(),
					},

					"ramp_up_load_balancing_algorithm": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForSessionHostLoadBalancingAlgorithm(), false),
					},

					"ramp_up_minimum_hosts_percent": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 100),
					},

					"ramp_up_capacity_threshold_percent": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 100),
					},

					"peak_start_time": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validateTime(),
					},

					"peak_load_balancing_algorithm": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForSessionHostLoadBalancingAlgorithm(), false),
					},

					"ramp_down_start_time": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validateTime(),
					},

					"ramp_down_load_balancing_algorithm": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForSessionHostLoadBalancingAlgorithm(), false),
					},

					"ramp_down_minimum_hosts_percent": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 100),
					},

					"ramp_down_capacity_threshold_percent": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 100),
					},

					"ramp_down_force_logoff_users": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"ramp_down_stop_hosts_when": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForStopHostsWhen(), false),
					},

					"ramp_down_wait_time_minutes": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"ramp_down_notification_message": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"off_peak_start_time": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validateTime(),
					},

					"off_peak_load_balancing_algorithm": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForSessionHostLoadBalancingAlgorithm(), false),
					},
				},
			},
		},

		"host_pool": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"hostpool_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: hostpool.ValidateHostPoolID,
					},
					"scaling_plan_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
