// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_gateway_nat_rule

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVPNGatewayNatRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := virtualwans.ParseNatRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.NatRulesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.NatRuleName)

	gatewayId := virtualwans.NewVpnGatewayID(id.SubscriptionId, id.ResourceGroupName, id.VpnGatewayName)
	d.Set("vpn_gateway_id", gatewayId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("ip_configuration_id", props.IPConfigurationId)
			d.Set("mode", pointer.From(props.Mode))
			d.Set("type", pointer.From(props.Type))

			if err := d.Set("external_mapping", flattenVpnGatewayNatRuleMappings(props.ExternalMappings)); err != nil {
				return fmt.Errorf("setting `external_mapping`: %+v", err)
			}

			if err := d.Set("internal_mapping", flattenVpnGatewayNatRuleMappings(props.InternalMappings)); err != nil {
				return fmt.Errorf("setting `internal_mapping`: %+v", err)
			}
		}
	}

	return nil
}

func flattenVpnGatewayNatRuleMappings(input *[]virtualwans.VpnNatRuleMapping) []any {
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
