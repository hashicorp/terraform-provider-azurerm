// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package horizondb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccHorizondbFirewallRule_listByPoolID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_firewall_rule", "list")
	checks := []querycheck.QueryResultCheck{
		querycheck.ExpectLengthAtLeast("azurerm_horizondb_firewall_rule.list", 3),
	}
	for index := range 3 {
		checks = append(checks, querycheck.ExpectIdentity("azurerm_horizondb_firewall_rule.list", map[string]knownvalue.Check{
			"name":                knownvalue.StringExact(fmt.Sprintf("acctest-rule-%d-%d", data.RandomInteger, index)),
			"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
			"resource_group_name": knownvalue.StringExact(fmt.Sprintf("acctestRG-horizondb-%d", data.RandomInteger)),
			"cluster_name":        knownvalue.StringExact(fmt.Sprintf("acctest-horizondb-%d", data.RandomInteger)),
			"pool_name":           knownvalue.StringExact("primary"),
		}))
	}
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_14_0)},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{Config: fmt.Sprintf(`
%s

resource "azurerm_horizondb_firewall_rule" "test" {
  count             = 3
  name              = "acctest-rule-%d-${count.index}"
  horizondb_pool_id = azurerm_horizondb_cluster.test.read_write_pool[0].id
  start_ip_address  = "192.0.2.${count.index}"
  end_ip_address    = "192.0.2.${count.index}"
}
`, HorizondbClusterResource{}.withPool(data, "primary", 1, "Disabled"), data.RandomInteger)},
			{
				Query: true,
				Config: `
list "azurerm_horizondb_firewall_rule" "list" {
  provider = azurerm
  config {
    horizondb_pool_id = azurerm_horizondb_cluster.test.read_write_pool[0].id
  }
}
`,
				QueryResultChecks: checks,
			},
		},
	})
}
