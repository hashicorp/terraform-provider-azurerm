// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/flowlogs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkwatchers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = NetworkWatcherFlowLogV0ToV1{}

type NetworkWatcherFlowLogV0ToV1 struct{}

func (NetworkWatcherFlowLogV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"network_watcher_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"network_security_group_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"storage_account_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"retention_policy": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"days": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
				},
			},
		},

		"traffic_analytics": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"workspace_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"workspace_region": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"workspace_resource_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"interval_in_minutes": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  60,
					},
				},
			},
		},

		"version": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (NetworkWatcherFlowLogV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old
		// 	/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/networkWatchers/watcher1/networkSecurityGroupId/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/networkSecurityGroups/group1
		// new:
		// 	/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/networkWatchers/watcher1/flowLogs/log1
		oldId := rawState["id"].(string)
		parts := strings.Split(oldId, "/networkSecurityGroupId")
		if len(parts) != 2 {
			return rawState, fmt.Errorf("Error: Network Watcher Flow Log ID could not be split on `/networkSecurityGroupId`: %s", oldId)
		}
		watcherId, err := networkwatchers.ParseNetworkWatcherIDInsensitively(parts[0])
		if err != nil {
			return rawState, err
		}

		var name string
		rawName, ok := rawState["name"]
		if ok {
			name = rawName.(string)
		} else {
			// The `name` is introduced as an attribute since 0e528be. If users have provisioned this resource prior to that commit, and didn't run a `refresh` for the flow log. Then the state won't have `name` included.
			// In this case, we will use the Portal way to construct the flow log name.
			nsgId, err := parse.NetworkSecurityGroupID(parts[1])
			if err != nil {
				return rawState, err
			}
			name = fmt.Sprintf("Microsoft.Network%s%s", watcherId.ResourceGroupName, nsgId.Name)
			if len(name) > 80 {
				name = name[:80]
			}
		}
		id := flowlogs.NewFlowLogID(watcherId.SubscriptionId, watcherId.ResourceGroupName, watcherId.NetworkWatcherName, name)
		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
