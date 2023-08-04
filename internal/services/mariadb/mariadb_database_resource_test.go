// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mariadb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MariaDbDatabaseResource struct{}

func TestAccMariaDbDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_database", "test")
	r := MariaDbDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("charset").HasValue("utf8"),
				check.That(data.ResourceName).Key("collation").HasValue("utf8_general_ci"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMariaDbDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_database", "test")
	r := MariaDbDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_mariadb_database"),
		},
	})
}

func (MariaDbDatabaseResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := databases.ParseDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MariaDB.DatabasesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (MariaDbDatabaseResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = %q
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "B_Gen5_2"

  storage_mb                   = 51200
  geo_redundant_backup_enabled = false
  backup_retention_days        = 7

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement_enabled      = true
}

resource "azurerm_mariadb_database" "test" {
  name                = "acctestmariadb_%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mariadb_server.test.name
  charset             = "utf8"
  collation           = "utf8_general_ci"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r MariaDbDatabaseResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mariadb_database" "import" {
  name                = azurerm_mariadb_database.test.name
  resource_group_name = azurerm_mariadb_database.test.resource_group_name
  server_name         = azurerm_mariadb_database.test.server_name
  charset             = azurerm_mariadb_database.test.charset
  collation           = azurerm_mariadb_database.test.collation
}
`, r.basic(data))
}
