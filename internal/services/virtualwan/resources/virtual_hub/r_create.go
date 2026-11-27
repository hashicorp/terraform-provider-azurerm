// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if _, ok := ctx.Deadline(); !ok {
		return errors.New("deadline is not properly set for Virtual Hub")
	}

	id := virtualwans.NewVirtualHubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.VirtualHubsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for present of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_hub", id.ID())
		}
	}

	parameters := virtualwans.VirtualHub{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &virtualwans.VirtualHubProperties{
			AllowBranchToBranchTraffic: pointer.To(d.Get("branch_to_branch_traffic_enabled").(bool)),
			RouteTable:                 expandVirtualHubRoute(d.Get("route").(*pluginsdk.Set).List()),
			HubRoutingPreference:       pointer.ToEnum[virtualwans.HubRoutingPreference](d.Get("hub_routing_preference").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("address_prefix"); ok {
		parameters.Properties.AddressPrefix = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("sku"); ok {
		parameters.Properties.Sku = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("virtual_wan_id"); ok {
		parameters.Properties.VirtualWAN = &virtualwans.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	if v, ok := d.GetOk("virtual_router_auto_scale_min_capacity"); ok {
		parameters.Properties.VirtualRouterAutoScaleConfiguration = &virtualwans.VirtualRouterAutoScaleConfiguration{
			MinCapacity: pointer.To(int64(v.(int))),
		}
	}

	if err := client.VirtualHubsCreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	// Hub returns provisioned while the routing state is still "provisioning". This might cause issues with following hubvnet connection operations.
	// https://github.com/Azure/azure-rest-api-specs/issues/10391
	// As a workaround, we will poll the routing state and ensure it is "Provisioned".

	// deadline is checked at the entry point of this function
	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Provisioning"},
		Target:                    []string{"Provisioned", "Failed", "None"},
		Refresh:                   virtualHubCreateRefreshFunc(ctx, client, id),
		PollInterval:              15 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   time.Until(timeout),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s provisioning route: %+v", id, err)
	}

	return resourceVirtualHubRead(d, meta)
}

func expandVirtualHubRoute(input []any) *virtualwans.VirtualHubRouteTable {
	if len(input) == 0 {
		return nil
	}

	results := make([]virtualwans.VirtualHubRoute, 0)
	for _, item := range input {
		if item == nil {
			continue
		}

		v := item.(map[string]any)
		addressPrefixes := v["address_prefixes"].([]any)
		nextHopIpAddress := v["next_hop_ip_address"].(string)

		results = append(results, virtualwans.VirtualHubRoute{
			AddressPrefixes:  pluginsdk.ExpandStringSlice(addressPrefixes),
			NextHopIPAddress: pointer.To(nextHopIpAddress),
		})
	}

	result := virtualwans.VirtualHubRouteTable{
		Routes: &results,
	}

	return &result
}
