// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type CdnFrontDoorBatchRuleSetDataSource struct{}

func TestAccCdnFrontDoorBatchRuleSetDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_frontdoor_batch_rule_set", "test")
	d := CdnFrontDoorBatchRuleSetDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("cdn_frontdoor_profile_id").MatchesOtherKey(check.That("azurerm_cdn_frontdoor_profile.test").Key("id")),
				check.That(data.ResourceName).Key("rules.#").HasValue("2"),
				check.That(data.ResourceName).Key("rules.0.order").HasValue("1"),
				check.That(data.ResourceName).Key("rules.1.order").HasValue("2"),
			),
		},
	})
}

func TestAccCdnFrontDoorBatchRuleSetDataSource_nonBatchRuleSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_frontdoor_batch_rule_set", "test")
	d := CdnFrontDoorBatchRuleSetDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config:      d.nonBatchRuleSet(data),
			ExpectError: regexp.MustCompile("was not provisioned using batch mode, and cannot be read by this data source"),
		},
	})
}

func (CdnFrontDoorBatchRuleSetDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  name                = azurerm_cdn_frontdoor_batch_rule_set.test.name
  profile_name        = azurerm_cdn_frontdoor_profile.test.name
  resource_group_name = azurerm_cdn_frontdoor_profile.test.resource_group_name
}
`, CdnFrontdoorBatchRuleSetResource{}.complete(data))
}

func (CdnFrontDoorBatchRuleSetDataSource) nonBatchRuleSet(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cdn_frontdoor_batch_rule_set" "test" {
  name                = azurerm_cdn_frontdoor_rule_set.test.name
  profile_name        = azurerm_cdn_frontdoor_profile.test.name
  resource_group_name = azurerm_cdn_frontdoor_profile.test.resource_group_name
}
`, CdnFrontDoorRuleSetResource{}.basic(data, true))
}
