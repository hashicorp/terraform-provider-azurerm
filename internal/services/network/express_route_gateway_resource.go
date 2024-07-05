// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/virtualwans"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressrouteconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutegateways"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExpressRouteGatewayCreate,
		Read:   resourceExpressRouteGatewayRead,
		Update: resourceExpressRouteGatewayUpdate,
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

func resourceExpressRouteGatewayCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteGateways
	connectionsClient := meta.(*clients.Client).Network.ExpressRouteConnections
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Println("[INFO] preparing arguments for ExpressRoute Gateway creation.")

	id := expressroutegateways.NewExpressRouteGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if resp.Model != nil && resp.Model.Id != nil && *resp.Model.Id != "" {
		return tf.ImportAsExistsError("azurerm_express_route_gateway", id.ID())
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
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteGatewayRead(d, meta)
}

func resourceExpressRouteGatewayUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteGateways
	connectionsClient := meta.(*clients.Client).Network.ExpressRouteConnections
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Println("[INFO] preparing arguments for ExpressRoute Gateway update.")

	id, err := expressroutegateways.ParseExpressRouteGatewayID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	gatewayId, err := expressrouteconnections.ParseExpressRouteGatewayID(d.Id())
	if err != nil {
		return err
	}

	respConnections, err := connectionsClient.List(ctx, *gatewayId)
	if err != nil && !response.WasNotFound(respConnections.HttpResponse) {
		return fmt.Errorf("retrieving %s: %+v", gatewayId, err)
	}

	payload := existing.Model

	var connections *[]expressroutegateways.ExpressRouteConnection
	if model := respConnections.Model; model != nil {
		connections = convertConnectionsToGatewayConnections(model.Value)
	}

	payload.Properties.ExpressRouteConnections = connections

	if d.HasChange("scale_units") {
		payload.Properties.AutoScaleConfiguration = &expressroutegateways.ExpressRouteGatewayPropertiesAutoScaleConfiguration{
			Bounds: &expressroutegateways.ExpressRouteGatewayPropertiesAutoScaleConfigurationBounds{
				Min: pointer.To(int64(d.Get("scale_units").(int))),
			},
		}
	}

	if d.HasChange("allow_non_virtual_wan_traffic") {
		payload.Properties.AllowNonVirtualWanTraffic = pointer.To(d.Get("allow_non_virtual_wan_traffic").(bool))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
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
				RoutingWeight:             props.RoutingWeight,
			}

			if routingConfiguration := props.RoutingConfiguration; routingConfiguration != nil {
				rc := &expressroutegateways.RoutingConfiguration{}

				if routingConfiguration.AssociatedRouteTable != nil {
					rc.AssociatedRouteTable = &expressroutegateways.SubResource{
						Id: routingConfiguration.AssociatedRouteTable.Id,
					}
				}

				if routingConfiguration.InboundRouteMap != nil {
					rc.InboundRouteMap = &expressroutegateways.SubResource{
						Id: routingConfiguration.InboundRouteMap.Id,
					}
				}

				if routingConfiguration.OutboundRouteMap != nil {
					rc.OutboundRouteMap = &expressroutegateways.SubResource{
						Id: routingConfiguration.OutboundRouteMap.Id,
					}
				}

				if routingConfiguration.PropagatedRouteTables != nil {
					rc.PropagatedRouteTables = &expressroutegateways.PropagatedRouteTable{
						Ids:    convertConnectionsSubresourceToGatewaySubResource(props.RoutingConfiguration.PropagatedRouteTables.Ids),
						Labels: routingConfiguration.PropagatedRouteTables.Labels,
					}
				}

				if vnet := routingConfiguration.VnetRoutes; vnet != nil {
					rc.VnetRoutes = &expressroutegateways.VnetRoute{
						BgpConnections: convertConnectionsSubresourceToGatewaySubResource(vnet.BgpConnections),
						StaticRoutes:   convertConnectionsStaticRouteToGatewayStaticRoute(vnet.StaticRoutes),
					}

					if src := vnet.StaticRoutesConfig; src != nil {
						rc.VnetRoutes.StaticRoutesConfig = &expressroutegateways.StaticRoutesConfig{
							PropagateStaticRoutes:          src.PropagateStaticRoutes,
							VnetLocalRouteOverrideCriteria: (*expressroutegateways.VnetLocalRouteOverrideCriteria)(src.VnetLocalRouteOverrideCriteria),
						}
					}
				}

				o.Properties.RoutingConfiguration = rc
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
