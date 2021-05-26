package iothub_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type IotHubDataSource struct {
}

func TestAccDataSourceIotHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_iothub", "test")
	r := IotHubDataSource{}

	name := fmt.Sprintf("acctestiothub-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(resourceGroupName),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("acctest"),
			),
		},
	})
}

func (IotHubDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestiothub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  tags = {
    environment = "acctest"
  }
}

data "azurerm_iothub" "test" {
  name                = azurerm_iothub.test.name
  resource_group_name = azurerm_resource_group.test.name
}
	`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
