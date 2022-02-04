package automation

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAutomationWebhook() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationWebhookCreateUpdate,
		Read:   resourceAutomationWebhookRead,
		Update: resourceAutomationWebhookCreateUpdate,
		Delete: resourceAutomationWebhookDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WebhookID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"automation_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AutomationAccount(),
			},

			"expiry_time": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validation.IsRFC3339Time,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"runbook_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.RunbookName(),
			},
			"run_on_worker_group": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
			"parameters": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"uri": {
				Type:         pluginsdk.TypeString,
				ConfigMode:   schema.SchemaConfigModeAttr,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				Sensitive:    true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
		},
	}
}

func resourceAutomationWebhookCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.WebhookClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewWebhookID(client.SubscriptionID, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))
	expiryTime := d.Get("expiry_time").(string)
	enabled := d.Get("enabled").(bool)
	runbookName := d.Get("runbook_name").(string)
	runOn := d.Get("run_on_worker_group").(string)
	webhookParameters := utils.ExpandMapStringPtrString(d.Get("parameters").(map[string]interface{}))

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for present of existing %s: %+v", id, err)
			}
		}

		if resp.ID != nil && *resp.ID != "" {
			return tf.ImportAsExistsError("azurerm_automation_webhook", *resp.ID)
		}
	}

	t, _ := time.Parse(time.RFC3339, expiryTime) // should be validated by the schema
	parameters := automation.WebhookCreateOrUpdateParameters{
		Name: utils.String(id.Name),
		WebhookCreateOrUpdateProperties: &automation.WebhookCreateOrUpdateProperties{
			IsEnabled:  utils.Bool(enabled),
			ExpiryTime: &date.Time{Time: t},
			Parameters: webhookParameters,
			Runbook: &automation.RunbookAssociationProperty{
				Name: utils.String(runbookName),
			},
			RunOn: utils.String(runOn),
		},
	}

	uri := ""
	if d.IsNewResource() {
		if v := d.Get("uri"); v != nil && v.(string) != "" {
			uri = v.(string)
		} else {
			resp, err := client.GenerateURI(ctx, id.ResourceGroup, id.AutomationAccountName)
			if err != nil {
				return fmt.Errorf("unable to generate URI for %s: %+v", id, err)
			}
			parameters.WebhookCreateOrUpdateProperties.URI = resp.Value
			uri = *resp.Value
		}
	} else {
		if d.Get("uri") != nil {
			parameters.WebhookCreateOrUpdateProperties.URI = utils.String(d.Get("uri").(string))
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	// URI is not present in the response from Azure, so it's set now, as there was no error returned
	if uri != "" {
		d.Set("uri", uri)
	}
	return resourceAutomationWebhookRead(d, meta)
}

func resourceAutomationWebhookRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.WebhookClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebhookID(d.Id())
	if err != nil {
		return err
	}
	resp, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Automation Webhook %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Automation Webhook %q (Automation Account Name %q / Resource Group %q): %+v", id.Name, id.AutomationAccountName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("automation_account_name", id.AutomationAccountName)
	d.Set("expiry_time", resp.ExpiryTime.String())
	d.Set("enabled", resp.IsEnabled)
	if resp.Runbook != nil {
		d.Set("runbook_name", resp.Runbook.Name)
	}
	d.Set("run_on_worker_group", resp.RunOn)
	if err = d.Set("parameters", utils.FlattenMapStringPtrString(resp.Parameters)); err != nil {
		return err
	}

	return nil
}

func resourceAutomationWebhookDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.WebhookClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebhookID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name); err != nil {
		return fmt.Errorf("deleting Automation Webhook %q (Automation Account Name %q / Resource Group %q): %+v", id.Name, id.AutomationAccountName, id.ResourceGroup, err)
	}

	return nil
}
