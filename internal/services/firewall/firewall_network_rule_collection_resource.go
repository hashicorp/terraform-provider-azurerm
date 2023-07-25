// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func resourceFirewallNetworkRuleCollection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFirewallNetworkRuleCollectionCreateUpdate,
		Read:   resourceFirewallNetworkRuleCollectionRead,
		Update: resourceFirewallNetworkRuleCollectionCreateUpdate,
		Delete: resourceFirewallNetworkRuleCollectionDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FirewallNetworkRuleCollectionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FirewallName,
			},

			"azure_firewall_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FirewallName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"priority": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(100, 65000),
			},

			"action": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.AzureFirewallRCActionTypeAllow),
					string(network.AzureFirewallRCActionTypeDeny),
				}, false),
			},

			"rule": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"description": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"source_addresses": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"source_ip_groups": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"destination_addresses": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"destination_ip_groups": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"destination_ports": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"destination_fqdns": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"protocols": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(network.AzureFirewallNetworkRuleProtocolAny),
									string(network.AzureFirewallNetworkRuleProtocolICMP),
									string(network.AzureFirewallNetworkRuleProtocolTCP),
									string(network.AzureFirewallNetworkRuleProtocolUDP),
								}, false),
							},
						},
					},
				},
			},
		},
	}
}

func resourceFirewallNetworkRuleCollectionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	firewallName := d.Get("azure_firewall_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(firewallName, AzureFirewallResourceName)
	defer locks.UnlockByName(firewallName, AzureFirewallResourceName)

	firewall, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		return fmt.Errorf("retrieving Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	if firewall.AzureFirewallPropertiesFormat == nil {
		return fmt.Errorf("expanding Firewall %q (Resource Group %q): `properties` was nil.", firewallName, resourceGroup)
	}
	props := *firewall.AzureFirewallPropertiesFormat

	if props.NetworkRuleCollections == nil {
		return fmt.Errorf("expanding Firewall %q (Resource Group %q): `properties.NetworkRuleCollections` was nil.", firewallName, resourceGroup)
	}
	ruleCollections := *props.NetworkRuleCollections

	networkRules, err := expandFirewallNetworkRules(d.Get("rule").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding Firewall Network Rules: %+v", err)
	}
	priority := d.Get("priority").(int)
	newRuleCollection := network.AzureFirewallNetworkRuleCollection{
		Name: utils.String(name),
		AzureFirewallNetworkRuleCollectionPropertiesFormat: &network.AzureFirewallNetworkRuleCollectionPropertiesFormat{
			Action: &network.AzureFirewallRCAction{
				Type: network.AzureFirewallRCActionType(d.Get("action").(string)),
			},
			Priority: utils.Int32(int32(priority)),
			Rules:    networkRules,
		},
	}

	index := -1
	var id string
	// determine if this already exists
	for i, v := range ruleCollections {
		if v.Name == nil || v.ID == nil {
			continue
		}

		if *v.Name == name {
			index = i
			id = *v.ID
			break
		}
	}

	if !d.IsNewResource() {
		if index == -1 {
			return fmt.Errorf("locating Network Rule Collection %q (Firewall %q / Resource Group %q)", name, firewallName, resourceGroup)
		}

		ruleCollections[index] = newRuleCollection
	} else {
		if d.IsNewResource() && index != -1 {
			return tf.ImportAsExistsError("azurerm_firewall_network_rule_collection", id)
		}

		// first double check it doesn't already exist
		ruleCollections = append(ruleCollections, newRuleCollection)
	}

	firewall.AzureFirewallPropertiesFormat.NetworkRuleCollections = &ruleCollections

	future, err := client.CreateOrUpdate(ctx, resourceGroup, firewallName, firewall)
	if err != nil {
		return fmt.Errorf("creating/updating Network Rule Collection %q in Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of Network Rule Collection %q of Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		return fmt.Errorf("retrieving Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	var collectionID string
	if props := read.AzureFirewallPropertiesFormat; props != nil {
		if collections := props.NetworkRuleCollections; collections != nil {
			for _, collection := range *collections {
				if collection.Name == nil {
					continue
				}

				if *collection.Name == name {
					collectionID = *collection.ID
					break
				}
			}
		}
	}

	if collectionID == "" {
		return fmt.Errorf("Cannot find ID for Network Rule Collection %q (Azure Firewall %q / Resource Group %q)", name, firewallName, resourceGroup)
	}
	d.SetId(collectionID)

	return resourceFirewallNetworkRuleCollectionRead(d, meta)
}

func resourceFirewallNetworkRuleCollectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallNetworkRuleCollectionID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.AzureFirewallName)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Azure Firewall %s  was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Azure Firewall %s : %+v", *id, err)
	}

	if read.AzureFirewallPropertiesFormat == nil {
		return fmt.Errorf("retrieving Network Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", id.NetworkRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}
	props := *read.AzureFirewallPropertiesFormat

	if props.NetworkRuleCollections == nil {
		return fmt.Errorf("retrieving Network Rule Collection %q (Firewall %q / Resource Group %q): `props.NetworkRuleCollections` was nil", id.NetworkRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	var rule *network.AzureFirewallNetworkRuleCollection
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

	if props := rule.AzureFirewallNetworkRuleCollectionPropertiesFormat; props != nil {
		if action := props.Action; action != nil {
			d.Set("action", string(action.Type))
		}

		if priority := props.Priority; priority != nil {
			d.Set("priority", int(*priority))
		}

		flattenedRules := flattenFirewallNetworkRuleCollectionRules(props.Rules)
		if err := d.Set("rule", flattenedRules); err != nil {
			return fmt.Errorf("setting `rule`: %+v", err)
		}
	}

	return nil
}

func resourceFirewallNetworkRuleCollectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallNetworkRuleCollectionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.AzureFirewallName, AzureFirewallResourceName)
	defer locks.UnlockByName(id.AzureFirewallName, AzureFirewallResourceName)

	firewall, err := client.Get(ctx, id.ResourceGroup, id.AzureFirewallName)
	if err != nil {
		if utils.ResponseWasNotFound(firewall.Response) {
			// assume deleted
			return nil
		}

		return fmt.Errorf("making Read request on Azure Firewall %q (Resource Group %q): %+v", id.AzureFirewallName, id.ResourceGroup, err)
	}

	props := firewall.AzureFirewallPropertiesFormat
	if props == nil {
		return fmt.Errorf("retrieving Network Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", id.NetworkRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}
	if props.NetworkRuleCollections == nil {
		return fmt.Errorf("retrieving Network Rule Collection %q (Firewall %q / Resource Group %q): `props.NetworkRuleCollections` was nil", id.NetworkRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	networkRules := make([]network.AzureFirewallNetworkRuleCollection, 0)
	for _, rule := range *props.NetworkRuleCollections {
		if rule.Name == nil {
			continue
		}

		if *rule.Name != id.NetworkRuleCollectionName {
			networkRules = append(networkRules, rule)
		}
	}
	props.NetworkRuleCollections = &networkRules

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.AzureFirewallName, firewall)
	if err != nil {
		return fmt.Errorf("deleting Network Rule Collection %q from Firewall %q (Resource Group %q): %+v", id.NetworkRuleCollectionName, id.AzureFirewallName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Network Rule Collection %q from Firewall %q (Resource Group %q): %+v", id.NetworkRuleCollectionName, id.AzureFirewallName, id.ResourceGroup, err)
	}

	return nil
}

func expandFirewallNetworkRules(input []interface{}) (*[]network.AzureFirewallNetworkRule, error) {
	rules := make([]network.AzureFirewallNetworkRule, 0)

	for _, nwRule := range input {
		rule := nwRule.(map[string]interface{})

		name := rule["name"].(string)
		description := rule["description"].(string)

		sourceAddresses := make([]string, 0)
		for _, v := range rule["source_addresses"].([]interface{}) {
			sourceAddresses = append(sourceAddresses, v.(string))
		}

		sourceIpGroups := make([]string, 0)
		for _, v := range rule["source_ip_groups"].([]interface{}) {
			sourceIpGroups = append(sourceIpGroups, v.(string))
		}

		if len(sourceAddresses) == 0 && len(sourceIpGroups) == 0 {
			return nil, fmt.Errorf("at least one of %q and %q must be specified for each rule", "source_addresses", "source_ip_groups")
		}

		destinationAddresses := make([]string, 0)
		for _, v := range rule["destination_addresses"].([]interface{}) {
			destinationAddresses = append(destinationAddresses, v.(string))
		}

		destinationIpGroups := make([]string, 0)
		for _, v := range rule["destination_ip_groups"].([]interface{}) {
			destinationIpGroups = append(destinationIpGroups, v.(string))
		}

		destinationFqdns := make([]string, 0)
		for _, v := range rule["destination_fqdns"].([]interface{}) {
			destinationFqdns = append(destinationFqdns, v.(string))
		}

		if len(destinationAddresses) == 0 && len(destinationIpGroups) == 0 && len(destinationFqdns) == 0 {
			return nil, fmt.Errorf("at least one of %q, %q and %q must be specified for each rule", "destination_addresses", "destination_ip_groups", "destination_fqdns")
		}

		destinationPorts := make([]string, 0)
		for _, v := range rule["destination_ports"].([]interface{}) {
			destinationPorts = append(destinationPorts, v.(string))
		}

		ruleToAdd := network.AzureFirewallNetworkRule{
			Name:                 utils.String(name),
			Description:          utils.String(description),
			SourceAddresses:      &sourceAddresses,
			SourceIPGroups:       &sourceIpGroups,
			DestinationAddresses: &destinationAddresses,
			DestinationIPGroups:  &destinationIpGroups,
			DestinationPorts:     &destinationPorts,
			DestinationFqdns:     &destinationFqdns,
		}

		nrProtocols := make([]network.AzureFirewallNetworkRuleProtocol, 0)
		for _, v := range rule["protocols"].([]interface{}) {
			s := network.AzureFirewallNetworkRuleProtocol(v.(string))
			nrProtocols = append(nrProtocols, s)
		}
		ruleToAdd.Protocols = &nrProtocols
		rules = append(rules, ruleToAdd)
	}

	return &rules, nil
}

func flattenFirewallNetworkRuleCollectionRules(rules *[]network.AzureFirewallNetworkRule) []interface{} {
	outputs := make([]interface{}, 0)
	if rules == nil {
		return outputs
	}

	for _, rule := range *rules {
		var (
			name            string
			description     string
			sourceAddresses []interface{}
			sourceIPGroups  []interface{}
			destAddresses   []interface{}
			destIPGroups    []interface{}
			destPorts       []interface{}
			destFqdns       []interface{}
		)

		if rule.Name != nil {
			name = *rule.Name
		}
		if rule.Description != nil {
			description = *rule.Description
		}
		if rule.SourceAddresses != nil {
			sourceAddresses = utils.FlattenStringSlice(rule.SourceAddresses)
		}
		if rule.SourceIPGroups != nil {
			sourceIPGroups = utils.FlattenStringSlice(rule.SourceIPGroups)
		}
		if rule.DestinationAddresses != nil {
			destAddresses = utils.FlattenStringSlice(rule.DestinationAddresses)
		}
		if rule.DestinationIPGroups != nil {
			destIPGroups = utils.FlattenStringSlice(rule.DestinationIPGroups)
		}
		if rule.DestinationPorts != nil {
			destPorts = utils.FlattenStringSlice(rule.DestinationPorts)
		}
		if rule.DestinationFqdns != nil {
			destFqdns = utils.FlattenStringSlice(rule.DestinationFqdns)
		}
		protocols := make([]string, 0)
		if rule.Protocols != nil {
			for _, protocol := range *rule.Protocols {
				protocols = append(protocols, string(protocol))
			}
		}
		outputs = append(outputs, map[string]interface{}{
			"name":                  name,
			"description":           description,
			"source_addresses":      sourceAddresses,
			"source_ip_groups":      sourceIPGroups,
			"destination_addresses": destAddresses,
			"destination_ip_groups": destIPGroups,
			"destination_ports":     destPorts,
			"destination_fqdns":     destFqdns,
			"protocols":             utils.FlattenStringSlice(&protocols),
		})
	}
	return outputs
}
