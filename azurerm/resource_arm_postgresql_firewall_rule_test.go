package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMPostgreSQLFirewallRule_basic(t *testing.T) {
	resourceName := "azurerm_postgresql_firewall_rule.test"
	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLFirewallRule_basic(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLFirewallRuleExists(resourceName),
				),
			},
		},
	})
}

func testCheckAzureRMPostgreSQLFirewallRuleExists(name string) resource.TestCheckFunc {
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
			return fmt.Errorf("Bad: no resource group found in state for PostgreSQL Firewall Rule: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).postgresqlFirewallRulesClient

		resp, err := client.Get(resourceGroup, serverName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on postgresqlFirewallRulesClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: PostgreSQL Firewall Rule %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMPostgreSQLFirewallRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).postgresqlDatabasesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_postgresql_firewall_rule" {
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
			return fmt.Errorf("PostgreSQL Firewall Rule still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMPostgreSQLFirewallRule_basic(rInt int) string {
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

resource "azurerm_postgresql_firewall_rule" "test" {
  name                = "acctestfwrule-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_postgresql_server.test.name}"
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "255.255.255.255"
}
`, rInt, rInt, rInt)
}
