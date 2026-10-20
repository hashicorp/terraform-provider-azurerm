// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package firewall_network_rule_collection

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

func resourceFirewallNetworkRuleCollectionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.AzureFirewalls
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallNetworkRuleCollectionID(d.Id())
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
		return fmt.Errorf("retrieving Network Rule Collection %q (Firewall %q / Resource Group %q): `model` was nil", id.NetworkRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	if read.Model.Properties == nil {
		return fmt.Errorf("retrieving Network Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", id.NetworkRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}
	props := *read.Model.Properties

	if props.NetworkRuleCollections == nil {
		return fmt.Errorf("retrieving Network Rule Collection %q (Firewall %q / Resource Group %q): `props.NetworkRuleCollections` was nil", id.NetworkRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	var rule *azurefirewalls.AzureFirewallNetworkRuleCollection
	for _, r := range *props.NetworkRuleCollections {
		if r.Name == nil {
			continue
		}

		if *r.Name == id.NetworkRuleCollectionName {
			rule = &r
			break
		}
	}

	if rule == nil {
		log.Printf("[DEBUG] Network Rule Collection %q was not found on Firewall %q (Resource Group %q) - removing from state!", id.NetworkRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("name", id.NetworkRuleCollectionName)
	d.Set("azure_firewall_name", id.AzureFirewallName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := rule.Properties; props != nil {
		if action := props.Action; action != nil {
			d.Set("action", pointer.FromEnum(action.Type))
		}

		if priority := props.Priority; priority != nil {
			d.Set("priority", int(*priority))
		}

		if err := d.Set("rule", flattenFirewallNetworkRuleCollectionRules(props.Rules)); err != nil {
			return fmt.Errorf("setting `rule`: %+v", err)
		}
	}

	return nil
}

func flattenFirewallNetworkRuleCollectionRules(rules *[]azurefirewalls.AzureFirewallNetworkRule) []any {
	outputs := make([]any, 0)
	if rules == nil {
		return outputs
	}

	for _, rule := range *rules {
		var (
			name            string
			description     string
			sourceAddresses []any
			sourceIPGroups  []any
			destAddresses   []any
			destIPGroups    []any
			destPorts       []any
			destFqdns       []any
		)

		if rule.Name != nil {
			name = *rule.Name
		}
		if rule.Description != nil {
			description = *rule.Description
		}
		if rule.SourceAddresses != nil {
			sourceAddresses = pluginsdk.FlattenSlice(rule.SourceAddresses)
		}
		if rule.SourceIPGroups != nil {
			sourceIPGroups = pluginsdk.FlattenSlice(rule.SourceIPGroups)
		}
		if rule.DestinationAddresses != nil {
			destAddresses = pluginsdk.FlattenSlice(rule.DestinationAddresses)
		}
		if rule.DestinationIPGroups != nil {
			destIPGroups = pluginsdk.FlattenSlice(rule.DestinationIPGroups)
		}
		if rule.DestinationPorts != nil {
			destPorts = pluginsdk.FlattenSlice(rule.DestinationPorts)
		}
		if rule.DestinationFqdns != nil {
			destFqdns = pluginsdk.FlattenSlice(rule.DestinationFqdns)
		}
		protocols := make([]string, 0)
		if rule.Protocols != nil {
			for _, protocol := range *rule.Protocols {
				protocols = append(protocols, string(protocol))
			}
		}
		outputs = append(outputs, map[string]any{
			"name":                  name,
			"description":           description,
			"source_addresses":      sourceAddresses,
			"source_ip_groups":      sourceIPGroups,
			"destination_addresses": destAddresses,
			"destination_ip_groups": destIPGroups,
			"destination_ports":     destPorts,
			"destination_fqdns":     destFqdns,
			"protocols":             pluginsdk.FlattenSlice(&protocols),
		})
	}
	return outputs
}
