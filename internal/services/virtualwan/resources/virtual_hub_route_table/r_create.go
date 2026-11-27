// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_route_table

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubRouteTableCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtHubId, err := virtualwans.ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(virtHubId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtHubId.VirtualHubName, virtualHubResourceName)

	id := virtualwans.NewHubRouteTableID(virtHubId.SubscriptionId, virtHubId.ResourceGroupName, virtHubId.VirtualHubName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.HubRouteTablesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_hub_route_table", id.ID())
		}
	}

	parameters := virtualwans.HubRouteTable{
		Name: pointer.To(d.Get("name").(string)),
		Properties: &virtualwans.HubRouteTableProperties{
			Labels: pluginsdk.ExpandStringSlice(d.Get("labels").(*pluginsdk.Set).List()),
			Routes: expandVirtualHubRouteTableHubRoutes(d.Get("route").(*pluginsdk.Set).List()),
		},
	}

	if err := client.HubRouteTablesCreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceVirtualHubRouteTableRead(d, meta)
}

func expandVirtualHubRouteTableHubRoutes(input []any) *[]virtualwans.HubRoute {
	results := make([]virtualwans.HubRoute, 0)

	for _, item := range input {
		v := item.(map[string]any)

		result := virtualwans.HubRoute{
			Name:            v["name"].(string),
			DestinationType: v["destinations_type"].(string),
			Destinations:    pointer.From(pluginsdk.ExpandStringSlice(v["destinations"].(*pluginsdk.Set).List())),
			NextHopType:     v["next_hop_type"].(string),
			NextHop:         v["next_hop"].(string),
		}

		results = append(results, result)
	}

	return &results
}
