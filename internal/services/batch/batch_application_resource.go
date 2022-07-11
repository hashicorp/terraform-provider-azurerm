package batch

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2022-01-01/batch"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceBatchApplication() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBatchApplicationCreate,
		Read:   resourceBatchApplicationRead,
		Update: resourceBatchApplicationUpdate,
		Delete: resourceBatchApplicationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ApplicationID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApplicationName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName,
			},

			"allow_updates": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"default_version": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.ApplicationVersion,
			},

			"display_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.ApplicationDisplayName,
			},
		},
	}
}

func resourceBatchApplicationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.ApplicationClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewApplicationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.BatchAccountName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
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

	if _, err := client.Create(ctx, id.ResourceGroup, id.BatchAccountName, id.Name, &parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceBatchApplicationRead(d, meta)
}

func resourceBatchApplicationRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("reading Batch Application %q (Account Name %q / Resource Group %q): %+v", id.Name, id.BatchAccountName, id.ResourceGroup, err)
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

func resourceBatchApplicationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("updating Batch Application %q (Account Name %q / Resource Group %q): %+v", id.Name, id.BatchAccountName, id.ResourceGroup, err)
	}

	return resourceBatchApplicationRead(d, meta)
}

func resourceBatchApplicationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.ApplicationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApplicationID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.BatchAccountName, id.Name); err != nil {
		return fmt.Errorf("deleting Batch Application %q (Account Name %q / Resource Group %q): %+v", id.Name, id.BatchAccountName, id.ResourceGroup, err)
	}

	return nil
}
