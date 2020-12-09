package firewall

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmFirewallNatRuleCollection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFirewallNatRuleCollectionCreateUpdate,
		Read:   resourceArmFirewallNatRuleCollectionRead,
		Update: resourceArmFirewallNatRuleCollectionCreateUpdate,
		Delete: resourceArmFirewallNatRuleCollectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateAzureFirewallName,
			},

			"azure_firewall_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateAzureFirewallName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(100, 65000),
			},

			"action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.Dnat),
					string(network.Snat),
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
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"translated_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"translated_port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source_addresses": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"source_ip_groups": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"destination_addresses": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"destination_ports": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"protocols": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(network.Any),
									string(network.ICMP),
									string(network.TCP),
									string(network.UDP),
								}, false),
							},
							Set: schema.HashString,
						},
					},
				},
			},
		},
	}
}

func resourceArmFirewallNatRuleCollectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	firewallName := d.Get("azure_firewall_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(firewallName, azureFirewallResourceName)
	defer locks.UnlockByName(firewallName, azureFirewallResourceName)

	firewall, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		return fmt.Errorf("retrieving Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	if firewall.AzureFirewallPropertiesFormat == nil {
		return fmt.Errorf("expanding Firewall %q (Resource Group %q): `properties` was nil.", firewallName, resourceGroup)
	}
	props := *firewall.AzureFirewallPropertiesFormat

	if props.NatRuleCollections == nil {
		return fmt.Errorf("expanding Firewall %q (Resource Group %q): `properties.NatRuleCollections` was nil.", firewallName, resourceGroup)
	}

	ruleCollections := *props.NatRuleCollections
	natRules, err := expandArmFirewallNatRules(d.Get("rule").(*schema.Set))
	if err != nil {
		return fmt.Errorf("expanding Firewall NAT Rules: %+v", err)
	}
	priority := d.Get("priority").(int)
	newRuleCollection := network.AzureFirewallNatRuleCollection{
		Name: utils.String(name),
		AzureFirewallNatRuleCollectionProperties: &network.AzureFirewallNatRuleCollectionProperties{
			Action: &network.AzureFirewallNatRCAction{
				Type: network.AzureFirewallNatRCActionType(d.Get("action").(string)),
			},
			Priority: utils.Int32(int32(priority)),
			Rules:    natRules,
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

	firewall.AzureFirewallPropertiesFormat.NatRuleCollections = &ruleCollections
	future, err := client.CreateOrUpdate(ctx, resourceGroup, firewallName, firewall)
	if err != nil {
		return fmt.Errorf("creating/updating NAT Rule Collection %q in Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of NAT Rule Collection %q of Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		return fmt.Errorf("retrieving Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	var collectionID string
	if props := read.AzureFirewallPropertiesFormat; props != nil {
		if collections := props.NatRuleCollections; collections != nil {
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
		return fmt.Errorf("Cannot find ID for NAT Rule Collection %q (Azure Firewall %q / Resource Group %q)", name, firewallName, resourceGroup)
	}
	d.SetId(collectionID)

	return resourceArmFirewallNatRuleCollectionRead(d, meta)
}

func resourceArmFirewallNatRuleCollectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	firewallName := id.Path["azureFirewalls"]
	name := id.Path["natRuleCollections"]

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
		return fmt.Errorf("retrieving NAT Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", name, firewallName, resourceGroup)
	}

	props := *read.AzureFirewallPropertiesFormat

	if props.NatRuleCollections == nil {
		return fmt.Errorf("retrieving NAT Rule Collection %q (Firewall %q / Resource Group %q): `props.NetworkRuleCollections` was nil", name, firewallName, resourceGroup)
	}

	var rule *network.AzureFirewallNatRuleCollection
	for _, r := range *props.NatRuleCollections {
		if r.Name == nil {
			continue
		}

		if *r.Name == name {
			rule = &r
			break
		}
	}

	if rule == nil {
		log.Printf("[DEBUG] NAT Rule Collection %q was not found on Firewall %q (Resource Group %q) - removing from state!", name, firewallName, resourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("name", rule.Name)
	d.Set("azure_firewall_name", firewallName)
	d.Set("resource_group_name", resourceGroup)

	if props := rule.AzureFirewallNatRuleCollectionProperties; props != nil {
		if action := props.Action; action != nil {
			d.Set("action", string(action.Type))
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

func resourceArmFirewallNatRuleCollectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	firewallName := id.Path["azureFirewalls"]
	name := id.Path["natRuleCollections"]

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
		return fmt.Errorf("retrieving NAT Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", name, firewallName, resourceGroup)
	}
	if props.NetworkRuleCollections == nil {
		return fmt.Errorf("retrieving NAT Rule Collection %q (Firewall %q / Resource Group %q): `props.NatRuleCollections` was nil", name, firewallName, resourceGroup)
	}

	natRules := make([]network.AzureFirewallNatRuleCollection, 0)
	for _, rule := range *props.NatRuleCollections {
		if rule.Name == nil {
			continue
		}

		if *rule.Name != name {
			natRules = append(natRules, rule)
		}
	}
	props.NatRuleCollections = &natRules

	future, err := client.CreateOrUpdate(ctx, resourceGroup, firewallName, firewall)
	if err != nil {
		return fmt.Errorf("deleting NAT Rule Collection %q from Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of NAT Rule Collection %q from Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	return nil
}

func expandArmFirewallNatRules(input *schema.Set) (*[]network.AzureFirewallNatRule, error) {
	nwRules := input.List()
	rules := make([]network.AzureFirewallNatRule, 0)

	for _, nwRule := range nwRules {
		rule := nwRule.(map[string]interface{})

		name := rule["name"].(string)
		description := rule["description"].(string)

		sourceAddresses := make([]string, 0)
		for _, v := range rule["source_addresses"].(*schema.Set).List() {
			sourceAddresses = append(sourceAddresses, v.(string))
		}

		sourceIpGroups := make([]string, 0)
		for _, v := range rule["source_ip_groups"].(*schema.Set).List() {
			sourceIpGroups = append(sourceIpGroups, v.(string))
		}

		if len(sourceAddresses) == 0 && len(sourceIpGroups) == 0 {
			return nil, fmt.Errorf("at least one of %q and %q must be specified for each rule", "source_addresses", "source_ip_groups")
		}

		destinationAddresses := make([]string, 0)
		for _, v := range rule["destination_addresses"].(*schema.Set).List() {
			destinationAddresses = append(destinationAddresses, v.(string))
		}

		destinationPorts := make([]string, 0)
		for _, v := range rule["destination_ports"].(*schema.Set).List() {
			destinationPorts = append(destinationPorts, v.(string))
		}

		translatedAddress := rule["translated_address"].(string)
		translatedPort := rule["translated_port"].(string)

		ruleToAdd := network.AzureFirewallNatRule{
			Name:                 utils.String(name),
			Description:          utils.String(description),
			SourceAddresses:      &sourceAddresses,
			SourceIPGroups:       &sourceIpGroups,
			DestinationAddresses: &destinationAddresses,
			DestinationPorts:     &destinationPorts,
			TranslatedAddress:    &translatedAddress,
			TranslatedPort:       &translatedPort,
		}

		nrProtocols := make([]network.AzureFirewallNetworkRuleProtocol, 0)
		protocols := rule["protocols"].(*schema.Set)
		for _, v := range protocols.List() {
			s := network.AzureFirewallNetworkRuleProtocol(v.(string))
			nrProtocols = append(nrProtocols, s)
		}
		ruleToAdd.Protocols = &nrProtocols
		rules = append(rules, ruleToAdd)
	}

	return &rules, nil
}

func flattenFirewallNatRuleCollectionRules(rules *[]network.AzureFirewallNatRule) []map[string]interface{} {
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
		if rule.TranslatedAddress != nil {
			output["translated_address"] = *rule.TranslatedAddress
		}
		if rule.TranslatedPort != nil {
			output["translated_port"] = *rule.TranslatedPort
		}
		if rule.SourceAddresses != nil {
			output["source_addresses"] = set.FromStringSlice(*rule.SourceAddresses)
		}
		if rule.SourceIPGroups != nil {
			output["source_ip_groups"] = set.FromStringSlice(*rule.SourceIPGroups)
		}
		if rule.DestinationAddresses != nil {
			output["destination_addresses"] = set.FromStringSlice(*rule.DestinationAddresses)
		}
		if rule.DestinationPorts != nil {
			output["destination_ports"] = set.FromStringSlice(*rule.DestinationPorts)
		}
		protocols := make([]string, 0)
		if rule.Protocols != nil {
			for _, protocol := range *rule.Protocols {
				protocols = append(protocols, string(protocol))
			}
		}
		output["protocols"] = set.FromStringSlice(protocols)
		outputs = append(outputs, output)
	}
	return outputs
}
