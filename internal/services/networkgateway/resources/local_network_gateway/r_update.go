// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package local_network_gateway

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/localnetworkgateways"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkgateway/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLocalNetworkGatewayUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetworkGateway.LocalNetworkGateways
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := localnetworkgateways.ParseLocalNetworkGatewayID(d.Id())
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

	payload := existing.Model

	if d.HasChange("gateway_address") {
		payload.Properties.GatewayIPAddress = pointer.To(d.Get("gateway_address").(string))
	}

	if d.HasChange("gateway_fqdn") {
		payload.Properties.Fqdn = pointer.To(d.Get("gateway_fqdn").(string))
	}

	if d.HasChange("bgp_settings") {
		payload.Properties.BgpSettings = expandLocalNetworkGatewayBGPSettings(d)
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	pollerType := custompollers.NewLocalNetworkGatewayPoller(client, *id)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)

	if d.HasChange("address_space") {
		payload.Properties.LocalNetworkAddressSpace = expandLocalNetworkGatewayAddressSpaces(d)
	}

	// This requires a custom poller because the API sometimes returns an async operation URL
	// that only returns 404s, causing the provider to poll until timeout when using the `ThenPoll` method.
	if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err := poller.PollUntilDone(ctx); err != nil {
		return err
	}

	return resourceLocalNetworkGatewayRead(d, meta)
}
