// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_runbook

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runbook"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	automationClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationRunbookDelete(d *pluginsdk.ResourceData, meta any) error {
	client := automationClient.From(meta).Runbook
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := runbook.ParseRunbookID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("issuing AzureRM delete request for Automation Runbook '%s': %+v", id.RunbookName, err)
	}

	return nil
}
