package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CapacityReservationGroupResource struct{}

func TestAccCapacityReservationGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_capacity_reservation_group", "test")
	r := CapacityReservationGroupResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCapacityReservationGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_capacity_reservation_group", "test")
	r := CapacityReservationGroupResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccCapacityReservationGroup_zones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_capacity_reservation_group", "test")
	r := CapacityReservationGroupResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.zones(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCapacityReservationGroup_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_capacity_reservation_group", "test")
	r := CapacityReservationGroupResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.tags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.tagsUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CapacityReservationGroupResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.CapacityReservationGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Compute.CapacityReservationGroupsClient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r CapacityReservationGroupResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-compute-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CapacityReservationGroupResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_capacity_reservation_group" "test" {
  name                = "acctest-crg-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func (r CapacityReservationGroupResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_capacity_reservation_group" "import" {
  name                = azurerm_capacity_reservation_group.test.name
  resource_group_name = azurerm_capacity_reservation_group.test.resource_group_name
  location            = azurerm_capacity_reservation_group.test.location
}
`, config)
}

func (r CapacityReservationGroupResource) zones(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_capacity_reservation_group" "test" {
  name                = "acctest-ccrg-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  zones               = ["1", "2"]
}
`, template, data.RandomInteger)
}

func (r CapacityReservationGroupResource) tags(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_capacity_reservation_group" "test" {
  name                = "acctest-ccrg-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func (r CapacityReservationGroupResource) tagsUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_capacity_reservation_group" "test" {
  name                = "acctest-ccrg-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  tags = {
    ENV2 = "Test2"
  }
}
`, template, data.RandomInteger)
}
