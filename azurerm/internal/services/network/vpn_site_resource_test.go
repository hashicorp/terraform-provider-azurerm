package network_test

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

type VPNSiteResource struct {
}

func TestAccVpnSite_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_site", "test")
	r := VPNSiteResource{}

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

func TestAccVpnSite_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_site", "test")
	r := VPNSiteResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVpnSite_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_site", "test")
	r := VPNSiteResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVpnSite_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_site", "test")
	r := VPNSiteResource{}

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

func (t VPNSiteResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.VpnSiteID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.VpnSitesClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading VPN Site (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r VPNSiteResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_site" "test" {
  name                = "acctest-VpnSite-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_wan_id      = azurerm_virtual_wan.test.id
  link {
    name       = "link1"
    ip_address = "10.0.0.1"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VPNSiteResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_site" "test" {
  name                = "acctest-VpnSite-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_cidrs       = ["10.0.0.0/24", "10.0.1.0/24"]

  device_vendor = "Cisco"
  device_model  = "foobar"

  link {
    name          = "link1"
    provider_name = "Verizon"
    speed_in_mbps = 50
    ip_address    = "10.0.0.1"
    bgp {
      asn             = 12345
      peering_address = "10.0.0.1"
    }
  }

  link {
    name = "link2"
    fqdn = "foo.com"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VPNSiteResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_site" "import" {
  name                = "acctest-VpnSite-%d"
  location            = azurerm_vpn_site.test.location
  resource_group_name = azurerm_vpn_site.test.resource_group_name
  virtual_wan_id      = azurerm_vpn_site.test.virtual_wan_id
  link {
    name       = "link1"
    ip_address = "10.0.0.1"
  }
}
`, r.basic(data), data.RandomInteger)
}

func (VPNSiteResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-rg-%d"
  location = "%s"
}


resource "azurerm_virtual_wan" "test" {
  name                = "acctest-vwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
