package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApplicationSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApplicationSecurityGroupCreateUpdate,
		Read:   resourceArmApplicationSecurityGroupRead,
		Update: resourceArmApplicationSecurityGroupCreateUpdate,
		Delete: resourceArmApplicationSecurityGroupDelete,
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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"tags": tags.Schema(),
		},
	}
}

func resourceArmApplicationSecurityGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.ApplicationSecurityGroupsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Application Security Group %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_application_security_group", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	securityGroup := network.ApplicationSecurityGroup{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}
	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, securityGroup)
	if err != nil {
		return fmt.Errorf("Error creating Application Security Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the Application Security Group %q (Resource Group %q) to finish creating: %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Application Security Group %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmApplicationSecurityGroupRead(d, meta)
}

func resourceArmApplicationSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.ApplicationSecurityGroupsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["applicationSecurityGroups"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Application Security Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmApplicationSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.ApplicationSecurityGroupsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["applicationSecurityGroups"]

	log.Printf("[DEBUG] Deleting Application Security Group %q (resource group %q)", name, resourceGroup)

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error issuing delete request for Application Security Group %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deletion of Application Security Group %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
