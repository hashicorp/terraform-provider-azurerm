// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type GridInfrastructureVersionsDataSource struct{}

func TestAccGridInfrastructureVersionsDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_grid_infrastructure_versions", "test")
	r := GridInfrastructureVersionsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.#").IsNotEmpty(),
				check.That(data.ResourceName).Key("versions.0.id").Exists(),
				check.That(data.ResourceName).Key("versions.0.name").Exists(),
				check.That(data.ResourceName).Key("versions.0.version").Exists(),
			),
		},
	})
}

func TestAccGridInfrastructureVersionsDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_grid_infrastructure_versions", "test")
	r := GridInfrastructureVersionsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.#").IsNotEmpty(),
				check.That(data.ResourceName).Key("versions.0.id").Exists(),
				check.That(data.ResourceName).Key("versions.0.name").Exists(),
				check.That(data.ResourceName).Key("versions.0.version").Exists(),
			),
		},
	})
}

func (d GridInfrastructureVersionsDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

data "azurerm_oracle_grid_infrastructure_versions" "test" {
  location = local.location
}
`, d.template(data))
}

func (d GridInfrastructureVersionsDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

data "azurerm_oracle_grid_infrastructure_versions" "test" {
  location = local.location
  shape    = "Exadata.X9M"
  zone     = local.zone
}
`, d.template(data))
}

func (a GridInfrastructureVersionsDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
locals {
  zone       = "1"
  location   = "%[1]s"
}

`, data.Locations.Primary)
}
