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

type KeyVaultSecretVersionsDataSource struct{}

func TestAccDataSourceKeyVaultSecretVersions_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_secret_versions", "test")
	r := KeyVaultSecretVersionsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.#").HasValue("1"),
				check.That(data.ResourceName).Key("versions.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("versions.0.created_date").MatchesRegex(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
				check.That(data.ResourceName).Key("versions.0.id").MatchesRegex(regexp.MustCompile(`^\w+$`)),
			),
		},
	})
}

func TestAccDataSourceKeyVaultSecretVersions_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_secret_versions", "test")
	r := KeyVaultSecretVersionsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.#").HasValue("1"),
				check.That(data.ResourceName).Key("versions.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("versions.0.created_date").MatchesRegex(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
				check.That(data.ResourceName).Key("versions.0.not_before_date").HasValue("2019-01-01T01:02:03Z"),
				check.That(data.ResourceName).Key("versions.0.expiration_date").HasValue("2020-01-01T01:02:03Z"),
				check.That(data.ResourceName).Key("versions.0.updated_date").MatchesRegex(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
				check.That(data.ResourceName).Key("versions.0.id").MatchesRegex(regexp.MustCompile(`^\w+$`)),
				check.That(data.ResourceName).Key("versions.0.uri").MatchesRegex(regexp.MustCompile(`^https://acctestkv-\w+.vault.azure.net/secrets/secret-\w+/\w+$`)),
			),
		},
	})
}

func (KeyVaultSecretVersionsDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_secret_versions" "test" {
  name         = azurerm_key_vault_secret.test.name
  key_vault_id = azurerm_key_vault.test.id
}
`, KeyVaultSecretResource{}.basic(data))
}

func (KeyVaultSecretVersionsDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_secret_versions" "test" {
  name         = azurerm_key_vault_secret.test.name
  key_vault_id = azurerm_key_vault.test.id
}
`, KeyVaultSecretResource{}.complete(data))
}
