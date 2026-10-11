// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_test

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

func TestAccCdnFrontDoorBatchRuleSet_listByProfileID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_batch_rule_set", "testlist1")
	r := CdnFrontdoorBatchRuleSetResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicAttachedRoute(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_cdn_frontdoor_batch_rule_set.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_cdn_frontdoor_batch_rule_set.list",
						map[string]knownvalue.Check{
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"profile_name":        knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"name":                knownvalue.StringExact(fmt.Sprintf("acctestBatchRuleSet%d", data.RandomInteger)),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r CdnFrontdoorBatchRuleSetResource) basicQuery() string {
	return `
list "azurerm_cdn_frontdoor_batch_rule_set" "list" {
  provider = azurerm
  config {
    cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  }
}`
}
