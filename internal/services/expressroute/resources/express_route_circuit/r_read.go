// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuits"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteports"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteCircuitRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ExpressRoute.ExpressRouteCircuits
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressroutecircuits.ParseExpressRouteCircuitID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ExpressRouteCircuitName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if err := d.Set("sku", flattenExpressRouteCircuitSku(model.Sku)); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}
		if props := model.Properties; props != nil {
			d.Set("bandwidth_in_gbps", props.BandwidthInGbps)

			if props.ExpressRoutePort != nil && props.ExpressRoutePort.Id != nil {
				portID, err := expressrouteports.ParseExpressRoutePortIDInsensitively(*props.ExpressRoutePort.Id)
				if err != nil {
					return err
				}
				d.Set("express_route_port_id", portID.ID())
			}

			d.Set("service_provider_provisioning_state", pointer.FromEnum(props.ServiceProviderProvisioningState))
			d.Set("service_key", props.ServiceKey)
			d.Set("allow_classic_operations", props.AllowClassicOperations)
			d.Set("rate_limiting_enabled", props.EnableDirectPortRateLimit)

			if serviceProviderProps := props.ServiceProviderProperties; serviceProviderProps != nil {
				d.Set("service_provider_name", serviceProviderProps.ServiceProviderName)
				d.Set("peering_location", serviceProviderProps.PeeringLocation)
				d.Set("bandwidth_in_mbps", serviceProviderProps.BandwidthInMbps)
			}
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}

func flattenExpressRouteCircuitSku(sku *expressroutecircuits.ExpressRouteCircuitSku) []any {
	if sku == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"tier":   pointer.FromEnum(sku.Tier),
			"family": pointer.FromEnum(sku.Family),
		},
	}
}

func expressRouteCircuitCreationRefreshFunc(ctx context.Context, client *expressroutecircuits.ExpressRouteCircuitsClient, id expressroutecircuits.ExpressRouteCircuitId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return nil, "NotFound", nil
			}

			return nil, "", fmt.Errorf("polling to check if the Express Route Circuit has been created: %+v", err)
		}

		if model := res.Model; model != nil && model.Properties != nil && model.Properties.ProvisioningState != nil {
			switch *model.Properties.ProvisioningState {
			case expressroutecircuits.ProvisioningStateFailed:
				return nil, "", fmt.Errorf("%s failed to provision: %+v", id.ID(), res)
			case expressroutecircuits.ProvisioningStateDeleting:
				return nil, "", fmt.Errorf("%s is in a deleting state: %+v", id.ID(), res)
			case expressroutecircuits.ProvisioningStateUpdating:
				return nil, "Waiting", nil
			}
		}

		return res, "Exists", nil
	}
}
