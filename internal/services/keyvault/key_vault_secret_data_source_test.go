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
