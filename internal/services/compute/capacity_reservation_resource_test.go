// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CapacityReservationResource struct{}

func TestAccCapacityReservation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_capacity_reservation", "test")
	r := CapacityReservationResource{}
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

func TestAccCapacityReservation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_capacity_reservation", "test")
	r := CapacityReservationResource{}
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

func TestAccCapacityReservation_zone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_capacity_reservation", "test")
	r := CapacityReservationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCapacityReservation_skuCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_capacity_reservation", "test")
	r := CapacityReservationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.skuCapacity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.skuCapacityUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCapacityReservation_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_capacity_reservation", "test")
	r := CapacityReservationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.tagsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CapacityReservationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := capacityreservations.ParseCapacityReservationID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Compute.CapacityReservationsClient.Get(ctx, *id, capacityreservations.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r CapacityReservationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-compute-%d"
  location = "%s"
}

resource "azurerm_capacity_reservation_group" "test" {
  name                = "acctest-ccrg-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CapacityReservationResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_capacity_reservation" "test" {
  name                          = "acctest-ccr-%d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id
  sku {
    name     = "Standard_F2"
    capacity = 2
  }
}
`, template, data.RandomInteger)
}

func (r CapacityReservationResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_capacity_reservation" "import" {
  name                          = azurerm_capacity_reservation.test.name
  capacity_reservation_group_id = azurerm_capacity_reservation.test.capacity_reservation_group_id
  sku {
    name     = "Standard_F2"
    capacity = 2
  }
}
`, config)
}

func (r CapacityReservationResource) zone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-compute-%d"
  location = "%s"
}

resource "azurerm_capacity_reservation_group" "test" {
  name                = "acctest-ccrg-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  zones               = ["1", "2"]
}

resource "azurerm_capacity_reservation" "test" {
  name                          = "acctest-ccr-%d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id
  zone                          = "2"
  sku {
    name     = "Standard_F2"
    capacity = 2
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CapacityReservationResource) skuCapacity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_capacity_reservation" "test" {
  name                          = "acctest-ccr-%d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id
  sku {
    name     = "Standard_F2"
    capacity = 0
  }
}
`, template, data.RandomInteger)
}

func (r CapacityReservationResource) skuCapacityUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_capacity_reservation" "test" {
  name                          = "acctest-ccr-%d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id
  sku {
    name     = "Standard_F2"
    capacity = 1
  }
}
`, template, data.RandomInteger)
}

func (r CapacityReservationResource) tags(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_capacity_reservation" "test" {
  name                          = "acctest-ccr-%d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id
  sku {
    name     = "Standard_F2"
    capacity = 2
  }
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func (r CapacityReservationResource) tagsUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_capacity_reservation" "test" {
  name                          = "acctest-ccr-%d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id
  sku {
    name     = "Standard_F2"
    capacity = 2
  }
  tags = {
    ENV2 = "Test2"
  }
}
`, template, data.RandomInteger)
}
