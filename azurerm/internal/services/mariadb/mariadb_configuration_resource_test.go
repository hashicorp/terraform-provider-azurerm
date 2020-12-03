package mariadb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccMariaDbConfiguration_characterSetServer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMariaDbConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariaDbConfiguration_characterSetServer(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbConfigurationValue(data.ResourceName, "hebrew"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccMariaDbConfiguration_empty(data),
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckMariaDbConfigurationValueReset(data.RandomInteger, "character_set_server"),
				),
			},
		},
	})
}

func TestAccMariaDbConfiguration_interactiveTimeout(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMariaDbConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariaDbConfiguration_interactiveTimeout(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbConfigurationValue(data.ResourceName, "30"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccMariaDbConfiguration_empty(data),
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckMariaDbConfigurationValueReset(data.RandomInteger, "interactive_timeout"),
				),
			},
		},
	})
}

func TestAccMariaDbConfiguration_logSlowAdminStatements(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMariaDbConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariaDbConfiguration_logSlowAdminStatements(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbConfigurationValue(data.ResourceName, "on"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccMariaDbConfiguration_empty(data),
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckMariaDbConfigurationValueReset(data.RandomInteger, "log_slow_admin_statements"),
				),
			},
		},
	})
}

func testCheckMariaDbConfigurationValue(resourceName string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MariaDB.ConfigurationsClient
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
			return fmt.Errorf("Bad: no resource group found in state for MariaDb Configuration: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MariaDb Configuration %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on mariadbConfigurationsClient: %+v", err)
		}

		if *resp.Value != value {
			return fmt.Errorf("MariaDb Configuration wasn't set. Expected '%s' - got '%s': \n%+v", value, *resp.Value, resp)
		}

		return nil
	}
}

func testCheckMariaDbConfigurationValueReset(rInt int, configurationName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MariaDB.ConfigurationsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resourceGroup := fmt.Sprintf("acctestRG-%d", rInt)
		serverName := fmt.Sprintf("acctestmariadbsvr-%d", rInt)

		resp, err := client.Get(ctx, resourceGroup, serverName, configurationName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MariaDb Configuration %q (server %q resource group: %q) does not exist", configurationName, serverName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on mariadbConfigurationsClient: %+v", err)
		}

		actualValue := *resp.Value
		defaultValue := *resp.DefaultValue

		if defaultValue != actualValue {
			return fmt.Errorf("MariaDb Configuration wasn't set to the default value. Expected '%s' - got '%s': \n%+v", defaultValue, actualValue, resp)
		}

		return nil
	}
}

func testCheckMariaDbConfigurationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MariaDB.ConfigurationsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mariadb_configuration" {
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

func testAccMariaDbConfiguration_characterSetServer(data acceptance.TestData) string {
	return testAccMariaDbConfiguration_template(data, "character_set_server", "hebrew")
}

func testAccMariaDbConfiguration_interactiveTimeout(data acceptance.TestData) string {
	return testAccMariaDbConfiguration_template(data, "interactive_timeout", "30")
}

func testAccMariaDbConfiguration_logSlowAdminStatements(data acceptance.TestData) string {
	return testAccMariaDbConfiguration_template(data, "log_slow_admin_statements", "on")
}

func testAccMariaDbConfiguration_template(data acceptance.TestData, name string, value string) string {
	server := testAccMariaDbConfiguration_empty(data)
	config := fmt.Sprintf(`
resource "azurerm_mariadb_configuration" "test" {
  name                = "%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mariadb_server.test.name}"
  value               = "%s"
}
`, name, value)
	return server + config
}

func testAccMariaDbConfiguration_empty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku_name            = "GP_Gen5_2"
  version             = "10.2"

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  ssl_enforcement_enabled      = true
  storage_mb                   = 51200
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
