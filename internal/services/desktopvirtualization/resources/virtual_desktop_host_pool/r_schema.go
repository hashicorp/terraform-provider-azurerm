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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(hostpool.HostPoolTypePersonal),
				string(hostpool.HostPoolTypePooled),
			}, false),
		},

		"load_balancer_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(hostpool.PossibleValuesForLoadBalancerType(), false),
		},

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

		"validate_environment": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"custom_rdp_properties": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"personal_desktop_assignment_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(hostpool.PossibleValuesForPersonalDesktopAssignmentType(), false),
		},

		"public_network_access": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(hostpool.PossibleValuesForHostpoolPublicNetworkAccess(), false),
			Default:      string(hostpool.HostpoolPublicNetworkAccessEnabled),
		},

		"maximum_sessions_allowed": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      999999,
			ValidateFunc: validation.IntBetween(0, 999999),
		},

		"start_vm_on_connect": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"preferred_app_group_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Description:  "Preferred App Group type to display",
			ValidateFunc: validation.StringInSlice(hostpool.PossibleValuesForPreferredAppGroupType(), false),
			Default:      string(hostpool.PreferredAppGroupTypeDesktop),
		},

		"scheduled_agent_updates": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"timezone": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "UTC",
					},

					"use_session_host_timezone": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"schedule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 2,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"day_of_week": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(hostpool.PossibleValuesForDayOfWeek(), false),
								},

								"hour_of_day": {
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
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"tags": commonschema.Tags(),
	}
}
