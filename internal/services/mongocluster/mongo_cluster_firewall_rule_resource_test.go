// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mongocluster_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2025-09-01/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MongoClusterFirewallRuleResource struct{}

func TestAccMongoClusterFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster_firewall_rule", "test")
	r := MongoClusterFirewallRuleResource{}

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

func TestAccMongoClusterFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster_firewall_rule", "test")
	r := MongoClusterFirewallRuleResource{}

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

func TestAccMongoClusterFirewallRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster_firewall_rule", "test")
	r := MongoClusterFirewallRuleResource{}

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

func (r MongoClusterFirewallRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := firewallrules.ParseFirewallRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MongoCluster.FirewallRulesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r MongoClusterFirewallRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mongo_cluster_firewall_rule" "test" {
  name             = "acctest-mcfr-%d"
  mongo_cluster_id = azurerm_mongo_cluster.test.id
  start_ip_address = "10.0.0.1"
  end_ip_address   = "10.0.0.255"
}
`, r.template(data), data.RandomInteger)
}

func (r MongoClusterFirewallRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mongo_cluster_firewall_rule" "import" {
  name             = azurerm_mongo_cluster_firewall_rule.test.name
  mongo_cluster_id = azurerm_mongo_cluster_firewall_rule.test.mongo_cluster_id
  start_ip_address = azurerm_mongo_cluster_firewall_rule.test.start_ip_address
  end_ip_address   = azurerm_mongo_cluster_firewall_rule.test.end_ip_address
}
`, r.basic(data))
}

func (r MongoClusterFirewallRuleResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mongo_cluster_firewall_rule" "test-2" {
  name             = "acctest-mcfr2-%[2]d"
  mongo_cluster_id = azurerm_mongo_cluster.test.id
  start_ip_address = "10.0.2.1"
  end_ip_address   = "10.0.2.255"
}

resource "azurerm_mongo_cluster_firewall_rule" "test" {
  name             = "acctest-mcfr-%[2]d"
  mongo_cluster_id = azurerm_mongo_cluster.test.id
  start_ip_address = "10.0.1.1"
  end_ip_address   = "10.0.1.255"
}
`, r.template(data), data.RandomInteger)
}

func (r MongoClusterFirewallRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-mc-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mongo_cluster" "test" {
  name                   = "acctest-mc-%[1]d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_username = "adminuser"
  administrator_password = "P@ssw0rd1234!"
  shard_count            = 1
  compute_tier           = "M10"
  high_availability_mode = "Disabled"
  storage_size_in_gb     = 32
  version                = "7.0"
}
`, data.RandomInteger, data.Locations.Primary)
}
