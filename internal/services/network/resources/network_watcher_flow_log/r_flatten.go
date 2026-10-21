// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_watcher_flow_log

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/flowlogs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceNetworkWatcherFlowLogFlatten(d *pluginsdk.ResourceData, id *flowlogs.FlowLogId, model *flowlogs.FlowLog) error {
	d.Set("name", id.FlowLogName)
	d.Set("network_watcher_name", id.NetworkWatcherName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("traffic_analytics", flattenNetworkWatcherFlowLogTrafficAnalytics(props.FlowAnalyticsConfiguration)); err != nil {
				return fmt.Errorf("setting `traffic_analytics`: %+v", err)
			}

			d.Set("enabled", props.Enabled)

			version := 0
			if format := props.Format; format != nil {
				version = int(*format.Version)
			}
			d.Set("version", version)

			// Azure API returns "" when flow log is disabled
			// Don't overwrite to prevent storage account ID diff when that is the case
			if props.StorageId != "" {
				d.Set("storage_account_id", props.StorageId)
			}

			targetResourceId := props.TargetResourceId
			if nsgId, err := networksecuritygroups.ParseNetworkSecurityGroupIDInsensitively(props.TargetResourceId); err == nil {
				targetResourceId = nsgId.ID()
			} else if vnetId, err := commonids.ParseVirtualNetworkIDInsensitively(props.TargetResourceId); err == nil {
				targetResourceId = vnetId.ID()
			} else if subnetId, err := commonids.ParseSubnetIDInsensitively(props.TargetResourceId); err == nil {
				targetResourceId = subnetId.ID()
			} else if nicId, err := commonids.ParseNetworkInterfaceIDInsensitively(props.TargetResourceId); err == nil {
				targetResourceId = nicId.ID()
			}

			d.Set("target_resource_id", targetResourceId)

			if err := d.Set("retention_policy", flattenNetworkWatcherFlowLogRetentionPolicy(props.RetentionPolicy)); err != nil {
				return fmt.Errorf("setting `retention_policy`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenNetworkWatcherFlowLogRetentionPolicy(input *flowlogs.RetentionPolicyParameters) []any {
	output := make([]any, 0)

	if input != nil {
		days := 0
		if input.Days != nil {
			days = int(*input.Days)
		}
		output = append(output, map[string]any{
			"days":    days,
			"enabled": pointer.From(input.Enabled),
		})
	}

	return output
}

func flattenNetworkWatcherFlowLogTrafficAnalytics(input *flowlogs.TrafficAnalyticsProperties) []any {
	output := make([]any, 0)
	if input != nil {
		if cfg := input.NetworkWatcherFlowAnalyticsConfiguration; cfg != nil {
			intervalInMinutes := 0
			if cfg.TrafficAnalyticsInterval != nil {
				intervalInMinutes = int(*cfg.TrafficAnalyticsInterval)
			}
			output = append(output, map[string]any{
				"enabled":               pointer.From(cfg.Enabled),
				"interval_in_minutes":   intervalInMinutes,
				"workspace_id":          pointer.From(cfg.WorkspaceId),
				"workspace_region":      pointer.From(cfg.WorkspaceRegion),
				"workspace_resource_id": pointer.From(cfg.WorkspaceResourceId),
			})
		}
	}

	return output
}
