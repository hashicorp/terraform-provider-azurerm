package mysql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMySQLConfiguration_characterSetServer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLConfiguration_characterSetServer(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLConfigurationValue(data.ResourceName, "hebrew"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMySQLConfiguration_empty(data),
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckAzureRMMySQLConfigurationValueReset(data, "character_set_server"),
				),
			},
		},
	})
}

func TestAccAzureRMMySQLConfiguration_interactiveTimeout(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLConfiguration_interactiveTimeout(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLConfigurationValue(data.ResourceName, "30"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMySQLConfiguration_empty(data),
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckAzureRMMySQLConfigurationValueReset(data, "interactive_timeout"),
				),
			},
		},
	})
}

func TestAccAzureRMMySQLConfiguration_logSlowAdminStatements(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLConfiguration_logSlowAdminStatements(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLConfigurationValue(data.ResourceName, "on"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMySQLConfiguration_empty(data),
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckAzureRMMySQLConfigurationValueReset(data, "log_slow_admin_statements"),
				),
			},
		},
	})
}

func testCheckAzureRMMySQLConfigurationValue(resourceName string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MySQL.ConfigurationsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for MySQL Configuration: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MySQL Configuration %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on mysqlConfigurationsClient: %+v", err)
		}

		if *resp.Value != value {
			return fmt.Errorf("MySQL Configuration wasn't set. Expected '%s' - got '%s': \n%+v", value, *resp.Value, resp)
		}

		return nil
	}
}

func testCheckAzureRMMySQLConfigurationValueReset(data acceptance.TestData, configurationName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MySQL.ConfigurationsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
		serverName := fmt.Sprintf("acctestmysqlsvr-%d", data.RandomInteger)

		resp, err := client.Get(ctx, resourceGroup, serverName, configurationName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MySQL Configuration %q (server %q resource group: %q) does not exist", configurationName, serverName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on mysqlConfigurationsClient: %+v", err)
		}

		actualValue := *resp.Value
		defaultValue := *resp.DefaultValue

		if defaultValue != actualValue {
			return fmt.Errorf("MySQL Configuration wasn't set to the default value. Expected '%s' - got '%s': \n%+v", defaultValue, actualValue, resp)
		}

		return nil
	}
}

func testCheckAzureRMMySQLConfigurationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MySQL.ConfigurationsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mysql_configuration" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}
	}

	return nil
}

func testAccAzureRMMySQLConfiguration_characterSetServer(data acceptance.TestData) string {
	return testAccAzureRMMySQLConfiguration_template(data, "character_set_server", "hebrew")
}

func testAccAzureRMMySQLConfiguration_interactiveTimeout(data acceptance.TestData) string {
	return testAccAzureRMMySQLConfiguration_template(data, "interactive_timeout", "30")
}

func testAccAzureRMMySQLConfiguration_logSlowAdminStatements(data acceptance.TestData) string {
	return testAccAzureRMMySQLConfiguration_template(data, "log_slow_admin_statements", "on")
}

func testAccAzureRMMySQLConfiguration_template(data acceptance.TestData, name string, value string) string {
	server := testAccAzureRMMySQLConfiguration_empty(data)
	config := fmt.Sprintf(`
resource "azurerm_mysql_configuration" "test" {
  name                = "%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mysql_server.test.name}"
  value               = "%s"
}
`, name, value)
	return server + config
}

func testAccAzureRMMySQLConfiguration_empty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.7"
  ssl_enforcement_enabled      = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
