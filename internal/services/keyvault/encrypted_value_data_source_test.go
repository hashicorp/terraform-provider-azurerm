package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type EncryptedValueDataSourceTest struct{}

func TestAccEncryptedValueDataSource_encryptAndDecrypt(t *testing.T) {
	// since this config includes both Encrypted and Decrypted we're testing both use-cases (and comparing the values below)
	// so we only need a single test here
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_encrypted_value", "decrypted")
	r := EncryptedValueDataSourceTest{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.decrypt(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("encrypted_value").MatchesOtherKey(check.That("data.azurerm_key_vault_encrypted_value.encrypted").Key("encrypted_value")),
				check.That(data.ResourceName).Key("plain_text_value").MatchesOtherKey(check.That("data.azurerm_key_vault_encrypted_value.encrypted").Key("plain_text_value")),
			),
		},
	})
}

func (t EncryptedValueDataSourceTest) decrypt(data acceptance.TestData) string {
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_key_vault_encrypted_value" "encrypted" {
  key_vault_key_id = azurerm_key_vault_key.test.id
  algorithm        = "RSA1_5"
  plain_text_value = "some-encrypted-value"
}

data "azurerm_key_vault_encrypted_value" "decrypted" {
  key_vault_key_id = azurerm_key_vault_key.test.id
  algorithm        = "RSA1_5"
  encrypted_data   = data.azurerm_key_vault_encrypted_value.encrypted.encrypted_data
}
`, template)
}

func (t EncryptedValueDataSourceTest) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%[3]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Delete",
      "Decrypt",
      "Encrypt",
      "Get",
      "Purge",
      "Recover",
      "Update",
      "GetRotationPolicy",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%[3]s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
