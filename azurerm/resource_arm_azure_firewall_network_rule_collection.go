package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
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
					string(network.AzureFirewallRCActionTypeAllow),
				}, true),
				StateFunc:        ignoreCaseStateFunc,
				DiffSuppressFunc: suppress.CaseDifference,
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
	fwName := d.Get("azure_firewall_name").(string)

	firewall, err := client.Get(ctx, resourceGroup, fwName)
	if err != nil {
		return fmt.Errorf("Error retrieving Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	ruleCollections := *firewall.AzureFirewallPropertiesFormat.NetworkRuleCollections
	existingCollection, exists := findArmAzureFirewallNetworkRuleCollectionByName(&firewall, name)
	if exists {
		existingCollection.AzureFirewallNetworkRuleCollectionPropertiesFormat.Action = &network.AzureFirewallRCAction{
			Type: action,
		}
		existingCollection.AzureFirewallNetworkRuleCollectionPropertiesFormat.Priority = priority
	} else {
		newFwRuleCol := expandArmAzureFirewallNetworkRuleCollection(d)
		ruleCollections = append(ruleCollections, newFwRuleCol)
	}
	firewall.NetworkRuleCollections = &ruleCollections
	future, err := client.CreateOrUpdate(ctx, resourceGroup, fwName, firewall)
	if err != nil {
		return fmt.Errorf("Error creating/updating Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation/update of Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Azure Firewall %q (Resource Group %q) ID", name, resourceGroup)
	}

	var collectionID string
	for _, collection := range *read.NetworkRuleCollections {
		if *collection.Name == name {
			collectionID = *collection.ID
		}
	}
	d.SetId(collectionID)

	return resourceArmAzureFirewallNetworkRuleCollectionRead(d, meta)
}

func expandArmAzureFirewallNetworkRuleCollection(d *schema.ResourceData) network.AzureFirewallNetworkRuleCollection {
	name := d.Get("name").(string)
	action := network.AzureFirewallRCActionType(d.Get("action").(string))
	properties := network.AzureFirewallNetworkRuleCollectionPropertiesFormat{
		Action: &network.AzureFirewallRCAction{
			Type: action,
		},
		Priority: utils.Int32(int32(d.Get("priority").(int))),
	}
	col := network.AzureFirewallNetworkRuleCollection{
		Name: &name,
		AzureFirewallNetworkRuleCollectionPropertiesFormat: &properties,
	}
	return col
}

func findArmAzureFirewallNetworkRuleCollectionByName(fw *network.AzureFirewall, name string) (*network.AzureFirewallNetworkRuleCollection, bool) {
	for _, collection := range *fw.NetworkRuleCollections {
		if collection.Name != nil && *collection.Name == name {
			return &collection, true
		}
	}
	return nil, false
}

func resourceArmAzureFirewallNetworkRuleCollectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	fwName := d.Get("azure_firewall_name").(string)

	firewall, err := client.Get(ctx, resourceGroup, fwName)
	if err != nil {
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
	return nil
}

func resourceArmAzureFirewallNetworkRuleCollectionDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}
