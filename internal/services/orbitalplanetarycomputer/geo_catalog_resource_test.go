// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package orbitalplanetarycomputer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbitalplanetarycomputer/2026-04-15/geocatalogs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type GeoCatalogResource struct{}

func TestAccGeoCatalog_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_geo_catalog", "test")
	r := GeoCatalogResource{}

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

func TestAccGeoCatalog_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_geo_catalog", "test")
	r := GeoCatalogResource{}

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

func TestAccGeoCatalog_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_geo_catalog", "test")
	r := GeoCatalogResource{}

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

func TestAccGeoCatalog_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_geo_catalog", "test")
	r := GeoCatalogResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
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
	})
}

func (r GeoCatalogResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := geocatalogs.ParseGeoCatalogID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.OrbitalPlanetaryComputer.GeoCatalogsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r GeoCatalogResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_geo_catalog" "test" {
  name                = "acctest-geocatalog-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, r.template(data), data.RandomString)
}

func (r GeoCatalogResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_geo_catalog" "import" {
  name                = azurerm_geo_catalog.test.name
  resource_group_name = azurerm_geo_catalog.test.resource_group_name
  location            = azurerm_geo_catalog.test.location
}
`, r.basic(data))
}

func (r GeoCatalogResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_geo_catalog" "test" {
  name                = "acctest-geocatalog-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  tags = {
    Environment = "Production"
    Label       = "Test"
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r GeoCatalogResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-planetarycomputer-%d"
  location = "uksouth" # Hardcode first as TeamCity location override does not work for newly added service
}
`, data.RandomInteger)
}
