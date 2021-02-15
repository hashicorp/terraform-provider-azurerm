package maps

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/maps/mgmt/2018-05-01/maps"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maps/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maps/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMapsAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceMapsAccountCreateUpdate,
		Read:   resourceMapsAccountRead,
		Update: resourceMapsAccountCreateUpdate,
		Delete: resourceMapsAccountDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AccountID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"S0",
					"S1",
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

func resourceMapsAccountCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maps.AccountsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Maps Account creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})
	sku := d.Get("sku_name").(string)

	if d.IsNewResource() {
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

	return resourceMapsAccountRead(d, meta)
}

func resourceMapsAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maps.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Maps Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}
	if props := resp.Properties; props != nil {
		d.Set("x_ms_client_id", props.XMsClientID)
	}

	keysResp, err := client.ListKeys(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error making Read Access Keys request on Maps Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	d.Set("primary_access_key", keysResp.PrimaryKey)
	d.Set("secondary_access_key", keysResp.SecondaryKey)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceMapsAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maps.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("Error deleting Maps Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}
