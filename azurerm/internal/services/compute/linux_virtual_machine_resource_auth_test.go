package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccLinuxVirtualMachine_authPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_authPassword(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccLinuxVirtualMachine_authPasswordAndSSH(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_authPasswordAndSSH(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccLinuxVirtualMachine_authSSH(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_authSSH(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLinuxVirtualMachine_authSSHMultipleKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_authSSHMultiple(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testLinuxVirtualMachine_authPassword(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctestVM-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "P@$$w0rd1234!"
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_authPasswordAndSSH(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_authSSH(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_authSSHMultiple(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.second_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, template, data.RandomInteger)
}
