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

type DataProtectionBackupInstanceDiskResource struct{}

func TestAccDataProtectionBackupInstanceDisk_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_disk", "test")
	r := DataProtectionBackupInstanceDiskResource{}
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

func TestAccDataProtectionBackupInstanceDisk_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_disk", "test")
	r := DataProtectionBackupInstanceDiskResource{}
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

func TestAccDataProtectionBackupInstanceDisk_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_disk", "test")
	r := DataProtectionBackupInstanceDiskResource{}
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

func TestAccDataProtectionBackupInstanceDisk_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_disk", "test")
	r := DataProtectionBackupInstanceDiskResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
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

func (r DataProtectionBackupInstanceDiskResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backupinstances.ParseBackupInstanceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupInstanceClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving DataProtection BackupInstance (%q): %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r DataProtectionBackupInstanceDiskResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-dataprotection-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctest-disk-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-dataprotection-vault-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
  soft_delete         = "Off"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_role_assignment" "test1" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Disk Snapshot Contributor"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "test2" {
  scope                = azurerm_managed_disk.test.id
  role_definition_name = "Disk Backup Reader"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}

resource "azurerm_data_protection_backup_policy_disk" "test" {
  name                            = "acctest-dbp-%d"
  vault_id                        = azurerm_data_protection_backup_vault.test.id
  backup_repeating_time_intervals = ["R/2021-05-20T04:54:23+00:00/PT4H"]
  default_retention_duration      = "P7D"
}

resource "azurerm_data_protection_backup_policy_disk" "another" {
  name                            = "acctest-dbp-other-%d"
  vault_id                        = azurerm_data_protection_backup_vault.test.id
  backup_repeating_time_intervals = ["R/2021-05-20T04:54:23+00:00/PT4H"]
  default_retention_duration      = "P10D"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(8), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r DataProtectionBackupInstanceDiskResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_disk" "test" {
  name                         = "acctest-dbi-%d"
  location                     = azurerm_resource_group.test.location
  vault_id                     = azurerm_data_protection_backup_vault.test.id
  disk_id                      = azurerm_managed_disk.test.id
  snapshot_resource_group_name = azurerm_resource_group.test.name
  backup_policy_id             = azurerm_data_protection_backup_policy_disk.test.id
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupInstanceDiskResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_disk" "import" {
  name                         = azurerm_data_protection_backup_instance_disk.test.name
  location                     = azurerm_data_protection_backup_instance_disk.test.location
  vault_id                     = azurerm_data_protection_backup_instance_disk.test.vault_id
  disk_id                      = azurerm_data_protection_backup_instance_disk.test.disk_id
  snapshot_resource_group_name = azurerm_data_protection_backup_instance_disk.test.snapshot_resource_group_name
  backup_policy_id             = azurerm_data_protection_backup_instance_disk.test.backup_policy_id
}
`, config)
}

func (r DataProtectionBackupInstanceDiskResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_disk" "test" {
  name                         = "acctest-dbi-%d"
  location                     = azurerm_resource_group.test.location
  vault_id                     = azurerm_data_protection_backup_vault.test.id
  disk_id                      = azurerm_managed_disk.test.id
  snapshot_resource_group_name = azurerm_resource_group.test.name
  backup_policy_id             = azurerm_data_protection_backup_policy_disk.test.id
}
`, template, data.RandomInteger)
}
