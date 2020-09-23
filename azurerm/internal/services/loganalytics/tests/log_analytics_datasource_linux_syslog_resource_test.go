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

func TestAccAzureRMLogAnalyticsDataSourceLinuxSyslog_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_linux_syslog", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxSyslog_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "syslog_name", "mail"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceLinuxSyslog_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_linux_syslog", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxSyslog_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "syslog_name", "auth"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceLinuxSyslog_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_linux_syslog", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxSyslog_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "syslog_name", "mail"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxSyslog_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "syslog_name", "auth"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxSyslog_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "syslog_name", "mail"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceLinuxSyslog_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_linux_syslog", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxSyslog_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMLogAnalyticsDataSourceLinuxSyslog_requiresImport),
		},
	})
}

func testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.DataSourcesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Log Analytics Data Source Linux syslog not found: %s", resourceName)
		}

		id, err := parse.LogAnalyticsDataSourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Workspace, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Log Analytics Data Source Linux syslog %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("failed to get on LogAnalytics.DataSources: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMLogAnalyticsDataSourceLinuxSyslogDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.DataSourcesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_datasource_linux_syslog" {
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

func testAccAzureRMLogAnalyticsDataSourceLinuxSyslog_basic(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceLinuxSyslog_template(data)
	return fmt.Sprintf(`%s

resource "azurerm_log_analytics_datasource_linux_syslog" "test" {
  name                = "acctestLADS-WE-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
  syslog_name      = "mail"
  syslog_severities         = ["emerg", "alert", "WARNING"]
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsDataSourceLinuxSyslog_complete(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceLinuxSyslog_template(data)
	return fmt.Sprintf(`%s

resource "azurerm_log_analytics_datasource_linux_syslog" "test" {
  name                = "acctestLADS-WE-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
  syslog_name      = "auth"
  syslog_severities         = ["emerg", "alert", "crit", "err", "warning", "notice", "info", "debug"]
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsDataSourceLinuxSyslog_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceLinuxSyslog_basic(data)
	return fmt.Sprintf(`%s

resource "azurerm_log_analytics_datasource_linux_syslog" "import" {
  name                = azurerm_log_analytics_datasource_linux_syslog.test.name
  resource_group_name = azurerm_log_analytics_datasource_linux_syslog.test.resource_group_name
  workspace_name      = azurerm_log_analytics_datasource_linux_syslog.test.workspace_name
  syslog_name         = azurerm_log_analytics_datasource_linux_syslog.test.syslog_name
  syslog_severities   = azurerm_log_analytics_datasource_linux_syslog.test.syslog_severities
}
`, template)
}

func testAccAzureRMLogAnalyticsDataSourceLinuxSyslog_template(data acceptance.TestData) string {
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
