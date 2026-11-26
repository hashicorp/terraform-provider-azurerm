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
			Description:  "The name which should be used for this Virtual Desktop Scaling Plan.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"friendly_name": {
			Description:  "Friendly name of the Scaling Plan.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
		},

		"description": {
			Description:  "A description of the Scaling Plan.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 512),
		},

		"time_zone": {
			Description: "The Time Zone which should be used by the Scaling Plan for time based events, [the possible values are defined here](https://jackstromberg.com/2017/01/list-of-time-zones-consumed-by-azure/).",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},

		"exclusion_tag": {
			Description: "The name of the tag associated with the VMs you want to exclude from autoscaling.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"schedule": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the schedule.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"days_of_week": {
						Description: "A list of days of the week on which this schedule will be used.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForDayOfWeek(), false),
						},
					},

					"ramp_up_start_time": {
						Description:  "The time at which Ramp-Up scaling will begin. This is also the end-time for the Ramp-Up period. The time must be specified in `HH:MM` format.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validateTime(),
					},

					"ramp_up_load_balancing_algorithm": {
						Description:  "The load balancing algorithm to use during the Ramp-Up period.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForSessionHostLoadBalancingAlgorithm(), false),
					},

					"ramp_up_minimum_hosts_percent": {
						Description:  "The minimum percentage of session host Virtual Machines to start during ramp-up for peak hours. For example, if Minimum percentage of hosts is specified as `10%` and total number of session hosts in your host pool is `10`, autoscale will ensure a minimum of `1` session host is available to take user connections.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 100),
					},

					"ramp_up_capacity_threshold_percent": {
						Description:  "This is the value of percentage of used host pool capacity that will be considered to evaluate whether to turn on/off Virtual Machines during the ramp-up and peak hours. For example, if capacity threshold is specified as `60%` and your total host pool capacity is `100` sessions, autoscale will turn on additional session hosts once the host pool exceeds a load of `60` sessions.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 100),
					},

					"peak_start_time": {
						Description:  "The time at which Peak scaling will begin. This is also the end-time for the Ramp-Up period. The time must be specified in `HH:MM` format.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validateTime(),
					},

					"peak_load_balancing_algorithm": {
						Description:  "The load balancing algorithm to use during Peak Hours.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForSessionHostLoadBalancingAlgorithm(), false),
					},

					"ramp_down_start_time": {
						Description:  "The time at which Ramp-Down scaling will begin. This is also the end-time for the Ramp-Up period. The time must be specified in `HH:MM` format.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validateTime(),
					},

					"ramp_down_load_balancing_algorithm": {
						Description:  "The load balancing algorithm to use during the Ramp-Down period.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForSessionHostLoadBalancingAlgorithm(), false),
					},

					"ramp_down_minimum_hosts_percent": {
						Description:  "The minimum percentage of session host Virtual Machines that you would like to get to for ramp-down and off-peak hours. For example, if Minimum percentage of hosts is specified as 10% and total number of session hosts in your host pool is 10, autoscale will ensure a minimum of 1 session host is available to take user connections.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 100),
					},

					"ramp_down_capacity_threshold_percent": {
						Description:  "This is the value in percentage of used host pool capacity that will be considered to evaluate whether to turn on/off Virtual Machines during the ramp-down and off-peak hours. For example, if capacity threshold is specified as 60% and your total host pool capacity is 100 sessions, autoscale will turn on additional session hosts once the host pool exceeds a load of 60 sessions.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 100),
					},

					"ramp_down_force_logoff_users": {
						Description: "Whether users will be forced to log-off session hosts once the `ramp_down_wait_time_minutes` value has been exceeded during the Ramp-Down period.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},

					"ramp_down_stop_hosts_when": {
						Description:  "Controls Session Host shutdown behaviour during Ramp-Down period. Session Hosts can either be shutdown when all sessions on the Session Host have ended, or when there are no Active sessions left on the Session Host.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForStopHostsWhen(), false),
					},

					"ramp_down_wait_time_minutes": {
						Description: "The number of minutes during Ramp-Down period that autoscale will wait after setting the session host VMs to drain mode, notifying any currently signed in users to save their work before forcing the users to logoff. Once all user sessions on the session host VM have been logged off, Autoscale will shut down the VM.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},

					"ramp_down_notification_message": {
						Description: "The notification message to send to users during Ramp-Down period when they are required to log-off.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},

					"off_peak_start_time": {
						Description:  "The time at which Off-Peak scaling will begin. This is also the end-time for the Ramp-Down period. The time must be specified in `HH:MM` format.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validateTime(),
					},

					"off_peak_load_balancing_algorithm": {
						Description:  "The load balancing algorithm to use during Off-Peak Hours.",
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
						Description:  "The ID of the HostPool to assign the Scaling Plan to.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: hostpool.ValidateHostPoolID,
					},
					"scaling_plan_enabled": {
						Description: "Whether the scaling plan is enabled or disabled for the HostPool.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
