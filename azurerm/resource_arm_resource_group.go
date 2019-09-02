package azurerm

import (
	"fmt"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmResourceGroupCreateUpdate,
		Read:   resourceArmResourceGroupRead,
		Update: resourceArmResourceGroupCreateUpdate,
		Exists: resourceArmResourceGroupExists,
		Delete: resourceArmResourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"tags": tags.Schema(),
		},
	}
}

func resourceArmResourceGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).resource.GroupsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing resource group: %+v", err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_resource_group", *existing.ID)
		}
	}

	parameters := resources.Group{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}

	if _, err := client.CreateOrUpdate(ctx, name, parameters); err != nil {
		return fmt.Errorf("Error creating resource group: %+v", err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("Error retrieving resource group: %+v", err)
	}

	d.SetId(*resp.ID)

	return resourceArmResourceGroupRead(d, meta)
}

func resourceArmResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).resource.GroupsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q: %+v", d.Id(), err)
	}

	name := id.ResourceGroup

	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading resource group %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading resource group: %+v", err)
	}

	d.Set("name", resp.Name)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmResourceGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*ArmClient).resource.GroupsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return false, fmt.Errorf("Error parsing Azure Resource ID %q: %+v", d.Id(), err)
	}

	name := id.ResourceGroup

	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return false, nil
		}

		return false, fmt.Errorf("Error reading resource group: %+v", err)
	}

	return true, nil
}

func resourceArmResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).resource.GroupsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q: %+v", d.Id(), err)
	}

	name := id.ResourceGroup

	deleteFuture, err := client.Delete(ctx, name)
	if err != nil {
		if response.WasNotFound(deleteFuture.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting Resource Group %q: %+v", name, err)
	}

	err = deleteFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(deleteFuture.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting Resource Group %q: %+v", name, err)
	}

	return nil
}
