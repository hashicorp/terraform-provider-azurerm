package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/msi/mgmt/2015-08-31-preview/msi"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmUserAssignedIdentity() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmUserAssignedIdentityCreateUpdate,
		Read:   resourceArmUserAssignedIdentityRead,
		Update: resourceArmUserAssignedIdentityCreateUpdate,
		Delete: resourceArmUserAssignedIdentityDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 30),
			Update: schema.DefaultTimeout(time.Minute * 30),
			Delete: schema.DefaultTimeout(time.Minute * 30),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 24),
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"tags": tagsSchema(),

			"principal_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmUserAssignedIdentityCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).userAssignedIdentitiesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM user identity creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		// first check if there's one in this subscription requiring import
		resp, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for the existence of User Assigned Identity %q (Resource Group %q): %+v", name, resGroup, err)
			}
		}

		if resp.ID != nil {
			return tf.ImportAsExistsError("azurerm_user_assigned_identity", *resp.ID)
		}
	}

	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})
	identity := msi.Identity{
		Name:     &name,
		Location: &location,
		Tags:     expandTags(tags),
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(tf.TimeoutForCreateUpdate(d)))
	defer cancel()
	_, err := client.CreateOrUpdate(waitCtx, resGroup, name, identity)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating User Assigned Identity %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read User Assigned Identity %q ID (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

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
	name := id.Path["userAssignedIdentities"]
	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on User Assigned Identity %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", resp.Location)

	if props := resp.IdentityProperties; props != nil {
		if principalId := props.PrincipalID; principalId != nil {
			d.Set("principal_id", principalId.String())
		}
	}

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
	name := id.Path["userAssignedIdentities"]

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	_, err = client.Delete(waitCtx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting User Assigned Identity %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}
