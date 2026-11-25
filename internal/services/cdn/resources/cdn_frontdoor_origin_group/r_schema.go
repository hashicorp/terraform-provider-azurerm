// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_origin_group

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/afdorigingroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cdnFrontdoorOriginGroupSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Front Door Origin Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorOriginGroupName,
		},

		"cdn_frontdoor_profile_id": {
			Description:  "The ID of the Front Door Profile within which this Front Door Origin Group should exist.",
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
						Description:  "The additional latency in milliseconds for probes to fall into the lowest latency bucket.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      50,
						ValidateFunc: validation.IntBetween(0, 1000),
					},

					"sample_size": {
						Description:  "The number of samples to consider for load balancing decisions.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      4,
						ValidateFunc: validation.IntBetween(0, 255),
					},

					"successful_samples_required": {
						Description:  "The number of samples within the sample period that must succeed.",
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
						Description: "The protocol to use for health probe.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(afdorigingroups.ProbeProtocolHTTP),
							string(afdorigingroups.ProbeProtocolHTTPS),
						}, false),
					},

					"request_type": {
						Description: "The type of health probe request that is made.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     string(afdorigingroups.HealthProbeRequestTypeHEAD),
						ValidateFunc: validation.StringInSlice([]string{
							string(afdorigingroups.HealthProbeRequestTypeGET),
							string(afdorigingroups.HealthProbeRequestTypeHEAD),
						}, false),
					},

					"interval_in_seconds": {
						Description:  "The number of seconds between health probes.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 255),
					},

					"path": {
						Description:  "The path relative to the origin that is used to determine the health of the origin.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "/",
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"session_affinity_enabled": {
			Description: "Whether session affinity should be enabled on this host.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"restore_traffic_time_to_healed_or_new_endpoint_in_minutes": {
			Description:  "The amount of time which should elapse before shifting traffic to another endpoint when a healthy endpoint becomes unhealthy or a new endpoint is added.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      10,
			ValidateFunc: validation.IntBetween(0, 50),
		},
	}
}
