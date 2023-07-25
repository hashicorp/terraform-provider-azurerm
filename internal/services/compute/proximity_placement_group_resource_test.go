// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ProximityPlacementGroupResource struct{}

func TestAccProximityPlacementGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_proximity_placement_group", "test")
	r := ProximityPlacementGroupResource{}

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

func TestAccProximityPlacementGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_proximity_placement_group", "test")
	r := ProximityPlacementGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_proximity_placement_group"),
		},
	})
}

func TestAccProximityPlacementGroup_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_proximity_placement_group", "test")
	r := ProximityPlacementGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccProximityPlacementGroup_allowedVmSizes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_proximity_placement_group", "test")
	r := ProximityPlacementGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithAllowedVmSizes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithAllowedVmSizesUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithAllowedVmSizes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccProximityPlacementGroup_zone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_proximity_placement_group", "test")
	r := ProximityPlacementGroupResource{}

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

func TestAccProximityPlacementGroup_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_proximity_placement_group", "test")
	r := ProximityPlacementGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
				check.That(data.ResourceName).Key("tags.cost_center").HasValue("MSFT"),
			),
		},
		{
			Config: r.withUpdatedTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("staging"),
			),
		},
		data.ImportStep(),
	})
}

func (t ProximityPlacementGroupResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := proximityplacementgroups.ParseProximityPlacementGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.ProximityPlacementGroupsClient.Get(ctx, *id, proximityplacementgroups.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (ProximityPlacementGroupResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := proximityplacementgroups.ParseProximityPlacementGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	if _, err := client.Compute.ProximityPlacementGroupsClient.Delete(ctx, *id); err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (ProximityPlacementGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctestPPG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ProximityPlacementGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_proximity_placement_group" "import" {
  name                = azurerm_proximity_placement_group.test.name
  location            = azurerm_proximity_placement_group.test.location
  resource_group_name = azurerm_proximity_placement_group.test.resource_group_name
}
`, r.basic(data))
}

func (ProximityPlacementGroupResource) basicWithAllowedVmSizes(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctestPPG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  allowed_vm_sizes = ["Standard_F1"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ProximityPlacementGroupResource) basicWithAllowedVmSizesUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctestPPG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  allowed_vm_sizes = ["Standard_F1", "Standard_F2"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ProximityPlacementGroupResource) zone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctestPPG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  allowed_vm_sizes = ["Standard_F2"]
  zone             = "1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ProximityPlacementGroupResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctestPPG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ProximityPlacementGroupResource) withUpdatedTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctestPPG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
