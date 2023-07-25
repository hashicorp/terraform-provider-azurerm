// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	firewallValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func resourceFirewallApplicationRuleCollection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFirewallApplicationRuleCollectionCreateUpdate,
		Read:   resourceFirewallApplicationRuleCollectionRead,
		Update: resourceFirewallApplicationRuleCollectionCreateUpdate,
		Delete: resourceFirewallApplicationRuleCollectionDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FirewallApplicationRuleCollectionID(id)
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
				ValidateFunc: firewallValidate.FirewallName,
			},

			"azure_firewall_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: firewallValidate.FirewallName,
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
							ValidateFunc: validation.NoZeroValues,
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
						"description": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"fqdn_tags": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"target_fqdns": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"protocol": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"type": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(network.AzureFirewallApplicationRuleProtocolTypeHTTP),
											string(network.AzureFirewallApplicationRuleProtocolTypeHTTPS),
											string(network.AzureFirewallApplicationRuleProtocolTypeMssql),
										}, false),
									},
									"port": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validate.PortNumber,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceFirewallApplicationRuleCollectionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	firewallName := d.Get("azure_firewall_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	applicationRules, err := expandFirewallApplicationRules(d.Get("rule").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding Firewall Application Rules: %+v", err)
	}

	locks.ByName(firewallName, AzureFirewallResourceName)
	defer locks.UnlockByName(firewallName, AzureFirewallResourceName)

	firewall, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		return fmt.Errorf("retrieving Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	if firewall.AzureFirewallPropertiesFormat == nil {
		return fmt.Errorf("retrieving Application Rule Collections (Firewall %q / Resource Group %q): `properties` was nil", firewallName, resourceGroup)
	}
	props := *firewall.AzureFirewallPropertiesFormat

	if props.ApplicationRuleCollections == nil {
		return fmt.Errorf("retrieving Application Rule Collections (Firewall %q / Resource Group %q): `properties.ApplicationRuleCollections` was nil", firewallName, resourceGroup)
	}
	ruleCollections := *props.ApplicationRuleCollections

	priority := d.Get("priority").(int)
	newRuleCollection := network.AzureFirewallApplicationRuleCollection{
		Name: utils.String(name),
		AzureFirewallApplicationRuleCollectionPropertiesFormat: &network.AzureFirewallApplicationRuleCollectionPropertiesFormat{
			Action: &network.AzureFirewallRCAction{
				Type: network.AzureFirewallRCActionType(d.Get("action").(string)),
			},
			Priority: utils.Int32(int32(priority)),
			Rules:    applicationRules,
		},
	}

	index := -1
	var id string
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
			return fmt.Errorf("locating Application Rule Collection %q (Firewall %q / Resource Group %q)", name, firewallName, resourceGroup)
		}

		ruleCollections[index] = newRuleCollection
	} else {
		if d.IsNewResource() && index != -1 {
			return tf.ImportAsExistsError("azurerm_firewall_application_rule_collection", id)
		}

		ruleCollections = append(ruleCollections, newRuleCollection)
	}

	firewall.AzureFirewallPropertiesFormat.ApplicationRuleCollections = &ruleCollections

	future, err := client.CreateOrUpdate(ctx, resourceGroup, firewallName, firewall)
	if err != nil {
		return fmt.Errorf("creating/updating Application Rule Collection %q in Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of Application Rule Collection %q of Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		return fmt.Errorf("retrieving Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	var collectionID string
	if props := read.AzureFirewallPropertiesFormat; props != nil {
		if collections := props.ApplicationRuleCollections; collections != nil {
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
		return fmt.Errorf("Cannot find ID for Application Rule Collection %q (Azure Firewall %q / Resource Group %q)", name, firewallName, resourceGroup)
	}
	d.SetId(collectionID)

	return resourceFirewallApplicationRuleCollectionRead(d, meta)
}

func resourceFirewallApplicationRuleCollectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallApplicationRuleCollectionID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.AzureFirewallName)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Azure Firewall %q (Resource Group %q) was not found - removing from state!", id.ApplicationRuleCollectionName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Azure Firewall %q (Resource Group %q): %+v", id.ApplicationRuleCollectionName, id.ResourceGroup, err)
	}

	if read.AzureFirewallPropertiesFormat == nil {
		return fmt.Errorf("retrieving Application Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", id.ApplicationRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}
	props := *read.AzureFirewallPropertiesFormat

	if props.ApplicationRuleCollections == nil {
		return fmt.Errorf("retrieving Application Rule Collection %q (Firewall %q / Resource Group %q): `props.ApplicationRuleCollections` was nil", id.ApplicationRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	var rule *network.AzureFirewallApplicationRuleCollection
	for _, r := range *props.ApplicationRuleCollections {
		if r.Name == nil {
			continue
		}

		if *r.Name == id.ApplicationRuleCollectionName {
			rule = &r
			break
		}
	}

	if rule == nil {
		log.Printf("[DEBUG] Application Rule Collection %q was not found on Firewall %q (Resource Group %q) - removing from state!", id.ApplicationRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("name", id.ApplicationRuleCollectionName)
	d.Set("azure_firewall_name", id.AzureFirewallName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := rule.AzureFirewallApplicationRuleCollectionPropertiesFormat; props != nil {
		if action := props.Action; action != nil {
			d.Set("action", string(action.Type))
		}

		if priority := props.Priority; priority != nil {
			d.Set("priority", int(*priority))
		}

		flattenedRules := flattenFirewallApplicationRuleCollectionRules(props.Rules)
		if err := d.Set("rule", flattenedRules); err != nil {
			return fmt.Errorf("setting `rule`: %+v", err)
		}
	}

	return nil
}

func resourceFirewallApplicationRuleCollectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallApplicationRuleCollectionID(d.Id())
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

		return fmt.Errorf("making Read request on Azure Firewall %s : %+v", *id, err)
	}

	props := firewall.AzureFirewallPropertiesFormat
	if props == nil {
		return fmt.Errorf("retrieving Application Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", id.ApplicationRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}
	if props.ApplicationRuleCollections == nil {
		return fmt.Errorf("retrieving Application Rule Collection %q (Firewall %q / Resource Group %q): `props.ApplicationRuleCollections` was nil", id.ApplicationRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	applicationRules := make([]network.AzureFirewallApplicationRuleCollection, 0)
	for _, rule := range *props.ApplicationRuleCollections {
		if rule.Name == nil {
			continue
		}

		if *rule.Name != id.ApplicationRuleCollectionName {
			applicationRules = append(applicationRules, rule)
		}
	}
	props.ApplicationRuleCollections = &applicationRules

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.AzureFirewallName, firewall)
	if err != nil {
		return fmt.Errorf("deleting Application Rule Collection %q from Firewall %q (Resource Group %q): %+v", id.ApplicationRuleCollectionName, id.AzureFirewallName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Application Rule Collection %q from Firewall %q (Resource Group %q): %+v", id.ApplicationRuleCollectionName, id.AzureFirewallName, id.ResourceGroup, err)
	}

	return nil
}

func expandFirewallApplicationRules(inputs []interface{}) (*[]network.AzureFirewallApplicationRule, error) {
	outputs := make([]network.AzureFirewallApplicationRule, 0)

	for _, input := range inputs {
		rule := input.(map[string]interface{})

		ruleName := rule["name"].(string)
		ruleDescription := rule["description"].(string)
		ruleSourceAddresses := rule["source_addresses"].([]interface{})
		ruleSourceIpGroups := rule["source_ip_groups"].([]interface{})
		ruleFqdnTags := rule["fqdn_tags"].([]interface{})
		ruleTargetFqdns := rule["target_fqdns"].([]interface{})

		output := network.AzureFirewallApplicationRule{
			Name:            utils.String(ruleName),
			Description:     utils.String(ruleDescription),
			SourceAddresses: utils.ExpandStringSlice(ruleSourceAddresses),
			SourceIPGroups:  utils.ExpandStringSlice(ruleSourceIpGroups),
			FqdnTags:        utils.ExpandStringSlice(ruleFqdnTags),
			TargetFqdns:     utils.ExpandStringSlice(ruleTargetFqdns),
		}

		ruleProtocols := make([]network.AzureFirewallApplicationRuleProtocol, 0)
		for _, v := range rule["protocol"].([]interface{}) {
			protocol := v.(map[string]interface{})
			port := protocol["port"].(int)
			ruleProtocol := network.AzureFirewallApplicationRuleProtocol{
				Port:         utils.Int32(int32(port)),
				ProtocolType: network.AzureFirewallApplicationRuleProtocolType(protocol["type"].(string)),
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

func flattenFirewallApplicationRuleCollectionRules(rules *[]network.AzureFirewallApplicationRule) []interface{} {
	outputs := make([]interface{}, 0)
	if rules == nil {
		return outputs
	}

	for _, rule := range *rules {
		output := make(map[string]interface{})
		if ruleName := rule.Name; ruleName != nil {
			output["name"] = *ruleName
		}
		if ruleDescription := rule.Description; ruleDescription != nil {
			output["description"] = *ruleDescription
		}
		if ruleSourceAddresses := rule.SourceAddresses; ruleSourceAddresses != nil {
			output["source_addresses"] = utils.FlattenStringSlice(ruleSourceAddresses)
		}
		if ruleSourceIpGroups := rule.SourceIPGroups; ruleSourceIpGroups != nil {
			output["source_ip_groups"] = utils.FlattenStringSlice(ruleSourceIpGroups)
		}
		if ruleFqdnTags := rule.FqdnTags; ruleFqdnTags != nil {
			output["fqdn_tags"] = utils.FlattenStringSlice(ruleFqdnTags)
		}
		if ruleTargetFqdns := rule.TargetFqdns; ruleTargetFqdns != nil {
			output["target_fqdns"] = utils.FlattenStringSlice(ruleTargetFqdns)
		}
		protocols := make([]map[string]interface{}, 0)
		if ruleProtocols := rule.Protocols; ruleProtocols != nil {
			for _, p := range *ruleProtocols {
				protocol := make(map[string]interface{})
				if port := p.Port; port != nil {
					protocol["port"] = int(*port)
				}
				protocol["type"] = string(p.ProtocolType)
				protocols = append(protocols, protocol)
			}
		}
		output["protocol"] = protocols
		outputs = append(outputs, output)
	}
	return outputs
}
