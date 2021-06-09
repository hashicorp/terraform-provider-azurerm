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

func TestAccDataSourceKeyVaultCertificateData_ecdsa_PFX(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_certificate_data", "test")
	r := KeyVaultCertificateDataDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.ecdsa_PFX(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("hex").Exists(),
				check.That(data.ResourceName).Key("pem").Exists(),
				check.That(data.ResourceName).Key("key").Exists(),
			),
		},
	})
}

func TestAccDataSourceKeyVaultCertificateData_ecdsa_PEM(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_certificate_data", "test")
	r := KeyVaultCertificateDataDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.ecdsa_PEM(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("hex").Exists(),
				check.That(data.ResourceName).Key("pem").Exists(),
				check.That(data.ResourceName).Key("key").Exists(),
			),
		},
	})
}

func TestAccDataSourceKeyVaultCertificateData_rsa_bundle_PEM(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_certificate_data", "test")
	r := KeyVaultCertificateDataDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.rsa_bundle_PEM(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("hex").Exists(),
				check.That(data.ResourceName).Key("pem").Exists(),
				check.That(data.ResourceName).Key("key").Exists(),
				check.That(data.ResourceName).Key("certificates_count").HasValue("2"),
			),
		},
	})
}

func TestAccDataSourceKeyVaultCertificateData_rsa_single_PEM(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_certificate_data", "test")
	r := KeyVaultCertificateDataDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.rsa_single_PEM(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("hex").Exists(),
				check.That(data.ResourceName).Key("pem").Exists(),
				check.That(data.ResourceName).Key("key").Exists(),
				check.That(data.ResourceName).Key("certificates_count").HasValue("1"),
			),
		},
	})
}

func TestAccDataSourceKeyVaultCertificateData_rsa_bundle_PFX(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_certificate_data", "test")
	r := KeyVaultCertificateDataDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.rsa_bundle_PFX(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("hex").Exists(),
				check.That(data.ResourceName).Key("pem").Exists(),
				check.That(data.ResourceName).Key("key").Exists(),
				check.That(data.ResourceName).Key("certificates_count").HasValue("2"),
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

func (KeyVaultCertificateDataDataSource) ecdsa_PFX(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_certificate_data" "test" {
  name         = azurerm_key_vault_certificate.test.name
  key_vault_id = azurerm_key_vault.test.id
  version      = azurerm_key_vault_certificate.test.version
}
`, KeyVaultCertificateResource{}.basicImportPFX_ECDSA(data))
}

func (KeyVaultCertificateDataDataSource) ecdsa_PEM(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_certificate_data" "test" {
  name         = azurerm_key_vault_certificate.test.name
  key_vault_id = azurerm_key_vault.test.id
  version      = azurerm_key_vault_certificate.test.version
}
`, KeyVaultCertificateResource{}.basicImportPEM_ECDSA(data))
}

func (KeyVaultCertificateDataDataSource) rsa_bundle_PEM(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_certificate_data" "test" {
  name         = azurerm_key_vault_certificate.test.name
  key_vault_id = azurerm_key_vault.test.id
  version      = azurerm_key_vault_certificate.test.version
}
`, KeyVaultCertificateResource{}.basicImportPEM_RSA_bundle(data))
}

func (KeyVaultCertificateDataDataSource) rsa_single_PEM(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_certificate_data" "test" {
  name         = azurerm_key_vault_certificate.test.name
  key_vault_id = azurerm_key_vault.test.id
  version      = azurerm_key_vault_certificate.test.version
}
`, KeyVaultCertificateResource{}.basicImportPEM_RSA(data))
}

func (KeyVaultCertificateDataDataSource) rsa_bundle_PFX(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_certificate_data" "test" {
  name         = azurerm_key_vault_certificate.test.name
  key_vault_id = azurerm_key_vault.test.id
  version      = azurerm_key_vault_certificate.test.version
}
`, KeyVaultCertificateResource{}.basicImportPFX_RSA_bundle(data))
}
