package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMStorageEncryptionScope_keyVaultKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_encryption_scope", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageEncryptionScopeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStorageEncryptionScope_keyVaultKey(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageEncryptionScopeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "source", "Microsoft.KeyVault"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "key_vault_key_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMStorageEncryptionScope_microsoftManagedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_encryption_scope", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageEncryptionScopeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStorageEncryptionScope_microsoftManagedKey(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageEncryptionScopeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "source", "Microsoft.Storage"),
					resource.TestCheckResourceAttr(data.ResourceName, "key_vault_key_id", ""),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMStorageEncryptionScope_keyVaultKey(data acceptance.TestData) string {
	basic := testAccAzureRMStorageEncryptionScope_keyVaultKey(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_encryption_scope" "test" {
  name               = azurerm_storage_encryption_scope.test.name
  storage_account_id = azurerm_storage_encryption_scope.test.storage_account_id
}
`, basic)
}

func testAccDataSourceAzureRMStorageEncryptionScope_microsoftManagedKey(data acceptance.TestData) string {
	basic := testAccAzureRMStorageEncryptionScope_microsoftManagedKey(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_encryption_scope" "test" {
  name               = azurerm_storage_encryption_scope.test.name
  storage_account_id = azurerm_storage_encryption_scope.test.storage_account_id
}
`, basic)
}
