package mariadb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MariaDbFirewallRuleResource struct {
}

func TestAccMariaDbFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_firewall_rule", "test")
	r := MariaDbFirewallRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMariaDbFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_firewall_rule", "test")
	r := MariaDbFirewallRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_mariadb_firewall_rule"),
		},
	})
}

func (MariaDbFirewallRuleResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	serverName := id.Path["servers"]
	name := id.Path["firewallRules"]

	resp, err := clients.MariaDB.FirewallRulesClient.Get(ctx, id.ResourceGroup, serverName, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving MariaDB Firewall Rule %q (Server %q / Resource Group %q): %v", name, serverName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.FirewallRuleProperties != nil), nil
}

func (MariaDbFirewallRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "GP_Gen5_2"

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

resource "azurerm_mariadb_firewall_rule" "test" {
  name                = "acctestFWRule_01-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mariadb_server.test.name}"
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "255.255.255.255"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r MariaDbFirewallRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mariadb_firewall_rule" "import" {
  name                = azurerm_mariadb_firewall_rule.test.name
  resource_group_name = azurerm_mariadb_firewall_rule.test.resource_group_name
  server_name         = azurerm_mariadb_firewall_rule.test.server_name
  start_ip_address    = azurerm_mariadb_firewall_rule.test.start_ip_address
  end_ip_address      = azurerm_mariadb_firewall_rule.test.end_ip_address
}
`, r.basic(data))
}
