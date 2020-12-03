package mariadb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMariaDbDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMariaDbDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "charset", "utf8"),
					resource.TestCheckResourceAttr(data.ResourceName, "collation", "utf8_general_ci"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMariaDbDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMariaDbDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbDatabaseExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMMariaDbDatabase_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_mariadb_database"),
			},
		},
	})
}

func testCheckAzureRMMariaDbDatabaseExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MariaDB.DatabasesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("bad: no resource group found in state for MariaDB database: %q", name)
		}

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: MariaDB database %q (Server %q Resource Group: %q) does not exist", name, serverName, resourceGroup)
			}
			return fmt.Errorf("bad: get on mariadbDatabasesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMMariaDbDatabaseDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MariaDB.DatabasesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mariadb_database" {
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
			return fmt.Errorf("error MariaDB database %q (Resource Group %q) still exists:\n%+v", name, resourceGroup, err)
		}
		return fmt.Errorf("MariaDB database %q (Resource Group %q) still exists:\n%#+v", name, resourceGroup, resp)
	}

	return nil
}

func testAccAzureRMMariaDbDatabase_basic(data acceptance.TestData) string {
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

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

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

func testAccAzureRMMariaDbDatabase_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMariaDbDatabase_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mariadb_database" "import" {
  name                = azurerm_mariadb_database.test.name
  resource_group_name = azurerm_mariadb_database.test.resource_group_name
  server_name         = azurerm_mariadb_database.test.server_name
  charset             = azurerm_mariadb_database.test.charset
  collation           = azurerm_mariadb_database.test.collation
}
`, template)
}
