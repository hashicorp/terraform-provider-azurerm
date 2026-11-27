// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_connection

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubConnectionRead(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := virtualwans.ParseHubVirtualNetworkConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.HubVirtualNetworkConnectionsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return resourceVirtualHubConnectionFlatten(d, id, resp.Model)
}

func virtualHubConnectionProvisioningStateRefreshFunc(ctx context.Context, client *virtualwans.VirtualWANsClient, id virtualwans.HubVirtualNetworkConnectionId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.HubVirtualNetworkConnectionsGet(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if res.Model == nil {
			return res, "", fmt.Errorf("retrieving %s: `model` was nil", id)
		}
		if res.Model.Properties == nil {
			return res, "", fmt.Errorf("retrieving %s: `properties` was nil", id)
		}
		if res.Model.Properties.ProvisioningState == nil {
			return res, "", fmt.Errorf("retrieving %s: `properties.provisioningState` was nil", id)
		}
		return res, string(*res.Model.Properties.ProvisioningState), nil
	}
}
