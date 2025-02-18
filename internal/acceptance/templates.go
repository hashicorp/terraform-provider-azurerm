// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"fmt"
)

// WriteOnlyKeyVaultSecretTemplate is a testing template specific for write-only attributes.
// It provisions a Key Vault, a Key Vault secret, and references the secret using the Key Vault Secret
// ephemeral resource.
func WriteOnlyKeyVaultSecretTemplate(data TestData, secret string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%[1]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Get",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "List",
      "Purge",
      "Recover",
      "Set",
    ]
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[1]s"
  value        = "%[2]s"
  key_vault_id = azurerm_key_vault.test.id
}

ephemeral "azurerm_key_vault_secret" "test" {
  name         = azurerm_key_vault_secret.test.name
  key_vault_id = azurerm_key_vault.test.id
}
`, data.RandomString, secret)
}
