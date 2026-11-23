// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package traffic_manager_nested_endpoint

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func trafficManagerNestedEndpointSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"profile_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: profiles.ValidateTrafficManagerProfileID,
		},

		"target_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"weight": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 1000),
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"minimum_child_endpoints": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"minimum_required_child_endpoints_ipv4": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"minimum_required_child_endpoints_ipv6": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"custom_header": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.NoZeroValues,
					},
					"value": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.NoZeroValues,
					},
				},
			},
		},

		"priority": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C the API dynamically increments the default value for priority depending on the number of endpoints
			Computed:     true,
			ValidateFunc: validation.IntBetween(1, 1000),
		},

		"endpoint_location": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true, // azignore:AZS007 - pre-existing violation
			StateFunc:        location.StateFunc,
			DiffSuppressFunc: location.DiffSuppressFunc,
		},

		"geo_mappings": {
			Type:     pluginsdk.TypeList,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Optional: true,
		},

		"subnet": {
			Type:     pluginsdk.TypeList,
			ForceNew: true,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"first": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsIPv4Address,
					},
					"last": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsIPv4Address,
					},
					"scope": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 32),
					},
				},
			},
		},
	}
}
