// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VirtualHubBGPConnectionResource struct{}

func TestAccVirtualHubBgpConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_bgp_connection", "test")
	r := VirtualHubBGPConnectionResource{}
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

func TestAccVirtualHubBgpConnection_virtualWan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_bgp_connection", "test")
	r := VirtualHubBGPConnectionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualWan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.virtualWan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubBgpConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_bgp_connection", "test")
	r := VirtualHubBGPConnectionResource{}
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

func (t VirtualHubBGPConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseVirtualHubBGPConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.VirtualWANs.VirtualHubBgpConnectionGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (VirtualHubBGPConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-VHub-%d"
  location = "%s"
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-VHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-PIP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNet-%d"
  address_space       = ["10.5.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "RouteServerSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]
}

resource "azurerm_virtual_hub_ip" "test" {
  name                         = "acctest-VHub-IP-%d"
  virtual_hub_id               = azurerm_virtual_hub.test.id
  private_ip_address           = "10.5.1.18"
  private_ip_allocation_method = "Static"
  public_ip_address_id         = azurerm_public_ip.test.id
  subnet_id                    = azurerm_subnet.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VirtualHubBGPConnectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_bgp_connection" "test" {
  name           = "acctest-VHub-BgpConnection-%d"
  virtual_hub_id = azurerm_virtual_hub.test.id
  peer_asn       = 65514
  peer_ip        = "169.254.21.5"

  depends_on = [azurerm_virtual_hub_ip.test]
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualHubBGPConnectionResource) virtualWanTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.5.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-VHUB-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.2.0/24"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r VirtualHubBGPConnectionResource) virtualWan(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctestbasicvhubconn-%[2]d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
}

resource "azurerm_virtual_hub_bgp_connection" "test" {
  name                          = "acctest-VHub-BgpConnection-%[2]d"
  virtual_hub_id                = azurerm_virtual_hub.test.id
  peer_asn                      = 65514
  peer_ip                       = "10.5.0.1"
  virtual_network_connection_id = azurerm_virtual_hub_connection.test.id
}
`, r.virtualWanTemplate(data), data.RandomInteger)
}

func (r VirtualHubBGPConnectionResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctestbasicvhubconn-%[2]d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvirtnet2%[2]d"
  address_space       = ["10.6.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test2" {
  name                = "acctestnsg2%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test2" {
  name                 = "acctestsubnet2%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test2.name
  address_prefixes     = ["10.6.1.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "test2" {
  subnet_id                 = azurerm_subnet.test2.id
  network_security_group_id = azurerm_network_security_group.test2.id
}

resource "azurerm_virtual_hub" "test2" {
  name                = "acctest-VHUB2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.1.2.0/24"
}

resource "azurerm_virtual_hub_connection" "test2" {
  name                      = "acctestbasicvhubconn2-%[2]d"
  virtual_hub_id            = azurerm_virtual_hub.test2.id
  remote_virtual_network_id = azurerm_virtual_network.test2.id
}

resource "azurerm_virtual_hub_bgp_connection" "test" {
  name                          = "acctest-VHub-BgpConnection-%[2]d"
  virtual_hub_id                = azurerm_virtual_hub.test2.id
  peer_asn                      = 65514
  peer_ip                       = "10.6.0.1"
  virtual_network_connection_id = azurerm_virtual_hub_connection.test2.id
}
`, r.virtualWanTemplate(data), data.RandomInteger)
}

func (r VirtualHubBGPConnectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_bgp_connection" "import" {
  name           = azurerm_virtual_hub_bgp_connection.test.name
  virtual_hub_id = azurerm_virtual_hub_bgp_connection.test.virtual_hub_id
  peer_asn       = azurerm_virtual_hub_bgp_connection.test.peer_asn
  peer_ip        = azurerm_virtual_hub_bgp_connection.test.peer_ip
}
`, r.basic(data))
}
