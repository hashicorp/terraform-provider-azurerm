// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package horizondb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HorizondbFirewallRuleResource struct{}

func TestAccHorizondbFirewallRule_basic(t *testing.T) {
	resource := HorizondbFirewallRuleResource{}
	data := acceptance.BuildTestData(t, "azurerm_horizondb_firewall_rule", "test")
	data.ResourceTest(t, resource, []acceptance.TestStep{
		{
			Config: resource.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(resource),
		},
		data.ImportStep(),
	})
}

func TestAccHorizondbFirewallRule_complete(t *testing.T) {
	resource := HorizondbFirewallRuleResource{}
	data := acceptance.BuildTestData(t, "azurerm_horizondb_firewall_rule", "test")
	data.ResourceTest(t, resource, []acceptance.TestStep{
		{
			Config: resource.complete(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(resource),
		},
		data.ImportStep(),
	})
}

func TestAccHorizondbFirewallRule_requiresImport(t *testing.T) {
	resource := HorizondbFirewallRuleResource{}
	data := acceptance.BuildTestData(t, "azurerm_horizondb_firewall_rule", "test")
	data.ResourceTest(t, resource, []acceptance.TestStep{
		{
			Config: resource.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(resource),
		},
		data.RequiresImportErrorStep(resource.requiresImport),
	})
}

func TestAccHorizondbFirewallRule_update(t *testing.T) {
	resource := HorizondbFirewallRuleResource{}
	data := acceptance.BuildTestData(t, "azurerm_horizondb_firewall_rule", "test")
	data.ResourceTest(t, resource, []acceptance.TestStep{
		{
			Config: resource.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(resource),
		},
		data.ImportStep(),
		{
			Config: resource.complete(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(resource),
		},
		data.ImportStep(),
		{
			Config: resource.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(resource),
		},
		data.ImportStep(),
	})
}

func (HorizondbFirewallRuleResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := firewallrules.ParseFirewallRuleID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.HorizonDB.FirewallRulesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (HorizondbFirewallRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_horizondb_firewall_rule" "test" {
  name              = "acctest-rule-%d"
  horizondb_pool_id = azurerm_horizondb_cluster.test.read_write_pool[0].id
  start_ip_address  = "192.0.2.1"
  end_ip_address    = "192.0.2.1"
}
`, HorizondbClusterResource{}.basic(data), data.RandomInteger)
}

func (HorizondbFirewallRuleResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_horizondb_firewall_rule" "test" {
  name              = "acctest-rule-%d"
  horizondb_pool_id = azurerm_horizondb_cluster.test.read_write_pool[0].id
  start_ip_address  = "192.0.2.10"
  end_ip_address    = "192.0.2.20"
}
`, HorizondbClusterResource{}.basic(data), data.RandomInteger)
}

func (resource HorizondbFirewallRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_horizondb_firewall_rule" "import" {
  name              = azurerm_horizondb_firewall_rule.test.name
  horizondb_pool_id = azurerm_horizondb_firewall_rule.test.horizondb_pool_id
  start_ip_address  = azurerm_horizondb_firewall_rule.test.start_ip_address
  end_ip_address    = azurerm_horizondb_firewall_rule.test.end_ip_address
}
`, resource.basic(data))
}
