package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MySQLConfigurationResource struct {
}

func TestAccMySQLConfiguration_characterSetServer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_configuration", "test")
	r := MySQLConfigurationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.characterSetServer(data),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClient(r.checkValue("hebrew")),
			),
		},
		data.ImportStep(),
		{
			Config: r.empty(data),
			Check: resource.ComposeTestCheckFunc(
				// "delete" resets back to the default value
				data.CheckWithClientForResource(r.checkReset("character_set_server"), "azurerm_mysql_server.test"),
			),
		},
	})
}

func TestAccMySQLConfiguration_interactiveTimeout(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_configuration", "test")
	r := MySQLConfigurationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.interactiveTimeout(data),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClient(r.checkValue("30")),
			),
		},
		data.ImportStep(),
		{
			Config: r.empty(data),
			Check: resource.ComposeTestCheckFunc(
				// "delete" resets back to the default value
				data.CheckWithClientForResource(r.checkReset("interactive_timeout"), "azurerm_mysql_server.test"),
			),
		},
	})
}

func TestAccMySQLConfiguration_logSlowAdminStatements(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_configuration", "test")
	r := MySQLConfigurationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.logSlowAdminStatements(data),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClient(r.checkValue("on")),
			),
		},
		data.ImportStep(),
		{
			Config: r.empty(data),
			Check: resource.ComposeTestCheckFunc(
				// "delete" resets back to the default value
				data.CheckWithClientForResource(r.checkReset("log_slow_admin_statements"), "azurerm_mysql_server.test"),
			),
		},
	})
}

func (t MySQLConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MySQL.ConfigurationsClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading MySQL Configuration (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r MySQLConfigurationResource) checkReset(configurationName string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
		id, err := parse.ServerID(state.Attributes["id"])
		if err != nil {
			return err
		}

		resp, err := clients.MySQL.ConfigurationsClient.Get(ctx, id.ResourceGroup, id.Name, configurationName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MySQL Configuration %q (server %q resource group: %q) does not exist", configurationName, id.Name, id.ResourceGroup)
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

func (r MySQLConfigurationResource) checkValue(value string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
		id, err := parse.ConfigurationID(state.Attributes["id"])
		if err != nil {
			return err
		}

		resp, err := clients.MySQL.ConfigurationsClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MySQL Configuration %q (server %q resource group: %q) does not exist", id.Name, id.ServerName, id.ResourceGroup)
			}

			return fmt.Errorf("Bad: Get on mysqlConfigurationsClient: %+v", err)
		}

		if *resp.Value != value {
			return fmt.Errorf("MySQL Configuration wasn't set. Expected '%s' - got '%s': \n%+v", value, *resp.Value, resp)
		}

		return nil
	}
}

func (r MySQLConfigurationResource) characterSetServer(data acceptance.TestData) string {
	return r.template(data, "character_set_server", "hebrew")
}

func (r MySQLConfigurationResource) interactiveTimeout(data acceptance.TestData) string {
	return r.template(data, "interactive_timeout", "30")
}

func (r MySQLConfigurationResource) logSlowAdminStatements(data acceptance.TestData) string {
	return r.template(data, "log_slow_admin_statements", "on")
}

func (r MySQLConfigurationResource) template(data acceptance.TestData, name string, value string) string {
	config := fmt.Sprintf(`
resource "azurerm_mysql_configuration" "test" {
  name                = "%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mysql_server.test.name}"
  value               = "%s"
}
`, name, value)
	return r.empty(data) + config
}

func (MySQLConfigurationResource) empty(data acceptance.TestData) string {
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
