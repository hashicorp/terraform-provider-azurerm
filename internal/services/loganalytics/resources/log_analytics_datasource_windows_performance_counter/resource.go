// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_datasource_windows_performance_counter

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/datasources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_datasource_windows_performance_counter/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsDataSourceWindowsPerformanceCounterCreateUpdate,
		Read:   resourceLogAnalyticsDataSourceWindowsPerformanceCounterRead,
		Update: resourceLogAnalyticsDataSourceWindowsPerformanceCounterCreateUpdate,
		Delete: resourceLogAnalyticsDataSourceWindowsPerformanceCounterDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := datasources.ParseDataSourceID(id)
			return err
		}, helpers.ImportLogAnalyticsDataSource(datasources.DataSourceKindWindowsPerformanceCounter)),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WindowsPerformanceCounterV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: logAnalyticsDatasourceWindowsPerformanceCounterSchema(),
	}
}
