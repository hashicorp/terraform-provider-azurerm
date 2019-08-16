package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSharedImageGallery() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSharedImageGalleryCreateUpdate,
		Read:   resourceArmSharedImageGalleryRead,
		Update: resourceArmSharedImageGalleryCreateUpdate,
		Delete: resourceArmSharedImageGalleryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SharedImageGalleryName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tagsSchema(),

			"unique_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmSharedImageGalleryCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).compute.GalleriesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Image Gallery creation.")

	// TODO: support for Timeouts/Requiring Import
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	description := d.Get("description").(string)
	tags := d.Get("tags").(map[string]interface{})

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Shared Image Gallery %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_shared_image_gallery", *existing.ID)
		}
	}

	gallery := compute.Gallery{
		Location: utils.String(location),
		GalleryProperties: &compute.GalleryProperties{
			Description: utils.String(description),
		},
		Tags: expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, gallery)
	if err != nil {
		return fmt.Errorf("Error creating/updating Shared Image Gallery %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of Shared Image Gallery %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Shared Image Gallery %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Shared Image Gallery %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmSharedImageGalleryRead(d, meta)
}

func resourceArmSharedImageGalleryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).compute.GalleriesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["galleries"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Shared Image Gallery %q (Resource Group %q) was not found - removing from state", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Shared Image Gallery %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.GalleryProperties; props != nil {
		d.Set("description", props.Description)
		if identifier := props.Identifier; identifier != nil {
			d.Set("unique_name", identifier.UniqueName)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmSharedImageGalleryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).compute.GalleriesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["galleries"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		// deleted outside of Terraform
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting Shared Image Gallery %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for the deletion of Shared Image Gallery %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
