package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccWindowsVirtualMachine_dataDiskBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_dataDiskBasic(data, false),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("admin_password"),
		},
	})
}

func TestAccWindowsVirtualMachine_dataDiskDeleteOnTermination(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_dataDiskBasic(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("admin_password"),
			{
				Config: testWindowsVirtualMachine_dataDiskVMRemoved(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkVirtualMachineManagedDiskIsDeleted(fmt.Sprintf("acctestRG-%d", data.RandomInteger), "testdatadisk", fmt.Sprintf("acctestvm%s", data.RandomString)),
				),
			},
		},
	})
}

func TestAccWindowsVirtualMachine_dataDiskMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_dataDiskMultiple(data, false),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("admin_password"),
		},
	})
}

func TestAccWindowsVirtualMachine_dataDiskUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_dataDiskBasic(data, false),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("admin_password"),
			{
				Config: testWindowsVirtualMachine_dataDiskMultiple(data, false),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("admin_password"),
			{
				Config: testWindowsVirtualMachine_dataDiskRemoveFirst(data, false),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("admin_password"),
		},
	})
}

func TestAccWindowsVirtualMachine_dataDiskUpdateWithDeleteDataDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_dataDiskBasic(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("admin_password"),
			{
				Config: testWindowsVirtualMachine_dataDiskMultiple(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("admin_password"),
			{
				Config: testWindowsVirtualMachine_dataDiskRemoveFirst(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					checkVirtualMachineManagedDiskIsDeleted(fmt.Sprintf("acctestRG-%d", data.RandomInteger), "testdatadisk", fmt.Sprintf("acctestvm%s", data.RandomString)),
				),
			},
			data.ImportStep("admin_password"),
		},
	})
}

func testWindowsVirtualMachine_dataDiskBasic(data acceptance.TestData, deleteDataDisks bool) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    virtual_machine {
      delete_data_disk_on_deletion = %t
    }
  }
}

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

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
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template, deleteDataDisks)
}

func testWindowsVirtualMachine_dataDiskMultiple(data acceptance.TestData, deleteDataDisk bool) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    virtual_machine {
      delete_data_disk_on_deletion = %t
    }
  }
}

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

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
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template, deleteDataDisk)
}

func testWindowsVirtualMachine_dataDiskRemoveFirst(data acceptance.TestData, deleteDataDisk bool) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    virtual_machine {
      delete_data_disk_on_deletion = %t
    }
  }
}

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  data_disk {
    name                 = "testdatadisk2"
    lun                  = 2
    caching              = "None"
    storage_account_type = "Standard_LRS"
    disk_size_gb         = 2
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template, deleteDataDisk)
}

func testWindowsVirtualMachine_dataDiskVMRemoved(data acceptance.TestData, deleteDataDisks bool) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    virtual_machine {
      delete_data_disk_on_deletion = %t
    }
  }
}

`, template, deleteDataDisks)
}
