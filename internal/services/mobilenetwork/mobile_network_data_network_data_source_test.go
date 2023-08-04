// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MobileNetworkDataNetworkDataSource struct{}

func TestAccMobileNetworkDataNetworkDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_data_network", "test")
	d := MobileNetworkDataNetworkDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key(`description`).HasValue("my favourite data network"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func (r MobileNetworkDataNetworkDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

data "azurerm_mobile_network_data_network" "test" {
  name              = azurerm_mobile_network_data_network.test.name
  mobile_network_id = azurerm_mobile_network_data_network.test.mobile_network_id
}
`, MobileNetworkDataNetworkResource{}.complete(data))
}
