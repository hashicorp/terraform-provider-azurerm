package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MySQLFlexibleServerConfigurationResource struct {
}

func TestAccMySQLFlexibleServerConfiguration_characterSetServer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server_configuration", "test")
	r := MySQLFlexibleServerConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.characterSetServer(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClient(r.checkValue("hebrew")),
			),
		},
		data.ImportStep(),
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				// "delete" resets back to the default value
				data.CheckWithClientForResource(r.checkReset("character_set_server"), "azurerm_mysql_flexible_server.test"),
			),
		},
	})
}

func TestAccMySQLFlexibleServerConfiguration_interactiveTimeout(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server_configuration", "test")
	r := MySQLFlexibleServerConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.interactiveTimeout(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClient(r.checkValue("30")),
			),
		},
		data.ImportStep(),
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				// "delete" resets back to the default value
				data.CheckWithClientForResource(r.checkReset("interactive_timeout"), "azurerm_mysql_flexible_server.test"),
			),
		},
	})
}

func TestAccMySQLFlexibleServerConfiguration_logSlowAdminStatements(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server_configuration", "test")
	r := MySQLFlexibleServerConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logSlowAdminStatements(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClient(r.checkValue("on")),
			),
		},
		data.ImportStep(),
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				// "delete" resets back to the default value
				data.CheckWithClientForResource(r.checkReset("log_slow_admin_statements"), "azurerm_mysql_flexible_server.test"),
			),
		},
	})
}

func (t MySQLFlexibleServerConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FlexibleServerConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MySQL.FlexibleServerConfigurationsClient.Get(ctx, id.ResourceGroup, id.FlexibleServerName, id.ConfigurationName)
	if err != nil {
		return nil, fmt.Errorf("reading %q: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r MySQLFlexibleServerConfigurationResource) checkReset(configurationName string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := parse.FlexibleServerID(state.Attributes["id"])
		if err != nil {
			return err
		}

		resp, err := clients.MySQL.FlexibleServerConfigurationsClient.Get(ctx, id.ResourceGroup, id.Name, configurationName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("%q does not exist", id)
			}
			return fmt.Errorf("Bad: Get on mysqlConfigurationsClient: %+v", err)
		}

		actualValue := *resp.Value
		defaultValue := *resp.DefaultValue

		if defaultValue != actualValue {
			return fmt.Errorf("MySQL Flexible Server Configuration wasn't set to the default value. Expected '%s' - got '%s': \n%+v", defaultValue, actualValue, resp)
		}

		return nil
	}
}

func (r MySQLFlexibleServerConfigurationResource) checkValue(value string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := parse.FlexibleServerConfigurationID(state.Attributes["id"])
		if err != nil {
			return err
		}

		resp, err := clients.MySQL.FlexibleServerConfigurationsClient.Get(ctx, id.ResourceGroup, id.FlexibleServerName, id.ConfigurationName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("%q does not exist", id.ConfigurationName)
			}

			return fmt.Errorf("Get on MySQL.FlexibleServerConfigurationsClient: %+v", err)
		}

		if *resp.Value != value {
			return fmt.Errorf("MySQL Flexible Server Configuration wasn't set. Expected '%s' - got '%s': \n%+v", value, *resp.Value, resp)
		}

		return nil
	}
}

func (r MySQLFlexibleServerConfigurationResource) characterSetServer(data acceptance.TestData) string {
	return r.template(data, "character_set_server", "hebrew")
}

func (r MySQLFlexibleServerConfigurationResource) interactiveTimeout(data acceptance.TestData) string {
	return r.template(data, "interactive_timeout", "30")
}

func (r MySQLFlexibleServerConfigurationResource) logSlowAdminStatements(data acceptance.TestData) string {
	return r.template(data, "log_slow_admin_statements", "on")
}

func (r MySQLFlexibleServerConfigurationResource) template(data acceptance.TestData, name string, value string) string {
	config := fmt.Sprintf(`
resource "azurerm_mysql_flexible_server_configuration" "test" {
  name                = "%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mysql_flexible_server.test.name}"
  value               = "%s"
}
`, name, value)
	return r.empty(data) + config
}

func (MySQLFlexibleServerConfigurationResource) empty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  sku_name               = "B_Standard_B1s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
