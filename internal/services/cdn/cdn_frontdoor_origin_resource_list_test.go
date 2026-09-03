// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_test

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

func TestAccCdnFrontDoorOrigin_listByOriginGroupID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_origin", "testlist1")
	r := CdnFrontdoorOriginResource{}

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
					querycheck.ExpectLengthAtLeast("azurerm_cdn_frontdoor_origin.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_cdn_frontdoor_origin.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringExact(fmt.Sprintf("acctest-cdnfdorigin-%d", data.RandomInteger)),
							"resource_group_name": knownvalue.StringExact(fmt.Sprintf("acctestrg-cdn-afdx-%d", data.RandomInteger)),
							"profile_name":        knownvalue.StringExact(fmt.Sprintf("acctest-cdnfdprofile-%d", data.RandomInteger)),
							"origin_group_name":   knownvalue.StringExact(fmt.Sprintf("acctest-cdnfd-group-%d", data.RandomInteger)),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r CdnFrontdoorOriginResource) basicQuery(data acceptance.TestData) string {
	return `
list "azurerm_cdn_frontdoor_origin" "list" {
  provider = azurerm
  config {
    cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  }
}
`
}
