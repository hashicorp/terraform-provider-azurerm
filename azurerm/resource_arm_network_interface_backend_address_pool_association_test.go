package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_basic(t *testing.T) {
	resourceName := "azurerm_network_interface_backend_address_pool_association.test"
	rInt := tf.AccRandTimeInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_basic(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceBackendAddressPoolAssociationExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_network_interface_backend_address_pool_association.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceBackendAddressPoolAssociationExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_requiresImport(rInt, location),
				ExpectError: testRequiresImportError("azurerm_network_interface_backend_address_pool_association"),
			},
		},
	})
}

func TestAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_deleted(t *testing.T) {
	resourceName := "azurerm_network_interface_backend_address_pool_association.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceBackendAddressPoolAssociationExists(resourceName),
					testCheckAzureRMNetworkInterfaceBackendAddressPoolAssociationDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMNetworkInterfaceBackendAddressPoolAssociationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		nicID, err := azure.ParseAzureResourceID(rs.Primary.Attributes["network_interface_id"])
		if err != nil {
			return err
		}

		nicName := nicID.Path["networkInterfaces"]
		resourceGroup := nicID.ResourceGroup
		backendAddressPoolId := rs.Primary.Attributes["backend_address_pool_id"]
		ipConfigurationName := rs.Primary.Attributes["ip_configuration_name"]

		client := testAccProvider.Meta().(*ArmClient).network.InterfacesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		read, err := client.Get(ctx, resourceGroup, nicName, "")
		if err != nil {
			return fmt.Errorf("Error retrieving Network Interface %q (Resource Group %q): %+v", nicName, resourceGroup, err)
		}

		c := azure.FindNetworkInterfaceIPConfiguration(read.InterfacePropertiesFormat.IPConfigurations, ipConfigurationName)
		if c == nil {
			return fmt.Errorf("IP Configuration %q wasn't found for Network Interface %q (Resource Group %q)", ipConfigurationName, nicName, resourceGroup)
		}
		config := *c

		found := false
		if config.InterfaceIPConfigurationPropertiesFormat.LoadBalancerBackendAddressPools != nil {
			for _, pool := range *config.InterfaceIPConfigurationPropertiesFormat.LoadBalancerBackendAddressPools {
				if *pool.ID == backendAddressPoolId {
					found = true
					break
				}
			}
		}

		if !found {
			return fmt.Errorf("Association between NIC %q and LB Backend Address Pool %q was not found!", nicName, backendAddressPoolId)
		}

		return nil
	}
}

func testCheckAzureRMNetworkInterfaceBackendAddressPoolAssociationDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		nicID, err := azure.ParseAzureResourceID(rs.Primary.Attributes["network_interface_id"])
		if err != nil {
			return err
		}

		nicName := nicID.Path["networkInterfaces"]
		resourceGroup := nicID.ResourceGroup
		backendAddressPoolId := rs.Primary.Attributes["backend_address_pool_id"]
		ipConfigurationName := rs.Primary.Attributes["ip_configuration_name"]

		client := testAccProvider.Meta().(*ArmClient).network.InterfacesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		read, err := client.Get(ctx, resourceGroup, nicName, "")
		if err != nil {
			return fmt.Errorf("Error retrieving Network Interface %q (Resource Group %q): %+v", nicName, resourceGroup, err)
		}

		c := azure.FindNetworkInterfaceIPConfiguration(read.InterfacePropertiesFormat.IPConfigurations, ipConfigurationName)
		if c == nil {
			return fmt.Errorf("IP Configuration %q wasn't found for Network Interface %q (Resource Group %q)", ipConfigurationName, nicName, resourceGroup)
		}
		config := *c

		updatedPools := make([]network.BackendAddressPool, 0)
		if config.InterfaceIPConfigurationPropertiesFormat.LoadBalancerBackendAddressPools != nil {
			for _, pool := range *config.InterfaceIPConfigurationPropertiesFormat.LoadBalancerBackendAddressPools {
				if *pool.ID != backendAddressPoolId {
					updatedPools = append(updatedPools, pool)
				}
			}
		}
		config.InterfaceIPConfigurationPropertiesFormat.LoadBalancerBackendAddressPools = &updatedPools

		future, err := client.CreateOrUpdate(ctx, resourceGroup, nicName, read)
		if err != nil {
			return fmt.Errorf("Error removing Backend Address Pool Association for Network Interface %q (Resource Group %q): %+v", nicName, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for removal of Backend Address Pool Association for NIC %q (Resource Group %q): %+v", nicName, resourceGroup, err)
		}

		return nil
	}
}

func testAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "primary"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "acctestpool"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_network_interface_backend_address_pool_association" "test" {
  network_interface_id    = "${azurerm_network_interface.test.id}"
  ip_configuration_name   = "testconfiguration1"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.test.id}"
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_requiresImport(rInt int, location string) string {
	template := testAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface_backend_address_pool_association" "import" {
  network_interface_id    = "${azurerm_network_interface_backend_address_pool_association.test.network_interface_id}"
  ip_configuration_name   = "${azurerm_network_interface_backend_address_pool_association.test.ip_configuration_name}"
  backend_address_pool_id = "${azurerm_network_interface_backend_address_pool_association.test.backend_address_pool_id}"
}
`, template)
}
