// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PrivateDNSResolverDnsForwardingRulesetDataSource struct{}

func TestAccPrivateDNSResolverDnsForwardingRulesetDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_resolver_dns_forwarding_ruleset", "test")
	d := PrivateDNSResolverDnsForwardingRulesetDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("private_dns_resolver_outbound_endpoint_ids.0").Exists(),
			),
		},
	})
}

func (d PrivateDNSResolverDnsForwardingRulesetDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_private_dns_resolver_dns_forwarding_ruleset" "test" {
  name                = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, PrivateDNSResolverDnsForwardingRulesetResource{}.basic(data))
}
