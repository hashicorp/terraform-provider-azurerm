package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMKeyVaultKeyDecrypt_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_key_decrypt", "test")
	plaintext := "testData"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultKeyDecrypt_basic(data, plaintext),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "plaintext", plaintext),
				),
			},
		},
	})
}

func testAccDataSourceKeyVaultKeyDecrypt_basic(data acceptance.TestData, plaintext string) string {
	t := testAccAzureRMKeyVaultKeyEncrypt_basic(data, plaintext)
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_key_decrypt" "test" {
  key_vault_key_id = azurerm_key_vault_key_encrypt.test.key_vault_key_id
  payload          = azurerm_key_vault_key_encrypt.test.cipher_text
  algorithm        = azurerm_key_vault_key_encrypt.test.algorithm
}
`, t)
}
