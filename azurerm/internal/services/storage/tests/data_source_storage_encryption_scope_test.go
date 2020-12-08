package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type StorageEncryptionScopeDataSourceTests struct{}

func TestAccDataSourceStorageEncryptionScope_keyVaultKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_encryption_scope", "test")
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: StorageEncryptionScopeDataSourceTests{}.keyVaultKey(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("source").HasValue("Microsoft.KeyVault"),
			),
		},
	})
}

func TestAccDataSourceStorageEncryptionScope_microsoftManagedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_encryption_scope", "test")
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: StorageEncryptionScopeDataSourceTests{}.microsoftManagedKey(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("source").HasValue("Microsoft.Storage"),
				check.That(data.ResourceName).Key("key_vault_key_id").IsEmpty(),
			),
		},
	})
}

func (StorageEncryptionScopeDataSourceTests) keyVaultKey(data acceptance.TestData) string {
	basic := StorageEncryptionScopeResourceTests{}.keyVaultKey(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_encryption_scope" "test" {
  name               = azurerm_storage_encryption_scope.test.name
  storage_account_id = azurerm_storage_encryption_scope.test.storage_account_id
}
`, basic)
}

func (StorageEncryptionScopeDataSourceTests) microsoftManagedKey(data acceptance.TestData) string {
	basic := StorageEncryptionScopeResourceTests{}.microsoftManagedKey(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_encryption_scope" "test" {
  name               = azurerm_storage_encryption_scope.test.name
  storage_account_id = azurerm_storage_encryption_scope.test.storage_account_id
}
`, basic)
}
