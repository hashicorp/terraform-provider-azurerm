// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_activity_log_alert

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2020-10-01/activitylogalertsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorActivityLogAlertRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.ActivityLogAlertsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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
