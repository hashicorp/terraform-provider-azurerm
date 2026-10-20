// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_webhook

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2015-10-31/webhook"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationWebhookUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.WebhookClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webhook.ParseWebHookID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: model or properties was nil", *id)
	}

	existing := resp.Model.Properties

	parameters := webhook.WebhookUpdateParameters{
		Name: pointer.To(id.WebHookName),
		Properties: &webhook.WebhookUpdateProperties{
			IsEnabled:  existing.IsEnabled,
			RunOn:      existing.RunOn,
			Parameters: existing.Parameters,
		},
	}

	if d.HasChange("enabled") {
		parameters.Properties.IsEnabled = pointer.To(d.Get("enabled").(bool))
	}

	if d.HasChange("run_on_worker_group") {
		parameters.Properties.RunOn = pointer.To(d.Get("run_on_worker_group").(string))
	}

	// NOTE: parameters is a map that is replaced as a whole
	if d.HasChange("parameters") {
		parameters.Properties.Parameters = pointer.To(helpers.ExpandStringInterfaceMap(d.Get("parameters").(map[string]any)))
	}

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceAutomationWebhookRead(d, meta)
}
