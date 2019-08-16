package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2017-10-12/cdn"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCdnProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCdnProfileCreate,
		Read:   resourceArmCdnProfileRead,
		Update: resourceArmCdnProfileUpdate,
		Delete: resourceArmCdnProfileDelete,
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

			"sku": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.StandardAkamai),
					string(cdn.StandardChinaCdn),
					string(cdn.StandardVerizon),
					string(cdn.StandardMicrosoft),
					string(cdn.PremiumVerizon),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmCdnProfileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cdn.ProfilesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM CDN Profile creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing CDN Profile %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_cdn_profile", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	sku := d.Get("sku").(string)
	tags := d.Get("tags").(map[string]interface{})

	cdnProfile := cdn.Profile{
		Location: &location,
		Tags:     expandTags(tags),
		Sku: &cdn.Sku{
			Name: cdn.SkuName(sku),
		},
	}

	future, err := client.Create(ctx, resGroup, name, cdnProfile)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read CDN Profile %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmCdnProfileRead(d, meta)
}

func resourceArmCdnProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cdn.ProfilesClient
	ctx := meta.(*ArmClient).StopContext

	if !d.HasChange("tags") {
		return nil
	}

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	newTags := d.Get("tags").(map[string]interface{})

	props := cdn.ProfileUpdateParameters{
		Tags: expandTags(newTags),
	}

	future, err := client.Update(ctx, resourceGroup, name, props)
	if err != nil {
		return fmt.Errorf("Error issuing update request for CDN Profile %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the update of CDN Profile %q (Resource Group %q) to commplete: %+v", name, resourceGroup, err)
	}

	return resourceArmCdnProfileRead(d, meta)
}

func resourceArmCdnProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cdn.ProfilesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["profiles"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure CDN Profile %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmCdnProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cdn.ProfilesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["profiles"]
	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing delete request for CDN Profile %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error waiting for CDN Profile %q (Resource Group %q) to be deleted: %+v", name, resourceGroup, err)
	}

	return err
}
