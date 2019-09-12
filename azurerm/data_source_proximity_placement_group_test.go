package azurerm

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccProximityPlacementGroupDataSource_basic(t *testing.T) {
	dataSourceName := "data.azurerm_proximity_placement_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProximityPlacementGroupDataSource_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "2"),
				),
			},
		},
	})
}

func testAccProximityPlacementGroupDataSource_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%s

data "azurerm_proximity_placement_group" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "${azurerm_proximity_placement_group.test.name}"
}
`, testAccProximityPlacementGroup_withTags(rInt, location))
}
