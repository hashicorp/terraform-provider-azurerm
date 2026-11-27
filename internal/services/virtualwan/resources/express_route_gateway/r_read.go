// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_gateway

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutegateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteGatewayRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.ExpressRouteGateways
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
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
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
