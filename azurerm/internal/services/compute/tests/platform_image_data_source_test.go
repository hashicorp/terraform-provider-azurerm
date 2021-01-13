package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type PlatformImageDataSource struct {
}

func TestAccDataSourcePlatformImage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_platform_image", "test")
	r := PlatformImageDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("sku").HasValue("16.04-LTS"),
			),
		},
	})
}

func TestAccDataSourcePlatformImage_withVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_platform_image", "test")
	r := PlatformImageDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.withVersion(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("offer").HasValue("UbuntuServer"),
				check.That(data.ResourceName).Key("sku").HasValue("16.04-LTS"),
				check.That(data.ResourceName).Key("version").HasValue("16.04.201811010"),
			),
		},
	})
}

func (PlatformImageDataSource) basic(data acceptance.TestData) string {
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

func (PlatformImageDataSource) withVersion(data acceptance.TestData) string {
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
