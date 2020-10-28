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

func TestAccAzureRMSubnetNetworkSecurityGroupAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_network_security_group_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetNetworkSecurityGroupAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetNetworkSecurityGroupAssociationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnetNetworkSecurityGroupAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_network_security_group_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetNetworkSecurityGroupAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetNetworkSecurityGroupAssociationExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMSubnetNetworkSecurityGroupAssociation_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_subnet_network_security_group_association"),
			},
		},
	})
}

func TestAccAzureRMSubnetNetworkSecurityGroupAssociation_updateSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_network_security_group_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetNetworkSecurityGroupAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetNetworkSecurityGroupAssociationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnetNetworkSecurityGroupAssociation_updateSubnet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetNetworkSecurityGroupAssociationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnetNetworkSecurityGroupAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_network_security_group_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetNetworkSecurityGroupAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetNetworkSecurityGroupAssociationExists(data.ResourceName),
					testCheckAzureRMSubnetNetworkSecurityGroupAssociationDisappears(data.ResourceName),
					testCheckAzureRMSubnetHasNoNetworkSecurityGroup(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMSubnetNetworkSecurityGroupAssociationExists(resourceName string) resource.TestCheckFunc {
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

		if props.NetworkSecurityGroup == nil || props.NetworkSecurityGroup.ID == nil {
			return fmt.Errorf("No Network Security Group association exists for Subnet %q (Virtual Network %q / Resource Group: %q)", subnetName, virtualNetworkName, resourceGroupName)
		}

		return nil
	}
}

func testCheckAzureRMSubnetNetworkSecurityGroupAssociationDisappears(resourceName string) resource.TestCheckFunc {
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

		read.SubnetPropertiesFormat.NetworkSecurityGroup = nil

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

func testCheckAzureRMSubnetHasNoNetworkSecurityGroup(resourceName string) resource.TestCheckFunc {
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

		if props.NetworkSecurityGroup != nil && ((props.NetworkSecurityGroup.ID == nil) || (props.NetworkSecurityGroup.ID != nil && *props.NetworkSecurityGroup.ID == "")) {
			return fmt.Errorf("No Network Security Group should exist for Subnet %q (Virtual Network %q / Resource Group: %q) but got %q", subnetName, virtualNetworkName, resourceGroupName, *props.NetworkSecurityGroup.ID)
		}

		return nil
	}
}

func testAccAzureRMSubnetNetworkSecurityGroupAssociation_basic(data acceptance.TestData) string {
	template := testAccAzureRMSubnetNetworkSecurityGroupAssociation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name               = "internal"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.0.2.0/24"
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}
`, template)
}

func testAccAzureRMSubnetNetworkSecurityGroupAssociation_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMSubnetNetworkSecurityGroupAssociation_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_network_security_group_association" "internal" {
  subnet_id                 = azurerm_subnet_network_security_group_association.test.subnet_id
  network_security_group_id = azurerm_subnet_network_security_group_association.test.network_security_group_id
}
`, template)
}

func testAccAzureRMSubnetNetworkSecurityGroupAssociation_updateSubnet(data acceptance.TestData) string {
	template := testAccAzureRMSubnetNetworkSecurityGroupAssociation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name               = "internal"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.0.2.0/24"

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}
`, template)
}

func testAccAzureRMSubnetNetworkSecurityGroupAssociation_template(data acceptance.TestData) string {
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

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule {
    name                       = "test123"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
