package compute_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func checkWindowsVirtualMachineIsDestroyed(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if rs.Type != "azurerm_windows_virtual_machine" {
			continue
		}

		id, err := parse.VirtualMachineID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func checkWindowsVirtualMachineExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.VirtualMachineID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Windows Virtual Machine %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
			}

			return fmt.Errorf("Bad: Get on VMClient: %+v", err)
		}

		return nil
	}
}

func testWindowsVirtualMachine_templateBase(data acceptance.TestData) string {
	return fmt.Sprintf(`
locals {
  vm_name = "acctestvm%s"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}
`, data.RandomString, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testWindowsVirtualMachine_template(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_templateBase(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, template, data.RandomInteger)
}
