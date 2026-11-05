// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_origin_group

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cdnFrontdoorOriginGroupSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorOriginGroupName,
		},

		"cdn_frontdoor_profile_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: afdorigingroups.ValidateProfileID,
		},

		"load_balancing": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"additional_latency_in_milliseconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      50,
						ValidateFunc: validation.IntBetween(0, 1000),
					},

					"sample_size": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      4,
						ValidateFunc: validation.IntBetween(0, 255),
					},

					"successful_samples_required": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      3,
						ValidateFunc: validation.IntBetween(0, 255),
					},
				},
			},
		},

		// Optional
		"health_probe": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"protocol": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(afdorigingroups.ProbeProtocolHTTP),
							string(afdorigingroups.ProbeProtocolHTTPS),
						}, false),
					},

					"request_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(afdorigingroups.HealthProbeRequestTypeHEAD),
						ValidateFunc: validation.StringInSlice([]string{
							string(afdorigingroups.HealthProbeRequestTypeGET),
							string(afdorigingroups.HealthProbeRequestTypeHEAD),
						}, false),
					},

					"interval_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 255),
					},

					"path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "/",
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"session_affinity_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"restore_traffic_time_to_healed_or_new_endpoint_in_minutes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      10,
			ValidateFunc: validation.IntBetween(0, 50),
		},
	}
}
