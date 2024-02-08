// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MonitorDataCollectionRuleDataSource struct{}

func TestAccMonitorDataCollectionRuleDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_data_collection_rule", "test")
	// https://learn.microsoft.com/en-us/azure/azure-monitor/logs/ingest-logs-event-hub#supported-regions
	data.Locations.Primary = "westeurope"
	d := MonitorDataCollectionRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("destinations.0.log_analytics.0.workspace_resource_id").Exists(),
				check.That(data.ResourceName).Key("destinations.0.log_analytics.0.name").HasValue("test-destination-log"),
				check.That(data.ResourceName).Key("data_flow.#").HasValue("4"),
				check.That(data.ResourceName).Key("data_flow.1.streams.#").HasValue("3"),
				check.That(data.ResourceName).Key("data_flow.1.destinations.#").HasValue("1"),
				check.That(data.ResourceName).Key("data_flow.2.streams.0").HasValue("Microsoft-Event"),
				check.That(data.ResourceName).Key("data_flow.2.destinations.0").HasValue("test-destination-log"),
				check.That(data.ResourceName).Key("data_flow.3.output_stream").HasValue("Microsoft-Syslog"),
				check.That(data.ResourceName).Key("data_flow.3.transform_kql").Exists(),
				check.That(data.ResourceName).Key("data_sources.0.syslog.0.facility_names.#").HasValue("5"),
				check.That(data.ResourceName).Key("data_sources.0.syslog.0.streams.#").HasValue("2"),
				check.That(data.ResourceName).Key("data_sources.0.iis_log.0.log_directories.#").HasValue("1"),
				check.That(data.ResourceName).Key("data_sources.0.log_file.0.format").HasValue("text"),
				check.That(data.ResourceName).Key("data_sources.0.data_import.0.event_hub_data_source.0.name").HasValue("test-datasource-import-event"),
				check.That(data.ResourceName).Key("data_sources.0.prometheus_forwarder.0.streams.0").HasValue("Microsoft-PrometheusMetrics"),
				check.That(data.ResourceName).Key("data_sources.0.performance_counter.#").HasValue("2"),
				check.That(data.ResourceName).Key("data_sources.0.performance_counter.1.sampling_frequency_in_seconds").HasValue("20"),
				check.That(data.ResourceName).Key("data_sources.0.performance_counter.1.name").HasValue("test-datasource-perfcounter2"),
				check.That(data.ResourceName).Key("data_sources.0.windows_event_log.0.x_path_queries.0").HasValue("System!*[System[EventID=4648]]"),
				check.That(data.ResourceName).Key("data_sources.0.extension.0.extension_json").Exists(),
				check.That(data.ResourceName).Key("immutable_id").Exists(),
				check.That(data.ResourceName).Key("stream_declaration.#").HasValue("2"),
				check.That(data.ResourceName).Key("stream_declaration.0.column.#").HasValue("3"),
			),
		},
	})
}

func (d MonitorDataCollectionRuleDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_monitor_data_collection_rule" "test" {
  name                = azurerm_monitor_data_collection_rule.test.name
  resource_group_name = azurerm_monitor_data_collection_rule.test.resource_group_name
}
`, MonitorDataCollectionRuleResource{}.complete(data))
}
