package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MsSqlFirewallRuleResource struct{}

func TestAccMsSqlFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_firewall_rule", "test")
	r := MsSqlFirewallRuleResource{}

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

func TestAccMsSqlFirewallRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_firewall_rule", "test")
	r := MsSqlFirewallRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_firewall_rule", "test")
	r := MsSqlFirewallRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r MsSqlFirewallRuleResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.FirewallRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.FirewallRulesClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving MSSQL %s: %+v", id.String(), err)
	}

	return utils.Bool(true), nil
}

func (MsSqlFirewallRuleResource) template(data acceptance.TestData) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r MsSqlFirewallRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_firewall_rule" "test" {
  name             = "acctestsqlserver%[2]d"
  server_id        = azurerm_mssql_server.test.id
  start_ip_address = "0.0.0.0"
  end_ip_address   = "255.255.255.255"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlFirewallRuleResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_firewall_rule" "test" {
  name             = "acctestsqlserver%[2]d"
  server_id        = azurerm_mssql_server.test.id
  start_ip_address = "10.0.17.62"
  end_ip_address   = "10.0.17.64"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlFirewallRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_firewall_rule" "import" {
  name             = azurerm_mssql_firewall_rule.test.name
  server_id        = azurerm_mssql_firewall_rule.test.server_id
  start_ip_address = azurerm_mssql_firewall_rule.test.start_ip_address
  end_ip_address   = azurerm_mssql_firewall_rule.test.end_ip_address
}
`, r.basic(data))
}
