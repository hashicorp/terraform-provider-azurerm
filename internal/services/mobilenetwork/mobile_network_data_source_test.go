// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MobileNetworkDataSource struct{}

func TestAccMobileNetworkDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network", "test")
	d := MobileNetworkDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key(`location`).Exists(),
				check.That(data.ResourceName).Key(`mobile_country_code`).HasValue("001"),
				check.That(data.ResourceName).Key(`mobile_network_code`).HasValue("01"),
			),
		},
		data.ImportStep(),
	})
}

func (r MobileNetworkDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

data "azurerm_mobile_network" "test" {
  name                = azurerm_mobile_network.test.name
  resource_group_name = azurerm_mobile_network.test.resource_group_name
}
`, MobileNetworkResource{}.complete(data))
}
