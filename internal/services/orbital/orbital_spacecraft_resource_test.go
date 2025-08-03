// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package orbital_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/spacecraft"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SpacecraftResource struct{}

func TestAccSpacecraft_basic(t *testing.T) {
	if features.FivePointOh() {
		t.Skipf("Skipping since `azurerm_orbital_spacecraft` is deprecated and will be removed in 5.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_orbital_spacecraft", "test")
	r := SpacecraftResource{}

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

func TestAccSpacecraft_update(t *testing.T) {
	if features.FivePointOh() {
		t.Skipf("Skipping since `azurerm_orbital_spacecraft` is deprecated and will be removed in 5.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_orbital_spacecraft", "test")
	r := SpacecraftResource{}

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

func TestAccSpacecraft_complete(t *testing.T) {
	if features.FivePointOh() {
		t.Skipf("Skipping since `azurerm_orbital_spacecraft` is deprecated and will be removed in 5.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_orbital_spacecraft", "test")
	r := SpacecraftResource{}

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

func (r SpacecraftResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := spacecraft.ParseSpacecraftID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Orbital.SpacecraftClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r SpacecraftResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_orbital_spacecraft" "test" {
  name                = "acctestspacecraft-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "westus"
  norad_id            = "12345"

  links {
    bandwidth_mhz        = 30
    center_frequency_mhz = 2050
    direction            = "Uplink"
    polarization         = "LHCP"
    name                 = "linkname"
  }

  two_line_elements = ["1 23455U 94089A   97320.90946019  .00000140  00000-0  10191-3 0  2621", "2 23455  99.0090 272.6745 0008546 223.1686 136.8816 14.11711747148495"]
  title_line        = "AQUA"
}
`, template, data.RandomInteger)
}

func (r SpacecraftResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_orbital_spacecraft" "test" {
  name                = "acctestspacecraft-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "westus"
  norad_id            = "23456"

  links {
    bandwidth_mhz        = 20
    center_frequency_mhz = 2045
    direction            = "Uplink"
    polarization         = "LHCP"
    name                 = "linkname"
  }

  two_line_elements = ["1 23455U 94089A   97320.90946019  .00000140  00000-0  10191-3 0  2621", "2 23455  99.0090 272.6745 0008546 223.1686 136.8816 14.11711747148495"]
  title_line        = "AQUB"
}
`, template, data.RandomInteger)
}

func (r SpacecraftResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_orbital_spacecraft" "test" {
  name                = "acctestspacecraft-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "eastus"
  norad_id            = "12345"

  links {
    bandwidth_mhz        = 30
    center_frequency_mhz = 2050
    direction            = "Uplink"
    polarization         = "LHCP"
    name                 = "linkname"
  }

  two_line_elements = ["1 23455U 94089A   97320.90946019  .00000140  00000-0  10191-3 0  2621", "2 23455  99.0090 272.6745 0008546 223.1686 136.8816 14.11711747148495"]
  title_line        = "AQUA"

  tags = {
    aks-managed-cluster-name = "9a57225d-a405-4d40-aa46-f13d2342abef"
  }
}
`, template, data.RandomInteger)
}

func (r SpacecraftResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}
