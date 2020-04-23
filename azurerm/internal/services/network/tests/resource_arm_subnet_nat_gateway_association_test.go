package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSubnetNatGatewayAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_nat_gateway_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional since this is a virtual resource
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetNatGatewayAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetNatGatewayAssociationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnetNatGatewayAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_nat_gateway_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional since this is a virtual resource
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetNatGatewayAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetNatGatewayAssociationExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMSubnetNatGatewayAssociation_requiresImport(data),
				ExpectError: acceptance.RequiresImportError(data.ResourceType),
			},
		},
	})
}

func TestAccAzureRMSubnetNatGatewayAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_nat_gateway_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional since this is virtual resource
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetNatGatewayAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetNatGatewayAssociationExists(data.ResourceName),
					testCheckAzureRMSubnetNatGatewayAssociationDisappears(data.ResourceName),
					testCheckAzureRMSubnetHasNoNatGateways(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMSubnetNatGatewayAssociation_updateSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_nat_gateway_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional since this is a virtual resource
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetNatGatewayAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetNatGatewayAssociationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnetNatGatewayAssociation_updateSubnet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetNatGatewayAssociationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMSubnetNatGatewayAssociationExists(resourceName string) resource.TestCheckFunc {
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

		if props.NatGateway == nil || props.NatGateway.ID == nil {
			return fmt.Errorf("No NAT Gateway association exists for Subnet %q (Virtual Network %q / Resource Group: %q)", subnetName, virtualNetworkName, resourceGroupName)
		}

		return nil
	}
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

func testAccAzureRMSubnetNatGatewayAssociation_basic(data acceptance.TestData) string {
	template := testAccAzureRMSubnetNatGatewayAssociation_template(data)
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
`, template)
}

func testAccAzureRMSubnetNatGatewayAssociation_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMSubnetNatGatewayAssociation_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_nat_gateway_association" "import" {
  subnet_id      = azurerm_subnet_nat_gateway_association.test.subnet_id
  nat_gateway_id = azurerm_subnet_nat_gateway_association.test.nat_gateway_id
}
`, template)
}

func testAccAzureRMSubnetNatGatewayAssociation_updateSubnet(data acceptance.TestData) string {
	template := testAccAzureRMSubnetNatGatewayAssociation_template(data)
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
`, template)
}

func testAccAzureRMSubnetNatGatewayAssociation_template(data acceptance.TestData) string {
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
