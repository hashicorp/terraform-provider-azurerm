// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package frontdoor

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFrontDoorUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := frontdoors.NewFrontDoorID(subscriptionId, resourceGroup, name)

	// remove in 3.0
	// due to a change in the RP, if a Frontdoor exists in a location other than 'Global' it may continue to
	// exist in that location, if this is a brand new Frontdoor it must be created in the 'Global' location
	var loc string

	exists, err := client.Get(ctx, id)
	if err != nil || exists.Model == nil {
		return fmt.Errorf("locating %s: %+v", id, err)
	} else {
		loc = location.Normalize(*exists.Model.Location)
	}

	cfgLocation, hasLocation := d.GetOk("location")
	if hasLocation {
		if loc != location.Normalize(cfgLocation.(string)) {
			return fmt.Errorf("the Front Door %q (Resource Group %q) already exists in %q and cannot be moved to the %q location", name, resourceGroup, loc, cfgLocation)
		}
	}

	existingModel := *exists.Model

	if d.HasChange("friendly_name") {
		existingModel.Properties.FriendlyName = pointer.To(d.Get("friendly_name").(string))
	}

	routingRules := d.Get("routing_rule").([]any)
	if d.HasChange("routing_rule") {
		rulesEngines := make(map[string]*frontdoors.SubResource)
		if existingModel.Properties != nil && existingModel.Properties.RoutingRules != nil {
			for _, rule := range *existingModel.Properties.RoutingRules {
				if rule.Properties != nil && rule.Properties.RulesEngine != nil {
					rulesEngines[*rule.Name] = rule.Properties.RulesEngine
				}
			}
		}
		existingModel.Properties.RoutingRules = expandFrontDoorRoutingRule(routingRules, id, &rulesEngines)
	}

	loadBalancingSettings := d.Get("backend_pool_load_balancing").([]any)
	if d.HasChange("backend_pool_load_balancing") {
		existingModel.Properties.LoadBalancingSettings = expandFrontDoorLoadBalancingSettingsModel(loadBalancingSettings, id)
	}

	healthProbeSettings := d.Get("backend_pool_health_probe").([]any)
	if d.HasChange("backend_pool_health_probe") {
		existingModel.Properties.HealthProbeSettings = expandFrontDoorHealthProbeSettingsModel(healthProbeSettings, id)
	}

	backendPools := d.Get("backend_pool").([]any)
	if d.HasChange("backend_pool") {
		existingModel.Properties.BackendPools = expandFrontDoorBackendPools(backendPools, id)
	}

	frontendEndpoints := d.Get("frontend_endpoint").([]any)
	if d.HasChange("frontend_endpoint") {
		existingModel.Properties.FrontendEndpoints = expandFrontDoorFrontendEndpoint(frontendEndpoints, id)
	}

	if d.HasChange("backend_pool_settings") {
		var backendCertNameCheck bool
		var backendPoolsSendReceiveTimeoutSeconds int64
		if bps, ok := d.Get("backend_pool_settings").([]any); ok && len(bps) > 0 {
			bpsMap := bps[0].(map[string]any)
			if v, ok := bpsMap["enforce_backend_pools_certificate_name_check"].(bool); ok {
				backendCertNameCheck = v
			}
			if v, ok := bpsMap["backend_pools_send_receive_timeout_seconds"].(int); ok {
				backendPoolsSendReceiveTimeoutSeconds = int64(v)
			}
			existingModel.Properties.BackendPoolsSettings = expandFrontDoorBackendPoolsSettings(backendCertNameCheck, backendPoolsSendReceiveTimeoutSeconds)
		}
	}

	if d.HasChange("load_balancer_enabled") {
		enabledState := expandFrontDoorEnabledState(d.Get("load_balancer_enabled").(bool))
		existingModel.Properties.EnabledState = &enabledState
	}

	if d.HasChanges("tags") {
		existingModel.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	// If the explicitResourceOrder is empty and it's not a new resource set the mapping table to the state file and return an error.
	// If the explicitResourceOrder is empty and it is a new resource it will run the CreateOrUpdate as expected
	// If the explicitResourceOrder is NOT empty and it is NOT a new resource it will run the CreateOrUpdate as expected
	explicitResourceOrder := d.Get("explicit_resource_order").([]any)
	if len(explicitResourceOrder) == 0 {
		d.Set("explicit_resource_order", flattenExplicitResourceOrder(backendPools, frontendEndpoints, routingRules, loadBalancingSettings, healthProbeSettings, id))
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, existingModel); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.Set("explicit_resource_order", flattenExplicitResourceOrder(backendPools, frontendEndpoints, routingRules, loadBalancingSettings, healthProbeSettings, id))
	}

	d.SetId(id.ID())
	return resourceFrontDoorRead(d, meta)
}
