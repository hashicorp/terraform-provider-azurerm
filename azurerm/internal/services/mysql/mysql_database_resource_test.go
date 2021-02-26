package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MySQLDatabaseResource struct {
}

func TestAccMySQLDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_database", "test")
	r := MySQLDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMySQLDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_database", "test")
	r := MySQLDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_mysql_database"),
		},
	})
}

func TestAccMySQLDatabase_charsetUppercase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_database", "test")
	r := MySQLDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.charsetUppercase(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("charset").HasValue("utf8"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMySQLDatabase_charsetMixedcase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_database", "test")
	r := MySQLDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.charsetMixedcase(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("charset").HasValue("utf8"),
			),
		},
		data.ImportStep(),
	})
}

func (t MySQLDatabaseResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["databases"]

	resp, err := clients.MySQL.DatabasesClient.Get(ctx, resourceGroup, serverName, name)
	if err != nil {
		return nil, fmt.Errorf("reading MySQL Database (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (MySQLDatabaseResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestpsqlsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.6"
  ssl_enforcement_enabled      = true
}

resource "azurerm_mysql_database" "test" {
  name                = "acctestdb_%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mysql_server.test.name
  charset             = "utf8"
  collation           = "utf8_unicode_ci"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r MySQLDatabaseResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_database" "import" {
  name                = azurerm_mysql_database.test.name
  resource_group_name = azurerm_mysql_database.test.resource_group_name
  server_name         = azurerm_mysql_database.test.server_name
  charset             = azurerm_mysql_database.test.charset
  collation           = azurerm_mysql_database.test.collation
}
`, r.basic(data))
}

func (MySQLDatabaseResource) charsetUppercase(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestpsqlsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.6"
  ssl_enforcement_enabled      = true
}

resource "azurerm_mysql_database" "test" {
  name                = "acctestdb_%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mysql_server.test.name
  charset             = "UTF8"
  collation           = "utf8_unicode_ci"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (MySQLDatabaseResource) charsetMixedcase(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestpsqlsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.6"
  ssl_enforcement_enabled      = true
}

resource "azurerm_mysql_database" "test" {
  name                = "acctestdb_%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mysql_server.test.name
  charset             = "Utf8"
  collation           = "utf8_unicode_ci"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
