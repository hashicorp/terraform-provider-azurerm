package tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"testing"
)

func TestAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_Enable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_linux_syslog_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_Enable(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "Enabled"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_Disable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_linux_syslog_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_Disable(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "Disabled"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_linux_syslog_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_Enable(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "Enabled"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_Disable(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "Disabled"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_Enable(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "Enabled"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_linux_syslog_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_Enable(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogCollectionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_requiresImport),
		},
	})
}

func testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogCollectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.DataSourcesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Log Analytics Data Source Linux syslog collection not found: %s", resourceName)
		}

		id, err := parse.LogAnalyticsDataSourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Workspace, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Log Analytics Data Source Linux syslog collection %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("failed to get on LogAnalytics.DataSources: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogCollectionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.DataSourcesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_datasource_linux_syslog_collection" {
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

func testAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_Enable(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_template(data)
	return fmt.Sprintf(` %s
resource "azurerm_log_analytics_datasource_linux_syslog_collection" "test" {
  name                = "acctestLADS-WE-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
  state               = "Enabled"
}
`,template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_Disable(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_template(data)
	return fmt.Sprintf(`%s

resource "azurerm_log_analytics_datasource_linux_syslog_collection" "test" {
  name                = "acctestLADS-WE-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
  state         = "Disabled"
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_Enable(data)
	return fmt.Sprintf(`%s

resource "azurerm_log_analytics_datasource_linux_syslog_collection" "import" {
  name                = azurerm_log_analytics_datasource_linux_syslog_collection.test.name
  resource_group_name = azurerm_log_analytics_datasource_linux_syslog_collection.test.resource_group_name
  workspace_name      = azurerm_log_analytics_datasource_linux_syslog_collection.test.workspace_name
  state      = azurerm_log_analytics_datasource_linux_syslog_collection.test.state
}
`, template)
}

func testAccAzureRMLogAnalyticsDataSourceLinuxSyslogCollection_template(data acceptance.TestData) string {
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
