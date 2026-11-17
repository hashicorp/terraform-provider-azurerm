// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package firewall_network_rule_collection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFirewallNetworkRuleCollectionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.AzureFirewalls
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	firewallName := d.Get("azure_firewall_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(firewallName, AzureFirewallResourceName)
	defer locks.UnlockByName(firewallName, AzureFirewallResourceName)

	firewallId := azurefirewalls.NewAzureFirewallID(subscriptionId, resourceGroup, firewallName)

	firewall, err := client.Get(ctx, firewallId)
	if err != nil {
		return fmt.Errorf("retrieving Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	if firewall.Model == nil {
		return fmt.Errorf("expanding Firewall %q (Resource Group %q): `model` was nil", firewallName, resourceGroup)
	}

	if firewall.Model.Properties == nil {
		return fmt.Errorf("expanding Firewall %q (Resource Group %q): `properties` was nil", firewallName, resourceGroup)
	}
	props := *firewall.Model.Properties

	if props.NetworkRuleCollections == nil {
		return fmt.Errorf("expanding Firewall %q (Resource Group %q): `properties.NetworkRuleCollections` was nil", firewallName, resourceGroup)
	}
	ruleCollections := *props.NetworkRuleCollections

	networkRules, err := expandFirewallNetworkRules(d.Get("rule").([]any))
	if err != nil {
		return fmt.Errorf("expanding Firewall Network Rules: %+v", err)
	}
	priority := d.Get("priority").(int)
	newRuleCollection := azurefirewalls.AzureFirewallNetworkRuleCollection{
		Name: pointer.To(name),
		Properties: &azurefirewalls.AzureFirewallNetworkRuleCollectionPropertiesFormat{
			Action: &azurefirewalls.AzureFirewallRCAction{
				Type: pointer.ToEnum[azurefirewalls.AzureFirewallRCActionType](d.Get("action").(string)),
			},
			Priority: pointer.To(int64(priority)),
			Rules:    networkRules,
		},
	}

	index := -1
	var id string
	// determine if this already exists
	for i, v := range ruleCollections {
		if v.Name == nil || v.Id == nil {
			continue
		}

		if *v.Name == name {
			index = i
			id = *v.Id
			break
		}
	}

	if !d.IsNewResource() {
		if index == -1 {
			return fmt.Errorf("locating Network Rule Collection %q (Firewall %q / Resource Group %q)", name, firewallName, resourceGroup)
		}

		ruleCollections[index] = newRuleCollection
	} else {
		exists := false
		if d.IsNewResource() && index != -1 {
			exists = true
			if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				return tf.ImportAsExistsError("azurerm_firewall_network_rule_collection", id)
			}
		}

		// first double check it doesn't already exist
		if !exists {
			ruleCollections = append(ruleCollections, newRuleCollection)
		}
	}

	firewall.Model.Properties.NetworkRuleCollections = &ruleCollections

	// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
	if err := client.CreateOrUpdateThenPoll(ctx, firewallId, *firewall.Model, azurefirewalls.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating Network Rule Collection %q in Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	read, err := client.Get(ctx, firewallId)
	if err != nil {
		return fmt.Errorf("retrieving Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", firewallId)
	}

	var collectionID string
	if props := read.Model.Properties; props != nil {
		if collections := props.NetworkRuleCollections; collections != nil {
			for _, collection := range *collections {
				if collection.Name == nil {
					continue
				}

				if *collection.Name == name {
					collectionID = *collection.Id
					break
				}
			}
		}
	}

	if collectionID == "" {
		return fmt.Errorf("cannot find ID for Network Rule Collection %q (Azure Firewall %q / Resource Group %q)", name, firewallName, resourceGroup)
	}
	d.SetId(collectionID)

	return resourceFirewallNetworkRuleCollectionRead(d, meta)
}

func expandFirewallNetworkRules(input []any) (*[]azurefirewalls.AzureFirewallNetworkRule, error) {
	rules := make([]azurefirewalls.AzureFirewallNetworkRule, 0)

	for _, nwRule := range input {
		rule := nwRule.(map[string]any)

		name := rule["name"].(string)
		description := rule["description"].(string)

		sourceAddresses := make([]string, 0)
		for _, v := range rule["source_addresses"].([]any) {
			sourceAddresses = append(sourceAddresses, v.(string))
		}

		sourceIpGroups := make([]string, 0)
		for _, v := range rule["source_ip_groups"].([]any) {
			sourceIpGroups = append(sourceIpGroups, v.(string))
		}

		if len(sourceAddresses) == 0 && len(sourceIpGroups) == 0 {
			return nil, fmt.Errorf("at least one of %q and %q must be specified for each rule", "source_addresses", "source_ip_groups")
		}

		destinationAddresses := make([]string, 0)
		for _, v := range rule["destination_addresses"].([]any) {
			destinationAddresses = append(destinationAddresses, v.(string))
		}

		destinationIpGroups := make([]string, 0)
		for _, v := range rule["destination_ip_groups"].([]any) {
			destinationIpGroups = append(destinationIpGroups, v.(string))
		}

		destinationFqdns := make([]string, 0)
		for _, v := range rule["destination_fqdns"].([]any) {
			destinationFqdns = append(destinationFqdns, v.(string))
		}

		if len(destinationAddresses) == 0 && len(destinationIpGroups) == 0 && len(destinationFqdns) == 0 {
			return nil, fmt.Errorf("at least one of %q, %q and %q must be specified for each rule", "destination_addresses", "destination_ip_groups", "destination_fqdns")
		}

		destinationPorts := make([]string, 0)
		for _, v := range rule["destination_ports"].([]any) {
			destinationPorts = append(destinationPorts, v.(string))
		}

		ruleToAdd := azurefirewalls.AzureFirewallNetworkRule{
			Name:                 pointer.To(name),
			Description:          pointer.To(description),
			SourceAddresses:      &sourceAddresses,
			SourceIPGroups:       &sourceIpGroups,
			DestinationAddresses: &destinationAddresses,
			DestinationIPGroups:  &destinationIpGroups,
			DestinationPorts:     &destinationPorts,
			DestinationFqdns:     &destinationFqdns,
		}

		nrProtocols := make([]azurefirewalls.AzureFirewallNetworkRuleProtocol, 0)
		for _, v := range rule["protocols"].([]any) {
			s := azurefirewalls.AzureFirewallNetworkRuleProtocol(v.(string))
			nrProtocols = append(nrProtocols, s)
		}
		ruleToAdd.Protocols = &nrProtocols
		rules = append(rules, ruleToAdd)
	}

	return &rules, nil
}
