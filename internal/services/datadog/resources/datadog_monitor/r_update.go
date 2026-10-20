// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/monitorsresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatadogMonitorUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Datadog.MonitorsResource
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := monitorsresource.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	payload := monitorsresource.DatadogMonitorResourceUpdateParameters{
		Properties: &monitorsresource.MonitorUpdateProperties{},
	}
	if d.HasChange("sku_name") {
		payload.Sku = &monitorsresource.ResourceSku{
			Name: d.Get("sku_name").(string),
		}
	}
	if d.HasChange("monitoring_enabled") {
		monitoringStatus := monitorsresource.MonitoringStatusDisabled
		if d.Get("monitoring_enabled").(bool) {
			monitoringStatus = monitorsresource.MonitoringStatusEnabled
		}
		payload.Properties.MonitoringStatus = pointer.To(monitoringStatus)
	}
	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.MonitorsUpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}
	return resourceDatadogMonitorRead(d, meta)
}
