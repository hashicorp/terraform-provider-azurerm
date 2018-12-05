package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMariaDbDatabase_basic(t *testing.T) {
	resourceName := "azurerm_mariadb_database.test"
	ri := acctest.RandInt()
	config := testAccAzureRMMariaDbDatabase_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "charset", "UTF8"),
					resource.TestCheckResourceAttr(resourceName, "collation", "utf8_general_ci"),
				),
			},
		},
	})
}

func TestAccAzureRMMariaDbDatabase_charsetLowercase(t *testing.T) {
	resourceName := "azurerm_mariadb_database.test"
	ri := acctest.RandInt()
	config := testAccAzureRMMariaDbDatabase_charsetLowercase(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "charset", "UTF8"),
					resource.TestCheckResourceAttr(resourceName, "collation", "utf8_general_ci"),
				),
			},
		},
	})
}

func TestAccAzureRMMariaDbDatabase_charsetMixedcase(t *testing.T) {
	resourceName := "azurerm_mariadb_database.test"
	ri := acctest.RandInt()
	config := testAccAzureRMMariaDbDatabase_charsetMixedcase(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "charset", "UTF8"),
					resource.TestCheckResourceAttr(resourceName, "collation", "utf8_general_ci"),
				),
			},
		},
	})
}

func testCheckAzureRMMariaDbDatabaseExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for MariaDB Database: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).postgresqlDatabasesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MariaDB Database %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on postgresqlDatabasesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMMariaDbDatabaseDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).postgresqlDatabasesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

			return err
		}

		return fmt.Errorf("MariaDB Database still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMMariaDbDatabase_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "B_Gen5_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}

resource "azurerm_mariadb_database" "test" {
  name                = "acctestmariadb-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mariadb_server.test.name}"
  charset             = "UTF8"
  collation           = "utf8_general_ci"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMMariaDbDatabase_charsetLowercase(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "B_Gen5_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}

resource "azurerm_mariadb_database" "test" {
  name                = "acctestmariadb-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mariadb_server.test.name}"
  charset             = "utf8"
  collation           = "utf8_general_ci"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMMariaDbDatabase_charsetMixedcase(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "B_Gen5_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}

resource "azurerm_mariadb_database" "test" {
  name                = "acctestmariadb-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mariadb_server.test.name}"
  charset             = "Utf8"
  collation           = "utf8_general_ci"
}
`, rInt, location, rInt, rInt)
}
