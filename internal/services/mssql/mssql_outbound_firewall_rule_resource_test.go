// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/outboundfirewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MsSqlOutboundFirewallRuleResource struct{}

func TestAccMsSqlOutboundFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_outbound_firewall_rule", "test")
	r := MsSqlOutboundFirewallRuleResource{}

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

func TestAccMsSqlOutboundFirewallRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_outbound_firewall_rule", "test")
	r := MsSqlOutboundFirewallRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlOutboundFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_outbound_firewall_rule", "test")
	r := MsSqlOutboundFirewallRuleResource{}

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

func (r MsSqlOutboundFirewallRuleResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := outboundfirewallrules.ParseOutboundFirewallRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.OutboundFirewallRulesClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(true), nil
}

func (MsSqlOutboundFirewallRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "msincredible"
  administrator_login_password = "P@55W0rD!!%[3]s"

  outbound_network_restriction_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r MsSqlOutboundFirewallRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_outbound_firewall_rule" "test" {
  name      = "sql%[2]d.database.windows.net"
  server_id = azurerm_mssql_server.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlOutboundFirewallRuleResource) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_outbound_firewall_rule" "test" {
  name      = "sql%[2]d.database.windows.net"
  server_id = azurerm_mssql_server.test.id
}

resource "azurerm_mssql_outbound_firewall_rule" "test2" {
  name      = "sql2%[2]d.database.windows.net"
  server_id = azurerm_mssql_server.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlOutboundFirewallRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_outbound_firewall_rule" "import" {
  name      = azurerm_mssql_outbound_firewall_rule.test.name
  server_id = azurerm_mssql_outbound_firewall_rule.test.server_id
}
`, r.basic(data))
}
