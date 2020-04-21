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

func TestAccAzureRMSubnetRouteTableAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_route_table_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional since this is a Virtual Resource
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetRouteTableAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetRouteTableAssociationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnetRouteTableAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_route_table_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional since this is a Virtual Resource
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetRouteTableAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetRouteTableAssociationExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMSubnetRouteTableAssociation_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_subnet_route_table_association"),
			},
		},
	})
}

func TestAccAzureRMSubnetRouteTableAssociation_updateSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_route_table_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional since this is a Virtual Resource
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetRouteTableAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetRouteTableAssociationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnetRouteTableAssociation_updateSubnet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetRouteTableAssociationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnetRouteTableAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_route_table_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional since this is a Virtual Resource
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetRouteTableAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetRouteTableAssociationExists(data.ResourceName),
					testCheckAzureRMSubnetRouteTableAssociationDisappears(data.ResourceName),
					testCheckAzureRMSubnetHasNoRouteTable(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMSubnetRouteTableAssociationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SubnetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		subnetId := rs.Primary.Attributes["subnet_id"]
		parsedId, err := azure.ParseAzureResourceID(subnetId)
		if err != nil {
			return err
		}

		resourceGroupName := parsedId.ResourceGroup
		virtualNetworkName := parsedId.Path["virtualNetworks"]
		subnetName := parsedId.Path["subnets"]

		resp, err := client.Get(ctx, resourceGroupName, virtualNetworkName, subnetName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Subnet %q (Virtual Network %q / Resource Group: %q) does not exist", subnetName, virtualNetworkName, resourceGroupName)
			}

			return fmt.Errorf("Bad: Get on subnetClient: %+v", err)
		}

		props := resp.SubnetPropertiesFormat
		if props == nil {
			return fmt.Errorf("Properties was nil for Subnet %q (Virtual Network %q / Resource Group: %q)", subnetName, virtualNetworkName, resourceGroupName)
		}

		if props.RouteTable == nil || props.RouteTable.ID == nil {
			return fmt.Errorf("No Route Table association exists for Subnet %q (Virtual Network %q / Resource Group: %q)", subnetName, virtualNetworkName, resourceGroupName)
		}

		return nil
	}
}

func testCheckAzureRMSubnetRouteTableAssociationDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SubnetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		subnetId := rs.Primary.Attributes["subnet_id"]
		parsedId, err := azure.ParseAzureResourceID(subnetId)
		if err != nil {
			return err
		}

		resourceGroup := parsedId.ResourceGroup
		virtualNetworkName := parsedId.Path["virtualNetworks"]
		subnetName := parsedId.Path["subnets"]

		read, err := client.Get(ctx, resourceGroup, virtualNetworkName, subnetName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(read.Response) {
				return fmt.Errorf("Error retrieving Subnet %q (Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
			}
		}

		read.SubnetPropertiesFormat.RouteTable = nil

		future, err := client.CreateOrUpdate(ctx, resourceGroup, virtualNetworkName, subnetName, read)
		if err != nil {
			return fmt.Errorf("Error updating Subnet %q (Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for completion of Subnet %q (Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMSubnetHasNoRouteTable(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SubnetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		subnetId := rs.Primary.Attributes["subnet_id"]
		parsedId, err := azure.ParseAzureResourceID(subnetId)
		if err != nil {
			return err
		}

		resourceGroupName := parsedId.ResourceGroup
		virtualNetworkName := parsedId.Path["virtualNetworks"]
		subnetName := parsedId.Path["subnets"]

		resp, err := client.Get(ctx, resourceGroupName, virtualNetworkName, subnetName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Subnet %q (Virtual Network %q / Resource Group: %q) does not exist", subnetName, virtualNetworkName, resourceGroupName)
			}

			return fmt.Errorf("Bad: Get on subnetClient: %+v", err)
		}

		props := resp.SubnetPropertiesFormat
		if props == nil {
			return fmt.Errorf("Properties was nil for Subnet %q (Virtual Network %q / Resource Group: %q)", subnetName, virtualNetworkName, resourceGroupName)
		}

		if props.RouteTable != nil && ((props.RouteTable.ID == nil) || (props.RouteTable.ID != nil && *props.RouteTable.ID == "")) {
			return fmt.Errorf("No Route Table should exist for Subnet %q (Virtual Network %q / Resource Group: %q) but got %q", subnetName, virtualNetworkName, resourceGroupName, *props.RouteTable.ID)
		}

		return nil
	}
}

func testAccAzureRMSubnetRouteTableAssociation_basic(data acceptance.TestData) string {
	template := testAccAzureRMSubnetRouteTableAssociation_template(data)
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
`, template)
}

func testAccAzureRMSubnetRouteTableAssociation_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMSubnetRouteTableAssociation_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_route_table_association" "import" {
  subnet_id      = azurerm_subnet_route_table_association.test.subnet_id
  route_table_id = azurerm_subnet_route_table_association.test.route_table_id
}
`, template)
}

func testAccAzureRMSubnetRouteTableAssociation_updateSubnet(data acceptance.TestData) string {
	template := testAccAzureRMSubnetRouteTableAssociation_template(data)
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
`, template)
}

func testAccAzureRMSubnetRouteTableAssociation_template(data acceptance.TestData) string {
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
