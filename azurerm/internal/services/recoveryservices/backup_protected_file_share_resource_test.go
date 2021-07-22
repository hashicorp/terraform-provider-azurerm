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

type BackupProtectedFileShareResource struct {
}

// TODO: These tests fail because enabling backup on file shares with no content
func TestAccBackupProtectedFileShare_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_file_share", "test")
	r := BackupProtectedFileShareResource{}

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

func TestAccBackupProtectedFileShare_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_file_share", "test")
	r := BackupProtectedFileShareResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			// vault cannot be deleted unless we unregister all backups
			Config: r.baseMultiple(data),
		},
	})
}

func TestAccBackupProtectedFileShare_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_file_share", "test")
	r := BackupProtectedFileShareResource{}

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

func TestAccBackupProtectedFileShare_updateBackupPolicyId(t *testing.T) {
	fBackupPolicyResourceName := "azurerm_backup_policy_file_share.test1"
	sBackupPolicyResourceName := "azurerm_backup_policy_file_share.test2"

	data := acceptance.BuildTestData(t, "azurerm_backup_protected_file_share", "test")
	r := BackupProtectedFileShareResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Create resources and link first backup policy id
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttrPair(data.ResourceName, "backup_policy_id", fBackupPolicyResourceName, "id"),
			),
		},
		{
			// Modify backup policy id to the second one
			// Set Destroy false to prevent error from cleaning up dangling resource
			Config: r.updatePolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttrPair(data.ResourceName, "backup_policy_id", sBackupPolicyResourceName, "id"),
			),
		},
		{
			// Remove protected items first before the associated policies are deleted
			Config: r.base(data),
		},
	})
}

func (t BackupProtectedFileShareResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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
		return nil, fmt.Errorf("reading Recovery Service Protected File Share (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (BackupProtectedFileShareResource) base(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[3]s"
  location                 = "${azurerm_resource_group.test.location}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
  name                 = "acctest-ss-%[1]d"
  storage_account_name = "${azurerm_storage_account.test.name}"
  metadata             = {}

  lifecycle {
    ignore_changes = [metadata] // Ignore changes Azure Backup makes to the metadata
  }
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-VAULT-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"

  soft_delete_enabled = true
}

resource "azurerm_backup_policy_file_share" "test1" {
  name                = "acctest-PFS-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r BackupProtectedFileShareResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_container_storage_account" "test" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  storage_account_id  = azurerm_storage_account.test.id
}

resource "azurerm_backup_protected_file_share" "test" {
  resource_group_name       = azurerm_resource_group.test.name
  recovery_vault_name       = azurerm_recovery_services_vault.test.name
  source_storage_account_id = azurerm_backup_container_storage_account.test.storage_account_id
  source_file_share_name    = azurerm_storage_share.test.name
  backup_policy_id          = azurerm_backup_policy_file_share.test1.id
}
`, r.base(data))
}

func (r BackupProtectedFileShareResource) updatePolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_file_share" "test2" {
  name                = "acctest-%d-Secondary"
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

resource "azurerm_backup_container_storage_account" "test" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  storage_account_id  = azurerm_storage_account.test.id
}

resource "azurerm_backup_protected_file_share" "test" {
  resource_group_name       = azurerm_resource_group.test.name
  recovery_vault_name       = azurerm_recovery_services_vault.test.name
  source_storage_account_id = azurerm_backup_container_storage_account.test.storage_account_id
  source_file_share_name    = azurerm_storage_share.test.name
  backup_policy_id          = azurerm_backup_policy_file_share.test2.id
}
`, r.base(data), data.RandomInteger)
}

func (r BackupProtectedFileShareResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_file_share" "test_import" {
  resource_group_name       = azurerm_resource_group.test.name
  recovery_vault_name       = azurerm_recovery_services_vault.test.name
  source_storage_account_id = azurerm_backup_container_storage_account.test.storage_account_id
  source_file_share_name    = azurerm_storage_share.test.name
  backup_policy_id          = azurerm_backup_policy_file_share.test1.id
}
`, r.basic(data))
}

func (r BackupProtectedFileShareResource) baseMultiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test1" {
  name                     = "acctest%[3]s1"
  location                 = "${azurerm_resource_group.test.location}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest%[3]s2"
  location                 = "${azurerm_resource_group.test.location}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "testshare1" {
  name                 = "acctest-ss-%[1]d-1"
  storage_account_name = "${azurerm_storage_account.test1.name}"
  metadata             = {}

  lifecycle {
    ignore_changes = [metadata] // Ignore changes Azure Backup makes to the metadata
  }
}

resource "azurerm_storage_share" "testshare2" {
  name                 = "acctest-ss-%[1]d-2"
  storage_account_name = "${azurerm_storage_account.test1.name}"
  metadata             = {}

  lifecycle {
    ignore_changes = [metadata] // Ignore changes Azure Backup makes to the metadata
  }
}

resource "azurerm_storage_share" "testshare3" {
  name                 = "acctest-ss-%[1]d-1"
  storage_account_name = "${azurerm_storage_account.test2.name}"
  metadata             = {}

  lifecycle {
    ignore_changes = [metadata] // Ignore changes Azure Backup makes to the metadata
  }
}

resource "azurerm_storage_share" "testshare4" {
  name                 = "acctest-ss-%[1]d-2"
  storage_account_name = "${azurerm_storage_account.test2.name}"
  metadata             = {}

  lifecycle {
    ignore_changes = [metadata] // Ignore changes Azure Backup makes to the metadata
  }
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-VAULT-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"

  soft_delete_enabled = true
}

resource "azurerm_backup_policy_file_share" "test" {
  name                = "acctest-PFS-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r BackupProtectedFileShareResource) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_container_storage_account" "test1" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  storage_account_id  = azurerm_storage_account.test1.id
}

resource "azurerm_backup_container_storage_account" "test2" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  storage_account_id  = azurerm_storage_account.test2.id
}

resource "azurerm_backup_protected_file_share" "test" {
  resource_group_name       = azurerm_resource_group.test.name
  recovery_vault_name       = azurerm_recovery_services_vault.test.name
  source_storage_account_id = azurerm_backup_container_storage_account.test2.storage_account_id
  source_file_share_name    = azurerm_storage_share.testshare3.name
  backup_policy_id          = azurerm_backup_policy_file_share.test.id
}
`, r.baseMultiple(data))
}
