// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package local_network_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/localnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	networkgatewayClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/networkgateway/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLocalNetworkGatewayRead(d *pluginsdk.ResourceData, meta any) error {
	client := networkgatewayClient.From(meta).LocalNetworkGateways
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := localnetworkgateways.ParseLocalNetworkGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.LocalNetworkGatewayName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		props := model.Properties
		d.Set("gateway_address", props.GatewayIPAddress)
		d.Set("gateway_fqdn", props.Fqdn)

		if lnas := props.LocalNetworkAddressSpace; lnas != nil {
			d.Set("address_space", lnas.AddressPrefixes)
		}
		if err := d.Set("bgp_settings", flattenLocalNetworkGatewayBGPSettings(props.BgpSettings)); err != nil {
			return err
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenLocalNetworkGatewayBGPSettings(input *localnetworkgateways.BgpSettings) []any {
	output := make(map[string]any)

	if input == nil {
		return []any{}
	}

	output["asn"] = int(*input.Asn)
	output["bgp_peering_address"] = *input.BgpPeeringAddress
	output["peer_weight"] = int(*input.PeerWeight)

	return []any{output}
}
