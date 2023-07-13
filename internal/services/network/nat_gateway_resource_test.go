// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NatGatewayResource struct{}

func TestAccNatGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway", "test")
	r := NatGatewayResource{}

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

func TestAccNatGateway_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway", "test")
	r := NatGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku_name").HasValue("Standard"),
				check.That(data.ResourceName).Key("idle_timeout_in_minutes").HasValue("10"),
				check.That(data.ResourceName).Key("zones.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNatGateway_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway", "test")
	r := NatGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku_name").HasValue("Standard"),
				check.That(data.ResourceName).Key("idle_timeout_in_minutes").HasValue("10"),
				check.That(data.ResourceName).Key("zones.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func (t NatGatewayResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NatGatewayID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.NatGatewayClient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading NAT Gateway (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

// Using alt location because the resource currently in private preview and is only available in eastus2.
func (NatGatewayResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_nat_gateway" "test" {
  name                = "acctestnatGateway-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger)
}

// Using alt location because the resource currently in private preview and is only available in eastus2.
func (NatGatewayResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicIP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicIPPrefix-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  prefix_length       = 30
  zones               = ["1"]
}

resource "azurerm_nat_gateway" "test" {
  name                    = "acctestnatGateway-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  sku_name                = "Standard"
  idle_timeout_in_minutes = 10
  zones                   = ["1"]
}

resource "azurerm_nat_gateway_public_ip_association" "test" {
  nat_gateway_id       = azurerm_nat_gateway.test.id
  public_ip_address_id = azurerm_public_ip.test.id
}

resource "azurerm_nat_gateway_public_ip_prefix_association" "test" {
  nat_gateway_id      = azurerm_nat_gateway.test.id
  public_ip_prefix_id = azurerm_public_ip_prefix.test.id
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
