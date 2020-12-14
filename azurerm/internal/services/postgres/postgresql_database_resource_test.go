package postgres_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPostgreSQLDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_database", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "charset", "UTF8"),
					resource.TestCheckResourceAttr(data.ResourceName, "collation", "English_United States.1252"),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_database", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "charset", "UTF8"),
					resource.TestCheckResourceAttr(data.ResourceName, "collation", "English_United States.1252"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMPostgreSQLDatabase_requiresImport),
		},
	})
}

func TestAccAzureRMPostgreSQLDatabase_collationWithHyphen(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_database", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLDatabase_collationWithHyphen(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "charset", "UTF8"),
					resource.TestCheckResourceAttr(data.ResourceName, "collation", "En-US"),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLDatabase_charsetLowercase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_database", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLDatabase_charsetLowercase(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "charset", "UTF8"),
					resource.TestCheckResourceAttr(data.ResourceName, "collation", "English_United States.1252"),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLDatabase_charsetMixedcase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_database", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLDatabase_charsetMixedcase(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "charset", "UTF8"),
					resource.TestCheckResourceAttr(data.ResourceName, "collation", "English_United States.1252"),
				),
			},
		},
	})
}

func testCheckAzureRMPostgreSQLDatabaseExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.DatabasesClient
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
			return fmt.Errorf("Bad: no resource group found in state for PostgreSQL Database: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: PostgreSQL Database %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on postgresqlDatabasesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPostgreSQLDatabaseDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.DatabasesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_postgresql_database" {
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

			return err
		}

		return fmt.Errorf("PostgreSQL Database still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMPostgreSQLDatabase_basic(data acceptance.TestData) string {
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

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement_enabled      = true
}

resource "azurerm_postgresql_database" "test" {
  name                = "acctest_PSQL_db_%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  charset             = "UTF8"
  collation           = "English_United States.1252"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPostgreSQLDatabase_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMPostgreSQLDatabase_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_database" "import" {
  name                = azurerm_postgresql_database.test.name
  resource_group_name = azurerm_postgresql_database.test.resource_group_name
  server_name         = azurerm_postgresql_database.test.server_name
  charset             = azurerm_postgresql_database.test.charset
  collation           = azurerm_postgresql_database.test.collation
}
`, template)
}

func testAccAzureRMPostgreSQLDatabase_collationWithHyphen(data acceptance.TestData) string {
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

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement_enabled      = true
}

resource "azurerm_postgresql_database" "test" {
  name                = "acctest_PSQL_db_%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  charset             = "UTF8"
  collation           = "En-US"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPostgreSQLDatabase_charsetLowercase(data acceptance.TestData) string {
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

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement_enabled      = true
}

resource "azurerm_postgresql_database" "test" {
  name                = "acctest_PSQL_db_%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  charset             = "utf8"
  collation           = "English_United States.1252"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPostgreSQLDatabase_charsetMixedcase(data acceptance.TestData) string {
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

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement_enabled      = true
}

resource "azurerm_postgresql_database" "test" {
  name                = "acctest_PSQL_db_%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  charset             = "Utf8"
  collation           = "English_United States.1252"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
