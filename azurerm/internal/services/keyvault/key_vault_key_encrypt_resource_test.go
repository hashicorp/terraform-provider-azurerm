package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

func TestAccKeyVaultKeyEncrypt_basic(t *testing.T) {
	plaintext := "testData"
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key_encrypt", "test")

	testCase := resource.TestCase{
		PreCheck: func() { acceptance.PreCheck(t) },
		ProviderFactories: map[string]terraform.ResourceProviderFactory{
			"azurerm": func() (terraform.ResourceProvider, error) {
				azurerm := provider.TestAzureProvider()
				return azurerm, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccKeyVaultKeyEncrypt_basic(data, plaintext),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "encrypted_data_in_base64url"),
				),
			},
		},
	}

	resource.ParallelTest(t, testCase)
}

func testAccKeyVaultKeyEncrypt_basic(data acceptance.TestData, plaintext string) string {
	t := testAccKeyVaultKeyEncrypt_template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = false
    }
  }
}

%s

resource "azurerm_key_vault_key_encrypt" "test" {
  name             = "acctest_encrypt%d"
  key_vault_key_id = azurerm_key_vault_key.test.id
  algorithm        = "RSA1_5"
  plaintext        = "%s"
}
`, t, data.RandomInteger, plaintext)
}

func testAccKeyVaultKeyEncrypt_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
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
