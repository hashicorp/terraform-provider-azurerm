// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RouteMapResource struct{}

func TestAccRouteMap_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_map", "test")
	r := RouteMapResource{}
	nameSuffix := randString()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, nameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRouteMap_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_map", "test")
	r := RouteMapResource{}
	nameSuffix := randString()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, nameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, nameSuffix),
			ExpectError: acceptance.RequiresImportError("azurerm_route_map"),
		},
	})
}

func TestAccRouteMap_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_map", "test")
	r := RouteMapResource{}
	nameSuffix := randString()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, nameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRouteMap_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_map", "test")
	r := RouteMapResource{}
	nameSuffix := randString()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, nameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, nameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, nameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r RouteMapResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.RouteMapID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.RouteMapsClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func randString() string {
	charSet := "abcdefghijklmnopqrstuvwxyz"
	strlen := 5
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(result)
}

func (r RouteMapResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-rm-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctestvhub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r RouteMapResource) basic(data acceptance.TestData, nameSuffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_route_map" "test" {
  name           = "acctestrm-%s"
  virtual_hub_id = azurerm_virtual_hub.test.id
}
`, r.template(data), nameSuffix)
}

func (r RouteMapResource) requiresImport(data acceptance.TestData, nameSuffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_route_map" "import" {
  name           = azurerm_route_map.test.name
  virtual_hub_id = azurerm_route_map.test.virtual_hub_id
}
`, r.basic(data, nameSuffix))
}

func (r RouteMapResource) complete(data acceptance.TestData, nameSuffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_route_map" "test" {
  name           = "acctestrm-%s"
  virtual_hub_id = azurerm_virtual_hub.test.id

  rule {
    name                 = "rule1"
    next_step_if_matched = "Continue"

    action {
      type = "Add"

      parameter {
        as_path = ["22334"]
      }
    }

    match_criterion {
      match_condition = "Contains"
      route_prefix    = ["10.0.0.0/8"]
    }
  }
}
`, r.template(data), nameSuffix)
}

func (r RouteMapResource) update(data acceptance.TestData, nameSuffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_route_map" "test" {
  name           = "acctestrm-%s"
  virtual_hub_id = azurerm_virtual_hub.test.id

  rule {
    name                 = "rule2"
    next_step_if_matched = "Terminate"

    action {
      type = "Replace"

      parameter {
        route_prefix = ["10.0.1.0/8"]
      }
    }

    match_criterion {
      match_condition = "NotContains"
      as_path         = ["223345"]
    }
  }
}
`, r.template(data), nameSuffix)
}
