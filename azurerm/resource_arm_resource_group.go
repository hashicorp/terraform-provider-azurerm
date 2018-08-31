package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 30),
			Update: schema.DefaultTimeout(time.Minute * 30),
			// on first glance this is excessive, however since a RG can contain things unmanaged by TF
			// which causes this delete to wait for those things
			Delete: schema.DefaultTimeout(time.Minute * 60),
		},

		Schema: map[string]*schema.Schema{
			"name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"tags": tagsSchema(),
		},
	}
}

func resourceArmResourceGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).resourceGroupsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)

	if d.IsNewResource() {
		// first check if there's one in this subscription requiring import
		resp, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for the existence of Resource Group %q: %+v", name, err)
			}
		}

		if resp.ID != nil {
			return tf.ImportAsExistsError("azurerm_resource_group", *resp.ID)
		}
	}

	location := azureRMNormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})
	parameters := resources.Group{
		Location: utils.String(location),
		Tags:     expandTags(tags),
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(tf.TimeoutForCreateUpdate(d)))
	defer cancel()
	_, err := client.CreateOrUpdate(waitCtx, name, parameters)
	if err != nil {
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
	client := meta.(*ArmClient).resourceGroupsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmResourceGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*ArmClient).resourceGroupsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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
	client := meta.(*ArmClient).resourceGroupsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err = deleteFuture.WaitForCompletionRef(waitCtx, client.Client)
	if err != nil {
		if response.WasNotFound(deleteFuture.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting Resource Group %q: %+v", name, err)
	}

	return nil
}
