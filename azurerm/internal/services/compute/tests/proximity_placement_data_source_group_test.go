package tests

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccProximityPlacementGroupDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_proximity_placement_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProximityPlacementGroupDataSource_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
		},
	})
}

func testAccProximityPlacementGroupDataSource_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_proximity_placement_group" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = azurerm_proximity_placement_group.test.name
}
`, testAccProximityPlacementGroup_withTags(data))
}
