// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MysqlFlexibleDatabaseResource struct{}

func TestAccMySQLFlexibleDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_database", "test")
	r := MysqlFlexibleDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMySQLFlexibleDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_database", "test")
	r := MysqlFlexibleDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_mysql_flexible_database"),
		},
	})
}

func TestAccMySQLFlexibleDatabase_charsetUppercase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_database", "test")
	r := MysqlFlexibleDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.charsetUppercase(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("charset").HasValue("utf8mb3_unicode_ci"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMySQLFlexibleDatabase_charsetMixedcase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_database", "test")
	r := MysqlFlexibleDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.charsetMixedcase(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("charset").HasValue("utf8mb3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMySQLFlexibleDatabase_collationUppercase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_database", "test")
	r := MysqlFlexibleDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.collationUppercase(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("collation").HasValue("utf8_unicode_ci"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMySQLFlexibleDatabase_utf8Aliases(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_database", "test")
	r := MysqlFlexibleDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.charsetAndCollation(data, "utf8mb3", "utf8mb3_unicode_ci"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:             r.charsetAndCollation(data, "utf8", "utf8_unicode_ci"),
			PlanOnly:           true,
			ExpectNonEmptyPlan: false,
		},
	})
}

func (r MysqlFlexibleDatabaseResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := databases.ParseDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MySQL.FlexibleServers.Databases.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r MysqlFlexibleDatabaseResource) basic(data acceptance.TestData) string {
	return r.charsetAndCollation(data, "utf8", "utf8_unicode_ci")
}

func (r MysqlFlexibleDatabaseResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_database" "import" {
  name                = azurerm_mysql_flexible_database.test.name
  resource_group_name = azurerm_mysql_flexible_database.test.resource_group_name
  server_name         = azurerm_mysql_flexible_database.test.server_name
  charset             = azurerm_mysql_flexible_database.test.charset
  collation           = azurerm_mysql_flexible_database.test.collation
}
`, r.basic(data))
}

func (r MysqlFlexibleDatabaseResource) charsetUppercase(data acceptance.TestData) string {
	return r.charsetAndCollation(data, "UTF8", "UTF8_UNICODE_CI")
}

func (r MysqlFlexibleDatabaseResource) charsetMixedcase(data acceptance.TestData) string {
	return r.charsetAndCollation(data, "Utf8", "utf8_unicode_ci")
}

func (r MysqlFlexibleDatabaseResource) collationUppercase(data acceptance.TestData) string {
	return r.charsetAndCollation(data, "utf8", "UTF8_UNICODE_CI")
}

func (MysqlFlexibleDatabaseResource) charsetAndCollation(data acceptance.TestData, charset, collation string) string {
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

resource "azurerm_mysql_flexible_database" "test" {
  name                = "acctestdb_%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mysql_flexible_server.test.name
  charset             = "%s"
  collation           = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, charset, collation)
}
