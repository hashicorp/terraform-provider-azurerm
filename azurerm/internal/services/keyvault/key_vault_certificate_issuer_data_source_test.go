package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type KeyVaultCertificateIssuerDataSource struct {
}

func TestAccDataSourceKeyVaultCertificateIssuer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_certificate_issuer", "test")
	r := KeyVaultCertificateIssuerDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("account_id").HasValue("test-account"),
				check.That(data.ResourceName).Key("provider_name").HasValue("DigiCert"),
				check.That(data.ResourceName).Key("org_id").HasValue("accTestOrg"),
				check.That(data.ResourceName).Key("admin.0.email_address").HasValue("admin@contoso.com"),
				check.That(data.ResourceName).Key("admin.0.first_name").HasValue("First"),
				check.That(data.ResourceName).Key("admin.0.last_name").HasValue("Last"),
				check.That(data.ResourceName).Key("admin.0.phone").HasValue("01234567890"),
			),
		},
	})
}

func (KeyVaultCertificateIssuerDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_certificate_issuer" "test" {
  name         = azurerm_key_vault_certificate_issuer.test.name
  key_vault_id = azurerm_key_vault.test.id
}
`, KeyVaultCertificateIssuerResource{}.complete(data))
}
