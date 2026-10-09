// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataProtectionBackupVaultDataSource struct{}

func TestAccDataProtectionBackupVaultDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_protection_backup_vault", "test")
	r := DataProtectionBackupVaultDataSource{}
	storageSettingsCheck := acceptance.ComposeTestCheckFunc(
		check.That(data.ResourceName).Key("storage_settings.#").HasValue("1"),
		check.That(data.ResourceName).Key("storage_settings.0.datastore_type").HasValue("VaultStore"),
		check.That(data.ResourceName).Key("storage_settings.0.redundancy").HasValue("LocallyRedundant"),
	)
	if !features.SixPointOh() {
		storageSettingsCheck = acceptance.ComposeTestCheckFunc(
			check.That(data.ResourceName).Key("datastore_type").HasValue("VaultStore"),
			check.That(data.ResourceName).Key("redundancy").HasValue("LocallyRedundant"),
		)
	}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				storageSettingsCheck,
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("Test"),
			),
		},
	})
}

func TestAccDataProtectionBackupVaultDataSource_archiveStore(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_protection_backup_vault", "test")
	r := DataProtectionBackupVaultDataSource{}
	storageSettingsCheck := acceptance.ComposeTestCheckFunc(
		check.That(data.ResourceName).Key("storage_settings.#").HasValue("2"),
		check.That(data.ResourceName).Key("storage_settings.0.datastore_type").HasValue("ArchiveStore"),
		check.That(data.ResourceName).Key("storage_settings.0.redundancy").HasValue("GeoRedundant"),
		check.That(data.ResourceName).Key("storage_settings.1.datastore_type").HasValue("VaultStore"),
		check.That(data.ResourceName).Key("storage_settings.1.redundancy").HasValue("GeoRedundant"),
	)
	if !features.SixPointOh() {
		storageSettingsCheck = acceptance.ComposeTestCheckFunc(
			check.That(data.ResourceName).Key("datastore_type").HasValue("ArchiveStore"),
			check.That(data.ResourceName).Key("redundancy").HasValue("GeoRedundant"),
		)
	}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.archiveStore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				storageSettingsCheck,
			),
		},
	})
}

func (r DataProtectionBackupVaultDataSource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backupvaultresources.ParseBackupVaultID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupVaultClient.BackupVaultsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving DataProtection BackupVault (%q): %+v", id, err)
	}
	return pointer.To(true), nil
}

func (r DataProtectionBackupVaultDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_protection_backup_vault" "test" {
  name                = azurerm_data_protection_backup_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, DataProtectionBackupVaultResource{}.complete(data))
}

func (DataProtectionBackupVaultDataSource) archiveStore(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_protection_backup_vault" "test" {
  name                = azurerm_data_protection_backup_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, DataProtectionBackupVaultResource{}.datastoreType(data, "ArchiveStore", "GeoRedundant", 14))
}
