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

func TestAccAzureRMLogAnalyticsDataSourceIIS_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_iis", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceIISDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceIIS_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceIISExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceIIS_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_iis", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceIISDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceIIS_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceIISExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceIIS_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_iis", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceIISDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceIIS_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceIISExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsDataSourceIIS_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceIISExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsDataSourceIIS_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceIISExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceIIS_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_iis", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceIISDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceIIS_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceIISExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMLogAnalyticsDataSourceIIS_requiresImport),
		},
	})
}

func testCheckAzureRMLogAnalyticsDataSourceIISExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.DataSourcesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Log Analytics Data Source IIS not found: %s", resourceName)
		}

		id, err := parse.LogAnalyticsDataSourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Workspace, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Log Analytics Data Source IIS %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("failed to get on LogAnalytics.DataSources: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMLogAnalyticsDataSourceIISDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.DataSourcesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_datasource_iis" {
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

func testAccAzureRMLogAnalyticsDataSourceIIS_basic(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceIIS_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_datasource_iis" "test" {
  name = "acctestLADS-IIS-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_name = azurerm_log_analytics_workspace.test.name
  state = "OnPremiseEnabled"
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsDataSourceIIS_complete(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceIIS_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_datasource_iis" "test" {
  name = "acctestLADS-IIS-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_name = azurerm_log_analytics_workspace.test.name
  state = "OnPremiseDisabled"
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsDataSourceIIS_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceIIS_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_datasource_iis" "import" {
  name = azurerm_log_analytics_datasource_iis.test.name
  resource_group_name = azurerm_log_analytics_datasource_iis.test.resource_group_name
  workspace_name = azurerm_log_analytics_datasource_iis.test.workspace_name
  state = azurerm_log_analytics_datasource_iis.test.state
}
`, template)
}

func testAccAzureRMLogAnalyticsDataSourceIIS_template(data acceptance.TestData) string {
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
