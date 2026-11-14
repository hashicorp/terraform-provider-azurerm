// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway_nat_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkGatewayNatRuleUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetworkGateway.NetworkGateways
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := networkgateways.ParseVirtualNetworkGatewayNatRuleID(d.Id())
	if err != nil {
		return err
	}

	props := networkgateways.VirtualNetworkGatewayNatRule{
		Name: pointer.To(d.Get("name").(string)),
		Properties: &networkgateways.VirtualNetworkGatewayNatRuleProperties{
			ExternalMappings: expandVirtualNetworkGatewayNatRuleMappings(d.Get("external_mapping").([]any)),
			InternalMappings: expandVirtualNetworkGatewayNatRuleMappings(d.Get("internal_mapping").([]any)),
			Mode:             pointer.ToEnum[networkgateways.VpnNatRuleMode](d.Get("mode").(string)),
			Type:             pointer.ToEnum[networkgateways.VpnNatRuleType](d.Get("type").(string)),
		},
	}

	if v, ok := d.GetOk("ip_configuration_id"); ok {
		props.Properties.IPConfigurationId = pointer.To(v.(string))
	}

	if err := client.VirtualNetworkGatewayNatRulesCreateOrUpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceVirtualNetworkGatewayNatRuleRead(d, meta)
}
