// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_origin_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorOriginGroupUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginGroupsClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdorigingroups.ParseOriginGroupID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties was nil", id)
	}
	props := existing.Model.Properties

	// The API requires that an explicit null be passed as the 'health_probe' value to disable the health probe
	// e.g. {"properties":{"healthProbeSettings":null}}
	if d.HasChange("health_probe") {
		props.HealthProbeSettings = expandCdnFrontDoorOriginGroupHealthProbeParameters(d.Get("health_probe").([]any))
	}

	if d.HasChange("load_balancing") {
		props.LoadBalancingSettings = expandCdnFrontDoorOriginGroupLoadBalancingSettingsParameters(d.Get("load_balancing").([]any))
	}

	if d.HasChange("restore_traffic_time_to_healed_or_new_endpoint_in_minutes") {
		props.TrafficRestorationTimeToHealedOrNewEndpointsInMinutes = pointer.To(int64(d.Get("restore_traffic_time_to_healed_or_new_endpoint_in_minutes").(int)))
	}

	if d.HasChange("session_affinity_enabled") {
		props.SessionAffinityState = pointer.To(afdorigingroups.EnabledStateDisabled)
		if d.Get("session_affinity_enabled").(bool) {
			props.SessionAffinityState = pointer.To(afdorigingroups.EnabledStateEnabled)
		}
	}

	if err := client.CreateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorOriginGroupRead(d, meta)
}
