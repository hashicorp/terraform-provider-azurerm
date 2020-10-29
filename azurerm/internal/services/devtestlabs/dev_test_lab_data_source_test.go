package devtestlabs_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMDevTestLab_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dev_test_lab", "test")

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDevTestLab_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "storage_type", "Premium"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMDevTestLab_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dev_test_lab", "test")

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDevTestLab_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "storage_type", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "World"),
				),
			},
		},
	})
}

func testAccDataSourceDevTestLab_basic(data acceptance.TestData) string {
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

func testAccDataSourceDevTestLab_complete(data acceptance.TestData) string {
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
