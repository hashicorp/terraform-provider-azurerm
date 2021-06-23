package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type KeyVaultSecretsDataSource struct {
}

func TestAccDataSourceKeyVaultSecrets_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_secrets", "test")
	r := KeyVaultSecretsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("names.#").HasValue("31"),
			),
		},
	})
}

func (KeyVaultSecretsDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_secret" "test2" {
  count        = 30
  name         = "secret-${count.index}"
  value        = "rick-and-morty"
  key_vault_id = azurerm_key_vault.test.id
}

data "azurerm_key_vault_secrets" "test" {
  key_vault_id = azurerm_key_vault.test.id

  depends_on = [azurerm_key_vault_secret.test, azurerm_key_vault_secret.test2]
}
`, KeyVaultSecretResource{}.basic(data))
}
