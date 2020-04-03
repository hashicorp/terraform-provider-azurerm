package datashare

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datashare/mgmt/2019-11-01/datashare"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataShareAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataShareAccountCreate,
		Read:   resourceArmDataShareAccountRead,
		Update: resourceArmDataShareAccountUpdate,
		Delete: resourceArmDataShareAccountDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DataShareAccountID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"location":            azure.SchemaLocation(),
			"resource_group_name": azure.SchemaResourceGroupName(),
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Accepted",
					"InProgress",
					"TransientFailure",
					"Succeeded",
					"Failed",
					"Canceled",
				}, false),
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"error": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:     schema.TypeString,
							Required: true,
						},
						"message": {
							Type:     schema.TypeString,
							Required: true,
						},
						"data_share_error_info_details": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validation.StringIsJSON,
							DiffSuppressFunc: structure.SuppressJsonDiff,
						},
						"target": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tags.Schema(),
		},
	}
}
func resourceArmDataShareAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.AccountClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing  DataShare Account %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_share_account", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	account := datashare.Account{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}

	props := datashare.AccountProperties{}
	account.AccountProperties = &props
	if _, err := client.Create(ctx, resourceGroup, name, account); err != nil {
		return fmt.Errorf("Error creating/updating DataShare Account %q (Resource Group %q / account %q): %+v", name, resourceGroup, account, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving DataShare Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read  DataShare Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	return resourceArmDataShareAccountRead(d, meta)
}

func resourceArmDataShareAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.AccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataShareAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] DataShare %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving DataShare Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("name", id.Name)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if name := resp.Name; name != nil {
		d.Set("name", name)
	}
	if props := resp.AccountProperties; props != nil {
		d.Set("created_at", props.CreatedAt.Format(time.RFC3339))
		d.Set("user_email", props.UserEmail)
		d.Set("user_name", props.UserName)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDataShareAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.AccountClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	t := d.Get("tags").(map[string]interface{})

	accountUpdateParameters := datashare.AccountUpdateParameters{
		Tags: tags.Expand(t),
	}

	if _, err := client.Update(ctx, resourceGroup, name, accountUpdateParameters); err != nil {
		return fmt.Errorf("Error creating/updating DataShare Account %q (Resource Group %q / accountUpdateParameters %q): %+v", name, resourceGroup, accountUpdateParameters, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving DataShare Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read  DataShare Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	return resourceArmDataShareAccountRead(d, meta)
}

func resourceArmDataShareAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.AccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataShareAccountID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("Error deleting DataShare Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return nil
}
