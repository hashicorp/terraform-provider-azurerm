package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMKeyVaultSecret_basic(t *testing.T) {
	dataSourceName := "data.azurerm_key_vault_secret.test"

	rString := acctest.RandString(8)
	location := acceptance.Location()
	config := testAccDataSourceKeyVaultSecret_basic(rString, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
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
	location := acceptance.Location()
	config := testAccDataSourceKeyVaultSecret_complete(rString, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
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
	r := testAccAzureRMKeyVaultSecret_basic(rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_secret" "test" {
  name         = "${azurerm_key_vault_secret.test.name}"
  key_vault_id = "${azurerm_key_vault.test.id}"
}
`, r)
}

func testAccDataSourceKeyVaultSecret_complete(rString string, location string) string {
	r := testAccAzureRMKeyVaultSecret_complete(rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_secret" "test" {
  name         = "${azurerm_key_vault_secret.test.name}"
  key_vault_id = "${azurerm_key_vault.test.id}"
}
`, r)
}
