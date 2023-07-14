// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SubnetRouteTableAssociationResource struct{}

func TestAccSubnetRouteTableAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_route_table_association", "test")
	r := SubnetRouteTableAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional since this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnetRouteTableAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_route_table_association", "test")
	r := SubnetRouteTableAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional since this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_subnet_route_table_association"),
		},
	})
}

func TestAccSubnetRouteTableAssociation_updateSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_route_table_association", "test")
	r := SubnetRouteTableAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		// intentional since this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateSubnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnetRouteTableAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_route_table_association", "test")
	r := SubnetRouteTableAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// NOTE: intentionally not using a DisappearsStep as this is a Virtual Resource
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.destroy),
				data.CheckWithClientForResource(SubnetResource{}.hasNoRouteTable, "azurerm_subnet.test"),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func (SubnetRouteTableAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseSubnetID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroupName
	virtualNetworkName := id.VirtualNetworkName
	subnetName := id.SubnetName

	resp, err := clients.Network.SubnetsClient.Get(ctx, resourceGroup, virtualNetworkName, subnetName, "")
	if err != nil {
		return nil, fmt.Errorf("reading Subnet Route Table Association (%s): %+v", id, err)
	}

	props := resp.SubnetPropertiesFormat
	if props == nil || props.RouteTable == nil {
		return nil, fmt.Errorf("properties was nil for Subnet %q (Virtual Network %q / Resource Group: %q)", subnetName, virtualNetworkName, resourceGroup)
	}

	return utils.Bool(props.RouteTable.ID != nil), nil
}

func (SubnetRouteTableAssociationResource) destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	parsedId, err := commonids.ParseSubnetID(state.Attributes["subnet_id"])
	if err != nil {
		return err
	}

	resourceGroup := parsedId.ResourceGroupName
	virtualNetworkName := parsedId.VirtualNetworkName
	subnetName := parsedId.SubnetName

	read, err := client.Network.SubnetsClient.Get(ctx, resourceGroup, virtualNetworkName, subnetName, "")
	if err != nil {
		if !utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf("retrieving Subnet %q (Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
		}
	}

	read.SubnetPropertiesFormat.RouteTable = nil

	future, err := client.Network.SubnetsClient.CreateOrUpdate(ctx, resourceGroup, virtualNetworkName, subnetName, read)
	if err != nil {
		return fmt.Errorf("updating Subnet %q (Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Network.SubnetsClient.Client); err != nil {
		return fmt.Errorf("waiting for completion of Subnet %q (Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
	}

	return nil
}

func (r SubnetRouteTableAssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  route_table_id = azurerm_route_table.test.id
}
`, r.template(data))
}

func (r SubnetRouteTableAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_route_table_association" "import" {
  subnet_id      = azurerm_subnet_route_table_association.test.subnet_id
  route_table_id = azurerm_subnet_route_table_association.test.route_table_id
}
`, r.basic(data))
}

func (r SubnetRouteTableAssociationResource) updateSubnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  route_table_id = azurerm_route_table.test.id
}
`, r.template(data))
}

func (SubnetRouteTableAssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_route_table" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  route {
    name                   = "first"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
