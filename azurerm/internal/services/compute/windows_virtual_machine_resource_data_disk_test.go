package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccWindowsVirtualMachine_dataDisksBasic(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachine_dataDisksComplete(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksComplete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config: r.dataDisksCompleteUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachine_dataDisksBasicUpdateDataDiskSize(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config: r.dataDisksBasicUpdateDataDiskSize(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachine_dataDisksBasicUpdateEncryptionSet(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksBasicWithEncryption(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config: r.dataDisksBasicWithEncryptionUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachine_dataDisksAddLocalDataDisk(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksAbsent(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config: r.dataDisksBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config: r.dataDisksAbsent(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachine_dataDisksAddExistingDataDisk(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksAbsent(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config: r.dataDisksExistingDisk(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachine_dataDisksMultipleUpdateSize(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config: r.dataDisksMultipleUpdateSize(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachine_dataDisksExistingDisk(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksExistingDisk(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachine_dataDisksExistingDiskResize(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksExistingDisk(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config: r.dataDisksExistingDiskResize(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func (r WindowsVirtualMachineResource) dataDisksBasic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
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

  data_disks {
    create {
      name                 = "acctest-localdisk"
      lun                  = 1
      caching              = "None"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 1
    }
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func (r WindowsVirtualMachineResource) dataDisksAbsent(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
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

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}

`, template)
}

func (r WindowsVirtualMachineResource) dataDisksComplete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}
resource "azurerm_managed_disk" "test1" {
  name                 = "acctested-%[2]d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}

resource "azurerm_managed_disk" "test2" {
  name                 = "acctested2-%[2]d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  tags = {
    environment = "acctest"
    cost-center = "ops"
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

  data_disks {
    create {
      name                 = "acctest-localdisk"
      lun                  = 1
      caching              = "None"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 1
    }

    create {
      name                 = "acctest-localdisk2"
      lun                  = 2
      caching              = "ReadOnly"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 2
    }

    create {
      name                 = "acctest-localdisk3"
      lun                  = 3
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 3
    }

    attach {
      managed_disk_id      = azurerm_managed_disk.test1.id
      lun                  = 10
      caching              = "None"
      storage_account_type = "Standard_LRS"
    }

    attach {
      managed_disk_id      = azurerm_managed_disk.test2.id
      lun                  = 11
      caching              = "ReadOnly"
      storage_account_type = "Standard_LRS"
    }
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}

`, template, data.RandomInteger)
}

func (r WindowsVirtualMachineResource) dataDisksCompleteUpdate(data acceptance.TestData) string {
	template := r.template(data)
	// Updates localdisk2 to 4GiB, removes localdisk3, and updates existingdisk1 to LUN-15, existingdisk2 to ReadWrite caching,
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_managed_disk" "test1" {
  name                 = "acctested-%[2]d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}

resource "azurerm_managed_disk" "test2" {
  name                 = "acctested2-%[2]d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  tags = {
    environment = "acctest"
    cost-center = "ops"
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

  data_disks {
    create {
      name                 = "acctest-localdisk"
      lun                  = 1
      caching              = "None"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 1
    }

    create {
      name                 = "acctest-localdisk2"
      lun                  = 2
      caching              = "ReadOnly"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 4
    }

    attach {
      managed_disk_id      = azurerm_managed_disk.test1.id
      lun                  = 10
      caching              = "None"
      storage_account_type = "Standard_LRS"
    }

    attach {
      managed_disk_id      = azurerm_managed_disk.test2.id
      lun                  = 11
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}

`, template, data.RandomInteger)
}

func (r WindowsVirtualMachineResource) dataDisksBasicUpdateDataDiskSize(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
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

  data_disks {
    create {
      name                 = "acctest-localdisk"
      lun                  = 1
      caching              = "ReadOnly"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 2
    }
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}

`, template)
}

func (r WindowsVirtualMachineResource) dataDisksBasicWithEncryption(data acceptance.TestData) string {
	template := r.diskOSDiskDiskEncryptionSetResource(data)
	return fmt.Sprintf(`
%s

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

  data_disks {
    create {
      name                   = "acctest-localdisk"
      lun                    = 1
      caching                = "ReadOnly"
      storage_account_type   = "Standard_LRS"
      disk_encryption_set_id = azurerm_disk_encryption_set.test.id
      disk_size_gb           = 1
    }
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}

`, template)
}

func (r WindowsVirtualMachineResource) dataDisksBasicWithEncryptionUpdate(data acceptance.TestData) string {
	template := r.diskOSDiskDiskEncryptionSetResource(data)
	return fmt.Sprintf(`
%s

resource "azurerm_disk_encryption_set" "update" {
  name                = "acctestdes-%d-u"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned"
  }
}


resource "azurerm_key_vault_access_policy" "disk-encryption" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "get",
    "wrapkey",
    "unwrapkey",
  ]

  tenant_id = azurerm_disk_encryption_set.update.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.update.identity.0.principal_id
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

  data_disks {
    create {
      name                   = "acctest-localdisk"
      lun                    = 1
      caching                = "ReadOnly"
      storage_account_type   = "Standard_LRS"
      disk_encryption_set_id = azurerm_disk_encryption_set.update.id
      disk_size_gb           = 1
    }
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}

`, template, data.RandomInteger)
}

func (r WindowsVirtualMachineResource) dataDisksMultipleUpdateSize(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
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

  data_disks {
    create {
      name                 = "acctest-localdisk"
      lun                  = 1
      caching              = "ReadOnly"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 2
    }

    create {
      name                 = "acctest-localdisk2"
      lun                  = 2
      caching              = "ReadOnly"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 1
    }
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}

`, template)
}

func (r WindowsVirtualMachineResource) dataDisksExistingDisk(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctested-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  tags = {
    environment = "acctest"
    cost-center = "ops"
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

  data_disks {
    attach {
      managed_disk_id      = azurerm_managed_disk.test.id
      lun                  = 1
      caching              = "None"
      storage_account_type = "Standard_LRS"
    }
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}

`, template, data.RandomInteger)
}

func (r WindowsVirtualMachineResource) dataDisksExistingDiskResize(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctested-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "2"

  tags = {
    environment = "acctest"
    cost-center = "ops"
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

  data_disks {
    attach {
      managed_disk_id      = azurerm_managed_disk.test.id
      lun                  = 1
      caching              = "None"
      storage_account_type = "Standard_LRS"
    }
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}

`, template, data.RandomInteger)
}
