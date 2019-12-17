package azurerm

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccProximityPlacementGroupDataSource_basic(t *testing.T) {
	dataSourceName := "data.azurerm_proximity_placement_group.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
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
