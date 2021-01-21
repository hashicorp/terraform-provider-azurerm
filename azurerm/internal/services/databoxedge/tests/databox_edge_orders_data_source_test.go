package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceDataboxEdgeOrder_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_databox_edge_order", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataboxEdgeOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcedataboxedgeOrder_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeOrderExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourcedataboxedgeOrder_basic(data acceptance.TestData) string {
	config := testAccDataboxEdgeOrder_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_databox_edge_order" "test" {
  resource_group_name = azurerm_databox_edge_order.test.resource_group_name
  name = azurerm_databox_edge_order.test.name
}
`, config)
}
