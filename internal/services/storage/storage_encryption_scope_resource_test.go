// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/encryptionscopes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageEncryptionScopeResource struct{}

func TestAccStorageEncryptionScope_keyVaultKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_encryption_scope", "test")
	r := StorageEncryptionScopeResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source").HasValue("Microsoft.KeyVault"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageEncryptionScope_keyVaultKeyVersionless(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_encryption_scope", "test")
	r := StorageEncryptionScopeResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultKeyVersionless(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source").HasValue("Microsoft.KeyVault"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageEncryptionScope_keyVaultKeyRequireInfrastructureEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_encryption_scope", "test")
	r := StorageEncryptionScopeResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultKeyRequireInfrastructureEncryption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source").HasValue("Microsoft.KeyVault"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageEncryptionScope_keyVaultKeyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_encryption_scope", "test")
	r := StorageEncryptionScopeResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source").HasValue("Microsoft.KeyVault"),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultKeyUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source").HasValue("Microsoft.KeyVault"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageEncryptionScope_keyVaultKeyToMicrosoftManagedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_encryption_scope", "test")
	r := StorageEncryptionScopeResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source").HasValue("Microsoft.KeyVault"),
			),
		},
		data.ImportStep(),
		{
			Config: r.microsoftManagedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source").HasValue("Microsoft.Storage"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageEncryptionScope_microsoftManagedKeyToKeyVaultManagedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_encryption_scope", "test")

	r := StorageEncryptionScopeResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.microsoftManagedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source").HasValue("Microsoft.Storage"),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source").HasValue("Microsoft.KeyVault"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageEncryptionScope_microsoftManagedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_encryption_scope", "test")

	r := StorageEncryptionScopeResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.microsoftManagedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source").HasValue("Microsoft.Storage"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageEncryptionScope_microsoftManagedKeyRequireInfrastructureEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_encryption_scope", "test")

	r := StorageEncryptionScopeResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.microsoftManagedKeyRequireInfrastructureEncryption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source").HasValue("Microsoft.Storage"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageEncryptionScope_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_encryption_scope", "test")

	r := StorageEncryptionScopeResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.microsoftManagedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source").HasValue("Microsoft.Storage"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t StorageEncryptionScopeResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := encryptionscopes.ParseEncryptionScopeID(state.Attributes["id"])
	if err != nil {
		return nil, err
	}

	resp, err := clients.Storage.ResourceManager.EncryptionScopes.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	enabled := false
	if model := resp.Model; model != nil && model.Properties != nil && model.Properties.State != nil {
		enabled = *model.Properties.State == encryptionscopes.EncryptionScopeStateEnabled
	}

	return utils.Bool(enabled), nil
}

func (t StorageEncryptionScopeResource) keyVaultKey(data acceptance.TestData) string {
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

%s

resource "azurerm_storage_encryption_scope" "test" {
  name               = "acctestES%d"
  storage_account_id = azurerm_storage_account.test.id
  source             = "Microsoft.KeyVault"
  key_vault_key_id   = azurerm_key_vault_key.first.id
}
`, template, data.RandomInteger)
}

func (t StorageEncryptionScopeResource) keyVaultKeyVersionless(data acceptance.TestData) string {
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

%s

resource "azurerm_storage_encryption_scope" "test" {
  name               = "acctestES%d"
  storage_account_id = azurerm_storage_account.test.id
  source             = "Microsoft.KeyVault"
  key_vault_key_id   = azurerm_key_vault_key.first.versionless_id
}
`, template, data.RandomInteger)
}

func (t StorageEncryptionScopeResource) keyVaultKeyUpdated(data acceptance.TestData) string {
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

%s

resource "azurerm_key_vault_key" "second" {
  name         = "second"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.storage,
  ]
}

resource "azurerm_storage_encryption_scope" "test" {
  name               = "acctestES%d"
  storage_account_id = azurerm_storage_account.test.id
  source             = "Microsoft.KeyVault"
  key_vault_key_id   = azurerm_key_vault_key.second.id
}
`, template, data.RandomInteger)
}

func (t StorageEncryptionScopeResource) keyVaultKeyRequireInfrastructureEncryption(data acceptance.TestData) string {
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

%s

resource "azurerm_storage_encryption_scope" "test" {
  name                               = "acctestES%d"
  storage_account_id                 = azurerm_storage_account.test.id
  source                             = "Microsoft.KeyVault"
  key_vault_key_id                   = azurerm_key_vault_key.first.id
  infrastructure_encryption_required = true
}
`, template, data.RandomInteger)
}

func (t StorageEncryptionScopeResource) microsoftManagedKey(data acceptance.TestData) string {
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

%s
resource "azurerm_storage_encryption_scope" "test" {
  name               = "acctestES%d"
  storage_account_id = azurerm_storage_account.test.id
  source             = "Microsoft.Storage"
}
`, template, data.RandomInteger)
}

func (t StorageEncryptionScopeResource) microsoftManagedKeyRequireInfrastructureEncryption(data acceptance.TestData) string {
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

%s
resource "azurerm_storage_encryption_scope" "test" {
  name                               = "acctestES%d"
  storage_account_id                 = azurerm_storage_account.test.id
  source                             = "Microsoft.Storage"
  infrastructure_encryption_required = true
}
`, template, data.RandomInteger)
}

func (t StorageEncryptionScopeResource) requiresImport(data acceptance.TestData) string {
	template := t.microsoftManagedKey(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_encryption_scope" "import" {
  name               = azurerm_storage_encryption_scope.test.name
  storage_account_id = azurerm_storage_encryption_scope.test.storage_account_id
  source             = azurerm_storage_encryption_scope.test.source
}
`, template)
}

func (StorageEncryptionScopeResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "storage" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_storage_account.test.identity.0.principal_id

  key_permissions = ["Get", "UnwrapKey", "WrapKey"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
}

resource "azurerm_key_vault_key" "first" {
  name         = "first"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.storage,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
