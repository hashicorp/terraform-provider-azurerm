package mysql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMySQLDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMySQLDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLDatabaseExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMMySQLDatabase_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_mysql_database"),
			},
		},
	})
}

func TestAccAzureRMMySQLDatabase_charsetUppercase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLDatabase_charsetUppercase(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "charset", "utf8"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMySQLDatabase_charsetMixedcase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLDatabase_charsetMixedcase(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "charset", "utf8"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMMySQLDatabaseExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MySQL.DatabasesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for MySQL Database: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MySQL Database %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on mysqlDatabasesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMMySQLDatabaseDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MySQL.DatabasesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mysql_database" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("MySQL Database still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMMySQLDatabase_basic(data acceptance.TestData) string {
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

func testAccAzureRMMySQLDatabase_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_database" "import" {
  name                = azurerm_mysql_database.test.name
  resource_group_name = azurerm_mysql_database.test.resource_group_name
  server_name         = azurerm_mysql_database.test.server_name
  charset             = azurerm_mysql_database.test.charset
  collation           = azurerm_mysql_database.test.collation
}
`, testAccAzureRMMySQLDatabase_basic(data))
}

func testAccAzureRMMySQLDatabase_charsetUppercase(data acceptance.TestData) string {
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

func testAccAzureRMMySQLDatabase_charsetMixedcase(data acceptance.TestData) string {
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
