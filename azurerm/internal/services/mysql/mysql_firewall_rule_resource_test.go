package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MySQLFirewallRuleResource struct {
}

func TestAccMySQLFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_firewall_rule", "test")
	r := MySQLFirewallRuleResource{}

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

func TestAccMySQLFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_firewall_rule", "test")
	r := MySQLFirewallRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_mysql_firewall_rule"),
		},
	})
}

func (t MySQLFirewallRuleResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.FirewallRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MySQL.FirewallRulesClient.Get(ctx, id.ResourceGroup, id.SubscriptionId, id.Name)
	if err != nil {
		return nil, err
	}

	return utils.Bool(resp.ID != nil), nil
}

func (MySQLFirewallRuleResource) basic(data acceptance.TestData) string {
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

func (r MySQLFirewallRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_firewall_rule" "import" {
  name                = azurerm_mysql_firewall_rule.test.name
  resource_group_name = azurerm_mysql_firewall_rule.test.resource_group_name
  server_name         = azurerm_mysql_firewall_rule.test.server_name
  start_ip_address    = azurerm_mysql_firewall_rule.test.start_ip_address
  end_ip_address      = azurerm_mysql_firewall_rule.test.end_ip_address
}
`, r.basic(data))
}
