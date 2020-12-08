package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMKeyVaultCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKeyVaultCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "certificate_data"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.key_properties.0.key_size", "2048"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.key_properties.0.key_type", "RSA"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVaultCertificate_generated(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKeyVaultCertificate_generated(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "certificate_data"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.issuer_parameters.0.name", "Self"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.key_properties.0.exportable", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.key_properties.0.key_size", "2048"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.key_properties.0.key_type", "RSA"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.key_properties.0.reuse_key", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.lifetime_action.0.action.0.action_type", "AutoRenew"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.lifetime_action.0.trigger.0.days_before_expiry", "30"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.secret_properties.0.content_type", "application/x-pkcs12"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.x509_certificate_properties.0.subject", "CN=hello-world"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_policy.0.x509_certificate_properties.0.validity_in_months", "12"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKeyVaultCertificate_basic(data acceptance.TestData) string {
	template := testAccAzureRMKeyVaultCertificate_basicImportPFX(data)
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_certificate" "test" {
  name         = azurerm_key_vault_certificate.test.name
  key_vault_id = azurerm_key_vault.test.id
}
`, template)
}

func testAccDataSourceAzureRMKeyVaultCertificate_generated(data acceptance.TestData) string {
	template := testAccAzureRMKeyVaultCertificate_basicGenerate(data)
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_certificate" "test" {
  name         = azurerm_key_vault_certificate.test.name
  key_vault_id = azurerm_key_vault.test.id
}
`, template)
}
