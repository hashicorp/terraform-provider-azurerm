package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMPlatformImage_basic(t *testing.T) {
	dataSourceName := "data.azurerm_platform_image.test"
	config := testAccDataSourceAzureRMPlatformImageBasic(acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "version"),
					resource.TestCheckResourceAttr(dataSourceName, "publisher", "Canonical"),
					resource.TestCheckResourceAttr(dataSourceName, "offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(dataSourceName, "sku", "16.04-LTS"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMPlatformImageBasic(location string) string {
	return fmt.Sprintf(`
data "azurerm_platform_image" "test" {
  location  = "%s"
  publisher = "Canonical"
  offer     = "UbuntuServer"
  sku       = "16.04-LTS"
}
`, location)
}
