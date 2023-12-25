// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/attacheddatanetwork"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MobileNetworkAttachedDataNetworkResource struct{}

func TestAccMobileNetworkAttachedDataNetwork_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_attached_data_network", "test")
	r := MobileNetworkAttachedDataNetworkResource{}
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

func TestAccMobileNetworkAttachedDataNetwork_withDataAccess(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_attached_data_network", "test")
	r := MobileNetworkAttachedDataNetworkResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withDataAccess(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMobileNetworkAttachedDataNetwork_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_attached_data_network", "test")
	r := MobileNetworkAttachedDataNetworkResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMobileNetworkAttachedDataNetwork_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_attached_data_network", "test")
	r := MobileNetworkAttachedDataNetworkResource{}
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

func TestAccMobileNetworkAttachedDataNetwork_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_attached_data_network", "test")
	r := MobileNetworkAttachedDataNetworkResource{}
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r MobileNetworkAttachedDataNetworkResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := attacheddatanetwork.ParseAttachedDataNetworkID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.MobileNetwork.AttachedDataNetworkClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r MobileNetworkAttachedDataNetworkResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mobile_network_data_network" "test" {
  name              = "acctest-mnadn-%[2]d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = azurerm_resource_group.test.location
}

`, MobileNetworkPacketCoreDataPlaneResource{}.basic(data), data.RandomInteger)
}

func (r MobileNetworkAttachedDataNetworkResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_attached_data_network" "test" {
  mobile_network_data_network_name         = azurerm_mobile_network_data_network.test.name
  mobile_network_packet_core_data_plane_id = azurerm_mobile_network_packet_core_data_plane.test.id
  location                                 = "%s"
  dns_addresses                            = ["1.1.1.1"]
  user_equipment_address_pool_prefixes     = ["2.4.0.0/16"]
}
`, r.template(data), data.Locations.Primary)
}

func (r MobileNetworkAttachedDataNetworkResource) withDataAccess(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_attached_data_network" "test" {
  mobile_network_data_network_name         = azurerm_mobile_network_data_network.test.name
  mobile_network_packet_core_data_plane_id = azurerm_mobile_network_packet_core_data_plane.test.id
  location                                 = "%s"
  dns_addresses                            = ["1.1.1.1"]
  user_equipment_address_pool_prefixes     = ["2.4.0.0/16"]
  user_plane_access_name                   = "test"
  user_plane_access_ipv4_address           = "10.204.141.4"
  user_plane_access_ipv4_gateway           = "10.204.141.1"
  user_plane_access_ipv4_subnet            = "10.204.141.0/24"
}
`, r.template(data), data.Locations.Primary)
}

func (r MobileNetworkAttachedDataNetworkResource) requiresImport(data acceptance.TestData) string {
	config := r.complete(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_attached_data_network" "import" {
  mobile_network_data_network_name         = azurerm_mobile_network_attached_data_network.test.mobile_network_data_network_name
  mobile_network_packet_core_data_plane_id = azurerm_mobile_network_attached_data_network.test.mobile_network_packet_core_data_plane_id
  location                                 = azurerm_mobile_network_attached_data_network.test.location
  dns_addresses                            = azurerm_mobile_network_attached_data_network.test.dns_addresses
  user_equipment_address_pool_prefixes     = azurerm_mobile_network_attached_data_network.test.user_equipment_address_pool_prefixes
  user_plane_access_name                   = azurerm_mobile_network_attached_data_network.test.user_plane_access_name
  user_plane_access_ipv4_address           = azurerm_mobile_network_attached_data_network.test.user_plane_access_ipv4_address
  user_plane_access_ipv4_gateway           = azurerm_mobile_network_attached_data_network.test.user_plane_access_ipv4_gateway
  user_plane_access_ipv4_subnet            = azurerm_mobile_network_attached_data_network.test.user_plane_access_ipv4_subnet
}
`, config)
}

func (r MobileNetworkAttachedDataNetworkResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_attached_data_network" "test" {
  mobile_network_data_network_name            = azurerm_mobile_network_data_network.test.name
  mobile_network_packet_core_data_plane_id    = azurerm_mobile_network_packet_core_data_plane.test.id
  location                                    = "%s"
  dns_addresses                               = ["1.1.1.1"]
  user_equipment_address_pool_prefixes        = ["2.4.1.0/24"]
  user_equipment_static_address_pool_prefixes = ["2.4.2.0/24"]
  user_plane_access_name                      = "test"
  user_plane_access_ipv4_address              = "10.204.141.4"
  user_plane_access_ipv4_gateway              = "10.204.141.1"
  user_plane_access_ipv4_subnet               = "10.204.141.0/24"

  network_address_port_translation {
    pinhole_maximum_number                      = 65536
    icmp_pinhole_timeout_in_seconds             = 30
    tcp_pinhole_timeout_in_seconds              = 100
    udp_pinhole_timeout_in_seconds              = 39
    tcp_port_reuse_minimum_hold_time_in_seconds = 120
    udp_port_reuse_minimum_hold_time_in_seconds = 60

    port_range {
      maximum = 49999
      minimum = 1024
    }
  }

  tags = {
    key = "value"
  }


}
`, r.template(data), data.Locations.Primary)
}

func (r MobileNetworkAttachedDataNetworkResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_attached_data_network" "test" {
  mobile_network_data_network_name            = azurerm_mobile_network_data_network.test.name
  mobile_network_packet_core_data_plane_id    = azurerm_mobile_network_packet_core_data_plane.test.id
  location                                    = "%s"
  dns_addresses                               = ["1.1.1.1"]
  user_equipment_address_pool_prefixes        = ["2.4.3.0/24"]
  user_equipment_static_address_pool_prefixes = ["2.4.4.0/24"]
  user_plane_access_name                      = "test"
  user_plane_access_ipv4_address              = "10.204.141.4"
  user_plane_access_ipv4_gateway              = "10.204.141.1"
  user_plane_access_ipv4_subnet               = "10.204.141.0/24"

  network_address_port_translation {
    pinhole_maximum_number                      = 65536
    icmp_pinhole_timeout_in_seconds             = 30
    tcp_pinhole_timeout_in_seconds              = 100
    udp_pinhole_timeout_in_seconds              = 39
    tcp_port_reuse_minimum_hold_time_in_seconds = 120
    udp_port_reuse_minimum_hold_time_in_seconds = 60

    port_range {
      maximum = 49999
      minimum = 1024
    }
  }

  tags = {
    key = "value"
  }


}
`, r.template(data), data.Locations.Primary)
}
