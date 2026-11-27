// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_activity_log_alert

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2020-10-01/activitylogalertsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	monitorClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorActivityLogAlertRead(d *pluginsdk.ResourceData, meta any) error {
	client := monitorClient.From(meta).ActivityLogAlertsClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := activitylogalertsapis.ParseActivityLogAlertID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ActivityLogAlertsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("getting Monitor %s: %+v", *id, err)
	}

	return resourceMonitorActivityLogAlertFlatten(d, id, resp.Model)
}
