package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBastionHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBastionHostCreateUpdate,
		Read:   resourceArmBastionHostRead,
		Update: resourceArmBastionHostCreateUpdate,
		Delete: resourceArmBastionHostDelete,
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
				ValidateFunc: validateAzureRMBastionHostName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"ip_configuration": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAzureRMBastionIPConfigName,
						},
						"subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"public_ip_address_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmBastionHostCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.BastionHostsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Println("[INFO] preparing arguments for Azure Bastion Host creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Bastion Host %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_bastion_host", *existing.ID)
		}
	}

	parameters := network.BastionHost{
		Location: &location,
		BastionHostPropertiesFormat: &network.BastionHostPropertiesFormat{
			IPConfigurations: expandArmBastionHostIPConfiguration(d.Get("ip_configuration").([]interface{})),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating Bastion Host %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of Bastion Host %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Bastion Host %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceArmBastionHostRead(d, meta)
}

func resourceArmBastionHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.BastionHostsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["bastionHosts"]
	resourceGroup := id.ResourceGroup

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			log.Printf("[DEBUG] Bastion Host %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			return nil
		}
		return fmt.Errorf("Error reading the state of Bastion Host %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.BastionHostPropertiesFormat; props != nil {
		d.Set("dns_name", props.DNSName)

		if ipConfigs := props.IPConfigurations; ipConfigs != nil {
			if err := d.Set("ip_configuration", flattenArmBastionHostIPConfiguration(ipConfigs)); err != nil {
				return fmt.Errorf("Error flattening `ip_configuration`: %+v", err)
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmBastionHostDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.BastionHostsClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["bastionHosts"]
	resourceGroup := id.ResourceGroup

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Bastion Host %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deletion of Bastion Host %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func validateAzureRMBastionHostName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[A-Za-z0-9][a-zA-Z0-9_.-]+[a-zA-Z0-9_]$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens. %q: %q", k, value))
	}

	if 1 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 characters: %q", k, value))
	}

	if len(value) > 80 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 80 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}

func validateAzureRMBastionIPConfigName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[A-Za-z0-9][a-zA-Z0-9_.-]+[a-zA-Z0-9_]$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens. %q: %q", k, value))
	}

	if 1 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 characters: %q", k, value))
	}

	if len(value) > 80 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 80 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}

func expandArmBastionHostIPConfiguration(input []interface{}) (ipConfigs *[]network.BastionHostIPConfiguration) {
	if len(input) == 0 {
		return nil
	}

	property := input[0].(map[string]interface{})
	ipConfName := property["name"].(string)
	subID := property["subnet_id"].(string)
	pipID := property["public_ip_address_id"].(string)

	return &[]network.BastionHostIPConfiguration{
		{
			Name: &ipConfName,
			BastionHostIPConfigurationPropertiesFormat: &network.BastionHostIPConfigurationPropertiesFormat{
				Subnet: &network.SubResource{
					ID: &subID,
				},
				PublicIPAddress: &network.SubResource{
					ID: &pipID,
				},
			},
		},
	}
}

func flattenArmBastionHostIPConfiguration(ipConfigs *[]network.BastionHostIPConfiguration) []interface{} {
	result := make([]interface{}, 0)
	if ipConfigs == nil {
		return result
	}

	for _, config := range *ipConfigs {
		ipConfig := make(map[string]interface{})

		if config.Name != nil {
			ipConfig["name"] = *config.Name
		}

		if props := config.BastionHostIPConfigurationPropertiesFormat; props != nil {
			if subnet := props.Subnet; subnet != nil {
				ipConfig["subnet_id"] = *subnet.ID
			}

			if pip := props.PublicIPAddress; pip != nil {
				ipConfig["public_ip_address_id"] = *pip.ID
			}
		}

		result = append(result, ipConfig)
	}
	return result
}
