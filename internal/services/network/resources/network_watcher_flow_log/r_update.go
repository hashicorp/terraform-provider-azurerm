// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_watcher_flow_log

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/flowlogs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkWatcherFlowLogUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.FlowLogs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := flowlogs.ParseFlowLogID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	targetResourceId := d.Get("target_resource_id").(string)

	locks.ByID(targetResourceId)
	defer locks.UnlockByID(targetResourceId)

	if d.HasChange("storage_account_id") {
		payload.Properties.StorageId = d.Get("storage_account_id").(string)
	}

	if d.HasChange("enabled") {
		payload.Properties.Enabled = pointer.To(d.Get("enabled").(bool))
	}

	if d.HasChange("retention_policy") {
		payload.Properties.RetentionPolicy = expandNetworkWatcherFlowLogRetentionPolicy(d.Get("retention_policy").([]any))
	}

	if d.HasChange("traffic_analytics") {
		payload.Properties.FlowAnalyticsConfiguration = expandNetworkWatcherFlowLogTrafficAnalytics(d)
	}

	if d.HasChange("version") {
		if version, ok := d.GetOk("version"); ok && version != 0 {
			payload.Properties.Format = &flowlogs.FlowLogFormatParameters{
				Version: pointer.To(int64(d.Get("version").(int))),
			}
		} else {
			payload.Properties.Format = &flowlogs.FlowLogFormatParameters{}
		}
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("target_resource_id") {
		payload.Properties.TargetResourceId = targetResourceId
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceNetworkWatcherFlowLogRead(d, meta)
}
