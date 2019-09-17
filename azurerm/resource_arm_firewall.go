package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"ip_configuration": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validateAzureFirewallSubnetName,
						},
						"internal_public_ip_address_id": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ValidateFunc:  azure.ValidateResourceID,
							Deprecated:    "This field has been deprecated. Use `public_ip_address_id` instead.",
							ConflictsWith: []string{"ip_configuration.0.public_ip_address_id"},
						},
						"public_ip_address_id": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ValidateFunc:  azure.ValidateResourceID,
							ConflictsWith: []string{"ip_configuration.0.internal_public_ip_address_id"},
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmFirewallCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.AzureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Azure Firewall creation")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Firewall %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_firewall", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	ipConfigs, subnetToLock, vnetToLock, err := expandArmFirewallIPConfigurations(d)
	if err != nil {
		return fmt.Errorf("Error Building list of Azure Firewall IP Configurations: %+v", err)
	}

	locks.ByName(name, azureFirewallResourceName)
	defer locks.UnlockByName(name, azureFirewallResourceName)

	locks.MultipleByName(subnetToLock, subnetResourceName)
	defer locks.UnlockMultipleByName(subnetToLock, subnetResourceName)

	locks.MultipleByName(vnetToLock, virtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetToLock, virtualNetworkResourceName)

	parameters := network.AzureFirewall{
		Location: &location,
		Tags:     tags.Expand(t),
		AzureFirewallPropertiesFormat: &network.AzureFirewallPropertiesFormat{
			IPConfigurations: ipConfigs,
		},
	}

	if !d.IsNewResource() {
		exists, err2 := client.Get(ctx, resourceGroup, name)
		if err2 != nil {
			if utils.ResponseWasNotFound(exists.Response) {
				return fmt.Errorf("Error retrieving existing Firewall %q (Resource Group %q): firewall not found in resource group", name, resourceGroup)
			}
			return fmt.Errorf("Error retrieving existing Firewall %q (Resource Group %q): %s", name, resourceGroup, err2)
		}
		if exists.AzureFirewallPropertiesFormat == nil {
			return fmt.Errorf("Error retrieving existing rules (Firewall %q / Resource Group %q): `props` was nil", name, resourceGroup)
		}
		props := *exists.AzureFirewallPropertiesFormat
		parameters.AzureFirewallPropertiesFormat.ApplicationRuleCollections = props.ApplicationRuleCollections
		parameters.AzureFirewallPropertiesFormat.NetworkRuleCollections = props.NetworkRuleCollections
		parameters.AzureFirewallPropertiesFormat.NatRuleCollections = props.NatRuleCollections
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
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
	client := meta.(*ArmClient).network.AzureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
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
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := read.AzureFirewallPropertiesFormat; props != nil {
		ipConfigs := flattenArmFirewallIPConfigurations(props.IPConfigurations)
		if err := d.Set("ip_configuration", ipConfigs); err != nil {
			return fmt.Errorf("Error setting `ip_configuration`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, read.Tags)
}

func resourceArmFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.AzureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
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

				parsedSubnetId, err2 := azure.ParseAzureResourceID(*config.Subnet.ID)
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

	locks.ByName(name, azureFirewallResourceName)
	defer locks.UnlockByName(name, azureFirewallResourceName)

	locks.MultipleByName(&subnetNamesToLock, subnetResourceName)
	defer locks.UnlockMultipleByName(&subnetNamesToLock, subnetResourceName)

	locks.MultipleByName(&virtualNetworkNamesToLock, virtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNamesToLock, virtualNetworkResourceName)

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
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

		pubID, exist := data["internal_public_ip_address_id"].(string)
		if !exist || pubID == "" {
			pubID, exist = data["public_ip_address_id"].(string)
		}

		if !exist || pubID == "" {
			return nil, nil, nil, fmt.Errorf("one of `ip_configuration.0.internal_public_ip_address_id` or `ip_configuration.0.public_ip_address_id` must be set")
		}

		subnetID, err := azure.ParseAzureResourceID(subnetId)
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
				PublicIPAddress: &network.SubResource{
					ID: utils.String(pubID),
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
				afIPConfig["public_ip_address_id"] = *id
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

func validateAzureFirewallSubnetName(v interface{}, k string) (warnings []string, errors []error) {
	parsed, err := azure.ParseAzureResourceID(v.(string))
	if err != nil {
		errors = append(errors, fmt.Errorf("Error parsing Azure Resource ID %q", v.(string)))
		return warnings, errors
	}
	subnetName := parsed.Path["subnets"]
	if subnetName != "AzureFirewallSubnet" {
		errors = append(errors, fmt.Errorf("The name of the Subnet for %q must be exactly 'AzureFirewallSubnet' to be used for the Azure Firewall resource", k))

	}

	return warnings, errors
}
