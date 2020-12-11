package compute_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAvailabilitySet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_availability_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAvailabilitySet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func testAccDataSourceAvailabilitySet_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "foo" = "bar"
  }
}

data "azurerm_availability_set" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = azurerm_availability_set.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}
