package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type KeyVaultCertificateDataDataSource struct {
}

func TestAccDataSourceKeyVaultCertificateData_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_certificate_data", "test")
	r := KeyVaultCertificateDataDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("hex").Exists(),
				check.That(data.ResourceName).Key("pem").Exists(),
				check.That(data.ResourceName).Key("key").Exists(),
				check.That(data.ResourceName).Key("expires").HasValue("2027-10-08T08:27:55Z"),
			),
		},
	})
}

func (KeyVaultCertificateDataDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_certificate_data" "test" {
  name         = azurerm_key_vault_certificate.test.name
  key_vault_id = azurerm_key_vault.test.id
  version      = azurerm_key_vault_certificate.test.version
}
`, KeyVaultCertificateResource{}.basicImportPFX(data))
}
