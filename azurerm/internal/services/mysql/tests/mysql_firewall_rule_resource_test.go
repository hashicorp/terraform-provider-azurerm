package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMySQLFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_firewall_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLFirewallRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLFirewallRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMySQLFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_firewall_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLFirewallRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLFirewallRuleExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMMySQLFirewallRule_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_mysql_firewall_rule"),
			},
		},
	})
}

func testCheckAzureRMMySQLFirewallRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MySQL.FirewallRulesClient
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
			return fmt.Errorf("Bad: no resource group found in state for MySQL Firewall Rule: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MySQL Firewall Rule %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on mysqlFirewallRulesClient: %s", err)
		}

		return nil
	}
}

func testCheckAzureRMMySQLFirewallRuleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MySQL.DatabasesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mysql_firewall_rule" {
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

			return fmt.Errorf("MySQL Firewall Rule still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMMySQLFirewallRule_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
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
  version                      = "5.6"
  ssl_enforcement_enabled      = true
}

resource "azurerm_mysql_firewall_rule" "test" {
  name                = "acctestfwrule-%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mysql_server.test.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "255.255.255.255"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMMySQLFirewallRule_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_firewall_rule" "import" {
  name                = azurerm_mysql_firewall_rule.test.name
  resource_group_name = azurerm_mysql_firewall_rule.test.resource_group_name
  server_name         = azurerm_mysql_firewall_rule.test.server_name
  start_ip_address    = azurerm_mysql_firewall_rule.test.start_ip_address
  end_ip_address      = azurerm_mysql_firewall_rule.test.end_ip_address
}
`, testAccAzureRMMySQLFirewallRule_basic(data))
}
