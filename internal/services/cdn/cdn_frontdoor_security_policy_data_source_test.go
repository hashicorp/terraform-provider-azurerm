// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type CdnFrontDoorSecurityPolicyDataSource struct{}

func TestAccCdnFrontDoorSecurityPolicyDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_frontdoor_security_policy", "test")
	d := CdnFrontDoorSecurityPolicyDataSource{}
	d.preCheck(t)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("cdn_frontdoor_profile_id").MatchesOtherKey(check.That("azurerm_cdn_frontdoor_security_policy.test").Key("cdn_frontdoor_profile_id")),
				check.That(data.ResourceName).Key("security_policies.0.firewall.0.cdn_frontdoor_firewall_policy_id").MatchesOtherKey(check.That("azurerm_cdn_frontdoor_security_policy.test").Key("security_policies.0.firewall.0.cdn_frontdoor_firewall_policy_id")),
				check.That(data.ResourceName).Key("security_policies.0.firewall.0.association.0.domain.0.cdn_frontdoor_domain_id").MatchesOtherKey(check.That("azurerm_cdn_frontdoor_security_policy.test").Key("security_policies.0.firewall.0.association.0.domain.0.cdn_frontdoor_domain_id")),
				check.That(data.ResourceName).Key("security_policies.0.firewall.0.association.0.patterns_to_match.0").HasValue("/*"),
			),
		},
	})
}

func (CdnFrontDoorSecurityPolicyDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cdn_frontdoor_security_policy" "test" {
  name                = azurerm_cdn_frontdoor_security_policy.test.name
  profile_name        = azurerm_cdn_frontdoor_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, CdnFrontDoorSecurityPolicyResource{}.basic(data))
}

func (CdnFrontDoorSecurityPolicyDataSource) preCheck(t *testing.T) {
	if os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skipf("`ARM_TEST_DATA_RESOURCE_GROUP` must be set for acceptance tests!")
	}
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" {
		t.Skipf("`ARM_TEST_DNS_ZONE` must be set for acceptance tests!")
	}
}
