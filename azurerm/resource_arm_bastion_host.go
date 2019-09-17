package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
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

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			// TODO: make this case sensitive once this API bug has been fixed:
			// https://github.com/Azure/azure-rest-api-specs/issues/5574
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"dns_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},

			"ip_configuration": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"public_ip_address_id": {
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

func resourceArmBastionHostCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.BastionHostsClient
	ctx := meta.(*ArmClient).StopContext

	log.Println("[INFO] preparing arguments for Azure Bastion Host creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	dnsName := d.Get("dns_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	properties := d.Get("ip_configuration").([]interface{})
	firstProperty := properties[0].(map[string]interface{})
	ipConfName := firstProperty["name"].(string)
	subID := firstProperty["subnet_id"].(string)
	pipID := firstProperty["public_ip_address_id"].(string)

	// Check if resources are to be imported
	if requireResourcesToBeImported && d.IsNewResource() {
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

	// subnet and public ip resources
	subnetID := network.SubResource{
		ID: &subID,
	}

	publicIPAddressID := network.SubResource{
		ID: &pipID,
	}

	bastionHostIPConfigurationPropertiesFormat := network.BastionHostIPConfigurationPropertiesFormat{
		Subnet:          &subnetID,
		PublicIPAddress: &publicIPAddressID,
	}

	bastionHostIPConfiguration := []network.BastionHostIPConfiguration{
		{
			Name: &ipConfName,
			BastionHostIPConfigurationPropertiesFormat: &bastionHostIPConfigurationPropertiesFormat,
		},
	}

	bastionHostProperties := network.BastionHostPropertiesFormat{
		IPConfigurations: &bastionHostIPConfiguration,
		DNSName:          &dnsName,
	}

	parameters := network.BastionHost{
		Location:                    &location,
		BastionHostPropertiesFormat: &bastionHostProperties,
		Tags:                        expandTags(tags),
	}

	// creation
	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating Bastion Host %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of Bastion Host %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// check presence
	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Bastion Host %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceArmBastionHostRead(d, meta)
}

func resourceArmBastionHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.BastionHostsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmBastionHostDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.BastionHostsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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
