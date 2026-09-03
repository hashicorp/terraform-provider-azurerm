// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccLbBackendAddressPool_listByLoadBalancerID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "testlist1")
	r := LbBackendAddressPoolResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_lb_backend_address_pool.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_lb_backend_address_pool.list",
						map[string]knownvalue.Check{
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"load_balancer_name":  knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"name":                knownvalue.StringExact("pool"),
						},
					),
				},
			},
		},
	})
}

func (r LbBackendAddressPoolResource) basicQuery(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_lb_backend_address_pool" "list" {
  provider = azurerm
  config {
    loadbalancer_id = "/subscriptions/%s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.Network/loadBalancers/acctestlb-%[2]d"
  }
}
`, data.Subscriptions.Primary, data.RandomInteger)
}
