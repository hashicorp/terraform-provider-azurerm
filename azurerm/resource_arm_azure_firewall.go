package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmAzureFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAzureFirewallCreate,
		Read:   resourceArmAzureFirewallRead,
		Update: resourceArmAzureFirewallCreate,
		Delete: resourceArmAzureFirewallDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"ip_configuration": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"internal_public_ip_address_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmAzureFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure Firewall creation.")
	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})
	ipConfigs := expandArmAzureFirewallIPConfigurations(d)

	parameters := network.AzureFirewall{
		Location: &location,
		Tags:     expandTags(tags),
		AzureFirewallPropertiesFormat: &network.AzureFirewallPropertiesFormat{
			IPConfigurations: &ipConfigs,
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation of Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Azure Firewall %q (resource group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAzureFirewallRead(d, meta)
}

func resourceArmAzureFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["azureFirewalls"]

	return nil
}

func resourceArmAzureFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	return nil
}

func expandArmAzureFirewallIPConfigurations(d *schema.ResourceData) []network.AzureFirewallIPConfiguration {
	configs := d.Get("ip_configuration").([]interface{})
	ipConfigs := make([]network.AzureFirewallIPConfiguration, 0, len(configs))
	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})
		name := data["name"].(string)
		subnetID := data["subnet_id"].(string)
		intPubID := data["private_ip_address"].(string)

		properties := network.AzureFirewallIPConfigurationPropertiesFormat{
			Subnet: &network.SubResource{
				ID: &subnetID,
			},
			InternalPublicIPAddress: &network.SubResource{
				ID: &intPubID,
			},
		}

		if v := data["private_ip_address"].(string); v != "" {
			properties.PrivateIPAddress = &v
		}

		ipConfig := network.AzureFirewallIPConfiguration{
			Name: &name,
			AzureFirewallIPConfigurationPropertiesFormat: &properties,
		}
		ipConfigs = append(ipConfigs, ipConfig)
	}
	return ipConfigs
}
