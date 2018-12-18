package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceArmManagementGroup_basic(t *testing.T) {
	dataSourceName := "data.azurerm_management_group.test"
	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArmManagementGroup_basic(ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "display_name", fmt.Sprintf("acctestmg-%d", ri)),
					resource.TestCheckResourceAttr(dataSourceName, "subscription_ids.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceArmManagementGroup_basic(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_management_group" "test" {
  display_name = "acctestmg-%d"
}

data "azurerm_management_group" "test" {
  group_id = "${azurerm_management_group.test.group_id}"
}
`, rInt)
}
