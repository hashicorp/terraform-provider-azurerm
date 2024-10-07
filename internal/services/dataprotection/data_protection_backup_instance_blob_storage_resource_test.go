// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DataProtectionBackupInstanceBlobStorageResource struct{}

func TestAccDataProtectionBackupInstanceBlobStorage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_blob_storage", "test")
	r := DataProtectionBackupInstanceBlobStorageResource{}
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

func TestAccDataProtectionBackupInstanceBlobStorage_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_blob_storage", "test")
	r := DataProtectionBackupInstanceBlobStorageResource{}
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

func TestAccDataProtectionBackupInstanceBlobStorage_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_blob_storage", "test")
	r := DataProtectionBackupInstanceBlobStorageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupInstanceBlobStorage_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_blob_storage", "test")
	r := DataProtectionBackupInstanceBlobStorageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r DataProtectionBackupInstanceBlobStorageResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backupinstances.ParseBackupInstanceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupInstanceClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r DataProtectionBackupInstanceBlobStorageResource) template(data acceptance.TestData) string {
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
}

resource "azurerm_storage_container" "test" {
  name                  = "testaccsc%[3]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_storage_container" "another" {
  name                  = "testaccsc2%[3]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
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

resource "azurerm_data_protection_backup_policy_blob_storage" "test" {
  name                                   = "acctest-dbp-%[1]d"
  vault_id                               = azurerm_data_protection_backup_vault.test.id
  operational_default_retention_duration = "P30D"
}

resource "azurerm_data_protection_backup_policy_blob_storage" "another" {
  name                                   = "acctest-dbp-other-%[1]d"
  vault_id                               = azurerm_data_protection_backup_vault.test.id
  operational_default_retention_duration = "P30D"
}

resource "azurerm_data_protection_backup_policy_blob_storage" "hybrid" {
  name                                   = "acctest-dbp-hybrid-%[1]d"
  vault_id                               = azurerm_data_protection_backup_vault.test.id
  operational_default_retention_duration = "P30D"

  backup_repeating_time_intervals  = ["R/2024-05-08T11:30:00+00:00/P1W"]
  vault_default_retention_duration = "P7D"

  retention_rule {
    name     = "Monthly"
    priority = 15
    life_cycle {
      duration        = "P6M"
      data_store_type = "VaultStore"
    }
    criteria {
      days_of_month = [1, 2, 0]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(8))
}

func (r DataProtectionBackupInstanceBlobStorageResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_data_protection_backup_instance_blob_storage" "test" {
  name               = "acctest-dbi-%d"
  location           = azurerm_resource_group.test.location
  vault_id           = azurerm_data_protection_backup_vault.test.id
  storage_account_id = azurerm_storage_account.test.id
  backup_policy_id   = azurerm_data_protection_backup_policy_blob_storage.test.id

  depends_on = [azurerm_role_assignment.test]
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupInstanceBlobStorageResource) basicUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_data_protection_backup_instance_blob_storage" "test" {
  name               = "acctest-dbi-%d"
  location           = azurerm_resource_group.test.location
  vault_id           = azurerm_data_protection_backup_vault.test.id
  storage_account_id = azurerm_storage_account.test.id
  backup_policy_id   = azurerm_data_protection_backup_policy_blob_storage.another.id

  depends_on = [azurerm_role_assignment.test]
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupInstanceBlobStorageResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_data_protection_backup_instance_blob_storage" "import" {
  name               = azurerm_data_protection_backup_instance_blob_storage.test.name
  location           = azurerm_data_protection_backup_instance_blob_storage.test.location
  vault_id           = azurerm_data_protection_backup_instance_blob_storage.test.vault_id
  storage_account_id = azurerm_data_protection_backup_instance_blob_storage.test.storage_account_id
  backup_policy_id   = azurerm_data_protection_backup_instance_blob_storage.test.backup_policy_id

  depends_on = [azurerm_role_assignment.test]
}
`, config)
}

func (r DataProtectionBackupInstanceBlobStorageResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_data_protection_backup_instance_blob_storage" "test" {
  name                            = "acctest-dbi-%d"
  location                        = azurerm_resource_group.test.location
  vault_id                        = azurerm_data_protection_backup_vault.test.id
  storage_account_id              = azurerm_storage_account.test.id
  backup_policy_id                = azurerm_data_protection_backup_policy_blob_storage.hybrid.id
  storage_account_container_names = [azurerm_storage_container.test.name]

  depends_on = [azurerm_role_assignment.test]
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupInstanceBlobStorageResource) completeUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_data_protection_backup_instance_blob_storage" "test" {
  name                            = "acctest-dbi-%d"
  location                        = azurerm_resource_group.test.location
  vault_id                        = azurerm_data_protection_backup_vault.test.id
  storage_account_id              = azurerm_storage_account.test.id
  backup_policy_id                = azurerm_data_protection_backup_policy_blob_storage.hybrid.id
  storage_account_container_names = [azurerm_storage_container.another.name]

  depends_on = [azurerm_role_assignment.test]
}
`, template, data.RandomInteger)
}
