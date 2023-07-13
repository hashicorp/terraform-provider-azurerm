// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databoxedge_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DataboxEdgeDeviceDataSource struct{}

func TestAccDataboxEdgeDeviceDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_edge_device", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: DataboxEdgeDeviceDataSource{}.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("sku_name").HasValue("EdgeP_Base-Standard"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		data.ImportStep(),
	},
	)
}

func (r DataboxEdgeDeviceDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_databox_edge_device" "test" {
  name                = azurerm_databox_edge_device.test.name
  resource_group_name = azurerm_databox_edge_device.test.resource_group_name
}

`, DataboxEdgeDeviceResource{}.complete(data))
}
