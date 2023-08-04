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

type KeyVaultSecretDataSource struct{}

func TestAccDataSourceKeyVaultSecret_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_secret", "test")
	r := KeyVaultSecretDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("value").HasValue("rick-and-morty"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("resource_id").MatchesRegex(regexp.MustCompile(`^/subscriptions/[\w-]+/resourceGroups/[\w-]+/providers/Microsoft.KeyVault/vaults/[\w-]+/secrets/[\w-]+/versions/[\w-]+$`)),
				check.That(data.ResourceName).Key("resource_versionless_id").MatchesRegex(regexp.MustCompile(`^/subscriptions/[\w-]+/resourceGroups/[\w-]+/providers/Microsoft.KeyVault/vaults/[\w-]+/secrets/[\w-]+$`)),
			),
		},
	})
}

func TestAccDataSourceKeyVaultSecret_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_secret", "test")
	r := KeyVaultSecretDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("value").HasValue("<rick><morty /></rick>"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.hello").HasValue("world"),
				check.That(data.ResourceName).Key("versionless_id").HasValue(fmt.Sprintf("https://acctestkv-%s.vault.azure.net/secrets/secret-%s", data.RandomString, data.RandomString)),
				check.That(data.ResourceName).Key("not_before_date").HasValue("2019-01-01T01:02:03Z"),
				check.That(data.ResourceName).Key("expiration_date").HasValue("2020-01-01T01:02:03Z"),
			),
		},
	})
}

func TestAccDataSourceKeyVaultSecret_specifyVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_secret", "test")
	r := KeyVaultSecretDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("value").HasValue("rick-and-morty"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("resource_id").MatchesRegex(regexp.MustCompile(`^/subscriptions/[\w-]+/resourceGroups/[\w-]+/providers/Microsoft.KeyVault/vaults/[\w-]+/secrets/[\w-]+/versions/[\w-]+$`)),
				check.That(data.ResourceName).Key("resource_versionless_id").MatchesRegex(regexp.MustCompile(`^/subscriptions/[\w-]+/resourceGroups/[\w-]+/providers/Microsoft.KeyVault/vaults/[\w-]+/secrets/[\w-]+$`)),
			),
		},
		{
			Config: r.specifyOldVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("value").HasValue("rick-and-morty"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("resource_id").MatchesRegex(regexp.MustCompile(`^/subscriptions/[\w-]+/resourceGroups/[\w-]+/providers/Microsoft.KeyVault/vaults/[\w-]+/secrets/[\w-]+/versions/[\w-]+$`)),
				check.That(data.ResourceName).Key("resource_versionless_id").MatchesRegex(regexp.MustCompile(`^/subscriptions/[\w-]+/resourceGroups/[\w-]+/providers/Microsoft.KeyVault/vaults/[\w-]+/secrets/[\w-]+$`)),
			),
		},
	})
}

func (KeyVaultSecretDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_secret" "test" {
  name         = azurerm_key_vault_secret.test.name
  key_vault_id = azurerm_key_vault.test.id
  version      = azurerm_key_vault_secret.test.version
}
`, KeyVaultSecretResource{}.basic(data))
}

func (KeyVaultSecretDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_secret" "test" {
  name         = azurerm_key_vault_secret.test.name
  key_vault_id = azurerm_key_vault.test.id
}
`, KeyVaultSecretResource{}.complete(data))
}

func (KeyVaultSecretDataSource) specifyOldVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`


data "azurerm_key_vault_secret" "test" {
  name         = "${local.secret_name}"
  key_vault_id = azurerm_key_vault.test.id
}

locals {
  old_version = "${data.azurerm_key_vault_secret.test.version}"
  secret_name = "secret-%s"
}

%s

data "azurerm_key_vault_secret" "test2" {
  name         = azurerm_key_vault_secret.test.name
  key_vault_id = azurerm_key_vault.test.id
  version      = "${local.old_version}"
  depends_on   = [data.azurerm_key_vault_secret.test]
}


`, data.RandomString, KeyVaultSecretResource{}.basicUpdated(data))
}
