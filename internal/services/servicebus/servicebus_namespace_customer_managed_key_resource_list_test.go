// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_test

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

func TestAccServicebusNamespaceCustomerManagedKey_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_customer_managed_key", "testlist")
	r := ServicebusNamespaceCustomerManagedKeyResource{}

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
				Config: r.basicListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_servicebus_namespace_customer_managed_key.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_servicebus_namespace_customer_managed_key.list",
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
				Config: r.basicListQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_servicebus_namespace_customer_managed_key.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_servicebus_namespace_customer_managed_key.list",
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

func (r ServicebusNamespaceCustomerManagedKeyResource) basicListQuery() string {
	return `
list "azurerm_servicebus_namespace_customer_managed_key" "list" {
  provider = azurerm
  config {}
}
`
}

func (r ServicebusNamespaceCustomerManagedKeyResource) basicListQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_servicebus_namespace_customer_managed_key" "list" {
  provider = azurerm
  config {
    subscription_id     = "%s"
    resource_group_name = "acctest-sb-%d"
  }
}
`, data.Subscriptions.Primary, data.RandomInteger)
}
