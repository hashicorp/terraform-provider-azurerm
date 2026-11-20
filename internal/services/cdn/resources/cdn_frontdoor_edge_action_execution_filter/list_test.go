// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action_execution_filter_test

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

func TestAccCdnFrontdoorEdgeActionExecutionFilter_list_basic(t *testing.T) {
	r := CdnFrontdoorEdgeActionExecutionFilterResource{}
	listResourceAddress := "azurerm_cdn_frontdoor_edge_action_execution_filter.list"

	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_edge_action_execution_filter", "test")

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
					querycheck.ExpectLength(listResourceAddress, 2),
					querycheck.ExpectIdentity(
						listResourceAddress,
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"edge_action_name":    knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r CdnFrontdoorEdgeActionExecutionFilterResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_edge_action_execution_filter" "test" {
  count = 2

  name                                 = "acctesteaf${count.index}%d"
  cdn_frontdoor_edge_action_id         = azurerm_cdn_frontdoor_edge_action.test.id
  cdn_frontdoor_edge_action_version_id = azurerm_cdn_frontdoor_edge_action_version.first.id
  header_name                          = "x-canary-${count.index}"
  header_value                         = "true"
}
`, r.template(data), data.RandomInteger)
}

func (CdnFrontdoorEdgeActionExecutionFilterResource) basicQuery() string {
	return `
list "azurerm_cdn_frontdoor_edge_action_execution_filter" "list" {
  provider = azurerm
  config {
    cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.test.id
  }
}
`
}
