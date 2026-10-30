// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
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

func TestAccApiManagementApiPolicy_listByApiID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_policy", "testlist1")
	r := ApiManagementApiPolicyResource{}

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
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_api_management_api_policy.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_api_management_api_policy.list",
						map[string]knownvalue.Check{
							"api_id":              knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"service_name":        knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r ApiManagementApiPolicyResource) basicQuery() string {
	return `
list "azurerm_api_management_api_policy" "list" {
  provider = azurerm
  config {
    api_management_api_id = azurerm_api_management_api.test.id
  }
}
`
}
