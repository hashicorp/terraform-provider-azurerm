package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMPostgreSQLDatabase_basic(t *testing.T) {
	resourceName := "azurerm_postgresql_database.test"
	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLDatabase_basic(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLDatabaseExists(resourceName),
				),
			},
		},
	})
}

func testCheckAzureRMPostgreSQLDatabaseExists(name string) resource.TestCheckFunc {
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
			return fmt.Errorf("Bad: no resource group found in state for PostgreSQL Database: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).postgresqlDatabasesClient

		resp, err := client.Get(resourceGroup, serverName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on postgresqlDatabasesClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: PostgreSQL Database %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMPostgreSQLDatabaseDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).postgresqlDatabasesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_postgresql_database" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(resourceGroup, serverName, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("PostgreSQL Database still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMPostgreSQLDatabase_basic(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}
resource "azurerm_postgresql_server" "test" {
  name = "acctestpsqlsvr-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "PGSQLB50"
    capacity = 50
    tier = "Basic"
  }

  administrator_login = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version = "9.6"
  storage_mb = 51200
  ssl_enforcement = "Enabled"
}

resource "azurerm_postgresql_database" "test" {
  name                = "acctestdb_%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_postgresql_server.test.name}"
  charset             = "UTF8"
  collation           = "English_United States.1252"
}
`, rInt, rInt, rInt)
}
