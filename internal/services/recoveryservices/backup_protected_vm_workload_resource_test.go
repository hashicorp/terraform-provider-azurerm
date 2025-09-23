// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2025-02-01/protecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type BackupProtectedVMWorkloadResource struct{}

func TestAccBackupProtectedVMWorkloadSequential(t *testing.T) {
	// The dependent VM and database require complex SAP workload configurations. Tests require pre-configured resources.
	if os.Getenv("ARM_TEST_SAP_VM_ID") == "" || os.Getenv("ARM_TEST_SAP_DATABASE_NAME") == "" {
		t.Skip("Skipping as `ARM_TEST_SAP_VM_ID` and `ARM_TEST_SAP_DATABASE_NAME` are not specified")
	}

	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"backupProtectedVMWorkload": {
			"basic":                        testAccBackupProtectedVMWorkload_Basic,
			"requiresImport":               testAccBackupProtectedVMWorkload_RequiresImport,
			"update":                       testAccBackupProtectedVMWorkload_Update,
			"protectionStopped":            testAccBackupProtectedVMWorkload_ProtectionStopped,
			"protectionStoppedOnDestroy":   testAccBackupProtectedVMWorkload_ProtectionStoppedOnDestroy,
			"recoverSoftDeletedWorkload":   testAccBackupProtectedVMWorkload_RecoverSoftDeletedWorkload,
		},
	})
}

func testAccBackupProtectedVMWorkload_Basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm_workload", "test")
	r := BackupProtectedVMWorkloadResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("workload_type").HasValue("SAPAseDatabase"),
				check.That(data.ResourceName).Key("protected_item_name").HasValue(os.Getenv("ARM_TEST_SAP_DATABASE_NAME")),
				check.That(data.ResourceName).Key("source_vm_id").HasValue(os.Getenv("ARM_TEST_SAP_VM_ID")),
			),
		},
		data.ImportStep(),
	})
}

func testAccBackupProtectedVMWorkload_RequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm_workload", "test")
	r := BackupProtectedVMWorkloadResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccBackupProtectedVMWorkload_Update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm_workload", "test")
	r := BackupProtectedVMWorkloadResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccBackupProtectedVMWorkload_ProtectionStopped(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm_workload", "test")
	r := BackupProtectedVMWorkloadResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.protectionStopped(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("protection_state").HasValue("ProtectionStopped"),
			),
		},
		data.ImportStep(),
		{
			Config: r.base(data),
		},
	})
}

func testAccBackupProtectedVMWorkload_ProtectionStoppedOnDestroy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm_workload", "test")
	r := BackupProtectedVMWorkloadResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.protectionStoppedOnDestroy(data),
		},
	})
}

func testAccBackupProtectedVMWorkload_RecoverSoftDeletedWorkload(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm_workload", "test")
	r := BackupProtectedVMWorkloadResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithSoftDelete(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
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
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
		},
		{
			Config: r.base(data),
		},
	})
}

func (t BackupProtectedVMWorkloadResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := protecteditems.ParseProtectedItemID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ProtectedItemsClient.Get(ctx, *id, protecteditems.GetOperationOptions{})
	if err != nil {
		return nil, fmt.Errorf("reading Backup Protected VM Workload (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r BackupProtectedVMWorkloadResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-vmworkload-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}

resource "azurerm_backup_policy_vm_workload" "test" {
  name                = "acctest-bpvmw-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  workload_type = "SAPAseDatabase"

  protection_policy {
    policy_type = "Full"

    backup {
      frequency = "Daily"
      time      = "15:00"
    }

    retention_daily {
      count = 30
    }
  }
}

resource "azurerm_backup_container_vm_app" "test" {
  source_resource_id  = "%s"
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
  workload_type       = "SAPAseDatabase"
}

resource "azurerm_backup_protected_vm_workload" "test" {
  resource_group_name   = azurerm_resource_group.test.name
  recovery_vault_name   = azurerm_recovery_services_vault.test.name
  source_vm_id          = "%s"
  backup_policy_id      = azurerm_backup_policy_vm_workload.test.id
  workload_type         = "SAPAseDatabase"
  protected_item_name   = "%s"

  depends_on = [azurerm_backup_container_vm_app.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, os.Getenv("ARM_TEST_SAP_VM_ID"), os.Getenv("ARM_TEST_SAP_VM_ID"), os.Getenv("ARM_TEST_SAP_DATABASE_NAME"))
}

func (r BackupProtectedVMWorkloadResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm_workload" "import" {
  resource_group_name = azurerm_backup_protected_vm_workload.test.resource_group_name
  recovery_vault_name = azurerm_backup_protected_vm_workload.test.recovery_vault_name
  source_vm_id        = azurerm_backup_protected_vm_workload.test.source_vm_id
  backup_policy_id    = azurerm_backup_protected_vm_workload.test.backup_policy_id
  workload_type       = azurerm_backup_protected_vm_workload.test.workload_type
  protected_item_name = azurerm_backup_protected_vm_workload.test.protected_item_name
}
`, r.basic(data))
}

func (r BackupProtectedVMWorkloadResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-vmworkload-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}

resource "azurerm_backup_policy_vm_workload" "test" {
  name                = "acctest-bpvmw-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  workload_type = "SAPAseDatabase"

  protection_policy {
    policy_type = "Full"

    backup {
      frequency = "Daily"
      time      = "15:00"
    }

    retention_daily {
      count = 30
    }
  }
}

resource "azurerm_backup_policy_vm_workload" "updated" {
  name                = "acctest-bpvmw-updated-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  workload_type = "SAPAseDatabase"

  protection_policy {
    policy_type = "Full"

    backup {
      frequency = "Daily"
      time      = "16:00"
    }

    retention_daily {
      count = 60
    }
  }
}

resource "azurerm_backup_container_vm_app" "test" {
  source_resource_id  = "%s"
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
  workload_type       = "SAPAseDatabase"
}

resource "azurerm_backup_protected_vm_workload" "test" {
  resource_group_name   = azurerm_resource_group.test.name
  recovery_vault_name   = azurerm_recovery_services_vault.test.name
  source_vm_id          = "%s"
  backup_policy_id      = azurerm_backup_policy_vm_workload.updated.id
  workload_type         = "SAPAseDatabase"
  protected_item_name   = "%s"

  depends_on = [azurerm_backup_container_vm_app.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, os.Getenv("ARM_TEST_SAP_VM_ID"), os.Getenv("ARM_TEST_SAP_VM_ID"), os.Getenv("ARM_TEST_SAP_DATABASE_NAME"))
}

func (r BackupProtectedVMWorkloadResource) base(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-vmworkload-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}

resource "azurerm_backup_policy_vm_workload" "test" {
  name                = "acctest-bpvmw-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  workload_type = "SAPAseDatabase"

  protection_policy {
    policy_type = "Full"

    backup {
      frequency = "Daily"
      time      = "15:00"
    }

    retention_daily {
      count = 30
    }
  }
}

resource "azurerm_backup_container_vm_app" "test" {
  source_resource_id  = "%s"
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
  workload_type       = "SAPAseDatabase"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, os.Getenv("ARM_TEST_SAP_VM_ID"))
}

func (r BackupProtectedVMWorkloadResource) protectionStopped(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm_workload" "test" {
  resource_group_name   = azurerm_resource_group.test.name
  recovery_vault_name   = azurerm_recovery_services_vault.test.name
  source_vm_id          = "%s"
  backup_policy_id      = azurerm_backup_policy_vm_workload.test.id
  workload_type         = "SAPAseDatabase"
  protected_item_name   = "%s"
  protection_state      = "ProtectionStopped"

  depends_on = [azurerm_backup_container_vm_app.test]
}
`, r.base(data), os.Getenv("ARM_TEST_SAP_VM_ID"), os.Getenv("ARM_TEST_SAP_DATABASE_NAME"))
}

func (r BackupProtectedVMWorkloadResource) protectionStoppedOnDestroy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    recovery_service {
      vm_workload_backup_stop_protection_and_retain_data_on_destroy = true
      purge_protected_items_from_vault_on_destroy                   = true
    }
  }
}

%s
`, r.baseWithOutProvider(data))
}

func (r BackupProtectedVMWorkloadResource) baseWithOutProvider(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-vmworkload-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}

resource "azurerm_backup_policy_vm_workload" "test" {
  name                = "acctest-bpvmw-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  workload_type = "SAPAseDatabase"

  protection_policy {
    policy_type = "Full"

    backup {
      frequency = "Daily"
      time      = "15:00"
    }

    retention_daily {
      count = 30
    }
  }
}

resource "azurerm_backup_container_vm_app" "test" {
  source_resource_id  = "%s"
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
  workload_type       = "SAPAseDatabase"
}

resource "azurerm_backup_protected_vm_workload" "test" {
  resource_group_name   = azurerm_resource_group.test.name
  recovery_vault_name   = azurerm_recovery_services_vault.test.name
  source_vm_id          = "%s"
  backup_policy_id      = azurerm_backup_policy_vm_workload.test.id
  workload_type         = "SAPAseDatabase"
  protected_item_name   = "%s"

  depends_on = [azurerm_backup_container_vm_app.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, os.Getenv("ARM_TEST_SAP_VM_ID"), os.Getenv("ARM_TEST_SAP_VM_ID"), os.Getenv("ARM_TEST_SAP_DATABASE_NAME"))
}

func (r BackupProtectedVMWorkloadResource) baseWithSoftDelete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-vmworkload-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = true
}

resource "azurerm_backup_policy_vm_workload" "test" {
  name                = "acctest-bpvmw-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  workload_type = "SAPAseDatabase"

  protection_policy {
    policy_type = "Full"

    backup {
      frequency = "Daily"
      time      = "15:00"
    }

    retention_daily {
      count = 30
    }
  }
}

resource "azurerm_backup_container_vm_app" "test" {
  source_resource_id  = "%s"
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
  workload_type       = "SAPAseDatabase"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, os.Getenv("ARM_TEST_SAP_VM_ID"))
}

func (r BackupProtectedVMWorkloadResource) basicWithSoftDelete(data acceptance.TestData, deleted bool) string {
	protectedWorkloadBlock := `
resource "azurerm_backup_protected_vm_workload" "test" {
  resource_group_name   = azurerm_resource_group.test.name
  recovery_vault_name   = azurerm_recovery_services_vault.test.name
  source_vm_id          = "%s"
  backup_policy_id      = azurerm_backup_policy_vm_workload.test.id
  workload_type         = "SAPAseDatabase"
  protected_item_name   = "%s"

  depends_on = [azurerm_backup_container_vm_app.test]
}
`
	if deleted {
		protectedWorkloadBlock = ""
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {
    recovery_services_vaults {
      recover_soft_deleted_backup_protected_vm_workload = true
    }
  }
}

%s

%s
`, r.baseWithSoftDelete(data), fmt.Sprintf(protectedWorkloadBlock, os.Getenv("ARM_TEST_SAP_VM_ID"), os.Getenv("ARM_TEST_SAP_DATABASE_NAME")))
}
