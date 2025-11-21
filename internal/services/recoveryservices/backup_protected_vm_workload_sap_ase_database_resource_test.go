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

type BackupProtectedVMWorkloadSAPAseDatabaseResource struct{}

func TestAccBackupProtectedVMWorkloadSAPAseDatabaseSequential(t *testing.T) {
	// The dependent VM and database require complex SAP workload configurations. Tests require pre-configured resources.
	if os.Getenv("ARM_TEST_SAP_VM_ID") == "" || os.Getenv("ARM_TEST_SAP_DATABASE_NAME") == "" || os.Getenv("ARM_TEST_SAP_DATABASE_INSTANCE_NAME") == "" {
		t.Skip("Skipping test as ARM_TEST_SAP_VM_ID, ARM_TEST_SAP_DATABASE_NAME, or ARM_TEST_SAP_DATABASE_INSTANCE_NAME environment variables are not set")
	}

	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"backupProtectedVMWorkloadSAPAseDatabase": {
			"basic":                      testAccBackupProtectedVMWorkloadSAPAseDatabase_Basic,
			"requiresImport":             testAccBackupProtectedVMWorkloadSAPAseDatabase_RequiresImport,
			"update":                     testAccBackupProtectedVMWorkloadSAPAseDatabase_Update,
			"protectionStopped":          testAccBackupProtectedVMWorkloadSAPAseDatabase_ProtectionStopped,
			"protectionStoppedOnDestroy": testAccBackupProtectedVMWorkloadSAPAseDatabase_ProtectionStoppedOnDestroy,
			"recoverSoftDeletedWorkload": testAccBackupProtectedVMWorkloadSAPAseDatabase_RecoverSoftDeletedWorkload,
		},
	})
}

func testAccBackupProtectedVMWorkloadSAPAseDatabase_Basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm_workload_sap_ase_database", "test")
	r := BackupProtectedVMWorkloadSAPAseDatabaseResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
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

func testAccBackupProtectedVMWorkloadSAPAseDatabase_RequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm_workload_sap_ase_database", "test")
	r := BackupProtectedVMWorkloadSAPAseDatabaseResource{}

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

func testAccBackupProtectedVMWorkloadSAPAseDatabase_Update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm_workload_sap_ase_database", "test")
	r := BackupProtectedVMWorkloadSAPAseDatabaseResource{}

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

func testAccBackupProtectedVMWorkloadSAPAseDatabase_ProtectionStopped(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm_workload_sap_ase_database", "test")
	r := BackupProtectedVMWorkloadSAPAseDatabaseResource{}

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

func testAccBackupProtectedVMWorkloadSAPAseDatabase_ProtectionStoppedOnDestroy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm_workload_sap_ase_database", "test")
	r := BackupProtectedVMWorkloadSAPAseDatabaseResource{}

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

func testAccBackupProtectedVMWorkloadSAPAseDatabase_RecoverSoftDeletedWorkload(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_protected_vm_workload_sap_ase_database", "test")
	r := BackupProtectedVMWorkloadSAPAseDatabaseResource{}

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

func (t BackupProtectedVMWorkloadSAPAseDatabaseResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := protecteditems.ParseProtectedItemID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ProtectedItemsClient.Get(ctx, *id, protecteditems.GetOperationOptions{})
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) basic(data acceptance.TestData) string {
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

resource "azurerm_backup_protected_vm_workload_sap_ase_database" "test" {
  resource_group_name    = azurerm_resource_group.test.name
  recovery_vault_name    = azurerm_recovery_services_vault.test.name
  source_vm_id           = "%s"
  backup_policy_id       = azurerm_backup_policy_vm_workload.test.id
  database_name          = "%s"
  database_instance_name = "%s"

  depends_on = [azurerm_backup_container_vm_app.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, os.Getenv("ARM_TEST_SAP_VM_ID"), os.Getenv("ARM_TEST_SAP_VM_ID"), os.Getenv("ARM_TEST_SAP_DATABASE_NAME"), os.Getenv("ARM_TEST_SAP_DATABASE_INSTANCE_NAME"))
}

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm_workload_sap_ase_database" "import" {
  resource_group_name    = azurerm_backup_protected_vm_workload_sap_ase_database.test.resource_group_name
  recovery_vault_name    = azurerm_backup_protected_vm_workload_sap_ase_database.test.recovery_vault_name
  source_vm_id           = azurerm_backup_protected_vm_workload_sap_ase_database.test.source_vm_id
  backup_policy_id       = azurerm_backup_protected_vm_workload_sap_ase_database.test.backup_policy_id
  database_name          = azurerm_backup_protected_vm_workload_sap_ase_database.test.database_name
  database_instance_name = azurerm_backup_protected_vm_workload_sap_ase_database.test.database_instance_name
}
`, r.basic(data))
}

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) update(data acceptance.TestData) string {
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

resource "azurerm_backup_protected_vm_workload_sap_ase_database" "test" {
  resource_group_name    = azurerm_resource_group.test.name
  recovery_vault_name    = azurerm_recovery_services_vault.test.name
  source_vm_id           = "%s"
  backup_policy_id       = azurerm_backup_policy_vm_workload.updated.id
  database_name          = "%s"
  database_instance_name = "%s"

  depends_on = [azurerm_backup_container_vm_app.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, os.Getenv("ARM_TEST_SAP_VM_ID"), os.Getenv("ARM_TEST_SAP_VM_ID"), os.Getenv("ARM_TEST_SAP_DATABASE_NAME"), os.Getenv("ARM_TEST_SAP_DATABASE_INSTANCE_NAME"))
}

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) base(data acceptance.TestData) string {
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

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) protectionStopped(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm_workload_sap_ase_database" "test" {
  resource_group_name    = azurerm_resource_group.test.name
  recovery_vault_name    = azurerm_recovery_services_vault.test.name
  source_vm_id           = "%s"
  backup_policy_id       = azurerm_backup_policy_vm_workload.test.id
  database_name          = "%s"
  database_instance_name = "%s"
  protection_state       = "ProtectionStopped"

  depends_on = [azurerm_backup_container_vm_app.test]
}
`, r.base(data), os.Getenv("ARM_TEST_SAP_VM_ID"), os.Getenv("ARM_TEST_SAP_DATABASE_NAME"), os.Getenv("ARM_TEST_SAP_DATABASE_INSTANCE_NAME"))
}

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) protectionStoppedOnDestroy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    recovery_service {
      vm_workload_sap_ase_database_backup_stop_protection_and_retain_data_on_destroy = true
    }
  }
}

%s

resource "azurerm_backup_protected_vm_workload_sap_ase_database" "test" {
  resource_group_name    = azurerm_resource_group.test.name
  recovery_vault_name    = azurerm_recovery_services_vault.test.name
  source_vm_id           = "%s"
  backup_policy_id       = azurerm_backup_policy_vm_workload.test.id
  database_name          = "%s"
  database_instance_name = "%s"

  depends_on = [azurerm_backup_container_vm_app.test]
}
`, r.baseWithOutProvider(data), os.Getenv("ARM_TEST_SAP_VM_ID"), os.Getenv("ARM_TEST_SAP_DATABASE_NAME"), os.Getenv("ARM_TEST_SAP_DATABASE_INSTANCE_NAME"))
}

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) baseWithOutProvider(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, os.Getenv("ARM_TEST_SAP_VM_ID"))
}

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) baseWithSoftDelete(data acceptance.TestData) string {
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

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) basicWithSoftDelete(data acceptance.TestData, deleted bool) string {
	if deleted {
		return r.baseWithSoftDelete(data)
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {
    recovery_service {
      vm_workload_backup_purge_soft_delete_data_on_destroy = false
    }
  }
}

%s

resource "azurerm_backup_protected_vm_workload_sap_ase_database" "test" {
  resource_group_name    = azurerm_resource_group.test.name
  recovery_vault_name    = azurerm_recovery_services_vault.test.name
  source_vm_id           = "%s"
  backup_policy_id       = azurerm_backup_policy_vm_workload.test.id
  database_name          = "%s"
  database_instance_name = "%s"

  depends_on = [azurerm_backup_container_vm_app.test]
}
`, r.baseWithSoftDelete(data), os.Getenv("ARM_TEST_SAP_VM_ID"), os.Getenv("ARM_TEST_SAP_DATABASE_NAME"), os.Getenv("ARM_TEST_SAP_DATABASE_INSTANCE_NAME"))
}
