// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MobileNetworkSliceDataSource struct{}

func TestAccMobileNetworkSliceDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_slice", "test")

	d := MobileNetworkSliceDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key(`location`).Exists(),
				check.That(data.ResourceName).Key(`description`).HasValue("my favorite slice"),
				check.That(data.ResourceName).Key(`single_network_slice_selection_assistance_information.0.slice_service_type`).HasValue("1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func (r MobileNetworkSliceDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

data "azurerm_mobile_network_slice" "test" {
  name              = azurerm_mobile_network_slice.test.name
  mobile_network_id = azurerm_mobile_network_slice.test.mobile_network_id
}
`, MobileNetworkSliceResource{}.complete(data))
}
