package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccLinuxVirtualMachine_dataDiskBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_dataDiskBasic(data, false),
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
				Config: testLinuxVirtualMachine_dataDiskBasic(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testLinuxVirtualMachine_dataDiskBasicVMRemoved(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkVirtualMachineManagedDiskIsDeleted(fmt.Sprintf("acctestRG-%d", data.RandomInteger), "testdatadisk", fmt.Sprintf("acctestVM-%d", data.RandomInteger)),
				),
			},
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
				Config: testLinuxVirtualMachine_dataDiskMultiple(data, false),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLinuxVirtualMachine_dataDiskUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_dataDiskBasic(data, false),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testLinuxVirtualMachine_dataDiskMultiple(data, false),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testLinuxVirtualMachine_dataDiskRemoveFirst(data, false),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					checkVirtualMachineManagedDiskIsDeleted(fmt.Sprintf("acctestRG-%d", data.RandomInteger), "testdatadisk", fmt.Sprintf("acctestVM-%d", data.RandomInteger)),
				),
				// Data disk should not be deleted
				ExpectError: regexp.MustCompile("bad: Data Disk \"testdatadisk\" for Virtual Machine"),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLinuxVirtualMachine_dataDiskUpdateWithDeleteDataDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_dataDiskBasic(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testLinuxVirtualMachine_dataDiskMultiple(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testLinuxVirtualMachine_dataDiskRemoveFirst(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					checkVirtualMachineManagedDiskIsDeleted(fmt.Sprintf("acctestRG-%d", data.RandomInteger), "testdatadisk", fmt.Sprintf("acctestVM-%d", data.RandomInteger)),
				),
			},
			data.ImportStep(),
		},
	})
}

// excluding from linter - can be reused for checking OS disk tests later
//nolint unparam
func checkVirtualMachineManagedDiskIsDeleted(resourceGroup string, dataDiskName string, vmName string) resource.TestCheckFunc {
	// Since we cannot rely on the state to provide the information, as it may be deleted, to find the disk we expect it to follow the acc test pattern
	return func(s *terraform.State) error {
		disksClient := acceptance.AzureProvider.Meta().(*clients.Client).Compute.DisksClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := disksClient.Get(ctx, resourceGroup, dataDiskName); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			} else {
				return fmt.Errorf("bad: failed to check Data Disk %q for Virtual Machine %q (resource group %q) was deleted: %+v", dataDiskName, vmName, resourceGroup, err)
			}
		}
		return fmt.Errorf("bad: Data Disk %q for Virtual Machine %q (resource group %q) was not deleted", dataDiskName, vmName, resourceGroup)
	}
}

func testLinuxVirtualMachine_dataDiskBasic(data acceptance.TestData, deleteDataDisks bool) string {
	template := testLinuxVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    virtual_machine {
      delete_data_disk_on_deletion = %t
    }
  }
}

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
`, template, deleteDataDisks, data.RandomInteger)
}

func testLinuxVirtualMachine_dataDiskBasicVMRemoved(data acceptance.TestData, deleteDataDisks bool) string {
	template := testLinuxVirtualMachine_template(data)
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

func testLinuxVirtualMachine_dataDiskMultiple(data acceptance.TestData, deleteDataDisks bool) string {
	template := testLinuxVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    virtual_machine {
      delete_data_disk_on_deletion = %t
    }
  }
}

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
`, template, deleteDataDisks, data.RandomInteger)
}

func testLinuxVirtualMachine_dataDiskRemoveFirst(data acceptance.TestData, deleteDateDisks bool) string {
	template := testLinuxVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    virtual_machine {
      delete_data_disk_on_deletion = %t
    }
  }
}

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
`, template, deleteDateDisks, data.RandomInteger)
}
