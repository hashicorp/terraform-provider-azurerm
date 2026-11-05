// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_origin_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceCdnFrontDoorOriginGroupFlatten(d *pluginsdk.ResourceData, id *afdorigingroups.OriginGroupId, model *afdorigingroups.AFDOriginGroup) error {
	d.Set("name", id.OriginGroupName)
	d.Set("cdn_frontdoor_profile_id", afdorigingroups.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID())

	if model != nil {
		if props := model.Properties; props != nil {
			if err := d.Set("health_probe", flattenCdnFrontDoorOriginGroupHealthProbeParameters(props.HealthProbeSettings)); err != nil {
				return fmt.Errorf("setting 'health_probe': %+v", err)
			}

			if err := d.Set("load_balancing", flattenCdnFrontDoorOriginGroupLoadBalancingSettingsParameters(props.LoadBalancingSettings)); err != nil {
				return fmt.Errorf("setting 'load_balancing': %+v", err)
			}

			d.Set("session_affinity_enabled", pointer.From(props.SessionAffinityState) == afdorigingroups.EnabledStateEnabled)
			d.Set("restore_traffic_time_to_healed_or_new_endpoint_in_minutes", props.TrafficRestorationTimeToHealedOrNewEndpointsInMinutes)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenCdnFrontDoorOriginGroupLoadBalancingSettingsParameters(input *afdorigingroups.LoadBalancingSettingsParameters) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"additional_latency_in_milliseconds": pointer.From(input.AdditionalLatencyInMilliseconds),
			"sample_size":                        pointer.From(input.SampleSize),
			"successful_samples_required":        pointer.From(input.SuccessfulSamplesRequired),
		},
	}
}

func flattenCdnFrontDoorOriginGroupHealthProbeParameters(input *afdorigingroups.HealthProbeParameters) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"interval_in_seconds": pointer.From(input.ProbeIntervalInSeconds),
			"path":                pointer.From(input.ProbePath),
			"protocol":            pointer.FromEnum(input.ProbeProtocol),
			"request_type":        pointer.FromEnum(input.ProbeRequestType),
		},
	}
}
