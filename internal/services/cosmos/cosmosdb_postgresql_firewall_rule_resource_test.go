// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CosmosDbPostgreSQLFirewallRuleResource struct{}

func TestCosmosDbPostgreSQLFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_postgresql_firewall_rule", "test")
	r := CosmosDbPostgreSQLFirewallRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "10.0.17.62", "10.0.17.64"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestCosmosDbPostgreSQLFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_postgresql_firewall_rule", "test")
	r := CosmosDbPostgreSQLFirewallRuleResource{}
	startIPAddress := "10.0.17.62"
	endIPAddress := "10.0.17.64"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, startIPAddress, endIPAddress),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport(data, startIPAddress, endIPAddress)
		}),
	})
}

func TestCosmosDbPostgreSQLFirewallRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_postgresql_firewall_rule", "test")
	r := CosmosDbPostgreSQLFirewallRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "10.0.17.62", "10.0.17.64"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "10.0.17.65", "10.0.17.67"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CosmosDbPostgreSQLFirewallRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := firewallrules.ParseFirewallRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cosmos.FirewallRulesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r CosmosDbPostgreSQLFirewallRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-pshscfwr-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_postgresql_cluster" "test" {
  name                            = "acctestcluster%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  administrator_login_password    = "H@Sh1CoR3!"
  coordinator_storage_quota_in_mb = 131072
  coordinator_vcore_count         = 2
  node_count                      = 0
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CosmosDbPostgreSQLFirewallRuleResource) basic(data acceptance.TestData, startIPAddress, endIPAddress string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_postgresql_firewall_rule" "test" {
  name             = "acctest-pshscfwr-%d"
  cluster_id       = azurerm_cosmosdb_postgresql_cluster.test.id
  start_ip_address = "%s"
  end_ip_address   = "%s"
}
`, r.template(data), data.RandomInteger, startIPAddress, endIPAddress)
}

func (r CosmosDbPostgreSQLFirewallRuleResource) requiresImport(data acceptance.TestData, startIPAddress, endIPAddress string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_postgresql_firewall_rule" "import" {
  name             = azurerm_cosmosdb_postgresql_firewall_rule.test.name
  cluster_id       = azurerm_cosmosdb_postgresql_firewall_rule.test.cluster_id
  start_ip_address = azurerm_cosmosdb_postgresql_firewall_rule.test.start_ip_address
  end_ip_address   = azurerm_cosmosdb_postgresql_firewall_rule.test.end_ip_address
}
`, r.basic(data, startIPAddress, endIPAddress))
}
