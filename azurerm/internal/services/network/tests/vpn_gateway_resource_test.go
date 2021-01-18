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

type VPNGatewayResource struct {
}

func TestAccVPNGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")
	r := VPNGatewayResource{}

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

func TestAccVPNGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")
	r := VPNGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_vpn_gateway"),
		},
	})
}

func TestAccVPNGateway_bgpSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")
	r := VPNGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.bgpSettings(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVPNGateway_scaleUnit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")
	r := VPNGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.scaleUnit(data, 2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.scaleUnit(data, 3),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVPNGateway_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")
	r := VPNGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.tags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.tagsUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.tags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t VPNGatewayResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.VpnGatewayID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.VpnGatewaysClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading VPN Gateway (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r VPNGatewayResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r VPNGatewayResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "import" {
  name                = azurerm_vpn_gateway.test.name
  location            = azurerm_vpn_gateway.test.location
  resource_group_name = azurerm_vpn_gateway.test.resource_group_name
  virtual_hub_id      = azurerm_vpn_gateway.test.virtual_hub_id
}
`, r.basic(data))
}

func (r VPNGatewayResource) bgpSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id

  bgp_settings {
    asn         = 65515
    peer_weight = 0

    instance_0_bgp_peering_address {
      custom_ips = ["169.254.21.5"]
    }

    instance_1_bgp_peering_address {
      custom_ips = ["169.254.21.10"]
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VPNGatewayResource) scaleUnit(data acceptance.TestData, scaleUnit int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id
  scale_unit          = %d
}
`, r.template(data), data.RandomInteger, scaleUnit)
}

func (r VPNGatewayResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id

  tags = {
    Hello = "World"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VPNGatewayResource) tagsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id

  tags = {
    Hello = "World"
    Rick  = "C-137"
  }
}
`, r.template(data), data.RandomInteger)
}

func (VPNGatewayResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctestvh-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_prefix      = "10.0.1.0/24"
  virtual_wan_id      = azurerm_virtual_wan.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
