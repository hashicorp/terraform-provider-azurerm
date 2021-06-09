package storage_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type StorageEncryptionScopeDataSource struct{}

func TestAccDataSourceStorageEncryptionScope_keyVaultKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_encryption_scope", "test")
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageEncryptionScopeDataSource{}.keyVaultKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("source").HasValue("Microsoft.KeyVault"),
			),
		},
	})
}

func TestAccDataSourceStorageEncryptionScope_microsoftManagedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_encryption_scope", "test")
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageEncryptionScopeDataSource{}.microsoftManagedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("source").HasValue("Microsoft.Storage"),
				check.That(data.ResourceName).Key("key_vault_key_id").IsEmpty(),
			),
		},
	})
}

func (StorageEncryptionScopeDataSource) keyVaultKey(data acceptance.TestData) string {
	basic := StorageEncryptionScopeResource{}.keyVaultKey(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_encryption_scope" "test" {
  name               = azurerm_storage_encryption_scope.test.name
  storage_account_id = azurerm_storage_encryption_scope.test.storage_account_id
}
`, basic)
}

func (StorageEncryptionScopeDataSource) microsoftManagedKey(data acceptance.TestData) string {
	basic := StorageEncryptionScopeResource{}.microsoftManagedKey(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_encryption_scope" "test" {
  name               = azurerm_storage_encryption_scope.test.name
  storage_account_id = azurerm_storage_encryption_scope.test.storage_account_id
}
`, basic)
}
