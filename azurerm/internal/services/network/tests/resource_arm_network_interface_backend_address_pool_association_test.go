package tests

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_backend_address_pool_association", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceBackendAddressPoolAssociationExists(data.ResourceName),
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

	data := acceptance.BuildTestData(t, "azurerm_network_interface_backend_address_pool_association", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceBackendAddressPoolAssociationExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_network_interface_backend_address_pool_association"),
			},
		},
	})
}

func TestAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_backend_address_pool_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceBackendAddressPoolAssociationExists(data.ResourceName),
					testCheckAzureRMNetworkInterfaceBackendAddressPoolAssociationDisappears(data.ResourceName),
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

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.InterfacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.InterfacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterfaceBackendAddressPoolAssociation_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface_backend_address_pool_association" "import" {
  network_interface_id    = "${azurerm_network_interface_backend_address_pool_association.test.network_interface_id}"
  ip_configuration_name   = "${azurerm_network_interface_backend_address_pool_association.test.ip_configuration_name}"
  backend_address_pool_id = "${azurerm_network_interface_backend_address_pool_association.test.backend_address_pool_id}"
}
`, template)
}
