// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package maps_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maps/2023-06-01/creators"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MapsCreatorResource struct{}

func TestAccMapsCreator_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maps_creator", "test")
	r := MapsCreatorResource{}
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

func TestAccMapsCreator_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maps_creator", "test")
	r := MapsCreatorResource{}
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

func TestAccMapsCreator_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maps_creator", "test")
	r := MapsCreatorResource{}
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

func TestAccMapsCreator_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maps_creator", "test")
	r := MapsCreatorResource{}
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

func (r MapsCreatorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := creators.ParseCreatorID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Maps.CreatorsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r MapsCreatorResource) template(data acceptance.TestData) string {
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
  location            = azurerm_resource_group.test.location
  sku_name            = "G2"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r MapsCreatorResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_maps_creator" "test" {
  name            = "accMapsCreator-%d"
  maps_account_id = azurerm_maps_account.test.id
  location        = "%s"
  storage_units   = 1
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MapsCreatorResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_maps_creator" "import" {
  name            = azurerm_maps_creator.test.name
  maps_account_id = azurerm_maps_account.test.id
  location        = "%s"
  storage_units   = 1
}
`, config, data.Locations.Primary)
}

func (r MapsCreatorResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_maps_creator" "test" {
  name            = "accMapsCreator-%d"
  maps_account_id = azurerm_maps_account.test.id
  location        = "%s"
  storage_units   = 1

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MapsCreatorResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_maps_creator" "test" {
  name            = "accMapsCreator-%d"
  maps_account_id = azurerm_maps_account.test.id
  location        = "%s"
  storage_units   = 2

  tags = {
    ENV  = "Test",
    ENV2 = "Test2"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}
