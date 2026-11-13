// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway_nat_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkGatewayNatRuleCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).NetworkGateway.NetworkGateways
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vnetGatewayId, err := virtualnetworkgateways.ParseVirtualNetworkGatewayID(d.Get("virtual_network_gateway_id").(string))
	if err != nil {
		return err
	}

	id := networkgateways.NewVirtualNetworkGatewayNatRuleID(subscriptionId, d.Get("resource_group_name").(string), vnetGatewayId.VirtualNetworkGatewayName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.VirtualNetworkGatewayNatRulesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_network_gateway_nat_rule", id.ID())
		}
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

	if err := client.VirtualNetworkGatewayNatRulesCreateOrUpdateCallbackThenPoll(ctx, id, props, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVirtualNetworkGatewayNatRuleRead(d, meta)
}

func expandVirtualNetworkGatewayNatRuleMappings(input []any) *[]networkgateways.VpnNatRuleMapping {
	results := make([]networkgateways.VpnNatRuleMapping, 0)

	for _, item := range input {
		v := item.(map[string]any)

		result := networkgateways.VpnNatRuleMapping{
			AddressSpace: pointer.To(v["address_space"].(string)),
		}

		if portRange := v["port_range"].(string); portRange != "" {
			result.PortRange = pointer.To(portRange)
		}

		results = append(results, result)
	}

	return &results
}
