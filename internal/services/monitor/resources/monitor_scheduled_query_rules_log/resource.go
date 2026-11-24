// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_log

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_scheduled_query_rules_log/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorScheduledQueryRulesLogCreateUpdate,
		Read:   resourceMonitorScheduledQueryRulesLogRead,
		Update: resourceMonitorScheduledQueryRulesLogCreateUpdate,
		Delete: resourceMonitorScheduledQueryRulesLogDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := scheduledqueryrules.ParseScheduledQueryRuleID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ScheduledQueryRulesLogUpgradeV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: monitorScheduledQueryRulesLogSchema(),
	}
}
