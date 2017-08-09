package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMSqlFirewallRule_basic(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccAzureRMSqlFirewallRule_basic, ri, ri, ri)
	postConfig := fmt.Sprintf(testAccAzureRMSqlFirewallRule_withUpdates, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlFirewallRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlFirewallRuleExists("azurerm_sql_firewall_rule.test"),
					resource.TestCheckResourceAttr("azurerm_sql_firewall_rule.test", "start_ip_address", "0.0.0.0"),
					resource.TestCheckResourceAttr("azurerm_sql_firewall_rule.test", "end_ip_address", "255.255.255.255"),
				),
			},

			resource.TestStep{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlFirewallRuleExists("azurerm_sql_firewall_rule.test"),
					resource.TestCheckResourceAttr("azurerm_sql_firewall_rule.test", "start_ip_address", "10.0.17.62"),
					resource.TestCheckResourceAttr("azurerm_sql_firewall_rule.test", "end_ip_address", "10.0.17.62"),
				),
			},
		},
	})
}

func testCheckAzureRMSqlFirewallRuleExists(name string) resource.TestCheckFunc {
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
			return fmt.Errorf("Bad: no resource group found in state forSQL Firewall Rule: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).sqlFirewallRulesClient
		resp, err := conn.Get(resourceGroup, serverName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get SQL Firewall Rule: %v", err)
		}

		if resp.Response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad:SQL Firewall Rule %s (resource group: %s) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMSqlFirewallRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).sqlFirewallRulesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sql_firewall_rule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, serverName, name)

		if err != nil {
			return nil
		}

		if resp.Response.StatusCode != http.StatusNotFound {
			return fmt.Errorf("SQL Firewall Rule still exists:\n%#v", resp.FirewallRuleProperties)
		}

	}

	return nil
}

var testAccAzureRMSqlFirewallRule_basic = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_sql_server" "test" {
    name = "acctestsqlserver%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "West US"
    version = "12.0"
    administrator_login = "mradministrator"
    administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_firewall_rule" "test" {
    name = "acctestsqlserver%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    server_name = "${azurerm_sql_server.test.name}"
    start_ip_address = "0.0.0.0"
    end_ip_address = "255.255.255.255"
}
`

var testAccAzureRMSqlFirewallRule_withUpdates = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_sql_server" "test" {
    name = "acctestsqlserver%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "West US"
    version = "12.0"
    administrator_login = "mradministrator"
    administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_firewall_rule" "test" {
    name = "acctestsqlserver%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    server_name = "${azurerm_sql_server.test.name}"
    start_ip_address = "10.0.17.62"
    end_ip_address = "10.0.17.62"
}
`
