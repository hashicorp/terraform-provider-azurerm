// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_origin_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorOriginGroupCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profile, err := afdorigingroups.ParseProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	id := afdorigingroups.NewOriginGroupID(profile.SubscriptionId, profile.ResourceGroupName, profile.ProfileName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError(azureCdnFrontDoorOriginGroupResourceName, id.ID())
		}
	}

	sessionAffinity := afdorigingroups.EnabledStateDisabled
	if d.Get("session_affinity_enabled").(bool) {
		sessionAffinity = afdorigingroups.EnabledStateEnabled
	}

	props := afdorigingroups.AFDOriginGroup{
		Properties: &afdorigingroups.AFDOriginGroupProperties{
			HealthProbeSettings:   expandCdnFrontDoorOriginGroupHealthProbeParameters(d.Get("health_probe").([]any)),
			LoadBalancingSettings: expandCdnFrontDoorOriginGroupLoadBalancingSettingsParameters(d.Get("load_balancing").([]any)),
			SessionAffinityState:  pointer.To(sessionAffinity),
			TrafficRestorationTimeToHealedOrNewEndpointsInMinutes: pointer.To(int64(d.Get("restore_traffic_time_to_healed_or_new_endpoint_in_minutes").(int))),
		},
	}

	if err := client.CreateCallbackThenPoll(ctx, id, props, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceCdnFrontDoorOriginGroupRead(d, meta)
}

func expandCdnFrontDoorOriginGroupHealthProbeParameters(input []any) *afdorigingroups.HealthProbeParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)

	return &afdorigingroups.HealthProbeParameters{
		ProbeIntervalInSeconds: pointer.To(int64(v["interval_in_seconds"].(int))),
		ProbePath:              pointer.To(v["path"].(string)),
		ProbeProtocol:          pointer.ToEnum[afdorigingroups.ProbeProtocol](v["protocol"].(string)),
		ProbeRequestType:       pointer.ToEnum[afdorigingroups.HealthProbeRequestType](v["request_type"].(string)),
	}
}

func expandCdnFrontDoorOriginGroupLoadBalancingSettingsParameters(input []any) *afdorigingroups.LoadBalancingSettingsParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)

	return &afdorigingroups.LoadBalancingSettingsParameters{
		AdditionalLatencyInMilliseconds: pointer.To(int64(v["additional_latency_in_milliseconds"].(int))),
		SampleSize:                      pointer.To(int64(v["sample_size"].(int))),
		SuccessfulSamplesRequired:       pointer.To(int64(v["successful_samples_required"].(int))),
	}
}
