// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PostgreSQLConfigurationResource struct{}

func TestAccPostgreSQLConfiguration_backslashQuote(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_configuration", "test")
	r := PostgreSQLConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backslashQuote(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClient(r.checkValue("on")),
			),
		},
		data.ImportStep(),
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				// "delete" resets back to the default value
				data.CheckWithClientForResource(r.checkReset("backslash_quote"), "azurerm_postgresql_server.test"),
			),
		},
	})
}

func TestAccPostgreSQLConfiguration_clientMinMessages(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_configuration", "test")
	r := PostgreSQLConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clientMinMessages(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClient(r.checkValue("DEBUG5")),
			),
		},
		data.ImportStep(),
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				// "delete" resets back to the default value
				data.CheckWithClientForResource(r.checkReset("client_min_messages"), "azurerm_postgresql_server.test"),
			),
		},
	})
}

func TestAccPostgreSQLConfiguration_deadlockTimeout(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_configuration", "test")
	r := PostgreSQLConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.deadlockTimeout(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClient(r.checkValue("5000")),
			),
		},
		data.ImportStep(),
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				// "delete" resets back to the default value
				data.CheckWithClientForResource(r.checkReset("deadlock_timeout"), "azurerm_postgresql_server.test"),
			),
		},
	})
}

func TestAccPostgreSQLConfiguration_multiplePostgreSQLConfigurations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_configuration", "test")
	r := PostgreSQLConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiplePostgreSQLConfigurations(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r PostgreSQLConfigurationResource) checkReset(configurationName string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := configurations.ParseServerID(state.Attributes["id"])
		if err != nil {
			return err
		}

		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, 15*time.Minute)
			defer cancel()
		}

		configurationId := configurations.NewConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ServerName, configurationName)

		resp, err := clients.Postgres.ConfigurationsClient.Get(ctx, configurationId)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("%s does not exist", id)
			}
			return fmt.Errorf("Bad: Get on postgresqlConfigurationsClient: %+v", err)
		}

		if resp.Model != nil && resp.Model.Properties != nil {
			actualValue := *resp.Model.Properties.Value
			defaultValue := *resp.Model.Properties.DefaultValue

			if defaultValue != actualValue {
				return fmt.Errorf("PostgreSQL Configuration wasn't set to the default value. Expected '%s' - got '%s': \n%+v", defaultValue, actualValue, resp)
			}
		}

		return nil
	}
}

func (r PostgreSQLConfigurationResource) checkValue(value string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := configurations.ParseConfigurationID(state.Attributes["id"])
		if err != nil {
			return err
		}

		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, 15*time.Minute)
			defer cancel()
		}

		resp, err := clients.Postgres.ConfigurationsClient.Get(ctx, *id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("%s does not exist", id)
			}

			return fmt.Errorf("Bad: Get on postgresqlConfigurationsClient: %+v", err)
		}

		if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.Value != nil {
			if *resp.Model.Properties.Value != value {
				return fmt.Errorf("PostgreSQL Configuration wasn't set. Expected '%s' - got '%s': \n%+v", value, *resp.Model.Properties.Value, resp)
			}
		} else {
			return fmt.Errorf("bad get on %s", id)
		}

		return nil
	}
}

func (t PostgreSQLConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := configurations.ParseConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 15*time.Minute)
		defer cancel()
	}

	resp, err := clients.Postgres.ConfigurationsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
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

  storage_mb                   = 51200
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement_enabled      = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r PostgreSQLConfigurationResource) multiplePostgreSQLConfigurations(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_configuration" "test" {
  name                = "idle_in_transaction_session_timeout"
  resource_group_name = azurerm_postgresql_server.test.resource_group_name
  server_name         = azurerm_postgresql_server.test.name
  value               = "60"
}

resource "azurerm_postgresql_configuration" "test2" {
  name                = "log_autovacuum_min_duration"
  resource_group_name = azurerm_postgresql_server.test.resource_group_name
  server_name         = azurerm_postgresql_server.test.name
  value               = "10"
}

resource "azurerm_postgresql_configuration" "test3" {
  name                = "log_lock_waits"
  resource_group_name = azurerm_postgresql_server.test.resource_group_name
  server_name         = azurerm_postgresql_server.test.name
  value               = "on"
}

resource "azurerm_postgresql_configuration" "test4" {
  name                = "log_min_duration_statement"
  resource_group_name = azurerm_postgresql_server.test.resource_group_name
  server_name         = azurerm_postgresql_server.test.name
  value               = "10"
}

resource "azurerm_postgresql_configuration" "test5" {
  name                = "log_statement"
  resource_group_name = azurerm_postgresql_server.test.resource_group_name
  server_name         = azurerm_postgresql_server.test.name
  value               = "ddl"
}

resource "azurerm_postgresql_configuration" "test6" {
  name                = "pg_stat_statements.track"
  resource_group_name = azurerm_postgresql_server.test.resource_group_name
  server_name         = azurerm_postgresql_server.test.name
  value               = "top"
}

resource "azurerm_postgresql_configuration" "test7" {
  name                = "pg_qs.query_capture_mode"
  resource_group_name = azurerm_postgresql_server.test.resource_group_name
  server_name         = azurerm_postgresql_server.test.name
  value               = "top"
}

resource "azurerm_postgresql_configuration" "test8" {
  name                = "pgms_wait_sampling.query_capture_mode"
  resource_group_name = azurerm_postgresql_server.test.resource_group_name
  server_name         = azurerm_postgresql_server.test.name
  value               = "all"
}

resource "azurerm_postgresql_configuration" "test9" {
  name                = "pg_qs.max_query_text_length"
  resource_group_name = azurerm_postgresql_server.test.resource_group_name
  server_name         = azurerm_postgresql_server.test.name
  value               = 10000
}

resource "azurerm_postgresql_configuration" "test10" {
  name                = "pg_qs.retention_period_in_days"
  resource_group_name = azurerm_postgresql_server.test.resource_group_name
  server_name         = azurerm_postgresql_server.test.name
  value               = 30
}
`, r.empty(data))
}
