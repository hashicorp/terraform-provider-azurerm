// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_webhook

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2015-10-31/webhook"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationWebhookCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.WebhookClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := webhook.NewWebHookID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
			IsEnabled:  pointer.To(d.Get("enabled").(bool)),
			ExpiryTime: pointer.To(d.Get("expiry_time").(string)),
			Parameters: pointer.To(helpers.ExpandStringInterfaceMap(d.Get("parameters").(map[string]any))),
			Runbook: &webhook.RunbookAssociationProperty{
				Name: pointer.To(d.Get("runbook_name").(string)),
			},
			RunOn: pointer.To(d.Get("run_on_worker_group").(string)),
		},
	}

	uri := ""
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
