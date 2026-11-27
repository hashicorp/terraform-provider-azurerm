// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_gateway_nat_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVPNGatewayNatRuleCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := account.From(meta).SubscriptionId
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	vpnGatewayId, err := virtualwans.ParseVpnGatewayID(d.Get("vpn_gateway_id").(string))
	if err != nil {
		return err
	}

	id := virtualwans.NewNatRuleID(subscriptionId, vpnGatewayId.ResourceGroupName, vpnGatewayId.VpnGatewayName, d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.NatRulesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_vpn_gateway_nat_rule", id.ID())
		}
	}

	props := virtualwans.VpnGatewayNatRule{
		Name: pointer.To(d.Get("name").(string)),
		Properties: &virtualwans.VpnGatewayNatRuleProperties{
			Mode: pointer.ToEnum[virtualwans.VpnNatRuleMode](d.Get("mode").(string)),
			Type: pointer.ToEnum[virtualwans.VpnNatRuleType](d.Get("type").(string)),
		},
	}

	if v, ok := d.GetOk("external_mapping"); ok {
		props.Properties.ExternalMappings = expandVpnGatewayNatRuleMappings(v.([]any))
	}

	if v, ok := d.GetOk("internal_mapping"); ok {
		props.Properties.InternalMappings = expandVpnGatewayNatRuleMappings(v.([]any))
	}

	if v, ok := d.GetOk("ip_configuration_id"); ok {
		props.Properties.IPConfigurationId = pointer.To(v.(string))
	}

	if err := client.NatRulesCreateOrUpdateCallbackThenPoll(ctx, id, props, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVPNGatewayNatRuleRead(d, meta)
}

func expandVpnGatewayNatRuleMappings(input []any) *[]virtualwans.VpnNatRuleMapping {
	results := make([]virtualwans.VpnNatRuleMapping, 0)

	for _, item := range input {
		v := item.(map[string]any)

		result := virtualwans.VpnNatRuleMapping{
			AddressSpace: pointer.To(v["address_space"].(string)),
		}

		if portRange := v["port_range"].(string); portRange != "" {
			result.PortRange = pointer.To(portRange)
		}

		results = append(results, result)
	}

	return &results
}
