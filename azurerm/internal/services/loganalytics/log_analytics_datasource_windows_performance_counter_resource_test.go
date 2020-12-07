package loganalytics

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_windows_performance_counter", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceWindowsPerformanceCounterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceWindowsPerformanceCounterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "object_name", "CPU"),
					resource.TestCheckResourceAttr(data.ResourceName, "instance_name", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "counter_name", "CPU"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval_seconds", "10"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_windows_performance_counter", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceWindowsPerformanceCounterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceWindowsPerformanceCounterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "object_name", "Mem"),
					resource.TestCheckResourceAttr(data.ResourceName, "instance_name", "inst1"),
					resource.TestCheckResourceAttr(data.ResourceName, "counter_name", "Mem"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval_seconds", "20"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_windows_performance_counter", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceWindowsPerformanceCounterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceWindowsPerformanceCounterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "object_name", "CPU"),
					resource.TestCheckResourceAttr(data.ResourceName, "instance_name", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "counter_name", "CPU"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval_seconds", "10"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceWindowsPerformanceCounterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "object_name", "Mem"),
					resource.TestCheckResourceAttr(data.ResourceName, "instance_name", "inst1"),
					resource.TestCheckResourceAttr(data.ResourceName, "counter_name", "Mem"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval_seconds", "20"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_windows_performance_counter", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceWindowsPerformanceCounterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceWindowsPerformanceCounterExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_requiresImport),
		},
	})
}

func testCheckAzureRMLogAnalyticsDataSourceWindowsPerformanceCounterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.DataSourcesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Log Analytics Data Source Windows Performance Counter not found: %s", resourceName)
		}

		id, err := parse.LogAnalyticsDataSourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Workspace, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Log Analytics Data Source Windows Performance Counter %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("failed to get on LogAnalytics.DataSourcesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMLogAnalyticsDataSourceWindowsPerformanceCounterDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.DataSourcesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_datasource_windows_performance_counter" {
			continue
		}

		id, err := parse.LogAnalyticsDataSourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Workspace, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("failed to get on LogAnalytics.DataSourcesClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_basic(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_complete(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_basic(data)
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
`, template)
}

func testAccAzureRMLogAnalyticsDataSourceWindowsPerformanceCounter_template(data acceptance.TestData) string {
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
