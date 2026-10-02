// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataProtectionBackupInstanceDataLakeStorageResource struct{}

func TestAccDataProtectionBackupInstanceDataLakeStorage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_data_lake_storage", "test")
	r := DataProtectionBackupInstanceDataLakeStorageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupInstanceDataLakeStorage_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_data_lake_storage", "test")
	r := DataProtectionBackupInstanceDataLakeStorageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDataProtectionBackupInstanceDataLakeStorage_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_data_lake_storage", "test")
	r := DataProtectionBackupInstanceDataLakeStorageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupInstanceDataLakeStorage_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_data_lake_storage", "test")
	r := DataProtectionBackupInstanceDataLakeStorageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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

func (r DataProtectionBackupInstanceDataLakeStorageResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backupinstanceresources.ParseBackupInstanceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupInstanceClient.BackupInstancesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-dataprotection-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled           = true
}

resource "azurerm_storage_container" "test" {
  name               = "acctestsc%[3]d"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-dataprotection-vault-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
  soft_delete         = "Off"
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Account Backup Contributor"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}

resource "azurerm_data_protection_backup_policy_data_lake_storage" "test" {
  name                            = "acctest-dbp-%[1]d"
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  backup_schedule                 = ["R/2021-05-23T02:30:00+00:00/P1W"]

  default_retention_duration = "P4M"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(8))
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) templateComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "another" {
  name               = "acctestsc2%[2]d"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_data_protection_backup_policy_data_lake_storage" "another" {
  name                            = "acctest-dbp-other-%[3]d"
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  backup_schedule                 = ["R/2021-05-23T02:30:00+00:00/P1W", "R/2021-05-24T03:40:00+00:00/P1W"]
  time_zone                       = "Coordinated Universal Time"

  default_retention_duration = "P4M"

  retention_rule {
    name              = "weekly"
    duration          = "P6M"
    absolute_criteria = "FirstOfWeek"
  }

  retention_rule {
    name                   = "thursday"
    duration               = "P1W"
    days_of_week           = ["Thursday", "Friday"]
    months_of_year         = ["November", "December"]
    scheduled_backup_times = ["2021-05-23T02:30:00Z"]
  }

  retention_rule {
    name                   = "monthly"
    duration               = "P1D"
    weeks_of_month         = ["First", "Last"]
    days_of_week           = ["Tuesday"]
    scheduled_backup_times = ["2021-05-23T02:30:00Z", "2021-05-24T03:40:00Z"]
  }
}
`, r.template(data), data.RandomIntOfLength(8), data.RandomInteger)
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_data_lake_storage" "test" {
  name                               = "acctest-dbi-%d"
  data_protection_backup_vault_id    = azurerm_data_protection_backup_vault.test.id
  location                           = azurerm_resource_group.test.location
  storage_account_id                 = azurerm_storage_account.test.id
  backup_policy_data_lake_storage_id = azurerm_data_protection_backup_policy_data_lake_storage.test.id
  storage_container_names            = [azurerm_storage_container.test.name]

  depends_on = [azurerm_role_assignment.test]
}
`, r.template(data), data.RandomInteger)
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_data_lake_storage" "import" {
  name                               = azurerm_data_protection_backup_instance_data_lake_storage.test.name
  data_protection_backup_vault_id    = azurerm_data_protection_backup_instance_data_lake_storage.test.data_protection_backup_vault_id
  location                           = azurerm_data_protection_backup_instance_data_lake_storage.test.location
  storage_account_id                 = azurerm_data_protection_backup_instance_data_lake_storage.test.storage_account_id
  backup_policy_data_lake_storage_id = azurerm_data_protection_backup_instance_data_lake_storage.test.backup_policy_data_lake_storage_id
  storage_container_names            = azurerm_data_protection_backup_instance_data_lake_storage.test.storage_container_names

  depends_on = [azurerm_role_assignment.test]
}
`, r.basic(data))
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_data_lake_storage" "test" {
  name                               = "acctest-dbi-%d"
  data_protection_backup_vault_id    = azurerm_data_protection_backup_vault.test.id
  location                           = azurerm_resource_group.test.location
  storage_account_id                 = azurerm_storage_account.test.id
  backup_policy_data_lake_storage_id = azurerm_data_protection_backup_policy_data_lake_storage.another.id
  storage_container_names            = [azurerm_storage_container.test.name, azurerm_storage_container.another.name]

  depends_on = [azurerm_role_assignment.test]
}
`, r.templateComplete(data), data.RandomInteger)
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_data_lake_storage" "test" {
  name                               = "acctest-dbi-%d"
  data_protection_backup_vault_id    = azurerm_data_protection_backup_vault.test.id
  location                           = azurerm_resource_group.test.location
  storage_account_id                 = azurerm_storage_account.test.id
  backup_policy_data_lake_storage_id = azurerm_data_protection_backup_policy_data_lake_storage.another.id
  storage_container_names            = [azurerm_storage_container.test.name, azurerm_storage_container.another.name]

  depends_on = [azurerm_role_assignment.test]
}
`, r.templateComplete(data), data.RandomInteger)
}
