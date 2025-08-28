// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatedns_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

type PrivateDnsZoneVirtualNetworkLinkDataSource struct{}

func TestAccDataSourcePrivateDnsZoneVirtualNetworkLink_basic(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("Skipping test as it has deprecated fields")
	}
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_zone_virtual_network_link", "test")
	r := PrivateDnsZoneVirtualNetworkLinkDataSource{}

	resourceName := "azurerm_private_dns_zone_virtual_network_link.test"
	zoneName := "azurerm_private_dns_zone.test"
	vnetName := "azurerm_virtual_network.test"

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").MatchesOtherKey(check.That(resourceName).Key("id")),
				check.That(data.ResourceName).Key("name").MatchesOtherKey(check.That(resourceName).Key("name")),
				check.That(data.ResourceName).Key("virtual_network_id").MatchesOtherKey(check.That(vnetName).Key("id")),
				check.That(data.ResourceName).Key("private_dns_zone_id").MatchesOtherKey(check.That(zoneName).Key("id")),
				check.That(data.ResourceName).Key("registration_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourcePrivateDnsZoneVirtualNetworkLink_resolutionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_zone_virtual_network_link", "test")
	r := PrivateDnsZoneVirtualNetworkLinkDataSource{}

	resourceName := "azurerm_private_dns_zone_virtual_network_link.test"
	zoneName := "azurerm_private_dns_zone.test"
	vnetName := "azurerm_virtual_network.test"

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.resolutionPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").MatchesOtherKey(check.That(resourceName).Key("id")),
				check.That(data.ResourceName).Key("name").MatchesOtherKey(check.That(resourceName).Key("name")),
				check.That(data.ResourceName).Key("resource_group_name").MatchesOtherKey(check.That(resourceName).Key("resource_group_name")),
				check.That(data.ResourceName).Key("virtual_network_id").MatchesOtherKey(check.That(vnetName).Key("id")),
				check.That(data.ResourceName).Key("private_dns_zone_name").MatchesOtherKey(check.That(zoneName).Key("name")),
				check.That(data.ResourceName).Key("resolution_policy").HasValue("NxDomainRedirect"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func (PrivateDnsZoneVirtualNetworkLinkDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                = azurerm_private_dns_zone_virtual_network_link.test.name
  private_dns_zone_id = azurerm_private_dns_zone.test.id
}
`, PrivateDnsZoneVirtualNetworkLinkResource{}.basic(data))
}

func TestAccDataSourcePrivateDnsZoneVirtualNetworkLink_resource_group_name(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("Skipping test as it has deprecated fields")
	}
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_zone_virtual_network_link", "test")
	r := PrivateDnsZoneVirtualNetworkLinkDataSource{}

	resourceName := "azurerm_private_dns_zone_virtual_network_link.test"
	zoneName := "azurerm_private_dns_zone.test"
	vnetName := "azurerm_virtual_network.test"

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.resource_group_name(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").MatchesOtherKey(check.That(resourceName).Key("id")),
				check.That(data.ResourceName).Key("name").MatchesOtherKey(check.That(resourceName).Key("name")),
				check.That(data.ResourceName).Key("resource_group_name").MatchesOtherKey(check.That(resourceName).Key("resource_group_name")),
				check.That(data.ResourceName).Key("virtual_network_id").MatchesOtherKey(check.That(vnetName).Key("id")),
				check.That(data.ResourceName).Key("private_dns_zone_name").MatchesOtherKey(check.That(zoneName).Key("name")),
				check.That(data.ResourceName).Key("registration_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func (PrivateDnsZoneVirtualNetworkLinkDataSource) resource_group_name(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = azurerm_private_dns_zone_virtual_network_link.test.name
  resource_group_name   = azurerm_resource_group.test.name
  private_dns_zone_name = azurerm_private_dns_zone.test.name
}
`, PrivateDnsZoneVirtualNetworkLinkResource{}.basic(data))
}

func (PrivateDnsZoneVirtualNetworkLinkDataSource) resolutionPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = azurerm_private_dns_zone_virtual_network_link.test.name
  resource_group_name   = azurerm_resource_group.test.name
  private_dns_zone_name = azurerm_private_dns_zone.test.name
}
`, PrivateDnsZoneVirtualNetworkLinkResource{}.resolutionPolicy(data, "NxDomainRedirect"))
}
