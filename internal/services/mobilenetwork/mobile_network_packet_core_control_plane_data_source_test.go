// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MobileNetworkPacketCoreControlPlanDataSource struct{}

func TestAccMobileNetworkPacketCoreControlPlanDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_control_plane", "test")
	d := MobileNetworkPacketCoreControlPlanDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key(`location`).Exists(),
				check.That(data.ResourceName).Key(`sku`).HasValue("G0"),
				check.That(data.ResourceName).Key(`site_ids.#`).HasValue("1"),
				check.That(data.ResourceName).Key(`local_diagnostics_access.0.authentication_type`).HasValue("AAD"),
				check.That(data.ResourceName).Key(`platform.0.type`).HasValue("AKS-HCI"),
			),
		},
	})
}

func (r MobileNetworkPacketCoreControlPlanDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

data "azurerm_mobile_network_packet_core_control_plane" "test" {
  name                = azurerm_mobile_network_packet_core_control_plane.test.name
  resource_group_name = azurerm_mobile_network_packet_core_control_plane.test.resource_group_name
}
`, MobileNetworkPacketCoreControlPlaneResource{}.basic(data))
}
