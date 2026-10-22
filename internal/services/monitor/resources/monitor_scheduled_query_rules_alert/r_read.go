// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_alert

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorScheduledQueryRulesAlertRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := scheduledqueryrules.ParseScheduledQueryRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("getting Monitor %s: %+v", *id, err)
	}

	return resourceMonitorScheduledQueryRulesAlertFlatten(d, id, resp.Model)
}
