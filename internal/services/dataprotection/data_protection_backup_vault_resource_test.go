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

type DataProtectionBackupVaultResource struct{}

func TestAccDataProtectionBackupVault_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault", "test")
	r := DataProtectionBackupVaultResource{}
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

func TestAccDataProtectionBackupVault_crossRegionRestore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault", "test")
	r := DataProtectionBackupVaultResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.crossRegionRestore(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.crossRegionRestore(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupVault_zoneRedundant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault", "test")
	r := DataProtectionBackupVaultResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zoneRedundant(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupVault_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault", "test")
	r := DataProtectionBackupVaultResource{}
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

func TestAccDataProtectionBackupVault_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault", "test")
	r := DataProtectionBackupVaultResource{}
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

func TestAccDataProtectionBackupVault_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault", "test")
	r := DataProtectionBackupVaultResource{}
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
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupVault_updateIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault", "test")
	r := DataProtectionBackupVaultResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateIdentityToSystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateIdentityToSystemAndUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateIdentityToUserAssigned(data),
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

func TestAccDataProtectionBackupVault_datastoreTypeOperationalStore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault", "test")
	r := DataProtectionBackupVaultResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.datastoreType(data, "OperationalStore", "LocallyRedundant", 14),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupVault_datastoreTypeArchiveStore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault", "test")
	r := DataProtectionBackupVaultResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.datastoreType(data, "ArchiveStore", "LocallyRedundant", 14),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.datastoreType(data, "ArchiveStore", "LocallyRedundant", 30),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.datastoreType(data, "ArchiveStore", "LocallyRedundant", 14),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupVault_datastoreTypeArchiveStoreGeoRedundant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault", "test")
	r := DataProtectionBackupVaultResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.datastoreType(data, "ArchiveStore", "GeoRedundant", 14),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupVault_storageSettingsVaultFirst(t *testing.T) {
	if !features.SixPointOh() {
		t.Skip("storage_setting is only available in 6.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault", "test")
	r := DataProtectionBackupVaultResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageSettingsVaultFirst(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r DataProtectionBackupVaultResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r DataProtectionBackupVaultResource) templateResourceGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-dataprotection-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r DataProtectionBackupVaultResource) basic(data acceptance.TestData) string {
	template := r.templateResourceGroup(data)
	if !features.SixPointOh() {
		return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
}
`, template, data.RandomInteger)
	}

	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  storage_setting {
    datastore_type = "VaultStore"
    redundancy     = "LocallyRedundant"
  }
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupVaultResource) crossRegionRestore(data acceptance.TestData, enabled bool) string {
	template := r.templateResourceGroup(data)
	if !features.SixPointOh() {
		return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_vault" "test" {
  name                         = "acctest-bv-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  datastore_type               = "VaultStore"
  redundancy                   = "GeoRedundant"
  cross_region_restore_enabled = %t
}
`, template, data.RandomInteger, enabled)
	}

	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  storage_setting {
    datastore_type = "VaultStore"
    redundancy     = "GeoRedundant"
  }
  cross_region_restore_enabled = %t
}
`, template, data.RandomInteger, enabled)
}

func (r DataProtectionBackupVaultResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	if !features.SixPointOh() {
		return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_vault" "import" {
  name                = azurerm_data_protection_backup_vault.test.name
  resource_group_name = azurerm_data_protection_backup_vault.test.resource_group_name
  location            = azurerm_data_protection_backup_vault.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
}
`, config)
	}

	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_vault" "import" {
  name                = azurerm_data_protection_backup_vault.test.name
  resource_group_name = azurerm_data_protection_backup_vault.test.resource_group_name
  location            = azurerm_data_protection_backup_vault.test.location
  storage_setting {
    datastore_type = "VaultStore"
    redundancy     = "LocallyRedundant"
  }
}
`, config)
}

func (r DataProtectionBackupVaultResource) complete(data acceptance.TestData) string {
	template := r.templateResourceGroup(data)
	if !features.SixPointOh() {
		return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"

  identity {
    type = "SystemAssigned"
  }

  immutability               = "Disabled"
  soft_delete                = "Off"
  retention_duration_in_days = 14

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
	}

	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  storage_setting {
    datastore_type = "VaultStore"
    redundancy     = "LocallyRedundant"
  }

  identity {
    type = "SystemAssigned"
  }

  immutability               = "Disabled"
  soft_delete                = "Off"
  retention_duration_in_days = 14

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupVaultResource) completeUpdate(data acceptance.TestData) string {
	template := r.templateResourceGroup(data)
	if !features.SixPointOh() {
		return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctestBV-%[2]d"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  immutability               = "Locked"
  soft_delete                = "On"
  retention_duration_in_days = 15

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
	}

	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctestBV-%[2]d"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  storage_setting {
    datastore_type = "VaultStore"
    redundancy     = "LocallyRedundant"
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  immutability               = "Locked"
  soft_delete                = "On"
  retention_duration_in_days = 15

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r DataProtectionBackupVaultResource) updateIdentityToSystemAssigned(data acceptance.TestData) string {
	template := r.templateResourceGroup(data)
	if !features.SixPointOh() {
		return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
	}

	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  storage_setting {
    datastore_type = "VaultStore"
    redundancy     = "LocallyRedundant"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupVaultResource) updateIdentityToSystemAndUserAssigned(data acceptance.TestData) string {
	template := r.templateResourceGroup(data)
	if !features.SixPointOh() {
		return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctestBV-%[2]d"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, template, data.RandomInteger, data.RandomInteger)
	}

	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctestBV-%[2]d"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  storage_setting {
    datastore_type = "VaultStore"
    redundancy     = "LocallyRedundant"
  }

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r DataProtectionBackupVaultResource) updateIdentityToUserAssigned(data acceptance.TestData) string {
	template := r.templateResourceGroup(data)
	if !features.SixPointOh() {
		return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctestBV-%[2]d"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, template, data.RandomInteger, data.RandomInteger)
	}

	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctestBV-%[2]d"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  storage_setting {
    datastore_type = "VaultStore"
    redundancy     = "LocallyRedundant"
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r DataProtectionBackupVaultResource) zoneRedundant(data acceptance.TestData) string {
	template := r.templateResourceGroup(data)
	if !features.SixPointOh() {
		return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "ZoneRedundant"
}
`, template, data.RandomInteger)
	}

	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  storage_setting {
    datastore_type = "VaultStore"
    redundancy     = "ZoneRedundant"
  }
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupVaultResource) datastoreType(data acceptance.TestData, datastoreType, redundancy string, retentionDays int) string {
	if !features.SixPointOh() {
		return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_vault" "test" {
  name                       = "acctest-bv-%[2]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  datastore_type             = "%[4]s"
  redundancy                 = "%[5]s"
  retention_duration_in_days = %[3]d
}
`, r.templateResourceGroup(data), data.RandomInteger, retentionDays, datastoreType, redundancy)
	}

	if datastoreType == "ArchiveStore" {
		return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  storage_setting {
    datastore_type = "%[4]s"
    redundancy     = "%[5]s"
  }

  storage_setting {
    datastore_type = "VaultStore"
    redundancy     = "%[5]s"
  }
  retention_duration_in_days = %[3]d
}
`, r.templateResourceGroup(data), data.RandomInteger, retentionDays, datastoreType, redundancy)
	}

	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  storage_setting {
    datastore_type = "%[4]s"
    redundancy     = "%[5]s"
  }
  retention_duration_in_days = %[3]d
}
`, r.templateResourceGroup(data), data.RandomInteger, retentionDays, datastoreType, redundancy)
}

func (r DataProtectionBackupVaultResource) storageSettingsVaultFirst(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  storage_setting {
    datastore_type = "VaultStore"
    redundancy     = "LocallyRedundant"
  }

  storage_setting {
    datastore_type = "ArchiveStore"
    redundancy     = "LocallyRedundant"
  }
}
`, r.templateResourceGroup(data), data.RandomInteger)
}

func (DataProtectionBackupVaultResource) template(data acceptance.TestData, softDelete string) string {
	if !features.SixPointOh() {
		return fmt.Sprintf(`
resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
  soft_delete         = "%[2]s"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, softDelete)
	}

	return fmt.Sprintf(`
resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  soft_delete         = "%[2]s"

  storage_setting {
    datastore_type = "VaultStore"
    redundancy     = "LocallyRedundant"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, softDelete)
}
