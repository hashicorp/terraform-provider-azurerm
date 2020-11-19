package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPostgreSQLConfiguration_backslashQuote(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_configuration", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLConfiguration_backslashQuote(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLConfigurationValue(data.ResourceName, "on"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPostgreSQLConfiguration_empty(data),
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckAzureRMPostgreSQLConfigurationValueReset(data.RandomInteger, "backslash_quote"),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLConfiguration_clientMinMessages(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_configuration", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLConfiguration_clientMinMessages(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLConfigurationValue(data.ResourceName, "DEBUG5"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPostgreSQLConfiguration_empty(data),
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckAzureRMPostgreSQLConfigurationValueReset(data.RandomInteger, "client_min_messages"),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLConfiguration_deadlockTimeout(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_configuration", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLConfiguration_deadlockTimeout(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLConfigurationValue(data.ResourceName, "5000"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPostgreSQLConfiguration_empty(data),
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckAzureRMPostgreSQLConfigurationValueReset(data.RandomInteger, "deadlock_timeout"),
				),
			},
		},
	})
}

func testCheckAzureRMPostgreSQLConfigurationValue(resourceName string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.ConfigurationsClient
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
			return fmt.Errorf("Bad: no resource group found in state for PostgreSQL Configuration: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: PostgreSQL Configuration %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on postgresqlConfigurationsClient: %+v", err)
		}

		if *resp.Value != value {
			return fmt.Errorf("PostgreSQL Configuration wasn't set. Expected '%s' - got '%s': \n%+v", value, *resp.Value, resp)
		}

		return nil
	}
}

func testCheckAzureRMPostgreSQLConfigurationValueReset(rInt int, configurationName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.ConfigurationsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resourceGroup := fmt.Sprintf("acctestRG-psql-%d", rInt)
		serverName := fmt.Sprintf("acctest-psql-server-%d", rInt)

		resp, err := client.Get(ctx, resourceGroup, serverName, configurationName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: PostgreSQL Configuration %q (server %q resource group: %q) does not exist", configurationName, serverName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on postgresqlConfigurationsClient: %+v", err)
		}

		actualValue := *resp.Value
		defaultValue := *resp.DefaultValue

		if defaultValue != actualValue {
			return fmt.Errorf("PostgreSQL Configuration wasn't set to the default value. Expected '%s' - got '%s': \n%+v", defaultValue, actualValue, resp)
		}

		return nil
	}
}

func testCheckAzureRMPostgreSQLConfigurationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.ConfigurationsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_postgresql_configuration" {
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

func testAccAzureRMPostgreSQLConfiguration_backslashQuote(data acceptance.TestData) string {
	return testAccAzureRMPostgreSQLConfiguration_template(data, "backslash_quote", "on")
}

func testAccAzureRMPostgreSQLConfiguration_clientMinMessages(data acceptance.TestData) string {
	return testAccAzureRMPostgreSQLConfiguration_template(data, "client_min_messages", "DEBUG5")
}

func testAccAzureRMPostgreSQLConfiguration_deadlockTimeout(data acceptance.TestData) string {
	return testAccAzureRMPostgreSQLConfiguration_template(data, "deadlock_timeout", "5000")
}

func testAccAzureRMPostgreSQLConfiguration_template(data acceptance.TestData, name string, value string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_configuration" "test" {
  name                = "%s"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  value               = "%s"
}
`, testAccAzureRMPostgreSQLConfiguration_empty(data), name, value)
}

func testAccAzureRMPostgreSQLConfiguration_empty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
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
  version                      = "9.6"
  ssl_enforcement_enabled      = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
