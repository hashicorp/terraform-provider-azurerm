// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package firewall_application_rule_collection

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

func resourceFirewallApplicationRuleCollectionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.AzureFirewalls
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	firewallName := d.Get("azure_firewall_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	applicationRules, err := expandFirewallApplicationRules(d.Get("rule").([]any))
	if err != nil {
		return fmt.Errorf("expanding Firewall Application Rules: %+v", err)
	}

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
	ruleCollections := *props.ApplicationRuleCollections

	priority := d.Get("priority").(int)
	newRuleCollection := azurefirewalls.AzureFirewallApplicationRuleCollection{
		Name: pointer.To(name),
		Properties: &azurefirewalls.AzureFirewallApplicationRuleCollectionPropertiesFormat{
			Action: &azurefirewalls.AzureFirewallRCAction{
				Type: pointer.ToEnum[azurefirewalls.AzureFirewallRCActionType](d.Get("action").(string)),
			},
			Priority: pointer.To(int64(priority)),
			Rules:    applicationRules,
		},
	}

	index := -1
	var id string
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
			return fmt.Errorf("locating Application Rule Collection %q (Firewall %q / Resource Group %q)", name, firewallName, resourceGroup)
		}

		ruleCollections[index] = newRuleCollection
	} else {
		exists := false
		if d.IsNewResource() && index != -1 {
			exists = true
			if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				return tf.ImportAsExistsError("azurerm_firewall_application_rule_collection", id)
			}
		}

		if !exists {
			ruleCollections = append(ruleCollections, newRuleCollection)
		}
	}

	firewall.Model.Properties.ApplicationRuleCollections = &ruleCollections

	// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
	if err = client.CreateOrUpdateThenPoll(ctx, firewallId, *firewall.Model, azurefirewalls.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating Application Rule Collection %q in %s: %+v", name, firewallId, err)
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
		if collections := props.ApplicationRuleCollections; collections != nil {
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
		return fmt.Errorf("cannot find ID for Application Rule Collection %q (Azure Firewall %q / Resource Group %q)", name, firewallName, resourceGroup)
	}
	d.SetId(collectionID)

	return resourceFirewallApplicationRuleCollectionRead(d, meta)
}

func expandFirewallApplicationRules(inputs []any) (*[]azurefirewalls.AzureFirewallApplicationRule, error) {
	outputs := make([]azurefirewalls.AzureFirewallApplicationRule, 0)

	for _, input := range inputs {
		rule := input.(map[string]any)

		ruleName := rule["name"].(string)
		ruleDescription := rule["description"].(string)
		ruleSourceAddresses := rule["source_addresses"].([]any)
		ruleSourceIpGroups := rule["source_ip_groups"].([]any)
		ruleFqdnTags := rule["fqdn_tags"].([]any)
		ruleTargetFqdns := rule["target_fqdns"].([]any)

		output := azurefirewalls.AzureFirewallApplicationRule{
			Name:            pointer.To(ruleName),
			Description:     pointer.To(ruleDescription),
			SourceAddresses: pluginsdk.ExpandStringSlice(ruleSourceAddresses),
			SourceIPGroups:  pluginsdk.ExpandStringSlice(ruleSourceIpGroups),
			FqdnTags:        pluginsdk.ExpandStringSlice(ruleFqdnTags),
			TargetFqdns:     pluginsdk.ExpandStringSlice(ruleTargetFqdns),
		}

		ruleProtocols := make([]azurefirewalls.AzureFirewallApplicationRuleProtocol, 0)
		for _, v := range rule["protocol"].([]any) {
			protocol := v.(map[string]any)
			port := protocol["port"].(int)
			ruleProtocol := azurefirewalls.AzureFirewallApplicationRuleProtocol{
				Port:         pointer.To(int64(port)),
				ProtocolType: pointer.ToEnum[azurefirewalls.AzureFirewallApplicationRuleProtocolType](protocol["type"].(string)),
			}
			ruleProtocols = append(ruleProtocols, ruleProtocol)
		}
		output.Protocols = &ruleProtocols
		if len(*output.FqdnTags) > 0 {
			if len(*output.TargetFqdns) > 0 || len(*output.Protocols) > 0 {
				return nil, fmt.Errorf("`fqdn_tags` cannot be used with `target_fqdns` or `protocol`")
			}
		}

		if len(*output.SourceAddresses) == 0 && len(*output.SourceIPGroups) == 0 {
			return nil, fmt.Errorf("at least one of %q and %q must be specified for each rule", "source_addresses", "source_ip_groups")
		}
		outputs = append(outputs, output)
	}

	return &outputs, nil
}
