package cdn_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type CdnProfileDataSource struct{}

func TestAccCdnProfileDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_profile", "test")
	d := CdnProfileDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: d.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku").HasValue("Standard_Verizon"),
			),
		},
	})
}

func TestAccCdnProfileDataSource_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_profile", "test")
	d := CdnProfileDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: d.withTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku").HasValue("Standard_Verizon"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
				check.That(data.ResourceName).Key("tags.cost_center").HasValue("MSFT"),
			),
		},
	})
}

func (d CdnProfileDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

data "azurerm_cdn_profile" "test" {
  name                = azurerm_cdn_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (d CdnProfileDataSource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}

data "azurerm_cdn_profile" "test" {
  name                = azurerm_cdn_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
