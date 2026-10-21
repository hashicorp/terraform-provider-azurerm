// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_watcher_flow_log

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/flowlogs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkwatchers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkWatcherFlowLogCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.FlowLogs
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := flowlogs.NewFlowLogID(subscriptionId, d.Get("resource_group_name").(string), d.Get("network_watcher_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError(azureNetworkWatcherFlowLogResourceName, id.ID())
	}

	targetResourceId := d.Get("target_resource_id").(string)

	locks.ByID(targetResourceId)
	defer locks.UnlockByID(targetResourceId)

	loc := d.Get("location").(string)
	if loc == "" {
		// Get the containing network watcher in order to reuse its location if the "location" is not specified.
		watcherClient := meta.(*clients.Client).Network.NetworkWatchers
		watcherId := networkwatchers.NewNetworkWatcherID(id.SubscriptionId, id.ResourceGroupName, id.NetworkWatcherName)
		resp, err := watcherClient.Get(ctx, watcherId)
		if err != nil {
			return fmt.Errorf("retrieving %s: %v", watcherId, err)
		}
		if model := resp.Model; model != nil && model.Location != nil {
			loc = *model.Location
		}
	}

	parameters := flowlogs.FlowLog{
		Location: pointer.To(location.Normalize(loc)),
		Properties: &flowlogs.FlowLogPropertiesFormat{
			TargetResourceId: targetResourceId,
			StorageId:        d.Get("storage_account_id").(string),
			Enabled:          pointer.To(d.Get("enabled").(bool)),
			RetentionPolicy:  expandNetworkWatcherFlowLogRetentionPolicy(d.Get("retention_policy").([]any)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if _, ok := d.GetOk("traffic_analytics"); ok {
		parameters.Properties.FlowAnalyticsConfiguration = expandNetworkWatcherFlowLogTrafficAnalytics(d)
	}

	if version, ok := d.GetOk("version"); ok {
		parameters.Properties.Format = &flowlogs.FlowLogFormatParameters{
			Version: pointer.To(int64(version.(int))),
		}
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceNetworkWatcherFlowLogRead(d, meta)
}

func expandNetworkWatcherFlowLogRetentionPolicy(input []any) *flowlogs.RetentionPolicyParameters {
	if len(input) < 1 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)
	enabled := v["enabled"].(bool)
	days := v["days"].(int)

	return &flowlogs.RetentionPolicyParameters{
		Enabled: pointer.To(enabled),
		Days:    pointer.To(int64(days)),
	}
}

func expandNetworkWatcherFlowLogTrafficAnalytics(d *pluginsdk.ResourceData) *flowlogs.TrafficAnalyticsProperties {
	vs := d.Get("traffic_analytics").([]any)

	if len(vs) == 0 {
		return nil
	}

	v := vs[0].(map[string]any)
	enabled := v["enabled"].(bool)
	workspaceID := v["workspace_id"].(string)
	workspaceRegion := v["workspace_region"].(string)
	workspaceResourceID := v["workspace_resource_id"].(string)
	interval := v["interval_in_minutes"].(int)

	return &flowlogs.TrafficAnalyticsProperties{
		NetworkWatcherFlowAnalyticsConfiguration: &flowlogs.TrafficAnalyticsConfigurationProperties{
			Enabled:                  pointer.To(enabled),
			WorkspaceId:              pointer.To(workspaceID),
			WorkspaceRegion:          pointer.To(workspaceRegion),
			WorkspaceResourceId:      pointer.To(workspaceResourceID),
			TrafficAnalyticsInterval: pointer.To(int64(interval)),
		},
	}
}
