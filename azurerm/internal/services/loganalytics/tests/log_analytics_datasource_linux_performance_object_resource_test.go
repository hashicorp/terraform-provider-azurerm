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

func TestAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_linux_performance_object", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceLinuxPerfObjDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxPerfObjExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "object_name", "Logical Disk"),
					resource.TestCheckResourceAttr(data.ResourceName, "instance_name", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval_seconds", "10"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_linux_performance_object", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceLinuxPerfObjDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxPerfObjExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "object_name", "Processor"),
					resource.TestCheckResourceAttr(data.ResourceName, "instance_name", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval_seconds", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "performance_counters.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_linux_performance_object", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceLinuxPerfObjDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxPerfObjExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "object_name", "Logical Disk"),
					resource.TestCheckResourceAttr(data.ResourceName, "instance_name", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval_seconds", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "performance_counters.#", "1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxPerfObjExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "object_name", "Processor"),
					resource.TestCheckResourceAttr(data.ResourceName, "instance_name", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval_seconds", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "performance_counters.#", "1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxPerfObjExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "object_name", "Logical Disk"),
					resource.TestCheckResourceAttr(data.ResourceName, "instance_name", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval_seconds", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "performance_counters.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_linux_performance_object", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsDataSourceLinuxPerfObjDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsDataSourceLinuxPerfObjExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_requiresImport),
		},
	})
}

func testCheckAzureRMLogAnalyticsDataSourceLinuxPerfObjExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.DataSourcesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Log Analytics Data Source Linux performance objecy not found: %s", resourceName)
		}

		id, err := parse.LogAnalyticsDataSourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Workspace, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Log Analytics Data Source Linux performance object %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("failed to get on LogAnalytics.DataSources: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMLogAnalyticsDataSourceLinuxPerfObjDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.DataSourcesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_datasource_linux_performance_object" {
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

func testAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_basic(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_template(data)
	return fmt.Sprintf(`%s

resource "azurerm_log_analytics_datasource_linux_performance_object" "test" {
  name					= "acctestLADS-WE-%d"
  resource_group_name	= azurerm_resource_group.test.name
  workspace_name		= azurerm_log_analytics_workspace.test.name
  object_name			= "Logical Disk"
  instance_name			= "*"
  interval_seconds		= 10
  performance_counters	= ["%% Used Space"]
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_complete(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_template(data)
	return fmt.Sprintf(`%s

resource "azurerm_log_analytics_datasource_linux_performance_object" "test" {
  name					= "acctestLADS-WE-%d"
  resource_group_name	= azurerm_resource_group.test.name
  workspace_name		= azurerm_log_analytics_workspace.test.name
  object_name			= "Processor"
  instance_name			= "*"
  interval_seconds		= 10
  performance_counters	= ["%% Processor Time"]
}
`, template, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_basic(data)
	return fmt.Sprintf(`%s

resource "azurerm_log_analytics_datasource_linux_performance_object" "import" {
  name					= azurerm_log_analytics_datasource_linux_performance_object.test.name
  resource_group_name	= azurerm_log_analytics_datasource_linux_performance_object.test.resource_group_name
  workspace_name		= azurerm_log_analytics_datasource_linux_performance_object.test.workspace_name
  object_name			= azurerm_log_analytics_datasource_linux_performance_object.test.object_name
  instance_name			= azurerm_log_analytics_datasource_linux_performance_object.test.instance_name
  interval_seconds		= azurerm_log_analytics_datasource_linux_performance_object.test.interval_seconds
  performance_counters	= azurerm_log_analytics_datasource_linux_performance_object.test.performance_counters
}
`, template)
}

func testAccAzureRMLogAnalyticsDataSourceLinuxPerfObj_template(data acceptance.TestData) string {
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
