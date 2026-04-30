// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-10-15/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CosmosDbFleetspaceResource struct{}

func TestAccCosmosDbFleetspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_fleetspace", "test")
	r := CosmosDbFleetspaceResource{}

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

func TestAccCosmosDbFleetspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_fleetspace", "test")
	r := CosmosDbFleetspaceResource{}

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

func TestAccCosmosDbFleetspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_fleetspace", "test")
	r := CosmosDbFleetspaceResource{}

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

func TestAccCosmosDbFleetspace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_fleetspace", "test")
	r := CosmosDbFleetspaceResource{}

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

func (CosmosDbFleetspaceResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fleets.ParseFleetspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Cosmos.FleetsClient.FleetspaceGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r CosmosDbFleetspaceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctest-cosmos-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CosmosDbFleetspaceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
resource "azurerm_cosmosdb_fleet" "test" {
  name                = "acctest-cosfleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_cosmosdb_fleetspace" "test" {
  name                = "acctest-cosfleetspace-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  fleet_name          = azurerm_cosmosdb_fleet.test.name
  service_tier        = "GeneralPurpose"
  data_regions = [
    azurerm_resource_group.test.location
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r CosmosDbFleetspaceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_cosmosdb_fleetspace" "import" {
  name                = azurerm_cosmosdb_fleetspace.test.name
  resource_group_name = azurerm_cosmosdb_fleetspace.test.resource_group_name
  fleet_name          = azurerm_cosmosdb_fleetspace.test.fleet_name
  service_tier        = azurerm_cosmosdb_fleetspace.test.service_tier
  data_regions        = azurerm_cosmosdb_fleetspace.test.data_regions
}
`, r.basic(data))
}

func (r CosmosDbFleetspaceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
resource "azurerm_cosmosdb_fleet" "test" {
  name                = "acctest-cosfleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_cosmosdb_fleetspace" "test" {
  name                = "acctest-cosfleetspace-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  fleet_name          = azurerm_cosmosdb_fleet.test.name
  service_tier        = "BusinessCritical"
  minimum_throughput  = 100000
  maximum_throughput  = 110000
  data_regions = [
    azurerm_resource_group.test.location,
    "%[3]s",
    "%[4]s"
  ]
}
`, r.template(data), data.RandomInteger, data.Locations.Secondary, data.Locations.Ternary)
}

func (r CosmosDbFleetspaceResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
resource "azurerm_cosmosdb_fleet" "test" {
  name                = "acctest-cosfleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_cosmosdb_fleetspace" "test" {
  name                = "acctest-cosfleetspace-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  fleet_name          = azurerm_cosmosdb_fleet.test.name
  service_tier        = "BusinessCritical"
  minimum_throughput  = 110000
  maximum_throughput  = 120000
  data_regions = [
    azurerm_resource_group.test.location,
    "%[3]s",
    "%[4]s"
  ]
}
`, r.template(data), data.RandomInteger, data.Locations.Secondary, data.Locations.Ternary)
}
