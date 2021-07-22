package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type KeyVaultCertificateDataSource struct {
}

func TestAccDataSourceKeyVaultCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.key_size").HasValue("2048"),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.key_type").HasValue("RSA"),
			),
		},
	})
}

func TestAccDataSourceKeyVaultCertificate_generated(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.generated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
				check.That(data.ResourceName).Key("certificate_policy.0.issuer_parameters.0.name").HasValue("Self"),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.exportable").HasValue("true"),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.key_size").HasValue("2048"),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.key_type").HasValue("RSA"),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.reuse_key").HasValue("true"),
				check.That(data.ResourceName).Key("certificate_policy.0.lifetime_action.0.action.0.action_type").HasValue("AutoRenew"),
				check.That(data.ResourceName).Key("certificate_policy.0.lifetime_action.0.trigger.0.days_before_expiry").HasValue("30"),
				check.That(data.ResourceName).Key("certificate_policy.0.secret_properties.0.content_type").HasValue("application/x-pkcs12"),
				check.That(data.ResourceName).Key("certificate_policy.0.x509_certificate_properties.0.subject").HasValue("CN=hello-world"),
				check.That(data.ResourceName).Key("certificate_policy.0.x509_certificate_properties.0.validity_in_months").HasValue("12"),
			),
		},
	})
}

func TestAccDataSourceKeyVaultCertificate_generatedEllipticCurve(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.generatedEllipticCurve(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("certificate_data").Exists(),
				check.That(data.ResourceName).Key("certificate_data_base64").Exists(),
				check.That(data.ResourceName).Key("certificate_policy.0.issuer_parameters.0.name").HasValue("Self"),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.curve").HasValue("P-256K"),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.exportable").HasValue("true"),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.key_size").HasValue("256"),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.key_type").HasValue("EC"),
				check.That(data.ResourceName).Key("certificate_policy.0.key_properties.0.reuse_key").HasValue("true"),
				check.That(data.ResourceName).Key("certificate_policy.0.lifetime_action.0.action.0.action_type").HasValue("AutoRenew"),
				check.That(data.ResourceName).Key("certificate_policy.0.lifetime_action.0.trigger.0.days_before_expiry").HasValue("30"),
				check.That(data.ResourceName).Key("certificate_policy.0.secret_properties.0.content_type").HasValue("application/x-pkcs12"),
				check.That(data.ResourceName).Key("certificate_policy.0.x509_certificate_properties.0.subject").HasValue("CN=hello-world"),
				check.That(data.ResourceName).Key("certificate_policy.0.x509_certificate_properties.0.validity_in_months").HasValue("12"),
			),
		},
	})
}

func (KeyVaultCertificateDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_certificate" "test" {
  name         = azurerm_key_vault_certificate.test.name
  key_vault_id = azurerm_key_vault.test.id
}
`, KeyVaultCertificateResource{}.basicImportPFX(data))
}

func (KeyVaultCertificateDataSource) generated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_certificate" "test" {
  name         = azurerm_key_vault_certificate.test.name
  key_vault_id = azurerm_key_vault.test.id
}
`, KeyVaultCertificateResource{}.basicGenerate(data))
}

func (KeyVaultCertificateDataSource) generatedEllipticCurve(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_certificate" "test" {
  name         = azurerm_key_vault_certificate.test.name
  key_vault_id = azurerm_key_vault.test.id
}
`, KeyVaultCertificateResource{}.basicGenerateEllipticCurve(data))
}
