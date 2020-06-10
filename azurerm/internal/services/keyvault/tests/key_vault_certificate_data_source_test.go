package tests

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
