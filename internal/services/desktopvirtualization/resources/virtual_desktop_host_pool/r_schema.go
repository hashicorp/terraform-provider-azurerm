// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_host_pool

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualDesktopHostPoolSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Virtual Desktop Host Pool.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"type": {
			Description: "The type of the Virtual Desktop Host Pool.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(hostpool.HostPoolTypePersonal),
				string(hostpool.HostPoolTypePooled),
			}, false),
		},

		"load_balancer_type": {
			Description:  "The type of the load balancer. Refer to the [Azure documentation](https://learn.microsoft.com/azure/virtual-desktop/host-pool-load-balancing) for more information.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(hostpool.PossibleValuesForLoadBalancerType(), false),
		},

		"friendly_name": {
			Description:  "A friendly name for the Virtual Desktop Host Pool.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
		},

		"description": {
			Description:  "A description for the Virtual Desktop Host Pool.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 512),
		},

		"validate_environment": {
			Description: "Allows you to test service changes before they are deployed to production.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"custom_rdp_properties": {
			Description: "A valid custom RDP properties string for the Virtual Desktop Host Pool, available properties can be [found in this article](https://learn.microsoft.com/windows-server/remote/remote-desktop-services/clients/rdp-files).",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"personal_desktop_assignment_type": {
			Description:  "`Automatic` assignment - The service will select an available host and assign it to a user. `Direct` assignment - Admin selects a specific host to assign to a user.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(hostpool.PossibleValuesForPersonalDesktopAssignmentType(), false),
		},

		"public_network_access": {
			Description:  "Whether public network access is allowed for the Virtual Desktop Host Pool.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(hostpool.PossibleValuesForHostpoolPublicNetworkAccess(), false),
			Default:      string(hostpool.HostpoolPublicNetworkAccessEnabled),
		},

		"maximum_sessions_allowed": {
			Description:  "A valid integer value from 0 to 999999 for the maximum number of users that have concurrent sessions on a session host. Should only be set if the `type` of your Virtual Desktop Host Pool is `Pooled`.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      999999,
			ValidateFunc: validation.IntBetween(0, 999999),
		},

		"start_vm_on_connect": {
			Description: "Enables or disables the Start VM on Connection Feature.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"preferred_app_group_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Description:  "Option to specify the preferred Application Group type for the Virtual Desktop Host Pool.",
			ValidateFunc: validation.StringInSlice(hostpool.PossibleValuesForPreferredAppGroupType(), false),
			Default:      string(hostpool.PreferredAppGroupTypeDesktop),
		},

		"scheduled_agent_updates": {
			Description: "This enables control of when Agent Updates will be applied to Session Hosts.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Description: "Enables or disables scheduled updates of the AVD agent components (RDAgent, Geneva Monitoring agent, and side-by-side stack) on session hosts. If this is enabled then up to two `schedule` blocks must be defined.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"timezone": {
						Description: "The time zone in which the agent update schedule will apply, [the possible values are defined here](https://jackstromberg.com/2017/01/list-of-time-zones-consumed-by-azure/). If `use_session_host_timezone` is enabled then it will override this setting.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     "UTC",
					},

					"use_session_host_timezone": {
						Description: "Whether scheduled agent updates should be applied based on the timezone of the affected session host. If configured then this setting overrides `timezone`.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"schedule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 2,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"day_of_week": {
									Description:  "The day of the week on which agent updates should be performed.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(hostpool.PossibleValuesForDayOfWeek(), false),
								},

								"hour_of_day": {
									Description:  "The hour of day the update window should start. The update is a 2 hour period following the hour provided. A leading zero should not be used.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(0, 23),
								},
							},
						},
					},
				},
			},
		},

		"vm_template": {
			Description:      "A VM template for session hosts configuration within hostpool. This is a JSON string.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"tags": commonschema.Tags(),
	}
}
