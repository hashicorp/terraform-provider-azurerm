package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/msi/mgmt/2015-08-31-preview/msi"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmUserAssignedIdentity() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmUserAssignedIdentityCreate,
		Read:   resourceArmUserAssignedIdentityRead,
		Update: resourceArmUserAssignedIdentityCreate,
		Delete: resourceArmUserAssignedIdentityDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameSchema(),
			"location":            locationSchema(),
			"tags":                tagsSchema(),
		},
	}
}

func resourceArmUserAssignedIdentityCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).userAssignedIdentitiesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM user identity creation.")

	resourceName := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})
	identity := msi.Identity{
		Name:     &resourceName,
		Location: &location,
		Tags:     expandTags(tags),
	}

	identity, err := client.CreateOrUpdate(ctx, resGroup, resourceName, identity)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating User Assigned Identity %q (Resource Group %q): %+v", resourceName, resGroup, err)
	}

	if identity.ID == nil {
		return fmt.Errorf("Cannot read User Assigned Identity %q ID (resource group %q) ID", resourceName, resGroup)
	}

	d.SetId(*identity.ID)

	return resourceArmUserAssignedIdentityRead(d, meta)
}

func resourceArmUserAssignedIdentityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).userAssignedIdentitiesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	resourceName := id.Path["userAssignedIdentities"]

	resp, err := client.Get(ctx, resGroup, resourceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on User Assigned Identity %q (Resource Group %q): %+v", resourceName, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", resp.Location)

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmUserAssignedIdentityDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).userAssignedIdentitiesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	resourceName := id.Path["userAssignedIdentities"]

	_, err = client.Delete(ctx, resGroup, resourceName)
	if err != nil {
		return fmt.Errorf("Error deleting User Assigned Identity %q (Resource Group %q): %+v", resourceName, resGroup, err)
	}

	return nil
}
