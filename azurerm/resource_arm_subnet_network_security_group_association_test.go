package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSubnetNetworkSecurityGroupAssociation_basic(t *testing.T) {
	resourceName := "azurerm_subnet_network_security_group_association.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetNetworkSecurityGroupAssociation_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetNetworkSecurityGroupAssociationExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSubnetNetworkSecurityGroupAssociation_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_subnet_network_security_group_association.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetNetworkSecurityGroupAssociation_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetNetworkSecurityGroupAssociationExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMSubnetNetworkSecurityGroupAssociation_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_subnet_network_security_group_association"),
			},
		},
	})
}

func TestAccAzureRMSubnetNetworkSecurityGroupAssociation_deleted(t *testing.T) {
	resourceName := "azurerm_subnet_network_security_group_association.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetNetworkSecurityGroupAssociation_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetNetworkSecurityGroupAssociationExists(resourceName),
					testCheckAzureRMSubnetNetworkSecurityGroupAssociationDisappears(resourceName),
					testCheckAzureRMSubnetHasNoNetworkSecurityGroup(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMSubnetNetworkSecurityGroupAssociationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
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

		client := testAccProvider.Meta().(*ArmClient).network.SubnetsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

		client := testAccProvider.Meta().(*ArmClient).network.SubnetsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

		client := testAccProvider.Meta().(*ArmClient).network.SubnetsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

func testAccAzureRMSubnetNetworkSecurityGroupAssociation_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                      = "acctestsubnet%d"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  virtual_network_name      = "${azurerm_virtual_network.test.name}"
  address_prefix            = "10.0.2.0/24"
  network_security_group_id = "${azurerm_network_security_group.test.id}"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

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

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = "${azurerm_subnet.test.id}"
  network_security_group_id = "${azurerm_network_security_group.test.id}"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMSubnetNetworkSecurityGroupAssociation_requiresImport(rInt int, location string) string {
	template := testAccAzureRMSubnetNetworkSecurityGroupAssociation_basic(rInt, location)
	return fmt.Sprintf(`
%s

`, template)
}
