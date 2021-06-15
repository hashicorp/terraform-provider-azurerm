package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PostgresqlFlexibleServerFirewallRuleResource struct {
}

func TestAccPostgresqlFlexibleServerFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_firewall_rule", "test")
	r := PostgresqlFlexibleServerFirewallRuleResource{}
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

func TestAccPostgresqlFlexibleServerFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_firewall_rule", "test")
	r := PostgresqlFlexibleServerFirewallRuleResource{}
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

func TestAccPostgresqlFlexibleServerFirewallRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_firewall_rule", "test")
	r := PostgresqlFlexibleServerFirewallRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (PostgresqlFlexibleServerFirewallRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FlexibleServerFirewallRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.FlexibleServerFirewallRuleClient.Get(ctx, id.ResourceGroup, id.FlexibleServerName, id.FirewallRuleName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Postgresql Flexible Server Firewall Rule %q ( Flexible Server: %q / resource group: %q): %+v", id.FirewallRuleName, id.FlexibleServerName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.FirewallRuleProperties != nil), nil
}

func (PostgresqlFlexibleServerFirewallRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server_firewall_rule" "test" {
  name             = "acctest-FSFR-%d"
  server_id        = azurerm_postgresql_flexible_server.test.id
  start_ip_address = "122.122.0.0"
  end_ip_address   = "122.122.0.0"
}
`, PostgresqlFlexibleServerResource{}.basic(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerFirewallRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server_firewall_rule" "import" {
  name             = azurerm_postgresql_flexible_server_firewall_rule.test.name
  server_id        = azurerm_postgresql_flexible_server_firewall_rule.test.server_id
  start_ip_address = azurerm_postgresql_flexible_server_firewall_rule.test.start_ip_address
  end_ip_address   = azurerm_postgresql_flexible_server_firewall_rule.test.end_ip_address
}
`, r.basic(data))
}

func (r PostgresqlFlexibleServerFirewallRuleResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server_firewall_rule" "test" {
  name             = "acctest-FSFR-%d"
  server_id        = azurerm_postgresql_flexible_server.test.id
  start_ip_address = "123.0.0.0"
  end_ip_address   = "123.0.0.0"
}
`, PostgresqlFlexibleServerResource{}.basic(data), data.RandomInteger)
}
