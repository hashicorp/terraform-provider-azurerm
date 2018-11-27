package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var azureFirewallResourceName = "azurerm_firewall"

func resourceArmFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFirewallCreateUpdate,
		Read:   resourceArmFirewallRead,
		Update: resourceArmFirewallCreateUpdate,
		Delete: resourceArmFirewallDelete,
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

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"ip_configuration": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
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
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmFirewallCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Azure Firewall creation")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})
	ipConfigs, subnetToLock, vnetToLock, err := expandArmFirewallIPConfigurations(d)
	if err != nil {
		return fmt.Errorf("Error Building list of Azure Firewall IP Configurations: %+v", err)
	}

	azureRMLockByName(name, azureFirewallResourceName)
	defer azureRMUnlockByName(name, azureFirewallResourceName)

	azureRMLockMultipleByName(subnetToLock, subnetResourceName)
	defer azureRMUnlockMultipleByName(subnetToLock, subnetResourceName)

	azureRMLockMultipleByName(vnetToLock, virtualNetworkResourceName)
	defer azureRMUnlockMultipleByName(vnetToLock, virtualNetworkResourceName)

	parameters := network.AzureFirewall{
		Location: &location,
		Tags:     expandTags(tags),
		AzureFirewallPropertiesFormat: &network.AzureFirewallPropertiesFormat{
			IPConfigurations: ipConfigs,
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
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

	return resourceArmFirewallRead(d, meta)
}

func resourceArmFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["azureFirewalls"]

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Firewall %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", read.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := read.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := read.AzureFirewallPropertiesFormat; props != nil {
		ipConfigs := flattenArmFirewallIPConfigurations(props.IPConfigurations)
		if err := d.Set("ip_configuration", ipConfigs); err != nil {
			return fmt.Errorf("Error setting `ip_configuration`: %+v", err)
		}
	}

	flattenAndSetTags(d, read.Tags)

	return nil
}

func resourceArmFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["azureFirewalls"]

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			// deleted outside of TF
			log.Printf("[DEBUG] Firewall %q was not found in Resource Group %q - assuming removed!", name, resourceGroup)
			return nil
		}

		return fmt.Errorf("Error retrieving Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)
	if props := read.AzureFirewallPropertiesFormat; props != nil {
		if configs := props.IPConfigurations; configs != nil {
			for _, config := range *configs {
				if config.Subnet == nil || config.Subnet.ID == nil {
					continue
				}

				parsedSubnetId, err2 := parseAzureResourceID(*config.Subnet.ID)
				if err2 != nil {
					return err2
				}
				subnetName := parsedSubnetId.Path["subnets"]

				if !sliceContainsValue(subnetNamesToLock, subnetName) {
					subnetNamesToLock = append(subnetNamesToLock, subnetName)
				}

				virtualNetworkName := parsedSubnetId.Path["virtualNetworks"]
				if !sliceContainsValue(virtualNetworkNamesToLock, virtualNetworkName) {
					virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, virtualNetworkName)
				}
			}
		}
	}

	azureRMLockByName(name, azureFirewallResourceName)
	defer azureRMUnlockByName(name, azureFirewallResourceName)

	azureRMLockMultipleByName(&subnetNamesToLock, subnetResourceName)
	defer azureRMUnlockMultipleByName(&subnetNamesToLock, subnetResourceName)

	azureRMLockMultipleByName(&virtualNetworkNamesToLock, virtualNetworkResourceName)
	defer azureRMUnlockMultipleByName(&virtualNetworkNamesToLock, virtualNetworkResourceName)

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the deletion of Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return err
}

func expandArmFirewallIPConfigurations(d *schema.ResourceData) (*[]network.AzureFirewallIPConfiguration, *[]string, *[]string, error) {
	configs := d.Get("ip_configuration").([]interface{})
	ipConfigs := make([]network.AzureFirewallIPConfiguration, 0)
	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})
		name := data["name"].(string)
		subnetId := data["subnet_id"].(string)
		intPubID := data["internal_public_ip_address_id"].(string)

		subnetID, err := parseAzureResourceID(subnetId)
		if err != nil {
			return nil, nil, nil, err
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
			Name: utils.String(name),
			AzureFirewallIPConfigurationPropertiesFormat: &network.AzureFirewallIPConfigurationPropertiesFormat{
				Subnet: &network.SubResource{
					ID: utils.String(subnetId),
				},
				InternalPublicIPAddress: &network.SubResource{
					ID: utils.String(intPubID),
				},
			},
		}
		ipConfigs = append(ipConfigs, ipConfig)
	}
	return &ipConfigs, &subnetNamesToLock, &virtualNetworkNamesToLock, nil
}

func flattenArmFirewallIPConfigurations(input *[]network.AzureFirewallIPConfiguration) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		afIPConfig := make(map[string]interface{})
		props := v.AzureFirewallIPConfigurationPropertiesFormat
		if props == nil {
			continue
		}

		if name := v.Name; name != nil {
			afIPConfig["name"] = *name
		}

		if subnet := props.Subnet; subnet != nil {
			if id := subnet.ID; id != nil {
				afIPConfig["subnet_id"] = *id
			}
		}

		if ipAddress := props.PrivateIPAddress; ipAddress != nil {
			afIPConfig["private_ip_address"] = *ipAddress
		}

		if pip := props.PublicIPAddress; pip != nil {
			if id := pip.ID; id != nil {
				afIPConfig["internal_public_ip_address_id"] = *id
			}
		}
		result = append(result, afIPConfig)
	}

	return result
}

func validateAzureFirewallName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// From the Portal:
	// The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.
	if matched := regexp.MustCompile(`^[0-9a-zA-Z]([0-9a-zA-Z.\_-]{0,}[0-9a-zA-Z_])?$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.", k))
	}

	return warnings, errors
}
