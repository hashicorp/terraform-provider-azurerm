package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation_basic(t *testing.T) {
	resourceName := "azurerm_network_interface_application_gateway_backend_address_pool_association.test"
	rInt := tf.AccRandTimeInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation_basic(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_network_interface_application_gateway_backend_address_pool_association.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation_requiresImport(rInt, location),
				ExpectError: testRequiresImportError("azurerm_network_interface_application_gateway_backend_address_pool_association"),
			},
		},
	})
}

func TestAccAzureRMNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation_deleted(t *testing.T) {
	resourceName := "azurerm_network_interface_application_gateway_backend_address_pool_association.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationExists(resourceName),
					testCheckAzureRMNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationExists(resourceName string) resource.TestCheckFunc {
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
		if config.InterfaceIPConfigurationPropertiesFormat.ApplicationGatewayBackendAddressPools != nil {
			for _, pool := range *config.InterfaceIPConfigurationPropertiesFormat.ApplicationGatewayBackendAddressPools {
				if *pool.ID == backendAddressPoolId {
					found = true
					break
				}
			}
		}

		if !found {
			return fmt.Errorf("Association between NIC %q and Application Gateway Backend Address Pool %q was not found!", nicName, backendAddressPoolId)
		}

		return nil
	}
}

func testCheckAzureRMNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationDisappears(resourceName string) resource.TestCheckFunc {
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

		updatedPools := make([]network.ApplicationGatewayBackendAddressPool, 0)
		if config.InterfaceIPConfigurationPropertiesFormat.ApplicationGatewayBackendAddressPools != nil {
			for _, pool := range *config.InterfaceIPConfigurationPropertiesFormat.ApplicationGatewayBackendAddressPools {
				if *pool.ID != backendAddressPoolId {
					updatedPools = append(updatedPools, pool)
				}
			}
		}
		config.InterfaceIPConfigurationPropertiesFormat.ApplicationGatewayBackendAddressPools = &updatedPools

		future, err := client.CreateOrUpdate(ctx, resourceGroup, nicName, read)
		if err != nil {
			return fmt.Errorf("Error removing Application Gateway Backend Address Pool Association for Network Interface %q (Resource Group %q): %+v", nicName, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for removal of Application Gateway Backend Address Pool Association for NIC %q (Resource Group %q): %+v", nicName, resourceGroup, err)
		}

		return nil
	}
}

func testAccAzureRMNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation_basic(rInt int, location string) string {
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

resource "azurerm_subnet" "frontend" {
  name                 = "frontend"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_subnet" "backend" {
  name                 = "backend"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.4.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Dynamic"
}

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "apptestag%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.frontend.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "${local.request_routing_rule_name}"
    rule_type                  = "Basic"
    http_listener_name         = "${local.listener_name}"
    backend_address_pool_name  = "${local.backend_address_pool_name}"
    backend_http_settings_name = "${local.http_setting_name}"
  }
}

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.frontend.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_network_interface_application_gateway_backend_address_pool_association" "test" {
  network_interface_id    = "${azurerm_network_interface.test.id}"
  ip_configuration_name   = "testconfiguration1"
  backend_address_pool_id = "${azurerm_application_gateway.test.backend_address_pool.0.id}"
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation_requiresImport(rInt int, location string) string {
	template := testAccAzureRMNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface_application_gateway_backend_address_pool_association" "import" {
  network_interface_id    = "${azurerm_network_interface_application_gateway_backend_address_pool_association.test.network_interface_id}"
  ip_configuration_name   = "${azurerm_network_interface_application_gateway_backend_address_pool_association.test.ip_configuration_name}"
  backend_address_pool_id = "${azurerm_network_interface_application_gateway_backend_address_pool_association.test.backend_address_pool_id}"
}
`, template)
}
