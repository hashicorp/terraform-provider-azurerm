// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcoredataplane"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MobileNetworkPacketCoreDataPlaneResource struct{}

func TestAccMobileNetworkPacketCoreDataPlane_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_data_plane", "test")
	r := MobileNetworkPacketCoreDataPlaneResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMobileNetworkPacketCoreDataPlane_withAccessInterface(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_data_plane", "test")
	r := MobileNetworkPacketCoreDataPlaneResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAccessInterface(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMobileNetworkPacketCoreDataPlane_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_data_plane", "test")
	r := MobileNetworkPacketCoreDataPlaneResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMobileNetworkPacketCoreDataPlane_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_data_plane", "test")
	r := MobileNetworkPacketCoreDataPlaneResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMobileNetworkPacketCoreDataPlane_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_data_plane", "test")
	r := MobileNetworkPacketCoreDataPlaneResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r MobileNetworkPacketCoreDataPlaneResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := packetcoredataplane.ParsePacketCoreDataPlaneID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.MobileNetwork.PacketCoreDataPlaneClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r MobileNetworkPacketCoreDataPlaneResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_packet_core_data_plane" "test" {
  name                                        = "acctest-mnpcdp-%d"
  mobile_network_packet_core_control_plane_id = azurerm_mobile_network_packet_core_control_plane.test.id
  location                                    = "%s"
}
`, MobileNetworkPacketCoreControlPlaneResource{}.basic(data), data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkPacketCoreDataPlaneResource) withAccessInterface(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_packet_core_data_plane" "test" {
  name                                        = "acctest-mnpcdp-%d"
  mobile_network_packet_core_control_plane_id = azurerm_mobile_network_packet_core_control_plane.test.id
  location                                    = "%s"

  user_plane_access_name         = "default-interface"
  user_plane_access_ipv4_address = "192.168.1.199"
  user_plane_access_ipv4_gateway = "192.168.1.1"
  user_plane_access_ipv4_subnet  = "192.168.1.0/25"


}
`, MobileNetworkPacketCoreControlPlaneResource{}.basic(data), data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkPacketCoreDataPlaneResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_packet_core_data_plane" "import" {
  name                                        = azurerm_mobile_network_packet_core_data_plane.test.name
  mobile_network_packet_core_control_plane_id = azurerm_mobile_network_packet_core_control_plane.test.id
  location                                    = "%s"

  user_plane_access_name         = "default-interface"
  user_plane_access_ipv4_address = "192.168.1.199"
  user_plane_access_ipv4_gateway = "192.168.1.1"
  user_plane_access_ipv4_subnet  = "192.168.1.0/25"
}
`, config, data.Locations.Primary)
}

func (r MobileNetworkPacketCoreDataPlaneResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_packet_core_data_plane" "test" {
  name                                        = "acctest-mnpcdp-%d"
  mobile_network_packet_core_control_plane_id = azurerm_mobile_network_packet_core_control_plane.test.id
  location                                    = "%s"
  user_plane_access_name                      = "default-interface"
  user_plane_access_ipv4_address              = "192.168.1.199"
  user_plane_access_ipv4_gateway              = "192.168.1.1"
  user_plane_access_ipv4_subnet               = "192.168.1.0/25"

  tags = {
    key = "value"
  }

}
`, MobileNetworkPacketCoreControlPlaneResource{}.basic(data), data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkPacketCoreDataPlaneResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_packet_core_data_plane" "test" {
  name                                        = "acctest-mnpcdp-%d"
  mobile_network_packet_core_control_plane_id = azurerm_mobile_network_packet_core_control_plane.test.id
  location                                    = "%s"
  user_plane_access_name                      = "default-interface"
  user_plane_access_ipv4_address              = "192.168.1.199"
  user_plane_access_ipv4_gateway              = "192.168.1.1"
  user_plane_access_ipv4_subnet               = "192.168.1.0/25"

  tags = {
    key = "value 2"
  }
}
`, MobileNetworkPacketCoreControlPlaneResource{}.basic(data), data.RandomInteger, data.Locations.Primary)
}
