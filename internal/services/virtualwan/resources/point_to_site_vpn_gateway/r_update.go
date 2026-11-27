// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package point_to_site_vpn_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePointToSiteVPNGatewayUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualWANP2SVPNGatewayID(d.Id())
	if err != nil {
		return err
	}
	existing, err := client.P2sVpnGatewaysGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	props := virtualwans.P2SVpnGatewayProperties{}
	if existing.Model.Properties != nil {
		props = *existing.Model.Properties
	}

	if d.HasChange("connection_configuration") {
		props.P2SConnectionConfigurations = expandPointToSiteVPNGatewayConnectionConfiguration(d.Get("connection_configuration").([]any))
	}

	if d.HasChange("scale_unit") {
		props.VpnGatewayScaleUnit = pointer.To(int64(d.Get("scale_unit").(int)))
	}

	if d.HasChange("dns_servers") {
		customDNSServers := pluginsdk.ExpandStringSlice(d.Get("dns_servers").([]any))
		if len(*customDNSServers) != 0 {
			props.CustomDnsServers = customDNSServers
		}
	}

	if d.HasChange("tags") {
		existing.Model.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	existing.Model.Properties = &props

	if err := client.P2sVpnGatewaysCreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourcePointToSiteVPNGatewayRead(d, meta)
}
