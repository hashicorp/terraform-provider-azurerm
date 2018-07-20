package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAzureFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAzureFirewallCreateUpdate,
		Read:   resourceArmAzureFirewallRead,
		Update: resourceArmAzureFirewallCreateUpdate,
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
							Computed: true,
						},
						"subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"internal_public_ip_address_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmAzureFirewallCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})
	ipConfigs, subnetToLock, vnetToLock, sgErr := expandArmAzureFirewallIPConfigurations(d)
	if sgErr != nil {
		return fmt.Errorf("Error Building list of Azure Firewall IP Configurations: %+v", sgErr)
	}

	azureRMLockMultipleByName(subnetToLock, subnetResourceName)
	defer azureRMUnlockMultipleByName(subnetToLock, subnetResourceName)

	azureRMLockMultipleByName(vnetToLock, virtualNetworkResourceName)
	defer azureRMUnlockMultipleByName(vnetToLock, virtualNetworkResourceName)

	parameters := network.AzureFirewall{
		Location: &location,
		Tags:     expandTags(tags),
		AzureFirewallPropertiesFormat: &network.AzureFirewallPropertiesFormat{
			IPConfigurations: &ipConfigs,
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
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

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.IPConfigurations != nil {
		ipConfigs := flattenArmAzureFirewallIPConfigurations(resp.IPConfigurations)
		d.Set("ip_configuration", ipConfigs)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmAzureFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["azureFirewalls"]

	configs := d.Get("ip_configuration").([]interface{})
	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		subnet_id := data["subnet_id"].(string)
		subnetID, err := parseAzureResourceID(subnet_id)
		if err != nil {
			return err
		}
		subnetName := subnetID.Path["subnets"]
		if !sliceContainsValue(subnetNamesToLock, subnetName) {
			subnetNamesToLock = append(subnetNamesToLock, subnetName)
		}

		virtualNetworkName := subnetID.Path["virtualNetworks"]
		if !sliceContainsValue(virtualNetworkNamesToLock, virtualNetworkName) {
			virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, virtualNetworkName)
		}
	}

	azureRMLockMultipleByName(&subnetNamesToLock, subnetResourceName)
	defer azureRMUnlockMultipleByName(&subnetNamesToLock, subnetResourceName)

	azureRMLockMultipleByName(&virtualNetworkNamesToLock, virtualNetworkResourceName)
	defer azureRMUnlockMultipleByName(&virtualNetworkNamesToLock, virtualNetworkResourceName)

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the deletion of Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return err
}

func expandArmAzureFirewallIPConfigurations(d *schema.ResourceData) ([]network.AzureFirewallIPConfiguration, *[]string, *[]string, error) {
	configs := d.Get("ip_configuration").([]interface{})
	ipConfigs := make([]network.AzureFirewallIPConfiguration, 0, len(configs))
	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})
		name := data["name"].(string)
		subnet_id := data["subnet_id"].(string)
		intPubID := data["internal_public_ip_address_id"].(string)

		properties := network.AzureFirewallIPConfigurationPropertiesFormat{
			Subnet: &network.SubResource{
				ID: &subnet_id,
			},
			InternalPublicIPAddress: &network.SubResource{
				ID: &intPubID,
			},
		}

		subnetID, err := parseAzureResourceID(subnet_id)
		if err != nil {
			return []network.AzureFirewallIPConfiguration{}, nil, nil, err
		}

		subnetName := subnetID.Path["subnets"]
		virtualNetworkName := subnetID.Path["virtualNetworks"]

		if !sliceContainsValue(subnetNamesToLock, subnetName) {
			subnetNamesToLock = append(subnetNamesToLock, subnetName)
		}

		if !sliceContainsValue(virtualNetworkNamesToLock, virtualNetworkName) {
			virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, virtualNetworkName)
		}

		ipConfig := network.AzureFirewallIPConfiguration{
			Name: &name,
			AzureFirewallIPConfigurationPropertiesFormat: &properties,
		}
		ipConfigs = append(ipConfigs, ipConfig)
	}
	return ipConfigs, &subnetNamesToLock, &virtualNetworkNamesToLock, nil
}

func flattenArmAzureFirewallIPConfigurations(ipConfigs *[]network.AzureFirewallIPConfiguration) []interface{} {
	result := make([]interface{}, 0, len(*ipConfigs))
	for _, ipConfig := range *ipConfigs {
		afIPConfig := make(map[string]interface{})
		props := ipConfig.AzureFirewallIPConfigurationPropertiesFormat

		afIPConfig["name"] = *ipConfig.Name
		afIPConfig["subnet_id"] = *props.Subnet.ID
		afIPConfig["private_ip_address"] = *props.PrivateIPAddress
		afIPConfig["internal_public_ip_address_id"] = *props.PublicIPAddress.ID
		result = append(result, afIPConfig)
	}

	return result
}
