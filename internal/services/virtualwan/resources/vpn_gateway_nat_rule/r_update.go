// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_gateway_nat_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVPNGatewayNatRuleUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := virtualwans.ParseNatRuleID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.NatRulesGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	props := virtualwans.VpnGatewayNatRule{
		Name: pointer.To(d.Get("name").(string)),
		Properties: &virtualwans.VpnGatewayNatRuleProperties{
			Mode:             pointer.ToEnum[virtualwans.VpnNatRuleMode](d.Get("mode").(string)),
			Type:             pointer.ToEnum[virtualwans.VpnNatRuleType](d.Get("type").(string)),
			ExternalMappings: existing.Model.Properties.ExternalMappings,
			InternalMappings: existing.Model.Properties.InternalMappings,
		},
	}

	if ok := d.HasChange("external_mapping"); ok {
		props.Properties.ExternalMappings = expandVpnGatewayNatRuleMappings(d.Get("external_mapping").([]any))
	}

	if ok := d.HasChange("internal_mapping"); ok {
		props.Properties.InternalMappings = expandVpnGatewayNatRuleMappings(d.Get("internal_mapping").([]any))
	}

	if v, ok := d.GetOk("ip_configuration_id"); ok {
		props.Properties.IPConfigurationId = pointer.To(v.(string))
	}

	if err := client.NatRulesCreateOrUpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceVPNGatewayNatRuleRead(d, meta)
}
