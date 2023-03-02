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
	d := MonitorDataCollectionRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("destinations.0.log_analytics.0.workspace_resource_id").Exists(),
				check.That(data.ResourceName).Key("destinations.0.log_analytics.0.name").HasValue("test-destination-log1"),
				check.That(data.ResourceName).Key("destinations.0.log_analytics.1.workspace_resource_id").Exists(),
				check.That(data.ResourceName).Key("destinations.0.log_analytics.1.name").HasValue("test-destination-log2"),
				check.That(data.ResourceName).Key("data_flow.#").HasValue("3"),
				check.That(data.ResourceName).Key("data_flow.1.streams.#").HasValue("3"),
				check.That(data.ResourceName).Key("data_flow.1.destinations.#").HasValue("1"),
				check.That(data.ResourceName).Key("data_flow.2.streams.0").HasValue("Microsoft-Event"),
				check.That(data.ResourceName).Key("data_flow.2.destinations.0").HasValue("test-destination-log1"),
				check.That(data.ResourceName).Key("data_sources.0.syslog.0.facility_names.#").HasValue("5"),
				check.That(data.ResourceName).Key("data_sources.0.syslog.0.streams.#").HasValue("2"),
				check.That(data.ResourceName).Key("data_sources.0.performance_counter.#").HasValue("2"),
				check.That(data.ResourceName).Key("data_sources.0.performance_counter.1.sampling_frequency_in_seconds").HasValue("20"),
				check.That(data.ResourceName).Key("data_sources.0.performance_counter.1.name").HasValue("test-datasource-perfcounter2"),
				check.That(data.ResourceName).Key("data_sources.0.windows_event_log.0.x_path_queries.0").HasValue("System!*[System[EventID=4648]]"),
				check.That(data.ResourceName).Key("data_sources.0.extension.0.extension_json").Exists(),
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
