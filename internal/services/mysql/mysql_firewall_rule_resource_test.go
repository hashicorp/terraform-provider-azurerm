// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MySQLFirewallRuleResource struct{}

func TestAccMySQLFirewallRule_basic(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skipping since `azurerm_mysql_firewall_rule` is deprecated and will be removed in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_mysql_firewall_rule", "test")
	r := MySQLFirewallRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMySQLFirewallRule_requiresImport(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skipping since `azurerm_mysql_firewall_rule` is deprecated and will be removed in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_mysql_firewall_rule", "test")
	r := MySQLFirewallRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t MySQLFirewallRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := firewallrules.ParseFirewallRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MySQL.MySqlClient.FirewallRules.Get(ctx, *id)
	if err != nil {
		return nil, err
	}

	return utils.Bool(resp.Model != nil), nil
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

  storage_mb                   = 51200
  geo_redundant_backup_enabled = false
  backup_retention_days        = 7

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.7"
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
