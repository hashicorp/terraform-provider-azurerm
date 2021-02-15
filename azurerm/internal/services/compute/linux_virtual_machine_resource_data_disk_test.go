package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccLinuxVirtualMachine_dataDisksBasic(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_dataDisksComplete(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksComplete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.dataDisksCompleteUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_dataDisksBasicUpdateDataDiskSize(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.dataDisksBasicUpdateDataDiskSize(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_dataDisksBasicUpdateEncryptionSet(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksBasicWithEncryption(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.dataDisksBasicWithEncryptionUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_dataDisksAddLocalDataDisk(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksAbsent(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.dataDisksBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.dataDisksAbsent(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_dataDisksAddExistingDataDisk(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksAbsent(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.dataDisksExistingDisk(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_dataDisksMultipleUpdateSize(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.dataDisksMultipleUpdateSize(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_dataDisksExistingDisk(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksExistingDisk(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_dataDisksExistingDiskResize(t *testing.T) {
	if !features.VMDataDiskBeta() {
		t.Skip("skipping as In-line Data Disk beta is not enabled")
	}

	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataDisksExistingDisk(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.dataDisksExistingDiskResize(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LinuxVirtualMachineResource) dataDisksBasic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
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

  data_disks {
    local {
      name                 = "acctest-localdisk"
      lun                  = 1
      caching              = "None"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 1
    }
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

func (r LinuxVirtualMachineResource) dataDisksAbsent(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
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

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}

`, template, data.RandomInteger)
}

func (r LinuxVirtualMachineResource) dataDisksComplete(data acceptance.TestData) string {
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

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%[2]d"
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

  data_disks {
    local {
      name                 = "acctest-localdisk"
      lun                  = 1
      caching              = "None"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 1
    }

    local {
      name                 = "acctest-localdisk2"
      lun                  = 2
      caching              = "ReadOnly"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 2
    }

    local {
      name                 = "acctest-localdisk3"
      lun                  = 3
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 3
    }

    existing {
      managed_disk_id      = azurerm_managed_disk.test1.id
      lun                  = 10
      caching              = "None"
      storage_account_type = "Standard_LRS"
    }

    existing {
      managed_disk_id      = azurerm_managed_disk.test2.id
      lun                  = 11
      caching              = "ReadOnly"
      storage_account_type = "Standard_LRS"
    }
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

func (r LinuxVirtualMachineResource) dataDisksCompleteUpdate(data acceptance.TestData) string {
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

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%[2]d"
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

  data_disks {
    local {
      name                 = "acctest-localdisk"
      lun                  = 1
      caching              = "None"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 1
    }

    local {
      name                 = "acctest-localdisk2"
      lun                  = 2
      caching              = "ReadOnly"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 4
    }

    existing {
      managed_disk_id      = azurerm_managed_disk.test1.id
      lun                  = 10
      caching              = "None"
      storage_account_type = "Standard_LRS"
    }

    existing {
      managed_disk_id      = azurerm_managed_disk.test2.id
      lun                  = 11
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }
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

func (r LinuxVirtualMachineResource) dataDisksBasicUpdateDataDiskSize(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
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

  data_disks {
    local {
      name                 = "acctest-localdisk"
      lun                  = 1
      caching              = "ReadOnly"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 2
    }
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

func (r LinuxVirtualMachineResource) dataDisksBasicWithEncryption(data acceptance.TestData) string {
	template := r.diskOSDiskDiskEncryptionSetResource(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%[2]d"
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

  data_disks {
    local {
      name                   = "acctest-localdisk"
      lun                    = 1
      caching                = "ReadOnly"
      storage_account_type   = "Standard_LRS"
      disk_encryption_set_id = azurerm_disk_encryption_set.test.id
      disk_size_gb           = 1
    }
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

func (r LinuxVirtualMachineResource) dataDisksBasicWithEncryptionUpdate(data acceptance.TestData) string {
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


resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%[2]d"
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

  data_disks {
    local {
      name                   = "acctest-localdisk"
      lun                    = 1
      caching                = "ReadOnly"
      storage_account_type   = "Standard_LRS"
      disk_encryption_set_id = azurerm_disk_encryption_set.update.id
      disk_size_gb           = 1
    }
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

func (r LinuxVirtualMachineResource) dataDisksMultipleUpdateSize(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
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

  data_disks {
    local {
      name                 = "acctest-localdisk"
      lun                  = 1
      caching              = "ReadOnly"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 2
    }

    local {
      name                 = "acctest-localdisk2"
      lun                  = 2
      caching              = "ReadOnly"
      storage_account_type = "Standard_LRS"
      disk_size_gb         = 1
    }
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

func (r LinuxVirtualMachineResource) dataDisksExistingDisk(data acceptance.TestData) string {
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

  data_disks {
    existing {
      managed_disk_id      = azurerm_managed_disk.test.id
      lun                  = 1
      caching              = "None"
      storage_account_type = "Standard_LRS"
    }
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}

`, template, data.RandomInteger, data.RandomInteger)
}

func (r LinuxVirtualMachineResource) dataDisksExistingDiskResize(data acceptance.TestData) string {
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

  data_disks {
    existing {
      managed_disk_id      = azurerm_managed_disk.test.id
      lun                  = 1
      caching              = "None"
      storage_account_type = "Standard_LRS"
    }
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}

`, template, data.RandomInteger, data.RandomInteger)
}
