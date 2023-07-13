// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type VirtualHubDataSource struct{}

func TestAccDataSourceAzureRMVirtualHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_hub", "test")
	r := VirtualHubDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("address_prefix").Exists(),
				check.That(data.ResourceName).Key("virtual_wan_id").Exists(),
				check.That(data.ResourceName).Key("virtual_router_asn").Exists(),
				check.That(data.ResourceName).Key("virtual_router_ips.#").Exists(),
			),
		},
	})
}

func (VirtualHubDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_virtual_hub" "test" {
  name                = azurerm_virtual_hub.test.name
  resource_group_name = azurerm_virtual_hub.test.resource_group_name
}
`, VirtualHubResource{}.basic(data))
}
