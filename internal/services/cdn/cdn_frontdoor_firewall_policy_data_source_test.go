// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type CdnFrontDoorFirewallPolicyDataSource struct{}

func TestAccCdnFrontDoorFirewallPolicyDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_frontdoor_firewall_policy", "test")
	d := CdnFrontDoorFirewallPolicyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("redirect_url").MatchesOtherKey(check.That("azurerm_cdn_frontdoor_firewall_policy.test").Key("redirect_url")),
				check.That(data.ResourceName).Key("js_challenge_cookie_expiration_in_minutes").HasValue("30"),
			),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicyDataSource_jsChallengeBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_frontdoor_firewall_policy", "test")
	d := CdnFrontDoorFirewallPolicyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.jsChallengePolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("redirect_url").MatchesOtherKey(check.That("azurerm_cdn_frontdoor_firewall_policy.test").Key("redirect_url")),
				check.That(data.ResourceName).Key("js_challenge_cookie_expiration_in_minutes").HasValue("45"),
			),
		},
	})
}

func (CdnFrontDoorFirewallPolicyDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = azurerm_cdn_frontdoor_firewall_policy.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, CdnFrontDoorFirewallPolicyResource{}.basic(data))
}

func (CdnFrontDoorFirewallPolicyDataSource) jsChallengePolicyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = azurerm_cdn_frontdoor_firewall_policy.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, CdnFrontDoorFirewallPolicyResource{}.jSChallengePolicyBasic(data))
}
