// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package local_network_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/localnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLocalNetworkGatewayCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetworkGateway.LocalNetworkGateways
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := localnetworkgateways.NewLocalNetworkGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_local_network_gateway", id.ID())
		}
	}

	gateway := localnetworkgateways.LocalNetworkGateway{
		Name:     pointer.To(id.LocalNetworkGatewayName),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: localnetworkgateways.LocalNetworkGatewayPropertiesFormat{
			LocalNetworkAddressSpace: expandLocalNetworkGatewayAddressSpaces(d),
			BgpSettings:              expandLocalNetworkGatewayBGPSettings(d),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if ipAddress := d.Get("gateway_address").(string); ipAddress != "" {
		gateway.Properties.GatewayIPAddress = &ipAddress
	} else {
		gateway.Properties.Fqdn = pointer.To(d.Get("gateway_fqdn").(string))
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, gateway, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceLocalNetworkGatewayRead(d, meta)
}

func expandLocalNetworkGatewayBGPSettings(d *pluginsdk.ResourceData) *localnetworkgateways.BgpSettings {
	v, exists := d.GetOk("bgp_settings")
	if !exists {
		return nil
	}

	settings := v.([]any)
	setting := settings[0].(map[string]any)

	bgpSettings := localnetworkgateways.BgpSettings{
		Asn:               pointer.To(int64(setting["asn"].(int))),
		BgpPeeringAddress: pointer.To(setting["bgp_peering_address"].(string)),
		PeerWeight:        pointer.To(int64(setting["peer_weight"].(int))),
	}

	return &bgpSettings
}

func expandLocalNetworkGatewayAddressSpaces(d *pluginsdk.ResourceData) *localnetworkgateways.AddressSpace {
	prefixes := make([]string, 0)

	for _, pref := range d.Get("address_space").(*pluginsdk.Set).List() {
		prefixes = append(prefixes, pref.(string))
	}

	return &localnetworkgateways.AddressSpace{
		AddressPrefixes: &prefixes,
	}
}
