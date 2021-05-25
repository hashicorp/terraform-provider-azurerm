package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PostgresqlFlexibleServerConfigurationResource struct {
}

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
				data.CheckWithClientForResource(r.checkReset(name), "azurerm_postgresql_flexible_server.test"),
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
				data.CheckWithClientForResource(r.checkReset(name), "azurerm_postgresql_flexible_server.test"),
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
			Config: r.basic(data, name, "Test APP before"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("value").HasValue("Test APP before"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, name, "Test APP after"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("value").HasValue("Test APP after"),
			),
		},
	})
}

func (r PostgresqlFlexibleServerConfigurationResource) checkReset(configurationName string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := parse.FlexibleServerID(state.ID)
		if err != nil {
			return err
		}

		resp, err := clients.Postgres.FlexibleServersConfigurationsClient.Get(ctx, id.ResourceGroup, id.Name, configurationName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Azure Postgresql Flexible Server configuration %q (server %q resource group: %q) does not exist", configurationName, id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on postgresqlConfigurationsClient: %+v", err)
		}

		actualValue := *resp.Value
		defaultValue := *resp.DefaultValue

		if defaultValue != actualValue {
			return fmt.Errorf("Azure Postgresql Flexible Server configuration wasn't set to the default value. Expected '%s' - got '%s': \n%+v", defaultValue, actualValue, resp)
		}

		return nil
	}
}

// Helper functions for verification
func (r PostgresqlFlexibleServerConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FlexibleServerConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.FlexibleServersConfigurationsClient.Get(ctx, id.ResourceGroup, id.FlexibleServerName, state.Attributes["name"])
	if err != nil {
		return nil, fmt.Errorf("reading Postgresql Configuration (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (PostgresqlFlexibleServerConfigurationResource) template(data acceptance.TestData) string {
	return PostgresqlFlexibleServerResource{}.basic(data)
}

func (r PostgresqlFlexibleServerConfigurationResource) basic(data acceptance.TestData, name, value string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server_configuration" "test" {
  name      = "%s"
  server_id = azurerm_postgresql_flexible_server.test.id
  value     = "%s"
}
`, r.template(data), name, value)
}
