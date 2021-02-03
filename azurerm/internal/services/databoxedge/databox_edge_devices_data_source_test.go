package databoxedge_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type DataboxEdgeDeviceDataSource struct {
}

func TestAccDataSourceDataboxEdgeDevice_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_databox_edge_device", "test")
	r := DataboxEdgeDeviceDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("Test"),
			),
		},
	})
}

func (DataboxEdgeDeviceDataSource) basic(data acceptance.TestData) string {
	r := DataboxEdgeDeviceResource{}
	return fmt.Sprintf(`
%s

data "azurerm_databox_edge_device" "test" {
  name                = azurerm_databox_edge_device.test.name
  resource_group_name = azurerm_databox_edge_device.test.resource_group_name
}
`, r.complete(data))
}
