package mobilenetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/sim"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MobileNetworkSimResource struct{}

func TestAccMobileNetworkSim_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim", "test")
	r := MobileNetworkSimResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication_key", "operator_key_code"),
	})
}

func TestAccMobileNetworkSim_withStaticIpConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim", "test")
	r := MobileNetworkSimResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withStaticIpConfiguration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication_key", "operator_key_code"),
	})
}

func TestAccMobileNetworkSim_withSimPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim", "test")
	r := MobileNetworkSimResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSimPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication_key", "operator_key_code"),
	})
}

func TestAccMobileNetworkSim_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim", "test")
	r := MobileNetworkSimResource{}
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

func TestAccMobileNetworkSim_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim", "test")
	r := MobileNetworkSimResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication_key", "operator_key_code"),
	})
}

func TestAccMobileNetworkSim_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim", "test")
	r := MobileNetworkSimResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication_key", "operator_key_code"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication_key", "operator_key_code"),
	})
}

func (r MobileNetworkSimResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := sim.ParseSimID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.MobileNetwork.SIMClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r MobileNetworkSimResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-mn-%[1]d"
  location = "%[2]s"
}

resource "azurerm_databox_edge_device" "test" {
  name                = "acct%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku_name = "EdgeP_Base-Standard"
}

resource "azurerm_mobile_network" "test" {
  name                = "acctest-mn-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  mobile_country_code = "001"
  mobile_network_code = "01"
}

resource "azurerm_mobile_network_sim_group" "test" {
  name              = "acctest-mnsg-%[1]d"
  location          = azurerm_mobile_network.test.location
  mobile_network_id = azurerm_mobile_network.test.id
}

resource "azurerm_mobile_network_site" "test" {
  name              = "acctest-mns-%[1]d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = azurerm_mobile_network.test.location
}

resource "azurerm_mobile_network_packet_core_control_plane" "test" {
  name                = "acctest-mnpccp-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_mobile_network.test.location
  sku                 = "G0"
  site_ids            = [azurerm_mobile_network_site.test.id]

  local_diagnostics_access {
    authentication_type = "AAD"
  }

  platform {
    type           = "AKS-HCI"
    edge_device_id = azurerm_databox_edge_device.test.id
  }

  depends_on = [azurerm_mobile_network.test]
}

resource "azurerm_mobile_network_packet_core_data_plane" "test" {
  name                                        = "acctest-mnpcdp-%[1]d"
  mobile_network_packet_core_control_plane_id = azurerm_mobile_network_packet_core_control_plane.test.id
  location                                    = azurerm_mobile_network.test.location
}

resource "azurerm_mobile_network_data_network" "test" {
  name              = "acctest-mndn-%[1]d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = azurerm_mobile_network.test.location
}

resource "azurerm_mobile_network_attached_data_network" "test" {
  mobile_network_data_network_name            = azurerm_mobile_network_data_network.test.name
  mobile_network_packet_core_data_plane_id    = azurerm_mobile_network_packet_core_data_plane.test.id
  location                                    = azurerm_mobile_network.test.location
  dns_addresses                               = ["1.1.1.1"]
  user_equipment_address_pool_prefixes        = ["2.4.1.0/24"]
  user_equipment_static_address_pool_prefixes = ["2.4.0.0/24"]
}

resource "azurerm_mobile_network_slice" "test" {
  name              = "acctest-mns-%[1]d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = azurerm_mobile_network.test.location
  single_network_slice_selection_assistance_information {
    slice_service_type = 1
  }
}

resource "azurerm_mobile_network_service" "test" {
  name               = "acctest-mns-%[1]d"
  mobile_network_id  = azurerm_mobile_network.test.id
  location           = azurerm_mobile_network.test.location
  service_precedence = 0

  pcc_rule {
    name                    = "default-rule"
    precedence              = 1
    traffic_control_enabled = true

    service_data_flow_template {
      direction      = "Uplink"
      name           = "IP-to-server"
      ports          = []
      protocol       = ["ip"]
      remote_ip_list = ["10.3.4.0/24"]
    }
  }
}

resource "azurerm_mobile_network_sim_policy" "test" {
  name                                   = "acctest-mnsp-%[1]d"
  mobile_network_id                      = azurerm_mobile_network.test.id
  location                               = azurerm_mobile_network.test.location
  default_slice_id                       = azurerm_mobile_network_slice.test.id
  registration_timer_in_seconds          = 3240
  rat_frequency_selection_priority_index = 1

  slice {
    default_data_network_id = azurerm_mobile_network_data_network.test.id
    slice_id                = azurerm_mobile_network_slice.test.id
    data_network {
      allocation_and_retention_priority_level = 9
      default_session_type                    = "IPv4"
      qos_indicator                           = 9
      preemption_capability                   = "NotPreempt"
      preemption_vulnerability                = "Preemptable"
      allowed_services_ids                    = [azurerm_mobile_network_service.test.id]
      data_network_id                         = azurerm_mobile_network_data_network.test.id
      session_aggregate_maximum_bit_rate {
        downlink = "1 Gbps"
        uplink   = "500 Mbps"
      }
    }
  }

  user_equipment_aggregate_maximum_bit_rate {
    downlink = "1 Gbps"
    uplink   = "500 Mbps"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkSimResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%[1]s

resource "azurerm_mobile_network_sim" "test" {
  name                                     = "acctest-mns-%[2]d"
  mobile_network_sim_group_id              = azurerm_mobile_network_sim_group.test.id
  authentication_key                       = "d3c97eed4e5a00726ad8a26d5918aa2f"
  integrated_circuit_card_identifier       = "8900000000000000000"
  international_mobile_subscriber_identity = "000000000000000"
  operator_key_code                        = "d3c97eed4e5a00726ad8a26d5918aa2f"
}
`, r.template(data), data.RandomInteger)
}

func (r MobileNetworkSimResource) withSimPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%[1]s
resource "azurerm_mobile_network_sim" "test" {
  name                                     = "acctest-mns-%[2]d"
  mobile_network_sim_group_id              = azurerm_mobile_network_sim_group.test.id
  authentication_key                       = "d3c97eed4e5a00726ad8a26d5918aa2f"
  integrated_circuit_card_identifier       = "8900000000000000000"
  international_mobile_subscriber_identity = "000000000000000"
  operator_key_code                        = "d3c97eed4e5a00726ad8a26d5918aa2f"
  sim_policy_id                            = azurerm_mobile_network_sim_policy.test.id
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkSimResource) withStaticIpConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%[1]s

resource "azurerm_mobile_network_sim" "test" {
  name                                     = "acctest-mns-%[2]d"
  mobile_network_sim_group_id              = azurerm_mobile_network_sim_group.test.id
  authentication_key                       = "00000000000000000000000000000000"
  integrated_circuit_card_identifier       = "8900000000000000000"
  international_mobile_subscriber_identity = "000000000000000"
  operator_key_code                        = "00000000000000000000000000000000"
  static_ip_configuration {
    attached_data_network_id = azurerm_mobile_network_attached_data_network.test.id
    slice_id                 = azurerm_mobile_network_slice.test.id
    static_ipv4_address      = "2.4.0.1"
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkSimResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_sim" "import" {
  name                                     = azurerm_mobile_network_sim.test.name
  mobile_network_sim_group_id              = azurerm_mobile_network_sim_group.test.id
  international_mobile_subscriber_identity = azurerm_mobile_network_sim.test.international_mobile_subscriber_identity
  authentication_key                       = "00000000000000000000000000000000"
  integrated_circuit_card_identifier       = "8900000000000000000"
  operator_key_code                        = "00000000000000000000000000000000"
}


`, config)
}

func (r MobileNetworkSimResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%[1]s

resource "azurerm_mobile_network_sim" "test" {
  name                                     = "acctest-mns-%[2]d"
  mobile_network_sim_group_id              = azurerm_mobile_network_sim_group.test.id
  authentication_key                       = "00000000000000000000000000000000"
  operator_key_code                        = "00000000000000000000000000000000"
  integrated_circuit_card_identifier       = "8900000000000000000"
  international_mobile_subscriber_identity = "000000000000000"
  sim_policy_id                            = azurerm_mobile_network_sim_policy.test.id

  static_ip_configuration {
    attached_data_network_id = azurerm_mobile_network_attached_data_network.test.id
    slice_id                 = azurerm_mobile_network_slice.test.id
    static_ipv4_address      = "2.4.0.1"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MobileNetworkSimResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_sim" "test" {
  name                                     = "acctest-mns-%[2]d"
  mobile_network_sim_group_id              = azurerm_mobile_network_sim_group.test.id
  authentication_key                       = "00000000000000000000000000000001"
  integrated_circuit_card_identifier       = "8900000000000000000"
  international_mobile_subscriber_identity = "000000000000000"
  operator_key_code                        = "00000000000000000000000000000001"
  sim_policy_id                            = azurerm_mobile_network_sim_policy.test.id

  static_ip_configuration {
    attached_data_network_id = azurerm_mobile_network_attached_data_network.test.id
    slice_id                 = azurerm_mobile_network_slice.test.id
    static_ipv4_address      = "2.4.0.1"
  }
}
`, r.template(data), data.RandomInteger)
}
