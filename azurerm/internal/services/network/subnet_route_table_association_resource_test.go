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

type SubnetRouteTableAssociationResource struct {
}

func TestAccSubnetRouteTableAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_route_table_association", "test")
	r := SubnetRouteTableAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		// intentional since this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnetRouteTableAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_route_table_association", "test")
	r := SubnetRouteTableAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		// intentional since this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		// intentional since this is a Virtual Resource
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateSubnet(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubnetRouteTableAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_route_table_association", "test")
	r := SubnetRouteTableAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// NOTE: intentionally not using a DisappearsStep as this is a Virtual Resource
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.destroy),
				data.CheckWithClientForResource(SubnetResource{}.hasNoRouteTable, "azurerm_subnet.test"),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func (SubnetRouteTableAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.SubnetID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	virtualNetworkName := id.VirtualNetworkName
	subnetName := id.Name

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

func (SubnetRouteTableAssociationResource) destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) error {
	parsedId, err := parse.SubnetID(state.Attributes["subnet_id"])
	if err != nil {
		return err
	}

	resourceGroup := parsedId.ResourceGroup
	virtualNetworkName := parsedId.VirtualNetworkName
	subnetName := parsedId.Name

	read, err := client.Network.SubnetsClient.Get(ctx, resourceGroup, virtualNetworkName, subnetName, "")
	if err != nil {
		if !utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf("Error retrieving Subnet %q (Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
		}
	}

	read.SubnetPropertiesFormat.RouteTable = nil

	future, err := client.Network.SubnetsClient.CreateOrUpdate(ctx, resourceGroup, virtualNetworkName, subnetName, read)
	if err != nil {
		return fmt.Errorf("Error updating Subnet %q (Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Network.SubnetsClient.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Subnet %q (Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
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
  address_prefix       = "10.0.2.0/24"
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
  address_prefix       = "10.0.2.0/24"

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
