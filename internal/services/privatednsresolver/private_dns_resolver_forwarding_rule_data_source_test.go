// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PrivateDNSResolverForwardingRuleDataSource struct{}

func TestAccPrivateDNSResolverForwardingRuleDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_resolver_forwarding_rule", "test")
	d := PrivateDNSResolverForwardingRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("dns_forwarding_ruleset_id").Exists(),
				check.That(data.ResourceName).Key("domain_name").HasValue("onprem.local."),
				check.That(data.ResourceName).Key("target_dns_servers.0.ip_address").HasValue("10.10.0.1"),
				check.That(data.ResourceName).Key("target_dns_servers.0.port").HasValue("53"),
			),
		},
	})
}

func (d PrivateDNSResolverForwardingRuleDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_private_dns_resolver_forwarding_rule" "test" {
  name                      = azurerm_private_dns_resolver_forwarding_rule.test.name
  dns_forwarding_ruleset_id = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.id
}
`, PrivateDNSResolverForwardingRuleResource{}.basic(data))
}
