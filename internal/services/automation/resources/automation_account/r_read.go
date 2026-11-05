// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_account

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/agentregistrationinformation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/automationaccount"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationAccountRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.AutomationAccount
	registrationClient := meta.(*clients.Client).Automation.AgentRegistrationInfoClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := automationaccount.ParseAutomationAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	infoId := agentregistrationinformation.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName)
	keysResp, err := registrationClient.Get(ctx, infoId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Agent Registration Info for %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Registration Info for %s: %+v", *id, err)
	}

	return resourceAutomationAccountFlatten(d, id, resp.Model, keysResp.Model)
}
