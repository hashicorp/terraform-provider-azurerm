// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package loadbalancer_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccLbRule_listByLoadBalancerID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "testlist1")
	r := LbRuleResource{}

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
					querycheck.ExpectLength("azurerm_lb_rule.list", 3),
					querycheck.ExpectIdentity(
						"azurerm_lb_rule.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger % 100000000))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"load_balancer_name":  knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
					querycheck.ExpectResourceKnownValues(
						"azurerm_lb_rule.list",
						queryfilter.ByDisplayName(knownvalue.StringExact(fmt.Sprintf("acctest-lb-rule-%d0", data.RandomInteger%100000000))),
						[]querycheck.KnownValueCheck{
							{
								Path:       tfjsonpath.New("loadbalancer_id"),
								KnownValue: knownvalue.StringExact(fmt.Sprintf("/subscriptions/%s/resourceGroups/acctestRG-lb-%d/providers/Microsoft.Network/loadBalancers/acctest-loadbalancer-%d", data.Subscriptions.Primary, data.RandomInteger, data.RandomInteger)),
							},
						},
					),
				},
			},
		},
	})
}

func (r LbRuleResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "test" {
  count                          = 3
  name                           = "acctest-lb-rule-%d${count.index}"
  loadbalancer_id                = azurerm_lb.test.id
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
  protocol                       = "Tcp"
  frontend_port                  = "338${count.index}"
  backend_port                   = "338${count.index}"
}
`, r.template(data), data.RandomInteger%100000000)
}

func (r LbRuleResource) basicQuery() string {
	return `
list "azurerm_lb_rule" "list" {
  provider = azurerm
  include_resource = true
  config {
    load_balancer_id = azurerm_lb.test.id
  }
}
`
}
