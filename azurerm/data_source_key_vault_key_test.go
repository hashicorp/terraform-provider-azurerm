package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMKeyVaultKey_complete(t *testing.T) {
	dataSourceName := "data.azurerm_key_vault_key.test"

	rString := acctest.RandString(8)
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultKey_complete(rString, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "key_type", "RSA"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.hello", "world"),
				),
			},
		},
	})
}

func testAccDataSourceKeyVaultKey_complete(rString string, location string) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_key" "test" {
  name 		= "${azurerm_key_vault_key.test.name}"
  vault_uri = "${azurerm_key_vault_key.test.vault_uri}"
}
`, testAccAzureRMKeyVaultKey_complete(rString, location))
}
