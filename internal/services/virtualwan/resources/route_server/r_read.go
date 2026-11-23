// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package route_server

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRouteServerRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVirtualHubID(d.Id())
	if err != nil {
		return err
	}

	routeServer, err := client.VirtualHubsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(routeServer.HttpResponse) {
			log.Printf("[INFO] %s does not exists - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.VirtualHubName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := routeServer.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			d.Set("sku", props.Sku)
			var virtualRouterIps *[]string
			if props.VirtualRouterIPs != nil {
				virtualRouterIps = props.VirtualRouterIPs
			}
			d.Set("virtual_router_ips", virtualRouterIps)
			if props.AllowBranchToBranchTraffic != nil {
				d.Set("branch_to_branch_traffic_enabled", props.AllowBranchToBranchTraffic)
			}
			d.Set("hub_routing_preference", pointer.From(props.HubRoutingPreference))
			if props.VirtualRouterAsn != nil {
				d.Set("virtual_router_asn", props.VirtualRouterAsn)
			}
			d.Set("routing_state", pointer.FromEnum(props.RoutingState))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("flattening `tags`: %+v", err)
		}
	}

	ipConfig, err := client.VirtualHubIPConfigurationList(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving IP Config for %s: %+v", id, err)
	}

	if model := ipConfig.Model; model != nil {
		for _, config := range *model {
			if props := config.Properties; props != nil {
				if props.PublicIPAddress != nil {
					d.Set("public_ip_address_id", pointer.From(props.PublicIPAddress.Id))
				}
				if props.Subnet != nil {
					d.Set("subnet_id", pointer.From(props.Subnet.Id))
				}
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func deleteRouteServerIpConfiguration(ctx context.Context, client *virtualwans.VirtualWANsClient, id commonids.VirtualHubIPConfigurationId) error {
	if err := client.VirtualHubIPConfigurationDeleteThenPoll(ctx, id); err != nil {
		return fmt.Errorf("deleting IP Config %s: %+v", id, err)
	}

	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: ipConfigStateRefreshFunc(ctx, client, id),
		Timeout: time.Until(timeout),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s: %+v", id, err)
	}
	return nil
}

func routeServerCreateRefreshFunc(ctx context.Context, client *virtualwans.VirtualWANsClient, id virtualwans.VirtualHubId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		resp, err := client.VirtualHubsGet(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return nil, "", fmt.Errorf("%s does not exists", id)
			}
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Model != nil && resp.Model.Properties != nil {
			return resp, pointer.FromEnum(resp.Model.Properties.ProvisioningState), nil
		}
		return nil, "", fmt.Errorf("unable to read the provisioning state of this %s", id)
	}
}

func ipConfigStateRefreshFunc(ctx context.Context, client *virtualwans.VirtualWANsClient, id commonids.VirtualHubIPConfigurationId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		resp, err := client.VirtualHubIPConfigurationGet(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, strconv.Itoa(resp.HttpResponse.StatusCode), nil
			}
			return nil, "", fmt.Errorf("polling for the status of %s: %+v", id, err)
		}
		return resp, strconv.Itoa(resp.HttpResponse.StatusCode), nil
	}
}
