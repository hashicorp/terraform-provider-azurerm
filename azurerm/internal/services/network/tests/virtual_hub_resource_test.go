package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type VirtualHubResource struct {
}

func TestAccVirtualHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub", "test")
	r := VirtualHubResource{}

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

func TestAccVirtualHub_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub", "test")
	r := VirtualHubResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_virtual_hub"),
		},
	})
}

func TestAccVirtualHub_routes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub", "test")
	r := VirtualHubResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.route(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.routeUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHub_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub", "test")
	r := VirtualHubResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.tags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t VirtualHubResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.VirtualHubID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.VirtualHubClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Virtual Hub (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r VirtualHubResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "test" {
  name                = "acctestVHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualHubResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "import" {
  name                = azurerm_virtual_hub.test.name
  location            = azurerm_virtual_hub.test.location
  resource_group_name = azurerm_virtual_hub.test.resource_group_name
  virtual_wan_id      = azurerm_virtual_hub.test.virtual_wan_id
  address_prefix      = azurerm_virtual_hub.test.address_prefix
}
`, r.basic(data))
}

func (r VirtualHubResource) route(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "test" {
  name                = "acctestVHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"

  route {
    address_prefixes    = ["172.0.1.0/24"]
    next_hop_ip_address = "12.34.56.78"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualHubResource) routeUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "test" {
  name                = "acctestVHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"

  route {
    address_prefixes    = ["172.0.1.0/24"]
    next_hop_ip_address = "87.65.43.21"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualHubResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "test" {
  name                = "acctestVHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"

  tags = {
    Hello = "World"
  }
}
`, r.template(data), data.RandomInteger)
}

func (VirtualHubResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
