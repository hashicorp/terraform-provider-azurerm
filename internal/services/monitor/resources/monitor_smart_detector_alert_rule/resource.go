// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_smart_detector_alert_rule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-06-01/smartdetectoralertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_smart_detector_alert_rule/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorSmartDetectorAlertRuleCreateUpdate,
		Read:   resourceMonitorSmartDetectorAlertRuleRead,
		Update: resourceMonitorSmartDetectorAlertRuleCreateUpdate,
		Delete: resourceMonitorSmartDetectorAlertRuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SmartDetectorAlertRuleV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := smartdetectoralertrules.ParseSmartDetectorAlertRuleID(id)
			return err
		}),

		Schema: monitorSmartDetectorAlertRuleSchema(),
	}
}
