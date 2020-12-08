package cdn

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/migration"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceCdnProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceCdnProfileCreate,
		Read:   resourceCdnProfileRead,
		Update: resourceCdnProfileUpdate,
		Delete: resourceCdnProfileDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ProfileID(id)
			return err
		}),

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

			"tags": tags.Schema(),
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    migration.CdnProfileV0Schema().CoreConfigSchema().ImpliedType(),
				Upgrade: migration.CdnProfileV0ToV1,
				Version: 0,
			},
		},
	}
}

func resourceCdnProfileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.ProfilesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM CDN Profile creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
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
	t := d.Get("tags").(map[string]interface{})

	cdnProfile := cdn.Profile{
		Location: &location,
		Tags:     tags.Expand(t),
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

	id, err := parse.ProfileID(*read.ID)
	if err != nil {
		return err
	}

	d.SetId(id.ID(""))

	return resourceCdnProfileRead(d, meta)
}

func resourceCdnProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.ProfilesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if !d.HasChange("tags") {
		return nil
	}

	id, err := parse.ProfileID(d.Id())
	if err != nil {
		return err
	}

	// name := d.Get("name").(string)
	// resourceGroup := d.Get("resource_group_name").(string)
	newTags := d.Get("tags").(map[string]interface{})

	props := cdn.ProfileUpdateParameters{
		Tags: tags.Expand(newTags),
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.Name, props)
	if err != nil {
		return fmt.Errorf("Error issuing update request for CDN Profile %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the update of CDN Profile %q (Resource Group %q) to commplete: %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceCdnProfileRead(d, meta)
}

func resourceCdnProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.ProfilesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProfileID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure CDN Profile %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceCdnProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.ProfilesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProfileID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("deleting CDN Profile %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("waiting for deletion of CDN Profile %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return err
}
