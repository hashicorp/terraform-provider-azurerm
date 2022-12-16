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

func TestAccMobileNetworkAttachedDataNetwork_withDataInterface(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_attached_data_network", "test")
	r := MobileNetworkAttachedDataNetworkResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withDataInterface(data),
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
			Config: r.basic(data),
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
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_attached_data_network" "test" {
  name                                     = "acctest-mnadn-%d"
  mobile_network_packet_core_data_plane_id = azurerm_mobile_network_packet_core_data_plane.test.id
  location                                 = "%s"
  dns_addresses                            = ["1.1.1.1"]

  depends_on = [azurerm_mobile_network_data_network.test]
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkAttachedDataNetworkResource) withDataInterface(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_attached_data_network" "test" {
  name                                     = "acctest-mnadn-%d"
  mobile_network_packet_core_data_plane_id = azurerm_mobile_network_packet_core_data_plane.test.id
  location                                 = "%s"
  dns_addresses                            = ["1.1.1.1"]

  user_plane_data_interface {
    name         = "test"
    ipv4_address = "10.204.141.4"
    ipv4_gateway = "10.204.141.1"
    ipv4_subnet  = "10.204.141.0/24"
  }

  depends_on = [azurerm_mobile_network_data_network.test]
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkAttachedDataNetworkResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_attached_data_network" "import" {
  name                                     = azurerm_mobile_network_attached_data_network.test.name
  mobile_network_packet_core_data_plane_id = azurerm_mobile_network_packet_core_data_plane.test.id
  location                                 = "%s"
  dns_addresses                            = azurerm_mobile_network_attached_data_network.test.dns_addresses

  user_plane_data_interface {
    name         = "test"
    ipv4_address = "10.204.141.4"
    ipv4_gateway = "10.204.141.1"
    ipv4_subnet  = "10.204.141.0/24"
  }

  depends_on = [azurerm_mobile_network_data_network.test]
}
`, config, data.Locations.Primary)
}

func (r MobileNetworkAttachedDataNetworkResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_attached_data_network" "test" {
  name                                        = "acctest-mnadn-%d"
  mobile_network_packet_core_data_plane_id    = azurerm_mobile_network_packet_core_data_plane.test.id
  location                                    = "%s"
  dns_addresses                               = ["1.1.1.1"]
  user_equipment_address_pool_prefixes        = ["2.4.0.0/16"]
  user_equipment_static_address_pool_prefixes = ["2.4.0.0/16"]
  network_address_port_translation_configuration {
    enabled                = true
    pinhole_maximum_number = 65536

    pinhole_timeouts_in_seconds {
      icmp = 30
      tcp  = 100
      udp  = 39
    }

    port_range {
      max_port = 49999
      min_port = 1024
    }
    port_reuse_minimum_hold_time_in_seconds {
      tcp = 120
      udp = 60
    }
  }

  user_plane_data_interface {
    name         = "test"
    ipv4_address = "10.204.141.4"
    ipv4_gateway = "10.204.141.1"
    ipv4_subnet  = "10.204.141.0/24"
  }

  tags = {
    key = "value"
  }

  depends_on = [azurerm_mobile_network_data_network.test]

}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkAttachedDataNetworkResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_attached_data_network" "test" {
  name                                        = "acctest-mnadn-%d"
  mobile_network_packet_core_data_plane_id    = azurerm_mobile_network_packet_core_data_plane.test.id
  location                                    = "%s"
  dns_addresses                               = ["1.1.1.1"]
  user_equipment_address_pool_prefixes        = ["2.4.0.0/16"]
  user_equipment_static_address_pool_prefixes = ["2.4.0.0/16"]
  network_address_port_translation_configuration {
    enabled                = true
    pinhole_maximum_number = 65536
    pinhole_timeouts_in_seconds {
      icmp = 30
      tcp  = 100
      udp  = 39
    }
    port_range {
      max_port = 49999
      min_port = 1024
    }
    port_reuse_minimum_hold_time_in_seconds {
      tcp = 120
      udp = 60
    }
  }
  user_plane_data_interface {
    name         = "test"
    ipv4_address = "10.204.141.4"
    ipv4_gateway = "10.204.141.1"
    ipv4_subnet  = "10.204.141.0/24"
  }
  tags = {
    key = "value"
  }

  depends_on = [azurerm_mobile_network_data_network.test]

}
`, template, data.RandomInteger, data.Locations.Primary)
}
