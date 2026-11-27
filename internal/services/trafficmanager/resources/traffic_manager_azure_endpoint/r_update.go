// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package traffic_manager_azure_endpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/trafficmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAzureEndpointUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := profiles.ParseTrafficManagerProfileID(d.Get("profile_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `profile_id`: %+v", err)
	}

	id := trafficmanagers.NewEndpointTypeID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.TrafficManagerProfileName, trafficmanagers.EndpointTypeAzureEndpoints, d.Get("name").(string))

	existing, err := client.EndpointsGet(ctx, id)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %v", id, err)
	}

	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("model/properties was nil for %s", id)
	}

	params := *existing.Model

	if d.HasChange("enabled") {
		status := trafficmanagers.EndpointStatusEnabled
		if !d.Get("enabled").(bool) {
			status = trafficmanagers.EndpointStatusDisabled
		}
		params.Properties.EndpointStatus = pointer.To(status)
	}

	if d.HasChange("always_serve_enabled") {
		alwaysServe := trafficmanagers.AlwaysServeDisabled
		if d.Get("always_serve_enabled").(bool) {
			alwaysServe = trafficmanagers.AlwaysServeEnabled
		}
		params.Properties.AlwaysServe = pointer.To(alwaysServe)
	}

	if d.HasChange("custom_header") {
		params.Properties.CustomHeaders = helpers.ExpandEndpointCustomHeaderConfig(d.Get("custom_header").([]any))
	}

	if d.HasChange("target_resource_id") {
		params.Properties.TargetResourceId = pointer.To(d.Get("target_resource_id").(string))
	}

	if d.HasChange("subnet") {
		params.Properties.Subnets = helpers.ExpandEndpointSubnetConfig(d.Get("subnet").([]any))
	}

	if d.HasChange("priority") {
		if priority := d.Get("priority").(int); priority != 0 {
			params.Properties.Priority = pointer.To(int64(priority))
		}
	}

	if d.HasChange("weight") {
		if weight := d.Get("weight").(int); weight != 0 {
			params.Properties.Weight = pointer.To(int64(weight))
		}
	}

	if d.HasChange("geo_mappings") {
		inputMappings := d.Get("geo_mappings").([]any)
		geoMappings := make([]string, 0)
		for _, v := range inputMappings {
			geoMappings = append(geoMappings, v.(string))
		}
		if len(geoMappings) > 0 {
			params.Properties.GeoMapping = &geoMappings
		} else {
			params.Properties.GeoMapping = nil
		}
	}

	if _, err := client.EndpointsCreateOrUpdate(ctx, id, params); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceAzureEndpointRead(d, meta)
}
