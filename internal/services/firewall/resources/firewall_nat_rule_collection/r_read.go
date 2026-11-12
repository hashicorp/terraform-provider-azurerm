// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package firewall_nat_rule_collection

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFirewallNatRuleCollectionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.AzureFirewalls
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallNatRuleCollectionID(d.Id())
	if err != nil {
		return err
	}

	firewallId := azurefirewalls.NewAzureFirewallID(id.SubscriptionId, id.ResourceGroup, id.AzureFirewallName)

	read, err := client.Get(ctx, firewallId)
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[DEBUG] Azure Firewall %s  was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Azure Firewall %s : %+v", *id, err)
	}

	if read.Model == nil {
		return fmt.Errorf("retrieving NAT Rule Collection %q (Firewall %q / Resource Group %q): `model` was nil", id.NatRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	if read.Model.Properties == nil {
		return fmt.Errorf("retrieving NAT Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", id.NatRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	props := *read.Model.Properties

	if props.NatRuleCollections == nil {
		return fmt.Errorf("retrieving NAT Rule Collection %q (Firewall %q / Resource Group %q): `props.NetworkRuleCollections` was nil", id.NatRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	var rule *azurefirewalls.AzureFirewallNatRuleCollection
	for _, r := range *props.NatRuleCollections {
		if r.Name == nil {
			continue
		}

		if *r.Name == id.NatRuleCollectionName {
			rule = &r
			break
		}
	}

	if rule == nil {
		log.Printf("[DEBUG] NAT Rule Collection %q was not found on Firewall %q (Resource Group %q) - removing from state!", id.NatRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("name", id.NatRuleCollectionName)
	d.Set("azure_firewall_name", id.AzureFirewallName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := rule.Properties; props != nil {
		if action := props.Action; action != nil {
			d.Set("action", pointer.FromEnum(action.Type))
		}

		if priority := props.Priority; priority != nil {
			d.Set("priority", int(*priority))
		}

		if err := d.Set("rule", flattenFirewallNatRuleCollectionRules(props.Rules)); err != nil {
			return fmt.Errorf("setting `rule`: %+v", err)
		}
	}

	return nil
}

func flattenFirewallNatRuleCollectionRules(rules *[]azurefirewalls.AzureFirewallNatRule) []any {
	outputs := make([]any, 0)
	if rules == nil {
		return outputs
	}

	for _, rule := range *rules {
		output := make(map[string]any)
		if rule.Name != nil {
			output["name"] = *rule.Name
		}
		if rule.Description != nil {
			output["description"] = *rule.Description
		}
		if rule.TranslatedAddress != nil {
			output["translated_address"] = *rule.TranslatedAddress
		}
		if rule.TranslatedPort != nil {
			output["translated_port"] = *rule.TranslatedPort
		}
		if rule.SourceAddresses != nil {
			output["source_addresses"] = pluginsdk.FlattenSlice(rule.SourceAddresses)
		}
		if rule.SourceIPGroups != nil {
			output["source_ip_groups"] = pluginsdk.FlattenSlice(rule.SourceIPGroups)
		}
		if rule.DestinationAddresses != nil {
			output["destination_addresses"] = pluginsdk.FlattenSlice(rule.DestinationAddresses)
		}
		if rule.DestinationPorts != nil {
			output["destination_ports"] = pluginsdk.FlattenSlice(rule.DestinationPorts)
		}
		protocols := make([]string, 0)
		if rule.Protocols != nil {
			for _, protocol := range *rule.Protocols {
				protocols = append(protocols, string(protocol))
			}
		}
		output["protocols"] = pluginsdk.FlattenSlice(&protocols)
		outputs = append(outputs, output)
	}
	return outputs
}
