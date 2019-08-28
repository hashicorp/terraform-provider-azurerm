package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/maps/mgmt/2018-05-01/maps"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	mapsint "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maps"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMapsAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMapsAccountCreateUpdate,
		Read:   resourceArmMapsAccountRead,
		Update: resourceArmMapsAccountCreateUpdate,
		Delete: resourceArmMapsAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: mapsint.ValidateName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"s0",
					"s1",
				}, false),
			},

			"tags": tags.Schema(),

			"x_ms_client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceArmMapsAccountCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).maps.AccountsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Maps Account creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})
	sku := d.Get("sku_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Maps Account %q (Resource Group %q): %+v", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_maps_account", *existing.ID)
		}
	}

	parameters := maps.AccountCreateParameters{
		Location: utils.String("global"),
		Sku: &maps.Sku{
			Name: &sku,
		},
		Tags: tags.Expand(t),
	}

	if _, err := client.CreateOrUpdate(ctx, resGroup, name, parameters); err != nil {
		return fmt.Errorf("Error creating/updating Maps Account %q (Resource Group %q) %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Maps Account %q (Resource Group %q) %+v", name, resGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Maps Account %q (Resource Group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmMapsAccountRead(d, meta)
}

func resourceArmMapsAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).maps.AccountsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["accounts"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Maps Account %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}
	if props := resp.Properties; props != nil {
		d.Set("x_ms_client_id", props.XMsClientID)
	}

	keysResp, err := client.ListKeys(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read Access Keys request on Maps Account %q (Resource Group %q): %+v", name, resGroup, err)
	}
	d.Set("primary_access_key", keysResp.PrimaryKey)
	d.Set("secondary_access_key", keysResp.SecondaryKey)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMapsAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).maps.AccountsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["accounts"]

	if _, err := client.Delete(ctx, resGroup, name); err != nil {
		return fmt.Errorf("Error deleting Maps Account %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}
