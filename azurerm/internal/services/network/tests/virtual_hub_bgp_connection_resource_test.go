package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type VirtualHubBGPConnectionResource struct {
}

func TestAccVirtualHubBgpConnection_basic(t *testing.T) {
	if true {
		t.Skip("Skipping due to API issue preventing deletion")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_bgp_connection", "test")
	r := VirtualHubBGPConnectionResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubBgpConnection_requiresImport(t *testing.T) {
	if true {
		t.Skip("Skipping due to API issue preventing deletion")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_bgp_connection", "test")
	r := VirtualHubBGPConnectionResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t VirtualHubBGPConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.BgpConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.VirtualHubBgpConnectionClient.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Virtual Hub BGP Connectionn (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
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
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNet-%d"
  address_space       = ["10.5.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-Subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.5.1.0/24"
}

resource "azurerm_virtual_hub_ip" "test" {
  name                         = "acctest-VHub-IP-%d"
  virtual_hub_id               = azurerm_virtual_hub.test.id
  private_ip_address           = "10.5.1.18"
  private_ip_allocation_method = "Static"
  public_ip_address_id         = azurerm_public_ip.test.id
  subnet_id                    = azurerm_subnet.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
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
