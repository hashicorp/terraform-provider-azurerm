package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMKeyVaultKeyEncrypt_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key_encrypt", "test")
	plaintext := "testData"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultKeyEncrypt_basic(data, plaintext),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "cipher_text"),
				),
			},
		},
	})
}

func testAccAzureRMKeyVaultKeyEncrypt_basic(data acceptance.TestData, plaintext string) string {
	t := testAccAzureRMKeyVaultKeyEncrypt_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key_encrypt" "test" {
  key_vault_key_id = azurerm_key_vault_key.test.id
  algorithm        = "RSA1_5"
  plaintext        = "%s"
}
`, t, plaintext)
}

func testAccAzureRMKeyVaultKeyEncrypt_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "premium"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "create",
      "delete",
      "get",
      "update",
      "decrypt",
      "encrypt",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
