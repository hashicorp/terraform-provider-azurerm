// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PrivateDNSResolverVirtualNetworkLinkDataSource struct{}

func TestAccPrivateDNSResolverVirtualNetworkLinkDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_resolver_virtual_network_link", "test")
	d := PrivateDNSResolverVirtualNetworkLinkDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("dns_forwarding_ruleset_id").Exists(),
				check.That(data.ResourceName).Key("virtual_network_id").Exists(),
			),
		},
	})
}

func (d PrivateDNSResolverVirtualNetworkLinkDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_private_dns_resolver_virtual_network_link" "test" {
  name                      = azurerm_private_dns_resolver_virtual_network_link.test.name
  dns_forwarding_ruleset_id = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.id
}
`, PrivateDNSResolverVirtualNetworkLinkResource{}.basic(data))
}
