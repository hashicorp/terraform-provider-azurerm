// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_watcher_flow_log

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/flowlogs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/network_watcher_flow_log/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-sequential -parent-resource network_watcher

const azureNetworkWatcherFlowLogResourceName = "azurerm_network_watcher_flow_log"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkWatcherFlowLogCreate,
		Read:   resourceNetworkWatcherFlowLogRead,
		Update: resourceNetworkWatcherFlowLogUpdate,
		Delete: resourceNetworkWatcherFlowLogDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.NetworkWatcherFlowLogV0ToV1{},
		}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&flowlogs.FlowLogId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&flowlogs.FlowLogId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: networkWatcherFlowLogSchema(),

		CustomizeDiff: func(_ context.Context, d *pluginsdk.ResourceDiff, _ any) error {
			if d.Id() == "" {
				if _, err := networksecuritygroups.ParseNetworkSecurityGroupID(d.Get("target_resource_id").(string)); err == nil {
					return errors.New("creation of new NSG flow logs is no longer supported by Azure as of June 30, 2025. NSG flow logs will be retired on September 30, 2027. For more information, see https://learn.microsoft.com/azure/network-watcher/nsg-flow-logs-migrate")
				}
			}

			return nil
		},
	}
}
