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

type DNSResolverOutboundEndpointDataSource struct{}

func TestAccDNSResolverOutboundEndpointDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_resolver_outbound_endpoint", "test")
	d := DNSResolverOutboundEndpointDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("private_dns_resolver_id").Exists(),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
			),
		},
	})
}

func (d DNSResolverOutboundEndpointDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_private_dns_resolver_outbound_endpoint" "test" {
  name                    = azurerm_private_dns_resolver_outbound_endpoint.test.name
  private_dns_resolver_id = azurerm_private_dns_resolver.test.id
}
`, DNSResolverOutboundEndpointResource{}.basic(data))
}
