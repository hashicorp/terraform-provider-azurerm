// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2022-12-01/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PostgresqlFlexibleServerFirewallRulesResource struct{}

func TestAccPostgresqlFlexibleServerFirewallRules_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_firewall_rules", "test")
	r := PostgresqlFlexibleServerFirewallRulesResource{}
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

func TestAccPostgresqlFlexibleServerFirewallRules_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_firewall_rules", "test")
	r := PostgresqlFlexibleServerFirewallRulesResource{}
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

func TestAccPostgresqlFlexibleServerFirewallRules_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_firewall_rules", "test")
	r := PostgresqlFlexibleServerFirewallRulesResource{}
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

func (PostgresqlFlexibleServerFirewallRulesResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := firewallrules.ParseFirewallRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.FlexibleServerFirewallRuleClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (PostgresqlFlexibleServerFirewallRulesResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_postgresql_flexible_server_firewall_rules" "test" {
  server_id = azurerm_postgresql_flexible_server.test.id
  firewall_rule {
    name             = "acctest-FSFR-%d"
    start_ip_address = "122.122.0.0"
    end_ip_address   = "122.122.0.0"
  }
}
`, PostgresqlFlexibleServerResource{}.basic(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerFirewallRulesResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_postgresql_flexible_server_firewall_rules" "import" {
  server_id = azurerm_postgresql_flexible_server_firewall_rules.test.id
  firewall_rule {
    name             = "acctest-FSFR"
    start_ip_address = "122.122.0.0"
    end_ip_address   = "122.122.0.0"
  }
}
`, r.basic(data))
}

func (r PostgresqlFlexibleServerFirewallRulesResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_postgresql_flexible_server_firewall_rules" "test" {
  server_id = azurerm_postgresql_flexible_server.test.id
  firewall_rule {
    name             = "acctest-FSFR-%d"
    start_ip_address = "123.0.0.0"
    end_ip_address   = "123.0.0.0"
  }
}
`, PostgresqlFlexibleServerResource{}.basic(data), data.RandomInteger)
}
