package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jen20/riviera/sql"
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

func testCheckAzureRMSqlFirewallRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).rivieraClient

		readRequest := conn.NewRequestForURI(rs.Primary.ID)
		readRequest.Command = &sql.GetFirewallRule{}

		readResponse, err := readRequest.Execute()
		if err != nil {
			return fmt.Errorf("Bad: GetFirewallRule: %+v", err)
		}
		if !readResponse.IsSuccessful() {
			return fmt.Errorf("Bad: GetFirewallRule: %+v", readResponse.Error)
		}

		return nil
	}
}

func testCheckAzureRMSqlFirewallRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).rivieraClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sql_firewall_rule" {
			continue
		}

		readRequest := conn.NewRequestForURI(rs.Primary.ID)
		readRequest.Command = &sql.GetFirewallRule{}

		readResponse, err := readRequest.Execute()
		if err != nil {
			return fmt.Errorf("Bad: GetFirewallRule: %+v", err)
		}

		if readResponse.IsSuccessful() {
			return fmt.Errorf("Bad: SQL Server Firewall Rule still exists: %+v", readResponse.Error)
		}
	}

	return nil
}

func testAccAzureRMSqlFirewallRule_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
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
    name = "acctestRG_%d"
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
