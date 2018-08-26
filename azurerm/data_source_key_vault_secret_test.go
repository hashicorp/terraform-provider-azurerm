package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMKeyVaultSecret_basic(t *testing.T) {
	dataSourceName := "data.azurerm_key_vault_secret.test"

	rString := acctest.RandString(8)
	location := testLocation()
	config := testAccDataSourceKeyVaultSecret_basic(rString, location)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "value", "rick-and-morty"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVaultSecret_complete(t *testing.T) {
	dataSourceName := "data.azurerm_key_vault_secret.test"

	rString := acctest.RandString(8)
	location := testLocation()
	config := testAccDataSourceKeyVaultSecret_complete(rString, location)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "value", "<rick><morty /></rick>"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.hello", "world"),
				),
			},
		},
	})
}

func testAccDataSourceKeyVaultSecret_basic(rString string, location string) string {
	resource := testAccAzureRMKeyVaultSecret_basic(rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_secret" "test" {
  name 		= "${azurerm_key_vault_secret.test.name}"
  vault_uri = "${azurerm_key_vault_secret.test.vault_uri}"
}
`, resource)
}

func testAccDataSourceKeyVaultSecret_complete(rString string, location string) string {
	resource := testAccAzureRMKeyVaultSecret_complete(rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_secret" "test" {
  name 		= "${azurerm_key_vault_secret.test.name}"
  vault_uri = "${azurerm_key_vault_secret.test.vault_uri}"
}
`, resource)
}
