// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
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
			Config: r.basic(),
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
			Config: r.complete(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.#").IsNotEmpty(),
				check.That(data.ResourceName).Key("versions.0.id").Exists(),
				check.That(data.ResourceName).Key("versions.0.name").Exists(),
				check.That(data.ResourceName).Key("versions.0.version").Exists(),
			),
		},
	})
}

func (d GridInfrastructureVersionsDataSource) basic() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_oracle_grid_infrastructure_versions" "test" {
  location = "eastus"
}
`
}

func (d GridInfrastructureVersionsDataSource) complete() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_oracle_grid_infrastructure_versions" "test" {
  location = "eastus"
  shape    = "Exadata.X9M"
  zone     = "2"
}
`
}
