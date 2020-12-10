package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMKeyVaultKey_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultKey_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "key_type", "RSA"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.hello", "world"),
				),
			},
		},
	})
}

func testAccDataSourceKeyVaultKey_complete(data acceptance.TestData) string {
	t := testAccAzureRMKeyVaultKey_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_key" "test" {
  name         = azurerm_key_vault_key.test.name
  key_vault_id = azurerm_key_vault.test.id
}
`, t)
}
