// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver_test

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

func TestAccPrivateDnsResolverForwardingRule_listByDnsForwardingRulesetID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_forwarding_rule", "testlist1")
	r := PrivateDnsResolverForwardingRuleResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.list(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_private_dns_resolver_forwarding_rule.list", 2),
					querycheck.ExpectIdentity(
						"azurerm_private_dns_resolver_forwarding_rule.list",
						map[string]knownvalue.Check{
							"name":                        knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name":         knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"dns_forwarding_ruleset_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":             knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r PrivateDnsResolverForwardingRuleResource) list(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_resolver_forwarding_rule" "list" {
  count                     = 2
  name                      = "acctest-drfr-list-${count.index}-%d"
  dns_forwarding_ruleset_id = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.id
  domain_name               = "onprem${count.index}.local."
  target_dns_servers {
    ip_address = "10.10.0.1"
    port       = 53
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PrivateDnsResolverForwardingRuleResource) basicQuery(data acceptance.TestData) string {
	return `
list "azurerm_private_dns_resolver_forwarding_rule" "list" {
  provider = azurerm
  config {
    dns_forwarding_ruleset_id = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.id
  }
}
`
}
