// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package maps_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/maps/2021-02-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MapsAccountResource struct{}

func TestAccMapsAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maps_account", "test")
	r := MapsAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("x_ms_client_id").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("sku_name").HasValue("S0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMapsAccount_sku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maps_account", "test")
	r := MapsAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sku(data, "S1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("x_ms_client_id").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("sku_name").HasValue("S1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMapsAccount_skuG2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maps_account", "test")
	r := MapsAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sku(data, "G2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("x_ms_client_id").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("sku_name").HasValue("G2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMapsAccount_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maps_account", "test")
	r := MapsAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("testing"),
			),
		},
		data.ImportStep(),
	})
}

func (MapsAccountResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := accounts.ParseAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Maps.AccountsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (MapsAccountResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_maps_account" "test" {
  name                = "accMapsAccount-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "S0"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (MapsAccountResource) sku(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_maps_account" "test" {
  name                = "accMapsAccount-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, sku)
}

func (MapsAccountResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_maps_account" "test" {
  name                = "accMapsAccount-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "S0"

  tags = {
    environment = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
