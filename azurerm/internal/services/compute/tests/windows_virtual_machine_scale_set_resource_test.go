package tests

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testCheckAzureRMWindowsVirtualMachineScaleSetDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMScaleSetClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_windows_virtual_machine_scale_set" {
			continue
		}

		id, err := parse.VirtualMachineScaleSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testCheckAzureRMWindowsVirtualMachineScaleSetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMScaleSetClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.VirtualMachineScaleSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Windows Virtual Machine Scale Set %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
			}

			return fmt.Errorf("Bad: Get on VMScaleSetClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMWindowsVirtualMachineScaleSet_vmName(data acceptance.TestData) string {
	// windows VM names can be up to 15 chars, however the prefix can only be 9 chars
	return fmt.Sprintf("acctvm%s", fmt.Sprintf("%d", data.RandomInteger)[0:2])
}

func testAccAzureRMWindowsVirtualMachineScaleSet_template(data acceptance.TestData) string {
	name := testAccAzureRMWindowsVirtualMachineScaleSet_vmName(data)
	return fmt.Sprintf(`
locals {
  vm_name = "%s"
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
  name               = "internal"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefixes   = ["10.0.2.0/24"]
}
`, name, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
