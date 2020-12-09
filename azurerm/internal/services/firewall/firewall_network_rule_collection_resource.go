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

func resourceArmFirewallNetworkRuleCollection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFirewallNetworkRuleCollectionCreateUpdate,
		Read:   resourceArmFirewallNetworkRuleCollectionRead,
		Update: resourceArmFirewallNetworkRuleCollectionCreateUpdate,
		Delete: resourceArmFirewallNetworkRuleCollectionDelete,
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
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
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
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"destination_ip_groups": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"destination_ports": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"destination_fqdns": {
							Type:     schema.TypeSet,
							Optional: true,
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

func resourceArmFirewallNetworkRuleCollectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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

	if props.NetworkRuleCollections == nil {
		return fmt.Errorf("expanding Firewall %q (Resource Group %q): `properties.NetworkRuleCollections` was nil.", firewallName, resourceGroup)
	}
	ruleCollections := *props.NetworkRuleCollections

	networkRules, err := expandArmFirewallNetworkRules(d.Get("rule").(*schema.Set))
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

	return resourceArmFirewallNetworkRuleCollectionRead(d, meta)
}

func resourceArmFirewallNetworkRuleCollectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	firewallName := id.Path["azureFirewalls"]
	name := id.Path["networkRuleCollections"]

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
		return fmt.Errorf("retrieving Network Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", name, firewallName, resourceGroup)
	}
	props := *read.AzureFirewallPropertiesFormat

	if props.NetworkRuleCollections == nil {
		return fmt.Errorf("retrieving Network Rule Collection %q (Firewall %q / Resource Group %q): `props.NetworkRuleCollections` was nil", name, firewallName, resourceGroup)
	}

	var rule *network.AzureFirewallNetworkRuleCollection
	for _, r := range *props.NetworkRuleCollections {
		if r.Name == nil {
			continue
		}

		if *r.Name == name {
			rule = &r
			break
		}
	}

	if rule == nil {
		log.Printf("[DEBUG] Network Rule Collection %q was not found on Firewall %q (Resource Group %q) - removing from state!", name, firewallName, resourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("name", rule.Name)
	d.Set("azure_firewall_name", firewallName)
	d.Set("resource_group_name", resourceGroup)

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

func resourceArmFirewallNetworkRuleCollectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	firewallName := id.Path["azureFirewalls"]
	name := id.Path["networkRuleCollections"]

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
		return fmt.Errorf("retrieving Network Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", name, firewallName, resourceGroup)
	}
	if props.NetworkRuleCollections == nil {
		return fmt.Errorf("retrieving Network Rule Collection %q (Firewall %q / Resource Group %q): `props.NetworkRuleCollections` was nil", name, firewallName, resourceGroup)
	}

	networkRules := make([]network.AzureFirewallNetworkRuleCollection, 0)
	for _, rule := range *props.NetworkRuleCollections {
		if rule.Name == nil {
			continue
		}

		if *rule.Name != name {
			networkRules = append(networkRules, rule)
		}
	}
	props.NetworkRuleCollections = &networkRules

	future, err := client.CreateOrUpdate(ctx, resourceGroup, firewallName, firewall)
	if err != nil {
		return fmt.Errorf("deleting Network Rule Collection %q from Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Network Rule Collection %q from Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	return nil
}

func expandArmFirewallNetworkRules(input *schema.Set) (*[]network.AzureFirewallNetworkRule, error) {
	nwRules := input.List()
	rules := make([]network.AzureFirewallNetworkRule, 0)

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

		destinationIpGroups := make([]string, 0)
		for _, v := range rule["destination_ip_groups"].(*schema.Set).List() {
			destinationIpGroups = append(destinationIpGroups, v.(string))
		}

		destinationFqdns := make([]string, 0)
		for _, v := range rule["destination_fqdns"].(*schema.Set).List() {
			destinationFqdns = append(destinationFqdns, v.(string))
		}

		if len(destinationAddresses) == 0 && len(destinationIpGroups) == 0 && len(destinationFqdns) == 0 {
			return nil, fmt.Errorf("at least one of %q, %q and %q must be specified for each rule", "destination_addresses", "destination_ip_groups", "destination_fqdns")
		}

		destinationPorts := make([]string, 0)
		for _, v := range rule["destination_ports"].(*schema.Set).List() {
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

func flattenFirewallNetworkRuleCollectionRules(rules *[]network.AzureFirewallNetworkRule) []map[string]interface{} {
	outputs := make([]map[string]interface{}, 0)
	if rules == nil {
		return outputs
	}

	for _, rule := range *rules {
		var (
			name            string
			description     string
			sourceAddresses *schema.Set
			sourceIPGroups  *schema.Set
			destAddresses   *schema.Set
			destIPGroups    *schema.Set
			destPorts       *schema.Set
			destFqdns       *schema.Set
		)

		if rule.Name != nil {
			name = *rule.Name
		}
		if rule.Description != nil {
			description = *rule.Description
		}
		if rule.SourceAddresses != nil {
			sourceAddresses = set.FromStringSlice(*rule.SourceAddresses)
		}
		if rule.SourceIPGroups != nil {
			sourceIPGroups = set.FromStringSlice(*rule.SourceIPGroups)
		}
		if rule.DestinationAddresses != nil {
			destAddresses = set.FromStringSlice(*rule.DestinationAddresses)
		}
		if rule.DestinationIPGroups != nil {
			destIPGroups = set.FromStringSlice(*rule.DestinationIPGroups)
		}
		if rule.DestinationPorts != nil {
			destPorts = set.FromStringSlice(*rule.DestinationPorts)
		}
		if rule.DestinationFqdns != nil {
			destFqdns = set.FromStringSlice(*rule.DestinationFqdns)
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
			"protocols":             set.FromStringSlice(protocols),
		})
	}
	return outputs
}
