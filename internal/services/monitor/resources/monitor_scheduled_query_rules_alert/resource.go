// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_alert

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_scheduled_query_rules_alert/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-name "AlertingActionConfigComplete"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorScheduledQueryRulesAlertCreateUpdate,
		Read:   resourceMonitorScheduledQueryRulesAlertRead,
		Update: resourceMonitorScheduledQueryRulesAlertCreateUpdate,
		Delete: resourceMonitorScheduledQueryRulesAlertDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&scheduledqueryrules.ScheduledQueryRuleId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&scheduledqueryrules.ScheduledQueryRuleId{}),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ScheduledQueryRulesAlertUpgradeV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: monitorScheduledQueryRulesAlertSchema(),
	}
}
