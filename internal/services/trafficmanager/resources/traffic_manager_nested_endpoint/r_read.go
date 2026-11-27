// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package traffic_manager_nested_endpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/trafficmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	trafficmanagerClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNestedEndpointRead(d *pluginsdk.ResourceData, meta any) error {
	client := trafficmanagerClient.From(meta).EndpointsClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := trafficmanagers.ParseEndpointTypeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.EndpointsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.EndpointName)
	d.Set("profile_id", profiles.NewTrafficManagerProfileID(id.SubscriptionId, id.ResourceGroupName, id.TrafficManagerProfileName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			enabled := true
			if props.EndpointStatus != nil && *props.EndpointStatus == trafficmanagers.EndpointStatusDisabled {
				enabled = false
			}
			d.Set("enabled", enabled)
			d.Set("target_resource_id", props.TargetResourceId)
			d.Set("weight", props.Weight)
			d.Set("minimum_child_endpoints", props.MinChildEndpoints)
			d.Set("minimum_required_child_endpoints_ipv4", props.MinChildEndpointsIPv4)
			d.Set("minimum_required_child_endpoints_ipv6", props.MinChildEndpointsIPv6)
			d.Set("priority", props.Priority)
			d.Set("endpoint_location", props.EndpointLocation)
			d.Set("geo_mappings", props.GeoMapping)

			if err := d.Set("custom_header", helpers.FlattenEndpointCustomHeaderConfig(props.CustomHeaders)); err != nil {
				return fmt.Errorf("setting `custom_header`: %s", err)
			}
			if err := d.Set("subnet", helpers.FlattenEndpointSubnetConfig(props.Subnets)); err != nil {
				return fmt.Errorf("setting `subnet`: %s", err)
			}
		}
	}

	return nil
}
