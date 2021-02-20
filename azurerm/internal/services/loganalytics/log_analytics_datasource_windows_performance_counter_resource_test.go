package loganalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LogAnalyticsDataSourceWindowsPerformanceCounterResource struct {
}

func TestAccLogAnalyticsDataSourceWindowsPerformanceCounter_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_windows_performance_counter", "test")
	r := LogAnalyticsDataSourceWindowsPerformanceCounterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("object_name").HasValue("CPU"),
				check.That(data.ResourceName).Key("instance_name").HasValue("*"),
				check.That(data.ResourceName).Key("counter_name").HasValue("CPU"),
				check.That(data.ResourceName).Key("interval_seconds").HasValue("10"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsDataSourceWindowsPerformanceCounter_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_windows_performance_counter", "test")
	r := LogAnalyticsDataSourceWindowsPerformanceCounterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("object_name").HasValue("Mem"),
				check.That(data.ResourceName).Key("instance_name").HasValue("inst1"),
				check.That(data.ResourceName).Key("counter_name").HasValue("Mem"),
				check.That(data.ResourceName).Key("interval_seconds").HasValue("20"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsDataSourceWindowsPerformanceCounter_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_windows_performance_counter", "test")
	r := LogAnalyticsDataSourceWindowsPerformanceCounterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("object_name").HasValue("CPU"),
				check.That(data.ResourceName).Key("instance_name").HasValue("*"),
				check.That(data.ResourceName).Key("counter_name").HasValue("CPU"),
				check.That(data.ResourceName).Key("interval_seconds").HasValue("10"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("object_name").HasValue("Mem"),
				check.That(data.ResourceName).Key("instance_name").HasValue("inst1"),
				check.That(data.ResourceName).Key("counter_name").HasValue("Mem"),
				check.That(data.ResourceName).Key("interval_seconds").HasValue("20"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsDataSourceWindowsPerformanceCounter_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_windows_performance_counter", "test")
	r := LogAnalyticsDataSourceWindowsPerformanceCounterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t LogAnalyticsDataSourceWindowsPerformanceCounterResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LogAnalyticsDataSourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.LogAnalytics.DataSourcesClient.Get(ctx, id.ResourceGroup, id.Workspace, id.Name)
	if err != nil {
		return nil, fmt.Errorf("readingLog Analytics Data Source Windows Event (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r LogAnalyticsDataSourceWindowsPerformanceCounterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_datasource_windows_performance_counter" "test" {
  name                = "acctestLADS-WPC-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
  object_name         = "CPU"
  instance_name       = "*"
  counter_name        = "CPU"
  interval_seconds    = 10
}
`, r.template(data), data.RandomInteger)
}

func (r LogAnalyticsDataSourceWindowsPerformanceCounterResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_datasource_windows_performance_counter" "test" {
  name                = "acctestLADS-WPC-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
  object_name         = "Mem"
  instance_name       = "inst1"
  counter_name        = "Mem"
  interval_seconds    = 20
}
`, r.template(data), data.RandomInteger)
}

func (r LogAnalyticsDataSourceWindowsPerformanceCounterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_datasource_windows_performance_counter" "import" {
  name                = azurerm_log_analytics_datasource_windows_performance_counter.test.name
  resource_group_name = azurerm_log_analytics_datasource_windows_performance_counter.test.resource_group_name
  workspace_name      = azurerm_log_analytics_datasource_windows_performance_counter.test.workspace_name
  object_name         = azurerm_log_analytics_datasource_windows_performance_counter.test.object_name
  instance_name       = azurerm_log_analytics_datasource_windows_performance_counter.test.instance_name
  counter_name        = azurerm_log_analytics_datasource_windows_performance_counter.test.counter_name
  interval_seconds    = azurerm_log_analytics_datasource_windows_performance_counter.test.interval_seconds
}
`, r.basic(data))
}

func (LogAnalyticsDataSourceWindowsPerformanceCounterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-la-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
