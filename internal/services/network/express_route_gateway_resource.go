// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/expressrouteconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/expressroutegateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExpressRouteGatewayCreateUpdate,
		Read:   resourceExpressRouteGatewayRead,
		Update: resourceExpressRouteGatewayCreateUpdate,
		Delete: resourceExpressRouteGatewayDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := expressroutegateways.ParseExpressRouteGatewayID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"virtual_hub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: virtualwans.ValidateVirtualHubID,
			},

			"scale_units": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},

			"allow_non_virtual_wan_traffic": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceExpressRouteGatewayCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteGateways
	connectionsClient := meta.(*clients.Client).Network.ExpressRouteConnections
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Println("[INFO] preparing arguments for ExpressRoute Gateway creation.")

	id := expressroutegateways.NewExpressRouteGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for present of existing %s: %+v", id, err)
			}
		}
		if resp.Model != nil && resp.Model.Id != nil && *resp.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_express_route_gateway", id.ID())
		}
	}

	gatewayId := expressrouteconnections.NewExpressRouteGatewayID(id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteGatewayName)

	respConnections, err := connectionsClient.List(ctx, gatewayId)
	if err != nil {
		// service will return 404 error if Gateway not exist
		if v, ok := err.(autorest.DetailedError); ok && v.StatusCode == http.StatusNotFound {
			log.Printf("[Debug]: Gateway connection not found. HTTP Code 404.")
		} else {
			return fmt.Errorf("retrieving %s: %+v", gatewayId, err)
		}
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
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteGatewayRead(d, meta)
}

func resourceExpressRouteGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteGateways
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressroutegateways.ParseExpressRouteGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ExpressRouteGatewayName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			d.Set("virtual_hub_id", pointer.From(props.VirtualHub.Id))
			d.Set("allow_non_virtual_wan_traffic", props.AllowNonVirtualWanTraffic)

			scaleUnits := 0
			if props.AutoScaleConfiguration != nil && props.AutoScaleConfiguration.Bounds != nil && props.AutoScaleConfiguration.Bounds.Min != nil {
				scaleUnits = int(*props.AutoScaleConfiguration.Bounds.Min)
			}
			d.Set("scale_units", scaleUnits)
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceExpressRouteGatewayDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteGateways
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressroutegateways.ParseExpressRouteGatewayID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func convertConnectionsToGatewayConnections(input *[]expressrouteconnections.ExpressRouteConnection) *[]expressroutegateways.ExpressRouteConnection {
	output := make([]expressroutegateways.ExpressRouteConnection, 0)

	if input == nil || len(*input) == 0 {
		return &output
	}
	for _, i := range *input {
		o := expressroutegateways.ExpressRouteConnection{
			Id:   i.Id,
			Name: i.Name,
		}

		if props := i.Properties; props != nil {
			o.Properties = &expressroutegateways.ExpressRouteConnectionProperties{
				AuthorizationKey:          props.AuthorizationKey,
				EnableInternetSecurity:    props.EnableInternetSecurity,
				EnablePrivateLinkFastPath: props.EnablePrivateLinkFastPath,
				ExpressRouteCircuitPeering: expressroutegateways.ExpressRouteCircuitPeeringId{
					Id: props.ExpressRouteCircuitPeering.Id,
				},
				ExpressRouteGatewayBypass: props.ExpressRouteGatewayBypass,
				ProvisioningState:         (*expressroutegateways.ProvisioningState)(props.ProvisioningState),
				RoutingConfiguration: &expressroutegateways.RoutingConfiguration{
					AssociatedRouteTable: &expressroutegateways.SubResource{
						Id: props.RoutingConfiguration.AssociatedRouteTable.Id,
					},
					InboundRouteMap: &expressroutegateways.SubResource{
						Id: props.RoutingConfiguration.InboundRouteMap.Id,
					},
					OutboundRouteMap: &expressroutegateways.SubResource{
						Id: props.RoutingConfiguration.OutboundRouteMap.Id,
					},
					PropagatedRouteTables: &expressroutegateways.PropagatedRouteTable{
						Ids:    convertConnectionsSubresourceToGatewaySubResource(props.RoutingConfiguration.PropagatedRouteTables.Ids),
						Labels: props.RoutingConfiguration.PropagatedRouteTables.Labels,
					},
					VnetRoutes: &expressroutegateways.VnetRoute{
						BgpConnections: convertConnectionsSubresourceToGatewaySubResource(props.RoutingConfiguration.VnetRoutes.BgpConnections),
						StaticRoutes:   convertConnectionsStaticRouteToGatewayStaticRoute(props.RoutingConfiguration.VnetRoutes.StaticRoutes),
						StaticRoutesConfig: &expressroutegateways.StaticRoutesConfig{
							PropagateStaticRoutes:          i.Properties.RoutingConfiguration.VnetRoutes.StaticRoutesConfig.PropagateStaticRoutes,
							VnetLocalRouteOverrideCriteria: (*expressroutegateways.VnetLocalRouteOverrideCriteria)(i.Properties.RoutingConfiguration.VnetRoutes.StaticRoutesConfig.VnetLocalRouteOverrideCriteria),
						},
					},
				},
				RoutingWeight: props.RoutingWeight,
			}
		}
		output = append(output, o)
	}

	return &output
}

func convertConnectionsStaticRouteToGatewayStaticRoute(input *[]expressrouteconnections.StaticRoute) *[]expressroutegateways.StaticRoute {
	output := make([]expressroutegateways.StaticRoute, 0)

	if input == nil || len(*input) == 0 {
		return &output
	}

	for _, i := range *input {
		output = append(output, expressroutegateways.StaticRoute{
			AddressPrefixes:  i.AddressPrefixes,
			Name:             i.Name,
			NextHopIPAddress: i.NextHopIPAddress,
		})
	}

	return &output
}

func convertConnectionsSubresourceToGatewaySubResource(input *[]expressrouteconnections.SubResource) *[]expressroutegateways.SubResource {
	output := make([]expressroutegateways.SubResource, 0)

	if input == nil || len(*input) == 0 {
		return &output
	}

	for _, i := range *input {
		output = append(output, expressroutegateways.SubResource{
			Id: i.Id,
		})
	}

	return &output
}
