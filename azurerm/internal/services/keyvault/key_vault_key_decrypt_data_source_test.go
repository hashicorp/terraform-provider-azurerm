package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

func TestAccDataSourceKeyVaultKeyDecrypt_basic(t *testing.T) {
	plaintext := "testData"
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_key_decrypt", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { acceptance.PreCheck(t) },
		ProviderFactories: map[string]terraform.ResourceProviderFactory{
			"azurerm": func() (terraform.ResourceProvider, error) {
				azurerm := provider.TestAzureProvider()
				return azurerm, nil
			},
		},
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
	t := testAccKeyVaultKeyEncrypt_basic(data, plaintext)
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_key_decrypt" "test" {
  key_vault_key_id = azurerm_key_vault_key_encrypt.test.key_vault_key_id
  payload          = azurerm_key_vault_key_encrypt.test.cipher_text
  algorithm        = azurerm_key_vault_key_encrypt.test.algorithm
}
`, t)
}
