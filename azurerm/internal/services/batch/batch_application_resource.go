package batch

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2020-03-01/batch"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceBatchApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceBatchApplicationCreate,
		Read:   resourceBatchApplicationRead,
		Update: resourceBatchApplicationUpdate,
		Delete: resourceBatchApplicationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ApplicationID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApplicationName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName,
			},

			"allow_updates": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"default_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ApplicationVersion,
			},

			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ApplicationDisplayName,
			},
		},
	}
}

func resourceBatchApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.ApplicationClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)

	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("checking for presence of existing Batch Application %q (Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
		}
	}
	if !utils.ResponseWasNotFound(resp.Response) {
		return tf.ImportAsExistsError("azurerm_batch_application", *resp.ID)
	}

	allowUpdates := d.Get("allow_updates").(bool)
	defaultVersion := d.Get("default_version").(string)
	displayName := d.Get("display_name").(string)

	parameters := batch.Application{
		ApplicationProperties: &batch.ApplicationProperties{
			AllowUpdates:   utils.Bool(allowUpdates),
			DefaultVersion: utils.String(defaultVersion),
			DisplayName:    utils.String(displayName),
		},
	}

	if _, err := client.Create(ctx, resourceGroup, accountName, name, &parameters); err != nil {
		return fmt.Errorf("Error creating Batch Application %q (Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	resp, err = client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Batch Application %q (Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Batch Application %q (Account Name %q / Resource Group %q) ID", name, accountName, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceBatchApplicationRead(d, meta)
}

func resourceBatchApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.ApplicationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApplicationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BatchAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Batch Application %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Batch Application %q (Account Name %q / Resource Group %q): %+v", id.Name, id.BatchAccountName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.BatchAccountName)
	if applicationProperties := resp.ApplicationProperties; applicationProperties != nil {
		d.Set("allow_updates", applicationProperties.AllowUpdates)
		d.Set("default_version", applicationProperties.DefaultVersion)
		d.Set("display_name", applicationProperties.DisplayName)
	}

	return nil
}

func resourceBatchApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.ApplicationClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApplicationID(d.Id())
	if err != nil {
		return err
	}

	allowUpdates := d.Get("allow_updates").(bool)
	defaultVersion := d.Get("default_version").(string)
	displayName := d.Get("display_name").(string)

	parameters := batch.Application{
		ApplicationProperties: &batch.ApplicationProperties{
			AllowUpdates:   utils.Bool(allowUpdates),
			DefaultVersion: utils.String(defaultVersion),
			DisplayName:    utils.String(displayName),
		},
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.BatchAccountName, id.Name, parameters); err != nil {
		return fmt.Errorf("Error updating Batch Application %q (Account Name %q / Resource Group %q): %+v", id.Name, id.BatchAccountName, id.ResourceGroup, err)
	}

	return resourceBatchApplicationRead(d, meta)
}

func resourceBatchApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.ApplicationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApplicationID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.BatchAccountName, id.Name); err != nil {
		return fmt.Errorf("Error deleting Batch Application %q (Account Name %q / Resource Group %q): %+v", id.Name, id.BatchAccountName, id.ResourceGroup, err)
	}

	return nil
}
