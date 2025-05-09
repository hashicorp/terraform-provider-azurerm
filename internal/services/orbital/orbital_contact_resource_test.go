// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package orbital_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/contact"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContactResource struct{}

func TestAccContact_basic(t *testing.T) {
	if features.FivePointOh() {
		t.Skipf("Skipping since `azurerm_orbital_contact` is deprecated and will be removed in 5.0")
	}
	skipContact(t)

	data := acceptance.BuildTestData(t, "azurerm_orbital_contact", "test")
	r := ContactResource{}

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

func (r ContactResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := contact.ParseContactID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Orbital.ContactClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r ContactResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_orbital_contact" "test" {
  name                   = "testcontact-%[2]d"
  spacecraft_id          = "%[3]s"
  reservation_start_time = "2025-07-16T20:35:00Z"
  reservation_end_time   = "2025-07-16T20:55:00Z"
  ground_station_name    = "Microsoft_Quincy"
  contact_profile_id     = azurerm_orbital_contact_profile.test.id
}
`, template, data.RandomInteger, os.Getenv("ARM_TEST_SPACECRAFT_ID"))
}

func (r ContactResource) template(data acceptance.TestData) string {
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

resource "azurerm_orbital_contact_profile" "test" {
  name                              = "testcontactprofile-%[1]d"
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
        ip_address     = "10.0.1.0"
        port           = "49153"
        protocol       = "TCP"
      }
    }
    direction    = "Uplink"
    name         = "RHCP_UL"
    polarization = "RHCP"
  }
  network_configuration_subnet_id = azurerm_subnet.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func skipContact(t *testing.T) {
	if os.Getenv("ARM_TEST_SPACECRAFT_ID") == "" {
		t.Skip("Skipping as `ARM_TEST_SPACECRAFT_ID` was not specified")
	}
}
