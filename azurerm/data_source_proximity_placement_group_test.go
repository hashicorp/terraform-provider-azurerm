package azurerm

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceProximityPlacementGroup_basic(t *testing.T) {
	dataSourceName := "data.azurerm_proximity_placement_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProximityPlacementGroup_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func testAccDataSourceProximityPlacementGroup_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctestppg-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    "foo" = "bar"
  }
}

data "azurerm_proximity_placement_group" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "${azurerm_proximity_placement_group.test.name}"
}
`, rInt, location)
}
