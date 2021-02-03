package databoxedge_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type DataboxEdgeOrderDataSource struct {
}

func TestAccDataSourceDataboxEdgeOrder_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_databox_edge_order", "test")
	r := DataboxEdgeOrderDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("device_name").Exists(),
			),
		},
	},
	)
}

func (DataboxEdgeOrderDataSource) basic(data acceptance.TestData) string {
	r := DataboxEdgeOrderResource{}
	return fmt.Sprintf(`
%s

data "azurerm_databox_edge_order" "test" {
  resource_group_name = azurerm_databox_edge_order.test.resource_group_name
  device_name         = azurerm_databox_edge_order.test.device_name
}
`, r.basic(data))
}
