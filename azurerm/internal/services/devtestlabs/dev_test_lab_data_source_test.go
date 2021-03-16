package devtestlabs_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AzureRMDevTestLabDataSource struct {
}

func TestAccAzureRMDevTestLabDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dev_test_lab", "test")
	r := AzureRMDevTestLabDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("storage_type").HasValue("Premium"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccAzureRMDevTestLabDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dev_test_lab", "test")
	r := AzureRMDevTestLabDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("storage_type").HasValue("Standard"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Hello").HasValue("World"),
			),
		},
	})
}

func (AzureRMDevTestLabDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_dev_test_lab" "test" {
  name                = azurerm_dev_test_lab.test.name
  resource_group_name = azurerm_dev_test_lab.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AzureRMDevTestLabDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  storage_type        = "Standard"

  tags = {
    Hello = "World"
  }
}

data "azurerm_dev_test_lab" "test" {
  name                = azurerm_dev_test_lab.test.name
  resource_group_name = azurerm_dev_test_lab.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
