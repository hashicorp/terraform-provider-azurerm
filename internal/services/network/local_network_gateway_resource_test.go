// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/localnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LocalNetworkGatewayResource struct{}

func TestAccLocalNetworkGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gateway_address").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("address_space.0").HasValue("127.0.0.0/8"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLocalNetworkGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_local_network_gateway"),
		},
	})
}

func TestAccLocalNetworkGateway_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccLocalNetworkGateway_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("acctest"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLocalNetworkGateway_bgpSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.bgpSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gateway_address").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("address_space.0").HasValue("127.0.0.0/8"),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLocalNetworkGateway_bgpSettingsDisable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.bgpSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gateway_address").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("address_space.0").HasValue("127.0.0.0/8"),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.asn").HasValue("2468"),
				check.That(data.ResourceName).Key("bgp_settings.0.bgp_peering_address").HasValue("10.104.1.1"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gateway_address").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("address_space.0").HasValue("127.0.0.0/8"),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("0"),
			),
		},
	})
}

func TestAccLocalNetworkGateway_bgpSettingsEnable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gateway_address").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("address_space.0").HasValue("127.0.0.0/8"),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("0"),
			),
		},
		{
			Config: r.bgpSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gateway_address").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("address_space.0").HasValue("127.0.0.0/8"),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.asn").HasValue("2468"),
				check.That(data.ResourceName).Key("bgp_settings.0.bgp_peering_address").HasValue("10.104.1.1"),
			),
		},
	})
}

func TestAccLocalNetworkGateway_bgpSettingsComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.bgpSettingsComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gateway_address").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("address_space.0").HasValue("127.0.0.0/8"),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.asn").HasValue("2468"),
				check.That(data.ResourceName).Key("bgp_settings.0.bgp_peering_address").HasValue("10.104.1.1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peer_weight").HasValue("15"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLocalNetworkGateway_updateAddressSpace(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.singleAddressSpace(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleAddressSpace(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleAddressSpaceUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLocalNetworkGateway_fqdn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fqdn(data, "www.foo.com"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.fqdn(data, "www.bar.com"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t LocalNetworkGatewayResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := localnetworkgateways.ParseLocalNetworkGatewayID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.LocalNetworkGateways.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (LocalNetworkGatewayResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := localnetworkgateways.ParseLocalNetworkGatewayID(state.ID)
	if err != nil {
		return nil, err
	}

	ctx2, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	if err := client.Network.Client.LocalNetworkGateways.DeleteThenPoll(ctx2, *id); err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return pointer.To(true), nil
}

func (LocalNetworkGatewayResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-lngw-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r LocalNetworkGatewayResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_local_network_gateway" "import" {
  name                = azurerm_local_network_gateway.test.name
  location            = azurerm_local_network_gateway.test.location
  resource_group_name = azurerm_local_network_gateway.test.resource_group_name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]
}
`, r.basic(data))
}

func (LocalNetworkGatewayResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-lngw-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]

  tags = {
    environment = "acctest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LocalNetworkGatewayResource) bgpSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-lngw-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]

  bgp_settings {
    asn                 = 2468
    bgp_peering_address = "10.104.1.1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LocalNetworkGatewayResource) bgpSettingsComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-lngw-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]

  bgp_settings {
    asn                 = 2468
    bgp_peering_address = "10.104.1.1"
    peer_weight         = 15
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LocalNetworkGatewayResource) singleAddressSpace(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-lngw-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/24"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LocalNetworkGatewayResource) multipleAddressSpace(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-lngw-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/24", "127.0.1.0/24"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LocalNetworkGatewayResource) multipleAddressSpaceUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-lngw-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.1.0/24", "127.0.0.0/24"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LocalNetworkGatewayResource) fqdn(data acceptance.TestData, fqdn string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_fqdn        = %q
  address_space       = ["127.0.0.0/8"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, fqdn)
}
