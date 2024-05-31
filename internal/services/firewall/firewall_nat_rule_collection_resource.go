// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/azurefirewalls"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceFirewallNatRuleCollection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFirewallNatRuleCollectionCreateUpdate,
		Read:   resourceFirewallNatRuleCollectionRead,
		Update: resourceFirewallNatRuleCollectionCreateUpdate,
		Delete: resourceFirewallNatRuleCollectionDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FirewallNatRuleCollectionID(id)
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
					string(azurefirewalls.AzureFirewallNatRCActionTypeDnat),
					string(azurefirewalls.AzureFirewallNatRCActionTypeSnat),
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
						"translated_address": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"translated_port": {
							Type:     pluginsdk.TypeString,
							Required: true,
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
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"destination_ports": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"protocols": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(azurefirewalls.AzureFirewallNetworkRuleProtocolAny),
									string(azurefirewalls.AzureFirewallNetworkRuleProtocolICMP),
									string(azurefirewalls.AzureFirewallNetworkRuleProtocolTCP),
									string(azurefirewalls.AzureFirewallNetworkRuleProtocolUDP),
								}, false),
							},
						},
					},
				},
			},
		},
	}
}

func resourceFirewallNatRuleCollectionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("retrieving %s: `model` was nil", firewallId)
	}

	if firewall.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `props` was nil", firewallId)
	}
	props := *firewall.Model.Properties

	if props.ApplicationRuleCollections == nil {
		return fmt.Errorf("retrieving %s: `props.ApplicationRuleCollections` was nil", firewallId)
	}

	ruleCollections := *props.NatRuleCollections
	natRules, err := expandFirewallNatRules(d.Get("rule").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding Firewall NAT Rules: %+v", err)
	}
	priority := d.Get("priority").(int)
	newRuleCollection := azurefirewalls.AzureFirewallNatRuleCollection{
		Name: utils.String(name),
		Properties: &azurefirewalls.AzureFirewallNatRuleCollectionProperties{
			Action: &azurefirewalls.AzureFirewallNatRCAction{
				Type: pointer.To(azurefirewalls.AzureFirewallNatRCActionType(d.Get("action").(string))),
			},
			Priority: utils.Int64(int64(priority)),
			Rules:    natRules,
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
			return fmt.Errorf("locating NAT Rule Collection %q (Firewall %q / Resource Group %q)", name, firewallName, resourceGroup)
		}

		ruleCollections[index] = newRuleCollection
	} else {
		if d.IsNewResource() && index != -1 {
			return tf.ImportAsExistsError("azurerm_firewall_nat_rule_collection", id)
		}

		// first double check it doesn't already exist
		ruleCollections = append(ruleCollections, newRuleCollection)
	}

	firewall.Model.Properties.NatRuleCollections = &ruleCollections
	if err = client.CreateOrUpdateThenPoll(ctx, firewallId, *firewall.Model); err != nil {
		return fmt.Errorf("creating/updating NAT Rule Collection %q in Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
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
		if collections := props.NatRuleCollections; collections != nil {
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
		return fmt.Errorf("Cannot find ID for NAT Rule Collection %q (Azure Firewall %q / Resource Group %q)", name, firewallName, resourceGroup)
	}
	d.SetId(collectionID)

	return resourceFirewallNatRuleCollectionRead(d, meta)
}

func resourceFirewallNatRuleCollectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			d.Set("action", string(pointer.From(action.Type)))
		}

		if priority := props.Priority; priority != nil {
			d.Set("priority", int(*priority))
		}

		flattenedRules := flattenFirewallNatRuleCollectionRules(props.Rules)
		if err := d.Set("rule", flattenedRules); err != nil {
			return fmt.Errorf("setting `rule`: %+v", err)
		}
	}

	return nil
}

func resourceFirewallNatRuleCollectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.AzureFirewalls
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallNatRuleCollectionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.AzureFirewallName, AzureFirewallResourceName)
	defer locks.UnlockByName(id.AzureFirewallName, AzureFirewallResourceName)

	firewallId := azurefirewalls.NewAzureFirewallID(id.SubscriptionId, id.ResourceGroup, id.AzureFirewallName)

	firewall, err := client.Get(ctx, firewallId)
	if err != nil {
		if response.WasNotFound(firewall.HttpResponse) {
			// assume deleted
			return nil
		}

		return fmt.Errorf("making Read request on Azure Firewall %q (Resource Group %q): %+v", id.AzureFirewallName, id.ResourceGroup, err)
	}

	if firewall.Model == nil {
		return fmt.Errorf("retrieving NAT Rule Collection %q (Firewall %q / Resource Group %q): `model` was nil", id.NatRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	props := firewall.Model.Properties
	if props == nil {
		return fmt.Errorf("retrieving NAT Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", id.NatRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}
	if props.NetworkRuleCollections == nil {
		return fmt.Errorf("retrieving NAT Rule Collection %q (Firewall %q / Resource Group %q): `props.NatRuleCollections` was nil", id.NatRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	natRules := make([]azurefirewalls.AzureFirewallNatRuleCollection, 0)
	for _, rule := range *props.NatRuleCollections {
		if rule.Name == nil {
			continue
		}

		if *rule.Name != id.NatRuleCollectionName {
			natRules = append(natRules, rule)
		}
	}
	props.NatRuleCollections = &natRules

	if err := client.CreateOrUpdateThenPoll(ctx, firewallId, *firewall.Model); err != nil {
		return fmt.Errorf("deleting NAT Rule Collection %q from Firewall %q (Resource Group %q): %+v", id.NatRuleCollectionName, id.AzureFirewallName, id.ResourceGroup, err)
	}

	return nil
}

func expandFirewallNatRules(input []interface{}) (*[]azurefirewalls.AzureFirewallNatRule, error) {
	rules := make([]azurefirewalls.AzureFirewallNatRule, 0)

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

		destinationPorts := make([]string, 0)
		for _, v := range rule["destination_ports"].([]interface{}) {
			destinationPorts = append(destinationPorts, v.(string))
		}

		translatedAddress := rule["translated_address"].(string)
		translatedPort := rule["translated_port"].(string)

		ruleToAdd := azurefirewalls.AzureFirewallNatRule{
			Name:                 utils.String(name),
			Description:          utils.String(description),
			SourceAddresses:      &sourceAddresses,
			SourceIPGroups:       &sourceIpGroups,
			DestinationAddresses: &destinationAddresses,
			DestinationPorts:     &destinationPorts,
			TranslatedAddress:    &translatedAddress,
			TranslatedPort:       &translatedPort,
		}

		nrProtocols := make([]azurefirewalls.AzureFirewallNetworkRuleProtocol, 0)
		for _, v := range rule["protocols"].([]interface{}) {
			s := azurefirewalls.AzureFirewallNetworkRuleProtocol(v.(string))
			nrProtocols = append(nrProtocols, s)
		}
		ruleToAdd.Protocols = &nrProtocols
		rules = append(rules, ruleToAdd)
	}

	return &rules, nil
}

func flattenFirewallNatRuleCollectionRules(rules *[]azurefirewalls.AzureFirewallNatRule) []interface{} {
	outputs := make([]interface{}, 0)
	if rules == nil {
		return outputs
	}

	for _, rule := range *rules {
		output := make(map[string]interface{})
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
			output["source_addresses"] = utils.FlattenStringSlice(rule.SourceAddresses)
		}
		if rule.SourceIPGroups != nil {
			output["source_ip_groups"] = utils.FlattenStringSlice(rule.SourceIPGroups)
		}
		if rule.DestinationAddresses != nil {
			output["destination_addresses"] = utils.FlattenStringSlice(rule.DestinationAddresses)
		}
		if rule.DestinationPorts != nil {
			output["destination_ports"] = utils.FlattenStringSlice(rule.DestinationPorts)
		}
		protocols := make([]string, 0)
		if rule.Protocols != nil {
			for _, protocol := range *rule.Protocols {
				protocols = append(protocols, string(protocol))
			}
		}
		output["protocols"] = utils.FlattenStringSlice(&protocols)
		outputs = append(outputs, output)
	}
	return outputs
}
