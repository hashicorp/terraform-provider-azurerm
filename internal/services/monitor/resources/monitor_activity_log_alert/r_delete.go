// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_activity_log_alert

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2020-10-01/activitylogalertsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorActivityLogAlertDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.ActivityLogAlertsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := activitylogalertsapis.ParseActivityLogAlertID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.ActivityLogAlertsDelete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting Monitor %s: %+v", *id, err)
		}
	}

	return nil
}
