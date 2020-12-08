package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMPlatformImage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_platform_image", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMPlatformImageBasic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "version"),
					resource.TestCheckResourceAttr(data.ResourceName, "publisher", "Canonical"),
					resource.TestCheckResourceAttr(data.ResourceName, "offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "16.04-LTS"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMPlatformImage_withVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_platform_image", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMPlatformImageWithVersion(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "version"),
					resource.TestCheckResourceAttr(data.ResourceName, "publisher", "Canonical"),
					resource.TestCheckResourceAttr(data.ResourceName, "offer", "UbuntuServer"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "16.04-LTS"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "16.04.201811010"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMPlatformImageBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_platform_image" "test" {
  location  = "%s"
  publisher = "Canonical"
  offer     = "UbuntuServer"
  sku       = "16.04-LTS"
}
`, data.Locations.Primary)
}

func testAccDataSourceAzureRMPlatformImageWithVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_platform_image" "test" {
  location  = "%s"
  publisher = "Canonical"
  offer     = "UbuntuServer"
  sku       = "16.04-LTS"
  version   = "16.04.201811010"
}
`, data.Locations.Primary)
}
