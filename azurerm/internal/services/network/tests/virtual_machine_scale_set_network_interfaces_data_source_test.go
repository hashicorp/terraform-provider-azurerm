package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceArmVirtualMachineScaleSetNetworkInterfaces_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_machine_scale_set_network_interfaces", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVirtualMachineScaleSetNetworkInterfaces_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "network_interfaces.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_interfaces.0.name", "example"),
				),
			},
		},
	})
}

func testAccAzureRMWindowsVirtualMachineScaleSet_vmName(data acceptance.TestData) string {
	// windows VM names can be up to 15 chars, however the prefix can only be 9 chars
	return fmt.Sprintf("acctvm%s", fmt.Sprintf("%d", data.RandomInteger)[0:2])
}

func testAccDataSourceVirtualMachineScaleSetNetworkInterfaces_basic(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_template(data)
	name := testAccAzureRMWindowsVirtualMachineScaleSet_vmName(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = "%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_D2s_v3"
  instances           = 2
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  overprovision       = false

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name                          = "example"
    primary                       = true
    enable_accelerated_networking = false

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}

data "azurerm_virtual_machine_scale_set_network_interfaces" "test" {
  virtual_machine_scale_set_name = azurerm_windows_virtual_machine_scale_set.test.name
  resource_group_name            = azurerm_windows_virtual_machine_scale_set.test.resource_group_name
}
`, template, name)
}
