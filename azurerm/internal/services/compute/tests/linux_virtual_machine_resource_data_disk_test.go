package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccLinuxVirtualMachine_dataDiskBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_dataDiskBasic(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLinuxVirtualMachine_dataDiskDeleteOnTermination(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_dataDiskDeleteOnTermination(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("data_disk.0.delete_on_termination"),
		},
	})
}

func TestAccLinuxVirtualMachine_dataDiskMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_dataDiskMultiple(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testLinuxVirtualMachine_dataDiskBasic(data acceptance.TestData) string {
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

  data_disk {
    name                 = "testdatadisk"
    lun                  = 1
    caching              = "None"
    storage_account_type = "Standard_LRS"
    disk_size_gb         = 1
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

func testLinuxVirtualMachine_dataDiskDeleteOnTermination(data acceptance.TestData) string {
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

  data_disk {
    name                  = "testdatadisk"
    lun                   = 1
    caching               = "None"
    storage_account_type  = "Standard_LRS"
    disk_size_gb          = 1
	delete_on_termination = true
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

func testLinuxVirtualMachine_dataDiskMultiple(data acceptance.TestData) string {
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

  data_disk {
    name                 = "testdatadisk"
    lun                  = 1
    caching              = "None"
    storage_account_type = "Standard_LRS"
    disk_size_gb         = 1
  }

  data_disk {
    name                 = "testdatadisk2"
    lun                  = 2
    caching              = "None"
    storage_account_type = "Standard_LRS"
    disk_size_gb         = 2
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
