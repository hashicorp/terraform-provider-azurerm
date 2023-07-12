// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/serveradministrators"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PostgreSqlAdministratorResource struct{}

func TestAccPostgreSqlAdministrator_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_active_directory_administrator", "test")
	r := PostgreSqlAdministratorResource{}

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
	})
}

func TestAccPostgreSqlAdministrator_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_active_directory_administrator", "test")
	r := PostgreSqlAdministratorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("login").HasValue("sqladmin"),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_postgresql_active_directory_administrator"),
		},
	})
}

func TestAccPostgreSqlAdministrator_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_active_directory_administrator", "test")
	r := PostgreSqlAdministratorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func (r PostgreSqlAdministratorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := serveradministrators.ParseServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.ServerAdministratorsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Postgresql AAD Administrator (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r PostgreSqlAdministratorResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := serveradministrators.ParseServerID(state.ID)
	if err != nil {
		return nil, err
	}

	if _, err := client.Postgres.ServerAdministratorsClient.Delete(ctx, *id); err != nil {
		return nil, fmt.Errorf("deleting Postgresql AAD Administrator (%s): %+v", id.String(), err)
	}

	return utils.Bool(true), nil
}

func (PostgreSqlAdministratorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
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

resource "azurerm_postgresql_active_directory_administrator" "test" {
  server_name         = azurerm_postgresql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
  login               = "sqladmin"
  tenant_id           = data.azurerm_client_config.current.tenant_id
  object_id           = data.azurerm_client_config.current.client_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r PostgreSqlAdministratorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_active_directory_administrator" "import" {
  server_name         = azurerm_postgresql_active_directory_administrator.test.server_name
  resource_group_name = azurerm_postgresql_active_directory_administrator.test.resource_group_name
  login               = azurerm_postgresql_active_directory_administrator.test.login
  tenant_id           = azurerm_postgresql_active_directory_administrator.test.tenant_id
  object_id           = azurerm_postgresql_active_directory_administrator.test.object_id
}
`, r.basic(data))
}

func (PostgreSqlAdministratorResource) withUpdates(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
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

resource "azurerm_postgresql_active_directory_administrator" "test" {
  server_name         = azurerm_postgresql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
  login               = "sqladmin2"
  tenant_id           = data.azurerm_client_config.current.tenant_id
  object_id           = data.azurerm_client_config.current.client_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
