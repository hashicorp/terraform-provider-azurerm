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

func TestAccAzureRMLogAnalyticsDataSourceWindowsEvent_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_windows_event", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceWindowsEventDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceWindowsEvent_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceWindowsEventExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceWindowsEvent_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_windows_event", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceWindowsEventDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceWindowsEvent_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceWindowsEventExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceWindowsEvent_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_windows_event", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceWindowsEventDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceWindowsEvent_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceWindowsEventExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsDataSourceWindowsEvent_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceWindowsEventExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsDataSourceWindowsEvent_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceWindowsEventExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceWindowsEvent_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_windows_event", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceWindowsEventDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceWindowsEvent_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceWindowsEventExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMLogAnalyticsDataSourceWindowsEvent_requiresImport),
		},
	})
}

func testCheckAzureRMLogAnalyticsDataSourceWindowsEventExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.DataSourcesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Log Analytics Data Source Windows Event not found: %s", resourceName)
		}

		id, err := parse.LogAnalyticsDataSourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Workspace, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Log Analytics Data Source Windows Event %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("failed to get on LogAnalytics.DataSources: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMLogAnalyticsDataSourceWindowsEventDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.DataSourcesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_datasource_windows_event" {
			continue
		}

		id, err := parse.LogAnalyticsDataSourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Workspace, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("failed to get on LogAnalytics.DataSources: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMLogAnalyticsDataSourceWindowsEvent_basic(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceWindowsEvent_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_datasource_windows_event" "test" {
  name                = "acctestLADS-WE-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
  event_log_name      = "Application"
  event_types         = ["error"]
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsDataSourceWindowsEvent_complete(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceWindowsEvent_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_datasource_windows_event" "test" {
  name                = "acctestLADS-WE-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
  event_log_name      = "Application"
  event_types         = ["InforMation", "warning", "Error"]
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsDataSourceWindowsEvent_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceWindowsEvent_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_datasource_windows_event" "import" {
  name                = azurerm_log_analytics_datasource_windows_event.test.name
  resource_group_name = azurerm_log_analytics_datasource_windows_event.test.resource_group_name
  workspace_name      = azurerm_log_analytics_datasource_windows_event.test.workspace_name
  event_log_name      = azurerm_log_analytics_datasource_windows_event.test.event_log_name
  event_types         = azurerm_log_analytics_datasource_windows_event.test.event_types
}
`, template)
}

func testAccAzureRMLogAnalyticsDataSourceWindowsEvent_template(data acceptance.TestData) string {
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
