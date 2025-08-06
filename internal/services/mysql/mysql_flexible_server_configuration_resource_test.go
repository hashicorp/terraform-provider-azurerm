// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MySQLFlexibleServerConfigurationResource struct{}

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

func TestAccMySQLFlexibleServerConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server_configuration", "test")
	r := MySQLFlexibleServerConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.slowQueryLog(data, "OFF"),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClient(r.checkValue("OFF")),
			),
		},
		data.ImportStep(),
		{
			Config: r.slowQueryLog(data, "ON"),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClient(r.checkValue("ON")),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMySQLFlexibleServerConfiguration_multipleServerConfigurations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server_configuration", "test")
	r := MySQLFlexibleServerConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleServerConfigurations(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t MySQLFlexibleServerConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := configurations.ParseConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	ctx2, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	resp, err := clients.MySQL.FlexibleServers.Configurations.Get(ctx2, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %q: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r MySQLFlexibleServerConfigurationResource) checkReset(configurationName string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		serverId, err := servers.ParseFlexibleServerID(state.Attributes["id"])
		if err != nil {
			return err
		}

		configId := configurations.NewConfigurationID(serverId.SubscriptionId, serverId.ResourceGroupName, serverId.FlexibleServerName, configurationName)
		ctx2, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
		resp, err := clients.MySQL.FlexibleServers.Configurations.Get(ctx2, configId)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("%q does not exist", configId)
			}
			return fmt.Errorf("Bad: Get on mysqlConfigurationsClient: %+v", err)
		}

		actualValue := *resp.Model.Properties.Value
		defaultValue := *resp.Model.Properties.DefaultValue

		if defaultValue != actualValue {
			return fmt.Errorf("MySQL Flexible Server Configuration wasn't set to the default value. Expected '%s' - got '%s': \n%+v", defaultValue, actualValue, resp)
		}

		return nil
	}
}

func (r MySQLFlexibleServerConfigurationResource) checkValue(value string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := configurations.ParseConfigurationID(state.Attributes["id"])
		if err != nil {
			return err
		}

		ctx2, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
		resp, err := clients.MySQL.FlexibleServers.Configurations.Get(ctx2, *id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("%q does not exist", id.ConfigurationName)
			}

			return fmt.Errorf("Get on MySQL.FlexibleServerConfigurationsClient: %+v", err)
		}

		if *resp.Model.Properties.Value != value {
			return fmt.Errorf("MySQL Flexible Server Configuration wasn't set. Expected '%s' - got '%s': \n%+v", value, *resp.Model.Properties.Value, resp)
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

func (r MySQLFlexibleServerConfigurationResource) slowQueryLog(data acceptance.TestData, value string) string {
	return r.template(data, "slow_query_log", value)
}

func (r MySQLFlexibleServerConfigurationResource) multipleServerConfigurations(data acceptance.TestData) string {
	config := `
resource "azurerm_mysql_flexible_server_configuration" "test" {
  name                = "disconnect_on_expired_password"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mysql_flexible_server.test.name
  value               = "on"
}

resource "azurerm_mysql_flexible_server_configuration" "test2" {
  name                = "character_set_server"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mysql_flexible_server.test.name
  value               = "hebrew"
}

resource "azurerm_mysql_flexible_server_configuration" "test3" {
  name                = "interactive_timeout"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mysql_flexible_server.test.name
  value               = "30"
}

resource "azurerm_mysql_flexible_server_configuration" "test4" {
  name                = "log_slow_admin_statements"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mysql_flexible_server.test.name
  value               = "on"
}

resource "azurerm_mysql_flexible_server_configuration" "test5" {
  name                = "slow_query_log"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mysql_flexible_server.test.name
  value               = "on"
}
`
	return r.empty(data) + config
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
  sku_name               = "B_Standard_B1ms"
  zone                   = "1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
