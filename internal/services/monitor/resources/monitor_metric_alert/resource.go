// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_metric_alert

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-03-01/metricalerts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_metric_alert/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorMetricAlertCreateUpdate,
		Read:   resourceMonitorMetricAlertRead,
		Update: resourceMonitorMetricAlertCreateUpdate,
		Delete: resourceMonitorMetricAlertDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&metricalerts.MetricAlertId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&metricalerts.MetricAlertId{}),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.MetricAlertUpgradeV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: monitorMetricAlertSchema(),
	}
}
