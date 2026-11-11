// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PostgresqlFlexibleServerConfigurationResource struct{}

func TestAccFlexibleServerConfiguration_backslashQuote(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_configuration", "test")
	r := PostgresqlFlexibleServerConfigurationResource{}
	name := "backslash_quote"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, name, "on"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("value").HasValue("on"),
			),
		},
		data.ImportStep(),
		{
			Config: r.template(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.checkWith(name, resetToDefaultCheck), "azurerm_postgresql_flexible_server.test"),
			),
		},
	})
}

func TestAccFlexibleServerConfiguration_azureExtensions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_configuration", "test")
	r := PostgresqlFlexibleServerConfigurationResource{}
	name := "azure.extensions"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, name, "CUBE,CITEXT,BTREE_GIST"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("value").HasValue("CUBE,CITEXT,BTREE_GIST"),
			),
		},
		data.ImportStep(),
		{
			Config: r.template(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.checkWith(name, resetToDefaultCheck), "azurerm_postgresql_flexible_server.test"),
			),
		},
	})
}

func TestAccFlexibleServerConfiguration_pgbouncerEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_configuration", "test")
	r := PostgresqlFlexibleServerConfigurationResource{}
	name := "pgbouncer.enabled"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, name, "true"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("value").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.template(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.checkWith(name, resetToDefaultCheck), "azurerm_postgresql_flexible_server.test"),
			),
		},
	})
}

func TestAccFlexibleServerConfiguration_updateApplicationName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_configuration", "test")
	r := PostgresqlFlexibleServerConfigurationResource{}
	name := "application_name"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, name, "true"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("value").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, name, "false"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("value").HasValue("false"),
			),
		},
	})
}

func (r PostgresqlFlexibleServerConfigurationResource) checkWith(configurationName string, checker func(*configurations.ConfigurationProperties) error) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := configurations.ParseFlexibleServerID(state.ID)
		if err != nil {
			return err
		}

		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, 15*time.Minute)
			defer cancel()
		}

		configurationId := configurations.NewConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName, configurationName)

		resp, err := clients.Postgres.FlexibleServersConfigurationsClient.Get(ctx, configurationId)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("%s does not exist", id)
			}
			return fmt.Errorf("retrieving %s: %+v", configurationId, err)
		}

		if resp.Model == nil || resp.Model.Properties == nil {
			return fmt.Errorf("unable to retrieve configuration properties for %s", configurationId)
		}

		if err = checker(resp.Model.Properties); err != nil {
			return fmt.Errorf("%s: \n%+v", err, resp)
		}

		return nil
	}
}

func resetToDefaultCheck(props *configurations.ConfigurationProperties) error {
	if props.Value != nil && props.DefaultValue != nil {
		actualValue := *props.Value
		defaultValue := *props.DefaultValue

		if defaultValue != actualValue {
			return fmt.Errorf("Azure Postgresql Flexible Server configuration wasn't set to the default value. Expected '%s' - got '%s'", defaultValue, actualValue)
		}
	}
	return nil
}

func pendingRestartCheck(expectedValue bool) func(*configurations.ConfigurationProperties) error {
	return func(props *configurations.ConfigurationProperties) error {
		if props.Value != nil && props.IsConfigPendingRestart != nil {
			actualValue := *props.IsConfigPendingRestart

			if actualValue != expectedValue {
				return fmt.Errorf("Azure Postgresql Flexible Server configuration IsConfigPendingRestart is expected'%v', but got ='%v'", expectedValue, actualValue)
			}
		}
		return nil
	}
}

func configurationValueCheck(expectedValue string) func(*configurations.ConfigurationProperties) error {
	return func(props *configurations.ConfigurationProperties) error {
		if v := pointer.From(props.Value); v != expectedValue {
			return fmt.Errorf("configuration value does not match: expected `%s`, got `%s`", expectedValue, v)
		}
		return nil
	}
}

func TestAccFlexibleServerConfiguration_multiplePostgresqlFlexibleServerConfigurations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_configuration", "test")
	r := PostgresqlFlexibleServerConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiplePostgresqlFlexibleServerConfigurations(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFlexibleServerConfiguration_restartServerForStaticParameters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_configuration", "test")
	r := PostgresqlFlexibleServerConfigurationResource{}
	name := "cron.max_running_jobs"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, name, "5"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("value").HasValue("5"),
				// by default "static" parameter change triggers server restart
				data.CheckWithClientForResource(r.checkWith(name, pendingRestartCheck(false)), "azurerm_postgresql_flexible_server.test"),
			),
		},
		data.ImportStep(),
		{
			Config: r.template(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.checkWith(name, resetToDefaultCheck), "azurerm_postgresql_flexible_server.test"),
			),
		},
	})
}

func TestAccFlexibleServerConfiguration_staticParameterWithReplica(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_configuration", "test")
	r := PostgresqlFlexibleServerConfigurationResource{}
	name := "cron.max_running_jobs"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.staticParameterWithReplica(data, name, "5"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("value").HasValue("5"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFlexibleServerConfiguration_scaleParametersWithReplica(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_configuration", "replica_max_worker_processes")
	r := PostgresqlFlexibleServerConfigurationResource{}

	locksResource := "azurerm_postgresql_flexible_server_configuration.replica_max_locks_per_transaction"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scaleParametersWithReplica(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(locksResource).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data.ImportStepFor(locksResource),
		{
			Config: r.scaleParametersWithReplica(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				// removing the replica overrides should reset them to the primary config
				data.CheckWithClientForResource(r.checkWith("max_worker_processes", configurationValueCheck("32")), "azurerm_postgresql_flexible_server.replica"),
				data.CheckWithClientForResource(r.checkWith("max_locks_per_transaction", configurationValueCheck("128")), "azurerm_postgresql_flexible_server.replica"),
			),
		},
	})
}

func TestAccFlexibleServerConfiguration_doesNotRestartServer_whenFeatureIsDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_configuration", "test")
	r := PostgresqlFlexibleServerConfigurationResource{}
	name := "cron.max_running_jobs"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withDisabledServerRestarts(data, name, "5"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("value").HasValue("5"),
				// when the feature is disabled, "static" parameter change does not trigger server restart and config stays in "pending-restart" state
				data.CheckWithClientForResource(r.checkWith(name, pendingRestartCheck(true)), "azurerm_postgresql_flexible_server.test"),
			),
		},
		data.ImportStep(),
		{
			Config: r.template(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.checkWith(name, resetToDefaultCheck), "azurerm_postgresql_flexible_server.test"),
			),
		},
	})
}

// Helper functions for verification
func (r PostgresqlFlexibleServerConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := configurations.ParseConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 15*time.Minute)
		defer cancel()
	}

	resp, err := clients.Postgres.FlexibleServersConfigurationsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Postgresql Configuration (%s): %+v", id.String(), err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (PostgresqlFlexibleServerConfigurationResource) template(data acceptance.TestData) string {
	return PostgresqlFlexibleServerResource{}.basic(data)
}

func (r PostgresqlFlexibleServerConfigurationResource) basic(data acceptance.TestData, name, value string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_postgresql_flexible_server_configuration" "test" {
  name      = "%[2]s"
  server_id = azurerm_postgresql_flexible_server.test.id
  value     = "%[3]s"
}
`, r.template(data), name, value)
}

func (r PostgresqlFlexibleServerConfigurationResource) staticParameterWithReplica(data acceptance.TestData, name, value string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_postgresql_flexible_server_configuration" "test" {
  name      = "%[2]s"
  server_id = azurerm_postgresql_flexible_server.test.id
  value     = "%[3]s"
}

resource "azurerm_postgresql_flexible_server" "replica" {
  name                = "acctest-replica-%[4]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  zone                = azurerm_postgresql_flexible_server.test.zone
  create_mode         = "Replica"
  source_server_id    = azurerm_postgresql_flexible_server.test.id
}
`, r.template(data), name, value, data.RandomInteger)
}

func (r PostgresqlFlexibleServerConfigurationResource) scaleParametersWithReplica(data acceptance.TestData, overrides bool) string {
	configurationOverrides := ""
	if overrides {
		configurationOverrides = `
resource "azurerm_postgresql_flexible_server_configuration" "replica_max_worker_processes" {
  name      = "max_worker_processes"
  server_id = azurerm_postgresql_flexible_server.replica.id
  value     = "48"
}

resource "azurerm_postgresql_flexible_server_configuration" "replica_max_locks_per_transaction" {
  name      = "max_locks_per_transaction"
  server_id = azurerm_postgresql_flexible_server.replica.id
  value     = "192"
}
`
	}

	return fmt.Sprintf(`
%[1]s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%[2]d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  version                = "17"
  sku_name               = "GP_Standard_D4s_v3"
  storage_mb             = 32768
  zone                   = "2"
}

resource "azurerm_postgresql_flexible_server_configuration" "primary_max_worker_processes" {
  name      = "max_worker_processes"
  server_id = azurerm_postgresql_flexible_server.test.id
  value     = "32"
}

resource "azurerm_postgresql_flexible_server_configuration" "primary_max_locks_per_transaction" {
  name      = "max_locks_per_transaction"
  server_id = azurerm_postgresql_flexible_server.test.id
  value     = "128"
}

resource "azurerm_postgresql_flexible_server" "replica" {
  name                = "acctest-fs-replica-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  create_mode         = "Replica"
  source_server_id    = azurerm_postgresql_flexible_server.test.id
  zone                = azurerm_postgresql_flexible_server.test.zone
  version             = "17"
  sku_name            = "GP_Standard_D4s_v3"
  storage_mb          = 32768
}

%[3]s
`, PostgresqlFlexibleServerResource{}.template(data), data.RandomInteger, configurationOverrides)
}

func (PostgresqlFlexibleServerConfigurationResource) multiplePostgresqlFlexibleServerConfigurations(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_postgresql_flexible_server_configuration" "test" {
  name      = "idle_in_transaction_session_timeout"
  server_id = azurerm_postgresql_flexible_server.test.id
  value     = "60"
}

resource "azurerm_postgresql_flexible_server_configuration" "test2" {
  name      = "log_autovacuum_min_duration"
  server_id = azurerm_postgresql_flexible_server.test.id
  value     = "10"
}

resource "azurerm_postgresql_flexible_server_configuration" "test3" {
  name      = "log_lock_waits"
  server_id = azurerm_postgresql_flexible_server.test.id
  value     = "on"
}

resource "azurerm_postgresql_flexible_server_configuration" "test4" {
  name      = "log_min_duration_statement"
  server_id = azurerm_postgresql_flexible_server.test.id
  value     = "10"
}

resource "azurerm_postgresql_flexible_server_configuration" "test5" {
  name      = "log_statement"
  server_id = azurerm_postgresql_flexible_server.test.id
  value     = "ddl"
}
`, PostgresqlFlexibleServerResource{}.complete(data))
}

func (r PostgresqlFlexibleServerConfigurationResource) withDisabledServerRestarts(data acceptance.TestData, name, value string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    postgresql_flexible_server {
      restart_server_on_configuration_value_change = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-postgresql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%[1]d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage_mb             = 32768
  version                = "16"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"
}

resource "azurerm_postgresql_flexible_server_configuration" "test" {
  name      = "%[3]s"
  server_id = azurerm_postgresql_flexible_server.test.id
  value     = "%[4]s"
}
`, data.RandomInteger, data.Locations.Primary, name, value)
}
