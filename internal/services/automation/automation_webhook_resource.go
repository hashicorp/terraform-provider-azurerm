// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2015-10-31/webhook"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/migration"
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
			_, err := webhook.ParseWebHookID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AutomationWebhookV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := webhook.NewWebHookID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))
	expiryTime := d.Get("expiry_time").(string)
	enabled := d.Get("enabled").(bool)
	runbookName := d.Get("runbook_name").(string)
	runOn := d.Get("run_on_worker_group").(string)
	webhookParameters := expandStringInterfaceMap(d.Get("parameters").(map[string]interface{}))

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if resp.Model != nil && resp.Model.Id != nil && *resp.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_automation_webhook", *resp.Model.Id)
		}
	}

	parameters := webhook.WebhookCreateOrUpdateParameters{
		Name: id.WebHookName,
		Properties: webhook.WebhookCreateOrUpdateProperties{
			IsEnabled:  utils.Bool(enabled),
			ExpiryTime: &expiryTime,
			Parameters: &webhookParameters,
			Runbook: &webhook.RunbookAssociationProperty{
				Name: utils.String(runbookName),
			},
			RunOn: utils.String(runOn),
		},
	}

	uri := ""
	if d.IsNewResource() {
		if v := d.Get("uri"); v != nil && v.(string) != "" {
			uri = v.(string)
			parameters.Properties.Uri = &uri
		} else {
			automationAccountId := webhook.NewAutomationAccountID(subscriptionId, id.ResourceGroupName, id.AutomationAccountName)
			resp, err := client.GenerateUri(ctx, automationAccountId)
			if err != nil {
				return fmt.Errorf("unable to generate URI for %s: %+v", id, err)
			}

			parameters.Properties.Uri = resp.Model
			if resp.Model != nil {
				uri = *resp.Model
			}
		}
	} else {
		if d.Get("uri") != nil {
			parameters.Properties.Uri = utils.String(d.Get("uri").(string))
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
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

	id, err := webhook.ParseWebHookID(d.Id())
	if err != nil {
		return err
	}
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.WebHookName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("automation_account_name", id.AutomationAccountName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {

			d.Set("expiry_time", props.ExpiryTime)
			d.Set("enabled", props.IsEnabled)

			if props.Runbook != nil && props.Runbook.Name != nil {
				d.Set("runbook_name", props.Runbook.Name)
			}
			d.Set("run_on_worker_group", props.RunOn)

			if err = d.Set("parameters", utils.FlattenPtrMapStringString(props.Parameters)); err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceAutomationWebhookDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.WebhookClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webhook.ParseWebHookID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
