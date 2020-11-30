package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hsm/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDedicatedHardwareSecurityModule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_hardware_security_module", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDedicatedHardwareSecurityModuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDedicatedHardwareSecurityModule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHardwareSecurityModuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDedicatedHardwareSecurityModule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_hardware_security_module", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDedicatedHardwareSecurityModuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDedicatedHardwareSecurityModule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHardwareSecurityModuleExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMDedicatedHardwareSecurityModule_requiresImport),
		},
	})
}

func TestAccAzureRMDedicatedHardwareSecurityModule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_hardware_security_module", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDedicatedHardwareSecurityModuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDedicatedHardwareSecurityModule_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHardwareSecurityModuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDedicatedHardwareSecurityModule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_hardware_security_module", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDedicatedHardwareSecurityModuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDedicatedHardwareSecurityModule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHardwareSecurityModuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDedicatedHardwareSecurityModule_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHardwareSecurityModuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMDedicatedHardwareSecurityModule_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMDedicatedHardwareSecurityModule_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_hardware_security_module" "import" {
  name                = azurerm_dedicated_hardware_security_module.test.name
  resource_group_name = azurerm_dedicated_hardware_security_module.test.resource_group_name
  location            = azurerm_dedicated_hardware_security_module.test.location
  sku_name            = azurerm_dedicated_hardware_security_module.test.sku_name
  stamp_id            = azurerm_dedicated_hardware_security_module.test.stamp_id

  network_profile {
    network_interface_private_ip_addresses = azurerm_dedicated_hardware_security_module.test.network_profile[0].network_interface_private_ip_addresses
    subnet_id                              = azurerm_dedicated_hardware_security_module.test.network_profile[0].subnet_id
  }
}
`, config)
}

func testCheckAzureRMDedicatedHardwareSecurityModuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).HSM.DedicatedHsmClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("dedicated hardware security module not found: %s", resourceName)
		}
		id, err := parse.DedicatedHardwareSecurityModuleID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.DedicatedHSMName); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Dedicated HardwareSecurityModule %q does not exist", id.DedicatedHSMName)
			}
			return fmt.Errorf("bad: Get on HardwareSecurityModules.DedicatedHsmClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMDedicatedHardwareSecurityModuleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).HSM.DedicatedHsmClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dedicated_hardware_security_module" {
			continue
		}
		id, err := parse.DedicatedHardwareSecurityModuleID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.DedicatedHSMName); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on HardwareSecurityModules.DedicatedHsmClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMDedicatedHardwareSecurityModule_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-hsm-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  address_space       = ["10.2.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-computesubnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.2.0.0/24"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctest-hsmsubnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.2.1.0/24"]

  delegation {
    name = "first"

    service_delegation {
      name = "Microsoft.HardwareSecurityModules/dedicatedHSMs"

      actions = [
        "Microsoft.Network/networkinterfaces/*",
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_subnet" "test3" {
  name                 = "gatewaysubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.2.255.0/26"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-pip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-vnetgateway-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "ExpressRoute"
  vpn_type = "PolicyBased"
  sku      = "Standard"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test3.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDedicatedHardwareSecurityModule_basic(data acceptance.TestData) string {
	template := testAccAzureRMDedicatedHardwareSecurityModule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_hardware_security_module" "test" {
  name                = "acctest-hsm-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "SafeNet Luna Network HSM A790"

  network_profile {
    network_interface_private_ip_addresses = ["10.2.1.8"]
    subnet_id                              = azurerm_subnet.test2.id
  }

  stamp_id = "stamp2"

  depends_on = [azurerm_virtual_network_gateway.test]
}
`, template, data.RandomString)
}

func testAccAzureRMDedicatedHardwareSecurityModule_complete(data acceptance.TestData) string {
	template := testAccAzureRMDedicatedHardwareSecurityModule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_hardware_security_module" "test" {
  name                = "acctest-hsm-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "SafeNet Luna Network HSM A790"

  network_profile {
    network_interface_private_ip_addresses = ["10.2.1.8"]
    subnet_id                              = azurerm_subnet.test2.id
  }

  stamp_id = "stamp2"

  tags = {
    env = "Test"
  }

  depends_on = [azurerm_virtual_network_gateway.test]
}
`, template, data.RandomString)
}
