// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_webhook

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2015-10-31/webhook"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationWebhookRead(d *pluginsdk.ResourceData, meta any) error {
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

			if err = d.Set("parameters", pluginsdk.FlattenPtrMapStringString(props.Parameters)); err != nil {
				return err
			}
		}
	}

	return nil
}
