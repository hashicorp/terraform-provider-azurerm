// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway_nat_rule

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkGatewayNatRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetworkGateway.NetworkGateways
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := networkgateways.ParseVirtualNetworkGatewayNatRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.VirtualNetworkGatewayNatRulesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.NatRuleName)
	d.Set("resource_group_name", id.ResourceGroupName)

	vnetGatewayId := virtualnetworkgateways.NewVirtualNetworkGatewayID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkGatewayName)
	d.Set("virtual_network_gateway_id", vnetGatewayId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if err := d.Set("external_mapping", flattenVirtualNetworkGatewayNatRuleMappings(props.ExternalMappings)); err != nil {
				return fmt.Errorf("setting `external_mapping`: %+v", err)
			}

			if err := d.Set("internal_mapping", flattenVirtualNetworkGatewayNatRuleMappings(props.InternalMappings)); err != nil {
				return fmt.Errorf("setting `internal_mapping`: %+v", err)
			}

			d.Set("ip_configuration_id", props.IPConfigurationId)
			d.Set("mode", pointer.FromEnum(props.Mode))
			d.Set("type", pointer.FromEnum(props.Type))
		}
	}
	return nil
}

func flattenVirtualNetworkGatewayNatRuleMappings(input *[]networkgateways.VpnNatRuleMapping) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]any{
			"address_space": pointer.From(item.AddressSpace),
			"port_range":    pointer.From(item.PortRange),
		})
	}

	return results
}
