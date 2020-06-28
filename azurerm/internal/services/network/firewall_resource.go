package network

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateAzureFirewallName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"ip_configuration": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"subnet_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validateAzureFirewallSubnetName,
						},
						"public_ip_address_id": {
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

			"threat_intel_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(network.AzureFirewallThreatIntelModeAlert),
				ValidateFunc: validation.StringInSlice([]string{
					string(network.AzureFirewallThreatIntelModeOff),
					string(network.AzureFirewallThreatIntelModeAlert),
					string(network.AzureFirewallThreatIntelModeDeny),
				}, false),
			},

			"zones": azure.SchemaMultipleZones(),

			"tags": tags.Schema(),
		},
	}
}

func resourceArmFirewallCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.AzureFirewallsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

	if err := validateFirewallConfigurationSettings(d); err != nil {
		return fmt.Errorf("Error validating Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	ipConfigs, subnetToLock, vnetToLock, err := expandArmFirewallIPConfigurations(d)
	zones := azure.ExpandZones(d.Get("zones").([]interface{}))
	if err != nil {
		return fmt.Errorf("Error Building list of Azure Firewall IP Configurations: %+v", err)
	}

	locks.ByName(name, azureFirewallResourceName)
	defer locks.UnlockByName(name, azureFirewallResourceName)

	locks.MultipleByName(vnetToLock, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetToLock, VirtualNetworkResourceName)

	locks.MultipleByName(subnetToLock, SubnetResourceName)
	defer locks.UnlockMultipleByName(subnetToLock, SubnetResourceName)

	parameters := network.AzureFirewall{
		Location: &location,
		Tags:     tags.Expand(t),
		AzureFirewallPropertiesFormat: &network.AzureFirewallPropertiesFormat{
			IPConfigurations: ipConfigs,
			ThreatIntelMode:  network.AzureFirewallThreatIntelMode(d.Get("threat_intel_mode").(string)),
		},
		Zones: zones,
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
	client := meta.(*clients.Client).Network.AzureFirewallsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

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
		if err := d.Set("ip_configuration", flattenArmFirewallIPConfigurations(props.IPConfigurations)); err != nil {
			return fmt.Errorf("Error setting `ip_configuration`: %+v", err)
		}
		d.Set("threat_intel_mode", string(props.ThreatIntelMode))
	}

	if err := d.Set("zones", azure.FlattenZones(read.Zones)); err != nil {
		return fmt.Errorf("Error setting `zones`: %+v", err)
	}

	return tags.FlattenAndSet(d, read.Tags)
}

func resourceArmFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.AzureFirewallsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

				if !azure.SliceContainsValue(subnetNamesToLock, subnetName) {
					subnetNamesToLock = append(subnetNamesToLock, subnetName)
				}

				virtualNetworkName := parsedSubnetId.Path["virtualNetworks"]
				if !azure.SliceContainsValue(virtualNetworkNamesToLock, virtualNetworkName) {
					virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, virtualNetworkName)
				}
			}
		}
	}

	locks.ByName(name, azureFirewallResourceName)
	defer locks.UnlockByName(name, azureFirewallResourceName)

	locks.MultipleByName(&virtualNetworkNamesToLock, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNamesToLock, VirtualNetworkResourceName)

	locks.MultipleByName(&subnetNamesToLock, SubnetResourceName)
	defer locks.UnlockMultipleByName(&subnetNamesToLock, SubnetResourceName)

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
		pubID := data["public_ip_address_id"].(string)

		ipConfig := network.AzureFirewallIPConfiguration{
			Name: utils.String(name),
			AzureFirewallIPConfigurationPropertiesFormat: &network.AzureFirewallIPConfigurationPropertiesFormat{
				PublicIPAddress: &network.SubResource{
					ID: utils.String(pubID),
				},
			},
		}

		if subnetId != "" {
			subnetID, err := azure.ParseAzureResourceID(subnetId)
			if err != nil {
				return nil, nil, nil, err
			}

			subnetName := subnetID.Path["subnets"]
			virtualNetworkName := subnetID.Path["virtualNetworks"]

			if !azure.SliceContainsValue(subnetNamesToLock, subnetName) {
				subnetNamesToLock = append(subnetNamesToLock, subnetName)
			}

			if !azure.SliceContainsValue(virtualNetworkNamesToLock, virtualNetworkName) {
				virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, virtualNetworkName)
			}

			ipConfig.AzureFirewallIPConfigurationPropertiesFormat.Subnet = &network.SubResource{
				ID: utils.String(subnetId),
			}
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
				afIPConfig["public_ip_address_id"] = *id
			}
		}
		result = append(result, afIPConfig)
	}

	return result
}

func ValidateAzureFirewallName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// From the Portal:
	// The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.
	if matched := regexp.MustCompile(`^[0-9a-zA-Z]([0-9a-zA-Z._-]{0,}[0-9a-zA-Z_])?$`).Match([]byte(value)); !matched {
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

func validateFirewallConfigurationSettings(d *schema.ResourceData) error {
	configs := d.Get("ip_configuration").([]interface{})
	subnetNumber := 0

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})
		if subnet, exist := data["subnet_id"].(string); exist && subnet != "" {
			subnetNumber++
		}
	}

	if subnetNumber != 1 {
		return fmt.Errorf(`The "ip_configuration" is invalid, %d "subnet_id" have been set, one "subnet_id" should be set among all "ip_configuration" blocks`, subnetNumber)
	}

	return nil
}
