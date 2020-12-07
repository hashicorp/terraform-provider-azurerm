package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMKeyVaultCertificateIssuer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_certificate_issuer", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKeyVaultCertificateIssuer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "account_id", "test-account"),
					resource.TestCheckResourceAttr(data.ResourceName, "provider_name", "DigiCert"),
					resource.TestCheckResourceAttr(data.ResourceName, "org_id", "accTestOrg"),
					resource.TestCheckResourceAttr(data.ResourceName, "admin.0.email_address", "admin@contoso.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "admin.0.first_name", "First"),
					resource.TestCheckResourceAttr(data.ResourceName, "admin.0.last_name", "Last"),
					resource.TestCheckResourceAttr(data.ResourceName, "admin.0.phone", "01234567890"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKeyVaultCertificateIssuer_basic(data acceptance.TestData) string {
	template := testAccAzureRMKeyVaultCertificateIssuer_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_certificate_issuer" "test" {
  name         = azurerm_key_vault_certificate_issuer.test.name
  key_vault_id = azurerm_key_vault.test.id
}
`, template)
}
