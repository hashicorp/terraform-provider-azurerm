// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MySQLConfigurationResource struct{}

func TestAccMySQLConfiguration_characterSetServer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_configuration", "test")
	r := MySQLConfigurationResource{}

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
				data.CheckWithClientForResource(r.checkReset("character_set_server"), "azurerm_mysql_server.test"),
			),
		},
	})
}

func TestAccMySQLConfiguration_interactiveTimeout(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_configuration", "test")
	r := MySQLConfigurationResource{}

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
				data.CheckWithClientForResource(r.checkReset("interactive_timeout"), "azurerm_mysql_server.test"),
			),
		},
	})
}

func TestAccMySQLConfiguration_logSlowAdminStatements(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_configuration", "test")
	r := MySQLConfigurationResource{}

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
				data.CheckWithClientForResource(r.checkReset("log_slow_admin_statements"), "azurerm_mysql_server.test"),
			),
		},
	})
}

func (t MySQLConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := configurations.ParseConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MySQL.MySqlClient.Configurations.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading MySQL Configuration (%s): %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r MySQLConfigurationResource) checkReset(configurationName string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		serverId, err := configurations.ParseServerID(state.Attributes["id"])
		if err != nil {
			return err
		}

		configurationId := configurations.NewConfigurationID(serverId.SubscriptionId, serverId.ResourceGroupName, serverId.ServerName, configurationName)

		ctx2, cancel := context.WithTimeout(ctx, 15*time.Minute)
		defer cancel()
		resp, err := clients.MySQL.MySqlClient.Configurations.Get(ctx2, configurationId)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("%s does not exist", configurationId)
			}
			return fmt.Errorf("get on mysqlConfigurationsClient: %+v", err)
		}

		if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Value == nil || resp.Model.Properties.DefaultValue == nil {
			return fmt.Errorf("one of model/properties/value/defaultValue was nil for %s", configurationId)
		}
		actualValue := *resp.Model.Properties.Value
		defaultValue := *resp.Model.Properties.DefaultValue

		if defaultValue != actualValue {
			return fmt.Errorf("MySQL Configuration wasn't set to the default value. Expected '%s' - got '%s': \n%+v", defaultValue, actualValue, resp)
		}

		return nil
	}
}

func (r MySQLConfigurationResource) checkValue(value string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := configurations.ParseConfigurationID(state.Attributes["id"])
		if err != nil {
			return err
		}

		ctx2, cancel := context.WithTimeout(ctx, 15*time.Minute)
		defer cancel()
		resp, err := clients.MySQL.MySqlClient.Configurations.Get(ctx2, *id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("%s does not exist", id)
			}

			return fmt.Errorf("get on mysqlConfigurationsClient: %+v", err)
		}

		if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Value == nil {
			return fmt.Errorf("one of model/properties/value was nil for %s", id)
		}

		if *resp.Model.Properties.Value != value {
			return fmt.Errorf("MySQL Configuration wasn't set. Expected '%s' - got '%s': \n%+v", value, *resp.Model.Properties.Value, resp)
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

  storage_mb                   = 51200
  geo_redundant_backup_enabled = false
  backup_retention_days        = 7

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.7"
  ssl_enforcement_enabled      = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
