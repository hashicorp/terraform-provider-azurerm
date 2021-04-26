package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type KeyVaultSecretDataSource struct {
}

func TestAccDataSourceKeyVaultSecret_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_secret", "test")
	r := KeyVaultSecretDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("value").HasValue("rick-and-morty"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceKeyVaultSecret_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_secret", "test")
	r := KeyVaultSecretDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
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
