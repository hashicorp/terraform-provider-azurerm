// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type KeyVaultSecretEphemeral struct{}

func TestAccDataSourceKeyVaultEphemeral_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_secret", "test")
	r := KeyVaultSecretEphemeral{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("value").HasValue("rick-and-morty"),
			),
		},
	})
}

func (KeyVaultSecretEphemeral) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_secret" "test" {
  name         = azurerm_key_vault_secret.test.name
  key_vault_id = azurerm_key_vault.test.id
  version      = azurerm_key_vault_secret.test.version
}

ephemeral "azurerm_key_vault_secret" "test" {
  name         = azurerm_key_vault_secret.test.name
  key_vault_id = azurerm_key_vault.test.id
}
`, KeyVaultSecretResource{}.basic(data))
}
