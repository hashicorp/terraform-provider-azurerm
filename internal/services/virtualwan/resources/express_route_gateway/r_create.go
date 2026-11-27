// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutegateways"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteGatewayCreate(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).ExpressRouteGateways
	connectionsClient := virtualwanClient.From(meta).ExpressRouteConnections
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id := expressroutegateways.NewExpressRouteGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if resp.Model != nil && resp.Model.Id != nil && *resp.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_express_route_gateway", id.ID())
		}
	}

	gatewayId := expressrouteconnections.NewExpressRouteGatewayID(id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteGatewayName)

	respConnections, err := connectionsClient.List(ctx, gatewayId)
	if err != nil && !response.WasNotFound(respConnections.HttpResponse) {
		return fmt.Errorf("retrieving %s: %+v", gatewayId, err)
	}

	var connections *[]expressroutegateways.ExpressRouteConnection
	if model := respConnections.Model; model != nil {
		connections = convertConnectionsToGatewayConnections(model.Value)
	}

	parameters := expressroutegateways.ExpressRouteGateway{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &expressroutegateways.ExpressRouteGatewayProperties{
			AllowNonVirtualWanTraffic: pointer.To(d.Get("allow_non_virtual_wan_traffic").(bool)),
			AutoScaleConfiguration: &expressroutegateways.ExpressRouteGatewayPropertiesAutoScaleConfiguration{
				Bounds: &expressroutegateways.ExpressRouteGatewayPropertiesAutoScaleConfigurationBounds{
					Min: pointer.To(int64(d.Get("scale_units").(int))),
				},
			},
			VirtualHub: expressroutegateways.VirtualHubId{
				Id: pointer.To(d.Get("virtual_hub_id").(string)),
			},
			ExpressRouteConnections: connections,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteGatewayRead(d, meta)
}
