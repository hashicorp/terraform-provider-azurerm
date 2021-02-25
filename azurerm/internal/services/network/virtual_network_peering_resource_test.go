package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

type VirtualNetworkPeeringResource struct {
}

func TestAccVirtualNetworkPeering_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_peering", "test1")
	r := VirtualNetworkPeeringResource{}
	secondResourceName := "azurerm_virtual_network_peering.test2"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allow_virtual_network_access").HasValue("true"),
				resource.TestCheckResourceAttr(secondResourceName, "allow_virtual_network_access", "true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualNetworkPeering_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_peering", "test1")
	r := VirtualNetworkPeeringResource{}
	secondResourceName := "azurerm_virtual_network_peering.test2"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccVirtualNetworkPeering_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_peering", "test1")
	r := VirtualNetworkPeeringResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccVirtualNetworkPeering_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_peering", "test1")
	r := VirtualNetworkPeeringResource{}
	secondResourceName := "azurerm_virtual_network_peering.test2"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allow_virtual_network_access").HasValue("true"),
				resource.TestCheckResourceAttr(secondResourceName, "allow_virtual_network_access", "true"),
				check.That(data.ResourceName).Key("allow_forwarded_traffic").HasValue("false"),
				resource.TestCheckResourceAttr(secondResourceName, "allow_forwarded_traffic", "false"),
			),
		},

		{
			Config: r.basicUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allow_virtual_network_access").HasValue("true"),
				resource.TestCheckResourceAttr(secondResourceName, "allow_virtual_network_access", "true"),
				check.That(data.ResourceName).Key("allow_forwarded_traffic").HasValue("true"),
				resource.TestCheckResourceAttr(secondResourceName, "allow_forwarded_traffic", "true"),
			),
		},
	})
}

func (t VirtualNetworkPeeringResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resGroup := id.ResourceGroup
	vnetName := id.Path["virtualNetworks"]
	name := id.Path["virtualNetworkPeerings"]

	resp, err := clients.Network.VnetPeeringsClient.Get(ctx, resGroup, vnetName, name)
	if err != nil {
		return nil, fmt.Errorf("reading Virtual Network Peering (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r VirtualNetworkPeeringResource) Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resGroup := id.ResourceGroup
	vnetName := id.Path["virtualNetworks"]
	name := id.Path["virtualNetworkPeerings"]

	future, err := client.Network.VnetPeeringsClient.Delete(ctx, resGroup, vnetName, name)
	if err != nil {
		return nil, fmt.Errorf("deleting on virtual network peering: %+v", err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Network.VnetPeeringsClient.Client); err != nil {
		return nil, fmt.Errorf("waiting for deletion of Peering %q: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (VirtualNetworkPeeringResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test1" {
  name                = "acctestvirtnet-1-%d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.1.0/24"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvirtnet-2-%d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.2.0/24"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network_peering" "test1" {
  name                         = "acctestpeer-1-%d"
  resource_group_name          = azurerm_resource_group.test.name
  virtual_network_name         = azurerm_virtual_network.test1.name
  remote_virtual_network_id    = azurerm_virtual_network.test2.id
  allow_virtual_network_access = true
}

resource "azurerm_virtual_network_peering" "test2" {
  name                         = "acctestpeer-2-%d"
  resource_group_name          = azurerm_resource_group.test.name
  virtual_network_name         = azurerm_virtual_network.test2.name
  remote_virtual_network_id    = azurerm_virtual_network.test1.id
  allow_virtual_network_access = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VirtualNetworkPeeringResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network_peering" "import" {
  name                         = azurerm_virtual_network_peering.test1.name
  resource_group_name          = azurerm_virtual_network_peering.test1.resource_group_name
  virtual_network_name         = azurerm_virtual_network_peering.test1.virtual_network_name
  remote_virtual_network_id    = azurerm_virtual_network_peering.test1.remote_virtual_network_id
  allow_virtual_network_access = azurerm_virtual_network_peering.test1.allow_virtual_network_access
}
`, r.basic(data))
}

func (VirtualNetworkPeeringResource) basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test1" {
  name                = "acctestvirtnet-1-%d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.1.0/24"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvirtnet-2-%d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.2.0/24"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network_peering" "test1" {
  name                         = "acctestpeer-1-%d"
  resource_group_name          = azurerm_resource_group.test.name
  virtual_network_name         = azurerm_virtual_network.test1.name
  remote_virtual_network_id    = azurerm_virtual_network.test2.id
  allow_forwarded_traffic      = true
  allow_virtual_network_access = true
}

resource "azurerm_virtual_network_peering" "test2" {
  name                         = "acctestpeer-2-%d"
  resource_group_name          = azurerm_resource_group.test.name
  virtual_network_name         = azurerm_virtual_network.test2.name
  remote_virtual_network_id    = azurerm_virtual_network.test1.id
  allow_forwarded_traffic      = true
  allow_virtual_network_access = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
