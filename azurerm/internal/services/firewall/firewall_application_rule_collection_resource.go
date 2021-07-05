package firewall

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/firewall/parse"
	firewallValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/firewall/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceFirewallApplicationRuleCollection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFirewallApplicationRuleCollectionCreateUpdate,
		Read:   resourceFirewallApplicationRuleCollectionRead,
		Update: resourceFirewallApplicationRuleCollectionCreateUpdate,
		Delete: resourceFirewallApplicationRuleCollectionDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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

			"resource_group_name": azure.SchemaResourceGroupName(),

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
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},
						"source_ip_groups": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},
						"description": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"fqdn_tags": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},
						"target_fqdns": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
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
										Optional:     true,
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

	locks.ByName(firewallName, azureFirewallResourceName)
	defer locks.UnlockByName(firewallName, azureFirewallResourceName)

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

	resourceGroup := id.ResourceGroup
	firewallName := id.AzureFirewallName
	name := id.ApplicationRuleCollectionName

	read, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Azure Firewall %q (Resource Group %q) was not found - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.AzureFirewallPropertiesFormat == nil {
		return fmt.Errorf("retrieving Application Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", name, firewallName, resourceGroup)
	}
	props := *read.AzureFirewallPropertiesFormat

	if props.ApplicationRuleCollections == nil {
		return fmt.Errorf("retrieving Application Rule Collection %q (Firewall %q / Resource Group %q): `props.ApplicationRuleCollections` was nil", name, firewallName, resourceGroup)
	}

	var rule *network.AzureFirewallApplicationRuleCollection
	for _, r := range *props.ApplicationRuleCollections {
		if r.Name == nil {
			continue
		}

		if *r.Name == name {
			rule = &r
			break
		}
	}

	if rule == nil {
		log.Printf("[DEBUG] Application Rule Collection %q was not found on Firewall %q (Resource Group %q) - removing from state!", name, firewallName, resourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("name", rule.Name)
	d.Set("azure_firewall_name", firewallName)
	d.Set("resource_group_name", resourceGroup)

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

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	firewallName := id.Path["azureFirewalls"]
	name := id.Path["applicationRuleCollections"]

	locks.ByName(firewallName, azureFirewallResourceName)
	defer locks.UnlockByName(firewallName, azureFirewallResourceName)

	firewall, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		if utils.ResponseWasNotFound(firewall.Response) {
			// assume deleted
			return nil
		}

		return fmt.Errorf("making Read request on Azure Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	props := firewall.AzureFirewallPropertiesFormat
	if props == nil {
		return fmt.Errorf("retrieving Application Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", name, firewallName, resourceGroup)
	}
	if props.ApplicationRuleCollections == nil {
		return fmt.Errorf("retrieving Application Rule Collection %q (Firewall %q / Resource Group %q): `props.ApplicationRuleCollections` was nil", name, firewallName, resourceGroup)
	}

	applicationRules := make([]network.AzureFirewallApplicationRuleCollection, 0)
	for _, rule := range *props.ApplicationRuleCollections {
		if rule.Name == nil {
			continue
		}

		if *rule.Name != name {
			applicationRules = append(applicationRules, rule)
		}
	}
	props.ApplicationRuleCollections = &applicationRules

	future, err := client.CreateOrUpdate(ctx, resourceGroup, firewallName, firewall)
	if err != nil {
		return fmt.Errorf("deleting Application Rule Collection %q from Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Application Rule Collection %q from Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	return nil
}

func expandFirewallApplicationRules(inputs []interface{}) (*[]network.AzureFirewallApplicationRule, error) {
	outputs := make([]network.AzureFirewallApplicationRule, 0)

	for _, input := range inputs {
		rule := input.(map[string]interface{})

		ruleName := rule["name"].(string)
		ruleDescription := rule["description"].(string)
		ruleSourceAddresses := rule["source_addresses"].(*pluginsdk.Set).List()
		ruleSourceIpGroups := rule["source_ip_groups"].(*pluginsdk.Set).List()
		ruleFqdnTags := rule["fqdn_tags"].(*pluginsdk.Set).List()
		ruleTargetFqdns := rule["target_fqdns"].(*pluginsdk.Set).List()

		output := network.AzureFirewallApplicationRule{
			Name:            utils.String(ruleName),
			Description:     utils.String(ruleDescription),
			SourceAddresses: utils.ExpandStringSlice(ruleSourceAddresses),
			SourceIPGroups:  utils.ExpandStringSlice(ruleSourceIpGroups),
			FqdnTags:        utils.ExpandStringSlice(ruleFqdnTags),
			TargetFqdns:     utils.ExpandStringSlice(ruleTargetFqdns),
		}

		ruleProtocols := make([]network.AzureFirewallApplicationRuleProtocol, 0)
		protocols := rule["protocol"].([]interface{})
		for _, v := range protocols {
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

func flattenFirewallApplicationRuleCollectionRules(rules *[]network.AzureFirewallApplicationRule) []map[string]interface{} {
	outputs := make([]map[string]interface{}, 0)
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
			output["source_addresses"] = set.FromStringSlice(*ruleSourceAddresses)
		}
		if ruleSourceIpGroups := rule.SourceIPGroups; ruleSourceIpGroups != nil {
			output["source_ip_groups"] = set.FromStringSlice(*ruleSourceIpGroups)
		}
		if ruleFqdnTags := rule.FqdnTags; ruleFqdnTags != nil {
			output["fqdn_tags"] = set.FromStringSlice(*ruleFqdnTags)
		}
		if ruleTargetFqdns := rule.TargetFqdns; ruleTargetFqdns != nil {
			output["target_fqdns"] = set.FromStringSlice(*ruleTargetFqdns)
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
