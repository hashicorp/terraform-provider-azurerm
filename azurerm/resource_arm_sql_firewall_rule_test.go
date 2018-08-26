package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSqlFirewallRule_basic(t *testing.T) {
	resourceName := "azurerm_sql_firewall_rule.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMSqlFirewallRule_basic(ri, testLocation())
	postConfig := testAccAzureRMSqlFirewallRule_withUpdates(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlFirewallRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "start_ip_address", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "end_ip_address", "255.255.255.255"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlFirewallRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "start_ip_address", "10.0.17.62"),
					resource.TestCheckResourceAttr(resourceName, "end_ip_address", "10.0.17.62"),
				),
			},
		},
	})
}

func TestAccAzureRMSqlFirewallRule_disappears(t *testing.T) {
	resourceName := "azurerm_sql_firewall_rule.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSqlFirewallRule_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlFirewallRuleExists(resourceName),
					testCheckAzureRMSqlFirewallRuleDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMSqlFirewallRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		ruleName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlFirewallRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, ruleName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("SQL Firewall Rule %q (server %q / resource group %q) was not found", ruleName, serverName, resourceGroup)
			}

			return err
		}

		return nil
	}
}

func testCheckAzureRMSqlFirewallRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sql_firewall_rule" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		ruleName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlFirewallRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, ruleName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("SQL Firewall Rule %q (server %q / resource group %q) still exists: %+v", ruleName, serverName, resourceGroup, resp)
	}

	return nil
}

func testCheckAzureRMSqlFirewallRuleDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		ruleName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlFirewallRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Delete(ctx, resourceGroup, serverName, ruleName)
		if err != nil {
			if utils.ResponseWasNotFound(resp) {
				return nil
			}

			return fmt.Errorf("Bad: Delete on sqlFirewallRulesClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMSqlFirewallRule_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_sql_server" "test" {
    name = "acctestsqlserver%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
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
`, rInt, location, rInt, rInt)
}

func testAccAzureRMSqlFirewallRule_withUpdates(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_sql_server" "test" {
    name = "acctestsqlserver%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
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
`, rInt, location, rInt, rInt)
}
