// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type GridInfrastructureMinorVersionsDataSource struct{}

func TestAccGridInfrastructureMinorVersionsDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_grid_infrastructure_minor_versions", "test")
	r := GridInfrastructureMinorVersionsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.#").IsNotEmpty(),
				check.That(data.ResourceName).Key("versions.0.id").Exists(),
				check.That(data.ResourceName).Key("versions.0.name").Exists(),
				check.That(data.ResourceName).Key("versions.0.version").Exists(),
				check.That(data.ResourceName).Key("versions.0.grid_image_ocid").Exists(),
			),
		},
	})
}

func TestAccGridInfrastructureMinorVersionsDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_grid_infrastructure_minor_versions", "test")
	r := GridInfrastructureMinorVersionsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.#").IsNotEmpty(),
				check.That(data.ResourceName).Key("versions.0.id").Exists(),
				check.That(data.ResourceName).Key("versions.0.name").Exists(),
				check.That(data.ResourceName).Key("versions.0.version").Exists(),
				check.That(data.ResourceName).Key("versions.0.grid_image_ocid").Exists(),
			),
		},
	})
}

func (d GridInfrastructureMinorVersionsDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

data "azurerm_oracle_grid_infrastructure_versions" "test" {
  location = local.location
  shape    = "ExaDbXS"
  zone     = local.zone
}

data "azurerm_oracle_grid_infrastructure_minor_versions" "test" {
  location                    = local.location
  grid_infrastructure_version = one([for item in data.azurerm_oracle_grid_infrastructure_versions.test.versions : item.name if item.version == local.grid_infrastructure_version])
  shape_family                = "EXADB_XS"
  zone                        = local.zone
}
`, d.template(data))
}

func (d GridInfrastructureMinorVersionsDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

data "azurerm_oracle_grid_infrastructure_versions" "test" {
  location = local.location
  shape    = "ExaDbXS"
  zone     = local.zone
}

data "azurerm_oracle_grid_infrastructure_minor_versions" "test" {
  location                    = local.location
  grid_infrastructure_version = data.azurerm_oracle_grid_infrastructure_versions.test.versions[0].name
  shape_family                = "EXADB_XS"
  zone                        = local.zone
}
`, d.template(data))
}

func (a GridInfrastructureMinorVersionsDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
locals {
  zone       = "1"
  location   = "%[1]s"
  grid_infrastructure_version = "26.0.0.0"
}

`, data.Locations.Primary)
}
