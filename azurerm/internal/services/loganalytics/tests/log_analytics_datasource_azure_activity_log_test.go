package tests

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

func TestAccAzureRMLogAnalyticsDataSourceAzureActivityLog_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_azure_activity_log", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceAzureActivityLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceAzureActivityLog_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceAzureActivityLogExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceAzureActivityLog_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_azure_activity_log", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceAzureActivityLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceAzureActivityLog_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceAzureActivityLogExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMLogAnalyticsDataSourceAzureActivityLog_requiresImport),
		},
	})
}

func testCheckAzureRMLogAnalyticsDataSourceAzureActivityLogExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.DataSourcesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Log Analytics Data Source Azure Activity Log not found: %s", resourceName)
		}

		id, err := parse.LogAnalyticsDataSourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Workspace, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Log Analytics Data Source Azure Activity Log %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Getting on LogAnalytics.DataSources: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMLogAnalyticsDataSourceAzureActivityLogDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.DataSourcesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_datasource_azure_activity_log" {
			continue
		}

		id, err := parse.LogAnalyticsDataSourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Workspace, id.Name)
		if err == nil {
			return fmt.Errorf("LogAnalytics.DataSources still exists")
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Getting on LogAnalytics.DataSources: %+v", err)
		}
		return nil
	}

	return nil
}

func testAccAzureRMLogAnalyticsDataSourceAzureActivityLog_basic(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceAzureActivityLog_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_datasource_azure_activity_log" "test" {
  name                = "acctestLADS-AAL-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsDataSourceAzureActivityLog_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceAzureActivityLog_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_datasource_azure_activity_log" "import" {
  name                = azurerm_log_analytics_datasource_azure_activity_log.test.name
  resource_group_name = azurerm_log_analytics_datasource_azure_activity_log.test.resource_group_name
  workspace_name      = azurerm_log_analytics_datasource_azure_activity_log.test.workspace_name
}
`, template)
}

func testAccAzureRMLogAnalyticsDataSourceAzureActivityLog_template(data acceptance.TestData) string {
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
