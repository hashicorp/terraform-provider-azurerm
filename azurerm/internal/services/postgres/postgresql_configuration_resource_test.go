package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PostgreSQLConfigurationResource struct {
}

func TestAccPostgreSQLConfiguration_backslashQuote(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_configuration", "test")
	r := PostgreSQLConfigurationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.backslashQuote(data),
			Check: resource.ComposeTestCheckFunc(
				testCheckPostgreSQLConfigurationValue(data.ResourceName, "on"),
			),
		},
		data.ImportStep(),
		{
			Config: r.empty(data),
			Check: resource.ComposeTestCheckFunc(
				// "delete" resets back to the default value
				data.CheckWithClientForResource(r.checkReset("backslash_quote"), "azurerm_postgresql_server.test"),
			),
		},
	})
}

func TestAccPostgreSQLConfiguration_clientMinMessages(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_configuration", "test")
	r := PostgreSQLConfigurationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.clientMinMessages(data),
			Check: resource.ComposeTestCheckFunc(
				testCheckPostgreSQLConfigurationValue(data.ResourceName, "DEBUG5"),
			),
		},
		data.ImportStep(),
		{
			Config: r.empty(data),
			Check: resource.ComposeTestCheckFunc(
				// "delete" resets back to the default value
				data.CheckWithClientForResource(r.checkReset("client_min_messages"), "azurerm_postgresql_server.test"),
			),
		},
	})
}

func TestAccPostgreSQLConfiguration_deadlockTimeout(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_configuration", "test")
	r := PostgreSQLConfigurationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.deadlockTimeout(data),
			Check: resource.ComposeTestCheckFunc(
				testCheckPostgreSQLConfigurationValue(data.ResourceName, "5000"),
			),
		},
		data.ImportStep(),
		{
			Config: r.empty(data),
			Check: resource.ComposeTestCheckFunc(
				// "delete" resets back to the default value
				data.CheckWithClientForResource(r.checkReset("deadlock_timeout"), "azurerm_postgresql_server.test"),
			),
		},
	})
}

func testCheckPostgreSQLConfigurationValue(resourceName string, value string) resource.TestCheckFunc {
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

func (r PostgreSQLConfigurationResource) checkReset(configurationName string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
		id, err := parse.ServerID(state.Attributes["id"])
		if err != nil {
			return err
		}

		resp, err := clients.Postgres.ConfigurationsClient.Get(ctx, id.ResourceGroup, id.Name, configurationName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: PostgreSQL Configuration %q (server %q resource group: %q) does not exist", configurationName, id.Name, id.ResourceGroup)
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

func (t PostgreSQLConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.ConfigurationsClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Postgresql Configuration (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r PostgreSQLConfigurationResource) backslashQuote(data acceptance.TestData) string {
	return r.template(data, "backslash_quote", "on")
}

func (r PostgreSQLConfigurationResource) clientMinMessages(data acceptance.TestData) string {
	return r.template(data, "client_min_messages", "DEBUG5")
}

func (r PostgreSQLConfigurationResource) deadlockTimeout(data acceptance.TestData) string {
	return r.template(data, "deadlock_timeout", "5000")
}

func (r PostgreSQLConfigurationResource) template(data acceptance.TestData, name string, value string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_configuration" "test" {
  name                = "%s"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  value               = "%s"
}
`, r.empty(data), name, value)
}

func (PostgreSQLConfigurationResource) empty(data acceptance.TestData) string {
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
