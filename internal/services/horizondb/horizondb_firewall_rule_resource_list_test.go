// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package horizondb_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
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
	r := HorizondbFirewallRuleResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_horizondb_firewall_rule.list", 3),
					querycheck.ExpectIdentity(
						"azurerm_horizondb_firewall_rule.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"cluster_name":        knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"pool_name":           knownvalue.StringExact("primary"),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (HorizondbFirewallRuleResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_horizondb_firewall_rule" "test" {
  count            = 3
  name             = "acctest-rule-%d-${count.index}"
  horizondb_pool_id = azurerm_horizondb_cluster.test.read_write_pool[0].id
  start_ip_address = "192.0.2.${count.index}"
  end_ip_address   = "192.0.2.${count.index}"
}
`, HorizondbClusterResource{}.basic(data), data.RandomInteger)
}

func (HorizondbFirewallRuleResource) basicQuery() string {
	return `
list "azurerm_horizondb_firewall_rule" "list" {
  provider = azurerm
  config {
    horizondb_pool_id = azurerm_horizondb_cluster.test.read_write_pool[0].id
  }
}
`
}
