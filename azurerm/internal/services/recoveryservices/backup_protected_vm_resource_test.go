package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type BackupProtectedVmResource struct {
}

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

func (t BackupProtectedVmResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	protectedItemName := id.Path["protectedItems"]
	vaultName := id.Path["vaults"]
	resourceGroup := id.ResourceGroup
	containerName := id.Path["protectionContainers"]

	resp, err := clients.RecoveryServices.ProtectedItemsClient.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
	if err != nil {
		return nil, fmt.Errorf("reading Recovery Service Protected VM (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (BackupProtectedVmResource) base(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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
  vm_size               = "Standard_A0"
  network_interface_ids = [azurerm_network_interface.test.id]

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
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

func (r BackupProtectedVmResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  source_vm_id        = azurerm_virtual_machine.test.id
  backup_policy_id    = azurerm_backup_policy_vm.test.id
}
`, r.base(data))
}

// For update backup policy id test
func (BackupProtectedVmResource) basePolicyTest(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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
