// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PostgreSQLFirewallRuleResource struct{}

func TestAccPostgreSQLFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_firewall_rule", "test")
	r := PostgreSQLFirewallRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("start_ip_address").HasValue("0.0.0.0"),
				check.That(data.ResourceName).Key("end_ip_address").HasValue("255.255.255.255"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPostgreSQLFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_firewall_rule", "test")
	r := PostgreSQLFirewallRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("start_ip_address").HasValue("0.0.0.0"),
				check.That(data.ResourceName).Key("end_ip_address").HasValue("255.255.255.255"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t PostgreSQLFirewallRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := firewallrules.ParseFirewallRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.FirewallRulesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Postgresql Firewall Rule (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (PostgreSQLFirewallRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_2"

  storage_mb                   = 51200
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement_enabled      = true
}

resource "azurerm_postgresql_firewall_rule" "test" {
  name                = "acctest-PSQL-fwrule-%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "255.255.255.255"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r PostgreSQLFirewallRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_firewall_rule" "import" {
  name                = azurerm_postgresql_firewall_rule.test.name
  resource_group_name = azurerm_postgresql_firewall_rule.test.resource_group_name
  server_name         = azurerm_postgresql_firewall_rule.test.server_name
  start_ip_address    = azurerm_postgresql_firewall_rule.test.start_ip_address
  end_ip_address      = azurerm_postgresql_firewall_rule.test.end_ip_address
}
`, r.basic(data))
}
