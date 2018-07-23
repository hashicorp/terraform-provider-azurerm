package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAzureFirewallNetworkRuleCollection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAzureFirewallNetworkRuleCollectionCreateUpdate,
		Read:   resourceArmAzureFirewallNetworkRuleCollectionRead,
		Update: resourceArmAzureFirewallNetworkRuleCollectionCreateUpdate,
		Delete: resourceArmAzureFirewallNetworkRuleCollectionDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"azure_firewall_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.AzureFirewallRCActionTypeAllow),
					string(network.AzureFirewallRCActionTypeDeny),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"rule": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
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
								}, true),
							},
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							Set:              schema.HashString,
						},
					},
				},
			},
		},
	}
}

func resourceArmAzureFirewallNetworkRuleCollectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	action := network.AzureFirewallRCActionType(d.Get("action").(string))
	priority := utils.Int32(int32(d.Get("priority").(int)))
	firewallName := d.Get("azure_firewall_name").(string)

	azureRMLockByName(firewallName, azureFirewallResourceName)
	defer azureRMUnlockByName(firewallName, azureFirewallResourceName)

	firewall, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		return fmt.Errorf("Error retrieving Azure Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	ipConfigs := fixArmAzureFirewallIPConfiguration(&firewall)
	firewall.AzureFirewallPropertiesFormat.IPConfigurations = &ipConfigs

	ruleCollections := *firewall.AzureFirewallPropertiesFormat.NetworkRuleCollections
	existingCollection, exists := findArmAzureFirewallNetworkRuleCollectionByName(&firewall, name)
	if exists {
		log.Printf("[INFO] updating existing AzureRM Azure Firewall network rule collection")
		existingCollection.AzureFirewallNetworkRuleCollectionPropertiesFormat.Action = &network.AzureFirewallRCAction{
			Type: action,
		}
		existingCollection.AzureFirewallNetworkRuleCollectionPropertiesFormat.Priority = priority
		rules := expandArmAzureFirewallNetworkRules(d)
		existingCollection.Rules = &rules
	} else {
		log.Printf("[INFO] adding new AzureRM Azure Firewall network rule collection")
		newFwRuleCol := expandArmAzureFirewallNetworkRuleCollection(d)
		ruleCollections = append(ruleCollections, newFwRuleCol)
	}
	firewall.AzureFirewallPropertiesFormat.NetworkRuleCollections = &ruleCollections

	future, err := client.CreateOrUpdate(ctx, resourceGroup, firewallName, firewall)
	if err != nil {
		return fmt.Errorf("Error creating/updating Azure Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation/update of Azure Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		return fmt.Errorf("Error retrieving Azure Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Azure Firewall %q (Resource Group %q) ID", firewallName, resourceGroup)
	}

	var collectionID string
	for _, collection := range *read.AzureFirewallPropertiesFormat.NetworkRuleCollections {
		if *collection.Name == name {
			collectionID = *collection.ID
		}
	}
	d.SetId(collectionID)

	return resourceArmAzureFirewallNetworkRuleCollectionRead(d, meta)
}

func resourceArmAzureFirewallNetworkRuleCollectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	fwName := d.Get("azure_firewall_name").(string)

	firewall, err := client.Get(ctx, resourceGroup, fwName)
	if err != nil {
		if utils.ResponseWasNotFound(firewall.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	collection, exists := findArmAzureFirewallNetworkRuleCollectionByName(&firewall, name)
	if !exists {
		d.SetId("")
		return nil
	}

	d.Set("name", collection.Name)
	d.Set("action", string(collection.AzureFirewallNetworkRuleCollectionPropertiesFormat.Action.Type))
	d.Set("priority", collection.AzureFirewallNetworkRuleCollectionPropertiesFormat.Priority)
	if rules := collection.AzureFirewallNetworkRuleCollectionPropertiesFormat.Rules; rules != nil {
		flattenedRules := flattenArmAzureFirewallNetworkRules(rules)
		if err := d.Set("rule", flattenedRules); err != nil {
			return fmt.Errorf("Error setting `rule`: %+v", err)
		}
	}

	return nil
}

func resourceArmAzureFirewallNetworkRuleCollectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	resourceGroup := id.ResourceGroup
	firewallName := id.Path["azureFirewalls"]

	azureRMLockByName(firewallName, azureFirewallResourceName)
	defer azureRMUnlockByName(firewallName, azureFirewallResourceName)

	firewall, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}
	_, exists := findArmAzureFirewallNetworkRuleCollectionByName(&firewall, name)
	if !exists {
		return nil
	}
	updatedCollections := removeArmAzureFirewallNetworkRuleCollectionByName(&firewall, name)
	firewall.AzureFirewallPropertiesFormat.NetworkRuleCollections = updatedCollections

	ipConfigs := fixArmAzureFirewallIPConfiguration(&firewall)
	firewall.AzureFirewallPropertiesFormat.IPConfigurations = &ipConfigs

	future, err := client.CreateOrUpdate(ctx, resourceGroup, firewallName, firewall)
	if err != nil {
		return fmt.Errorf("Error creating/updating Azure Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation/update of Azure Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	return nil
}

func expandArmAzureFirewallNetworkRuleCollection(d *schema.ResourceData) network.AzureFirewallNetworkRuleCollection {
	name := d.Get("name").(string)
	action := network.AzureFirewallRCActionType(d.Get("action").(string))
	rules := expandArmAzureFirewallNetworkRules(d)
	properties := network.AzureFirewallNetworkRuleCollectionPropertiesFormat{
		Action: &network.AzureFirewallRCAction{
			Type: action,
		},
		Priority: utils.Int32(int32(d.Get("priority").(int))),
		Rules:    &rules,
	}
	col := network.AzureFirewallNetworkRuleCollection{
		Name: &name,
		AzureFirewallNetworkRuleCollectionPropertiesFormat: &properties,
	}
	return col
}

func findArmAzureFirewallNetworkRuleCollectionByName(firewall *network.AzureFirewall, name string) (*network.AzureFirewallNetworkRuleCollection, bool) {
	for _, collection := range *firewall.AzureFirewallPropertiesFormat.NetworkRuleCollections {
		if collection.Name != nil && *collection.Name == name {
			return &collection, true
		}
	}
	return nil, false
}

func removeArmAzureFirewallNetworkRuleCollectionByName(firewall *network.AzureFirewall, name string) *[]network.AzureFirewallNetworkRuleCollection {
	collections := *firewall.AzureFirewallPropertiesFormat.NetworkRuleCollections
	for i, collection := range collections {
		if collection.Name != nil && *collection.Name == name {
			collections = append(collections[:i], collections[i+1:]...)
			continue
		}
	}
	return &collections
}

func expandArmAzureFirewallNetworkRules(d *schema.ResourceData) []network.AzureFirewallNetworkRule {
	nwRules := d.Get("rule").(*schema.Set).List()
	rules := make([]network.AzureFirewallNetworkRule, 0)

	for _, nwRule := range nwRules {
		rule := nwRule.(map[string]interface{})

		name := rule["name"].(string)
		description := rule["description"].(string)
		sourceAddresses := rule["source_addresses"].(*schema.Set)
		destinationAddresses := rule["destination_addresses"].(*schema.Set)
		destinationPorts := rule["destination_ports"].(*schema.Set)
		protocols := rule["protocols"].(*schema.Set)
		ruleToAdd := network.AzureFirewallNetworkRule{
			Name: &name,
		}
		ruleToAdd.Description = &description
		ruleToAdd.SourceAddresses = expandArmAzureFirewallSet(sourceAddresses)
		ruleToAdd.DestinationAddresses = expandArmAzureFirewallSet(destinationAddresses)
		ruleToAdd.DestinationPorts = expandArmAzureFirewallSet(destinationPorts)
		nrProtocols := make([]network.AzureFirewallNetworkRuleProtocol, 0)
		for _, v := range protocols.List() {
			s := network.AzureFirewallNetworkRuleProtocol(v.(string))
			nrProtocols = append(nrProtocols, s)
		}
		ruleToAdd.Protocols = &nrProtocols
		rules = append(rules, ruleToAdd)
	}

	return rules
}

func flattenArmAzureFirewallNetworkRules(rules *[]network.AzureFirewallNetworkRule) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if rules == nil {
		return result
	}
	for _, rule := range *rules {
		fwRule := make(map[string]interface{})
		if rule.Name != nil {
			fwRule["name"] = *rule.Name
		}
		if rule.Description != nil {
			fwRule["description"] = *rule.Description
		}
		if rule.SourceAddresses != nil {
			fwRule["source_addresses"] = sliceToSet(*rule.SourceAddresses)
		}
		if rule.DestinationAddresses != nil {
			fwRule["destination_addresses"] = sliceToSet(*rule.DestinationAddresses)
		}
		if rule.DestinationPorts != nil {
			fwRule["destination_ports"] = sliceToSet(*rule.DestinationPorts)
		}
		protocols := make([]string, 0)
		if rule.Protocols != nil {
			for _, protocol := range *rule.Protocols {
				protocols = append(protocols, string(protocol))
			}
		}
		fwRule["protocols"] = sliceToSet(protocols)
		result = append(result, fwRule)
	}
	return result
}
