// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_metric_alert

import (
	"bytes"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-03-01/metricalerts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorMetricAlertRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.MetricAlertsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := metricalerts.ParseMetricAlertID(d.Id())
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
		return fmt.Errorf("getting Monitor %s: %+v", *id, err)
	}

	if err := resourceMonitorMetricAlertFlatten(d, id, resp.Model); err != nil {
		return err
	}

	return nil
}

func resourceMonitorMetricAlertActionHash(input any) int {
	var buf bytes.Buffer
	if v, ok := input.(map[string]any); ok {
		fmt.Fprintf(&buf, "%s-", v["action_group_id"].(string))
		if m, ok := v["webhook_properties"].(map[string]any); ok && m != nil {
			fmt.Fprintf(&buf, "%v-", m)
		}
	}
	return pluginsdk.HashString(buf.String())
}
