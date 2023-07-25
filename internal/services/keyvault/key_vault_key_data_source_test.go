// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type KeyVaultKeyDataSource struct{}

func TestAccDataSourceKeyVaultKey_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_key", "test")
	r := KeyVaultKeyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("key_type").HasValue("RSA"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.hello").HasValue("world"),
				check.That(data.ResourceName).Key("resource_id").MatchesRegex(regexp.MustCompile(`^/subscriptions/[\w-]+/resourceGroups/[\w-]+/providers/Microsoft.KeyVault/vaults/[\w-]+/keys/[\w-]+/versions/[\w-]+$`)),
				check.That(data.ResourceName).Key("resource_versionless_id").MatchesRegex(regexp.MustCompile(`^/subscriptions/[\w-]+/resourceGroups/[\w-]+/providers/Microsoft.KeyVault/vaults/[\w-]+/keys/[\w-]+$`)),
			),
		},
	})
}

func (KeyVaultKeyDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_key" "test" {
  name         = azurerm_key_vault_key.test.name
  key_vault_id = azurerm_key_vault.test.id
}
`, KeyVaultKeyResource{}.complete(data))
}
