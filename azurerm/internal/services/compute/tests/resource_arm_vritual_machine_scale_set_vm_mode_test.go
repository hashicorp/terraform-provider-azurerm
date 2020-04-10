package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMWindowsVirtualMachineScaleSetVMMode_basicLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_vm_mode", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetVMModeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSetVMMode_basicLinux(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetVMModeExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMWindowsVirtualMachineScaleSetVMModeDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMScaleSetClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_machine_scale_set_vm_mode" {
			continue
		}

		id, err := parse.VirtualMachineScaleSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Compute.VMScaleSetClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testCheckAzureRMWindowsVirtualMachineScaleSetVMModeExists(resourceName string) resource.TestCheckFunc {
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
				return fmt.Errorf("Bad: Virtual Machine Scale Set VM Mode %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
			}

			return fmt.Errorf("Bad: Get on Compute.VMScaleSetClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMWindowsVirtualMachineScaleSetVMMode_basicLinux(data acceptance.TestData) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSetVMMode_template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_machine_scale_set_vm_mode" "test" {
  name                 = "acctestVMO-%[2]d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  
  platform_fault_domain_count = 5

  zones = ["1"]

  tags = {
    environment = "Terraform Deployment"
  }
}

resource "azurerm_linux_virtual_machine" "main" {
  name                            = "acctestVM-%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "P@ssw0rd1234!"
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  virtual_machine_scale_set_id = azurerm_virtual_machine_scale_set_vm_mode.test.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMWindowsVirtualMachineScaleSetVMMode_template(data acceptance.TestData) string {
	// in VMSS VMO mode, the `platform_fault_domain_count` has different acceptable values for different locations,
	// therefore this location is fixed to EastUS2 to make sure the acceptance test can always pass
	location := "EastUS2"
	return fmt.Sprintf(`
provider "azurerm" {
features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-VMSS-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-network-%[1]d"
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

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.RandomInteger, location)
}
