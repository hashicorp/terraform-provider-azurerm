// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type BackupProtectedVmResource struct{}

func TestAccBackupProtectedVm_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm", "test")
	r := BackupProtectedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			// vault cannot be deleted unless we unregister all backups
			Config: r.base(data),
		},
	})
}

func TestAccBackupProtectedVm_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm", "test")
	r := BackupProtectedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
		{
			// vault cannot be deleted unless we unregister all backups
			Config: r.base(data),
		},
	})
}

func TestAccBackupProtectedVm_separateResourceGroups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm", "test")
	r := BackupProtectedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.separateResourceGroups(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			// vault cannot be deleted unless we unregister all backups
			Config: r.additionalVault(data),
		},
	})
}

func TestAccBackupProtectedVm_updateBackupPolicyId(t *testing.T) {
	virtualMachine := "azurerm_virtual_machine.test"
	fBackupPolicyResourceName := "azurerm_backup_policy_vm.test"
	sBackupPolicyResourceName := "azurerm_backup_policy_vm.test_change_backup"
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm", "test")
	r := BackupProtectedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{ // Create resources and link first backup policy id
			ResourceName: fBackupPolicyResourceName,
			Config:       r.linkFirstBackupPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttrPair(data.ResourceName, "backup_policy_id", fBackupPolicyResourceName, "id"),
			),
		},
		{ // Modify backup policy id to the second one
			// Set Destroy false to prevent error from cleaning up dangling resource
			ResourceName: sBackupPolicyResourceName,
			Config:       r.linkSecondBackupPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttrPair(data.ResourceName, "backup_policy_id", sBackupPolicyResourceName, "id"),
			),
		},
		{
			// Remove backup policy link
			// Backup policy link will need to be removed first so the VM's backup policy subsequently reverts to Default
			// Azure API is quite sensitive, adding the step to control resource cleanup order
			ResourceName: fBackupPolicyResourceName,
			Config:       r.withBasePolicy(data),
		},
		{
			// Then VM can be removed
			ResourceName: virtualMachine,
			Config:       r.withBasePolicy(data),
		},
		{
			// Remove backup policies and vault
			ResourceName: data.ResourceName,
			Config:       r.basePolicyTest(data),
		},
	})
}

func TestAccBackupProtectedVm_updateVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm", "test")
	r := BackupProtectedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updateVaultFirstBackupVm(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateVaultSecondBackupVm(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		{
			// vault cannot be deleted unless we unregister all backups
			Config: r.additionalVault(data),
		},
	})
}

func TestAccBackupProtectedVm_updateDiskExclusion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm", "test")
	r := BackupProtectedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateDiskExclusion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectedVm_removeVM(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm", "test")
	r := BackupProtectedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.removeVM(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectedVm_protectionStopped(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm", "test")
	r := BackupProtectedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.protectionStopped(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			// vault cannot be deleted unless we unregister all backups
			Config: r.base(data),
		},
	})
}

func TestAccBackupProtectedVm_protectionStoppedOnDestroy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm", "test")
	r := BackupProtectedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.protectionStoppedOnDestroy(data),
		},
	})
}

func TestAccBackupProtectedVm_recoverSoftDeletedVM(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm", "test")
	r := BackupProtectedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithSoftDelete(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithSoftDelete(data, true),
		},
		{
			Config: r.basicWithSoftDelete(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			// to disable soft delete feature
			Config: r.basic(data),
		},
		{
			// vault cannot be deleted unless we unregister all backups
			Config: r.base(data),
		},
	})
}

func (t BackupProtectedVmResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := protecteditems.ParseProtectedItemID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ProtectedItemsClient.Get(ctx, *id, protecteditems.GetOperationOptions{})
	if err != nil {
		return nil, fmt.Errorf("reading Recovery Service Protected VM (%s): %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (BackupProtectedVmResource) base(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet"
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest_subnet"
  virtual_network_name = azurerm_virtual_network.test.name
  resource_group_name  = azurerm_resource_group.test.name
  address_prefixes     = ["10.0.10.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctest_nic"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "acctestipconfig"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-ip"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  domain_name_label   = "acctestip%d"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctest-datadisk"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1023"
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctestvm"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  vm_size               = "Standard_D1_v2"
  network_interface_ids = [azurerm_network_interface.test.id]

  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  storage_os_disk {
    name              = "acctest-osdisk"
    managed_disk_type = "Standard_LRS"
    caching           = "ReadWrite"
    create_option     = "FromImage"
  }

  storage_data_disk {
    name              = "acctest-datadisk"
    managed_disk_id   = azurerm_managed_disk.test.id
    managed_disk_type = "Standard_LRS"
    disk_size_gb      = azurerm_managed_disk.test.disk_size_gb
    create_option     = "Attach"
    lun               = 0
  }

  storage_data_disk {
    name              = "acctest-another-datadisk"
    create_option     = "Empty"
    disk_size_gb      = "1"
    lun               = 1
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "acctest"
    admin_username = "vmadmin"
    admin_password = "Password123!@#"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  boot_diagnostics {
    enabled     = true
    storage_uri = azurerm_storage_account.test.primary_blob_endpoint
  }

}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (BackupProtectedVmResource) baseWithoutVM(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet"
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest_subnet"
  virtual_network_name = azurerm_virtual_network.test.name
  resource_group_name  = azurerm_resource_group.test.name
  address_prefixes     = ["10.0.10.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctest_nic"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "acctestipconfig"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-ip"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  domain_name_label   = "acctestip%d"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (BackupProtectedVmResource) baseWithSoftDelete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet"
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest_subnet"
  virtual_network_name = azurerm_virtual_network.test.name
  resource_group_name  = azurerm_resource_group.test.name
  address_prefixes     = ["10.0.10.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctest_nic"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "acctestipconfig"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-ip"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  domain_name_label   = "acctestip%d"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctest-datadisk"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1023"
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctestvm"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  vm_size               = "Standard_D1_v2"
  network_interface_ids = [azurerm_network_interface.test.id]

  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  storage_os_disk {
    name              = "acctest-osdisk"
    managed_disk_type = "Standard_LRS"
    caching           = "ReadWrite"
    create_option     = "FromImage"
  }

  storage_data_disk {
    name              = "acctest-datadisk"
    managed_disk_id   = azurerm_managed_disk.test.id
    managed_disk_type = "Standard_LRS"
    disk_size_gb      = azurerm_managed_disk.test.disk_size_gb
    create_option     = "Attach"
    lun               = 0
  }

  storage_data_disk {
    name              = "acctest-another-datadisk"
    create_option     = "Empty"
    disk_size_gb      = "1"
    lun               = 1
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "acctest"
    admin_username = "vmadmin"
    admin_password = "Password123!@#"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  boot_diagnostics {
    enabled     = true
    storage_uri = azurerm_storage_account.test.primary_blob_endpoint
  }

  lifecycle {
    ignore_changes = [tags]
  }
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = true
}

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r BackupProtectedVmResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  source_vm_id        = azurerm_virtual_machine.test.id
  backup_policy_id    = azurerm_backup_policy_vm.test.id

  include_disk_luns = [0]
}
`, r.base(data))
}

func (r BackupProtectedVmResource) updateDiskExclusion(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  source_vm_id        = azurerm_virtual_machine.test.id
  backup_policy_id    = azurerm_backup_policy_vm.test.id

  exclude_disk_luns = [0, 1]
}
`, r.base(data))
}

// For update backup policy id test
func (BackupProtectedVmResource) basePolicyTest(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-%d-1"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet"
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest_subnet"
  virtual_network_name = azurerm_virtual_network.test.name
  resource_group_name  = azurerm_resource_group.test.name
  address_prefixes     = ["10.0.10.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctest_nic"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "acctestipconfig"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-ip"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  domain_name_label   = "acctestip%d"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctest-datadisk"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1023"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

// For update backup policy id test
func (r BackupProtectedVmResource) withBasePolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_backup_policy_vm" "test_change_backup" {
  name                = "acctest2-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 15
  }
}
`, r.base(data), data.RandomInteger)
}

// For update backup policy id test
func (r BackupProtectedVmResource) linkFirstBackupPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  source_vm_id        = azurerm_virtual_machine.test.id
  backup_policy_id    = azurerm_backup_policy_vm.test.id
}
`, r.withBasePolicy(data))
}

// For update backup policy id test
func (r BackupProtectedVmResource) linkSecondBackupPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  source_vm_id        = azurerm_virtual_machine.test.id
  backup_policy_id    = azurerm_backup_policy_vm.test_change_backup.id
}
`, r.withBasePolicy(data))
}

func (r BackupProtectedVmResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm" "import" {
  resource_group_name = azurerm_backup_protected_vm.test.resource_group_name
  recovery_vault_name = azurerm_backup_protected_vm.test.recovery_vault_name
  source_vm_id        = azurerm_backup_protected_vm.test.source_vm_id
  backup_policy_id    = azurerm_backup_protected_vm.test.backup_policy_id
}
`, r.basic(data))
}

func (r BackupProtectedVmResource) additionalVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-backup-%d-2"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test2" {
  name                = "acctest2-%d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  sku                 = "Standard"

  soft_delete_enabled = false
}

resource "azurerm_backup_policy_vm" "test2" {
  name                = "acctest2-%d"
  resource_group_name = azurerm_resource_group.test2.name
  recovery_vault_name = azurerm_recovery_services_vault.test2.name

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
`, r.base(data), data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r BackupProtectedVmResource) separateResourceGroups(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = azurerm_resource_group.test2.name
  recovery_vault_name = azurerm_recovery_services_vault.test2.name
  backup_policy_id    = azurerm_backup_policy_vm.test2.id
  source_vm_id        = azurerm_virtual_machine.test.id
}
`, r.additionalVault(data))
}

func (r BackupProtectedVmResource) removeVM(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  backup_policy_id    = azurerm_backup_policy_vm.test.id

  include_disk_luns = [0]
}
`, r.baseWithoutVM(data))
}

func (r BackupProtectedVmResource) updateVaultFirstBackupVm(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  backup_policy_id    = azurerm_backup_policy_vm.test.id
  source_vm_id        = azurerm_virtual_machine.test.id
}
`, r.additionalVault(data))
}

func (r BackupProtectedVmResource) updateVaultSecondBackupVm(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = azurerm_resource_group.test2.name
  recovery_vault_name = azurerm_recovery_services_vault.test2.name
  backup_policy_id    = azurerm_backup_policy_vm.test2.id
  source_vm_id        = azurerm_virtual_machine.test.id
}
`, r.additionalVault(data))
}

func (r BackupProtectedVmResource) protectionStopped(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  source_vm_id        = azurerm_virtual_machine.test.id

  include_disk_luns = [0]
  protection_state  = "ProtectionStopped"
}
`, r.base(data))
}

func (r BackupProtectedVmResource) protectionStoppedOnDestroy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    recovery_service {
      vm_backup_stop_protection_and_retain_data_on_destroy = true
      purge_protected_items_from_vault_on_destroy          = true
    }
  }
}

%s
`, r.base(data))
}

func (r BackupProtectedVmResource) basicWithSoftDelete(data acceptance.TestData, deleted bool) string {
	protectedVMBlock := `
resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  source_vm_id        = azurerm_virtual_machine.test.id
  backup_policy_id    = azurerm_backup_policy_vm.test.id

  include_disk_luns = [0]
}
`
	if deleted {
		protectedVMBlock = ""
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {
    recovery_services_vaults {
      recover_soft_deleted_backup_protected_vm = true
    }
  }
}

%s

%s
`, r.baseWithSoftDelete(data), protectedVMBlock)
}
