package azurerm

import (
	"fmt"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-08-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmFirewallApplicationRuleCollection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFirewallApplicationRuleCollectionCreateUpdate,
		Read:   resourceArmFirewallApplicationRuleCollectionRead,
		Update: resourceArmFirewallApplicationRuleCollectionCreateUpdate,
		Delete: resourceArmFirewallApplicationRuleCollectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureFirewallName,
			},

			"azure_firewall_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureFirewallName,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(100, 65000),
			},

			"action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.AzureFirewallRCActionTypeAllow),
					string(network.AzureFirewallRCActionTypeDeny),
				}, false),
			},

			"rule": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source_addresses": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"fqdn_tags": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
							// ConflictsWith: []string{
							// 	"target_fqdns",
							// 	"protocol",
							// },
						},
						"target_fqdns": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
							// ConflictsWith: []string{"fqdn_tags"},
						},
						"protocol": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validate.PortNumber,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(network.AzureFirewallApplicationRuleProtocolTypeHTTP),
											string(network.AzureFirewallApplicationRuleProtocolTypeHTTPS),
										}, false),
									},
								},
							},
							// Set: resourceArmFirewallApplicationRuleCollectionRuleProtocolHash,
							// ConflictsWith: []string{"fqdn_tags"},
						},
					},
				},
			},
		},
	}
}

func resourceArmFirewallApplicationRuleCollectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	firewallName := d.Get("azure_firewall_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	azureRMLockByName(firewallName, azureFirewallResourceName)
	defer azureRMUnlockByName(firewallName, azureFirewallResourceName)

	firewall, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		return fmt.Errorf("Error retrieving Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	if firewall.AzureFirewallPropertiesFormat == nil {
		return fmt.Errorf("Error expanding Firewall %q (Resource Group %q): `properties` was nil.", firewallName, resourceGroup)
	}
	props := *firewall.AzureFirewallPropertiesFormat

	if props.ApplicationRuleCollections == nil {
		return fmt.Errorf("Error expanding Firewall %q (Resource Group %q): `properties.ApplicationRuleCollections` was nil.", firewallName, resourceGroup)
	}
	ruleCollections := *props.ApplicationRuleCollections

	applicationRules := expandArmFirewallApplicationRules(d.Get("rule").(*schema.Set))
	priority := d.Get("priority").(int)
	newRuleCollection := network.AzureFirewallApplicationRuleCollection{
		Name: utils.String(name),
		AzureFirewallApplicationRuleCollectionPropertiesFormat: &network.AzureFirewallApplicationRuleCollectionPropertiesFormat{
			Action: &network.AzureFirewallRCAction{
				Type: network.AzureFirewallRCActionType(d.Get("action").(string)),
			},
			Priority: utils.Int32(int32(priority)),
			Rules:    &applicationRules,
		},
	}

	if !d.IsNewResource() {
		index := -1
		for i, v := range ruleCollections {
			if v.Name == nil {
				continue
			}

			if *v.Name == name {
				index = i
				break
			}
		}

		if index == -1 {
			return fmt.Errorf("Error locating Application Rule Collection %q (Firewall %q / Resource Group %q)", name, firewallName, resourceGroup)
		}

		ruleCollections[index] = newRuleCollection
	} else {
		ruleCollections = append(ruleCollections, newRuleCollection)
	}

	firewall.AzureFirewallPropertiesFormat.ApplicationRuleCollections = &ruleCollections

	future, err := client.CreateOrUpdate(ctx, resourceGroup, firewallName, firewall)
	if err != nil {
		return fmt.Errorf("Error creating/updating Application Rule Collection %q in Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation/update of Application Rule Collection %q of Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		return fmt.Errorf("Error retrieving Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
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

	return resourceArmFirewallApplicationRuleCollectionRead(d, meta)
}

func resourceArmFirewallApplicationRuleCollectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	firewallName := id.Path["azureFirewalls"]
	name := id.Path["applicationRuleCollections"]

	read, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Azure Firewall %q (Resource Group %q) was not found - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.AzureFirewallPropertiesFormat == nil {
		return fmt.Errorf("Error retrieving Application Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", name, firewallName, resourceGroup)
	}
	props := *read.AzureFirewallPropertiesFormat

	if props.ApplicationRuleCollections == nil {
		return fmt.Errorf("Error retrieving Application Rule Collection %q (Firewall %q / Resource Group %q): `props.ApplicationRuleCollections` was nil", name, firewallName, resourceGroup)
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
		// if err := d.Set("rule", schema.NewSet(resourceArmFirewallApplicationRuleCollectionRuleProtocolHash, flattenedRules)); err != nil {
		if err := d.Set("rule", flattenedRules); err != nil {
			return fmt.Errorf("Error setting `rule`: %+v", err)
		}
	}

	return nil
}

func resourceArmFirewallApplicationRuleCollectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	firewallName := id.Path["azureFirewalls"]
	name := id.Path["applicationRuleCollections"]

	azureRMLockByName(firewallName, azureFirewallResourceName)
	defer azureRMUnlockByName(firewallName, azureFirewallResourceName)

	firewall, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		if utils.ResponseWasNotFound(firewall.Response) {
			// assume deleted
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	props := firewall.AzureFirewallPropertiesFormat
	if props == nil {
		return fmt.Errorf("Error retrieving Application Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", name, firewallName, resourceGroup)
	}
	if props.ApplicationRuleCollections == nil {
		return fmt.Errorf("Error retrieving Application Rule Collection %q (Firewall %q / Resource Group %q): `props.ApplicationRuleCollections` was nil", name, firewallName, resourceGroup)
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
		return fmt.Errorf("Error deleting Application Rule Collection %q from Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for deletion of Application Rule Collection %q from Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	return nil
}

func expandArmFirewallApplicationRules(input *schema.Set) []network.AzureFirewallApplicationRule {
	appRules := input.List()
	rules := make([]network.AzureFirewallApplicationRule, 0)

	for _, appRule := range appRules {
		rule := appRule.(map[string]interface{})

		name := rule["name"].(string)
		description := rule["description"].(string)

		sourceAddresses := make([]string, 0)
		for _, v := range rule["source_addresses"].(*schema.Set).List() {
			sourceAddresses = append(sourceAddresses, v.(string))
		}

		fqdnTags := make([]string, 0)
		for _, v := range rule["fqdn_tags"].(*schema.Set).List() {
			fqdnTags = append(fqdnTags, v.(string))
		}

		targetFqdns := make([]string, 0)
		for _, v := range rule["target_fqdns"].(*schema.Set).List() {
			targetFqdns = append(targetFqdns, v.(string))
		}

		ruleToAdd := network.AzureFirewallApplicationRule{
			Name:            utils.String(name),
			Description:     utils.String(description),
			SourceAddresses: &sourceAddresses,
			FqdnTags:        &fqdnTags,
			TargetFqdns:     &targetFqdns,
		}

		arProtocols := make([]network.AzureFirewallApplicationRuleProtocol, 0)
		protocols := rule["protocol"].([]interface{})
		for _, v := range protocols {
			protocol := v.(map[string]interface{})
			port := protocol["port"].(int)
			p := network.AzureFirewallApplicationRuleProtocol{
				Port:         utils.Int32(int32(port)),
				ProtocolType: network.AzureFirewallApplicationRuleProtocolType(protocol["type"].(string)),
			}
			arProtocols = append(arProtocols, p)
		}
		ruleToAdd.Protocols = &arProtocols
		rules = append(rules, ruleToAdd)
	}

	return rules
}

func flattenFirewallApplicationRuleCollectionRules(rules *[]network.AzureFirewallApplicationRule) []map[string]interface{} {
	outputs := make([]map[string]interface{}, 0)
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
		if rule.SourceAddresses != nil {
			output["source_addresses"] = sliceToSet(*rule.SourceAddresses)
		}
		if rule.FqdnTags != nil {
			output["fqdn_tags"] = sliceToSet(*rule.FqdnTags)
		}
		if rule.TargetFqdns != nil {
			output["target_fqdns"] = sliceToSet(*rule.TargetFqdns)
		}
		protocols := make([]map[string]interface{}, 0)
		if arProtocols := rule.Protocols; arProtocols != nil {
			for _, p := range *arProtocols {
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
