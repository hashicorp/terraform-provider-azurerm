// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package horizondb_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

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
		{Config: resource.basic(data), Check: check.That(data.ResourceName).ExistsInAzure(resource)},
		data.ImportStep(),
	})
}

func TestAccHorizondbFirewallRule_complete(t *testing.T) {
	resource := HorizondbFirewallRuleResource{}
	data := acceptance.BuildTestData(t, "azurerm_horizondb_firewall_rule", "test")
	data.ResourceTest(t, resource, []acceptance.TestStep{
		{Config: resource.complete(data), Check: check.That(data.ResourceName).ExistsInAzure(resource)},
		data.ImportStep(),
	})
}

func TestAccHorizondbFirewallRule_requiresImport(t *testing.T) {
	resource := HorizondbFirewallRuleResource{}
	data := acceptance.BuildTestData(t, "azurerm_horizondb_firewall_rule", "test")
	data.ResourceTest(t, resource, []acceptance.TestStep{
		{Config: resource.basic(data), Check: check.That(data.ResourceName).ExistsInAzure(resource)},
		data.RequiresImportErrorStep(resource.requiresImport),
	})
}

func TestAccHorizondbFirewallRule_update(t *testing.T) {
	resource := HorizondbFirewallRuleResource{}
	data := acceptance.BuildTestData(t, "azurerm_horizondb_firewall_rule", "test")
	data.ResourceTest(t, resource, []acceptance.TestStep{
		{Config: resource.basic(data), Check: check.That(data.ResourceName).ExistsInAzure(resource)},
		data.ImportStep(),
		{Config: resource.complete(data), Check: check.That(data.ResourceName).ExistsInAzure(resource)},
		data.ImportStep(),
		{Config: resource.config(data, "192.0.2.10", "192.0.2.20", "Updated acceptance test firewall rule"), Check: check.That(data.ResourceName).ExistsInAzure(resource)},
		data.ImportStep(),
		{Config: resource.config(data, "192.0.2.10", "192.0.2.20", ""), Check: check.That(data.ResourceName).ExistsInAzure(resource)},
		data.ImportStep(),
		{Config: resource.basic(data), Check: check.That(data.ResourceName).ExistsInAzure(resource)},
		data.ImportStep(),
	})
}

func TestAccHorizondbFirewallRule_invalidRange(t *testing.T) {
	resource := HorizondbFirewallRuleResource{}
	data := acceptance.BuildTestData(t, "azurerm_horizondb_firewall_rule", "test")
	data.ResourceTest(t, resource, []acceptance.TestStep{
		{
			Config:      resource.config(data, "192.0.2.10", "192.0.2.1", ""),
			ExpectError: regexp.MustCompile("`end_ip_address` must be greater than or equal to `start_ip_address`"),
		},
	})
}

func TestAccHorizondbFirewallRule_overwriteExisting(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_firewall_rule", "test")
	r := HorizondbFirewallRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.overwritePrerequisites(data),
			Check:  data.CheckWithClientForResource(r.prepareOverwrite(data), "azurerm_horizondb_cluster.test"),
		},
		{
			Config: r.overwriteExisting(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
	})
}

func (HorizondbFirewallRuleResource) prepareOverwrite(data acceptance.TestData) acceptance.ClientCheckFunc {
	return func(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
		poolId, err := firewallrules.ParsePoolID(state.Attributes["read_write_pool.0.id"])
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
		defer cancel()

		id := firewallrules.NewFirewallRuleID(poolId.SubscriptionId, poolId.ResourceGroupName, poolId.ClusterName, poolId.PoolName, fmt.Sprintf("acctest-rule-%d", data.RandomInteger))
		payload := firewallrules.FirewallRule{
			Properties: &firewallrules.FirewallRuleProperties{
				StartIPAddress: "192.0.2.1",
				EndIPAddress:   "192.0.2.1",
				Description:    pointer.To("Before overwrite"),
			},
		}
		if err := client.HorizonDB.FirewallRulesClient.CreateOrReplaceThenPoll(ctx, id, payload); err != nil {
			return fmt.Errorf("creating pre-existing %s: %+v", id, err)
		}
		return nil
	}
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

func (resource HorizondbFirewallRuleResource) basic(data acceptance.TestData) string {
	return resource.config(data, "192.0.2.1", "192.0.2.1", "")
}

func (resource HorizondbFirewallRuleResource) complete(data acceptance.TestData) string {
	return resource.config(data, "192.0.2.10", "192.0.2.20", "Acceptance test firewall rule")
}

func (HorizondbFirewallRuleResource) config(data acceptance.TestData, start, end, description string) string {
	descriptionConfig := ""
	if description != "" {
		descriptionConfig = fmt.Sprintf("description = %q", description)
	}
	return fmt.Sprintf(`
%s

resource "azurerm_horizondb_firewall_rule" "test" {
  name              = "acctest-rule-%d"
  horizondb_pool_id = azurerm_horizondb_cluster.test.read_write_pool[0].id
  start_ip_address  = %q
  end_ip_address    = %q
  %s
}
`, HorizondbClusterResource{}.withPool(data, "primary", 1, "Disabled"), data.RandomInteger, start, end, descriptionConfig)
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

func (HorizondbFirewallRuleResource) overwritePrerequisites(data acceptance.TestData) string {
	return HorizondbClusterResource{}.withPoolAndFeatures(data, "primary", 1, "Disabled", "skip_import_check_on_create_and_allow_overwriting_existing_resources = true")
}

func (r HorizondbFirewallRuleResource) overwriteExisting(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_horizondb_firewall_rule" "test" {
  name              = "acctest-rule-%d"
  horizondb_pool_id = azurerm_horizondb_cluster.test.read_write_pool[0].id
  start_ip_address  = "192.0.2.10"
  end_ip_address    = "192.0.2.20"
  description       = "Acceptance test firewall rule"
}
`, r.overwritePrerequisites(data), data.RandomInteger)
}
