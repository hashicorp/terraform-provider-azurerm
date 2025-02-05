// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package orbital_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/contactprofile"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContactProfileResource struct{}

func TestAccContactProfile_basic(t *testing.T) {
	if features.FivePointOh() {
		t.Skipf("Skipping since `azurerm_orbital_contact_profile` is deprecated and will be removed in 5.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_orbital_contact_profile", "test")
	r := ContactProfileResource{}

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

func TestAccContactProfile_multipleChannels(t *testing.T) {
	if features.FivePointOh() {
		t.Skipf("Skipping since `azurerm_orbital_contact_profile` is deprecated and will be removed in 5.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_orbital_contact_profile", "test")
	r := ContactProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleChannels(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContactProfile_addChannel(t *testing.T) {
	if features.FivePointOh() {
		t.Skipf("Skipping since `azurerm_orbital_contact_profile` is deprecated and will be removed in 5.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_orbital_contact_profile", "test")
	r := ContactProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleChannels(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContactProfile_update(t *testing.T) {
	if features.FivePointOh() {
		t.Skipf("Skipping since `azurerm_orbital_contact_profile` is deprecated and will be removed in 5.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_orbital_contact_profile", "test")
	r := ContactProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContactProfile_complete(t *testing.T) {
	if features.FivePointOh() {
		t.Skipf("Skipping since `azurerm_orbital_contact_profile` is deprecated and will be removed in 5.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_orbital_contact_profile", "test")
	r := ContactProfileResource{}

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

func (r ContactProfileResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := contactprofile.ParseContactProfileID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Orbital.ContactProfileClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retreiving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r ContactProfileResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_orbital_contact_profile" "test" {
  name                              = "testcontactprofile-%[2]d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  minimum_variable_contact_duration = "PT1M"
  auto_tracking                     = "disabled"
  links {
    channels {
      name                 = "channelname"
      bandwidth_mhz        = 100
      center_frequency_mhz = 101
      end_point {
        end_point_name = "AQUA_command"
        port           = "49513"
        protocol       = "TCP"
      }
    }
    direction    = "Uplink"
    name         = "RHCP_UL"
    polarization = "RHCP"
  }
  network_configuration_subnet_id = azurerm_subnet.test.id
}
`, template, data.RandomInteger)
}

func (r ContactProfileResource) multipleChannels(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_orbital_contact_profile" "test" {
  name                              = "testcontactprofile-%[2]d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  minimum_variable_contact_duration = "PT1M"
  auto_tracking                     = "disabled"
  links {
    channels {
      name                       = "channelname"
      bandwidth_mhz              = 100
      center_frequency_mhz       = 101
      demodulation_configuration = "aqua_direct_broadcast"
      modulation_configuration   = "AQUA_UPLINK_BPSK"
      end_point {
        end_point_name = "AQUA_command"
        port           = "49513"
        protocol       = "TCP"
      }
    }
    channels {
      name                     = "channelname2"
      bandwidth_mhz            = 102
      center_frequency_mhz     = 103
      modulation_configuration = "AQUA_UPLINK_BPSK"
      end_point {
        end_point_name = "AQUA_command"
        port           = "49514"
        protocol       = "TCP"
      }
    }
    direction    = "Uplink"
    name         = "RHCP_UL"
    polarization = "RHCP"
  }
  network_configuration_subnet_id = azurerm_subnet.test.id
}
`, template, data.RandomInteger)
}

func (r ContactProfileResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_orbital_contact_profile" "test" {
  name                              = "testcontactprofile-%[2]d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  minimum_variable_contact_duration = "PT2M"
  auto_tracking                     = "disabled"
  links {
    channels {
      name                 = "channelname"
      bandwidth_mhz        = 102
      center_frequency_mhz = 103
      end_point {
        end_point_name = "AQUA_command"
        ip_address     = "10.0.1.0"
        port           = "49515"
        protocol       = "TCP"
      }
    }
    direction    = "Downlink"
    name         = "RHCP_UL"
    polarization = "RHCP"
  }
  network_configuration_subnet_id = azurerm_subnet.test.id
}
`, template, data.RandomInteger)
}

func (r ContactProfileResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_eventhub_namespace" "test" {
  name                = "eventhubtest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = 1

  tags = {
    environment = "Production"
  }
}

resource "azurerm_eventhub" "test" {
  name                = "testeventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

data "azuread_service_principal" "test" {
  display_name = "Azure Orbital Resource Provider"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_eventhub.test.id
  role_definition_name = "Azure Event Hubs Data Sender"
  principal_id         = data.azuread_service_principal.test.object_id
}

resource "azurerm_orbital_contact_profile" "test" {
  name                              = "testcontactprofile-%d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  minimum_variable_contact_duration = "PT1M"
  auto_tracking                     = "disabled"
  event_hub_uri                     = azurerm_eventhub.test.id
  minimum_elevation_degrees         = 12.12
  links {
    channels {
      name                       = "channelname"
      bandwidth_mhz              = 100
      center_frequency_mhz       = 101
      demodulation_configuration = "aqua_direct_broadcast"
      modulation_configuration   = "AQUA_UPLINK_BPSK"
      end_point {
        end_point_name = "AQUA_command"
        ip_address     = "10.0.1.0"
        port           = "49513"
        protocol       = "TCP"
      }
    }
    direction    = "Uplink"
    name         = "RHCP_UL"
    polarization = "RHCP"
  }
  network_configuration_subnet_id = azurerm_subnet.test.id

  depends_on = [azurerm_role_assignment.test]
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ContactProfileResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "testvnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "orbitalgateway"

    service_delegation {
      name = "Microsoft.Orbital/orbitalGateways"
      actions = [
        "Microsoft.Network/publicIPAddresses/join/action",
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/read",
        "Microsoft.Network/publicIPAddresses/read",
      ]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
