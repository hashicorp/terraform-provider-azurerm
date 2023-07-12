// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatedns_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2020-06-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PrivateDnsZoneVirtualNetworkLinkResource struct{}

func TestAccPrivateDnsZoneVirtualNetworkLink_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_zone_virtual_network_link", "test")
	r := PrivateDnsZoneVirtualNetworkLinkResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateDnsZoneVirtualNetworkLink_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_zone_virtual_network_link", "test")
	r := PrivateDnsZoneVirtualNetworkLinkResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateDnsZoneVirtualNetworkLink_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_zone_virtual_network_link", "test")
	r := PrivateDnsZoneVirtualNetworkLinkResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccPrivateDnsZoneVirtualNetworkLink_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_zone_virtual_network_link", "test")
	r := PrivateDnsZoneVirtualNetworkLinkResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		{
			Config: r.withTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func (t PrivateDnsZoneVirtualNetworkLinkResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := virtualnetworklinks.ParseVirtualNetworkLinkID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.PrivateDns.VirtualNetworkLinksClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Private DNS Zone Virtual Network Link (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (PrivateDnsZoneVirtualNetworkLinkResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = "acctestVnetZone%d.com"
  private_dns_zone_name = azurerm_private_dns_zone.test.name
  virtual_network_id    = azurerm_virtual_network.test.id
  resource_group_name   = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (PrivateDnsZoneVirtualNetworkLinkResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = "acctestVnetZone%d.com"
  private_dns_zone_name = azurerm_private_dns_zone.test.name
  virtual_network_id    = azurerm_virtual_network.test.id
  resource_group_name   = azurerm_resource_group.test.name
  registration_enabled  = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r PrivateDnsZoneVirtualNetworkLinkResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_zone_virtual_network_link" "import" {
  name                  = azurerm_private_dns_zone_virtual_network_link.test.name
  private_dns_zone_name = azurerm_private_dns_zone_virtual_network_link.test.private_dns_zone_name
  virtual_network_id    = azurerm_private_dns_zone_virtual_network_link.test.virtual_network_id
  resource_group_name   = azurerm_private_dns_zone_virtual_network_link.test.resource_group_name
}
`, r.basic(data))
}

func (PrivateDnsZoneVirtualNetworkLinkResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = "acctestVnetZone%d.com"
  private_dns_zone_name = azurerm_private_dns_zone.test.name
  virtual_network_id    = azurerm_virtual_network.test.id
  resource_group_name   = azurerm_resource_group.test.name

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (PrivateDnsZoneVirtualNetworkLinkResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = "acctestVnetZone%d.com"
  private_dns_zone_name = azurerm_private_dns_zone.test.name
  virtual_network_id    = azurerm_virtual_network.test.id
  resource_group_name   = azurerm_resource_group.test.name

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
