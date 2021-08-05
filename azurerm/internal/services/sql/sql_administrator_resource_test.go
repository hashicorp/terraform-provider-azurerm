package sql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SqlAdministratorResource struct{}

func TestAccSqlAdministrator_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_active_directory_administrator", "test")
	r := SqlAdministratorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("login").HasValue("sqladmin"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withUpdates(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("login").HasValue("sqladmin2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSqlAdministrator_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_active_directory_administrator", "test")
	r := SqlAdministratorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("login").HasValue("sqladmin"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSqlAdministrator_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_active_directory_administrator", "test")
	r := SqlAdministratorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func (r SqlAdministratorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AzureActiveDirectoryAdministratorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Sql.ServerAzureADAdministratorsClient.Get(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving AAD Administrator %q (Server %q / Resource Group %q): %+v", id.AdministratorName, id.ServerName, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r SqlAdministratorResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AzureActiveDirectoryAdministratorID(state.ID)
	if err != nil {
		return nil, err
	}
	if _, err := client.Sql.ServerAzureADAdministratorsClient.Delete(ctx, id.ResourceGroup, id.ServerName); err != nil {
		return nil, err
	}
	return utils.Bool(true), nil
}

func (r SqlAdministratorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_active_directory_administrator" "test" {
  server_name         = azurerm_sql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
  login               = "sqladmin"
  tenant_id           = data.azurerm_client_config.current.tenant_id
  object_id           = data.azurerm_client_config.current.client_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r SqlAdministratorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sql_active_directory_administrator" "import" {
  server_name         = azurerm_sql_active_directory_administrator.test.server_name
  resource_group_name = azurerm_sql_active_directory_administrator.test.resource_group_name
  login               = azurerm_sql_active_directory_administrator.test.login
  tenant_id           = azurerm_sql_active_directory_administrator.test.tenant_id
  object_id           = azurerm_sql_active_directory_administrator.test.object_id
}
`, r.basic(data))
}

func (r SqlAdministratorResource) withUpdates(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_active_directory_administrator" "test" {
  server_name         = azurerm_sql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
  login               = "sqladmin2"
  tenant_id           = data.azurerm_client_config.current.tenant_id
  object_id           = data.azurerm_client_config.current.client_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
