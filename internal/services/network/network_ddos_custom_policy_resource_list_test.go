// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_test

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

func TestAccNetworkDDoSCustomPolicy_listBySubscriptionAndRG(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_custom_policy", "testlist")
	r := NetworkDdosCustomPolicyResource{}
	listResourceAddress := "azurerm_network_ddos_custom_policy.list"

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
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3), // expect at least the 3 we created
					querycheck.ExpectIdentity(
						listResourceAddress,
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3), // expect exactly the 3 we created in that resource group
					querycheck.ExpectIdentity(
						listResourceAddress,
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

// provision multiple Network DDoS Custom Policy resources for testing
func (r NetworkDdosCustomPolicyResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_network_ddos_custom_policy" "test" {
  count = 3

  name                = "acctest-ddoscp-${count.index}-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  detection_rule {
    name               = "detectionRuleTcp"
    packets_per_second = 1000000
    traffic_type       = "Tcp"
  }
}
`, r.template(data), data.RandomInteger)
}

// define the basic list query for testing
func (r NetworkDdosCustomPolicyResource) basicQuery() string {
	return `
list "azurerm_network_ddos_custom_policy" "list" {
  provider = azurerm
  config {
  }
}
`
}

// define the list query for testing by resource group name
func (r NetworkDdosCustomPolicyResource) basicQueryByResourceGroupName() string {
	return `
list "azurerm_network_ddos_custom_policy" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}
