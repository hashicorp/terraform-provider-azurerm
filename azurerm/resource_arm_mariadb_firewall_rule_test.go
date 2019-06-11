package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMariaDbFirewallRule_basic(t *testing.T) {
	resourceName := "azurerm_mariadb_firewall_rule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbFirewallRule_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbFirewallRuleExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMariaDbFirewallRule_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_mariadb_firewall_rule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbFirewallRule_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbFirewallRuleExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMMariaDbFirewallRule_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_mariadb_firewall_rule"),
			},
		},
	})
}

func testCheckAzureRMMariaDbFirewallRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for MariaDb Firewall Rule: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).mariadbFirewallRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MariaDb Firewall Rule %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on mariadbFirewallRulesClient: %s", err)
		}

		return nil
	}
}

func testCheckAzureRMMariaDbFirewallRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).mariadbDatabasesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mariadb_firewall_rule" {
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

			return fmt.Errorf("MariaDb Firewall Rule still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMMariaDbFirewallRule_basic(rInt int, location string) string {
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
    name     = "GP_Gen5_2"
    capacity = 2
    tier     = "GeneralPurpose"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.6"
  ssl_enforcement              = "Enabled"
}

resource "azurerm_mariadb_firewall_rule" "test" {
  name                = "acctestfwrule-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mariadb_server.test.name}"
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "255.255.255.255"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMMariaDbFirewallRule_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mariadb_firewall_rule" "import" {
  name                = "${azurerm_mariadb_firewall_rule.test.name}"
  resource_group_name = "${azurerm_mariadb_firewall_rule.test.resource_group_name}"
  server_name         = "${azurerm_mariadb_firewall_rule.test.server_name}"
  start_ip_address    = "${azurerm_mariadb_firewall_rule.test.start_ip_address}"
  end_ip_address      = "${azurerm_mariadb_firewall_rule.test.end_ip_address}"
}
`, testAccAzureRMMariaDbFirewallRule_basic(rInt, location))
}
