package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SubnetNatGatewayAssociationResource struct {
}

func TestAccAzureRMSubnetNatGatewayAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_nat_gateway_association", "test")
	r := SubnetNatGatewayAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		// intentional since this is a virtual resource
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMSubnetNatGatewayAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_nat_gateway_association", "test")
	r := SubnetNatGatewayAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		// intentional since this is a virtual resource
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func TestAccAzureRMSubnetNatGatewayAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_nat_gateway_association", "test")
	r := SubnetNatGatewayAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		// intentional since this is virtual resource
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				testCheckAzureRMSubnetNatGatewayAssociationDisappears(data.ResourceName),
				testCheckAzureRMSubnetHasNoNatGateways(data.ResourceName),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccAzureRMSubnetNatGatewayAssociation_updateSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_nat_gateway_association", "test")
	r := SubnetNatGatewayAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		// intentional since this is a virtual resource
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

func (t SubnetNatGatewayAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	virtualNetworkName := id.Path["virtualNetworks"]
	subnetName := id.Path["subnets"]

	resp, err := clients.Network.SubnetsClient.Get(ctx, resourceGroup, virtualNetworkName, subnetName, "")
	if err != nil {
		return nil, fmt.Errorf("reading Subnet Nat Gateway Association (%s): %+v", id, err)
	}

	props := resp.SubnetPropertiesFormat
	if props == nil || props.NatGateway == nil {
		return nil, fmt.Errorf("properties was nil for Subnet %q (Virtual Network %q / Resource Group: %q)", subnetName, virtualNetworkName, resourceGroup)
	}

	return utils.Bool(props.NatGateway.ID != nil), nil
}

func testCheckAzureRMSubnetNatGatewayAssociationDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SubnetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		subnetId := rs.Primary.Attributes["subnet_id"]
		parsedSubnetId, err := azure.ParseAzureResourceID(subnetId)
		if err != nil {
			return err
		}

		resourceGroup := parsedSubnetId.ResourceGroup
		virtualNetworkName := parsedSubnetId.Path["virtualNetworks"]
		subnetName := parsedSubnetId.Path["subnets"]

		subnet, err := client.Get(ctx, resourceGroup, virtualNetworkName, subnetName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(subnet.Response) {
				return fmt.Errorf("Error retrieving Subnet %q (Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
			}
			return fmt.Errorf("Bad: Get on subnetClient: %+v", err)
		}

		props := subnet.SubnetPropertiesFormat
		if props == nil {
			return fmt.Errorf("Properties was nil for Subnet %q (Virtual Network %q / Resource Group: %q)", subnetName, virtualNetworkName, resourceGroup)
		}
		props.NatGateway = nil

		future, err := client.CreateOrUpdate(ctx, resourceGroup, virtualNetworkName, subnetName, subnet)
		if err != nil {
			return fmt.Errorf("Error updating Subnet %q (Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for completion of Subnet %q (Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
		}
		return nil
	}
}

func testCheckAzureRMSubnetHasNoNatGateways(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SubnetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		subnetId := rs.Primary.Attributes["subnet_id"]
		parsedSubnetId, err := azure.ParseAzureResourceID(subnetId)
		if err != nil {
			return err
		}
		resourceGroupName := parsedSubnetId.ResourceGroup
		virtualNetworkName := parsedSubnetId.Path["virtualNetworks"]
		subnetName := parsedSubnetId.Path["subnets"]

		subnet, err := client.Get(ctx, resourceGroupName, virtualNetworkName, subnetName, "")
		if err != nil {
			if utils.ResponseWasNotFound(subnet.Response) {
				return fmt.Errorf("Bad: Subnet %q (Virtual Network %q / Resource Group: %q) does not exist", subnetName, virtualNetworkName, resourceGroupName)
			}
			return fmt.Errorf("Bad: Get on subnetClient: %+v", err)
		}

		props := subnet.SubnetPropertiesFormat
		if props == nil {
			return fmt.Errorf("Properties was nil for Subnet %q (Virtual Network %q / Resource Group: %q)", subnetName, virtualNetworkName, resourceGroupName)
		}

		if props.NatGateway != nil && ((props.NatGateway.ID == nil) || (props.NatGateway.ID != nil && *props.NatGateway.ID == "")) {
			return fmt.Errorf("No Route Table should exist for Subnet %q (Virtual Network %q / Resource Group: %q) but got %q", subnetName, virtualNetworkName, resourceGroupName, *props.RouteTable.ID)
		}
		return nil
	}
}

func (r SubnetNatGatewayAssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_subnet_nat_gateway_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  nat_gateway_id = azurerm_nat_gateway.test.id
}
`, r.template(data))
}

func (r SubnetNatGatewayAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_nat_gateway_association" "import" {
  subnet_id      = azurerm_subnet_nat_gateway_association.test.subnet_id
  nat_gateway_id = azurerm_subnet_nat_gateway_association.test.nat_gateway_id
}
`, r.basic(data))
}

func (r SubnetNatGatewayAssociationResource) updateSubnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_subnet_nat_gateway_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  nat_gateway_id = azurerm_nat_gateway.test.id
}
`, r.template(data))
}

func (SubnetNatGatewayAssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_nat_gateway" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
