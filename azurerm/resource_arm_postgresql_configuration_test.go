package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMPostgreSQLConfiguration_backslashQuote(t *testing.T) {
	resourceName := "azurerm_postgresql_configuration.test"
	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLConfiguration_backslashQuote(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			return testCheckAzureRMPostgreSQLConfigurationReset(s, "safe_encoding")
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLConfigurationValue(resourceName, "on"),
				),
			},
		},
	})
}

func testCheckAzureRMPostgreSQLConfigurationValue(name string, value string) resource.TestCheckFunc {
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
			return fmt.Errorf("Bad: no resource group found in state for PostgreSQL Configuration: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).postgresqlConfigurationsClient

		resp, err := client.Get(resourceGroup, serverName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on postgresqlConfigurationsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: PostgreSQL Configuration %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
		}

		if *resp.Value != value {
			return fmt.Errorf("PostgreSQL Configuration wasn't set. Expected '%s' - got '%s': \n%+v", value, *resp.Value, resp)
		}

		return nil
	}
}

func testCheckAzureRMPostgreSQLConfigurationReset(s *terraform.State, value string) error {
	client := testAccProvider.Meta().(*ArmClient).postgresqlConfigurationsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_postgresql_configuration" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(resourceGroup, serverName, name)

		if err != nil {
			return nil
		}

		if *resp.Value != value {
			return fmt.Errorf("PostgreSQL Configuration wasn't reset. Expected '%s' - got '%s': \n%+v", value, *resp.Value, resp)
		}
	}

	return nil
}

func testAccAzureRMPostgreSQLConfiguration_backslashQuote(rInt int) string {
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

resource "azurerm_postgresql_configuration" "test" {
  name                = "backslash_quote"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_postgresql_server.test.name}"
  value               = "on"
}
`, rInt, rInt)
}
